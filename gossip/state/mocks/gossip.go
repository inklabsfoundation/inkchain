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

package mocks

import (
	"github.com/inklabsfoundation/inkchain/gossip/api"
	"github.com/inklabsfoundation/inkchain/gossip/comm"
	"github.com/inklabsfoundation/inkchain/gossip/common"
	"github.com/inklabsfoundation/inkchain/gossip/discovery"
	proto "github.com/inklabsfoundation/inkchain/protos/gossip"
	"github.com/stretchr/testify/mock"
)

type GossipMock struct {
	mock.Mock
}

func (*GossipMock) SuspectPeers(s api.PeerSuspector) {
	panic("implement me")
}

func (*GossipMock) Send(msg *proto.GossipMessage, peers ...*comm.RemotePeer) {
	panic("implement me")
}

func (*GossipMock) Peers() []discovery.NetworkMember {
	panic("implement me")
}

func (*GossipMock) PeersOfChannel(common.ChainID) []discovery.NetworkMember {
	return nil
}

func (*GossipMock) UpdateMetadata(metadata []byte) {
	panic("implement me")
}

func (*GossipMock) UpdateChannelMetadata(metadata []byte, chainID common.ChainID) {

}

func (*GossipMock) Gossip(msg *proto.GossipMessage) {
	panic("implement me")
}

func (g *GossipMock) Accept(acceptor common.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan proto.ReceivedMessage) {
	args := g.Called(acceptor, passThrough)
	if args.Get(0) == nil {
		return nil, args.Get(1).(<-chan proto.ReceivedMessage)
	}
	return args.Get(0).(<-chan *proto.GossipMessage), nil
}

func (g *GossipMock) JoinChan(joinMsg api.JoinChannelMessage, chainID common.ChainID) {
}

func (*GossipMock) Stop() {
}
