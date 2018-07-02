/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

package shim

import (
	"errors"
	"fmt"
	"sync"

	"math/big"

	"encoding/json"

	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/core/wallet"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet/kvtranset"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"github.com/looplab/fsm"
	"github.com/inklabsfoundation/inkchain/protos/ledger/crosstranset/kvcrosstranset"
)

// PeerChaincodeStream interface for stream between Peer and chaincode instance.
type PeerChaincodeStream interface {
	Send(*pb.ChaincodeMessage) error
	Recv() (*pb.ChaincodeMessage, error)
	CloseSend() error
}

type nextStateInfo struct {
	msg      *pb.ChaincodeMessage
	sendToCC bool
}

func (handler *Handler) triggerNextState(msg *pb.ChaincodeMessage, send bool) {
	handler.nextState <- &nextStateInfo{msg, send}
}

// Handler handler implementation for shim side of chaincode.
type Handler struct {
	sync.RWMutex
	//shim to peer grpc serializer. User only in serialSend
	serialLock sync.Mutex
	To         string
	ChatStream PeerChaincodeStream
	FSM        *fsm.FSM
	cc         Chaincode
	// Multiple queries (and one transaction) with different txids can be executing in parallel for this chaincode
	// responseChannel is the channel on which responses are communicated by the shim to the chaincodeStub.
	responseChannel map[string]chan pb.ChaincodeMessage
	nextState       chan *nextStateInfo
}

func shorttxid(txid string) string {
	if len(txid) < 8 {
		return txid
	}
	return txid[0:8]
}

//serialSend serializes msgs so gRPC will be happy
func (handler *Handler) serialSend(msg *pb.ChaincodeMessage) error {
	handler.serialLock.Lock()
	defer handler.serialLock.Unlock()

	err := handler.ChatStream.Send(msg)

	return err
}

//serialSendAsync serves the same purpose as serialSend (serializ msgs so gRPC will
//be happy). In addition, it is also asynchronous so send-remoterecv--localrecv loop
//can be nonblocking. Only errors need to be handled and these are handled by
//communication on supplied error channel. A typical use will be a non-blocking or
//nil channel
func (handler *Handler) serialSendAsync(msg *pb.ChaincodeMessage, errc chan error) {
	go func() {
		err := handler.serialSend(msg)
		if errc != nil {
			errc <- err
		}
	}()
}

func (handler *Handler) createChannel(txid string) (chan pb.ChaincodeMessage, error) {
	handler.Lock()
	defer handler.Unlock()
	if handler.responseChannel == nil {
		return nil, fmt.Errorf("[%s]Cannot create response channel", shorttxid(txid))
	}
	if handler.responseChannel[txid] != nil {
		return nil, fmt.Errorf("[%s]Channel exists", shorttxid(txid))
	}
	c := make(chan pb.ChaincodeMessage)
	handler.responseChannel[txid] = c
	return c, nil
}

func (handler *Handler) sendChannel(msg *pb.ChaincodeMessage) error {
	handler.Lock()
	defer handler.Unlock()
	if handler.responseChannel == nil {
		return fmt.Errorf("[%s]Cannot send message response channel", shorttxid(msg.Txid))
	}
	if handler.responseChannel[msg.Txid] == nil {
		return fmt.Errorf("[%s]sendChannel does not exist", shorttxid(msg.Txid))
	}

	chaincodeLogger.Debugf("[%s]before send", shorttxid(msg.Txid))
	handler.responseChannel[msg.Txid] <- *msg
	chaincodeLogger.Debugf("[%s]after send", shorttxid(msg.Txid))

	return nil
}

