package main

import (
	"encoding/json"
	"fmt"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"math/big"
	"strings"
)

const (
	Unlock         = "unlock"         //public chain turn into
	Lock           = "lock"           //union chain turn out
	RegistPlatform = "registPlatform" //register a platform
	RemovePlatform = "removePlatform" //remove a platform
	QueryTxInfo    = "queryTxInfo"    //query transaction info
)

//turn out state struct
type turnOutMessage struct {
	FromUser   string   `json:"fromUser"`
	Value      *big.Int `json:"value"`
	ToPlatform string   `json:"toPlatform"`
	ToUser     string   `json:"toUser"`
	DateTime   string   `json:"dateTime"`
}

//turn in state struct
type turnInMessage struct {
	TxId         string   `json:"txId"`
	Value        *big.Int `json:"value"`
	FromUser     string   `json:"fromUser"`
	FromPlatform string   `json:"fromPlatForm"`
	ToUser       string   `json:"toUser"`
	DateTime     string   `json:"dateTime"`
}

//add platform event
type platformEvent struct {
	PlatName string `json:"platName"` //platform name
	Symbol   string `json:"symbol"`   //sign `+` add `-`remove
}

type XcChaincode struct {
	owner        string //chain code owner
	platName     string //platform name
	inkTokenAddr string //coin account
}

//init chain code
func (x *XcChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 2 {
		return shim.Error("Params Error")
	}
	x.owner = "4230a12f5b0693dd88bb35c79d7e56a68614b199"
	if len(x.owner) <= 0 || x.owner == "" {
		return shim.Error("Please input the right inkToken owner address")
	}
	x.platName = args[0]
	if x.platName == "" || x.platName == "nil" {
		return shim.Error("Please the right plat name")
	}
	x.inkTokenAddr = args[1]
	if x.inkTokenAddr == "" || len(x.inkTokenAddr) <= 0 {
		return shim.Error("Please input the right inkToken owner address")
	}
	return shim.Success([]byte("init success"))
}

func (x *XcChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case RegistPlatform:
		return x.registPlatform(stub, args)
	case RemovePlatform:
		return x.removePlatform(stub, args)
	case Unlock:
		return x.unlock(stub, args)
	case Lock:
		return x.lock(stub, args)
	case QueryTxInfo:
		return x.queryTxInfo(stub, args)
	}
	return shim.Success([]byte("invoke"))
}

//register a platform
//args platform  string supportCross bool
func (x *XcChaincode) registPlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Params Error")
	}

	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	}
	if sender != x.owner {
		return shim.Error("Sender must be chainCode's owner")
	}

	platform := strings.ToLower(args[0])
	//try to get platform state from book which key is platform's value
	platState, err := stub.GetState(platform)
	if err != nil {
		return shim.Error("Failed to get platform: " + err.Error())
	} else if platState != nil {
		return shim.Error("This platform existed")
	}
	//make json data and write to book
	state, _ := json.Marshal(map[string]bool{platform: true})
	err = stub.PutState(platform, state)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("platformEvent", x.buildPlatformEventMessage(platform, "+"))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("Operate Success"))
}

//remove one platform
//args platform string
func (x *XcChaincode) removePlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Params Error")
	}

	//validate operator's permission
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	} else if sender != x.owner {
		return shim.Error("Sender must be chainCode's owner")
	}

	platform := strings.ToLower(args[0])
	//try to get platform state from book which key is platform's value
	platState, err := stub.GetState(platform)
	if err != nil {
		return shim.Error("Failed to get platform: " + err.Error())
	} else if platState == nil {
		return shim.Error("This platform not existed")
	}
	//do remove
	err = stub.DelState(platform)
	if err != nil {
		return shim.Error("Failed to delete platform:" + err.Error())
	}
	//trigger event
	err = stub.SetEvent("platformEvent", x.buildPlatformEventMessage(platform, "-"))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("Operate Success"))
}

