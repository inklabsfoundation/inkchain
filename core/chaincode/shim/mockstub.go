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

// Package shim provides APIs for the chaincode to access its state
// variables, transaction context and call other chaincodes.
package shim

import (
	"container/list"
	"errors"
	"fmt"
	"strings"

	"encoding/json"
	"math/big"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/inklabsfoundation/inkchain/common/util"
	"github.com/inklabsfoundation/inkchain/core/wallet"
	"github.com/inklabsfoundation/inkchain/protos/ledger/queryresult"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet/kvtranset"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"github.com/op/go-logging"
)

// Logger for the shim package.
var mockLogger = logging.MustGetLogger("mock")

// MockStub is an implementation of ChaincodeStubInterface for unit testing chaincode.
// Use this instead of ChaincodeStub in your chaincode's unit test calls to Init or Invoke.
type MockStub struct {
	// arguments the stub was called with
	args [][]byte

	// A pointer back to the chaincode that will invoke this, set by constructor.
	// If a peer calls this stub, the chaincode will be invoked from here.
	cc Chaincode

	// A nice name that can be used for logging
	Name string

	// State keeps name value pairs
	State map[string][]byte

	// Keys stores the list of mapped values in lexical order
	Keys *list.List

	// registered list of other MockStub chaincodes that can be called from this MockStub
	Invokables map[string]*MockStub

	// stores a transaction uuid while being Invoked / Deployed
	// TODO if a chaincode uses recursion this may need to be a stack of TxIDs or possibly a reference counting map
	TxID string

	TxTimestamp *timestamp.Timestamp

	// mocked signedProposal
	signedProposal *pb.SignedProposal
}

func (stub *MockStub) GetTxID() string {
	return stub.TxID
}

func (stub *MockStub) GetArgs() [][]byte {
	return stub.args
}

func (stub *MockStub) GetStringArgs() []string {
	args := stub.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}

func (stub *MockStub) GetFunctionAndParameters() (function string, params []string) {
	allargs := stub.GetStringArgs()
	function = ""
	params = []string{}
	if len(allargs) >= 1 {
		function = allargs[0]
		params = allargs[1:]
	}
	return
}

// Used to indicate to a chaincode that it is part of a transaction.
// This is important when chaincodes invoke each other.
// MockStub doesn't support concurrent transactions at present.
func (stub *MockStub) MockTransactionStart(txid string) {
	stub.TxID = txid
	stub.setSignedProposal(&pb.SignedProposal{})
	stub.setTxTimestamp(util.CreateUtcTimestamp())
}

// End a mocked transaction, clearing the UUID.
func (stub *MockStub) MockTransactionEnd(uuid string) {
	stub.signedProposal = nil
	stub.TxID = ""
}

// Register a peer chaincode with this MockStub
// invokableChaincodeName is the name or hash of the peer
// otherStub is a MockStub of the peer, already intialised
func (stub *MockStub) MockPeerChaincode(invokableChaincodeName string, otherStub *MockStub) {
	stub.Invokables[invokableChaincodeName] = otherStub
}