//sends a message and selects
func (handler *Handler) sendReceive(msg *pb.ChaincodeMessage, c chan pb.ChaincodeMessage) (pb.ChaincodeMessage, error) {
	errc := make(chan error, 1)
	handler.serialSendAsync(msg, errc)

	//the serialsend above will send an err or nil
	//the select filters that first error(or nil)
	//and continues to wait for the response
	//it is possible that the response triggers first
	//in which case the errc obviously worked and is
	//ignored
	for {
		select {
		case err := <-errc:
			if err == nil {
				continue
			}
			//would have been logged, return false
			return pb.ChaincodeMessage{}, err
		case outmsg, val := <-c:
			if !val {
				return pb.ChaincodeMessage{}, fmt.Errorf("unexpected failure on receive")
			}
			return outmsg, nil
		}
	}
}

func (handler *Handler) deleteChannel(txid string) {
	handler.Lock()
	defer handler.Unlock()
	if handler.responseChannel != nil {
		delete(handler.responseChannel, txid)
	}
}

// NewChaincodeHandler returns a new instance of the shim side handler.
func newChaincodeHandler(peerChatStream PeerChaincodeStream, chaincode Chaincode) *Handler {
	v := &Handler{
		ChatStream: peerChatStream,
		cc:         chaincode,
	}
	v.responseChannel = make(map[string]chan pb.ChaincodeMessage)
	v.nextState = make(chan *nextStateInfo)

	// Create the shim side FSM
	v.FSM = fsm.NewFSM(
		"created",
		fsm.Events{
			{Name: pb.ChaincodeMessage_REGISTERED.String(), Src: []string{"created"}, Dst: "established"},
			{Name: pb.ChaincodeMessage_READY.String(), Src: []string{"established"}, Dst: "ready"},
			{Name: pb.ChaincodeMessage_ERROR.String(), Src: []string{"init"}, Dst: "established"},
			{Name: pb.ChaincodeMessage_RESPONSE.String(), Src: []string{"init"}, Dst: "init"},
			{Name: pb.ChaincodeMessage_INIT.String(), Src: []string{"ready"}, Dst: "ready"},
			{Name: pb.ChaincodeMessage_TRANSACTION.String(), Src: []string{"ready"}, Dst: "ready"},
			{Name: pb.ChaincodeMessage_RESPONSE.String(), Src: []string{"ready"}, Dst: "ready"},
			{Name: pb.ChaincodeMessage_ERROR.String(), Src: []string{"ready"}, Dst: "ready"},
			{Name: pb.ChaincodeMessage_COMPLETED.String(), Src: []string{"init"}, Dst: "ready"},
			{Name: pb.ChaincodeMessage_COMPLETED.String(), Src: []string{"ready"}, Dst: "ready"},
		},
		fsm.Callbacks{
			"before_" + pb.ChaincodeMessage_REGISTERED.String():  func(e *fsm.Event) { v.beforeRegistered(e) },
			"after_" + pb.ChaincodeMessage_RESPONSE.String():     func(e *fsm.Event) { v.afterResponse(e) },
			"after_" + pb.ChaincodeMessage_ERROR.String():        func(e *fsm.Event) { v.afterError(e) },
			"before_" + pb.ChaincodeMessage_INIT.String():        func(e *fsm.Event) { v.beforeInit(e) },
			"before_" + pb.ChaincodeMessage_TRANSACTION.String(): func(e *fsm.Event) { v.beforeTransaction(e) },
		},
	)
	return v
}

// beforeRegistered is called to handle the REGISTERED message.
func (handler *Handler) beforeRegistered(e *fsm.Event) {
	if _, ok := e.Args[0].(*pb.ChaincodeMessage); !ok {
		e.Cancel(fmt.Errorf("Received unexpected message type"))
		return
	}
	chaincodeLogger.Debugf("Received %s, ready for invocations", pb.ChaincodeMessage_REGISTERED)
}

