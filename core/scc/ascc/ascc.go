/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ascc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/inklabsfoundation/inkchain/common/flogging"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	"github.com/inklabsfoundation/inkchain/core/policy"
	"github.com/inklabsfoundation/inkchain/core/policyprovider"
	"github.com/inklabsfoundation/inkchain/msp/mgmt"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
)

// Create Logger
var tralogger = flogging.MustGetLogger("ascc")

// These are function names from Invoke first parameter
const (
	//invoke functions
	RegisterAndIssueToken string = "registerAndIssueToken"
	InvalidateToken       string = "invalidateToken"

	//token status
	Created    string = "created"
	Delivered  string = "issued"
	Invalidate string = "invalidated"
)

// type GenAccount
type Token struct {
	// token name
	Name string `json:"tokenName"`
	// total supply of the token
	totalSupply *big.Int `json:"totalSupply"`
	// initial address to issue
	Address string `json:"address"`
	// token status : Created, Delivered, Invalidate
	Status string `json:"status"`
	// token decimals
	Decimals int `json:"decimals"`
}

//-------------- the ascc ------------------
type AssetSysCC struct {
	// policyChecker is the interface used to perform
	// access control
	policyChecker policy.PolicyChecker
}

// Init initializes ascc
func (t *AssetSysCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	tralogger.Info("Init ascc")

	// Init policy checker for access control
	t.policyChecker = policyprovider.GetPolicyChecker()

	return shim.Success(nil)
}

// Invoke func
func (t *AssetSysCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	tralogger.Debugf("ascc starts: %d args", len(args))
	// Handle ACL of ascc:
	// 1. get the signed proposal
	//	sp, err := stub.GetSignedProposal()
	//	if err != nil {
	//		return shim.Error(fmt.Sprintf("Failed retrieving signed proposal on executing %s with error %s", function, err))
	//	}
	switch function {

	case RegisterAndIssueToken:
		if len(args) != 4 { //name, totalSupply, decimals, address
			returnMessage := fmt.Sprint("Incorrect number of arguments for IssueToken, %d", len(args))
			tralogger.Debugf("Incorrect number of arguments for IssueToken, %d", len(args))
			return shim.Error(returnMessage)
		}

		// 2. check local MSP Admins policy

		sp, err := stub.GetSignedProposal()
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed retrieving signed proposal on executing %s with error %s", function, err))
		}

		if err := t.policyChecker.CheckPolicyNoChannel(mgmt.Admins, sp); err != nil {
			tralogger.Debugf(fmt.Sprintf("Authorization for RegisterAndIssueToken has been denied (error-%s)", err))
			return shim.Error(fmt.Sprintf("Authorization for RegisterAndIssueToken has been denied (error-%s)", err))
		}
		tralogger.Debugf("Invoke function: %s", RegisterAndIssueToken)
		return t.registerAndIssueToken(stub, args)

	case InvalidateToken:
		if len(args) != 1 { //name
			returnMessage := fmt.Sprint("Incorrect number of arguments for InvalidateToken, %d", len(args))
			tralogger.Debugf("Incorrect number of arguments for InvalidateToken, %d", len(args))
			return shim.Error(returnMessage)
		}

		tralogger.Debugf("Invoke function: %s", InvalidateToken)

		return t.invalidateToken(stub, args)

	}

	return shim.Error("Invalid invoke function name. Expecting \"registerAndIssueToken\" or \"invalidateToken\".")
}

