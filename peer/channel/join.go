/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

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

package channel

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/inklabsfoundation/inkchain/core/scc/cscc"
	"github.com/inklabsfoundation/inkchain/peer/common"
	pcommon "github.com/inklabsfoundation/inkchain/protos/common"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	putils "github.com/inklabsfoundation/inkchain/protos/utils"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const commandDescription = "Joins the peer to a chain."

func joinCmd(cf *ChannelCmdFactory) *cobra.Command {
	// Set the flags on the channel start command.
	joinCmd := &cobra.Command{
		Use:   "join",
		Short: commandDescription,
		Long:  commandDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return join(cmd, args, cf)
		},
	}
	flagList := []string{
		"blockpath",
	}
	attachFlags(joinCmd, flagList)

	return joinCmd
}

//GBFileNotFoundErr genesis block file not found
type GBFileNotFoundErr string

func (e GBFileNotFoundErr) Error() string {
	return fmt.Sprintf("genesis block file not found %s", string(e))
}

//ProposalFailedErr proposal failed
type ProposalFailedErr string

func (e ProposalFailedErr) Error() string {
	return fmt.Sprintf("proposal failed (err: %s)", string(e))
}

func getJoinCCSpec() (*pb.ChaincodeSpec, error) {
	if genesisBlockPath == common.UndefinedParamValue {
		return nil, errors.New("Must supply genesis block file")
	}

	gb, err := ioutil.ReadFile(genesisBlockPath)
	if err != nil {
		return nil, GBFileNotFoundErr(err.Error())
	}
	// Build the spec
	input := &pb.ChaincodeInput{Args: [][]byte{[]byte(cscc.JoinChain), gb}}

	spec := &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]),
		ChaincodeId: &pb.ChaincodeID{Name: "cscc"},
		Input:       input,
	}

	return spec, nil
}

func executeJoin(cf *ChannelCmdFactory) (err error) {
	spec, err := getJoinCCSpec()
	if err != nil {
		return err
	}

	// Build the ChaincodeInvocationSpec message
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	creator, err := cf.Signer.Serialize()
	if err != nil {
		return fmt.Errorf("Error serializing identity for %s: %s", cf.Signer.GetIdentifier(), err)
	}

	var prop *pb.Proposal
	prop, _, err = putils.CreateProposalFromCIS(pcommon.HeaderType_CONFIG, "", invocation, creator)
	if err != nil {
		return fmt.Errorf("Error creating proposal for join %s", err)
	}

	var signedProp *pb.SignedProposal
	signedProp, err = putils.GetSignedProposal(prop, cf.Signer)
	if err != nil {
		return fmt.Errorf("Error creating signed proposal %s", err)
	}

	var proposalResp *pb.ProposalResponse
	proposalResp, err = cf.EndorserClient.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return ProposalFailedErr(err.Error())
	}

	if proposalResp == nil {
		return ProposalFailedErr("nil proposal response")
	}

	if proposalResp.Response.Status != 0 && proposalResp.Response.Status != 200 {
		return ProposalFailedErr(fmt.Sprintf("bad proposal response %d", proposalResp.Response.Status))
	}
	logger.Infof("Peer joined the channel!")
	return nil
}

func join(cmd *cobra.Command, args []string, cf *ChannelCmdFactory) error {
	if genesisBlockPath == common.UndefinedParamValue {
		return errors.New("Must supply genesis block path")
	}

	var err error
	if cf == nil {
		cf, err = InitCmdFactory(EndorserRequired, OrdererNotRequired)
		if err != nil {
			return err
		}
	}
	return executeJoin(cf)
}