// handleInit handles request to initialize chaincode.
func (handler *Handler) handleInit(msg *pb.ChaincodeMessage) {
	// The defer followed by triggering a go routine dance is needed to ensure that the previous state transition
	// is completed before the next one is triggered. The previous state transition is deemed complete only when
	// the beforeInit function is exited. Interesting bug fix!!
	go func() {
		var nextStateMsg *pb.ChaincodeMessage

		send := true

		defer func() {
			handler.triggerNextState(nextStateMsg, send)
		}()

		errFunc := func(err error, payload []byte, ce *pb.ChaincodeEvent, errFmt string, args ...string) *pb.ChaincodeMessage {
			if err != nil {
				// Send ERROR message to chaincode support and change state
				if payload == nil {
					payload = []byte(err.Error())
				}
				chaincodeLogger.Errorf(errFmt, args)
				return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: payload, Txid: msg.Txid, ChaincodeEvent: ce}
			}
			return nil
		}
		// Get the function and args from Payload
		input := &pb.ChaincodeInput{}
		unmarshalErr := proto.Unmarshal(msg.Payload, input)
		if nextStateMsg = errFunc(unmarshalErr, nil, nil, "[%s]Incorrect payload format. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ERROR.String()); nextStateMsg != nil {
			return
		}

		// Call chaincode's Run
		// Create the ChaincodeStub which the chaincode can use to callback
		stub := new(ChaincodeStub)

		err := stub.init(handler, msg.Txid, input, msg.Proposal)
		if nextStateMsg = errFunc(err, nil, stub.chaincodeEvent, "[%s]Init get error response [%s]. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ERROR.String()); nextStateMsg != nil {
			return
		}

		res := handler.cc.Init(stub)
		chaincodeLogger.Debugf("[%s]Init get response status: %d", shorttxid(msg.Txid), res.Status)

		if res.Status >= ERROR {
			err = fmt.Errorf("%s", res.Message)
			if nextStateMsg = errFunc(err, []byte(res.Message), stub.chaincodeEvent, "[%s]Init get error response [%s]. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ERROR.String()); nextStateMsg != nil {
				return
			}
		}

		resBytes, err := proto.Marshal(&res)
		if nextStateMsg = errFunc(err, nil, stub.chaincodeEvent, "[%s]Init marshal response error [%s]. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ERROR.String()); nextStateMsg != nil {
			return
		}

		// Send COMPLETED message to chaincode support and change state
		nextStateMsg = &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: resBytes, Txid: msg.Txid, ChaincodeEvent: stub.chaincodeEvent}
		chaincodeLogger.Debugf("[%s]Init succeeded. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_COMPLETED)
	}()
}

// beforeInit will initialize the chaincode if entering init from established.
func (handler *Handler) beforeInit(e *fsm.Event) {
	chaincodeLogger.Debugf("Entered state %s", handler.FSM.Current())
	msg, ok := e.Args[0].(*pb.ChaincodeMessage)
	if !ok {
		e.Cancel(fmt.Errorf("Received unexpected message type"))
		return
	}
	chaincodeLogger.Debugf("[%s]Received %s, initializing chaincode", shorttxid(msg.Txid), msg.Type.String())
	if msg.Type.String() == pb.ChaincodeMessage_INIT.String() {
		// Call the chaincode's Run function to initialize
		handler.handleInit(msg)
	}
}

