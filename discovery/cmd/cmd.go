/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"os"
	"time"

	"github.com/inklabsfoundation/inkchain/cmd/common"
	"github.com/inklabsfoundation/inkchain/discovery/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	PeersCommand     = "peers"
	ConfigCommand    = "config"
	EndorsersCommand = "endorsers"
)

var (
	// responseParserWriter defines the stdout
	responseParserWriter = os.Stdout
)

const (
	defaultTimeout = time.Second * 10
)

//go:generate mockery -dir . -name Stub -case underscore -output mocks/

// Stub represents the remote discovery service
type Stub interface {
	// Send sends the request, and receives a response
	Send(server string, conf common.Config, req *discovery.Request) (ServiceResponse, error)
}

//go:generate mockery -dir . -name ResponseParser -case underscore -output mocks/

// ResponseParser parses responses sent from the server
type ResponseParser interface {
	// ParseResponse parses the response and uses the given output when emitting data
	ParseResponse(channel string, response ServiceResponse) error
}

//go:generate mockery -dir . -name CommandRegistrar -case underscore -output mocks/

// CommandRegistrar registers commands
type CommandRegistrar interface {
	// Command adds a new top-level command to the CLI
	Command(name, help string, onCommand common.CLICommand) *kingpin.CmdClause
}

// AddCommands registers the discovery commands to the given CommandRegistrar
func AddCommands(cli CommandRegistrar) {
	peerCmd := NewPeerCmd(&ClientStub{}, &PeerResponseParser{Writer: responseParserWriter})
	peers := cli.Command(PeersCommand, "Discover peers", peerCmd.Execute)
	server := peers.Flag("server", "Sets the endpoint of the server to connect").String()
	channel := peers.Flag("channel", "Sets the channel the query is intended to").String()
	peerCmd.SetServer(server)
	peerCmd.SetChannel(channel)

	configCmd := NewConfigCmd(&ClientStub{}, &ConfigResponseParser{Writer: responseParserWriter})
	config := cli.Command(ConfigCommand, "Discover channel config", configCmd.Execute)

	configCmd.SetChannel(channel)

	endorserCmd := NewEndorsersCmd(&RawStub{}, &EndorserResponseParser{Writer: responseParserWriter})
	endorsers := cli.Command(EndorsersCommand, "Discover chaincode endorsers", endorserCmd.Execute)

	endorserCmd.SetServer(server)
	endorserCmd.SetChaincodes(chaincodes)
	endorserCmd.SetCollections(collections)
}
