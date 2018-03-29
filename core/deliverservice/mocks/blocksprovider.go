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
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	gossip_common "github.com/inklabsfoundation/inkchain/gossip/common"
	"github.com/inklabsfoundation/inkchain/gossip/discovery"
	"github.com/inklabsfoundation/inkchain/protos/common"
	gossip_proto "github.com/inklabsfoundation/inkchain/protos/gossip"
	"github.com/inklabsfoundation/inkchain/protos/orderer"
	"github.com/inklabsfoundation/inkchain/protos/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// MockGossipServiceAdapter mocking structure for gossip service, used to initialize
// the blocks providers implementation and asserts the number
// of function calls used.
type MockGossipServiceAdapter struct {
	AddPayloadsCnt int32

	GossipBlockDisseminations chan uint64
}

type MockAtomicBroadcastClient struct {
	BD *MockBlocksDeliverer
}

func (mabc *MockAtomicBroadcastClient) Broadcast(ctx context.Context, opts ...grpc.CallOption) (orderer.AtomicBroadcast_BroadcastClient, error) {
	panic("Should not be used")
}
func (mabc *MockAtomicBroadcastClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (orderer.AtomicBroadcast_DeliverClient, error) {
	return mabc.BD, nil
}

// PeersOfChannel returns the slice with peers participating in given channel
func (*MockGossipServiceAdapter) PeersOfChannel(gossip_common.ChainID) []discovery.NetworkMember {
	return []discovery.NetworkMember{}
}

// AddPayload adds gossip payload to the local state transfer buffer
func (mock *MockGossipServiceAdapter) AddPayload(chainID string, payload *gossip_proto.Payload) error {
	atomic.AddInt32(&mock.AddPayloadsCnt, 1)
	return nil
}

// Gossip message to the all peers
func (mock *MockGossipServiceAdapter) Gossip(msg *gossip_proto.GossipMessage) {
	mock.GossipBlockDisseminations <- msg.GetDataMsg().Payload.SeqNum
}

// MockBlocksDeliverer mocking structure of BlocksDeliverer interface to initialize
// the blocks provider implementation
type MockBlocksDeliverer struct {
	DisconnectCalled chan struct{}
	CloseCalled      chan struct{}
	Pos              uint64
	grpc.ClientStream
	RecvCnt  int32
	MockRecv func(mock *MockBlocksDeliverer) (*orderer.DeliverResponse, error)
}

// Recv gets responses from the ordering service, currently mocked to return
// only one response with empty block.
func (mock *MockBlocksDeliverer) Recv() (*orderer.DeliverResponse, error) {
	atomic.AddInt32(&mock.RecvCnt, 1)
	return mock.MockRecv(mock)
}

// MockRecv mock for the Recv function
func MockRecv(mock *MockBlocksDeliverer) (*orderer.DeliverResponse, error) {
	pos := mock.Pos

	// Advance position for the next call
	mock.Pos++
	return &orderer.DeliverResponse{
		Type: &orderer.DeliverResponse_Block{
			Block: &common.Block{
				Header: &common.BlockHeader{
					Number:       pos,
					DataHash:     []byte{},
					PreviousHash: []byte{},
				},
				Data: &common.BlockData{
					Data: [][]byte{},
				},
			}},
	}, nil
}

// Send sends the envelope with request for the blocks for ordering service
// currently mocked and not doing anything
func (mock *MockBlocksDeliverer) Send(env *common.Envelope) error {
	payload, _ := utils.GetPayload(env)
	seekInfo := &orderer.SeekInfo{}

	proto.Unmarshal(payload.Data, seekInfo)

	// Read starting position
	switch t := seekInfo.Start.Type.(type) {
	case *orderer.SeekPosition_Oldest:
		mock.Pos = 0
	case *orderer.SeekPosition_Specified:
		mock.Pos = t.Specified.Number
	}
	return nil
}

func (mock *MockBlocksDeliverer) Disconnect() {
	mock.DisconnectCalled <- struct{}{}
}

func (mock *MockBlocksDeliverer) Close() {
	if mock.CloseCalled == nil {
		return
	}
	mock.CloseCalled <- struct{}{}
}

// MockLedgerInfo mocking implementation of LedgerInfo interface, needed
// for test initialization purposes
type MockLedgerInfo struct {
	Height uint64
}

// LedgerHeight returns mocked value to the ledger height
func (li *MockLedgerInfo) LedgerHeight() (uint64, error) {
	return atomic.LoadUint64(&li.Height), nil
}
