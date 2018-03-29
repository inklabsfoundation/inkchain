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

package mock

import (
	"github.com/inklabsfoundation/inkchain/gossip/api"
	"github.com/inklabsfoundation/inkchain/gossip/comm"
	"github.com/inklabsfoundation/inkchain/gossip/common"
	"github.com/inklabsfoundation/inkchain/gossip/util"
	proto "github.com/inklabsfoundation/inkchain/protos/gossip"
)

// Mock which aims to simulate socket
type socketMock struct {
	// socket endpoint
	endpoint string

	// To simulate simple tcp socket
	socket chan interface{}
}

// Mock of primitive tcp packet structure
type packetMock struct {
	// Sender channel message sent from
	src *socketMock

	// Destination channel sent to
	dst *socketMock

	msg interface{}
}

type channelMock struct {
	accept common.MessageAcceptor

	channel chan proto.ReceivedMessage
}

type commMock struct {
	id string

	members map[string]*socketMock

	acceptors []*channelMock

	deadChannel chan common.PKIidType

	done chan struct{}
}

var logger = util.GetLogger(util.LoggingMockModule, "")

// NewCommMock creates mocked communication object
func NewCommMock(id string, members map[string]*socketMock) comm.Comm {
	res := &commMock{
		id: id,

		members: members,

		acceptors: make([]*channelMock, 0),

		done: make(chan struct{}),

		deadChannel: make(chan common.PKIidType),
	}
	// Start communication service
	go res.start()

	return res
}

// Respond sends a GossipMessage to the origin from which this ReceivedMessage was sent from
func (packet *packetMock) Respond(msg *proto.GossipMessage) {
	sMsg, _ := msg.NoopSign()
	packet.src.socket <- &packetMock{
		src: packet.dst,
		dst: packet.src,
		msg: sMsg,
	}
}

// GetSourceEnvelope Returns the Envelope the ReceivedMessage was
// constructed with
func (packet *packetMock) GetSourceEnvelope() *proto.Envelope {
	return nil
}

// GetGossipMessage returns the underlying GossipMessage
func (packet *packetMock) GetGossipMessage() *proto.SignedGossipMessage {
	return packet.msg.(*proto.SignedGossipMessage)
}

// GetConnectionInfo returns information about the remote peer
// that sent the message
func (packet *packetMock) GetConnectionInfo() *proto.ConnectionInfo {
	return nil
}

func (mock *commMock) start() {
	logger.Debug("Starting communication mock module...")
	for {
		select {
		case <-mock.done:
			{
				// Got final signal, exiting...
				logger.Debug("Exiting...")
				return
			}
		case msg := <-mock.members[mock.id].socket:
			{
				logger.Debug("Got new message", msg)
				packet := msg.(*packetMock)
				for _, channel := range mock.acceptors {
					// if message acceptor agrees to get
					// new message forward it to the received
					// messages channel
					if channel.accept(packet) {
						channel.channel <- packet
					}
				}
			}
		}
	}
}

// GetPKIid returns this instance's PKI id
func (mock *commMock) GetPKIid() common.PKIidType {
	return common.PKIidType(mock.id)
}

// Send sends a message to remote peers
func (mock *commMock) Send(msg *proto.SignedGossipMessage, peers ...*comm.RemotePeer) {
	for _, peer := range peers {
		logger.Debug("Sending message to peer ", peer.Endpoint, "from ", mock.id)
		mock.members[peer.Endpoint].socket <- &packetMock{
			src: mock.members[mock.id],
			dst: mock.members[peer.Endpoint],
			msg: msg,
		}
	}
}

// Probe probes a remote node and returns nil if its responsive,
// and an error if it's not.
func (mock *commMock) Probe(peer *comm.RemotePeer) error {
	return nil
}

// Handshake authenticates a remote peer and returns
// (its identity, nil) on success and (nil, error)
func (mock *commMock) Handshake(peer *comm.RemotePeer) (api.PeerIdentityType, error) {
	return nil, nil
}

// Accept returns a dedicated read-only channel for messages sent by other nodes that match a certain predicate.
// Each message from the channel can be used to send a reply back to the sender
func (mock *commMock) Accept(accept common.MessageAcceptor) <-chan proto.ReceivedMessage {
	ch := make(chan proto.ReceivedMessage)
	mock.acceptors = append(mock.acceptors, &channelMock{accept, ch})
	return ch
}

// PresumedDead returns a read-only channel for node endpoints that are suspected to be offline
func (mock *commMock) PresumedDead() <-chan common.PKIidType {
	return mock.deadChannel
}

// CloseConn closes a connection to a certain endpoint
func (mock *commMock) CloseConn(peer *comm.RemotePeer) {
	// NOOP
}

// Stop stops the module
func (mock *commMock) Stop() {
	logger.Debug("Stopping communication module, closing all accepting channels.")
	for _, accept := range mock.acceptors {
		close(accept.channel)
	}
	logger.Debug("[XXX]: Sending done signal to close the module.")
	mock.done <- struct{}{}
}
