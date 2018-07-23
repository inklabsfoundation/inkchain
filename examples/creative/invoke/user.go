package invoke

import (
	. "github.com/inklabsfoundation/inkchain/examples/creative/conf"
	. "github.com/inklabsfoundation/inkchain/examples/creative/util"
	. "github.com/inklabsfoundation/inkchain/examples/creative/model"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"fmt"
	"encoding/json"
	"strings"
)

type UserInvoke struct{}

func (*UserInvoke) AddUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("add user start.")
	username := args[0]
	email := args[1]
	//email = "mm@mm.com"
	// The user parameters will be extended. TODO
	user_key := GetUserKey(username)
	// get user's address
	address, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	address = strings.ToLower(address)
	// check if user exists
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This user already exists: " + username)
		return shim.Error("This user already exists: " + username)
	}
	// add user
	user := &User{Username: username, Email: email, Address: address}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(user_key, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("add user success."))
}

func (*UserInvoke) DeleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("delete user start.")
	username := args[0]
	user_key := GetUserKey(username)
	// step 1: get the user info
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}
	var userJSON User
	err = json.Unmarshal([]byte(userAsBytes), &userJSON)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + username + "\"}"
		return shim.Error(jsonResp)
	}
	// step 2: get user's address
	address, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	address = strings.ToLower(address)
	if userJSON.Address != address {
		return shim.Error("The sender's address doesn't correspond with the user's.")
	}
	// step 3: delete user's info
	err = stub.DelState(user_key)
	if err != nil {
		fmt.Println("Fail to delete: " + user_key)
		return shim.Error("Fail to delete" + user_key)
	}
	return shim.Success([]byte("delete user success."))
}

func (*UserInvoke) ModifyUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("modify user start.")
	username := args[0]
	user_key := GetUserKey(username)
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}
	var userJSON User
	err = json.Unmarshal([]byte(userAsBytes), &userJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = GetModifyUser(&userJSON, args[1:])
	if err != nil {
		return shim.Error(err.Error())
	}
	userJSONasBytes, err := json.Marshal(userJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(user_key, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(userJSONasBytes)
}

func (*UserInvoke) QueryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("query user start.")
	username := args[0]
	user_key := GetUserKey(username)
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}
	return shim.Success(userAsBytes)
}

func (*UserInvoke) ListOfUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("list of user start.")
	resultsIterator, err := stub.GetStateByRange(UserPrefix+StateStartSymbol, UserPrefix+StateEndSymbol)
	if err != nil {
		return shim.Error(err.Error())
	}
	list, err := GetListResult(resultsIterator)
	if err != nil {
		return shim.Error("getListResult failed")
	}
	return shim.Success(list)
}