// handleTransaction Handles request to execute a transaction.
func (handler *Handler) handleTransaction(msg *pb.ChaincodeMessage) {
	// The defer followed by triggering a go routine dance is needed to ensure that the previous state transition
	// is completed before the next one is triggered. The previous state transition is deemed complete only when
	// the beforeInit function is exited. Interesting bug fix!!
	go func() {
		//better not be nil
		var nextStateMsg *pb.ChaincodeMessage

		send := true

		defer func() {
			handler.triggerNextState(nextStateMsg, send)
		}()

		errFunc := func(err error, ce *pb.ChaincodeEvent, errStr string, args ...string) *pb.ChaincodeMessage {
			if err != nil {
				payload := []byte(err.Error())
				chaincodeLogger.Errorf(errStr, args)
				return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: payload, Txid: msg.Txid, ChaincodeEvent: ce}
			}
			return nil
		}

		// Get the function and args from Payload
		input := &pb.ChaincodeInput{}
		unmarshalErr := proto.Unmarshal(msg.Payload, input)
		if nextStateMsg = errFunc(unmarshalErr, nil, "[%s]Incorrect payload format. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ERROR.String()); nextStateMsg != nil {
			return
		}

		// Call chaincode's Run
		// Create the ChaincodeStub which the chaincode can use to callback
		stub := new(ChaincodeStub)
		err := stub.init(handler, msg.Txid, input, msg.Proposal)
		if nextStateMsg = errFunc(err, stub.chaincodeEvent, "[%s]Transaction execution failed. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ERROR.String()); nextStateMsg != nil {
			return
		}

		res := handler.cc.Invoke(stub)

		// Endorser will handle error contained in Response.
		resBytes, err := proto.Marshal(&res)
		if nextStateMsg = errFunc(err, stub.chaincodeEvent, "[%s]Transaction execution failed. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ERROR.String()); nextStateMsg != nil {
			return
		}

		// Send COMPLETED message to chaincode support and change state
		chaincodeLogger.Debugf("[%s]Transaction completed. Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_COMPLETED)
		nextStateMsg = &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: resBytes, Txid: msg.Txid, ChaincodeEvent: stub.chaincodeEvent}
	}()
}

// beforeTransaction will execute chaincode's Run if coming from a TRANSACTION event.
func (handler *Handler) beforeTransaction(e *fsm.Event) {
	msg, ok := e.Args[0].(*pb.ChaincodeMessage)
	if !ok {
		e.Cancel(fmt.Errorf("Received unexpected message type"))
		return
	}
	chaincodeLogger.Debugf("[%s]Received %s, invoking transaction on chaincode(Src:%s, Dst:%s)", shorttxid(msg.Txid), msg.Type.String(), e.Src, e.Dst)
	if msg.Type.String() == pb.ChaincodeMessage_TRANSACTION.String() {
		// Call the chaincode's Run function to invoke transaction
		handler.handleTransaction(msg)
	}
}

// afterResponse is called to deliver a response or error to the chaincode stub.
func (handler *Handler) afterResponse(e *fsm.Event) {
	msg, ok := e.Args[0].(*pb.ChaincodeMessage)
	if !ok {
		e.Cancel(fmt.Errorf("Received unexpected message type"))
		return
	}

	if err := handler.sendChannel(msg); err != nil {
		chaincodeLogger.Errorf("[%s]error sending %s (state:%s): %s", shorttxid(msg.Txid), msg.Type, handler.FSM.Current(), err)
	} else {
		chaincodeLogger.Debugf("[%s]Received %s, communicated (state:%s)", shorttxid(msg.Txid), msg.Type, handler.FSM.Current())
	}
}

func (handler *Handler) afterError(e *fsm.Event) {
	msg, ok := e.Args[0].(*pb.ChaincodeMessage)
	if !ok {
		e.Cancel(fmt.Errorf("Received unexpected message type"))
		return
	}

	/* TODO- revisit. This may no longer be needed with the serialized/streamlined messaging model
	 * There are two situations in which the ERROR event can be triggered:
	 * 1. When an error is encountered within handleInit or handleTransaction - some issue at the chaincode side; In this case there will be no responseChannel and the message has been sent to the validator.
	 * 2. The chaincode has initiated a request (get/put/del state) to the validator and is expecting a response on the responseChannel; If ERROR is received from validator, this needs to be notified on the responseChannel.
	 */
	if err := handler.sendChannel(msg); err == nil {
		chaincodeLogger.Debugf("[%s]Error received from validator %s, communicated(state:%s)", shorttxid(msg.Txid), msg.Type, handler.FSM.Current())
	}
}

