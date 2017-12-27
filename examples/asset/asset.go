/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"strconv"
	"encoding/json"
	"strings"
)

const (
	// invoke func name
	AddUser				string = "addUser"
	QueryUser			string = "queryUser"
	AddAsset			string = "addAsset"
	ReadAsset			string = "readAsset"
	DeleteAsset			string = "deleteAsset"
	TransferAsset		string = "transferAsset"
	QueryAssetsByOwner	string = "queryAssetsByOwner"
	GetHistoryForAsset	string = "getHistoryForAsset"
)

// Prefixes for user and asset separately
const (
	UserPrefix	= "USER_"
	AssetPrefix	= "ASSET_"
)

// Demo chaincode for asset registering, querying and transferring
type assetChaincode struct {
}

type user struct {
	Name	string `json:"name"`
	Age		int	   `json:"age"`
	Address string `json:"address"`	// the address actually decides a user
}

type asset struct {
	Name 	string `json:"name"`
	Type 	string `json:"type"`
	Content	string `json:"content"`
	Owner 	string `json:"owner"`	// store the name of the asset here
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(assetChaincode))
	if err != nil {
		fmt.Printf("Error starting assetChaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *assetChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("assetChaincode Init.")
	return shim.Success([]byte("Init success."))
}

// Invoke func
func (t *assetChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("assetChaincode Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case AddUser:
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2.")
		}
		// args[0]: user name
		// args[1]: user age
		// user address could be revealed from private key provided when invoking
		return t.addUser(stub, args)

	case QueryUser:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: user name
		return t.queryUser(stub, args)

	case AddAsset:
		if len(args) != 4 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: asset name
		// args[1]: type
		// args[2]: content
		// args[3]: owner
		return t.addAsset(stub, args)

	case ReadAsset:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: asset name
		return t.readAsset(stub, args)

	case DeleteAsset:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: asset name
		return t.delAsset(stub, args)
		// TODO:
	//case TransferAsset:	
	//case QueryAssetsByOwner:
	//case GetHistoryForAsset:
	}

	return shim.Error("Invalid invoke function name.")
}

// addUser: Register a new user
func (t *assetChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var new_name string
	var new_age int
	var err error

	new_name = args[0]
	new_age, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for user's age.")
	}

	// get user's address
	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)

	// check if user exists
	user_key := UserPrefix + new_name
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This user already exists: " + new_name)
		return shim.Error("This user already exists: " + new_name)
	}

	// register user
	user := &user{new_name, new_age, new_add}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(user_key, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("User register success."))
}

// queryUser: query the informatioin of a user
func (t *assetChaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	user_name := args[0]
	user_key := UserPrefix + user_name

	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + user_name)
		return shim.Error("This user doesn't exist: " + user_name)
	}

	return shim.Success(userAsBytes)
}

// addasset: add a new asset
func (t *assetChaincode) addAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	asset_name := args[0]
	asset_key := AssetPrefix + asset_name

	asset_type := args[1]
	asset_content := args[2]
	owner_name := args[3]

	// verify weather the owner exists
	owner_key := UserPrefix + owner_name
	userAsBytes, err := stub.GetState(owner_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This owner doesn't exist: " + owner_name)
		return shim.Error("This owner doesn't exist: " + owner_name)
	}

	// register asset
	asset := &asset{asset_name, asset_type,asset_content, owner_name}
	assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(asset_key, assetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("asset register success."))
}

// readassetcut: query the informatioin of a user
func (t *assetChaincode) readAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	asset_name := args[0]
	asset_key := AssetPrefix + asset_name

	assetAsBytes, err := stub.GetState(asset_key)
	if err != nil {
		return shim.Error("Fail to get asset: " + err.Error())
	}
	if assetAsBytes == nil {
		fmt.Println("This asset doesn't exist: " + asset_name)
		return shim.Error("This asset doesn't exist: " + asset_name)
	}

	return shim.Success(assetAsBytes)
}

// delassetcut: del a specific asset by name
func (t *assetChaincode) delAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	asset_name := args[0]
	asset_key := AssetPrefix + asset_name


	fmt.Println( " - delete asset begin. ")

	// step 1: get the asset info
	assetAsBytes, err := stub.GetState(asset_key)
	if err != nil {
		return shim.Error("Fail to get asset: " + err.Error())
	}
	if assetAsBytes == nil {
		fmt.Println("This asset doesn't exist: " + asset_name)
		return shim.Error("This asset doesn't exist: " + asset_name)
	}

	var assetJSON asset

	err = json.Unmarshal([]byte(assetAsBytes), &assetJSON)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + asset_name + "\"}"
		return shim.Error(jsonResp)
	}

	// test output
	fmt.Println( " - get asset info: name: " + assetJSON.Name + "; owner:" + assetJSON.Owner)

	// step 2: get the owner's address
	user_name := assetJSON.Owner
	user_key := UserPrefix + user_name
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + asset_name)
		return shim.Error("This user doesn't exist: " + asset_name)
	}

	var userJSON user

	err = json.Unmarshal([]byte(userAsBytes), &userJSON)
	owner_add := userJSON.Address
	// test output
	fmt.Println( " - get owner address: " + owner_add)

	// step 3: compare sender's address and owner's

	sender_add, err := stub.GetSender()
	// test output
	fmt.Println( " - get sender address: " + sender_add)


	// step 4: check address and then delete the asset

	if owner_add != sender_add {
		fmt.Println("Authorization denied. ")
		return shim.Error("Authorization denied. ")
	}

	err = stub.DelState(asset_key)
	if err != nil {
		fmt.Println("Fail to delete: " + asset_name)
		return shim.Error("Fail to delete" + asset_name)
	}

	return shim.Success([]byte("asset delete success."))
}
