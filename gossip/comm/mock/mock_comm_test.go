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
	"testing"

	"github.com/inklabsfoundation/inkchain/gossip/comm"
	"github.com/inklabsfoundation/inkchain/gossip/common"
	proto "github.com/inklabsfoundation/inkchain/protos/gossip"
	"github.com/stretchr/testify/assert"
)

func TestMockComm(t *testing.T) {
	first := &socketMock{"first", make(chan interface{})}
	second := &socketMock{"second", make(chan interface{})}
	members := make(map[string]*socketMock)

	members[first.endpoint] = first
	members[second.endpoint] = second

	comm1 := NewCommMock(first.endpoint, members)
	defer comm1.Stop()

	msgCh := comm1.Accept(func(message interface{}) bool {
		return message.(proto.ReceivedMessage).GetGossipMessage().GetStateRequest() != nil ||
			message.(proto.ReceivedMessage).GetGossipMessage().GetStateResponse() != nil
	})

	comm2 := NewCommMock(second.endpoint, members)
	defer comm2.Stop()

	sMsg, _ := (&proto.GossipMessage{
		Content: &proto.GossipMessage_StateRequest{&proto.RemoteStateRequest{
			StartSeqNum: 1,
			EndSeqNum:   3,
		}},
	}).NoopSign()
	comm2.Send(sMsg, &comm.RemotePeer{"first", common.PKIidType("first")})

	msg := <-msgCh

	assert.NotNil(t, msg.GetGossipMessage().GetStateRequest())
	assert.Equal(t, "first", string(comm1.GetPKIid()))
}

func TestMockComm_PingPong(t *testing.T) {
	members := make(map[string]*socketMock)

	members["peerA"] = &socketMock{"peerA", make(chan interface{})}
	members["peerB"] = &socketMock{"peerB", make(chan interface{})}

	peerA := NewCommMock("peerA", members)
	peerB := NewCommMock("peerB", members)

	all := func(interface{}) bool {
		return true
	}

	rcvChA := peerA.Accept(all)
	rcvChB := peerB.Accept(all)

	sMsg, _ := (&proto.GossipMessage{
		Content: &proto.GossipMessage_DataMsg{
			&proto.DataMessage{
				&proto.Payload{
					SeqNum: 1,
					Data:   []byte("Ping"),
				},
			}},
	}).NoopSign()
	peerA.Send(sMsg, &comm.RemotePeer{"peerB", common.PKIidType("peerB")})

	msg := <-rcvChB
	dataMsg := msg.GetGossipMessage().GetDataMsg()
	data := string(dataMsg.Payload.Data)
	assert.Equal(t, "Ping", data)

	msg.Respond(&proto.GossipMessage{
		Content: &proto.GossipMessage_DataMsg{
			&proto.DataMessage{
				&proto.Payload{
					SeqNum: 1,
					Data:   []byte("Pong"),
				},
			}},
	})

	msg = <-rcvChA
	dataMsg = msg.GetGossipMessage().GetDataMsg()
	data = string(dataMsg.Payload.Data)
	assert.Equal(t, "Pong", data)

}