// issue Tokens Invoke
// invoke function
func (t *AssetSysCC) registerAndIssueToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	tokenName := args[0]

	totalSupply := big.NewInt(0)
	_, good := totalSupply.SetString(args[1], 10)
	if !good {
		return shim.Error("Expecting integer value for totalSupply.")
	}
	dec, _ := strconv.Atoi(args[2])
	addr := args[3]

	//Get exist token
	var existToken Token
	existTokenBytes, err := stub.GetState(tokenName)
	if err != nil {
		msgCheck := "Check token existance error, fail to getState of "
		msgCheck += tokenName
		tralogger.Debug(msgCheck)
		return shim.Error(msgCheck)
	}

	//Get the information of token
	//If not exist, create a new token first
	if existTokenBytes == nil {
		//not exist
		//create the token
		existToken.Status = Created
		existToken.Name = tokenName
		existToken.totalSupply = totalSupply
		existToken.Address = addr
		existToken.Decimals = dec
	} else {
		//exist
		//unmarshal the jsonBytes & check token information
		err = json.Unmarshal(existTokenBytes, &existToken)
		if err != nil {
			msgUnmarshal := "Unmarshal exist tokenBytes err "
			msgUnmarshal += tokenName
			tralogger.Debug(msgUnmarshal)
			return shim.Error(msgUnmarshal)
		}
		//check the status of token
		if existToken.Status != Created {
			msgCheckTS := "Token status err, fail to issue token."
			tralogger.Debug(msgCheckTS)
			return shim.Error(msgCheckTS)
		}
		//check the information of token
		if existToken.Address != addr || existToken.totalSupply.Cmp(totalSupply) != 0 || existToken.Decimals != dec {
			msgCheckTInfo := "Token info err, check fialed."
			tralogger.Debug(msgCheckTInfo)
			return shim.Error(msgCheckTInfo)
		}
	}

	//set the token number to address

	//get account of Address
	account, err := stub.GetAccount(addr)
	//check if token has been issued before
	if err == nil {
		if account != nil {
			if _, ok := account.Balance[tokenName]; ok {
				msgBalanceCheck := "Token " + tokenName + " already exist in " + addr
				tralogger.Debug(msgBalanceCheck)
				return shim.Error(msgBalanceCheck)
			}
		}
	}
	//token hasnot been issued, then
	//issue token
	err = stub.IssueToken(addr, tokenName, totalSupply)
	if err != nil {
		return shim.Error(err.Error())
	}

	existToken.Status = Delivered

	//store the latest status for token in ascc
	existTokenJson, err := json.Marshal(&existToken)
	err = stub.PutState(tokenName, existTokenJson)

	if err != nil {
		msgUpdate := "Store the latest token status err."
		tralogger.Debug(msgUpdate)
		return shim.Error(msgUpdate)
	}

	return shim.Success([]byte("Token issued success!"))
}

// invalidate Tokens Invoke
// invoke function
func (t *AssetSysCC) invalidateToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	tokenName := args[0]

	//Get exist token
	var existToken Token
	existTokenBytes, err := stub.GetState(tokenName)
	if err != nil {
		msgCheck := "Check token existance error, fail to getState of "
		msgCheck += tokenName
		tralogger.Debug(msgCheck)
		return shim.Error(msgCheck)
	}

	//Check token status
	//If not exist, return err
	if existTokenBytes == nil {
		msgCheckExist := "Token not exist, fail to get status: "
		msgCheckExist += tokenName
		tralogger.Debug(msgCheckExist)
		return shim.Error(msgCheckExist)
	}

	//Unmarshal tokenBytes
	err = json.Unmarshal(existTokenBytes, &existToken)
	if err != nil {
		msgUnmarshal := "Unmarshal exist tokenBytes err "
		msgUnmarshal += tokenName
		tralogger.Debug(msgUnmarshal)
		return shim.Error(msgUnmarshal)
	}

	//check the status of token
	if existToken.Status == Invalidate {
		return shim.Error("Token already invalidated.")
	}

	existToken.Status = Invalidate

	//store the latest status for token
	existTokenJson, err := json.Marshal(&existToken)
	err = stub.PutState(tokenName, existTokenJson)

	if err != nil {
		msgUpdate := "Store the latest token status err."
		tralogger.Debug(msgUpdate)
		return shim.Error(msgUpdate)
	}

	return shim.Success([]byte("Token invalidate success!"))
}