// TODO: Implement method to get and put entire state map and not one key at a time?
// handleGetState communicates with the validator to fetch the requested state information from the ledger.
func (handler *Handler) handleGetState(key string, txid string) ([]byte, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return nil, err
	}

	defer handler.deleteChannel(txid)

	// Send GET_STATE message to validator chaincode support
	payload := []byte(key)
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: payload, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_STATE)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending GET_STATE %s", shorttxid(txid), err))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]GetState received payload %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		return responseMsg.Payload, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]GetState received error %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

// handlePutState communicates with the validator to put state information into the ledger.
func (handler *Handler) handlePutState(key string, value []byte, txid string) error {
	// Check if this is a transaction
	chaincodeLogger.Debugf("[%s]Inside putstate", shorttxid(txid))

	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.PutStateInfo{Key: key, Value: value})

	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return err
	}

	defer handler.deleteChannel(txid)

	// Send PUT_STATE message to validator chaincode support
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_PUT_STATE)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return errors.New(fmt.Sprintf("[%s]error sending PUT_STATE %s", msg.Txid, err))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully updated state", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		return nil
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s. Payload: %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR, responseMsg.Payload)
		return errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

// handleDelState communicates with the validator to delete a key from the state in the ledger.
func (handler *Handler) handleDelState(key string, txid string) error {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return err
	}

	defer handler.deleteChannel(txid)

	// Send DEL_STATE message to validator chaincode support
	payload := []byte(key)
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_DEL_STATE, Payload: payload, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_DEL_STATE)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return errors.New(fmt.Sprintf("[%s]error sending DEL_STATE %s", shorttxid(msg.Txid), pb.ChaincodeMessage_DEL_STATE))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully deleted state", msg.Txid, pb.ChaincodeMessage_RESPONSE)
		return nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s. Payload: %s", msg.Txid, pb.ChaincodeMessage_ERROR, responseMsg.Payload)
		return errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