//public chain turn in
func (x *XcChaincode) unlock(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 5 {
		return shim.Error("Params Error")
	}

	fromPlatform := strings.ToLower(args[0])
	amount := big.NewInt(0)
	toUser := strings.ToLower(args[2])
	_, ok := amount.SetString(args[1], 10)
	pubTxId := strings.ToLower(args[3])

	if !ok {
		return shim.Error("Expecting integer value for amount")
	}
	//try to get state from book which key is variable fromPlatform's value
	platState, err := stub.GetState(fromPlatform)
	if err != nil {
		return shim.Error("Failed to get platform: " + err.Error())
	} else if platState == nil {
		return shim.Error("The platform named " + fromPlatform + " is not registered")
	}

	//build state key
	key := fromPlatform + pubTxId
	//validate txId has not been processed
	xcMs, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	} else if xcMs != nil {
		return shim.Error("This transaction has been processed")
	}

	//do transfer  `wait to change`
	//@todo function change to the other function that used to transfer from token address to toUser
	err = stub.Transfer(toUser, "INK", amount)
	if err != nil {
		return shim.Error("transfer error " + err.Error())
	}

	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	//build turn in state and change to json
	state := x.buildTurnInMessage(stub.GetTxID(), "", fromPlatform, amount, toUser, pubTxId, txTimestamp.String())
	err = stub.PutState(key, state)
	if err != nil {
		return shim.Error(err.Error())
	}

	//build composite key
	indexName := "type～address~platform~key"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{"in", toUser, fromPlatform, key})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)

	return shim.Success([]byte("unlockSuccess"))
}

//union chain turn out
func (x *XcChaincode) lock(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error("Params Error")
	}
	//get operator
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	} else if sender == "" {
		return shim.Error("Account not exist")
	}
	toPlatform := strings.ToLower(args[0])
	toUser := strings.ToLower(args[1])
	amount := big.NewInt(0)
	_, ok := amount.SetString(args[2], 10)
	if !ok {
		return shim.Error("Expecting integer value for amount")
	}

	//try to get state from book which key is variable toPlatform's value
	platState, err := stub.GetState(toPlatform)
	if err != nil {
		return shim.Error("Failed to get platform: " + err.Error())
	} else if platState == nil {
		return shim.Error("The platform named " + toPlatform + " is not registered")
	}

	//set txId to be key
	key := stub.GetTxID()
	//do transfer
	err = stub.Transfer(x.inkTokenAddr, "INK", amount)
	if err != nil {
		return shim.Error("Transfer error " + err.Error())
	}
	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	//build turn out state
	state := x.buildTurnOutMessage(key, sender, toPlatform, toUser, amount, txTimestamp.String())
	err = stub.PutState(key, state)
	if err != nil {
		return shim.Error(err.Error())
	}
	//build composite key
	indexName := "type~address~datetime~platform~key"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{"out", sender, txTimestamp.String(), x.platName, key})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)

	//sign
	signJson, err := x.signJson([]byte("abc"), "60320b8a71bc314404ef7d194ad8cac0bee1e331")
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(signJson)
}

//query transaction info
func (x *XcChaincode) queryTxInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Params error")
	}
	key := strings.ToLower(args[0])
	if len(key) == 0 {
		return shim.Error("Please input a right key")
	}
	stateJson, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	} else if stateJson == nil {
		return shim.Error("Can't find state with named " + key)
	}
	return shim.Success(stateJson)
}

//build platform change event
func (x *XcChaincode) buildPlatformEventMessage(platform string, symbol string) []byte {
	msg := platformEvent{platform, symbol}
	msgJson, _ := json.Marshal(msg)
	return msgJson
}

//build turn in state and change to json
func (x *XcChaincode) buildTurnInMessage(txId string, fromUser string, fromPlatform string, value *big.Int, toUser string, pubTxId string, now string) []byte {
	state := turnInMessage{txId, value, fromUser, fromPlatform, toUser, now}
	stateJson, _ := json.Marshal(state)
	return stateJson
}

//build turn out state and change to json
func (x *XcChaincode) buildTurnOutMessage( fromUser string, toPlatform string, toUser string, value *big.Int, now string) []byte {
	state := turnOutMessage{fromUser, value, toPlatform, toUser, now}
	stateJson, _ := json.Marshal(state)
	return stateJson
}

func (x *XcChaincode) signJson(json []byte, priKey string) ([]byte, error) {
	return []byte("f4128988cbe7df8315440adde412a8955f7f5ff9a5468a791433727f82717a6753bd71882079522207060b681fbd3f5623ee7ed66e33fc8e581f442acbcf6ab800"), nil
}

func main() {
	err := shim.Start(new(XcChaincode))
	if err != nil {
		fmt.Printf("Error starting tokenChaincode: %s", err)
	}
}
