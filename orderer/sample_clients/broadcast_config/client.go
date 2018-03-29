/*
Copyright Ziggurat Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"

	"google.golang.org/grpc"

	genesisconfig "github.com/inklabsfoundation/inkchain/common/configtx/tool/localconfig"
	"github.com/inklabsfoundation/inkchain/msp"
	mspmgmt "github.com/inklabsfoundation/inkchain/msp/mgmt"
	"github.com/inklabsfoundation/inkchain/orderer/localconfig"
	cb "github.com/inklabsfoundation/inkchain/protos/common"
	ab "github.com/inklabsfoundation/inkchain/protos/orderer"
)

var conf *config.TopLevel
var genConf *genesisconfig.Profile
var signer msp.SigningIdentity

type broadcastClient struct {
	ab.AtomicBroadcast_BroadcastClient
}

func (bc *broadcastClient) broadcast(env *cb.Envelope) error {
	var err error
	var resp *ab.BroadcastResponse

	err = bc.Send(env)
	if err != nil {
		return err
	}

	resp, err = bc.Recv()
	if err != nil {
		return err
	}

	fmt.Println("Status:", resp)
	return nil
}

// cmdImpl holds the command and its arguments.
type cmdImpl struct {
	name string
	args argsImpl
}

// argsImpl holds all the possible arguments for all possible commands.
type argsImpl struct {
	consensusType  string
	creationPolicy string
	chainID        string
}

func init() {
	conf = config.Load()
	genConf = genesisconfig.Load(genesisconfig.SampleInsecureProfile)

	// Load local MSP
	err := mspmgmt.LoadLocalMsp(conf.General.LocalMSPDir, conf.General.BCCSP, conf.General.LocalMSPID)
	if err != nil {
		panic(fmt.Errorf("Failed to initialize local MSP: %s", err))
	}

	msp := mspmgmt.GetLocalMSP()
	signer, err = msp.GetDefaultSigningIdentity()
	if err != nil {
		panic(fmt.Errorf("Failed to initialize get default signer: %s", err))
	}
}

func main() {
	cmd := new(cmdImpl)
	var srv string

	flag.StringVar(&srv, "server", fmt.Sprintf("%s:%d", conf.General.ListenAddress, conf.General.ListenPort), "The RPC server to connect to.")
	flag.StringVar(&cmd.name, "cmd", "newChain", "The action that this client is requesting via the config transaction.")
	flag.StringVar(&cmd.args.consensusType, "consensusType", genConf.Orderer.OrdererType, "In case of a newChain command, the type of consensus the ordering service is running on.")
	flag.StringVar(&cmd.args.creationPolicy, "creationPolicy", "AcceptAllPolicy", "In case of a newChain command, the chain creation policy this request should be validated against.")
	flag.StringVar(&cmd.args.chainID, "chainID", "NewChannelId", "In case of a newChain command, the chain ID to create.")
	flag.Parse()

	conn, err := grpc.Dial(srv, grpc.WithInsecure())
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	client, err := ab.NewAtomicBroadcastClient(conn).Broadcast(context.TODO())
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	bc := &broadcastClient{client}

	switch cmd.name {
	case "newChain":
		env := newChainRequest(cmd.args.consensusType, cmd.args.creationPolicy, cmd.args.chainID)
		fmt.Println("Requesting the creation of chain", cmd.args.chainID)
		fmt.Println(bc.broadcast(env))
	default:
		panic("Invalid command given")
	}
}