func (handler *Handler) handleGetStateByRange(startKey, endKey string, txid string) (*pb.QueryResponse, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return nil, err
	}

	defer handler.deleteChannel(txid)

	// Send GET_STATE_BY_RANGE message to validator chaincode support
	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.GetStateByRange{StartKey: startKey, EndKey: endKey})

	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE_BY_RANGE, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_STATE_BY_RANGE)

	var responseMsg pb.ChaincodeMessage
	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_STATE_BY_RANGE))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully got range", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)

		rangeQueryResponse := &pb.QueryResponse{}
		if err = proto.Unmarshal(responseMsg.Payload, rangeQueryResponse); err != nil {
			return nil, errors.New(fmt.Sprintf("[%s]GetStateByRangeResponse unmarshall error", shorttxid(responseMsg.Txid)))
		}

		return rangeQueryResponse, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("Incorrect chaincode message %s received. Expecting %s or %s", responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

func (handler *Handler) handleQueryStateNext(id, txid string) (*pb.QueryResponse, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return nil, err
	}

	defer handler.deleteChannel(txid)

	// Send QUERY_STATE_NEXT message to validator chaincode support
	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.QueryStateNext{Id: id})

	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_QUERY_STATE_NEXT)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_QUERY_STATE_NEXT))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully got range", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)

		queryResponse := &pb.QueryResponse{}
		if err = proto.Unmarshal(responseMsg.Payload, queryResponse); err != nil {
			return nil, errors.New(fmt.Sprintf("[%s]unmarshall error", shorttxid(responseMsg.Txid)))
		}

		return queryResponse, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("Incorrect chaincode message %s received. Expecting %s or %s", responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

func (handler *Handler) handleQueryStateClose(id, txid string) (*pb.QueryResponse, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return nil, err
	}

	defer handler.deleteChannel(txid)

	// Send QUERY_STATE_CLOSE message to validator chaincode support
	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.QueryStateClose{Id: id})

	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_QUERY_STATE_CLOSE)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_QUERY_STATE_CLOSE))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully got range", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)

		queryResponse := &pb.QueryResponse{}
		if err = proto.Unmarshal(responseMsg.Payload, queryResponse); err != nil {
			return nil, errors.New(fmt.Sprintf("[%s]unmarshall error", shorttxid(responseMsg.Txid)))
		}

		return queryResponse, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("Incorrect chaincode message %s received. Expecting %s or %s", responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

func (handler *Handler) handleGetQueryResult(query string, txid string) (*pb.QueryResponse, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return nil, err
	}

	defer handler.deleteChannel(txid)

	// Send GET_QUERY_RESULT message to validator chaincode support
	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.GetQueryResult{Query: query})

	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_QUERY_RESULT, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_QUERY_RESULT)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_QUERY_RESULT))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully got range", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)

		executeQueryResponse := &pb.QueryResponse{}
		if err = proto.Unmarshal(responseMsg.Payload, executeQueryResponse); err != nil {
			return nil, errors.New(fmt.Sprintf("[%s]unmarshall error", shorttxid(responseMsg.Txid)))
		}

		return executeQueryResponse, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("Incorrect chaincode message %s received. Expecting %s or %s", responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

func (handler *Handler) handleGetHistoryForKey(key string, txid string) (*pb.QueryResponse, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return nil, err
	}

	defer handler.deleteChannel(txid)

	// Send GET_HISTORY_FOR_KEY message to validator chaincode support
	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.GetHistoryForKey{Key: key})

	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_HISTORY_FOR_KEY, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_HISTORY_FOR_KEY)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_HISTORY_FOR_KEY))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully got range", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)

		getHistoryForKeyResponse := &pb.QueryResponse{}
		if err = proto.Unmarshal(responseMsg.Payload, getHistoryForKeyResponse); err != nil {
			return nil, errors.New(fmt.Sprintf("[%s]unmarshall error", shorttxid(responseMsg.Txid)))
		}

		return getHistoryForKeyResponse, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("Incorrect chaincode message %s received. Expecting %s or %s", responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

//----------------------------inkchain wallet

func (handler *Handler) handleTransfer(trans *kvtranset.KVTranSet, txid string) error {
	// Check if this is a transaction
	chaincodeLogger.Debugf("[%s]Inside transfer", shorttxid(txid))
	//we constructed a valid object. No need to check for error
	var tranSet []*pb.Transfer
	for _, tran := range trans.Trans {
		protoTran := pb.Transfer{To: []byte(strings.ToLower(tran.To)), BalanceType: []byte(tran.BalanceType), Amount: tran.Amount}
		tranSet = append(tranSet, &protoTran)
	}
	transferInfo := &pb.TransferInfo{TranSet: tranSet}
	payloadBytes, err := proto.Marshal(transferInfo)
	// Create the channel on which to communicate the response from validating peer
	if err != nil {
		return err
	}
	var respChan chan pb.ChaincodeMessage
	if respChan, err = handler.createChannel(txid); err != nil {
		return err
	}

	defer handler.deleteChannel(txid)

	// Send Transfer message to validator chaincode support
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSFER, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_TRANSFER)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return errors.New(fmt.Sprintf("[%s]error sending TRANSFER %s", msg.Txid, err))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully transfered", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		return nil
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s. Payload: %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR, responseMsg.Payload)
		return errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