// Initialise this chaincode,  also starts and ends a transaction.
func (stub *MockStub) MockInit(uuid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	res := stub.cc.Init(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

// Invoke this chaincode, also starts and ends a transaction.
func (stub *MockStub) MockInvoke(uuid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

// Invoke this chaincode, also starts and ends a transaction.
func (stub *MockStub) MockInvokeWithSignedProposal(uuid string, args [][]byte, sp *pb.SignedProposal) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	stub.signedProposal = sp
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

// GetState retrieves the value for a given key from the ledger
func (stub *MockStub) GetState(key string) ([]byte, error) {
	value := stub.State[key]
	mockLogger.Debug("MockStub", stub.Name, "Getting", key, value)
	return value, nil
}

// PutState writes the specified `value` and `key` into the ledger.
func (stub *MockStub) PutState(key string, value []byte) error {
	if stub.TxID == "" {
		mockLogger.Error("Cannot PutState without a transactions - call stub.MockTransactionStart()?")
		return errors.New("Cannot PutState without a transactions - call stub.MockTransactionStart()?")
	}

	mockLogger.Debug("MockStub", stub.Name, "Putting", key, value)
	stub.State[key] = value

	// insert key into ordered list of keys
	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(string)
		comp := strings.Compare(key, elemValue)
		mockLogger.Debug("MockStub", stub.Name, "Compared", key, elemValue, " and got ", comp)
		if comp < 0 {
			// key < elem, insert it before elem
			stub.Keys.InsertBefore(key, elem)
			mockLogger.Debug("MockStub", stub.Name, "Key", key, " inserted before", elem.Value)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			mockLogger.Debug("MockStub", stub.Name, "Key", key, "already in State")
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.Keys.PushBack(key)
				mockLogger.Debug("MockStub", stub.Name, "Key", key, "appended")
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.Keys.Len() == 0 {
		stub.Keys.PushFront(key)
		mockLogger.Debug("MockStub", stub.Name, "Key", key, "is first element in list")
	}

	return nil
}

// DelState removes the specified `key` and its value from the ledger.
func (stub *MockStub) DelState(key string) error {
	mockLogger.Debug("MockStub", stub.Name, "Deleting", key, stub.State[key])
	delete(stub.State, key)

	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		if strings.Compare(key, elem.Value.(string)) == 0 {
			stub.Keys.Remove(elem)
		}
	}

	return nil
}

func (stub *MockStub) GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error) {
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, err
	}
	return NewMockStateRangeQueryIterator(stub, startKey, endKey), nil
}

// GetQueryResult function can be invoked by a chaincode to perform a
// rich query against state database.  Only supported by state database implementations
// that support rich query.  The query string is in the syntax of the underlying
// state database. An iterator is returned which can be used to iterate (next) over
// the query result set
func (stub *MockStub) GetQueryResult(query string) (StateQueryIteratorInterface, error) {
	// Not implemented since the mock engine does not have a query engine.
	// However, a very simple query engine that supports string matching
	// could be implemented to test that the framework supports queries
	return nil, errors.New("Not Implemented")
}

// GetHistoryForKey function can be invoked by a chaincode to return a history of
// key values across time. GetHistoryForKey is intended to be used for read-only queries.
func (stub *MockStub) GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error) {
	return nil, errors.New("Not Implemented")
}

//GetStateByPartialCompositeKey function can be invoked by a chaincode to query the
//state based on a given partial composite key. This function returns an
//iterator which can be used to iterate over all composite keys whose prefix
//matches the given partial composite key. This function should be used only for
//a partial composite key. For a full composite key, an iter with empty response
//would be returned.
func (stub *MockStub) GetStateByPartialCompositeKey(objectType string, attributes []string) (StateQueryIteratorInterface, error) {
	partialCompositeKey, err := stub.CreateCompositeKey(objectType, attributes)
	if err != nil {
		return nil, err
	}
	return NewMockStateRangeQueryIterator(stub, partialCompositeKey, partialCompositeKey+string(maxUnicodeRuneValue)), nil
}

// CreateCompositeKey combines the list of attributes
//to form a composite key.
func (stub *MockStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return createCompositeKey(objectType, attributes)
}

// SplitCompositeKey splits the composite key into attributes
// on which the composite key was formed.
func (stub *MockStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return splitCompositeKey(compositeKey)
}

// InvokeChaincode calls a peered chaincode.
// E.g. stub1.InvokeChaincode("stub2Hash", funcArgs, channel)
// Before calling this make sure to create another MockStub stub2, call stub2.MockInit(uuid, func, args)
// and register it with stub1 by calling stub1.MockPeerChaincode("stub2Hash", stub2)
func (stub *MockStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {
	// Internally we use chaincode name as a composite name
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}
	// TODO "args" here should possibly be a serialized pb.ChaincodeInput
	otherStub := stub.Invokables[chaincodeName]
	mockLogger.Debug("MockStub", stub.Name, "Invoking peer chaincode", otherStub.Name, args)
	//	function, strings := getFuncArgs(args)
	res := otherStub.MockInvoke(stub.TxID, args)
	mockLogger.Debug("MockStub", stub.Name, "Invoked peer chaincode", otherStub.Name, "got", fmt.Sprintf("%+v", res))
	return res
}

