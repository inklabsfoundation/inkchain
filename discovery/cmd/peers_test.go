/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/inklabsfoundation/inkchain/cmd/common"
	. "github.com/inklabsfoundation/inkchain/discovery/client"
	discovery "github.com/inklabsfoundation/inkchain/discovery/cmd"
	"github.com/inklabsfoundation/inkchain/discovery/cmd/mocks"
	"github.com/inklabsfoundation/inkchain/gossip/protoext"
	"github.com/inklabsfoundation/inkchain/protos/gossip"
	"github.com/inklabsfoundation/inkchain/protos/msp"
	"github.com/inklabsfoundation/inkchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPeerCmd(t *testing.T) {
	server := "peer0"
	stub := &mocks.Stub{}
	parser := &mocks.ResponseParser{}
	cmd := discovery.NewPeerCmd(stub, parser)

	t.Run("no server supplied", func(t *testing.T) {
		cmd.SetServer(nil)
		err := cmd.Execute(common.Config{})
		assert.Equal(t, err.Error(), "no server specified")
	})

	t.Run("Server return error", func(t *testing.T) {
		cmd.SetServer(&server)
		stub.On("Send", server, mock.Anything, mock.Anything).Return(nil, errors.New("deadline exceeded")).Once()
		err := cmd.Execute(common.Config{})
		assert.Contains(t, err.Error(), "deadline exceeded")
	})

	t.Run("Channel(less) peer query", func(t *testing.T) {
		stub.On("Send", server, mock.Anything, mock.Anything).Return(nil, nil).Twice()
		cmd.SetServer(&server)
		cmd.SetChannel(nil)

		var emptyChannel string
		parser.On("ParseResponse", emptyChannel, mock.Anything).Return(nil)
		err := cmd.Execute(common.Config{})
		assert.NoError(t, err)

		channel := "mychannel"
		cmd.SetChannel(&channel)
		parser.On("ParseResponse", channel, mock.Anything).Return(nil)
		err = cmd.Execute(common.Config{})
		assert.NoError(t, err)
	})
}

func TestParsePeers(t *testing.T) {
	buff := &bytes.Buffer{}
	parser := &discovery.PeerResponseParser{Writer: buff}
	res := &mocks.ServiceResponse{}

	sID := &msp.SerializedIdentity{
		Mspid:   "Org1MSP",
		IdBytes: []byte("identity"),
	}

	idBytes := protoutil.MarshalOrPanic(sID)

	validPeer := &Peer{
		MSPID:            "Org1MSP",
		Identity:         idBytes,
		AliveMessage:     aliveMessage(0),
		StateInfoMessage: stateInfoMessage(100),
	}
	invalidPeer := &Peer{
		MSPID: "Org2MSP",
	}

	chanRes := &mocks.ChannelResponse{}
	chanRes.On("Peers").Return([]*Peer{validPeer, invalidPeer}, nil)
	locRes := &mocks.LocalResponse{}
	locRes.On("Peers").Return([]*Peer{validPeer, invalidPeer}, nil)

	res.On("ForChannel", "mychannel").Return(chanRes)
	res.On("ForLocal").Return(locRes)

	channel2expected := map[string]string{
		"mychannel": "[\n\t{\n\t\t\"MSPID\": \"Org1MSP\",\n\t\t\"LedgerHeight\": 100,\n\t\t\"Endpoint\": \"p0\",\n\t\t\"Identity\": \"identity\",\n\t\t\"Chaincodes\": [\n\t\t\t\"mycc\",\n\t\t\t\"mycc2\"\n\t\t]\n\t},\n\t{\n\t\t\"MSPID\": \"Org2MSP\",\n\t\t\"LedgerHeight\": 0,\n\t\t\"Endpoint\": \"\",\n\t\t\"Identity\": \"\",\n\t\t\"Chaincodes\": null\n\t}\n]",
		"":          "[\n\t{\n\t\t\"MSPID\": \"Org1MSP\",\n\t\t\"Endpoint\": \"p0\",\n\t\t\"Identity\": \"identity\"\n\t},\n\t{\n\t\t\"MSPID\": \"Org2MSP\",\n\t\t\"Endpoint\": \"\",\n\t\t\"Identity\": \"\"\n\t}\n]",
	}

	for channel, expected := range channel2expected {
		buff.Reset()
		err := parser.ParseResponse(channel, res)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%s\n", expected), buff.String())
	}
}

func aliveMessage(id int) *protoext.SignedGossipMessage {
	g := &gossip.GossipMessage{
		Content: &gossip.GossipMessage_AliveMsg{
			AliveMsg: &gossip.AliveMessage{
				Timestamp: &gossip.PeerTime{
					SeqNum: uint64(id),
					IncNum: uint64(time.Now().UnixNano()),
				},
				Membership: &gossip.Member{
					Endpoint: fmt.Sprintf("p%d", id),
				},
			},
		},
	}
	sMsg, _ := protoext.NoopSign(g)
	return sMsg
}

func stateInfoMessage(height uint64) *protoext.SignedGossipMessage {
	g := &gossip.GossipMessage{
		Content: &gossip.GossipMessage_StateInfo{
			StateInfo: &gossip.StateInfo{
				Timestamp: &gossip.PeerTime{
					SeqNum: 5,
					IncNum: uint64(time.Now().UnixNano()),
				},
				Properties: &gossip.Properties{
					LedgerHeight: height,
					Chaincodes: []*gossip.Chaincode{
						{Name: "mycc"},
						{Name: "mycc2"},
					},
				},
			},
		},
	}
	sMsg, _ := protoext.NoopSign(g)
	return sMsg
}