func (handler *Handler) handleCrossTransfer(trans *kvcrosstranset.KVCrossTranSet, txId string) error {
	// Check if this is a transaction
	chaincodeLogger.Debugf("[%s]Inside cross transfer", shorttxid(txId))
	//we constructed a valid object. No need to check for error
	var tranSet []*pb.CrossTransfer
	for _, tran := range trans.Trans {
		protoTran := pb.CrossTransfer{To: []byte(strings.ToLower(tran.To)), Amount: tran.Amount}
		tranSet = append(tranSet, &protoTran)
	}
	crossTransferInfo := &pb.CrossTransferInfo{TranSet: tranSet, PubTxId: []byte(trans.PubTxId), FromPlatForm: []byte(trans.FromPlatForm)}
	payloadBytes, err := proto.Marshal(crossTransferInfo)
	// Create the channel on which to communicate the response from validating peer
	if err != nil {
		return err
	}
	var respChan chan pb.ChaincodeMessage
	if respChan, err = handler.createChannel(txId); err != nil {
		return err
	}

	defer handler.deleteChannel(txId)

	// Send Transfer message to validator chaincode support
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_CROSS_TRANSFER, Payload: payloadBytes, Txid: txId}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_CROSS_TRANSFER)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return errors.New(fmt.Sprintf("[%s]error sending CROSS_TRANSFER %s", msg.Txid, err))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully cross transferred", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		return nil
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s. Payload: %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR, responseMsg.Payload)
		return errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

func (handler *Handler) handleGetAccount(address string, txid string) (*wallet.Account, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return nil, err
	}
	defer handler.deleteChannel(txid)
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_ACCOUNT, Payload: []byte(address), Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_ACCOUNT)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending GET_ACCOUNT %s", shorttxid(txid), err))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]GetAccount received payload %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		account := &wallet.Account{}
		if jsonErr := json.Unmarshal(responseMsg.Payload, account); jsonErr != nil {
			return nil, jsonErr
		}
		return account, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]GetAccount received error %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

// handlePutState communicates with the validator to put state information into the ledger.
func (handler *Handler) handleIssueToken(address string, balanceType string, amount *big.Int, txid string) error {
	// Check if this is a transaction
	chaincodeLogger.Debugf("[%s]Inside issue token", shorttxid(txid))
	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.IssueTokenInfo{Address: []byte(address), BalanceType: []byte(balanceType), Amount: amount.Bytes()})

	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return err
	}

	defer handler.deleteChannel(txid)

	// Send PUT_STATE message to validator chaincode support
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ISSUE_TOKEN, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_ISSUE_TOKEN)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return errors.New(fmt.Sprintf("[%s]error sending SET_BALANCE %s", msg.Txid, err))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully set balance", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		return nil
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s. Payload: %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR, responseMsg.Payload)
		return errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

// handleCalcFee communicates to calculate the fee of content.
func (handler *Handler) handleCalcFee(content string, txId string) (*big.Int, error) {
	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txId); err != nil {
		return nil, err
	}

	defer handler.deleteChannel(txId)

	payload := []byte(content)
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_FEE, Payload: payload, Txid: txId}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_GET_FEE)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return nil, errors.New(fmt.Sprintf("[%s]error sending GET_FEE %s", shorttxid(txId), err))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]GetFee received payload %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		fee := big.NewInt(0)
		fee.SetBytes(responseMsg.Payload[:])
		return fee, nil
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]GetState received error %s", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return nil, errors.New(string(responseMsg.Payload[:]))
	}

	// Incorrect chaincode message received
	return nil, errors.New(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR))
}

//---------------------------

func (handler *Handler) createResponse(status int32, payload []byte) pb.Response {
	return pb.Response{Status: status, Payload: payload}
}