// warning!  This method can not produce right outputs cause the sender is obtained from ChaincodeInvokeSpec
func (stub *MockStub) Transfer(to string, balanceType string, amount *big.Int) error {
	if stub.TxID == "" {
		mockLogger.Error("Cannot Transfer without a transactions - call stub.MockTransactionStart()?")
		return errors.New("Cannot Transfer without a transactions - call stub.MockTransactionStart()?")
	}
	/*
		to = strings.ToLower(to)
		mockLogger.Debug("MockStub", stub.Name, "Transfer To", to, amount)

		toAddress, err := wallet.HexToAddress(to)
		if err != nil {
			return err
		}

			sig, err := wallet.SignatureStringToBytes(signature)
			if err != nil {
				return errors.New(err.Error())
			}
			fromAddressExtracted, err := wallet.GetSenderFromSignature(fromAddress, toAddress, amount, sig)
			if fromAddressExtracted.ToString() != from {
				return errors.New("sender mismatch")
			}

		toAccount := &wallet.Account{}
		if value, ok := stub.State[fromAddress.ToString()]; ok {
			jsonErr := json.Unmarshal(value, fromAccount)
			if jsonErr != nil {
				return jsonErr
			}
			if fromAccount.Balance == nil {
				return errors.New("balance" + balanceType + "not exists")
			}
			fromBalance, ok := fromAccount.Balance[balanceType]
			if !ok {
				return errors.New("balance" + wallet.MAIN_BALANCE_NAME + "not exists")
			}
			if fromBalance.Cmp(amount) < 0 {
				return errors.New("insufficient balance for sender")
			}
			fromBalance = fromBalance.Sub(fromBalance, amount)

			fromAccountBytes, jsonErr := json.Marshal(fromAccount)

			if jsonErr != nil {
				return errors.New("error marshaling sender account")
			}
			stub.State[fromAddress.ToString()] = fromAccountBytes
			key := fromAddress.ToString()
			// insert key into ordered list of keys
			for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
				elemValue := elem.Value.(string)
				comp := strings.Compare(key, elemValue)
				mockLogger.Debug("MockStub", stub.Name, "Compared", key, elemValue, " and got ", comp)
				if comp < 0 {
					// key < elem, insert it before elem
					stub.Keys.InsertBefore(key, elem)
					mockLogger.Debug("MockStub", stub.Name, "Key", key, " inserted before", elem.Value)
					break
				} else if comp == 0 {
					// keys exists, no need to change
					mockLogger.Debug("MockStub", stub.Name, "Key", key, "already in State")
					break
				} else { // comp > 0
					// key > elem, keep looking unless this is the end of the list
					if elem.Next() == nil {
						stub.Keys.PushBack(key)
						mockLogger.Debug("MockStub", stub.Name, "Key", key, "appended")
						break
					}
				}
			}

			// special case for empty Keys list
			if stub.Keys.Len() == 0 {
				stub.Keys.PushFront(key)
				mockLogger.Debug("MockStub", stub.Name, "Key", key, "is first element in list")
			}
			value, ok = stub.State[toAddress.ToString()]
			if ok {
				jsonErr = json.Unmarshal(value, &toAccount)
				if jsonErr != nil {
					return jsonErr
				}
			}
			if toAccount.Balance == nil {
				toAccount.Balance = make(map[string]*big.Int)
			}
			toBalance, ok := toAccount.Balance[balanceType]
			if !ok {
				toBalance = big.NewInt(0)
				toAccount.Balance[balanceType] = toBalance
			}
			toBalance.Add(toBalance, amount)
			toAccountBytes, jsonErr := json.Marshal(toAccount)
			stub.State[to] = toAccountBytes

			key = to
			// insert key into ordered list of keys
			for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
				elemValue := elem.Value.(string)
				comp := strings.Compare(key, elemValue)
				mockLogger.Debug("MockStub", stub.Name, "Compared", key, elemValue, " and got ", comp)
				if comp < 0 {
					// key < elem, insert it before elem
					stub.Keys.InsertBefore(key, elem)
					mockLogger.Debug("MockStub", stub.Name, "Key", key, " inserted before", elem.Value)
					break
				} else if comp == 0 {
					// keys exists, no need to change
					mockLogger.Debug("MockStub", stub.Name, "Key", key, "already in State")
					break
				} else { // comp > 0
					// key > elem, keep looking unless this is the end of the list
					if elem.Next() == nil {
						stub.Keys.PushBack(key)
						mockLogger.Debug("MockStub", stub.Name, "Key", key, "appended")
						break
					}
				}
			}

			// special case for empty Keys list
			if stub.Keys.Len() == 0 {
				stub.Keys.PushFront(key)
				mockLogger.Debug("MockStub", stub.Name, "Key", key, "is first element in list")
			}
			return nil
		}
	*/
	return errors.New(" this function could not be used in mock invocation")
}

