/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
)

const (
	//func name
	GetBalance string = "getBalance"
	GetAccount string = "getAccount"
	Transfer   string = "transfer"
	Counter    string = "counter"
	Sender     string = "sender"
)

// User chaincode for token operations
// After a token issued, users can use this chaincode to make query or transfer operations.
type tokenChaincode struct {
}

// Init func
func (t *tokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("token user chaincode Init.")
	return shim.Success([]byte("Init success."))
}

// Invoke func
func (t *tokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("token user chaincode Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case GetBalance:
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2.")
		}
		return t.getBalance(stub, args)

	case GetAccount:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		return t.getAccount(stub, args)

	case Transfer:
		if len(args) != 3 {
			return shim.Error("Incorrect number of arguments. Expecting 3")
		}
		return t.transfer(stub, args)

	case Counter:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}
		return t.getCounter(stub, args)

	case Sender:
		sender, err := stub.GetSender()
		if err != nil {
			return shim.Error("Get sender failed.")
		}
		return shim.Success([]byte(sender))

	}

	return shim.Error("Invalid invoke function name. Expecting \"getBalance\", \"getAccount\", \"transfer\", \"counter\" or \"sender\".")
}

// getBalance
// Get the balance of a specific token type in an account
func (t *tokenChaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Address
	var BalanceType string // Token type
	var err error

	A = strings.ToLower(args[0])
	BalanceType = args[1]
	// Get the state from the ledger
	account, err := stub.GetAccount(A)
	if err != nil {
		jsonResp := "{\"Error\":\"account not exists\"}"
		return shim.Error(jsonResp)
	}

	if account == nil || account.Balance[BalanceType] == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"" + BalanceType + "\":\"" + account.Balance[BalanceType].String() + "\"}"
	return shim.Success([]byte(jsonResp))
}

// getAccount
// Get the balances of all token types in an account
func (t *tokenChaincode) getAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Address
	var err error

	A = strings.ToLower(args[0])
	// Get the state from the ledger
	account, err := stub.GetAccount(A)
	if err != nil {
		jsonResp := "{\"Error\":\"account not exists\"}"
		return shim.Error(jsonResp)
	}

	if account == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}
	balanceJson, jsonErr := json.Marshal(account.Balance)
	if jsonErr != nil {
		return shim.Error(jsonErr.Error())
	}
	jsonResp := "{\"Name\":\"" + A + "\",\"Balance\":\"" + string(balanceJson[:]) + "\"}"
	return shim.Success([]byte(jsonResp))
}

// transfer
// Send tokens to the specified address
func (t *tokenChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var B string // To address
	var BalanceType string // Token type
	var err error

	B = strings.ToLower(args[0])
	BalanceType = args[1]

	// Amount
	amount := big.NewInt(0)
	_, good := amount.SetString(args[2], 10)
	if !good {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(B, BalanceType, amount)
	if err != nil {
		return shim.Error("transfer error" + err.Error())
	}
	return shim.Success(nil)
}

// counter
// Get current tx counter of an account
func (t *tokenChaincode) getCounter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Address
	var err error

	A = strings.ToLower(args[0])
	account, err := stub.GetAccount(A)
	if err != nil {
		jsonResp := "{\"Error\":\"account not exists\"}"
		return shim.Error(jsonResp)
	}

	if account == nil {
		jsonResp := "{\"Error\":\"account not exists for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"counter\":\"" + string(account.Counter) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(strconv.FormatUint(account.Counter, 10)))
}

func main() {
	err := shim.Start(new(tokenChaincode))
	if err != nil {
		fmt.Printf("Error starting tokenChaincode: %s", err)
	}
}