// handleInvokeChaincode communicates with the validator to invoke another chaincode.
func (handler *Handler) handleInvokeChaincode(chaincodeName string, args [][]byte, txid string) pb.Response {
	//we constructed a valid object. No need to check for error
	payloadBytes, _ := proto.Marshal(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: chaincodeName}, Input: &pb.ChaincodeInput{Args: args}})

	// Create the channel on which to communicate the response from validating peer
	var respChan chan pb.ChaincodeMessage
	var err error
	if respChan, err = handler.createChannel(txid); err != nil {
		return handler.createResponse(ERROR, []byte(err.Error()))
	}

	defer handler.deleteChannel(txid)

	// Send INVOKE_CHAINCODE message to validator chaincode support
	msg := &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: payloadBytes, Txid: txid}
	chaincodeLogger.Debugf("[%s]Sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_INVOKE_CHAINCODE)

	var responseMsg pb.ChaincodeMessage

	if responseMsg, err = handler.sendReceive(msg, respChan); err != nil {
		return handler.createResponse(ERROR, []byte(fmt.Sprintf("[%s]error sending %s", shorttxid(msg.Txid), pb.ChaincodeMessage_INVOKE_CHAINCODE)))
	}

	if responseMsg.Type.String() == pb.ChaincodeMessage_RESPONSE.String() {
		// Success response
		chaincodeLogger.Debugf("[%s]Received %s. Successfully invoked chaincode", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
		respMsg := &pb.ChaincodeMessage{}
		if err = proto.Unmarshal(responseMsg.Payload, respMsg); err != nil {
			return handler.createResponse(ERROR, []byte(err.Error()))
		}
		if respMsg.Type == pb.ChaincodeMessage_COMPLETED {
			// Success response
			chaincodeLogger.Debugf("[%s]Received %s. Successfully invoed chaincode", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_RESPONSE)
			res := &pb.Response{}
			if err = proto.Unmarshal(respMsg.Payload, res); err != nil {
				return handler.createResponse(ERROR, []byte(err.Error()))
			}
			return *res
		}
		chaincodeLogger.Errorf("[%s]Received %s. Error from chaincode", shorttxid(responseMsg.Txid), respMsg.Type.String())
		return handler.createResponse(ERROR, responseMsg.Payload)
	}
	if responseMsg.Type.String() == pb.ChaincodeMessage_ERROR.String() {
		// Error response
		chaincodeLogger.Errorf("[%s]Received %s.", shorttxid(responseMsg.Txid), pb.ChaincodeMessage_ERROR)
		return handler.createResponse(ERROR, responseMsg.Payload)
	}

	// Incorrect chaincode message received
	return handler.createResponse(ERROR, []byte(fmt.Sprintf("[%s]Incorrect chaincode message %s received. Expecting %s or %s", shorttxid(responseMsg.Txid), responseMsg.Type, pb.ChaincodeMessage_RESPONSE, pb.ChaincodeMessage_ERROR)))
}

// handleMessage message handles loop for shim side of chaincode/validator stream.
func (handler *Handler) handleMessage(msg *pb.ChaincodeMessage) error {
	if msg.Type == pb.ChaincodeMessage_KEEPALIVE {
		// Received a keep alive message, we don't do anything with it for now
		// and it does not touch the state machine
		return nil
	}
	chaincodeLogger.Debugf("[%s]Handling ChaincodeMessage of type: %s(state:%s)", shorttxid(msg.Txid), msg.Type, handler.FSM.Current())
	if handler.FSM.Cannot(msg.Type.String()) {
		err := errors.New(fmt.Sprintf("[%s]Chaincode handler FSM cannot handle message (%s) with payload size (%d) while in state: %s", msg.Txid, msg.Type.String(), len(msg.Payload), handler.FSM.Current()))
		handler.serialSend(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: []byte(err.Error()), Txid: msg.Txid})
		return err
	}
	err := handler.FSM.Event(msg.Type.String(), msg)
	return filterError(err)
}

// filterError filters the errors to allow NoTransitionError and CanceledError to not propagate for cases where embedded Err == nil.
func filterError(errFromFSMEvent error) error {
	if errFromFSMEvent != nil {
		if noTransitionErr, ok := errFromFSMEvent.(*fsm.NoTransitionError); ok {
			if noTransitionErr.Err != nil {
				// Only allow NoTransitionError's, all others are considered true error.
				return errFromFSMEvent
			}
		}
		if canceledErr, ok := errFromFSMEvent.(*fsm.CanceledError); ok {
			if canceledErr.Err != nil {
				// Only allow NoTransitionError's, all others are considered true error.
				return canceledErr
				//t.Error("expected only 'NoTransitionError'")
			}
			chaincodeLogger.Debugf("Ignoring CanceledError: %s", canceledErr)
		}
	}
	return nil
}