func (stub *MockStub) CrossTransfer(to string, balanceType string, amount *big.Int, pubTxId string, fromPlatform string) error {
	if stub.TxID == "" {
		mockLogger.Error("Cannot Transfer without a transactions - call stub.MockTransactionStart()?")
		return errors.New("Cannot Transfer without a transactions - call stub.MockTransactionStart()?")
	}
	return errors.New(" this function could not be used in mock invocation")
}
func (stub *MockStub) MultiTransfer(trans *kvtranset.KVTranSet) error {
	return errors.New(" this function could not be used in mock invocation")
}
func (stub *MockStub) GetAccount(address string) (*wallet.Account, error) {
	address = strings.ToLower(address)
	if accountBytes, ok := stub.State[address]; ok {
		account := &wallet.Account{}
		jsonErr := json.Unmarshal(accountBytes, account)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return account, nil
	}
	return nil, errors.New("getAccount error")
}

func (stub *MockStub) IssueToken(address string, balanceType string, amount *big.Int) error {
	address = strings.ToLower(address)
	if stub.TxID == "" {
		mockLogger.Error("Cannot issue token without a transactions - call stub.MockTransactionStart()?")
		return errors.New("Cannot issue token without a transactions - call stub.MockTransactionStart()?")
	}
	account := &wallet.Account{}
	accountJson, ok := stub.State[address]
	if ok {
		err := json.Unmarshal(accountJson, account)
		if err == nil {
			_, ok := account.Balance[balanceType]
			if ok {
				return errors.New("balance exists in this address")
			}
		}

	}
	if account.Balance == nil {
		account.Balance = make(map[string]*big.Int)
	}

	account.Balance[balanceType] = amount
	accountBytes, jsonErr := json.Marshal(account)
	if jsonErr != nil {
		return jsonErr
	}
	stub.State[address] = accountBytes
	// insert key into ordered list of keys
	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(string)
		comp := strings.Compare(address, elemValue)
		mockLogger.Debug("MockStub", stub.Name, "Compared", address, elemValue, " and got ", comp)
		if comp < 0 {
			// key < elem, insert it before elem
			stub.Keys.InsertBefore(address, elem)
			mockLogger.Debug("MockStub", stub.Name, "Key", address, " inserted before", elem.Value)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			mockLogger.Debug("MockStub", stub.Name, "Key", address, "already in State")
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.Keys.PushBack(address)
				mockLogger.Debug("MockStub", stub.Name, "Key", address, "appended")
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.Keys.Len() == 0 {
		stub.Keys.PushFront(address)
		mockLogger.Debug("MockStub", stub.Name, "Key", address, "is first element in list")
	}

	return nil
	return nil
}

// Not implemented
func (stub *MockStub) GetCreator() ([]byte, error) {
	return nil, nil
}

// Not implemented
func (stub *MockStub) GetTransient() (map[string][]byte, error) {
	return nil, nil
}

// Not implemented
func (stub *MockStub) GetBinding() ([]byte, error) {
	return nil, nil
}
func (stub *MockStub) GetSender() (string, error) {
	return "", nil
}

// Not implemented
func (stub *MockStub) GetSignedProposal() (*pb.SignedProposal, error) {
	return stub.signedProposal, nil
}

func (stub *MockStub) setSignedProposal(sp *pb.SignedProposal) {
	stub.signedProposal = sp
}

// Not implemented
func (stub *MockStub) GetArgsSlice() ([]byte, error) {
	return nil, nil
}

func (stub *MockStub) setTxTimestamp(time *timestamp.Timestamp) {
	stub.TxTimestamp = time
}

func (stub *MockStub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	if stub.TxTimestamp == nil {
		return nil, errors.New("TxTimestamp not set.")
	}
	return stub.TxTimestamp, nil
}

// Not implemented
func (stub *MockStub) SetEvent(name string, payload []byte) error {
	return nil
}

// Constructor to initialise the internal State map
func NewMockStub(name string, cc Chaincode) *MockStub {
	mockLogger.Debug("MockStub(", name, cc, ")")
	s := new(MockStub)
	s.Name = name
	s.cc = cc
	s.State = make(map[string][]byte)
	s.Invokables = make(map[string]*MockStub)
	s.Keys = list.New()

	return s
}

/*****************************
 Range Query Iterator
*****************************/

type MockStateRangeQueryIterator struct {
	Closed   bool
	Stub     *MockStub
	StartKey string
	EndKey   string
	Current  *list.Element
}

// HasNext returns true if the range query iterator contains additional keys
// and values.
func (iter *MockStateRangeQueryIterator) HasNext() bool {
	if iter.Closed {
		// previously called Close()
		mockLogger.Error("HasNext() but already closed")
		return false
	}

	if iter.Current == nil {
		mockLogger.Error("HasNext() couldn't get Current")
		return false
	}

	current := iter.Current
	for current != nil {
		// if this is an open-ended query for all keys, return true
		if iter.StartKey == "" && iter.EndKey == "" {
			return true
		}
		comp1 := strings.Compare(current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(current.Value.(string), iter.EndKey)
		if comp1 >= 0 {
			if comp2 <= 0 {
				mockLogger.Debug("HasNext() got next")
				return true
			} else {
				mockLogger.Debug("HasNext() but no next")
				return false

			}
		}
		current = current.Next()
	}

	// we've reached the end of the underlying values
	mockLogger.Debug("HasNext() but no next")
	return false
}

// Next returns the next key and value in the range query iterator.
func (iter *MockStateRangeQueryIterator) Next() (*queryresult.KV, error) {
	if iter.Closed == true {
		mockLogger.Error("MockStateRangeQueryIterator.Next() called after Close()")
		return nil, errors.New("MockStateRangeQueryIterator.Next() called after Close()")
	}

	if iter.HasNext() == false {
		mockLogger.Error("MockStateRangeQueryIterator.Next() called when it does not HaveNext()")
		return nil, errors.New("MockStateRangeQueryIterator.Next() called when it does not HaveNext()")
	}

	for iter.Current != nil {
		comp1 := strings.Compare(iter.Current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(iter.Current.Value.(string), iter.EndKey)
		// compare to start and end keys. or, if this is an open-ended query for
		// all keys, it should always return the key and value
		if (comp1 >= 0 && comp2 <= 0) || (iter.StartKey == "" && iter.EndKey == "") {
			key := iter.Current.Value.(string)
			value, err := iter.Stub.GetState(key)
			iter.Current = iter.Current.Next()
			return &queryresult.KV{Key: key, Value: value}, err
		}
		iter.Current = iter.Current.Next()
	}
	mockLogger.Error("MockStateRangeQueryIterator.Next() went past end of range")
	return nil, errors.New("MockStateRangeQueryIterator.Next() went past end of range")
}

// Close closes the range query iterator. This should be called when done
// reading from the iterator to free up resources.
func (iter *MockStateRangeQueryIterator) Close() error {
	if iter.Closed == true {
		mockLogger.Error("MockStateRangeQueryIterator.Close() called after Close()")
		return errors.New("MockStateRangeQueryIterator.Close() called after Close()")
	}

	iter.Closed = true
	return nil
}

func (iter *MockStateRangeQueryIterator) Print() {
	mockLogger.Debug("MockStateRangeQueryIterator {")
	mockLogger.Debug("Closed?", iter.Closed)
	mockLogger.Debug("Stub", iter.Stub)
	mockLogger.Debug("StartKey", iter.StartKey)
	mockLogger.Debug("EndKey", iter.EndKey)
	mockLogger.Debug("Current", iter.Current)
	mockLogger.Debug("HasNext?", iter.HasNext())
	mockLogger.Debug("}")
}

func NewMockStateRangeQueryIterator(stub *MockStub, startKey string, endKey string) *MockStateRangeQueryIterator {
	mockLogger.Debug("NewMockStateRangeQueryIterator(", stub, startKey, endKey, ")")
	iter := new(MockStateRangeQueryIterator)
	iter.Closed = false
	iter.Stub = stub
	iter.StartKey = startKey
	iter.EndKey = endKey
	iter.Current = stub.Keys.Front()

	iter.Print()

	return iter
}

func getBytes(function string, args []string) [][]byte {
	bytes := make([][]byte, 0, len(args)+1)
	bytes = append(bytes, []byte(function))
	for _, s := range args {
		bytes = append(bytes, []byte(s))
	}
	return bytes
}

func getFuncArgs(bytes [][]byte) (string, []string) {
	mockLogger.Debugf("getFuncArgs(%x)", bytes)
	function := string(bytes[0])
	args := make([]string, len(bytes)-1)
	for i := 1; i < len(bytes); i++ {
		mockLogger.Debugf("getFuncArgs - i:%x, len(bytes):%x", i, len(bytes))
		args[i-1] = string(bytes[i])
	}
	return function, args
}
