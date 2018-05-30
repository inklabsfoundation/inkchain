package xscc

import (
	"encoding/json"
	"fmt"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"github.com/inklabsfoundation/inkchain/common/flogging"
	"github.com/inklabsfoundation/inkchain/core/policy"
	"github.com/inklabsfoundation/inkchain/core/policyprovider"
	"math/big"
	"strings"
	"github.com/inklabsfoundation/inkchain/core/wallet"
	"errors"
)

// Create Logger
var logger = flogging.MustGetLogger("xscc")

const (
	Unlock         = "unlock"         //public chain turn into
	Lock           = "lock"           //union chain turn out
	RegistPlatform = "registPlatform" //register a platform
	RemovePlatform = "removePlatform" //remove a platform
	QueryTxInfo    = "queryTxInfo"    //query transaction info
	QuerySignature = "querySignature" //query transaction signature
)

//turn out state struct
type turnOutMessage struct {
	FromAccount string   `json:"fromAccount"`
	BalanceType string   `json:"balanceType"`
	Value       *big.Int `json:"value"`
	ToPlatform  string   `json:"toPlatform"`
	ToAccount   string   `json:"toAccount"`
	DateTime    string   `json:"dateTime"`
}

//turn in state struct
type turnInMessage struct {
	TxId         string   `json:"txId"`
	Value        *big.Int `json:"value"`
	FromAccount  string   `json:"fromAccount"`
	FromPlatform string   `json:"fromPlatForm"`
	ToAccount    string   `json:"toAccount"`
	DateTime     string   `json:"dateTime"`
}

type CrossTrainSysCC struct {
	owner        string            //chain code owner
	platName     string            //platform name
	tokenAddress map[string]string //coin account
	// policyChecker is the interface used to perform
	// access control
	policyChecker policy.PolicyChecker
}

//init chain code
func (c *CrossTrainSysCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("xscc init")
	// Init policy checker for access control
	c.policyChecker = policyprovider.GetPolicyChecker()
	c.owner = wallet.CrossChainManager
	c.platName = wallet.LocalPlatform
	c.tokenAddress = wallet.TokenAddress
	if c.owner == "" || c.platName == "" || len(c.tokenAddress) <= 0 {
		return shim.Error("init arg error")
	}
	return shim.Success([]byte("init success"))
}

func (c *CrossTrainSysCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	logger.Debugf("xscc starts: %d args ", len(args))
	switch funcName {
	case RegistPlatform:
		return c.registPlatform(stub, args)
	case RemovePlatform:
		return c.removePlatform(stub, args)
	case Unlock:
		return c.unlock(stub, args)
	case Lock:
		return c.lock(stub, args)
	case QueryTxInfo:
		return c.queryTxInfo(stub, args)
	case QuerySignature:
		return c.querySignature(stub, args)
	}
	return shim.Success([]byte("Invalid invoke function name. Expecting \"RegistPlatform\" or \"RemovePlatform\" or \"Unlock\" or \"Lock\" or \"QueryTxInfo\" or \"QueryTxInfo\"."))
}

//register a platform
//args platform  string supportCross bool
func (c *CrossTrainSysCC) registPlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Params Error")
	}

	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	}
	if sender != c.owner {
		return shim.Error("Sender must be chainCode's owner")
	}

	platform := strings.ToLower(args[0])
	//try to get platform state from book which key is platform's value
	platState, err := stub.GetState(platform)
	if err != nil {
		return shim.Error("Failed to get platform: " + err.Error())
	} else if platState != nil {
		return shim.Error("This platform exists")
	}
	//make json data and write to book
	state, _ := json.Marshal(map[string]bool{platform: true})
	err = stub.PutState(platform, state)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("Operation Success"))
}

//remove one platform
//args platform string
func (c *CrossTrainSysCC) removePlatform(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Params Error")
	}

	//validate operator's permission
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	} else if sender != c.owner {
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

	return shim.Success([]byte("Operation Success"))
}

//public chain turn in
func (c *CrossTrainSysCC) unlock(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 5 {
		return shim.Error("Params Error")
	}

	fromPlatform := strings.ToLower(args[0])
	fromAccount := strings.ToLower(args[1])
	amount := big.NewInt(0)
	_, ok := amount.SetString(args[2], 10)
	toAccount := strings.ToLower(args[3])
	pubTxId := strings.ToLower(args[4])

	if !ok {
		return shim.Error("Expecting integer value for amount")
	} else if amount.Cmp(big.NewInt(0)) <= 0 {
		return shim.Error("Amount must be more than zero")
	}
	//try to get state from book which key is variable fromPlatform's value
	platState, err := stub.GetState(fromPlatform)
	if err != nil {
		return shim.Error("Failed to get platform: " + err.Error())
	} else if platState == nil {
		return shim.Error("The platform named " + fromPlatform + " is not registered")
	}

	//build state key
	key := fromPlatform + "|" + pubTxId
	//validate txId has not been processed
	xcMs, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	} else if xcMs != nil {
		return shim.Error("This transaction has been processed")
	}

	//do transfer  `wait to change`
	err = stub.CrossTransfer(toAccount, amount, pubTxId, fromPlatform)
	if err != nil {
		return shim.Error("transfer error " + err.Error())
	}

	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	timeStr := fmt.Sprintf("%d", txTimestamp.GetSeconds())
	//build turn in state and change to json
	state := c.buildTurnInMessage(stub.GetTxID(), fromAccount, fromPlatform, amount, toAccount, pubTxId, timeStr)
	err = stub.PutState(key, state)
	if err != nil {
		return shim.Error(err.Error())
	}

	//build composite key
	indexName := "type～address~datetime~platform~key"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{"in", toAccount, timeStr, fromPlatform, key})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)

	return shim.Success([]byte("Unlock Success"))
}

//union chain turn out
func (c *CrossTrainSysCC) lock(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//return shim.Error("Unlock function had been moved to another chaincode")
	if len(args) < 4 {
		return shim.Error("Params Error")
	}
	//get operator
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	} else if sender == "" {
		return shim.Error("Account not exist")
	}
	toPlatform := args[0]
	platformLower := strings.ToLower(toPlatform)
	toAccount := strings.ToLower(args[1])
	amount := big.NewInt(0)
	_, ok := amount.SetString(args[2], 10)
	balanceType := args[3]
	tokenAddr := c.tokenAddress[balanceType]
	if tokenAddr == "" {
		return shim.Error("Token address not found")
	}

	if !ok {
		return shim.Error("Expecting integer value for amount")
	} else if amount.Cmp(big.NewInt(0)) <= 0 {
		return shim.Error("Amount must more than zero")
	}

	//try to get state from book which key is variable toPlatform's value
	platState, err := stub.GetState(platformLower)
	if err != nil {
		return shim.Error("Failed to get platform: " + err.Error())
	} else if platState == nil {
		return shim.Error("The platform named " + platformLower + " is not registered")
	}

	//set txId to be key
	key := stub.GetTxID()
	//do transfer
	err = stub.Transfer(tokenAddr, balanceType, amount)
	if err != nil {
		return shim.Error("Transfer error " + err.Error())
	}
	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	timeStr := fmt.Sprintf("%d", txTimestamp.GetSeconds())
	//build turn out state
	state := c.buildTurnOutMessage(sender, balanceType, toPlatform, toAccount, amount, timeStr)
	err = stub.PutState(key, state)
	if err != nil {
		return shim.Error(err.Error())
	}
	//build composite key
	indexName := "type~address~datetime~platform~key"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{"out", sender, timeStr, c.platName, key})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)
	return shim.Success([]byte("Operate Success"))
}

//query transaction info
func (c *CrossTrainSysCC) queryTxInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

//query transaction signature
func (c *CrossTrainSysCC) querySignature(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Params error")
	}
	key := strings.ToLower(args[0])
	if len(key) == 0 {
		return shim.Error("Please input a right key")
	}
	if strings.Contains(key, "|") {
		return shim.Error("QuerySign do not support turn in transaction")
	}
	stateJson, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	} else if stateJson == nil {
		return shim.Error("Can't find state with named " + key)
	}
	state := turnOutMessage{}
	err = json.Unmarshal(stateJson, &state)
	if err != nil {
		return shim.Error(err.Error())
	}
	//sign
	str := fmt.Sprintf("%s:0x%s:%s:0x%s:%d:%s:%s", wallet.LocalPlatform, state.FromAccount[1:], state.ToPlatform, state.ToAccount, state.Value, state.BalanceType, key)
	sign, err := c.signJson([]byte(str), strings.ToLower(state.ToPlatform))
	if err != nil {
		return shim.Error(err.Error())
	}
	result := map[string]interface{}{"sign": sign, "state": state}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultJson)
}

//build turn in state and change to json
func (c *CrossTrainSysCC) buildTurnInMessage(txId string, fromAccount string, fromPlatform string, value *big.Int, toAccount string, pubTxId string, now string) []byte {
	state := turnInMessage{txId, value, fromAccount, fromPlatform, toAccount, now}
	stateJson, _ := json.Marshal(state)
	return stateJson
}

//build turn out state and change to json
func (c *CrossTrainSysCC) buildTurnOutMessage(fromAccount string, balanceType string, toPlatform string, toAccount string, value *big.Int, now string) []byte {
	state := turnOutMessage{fromAccount, balanceType, value, toPlatform, toAccount, now}
	stateJson, _ := json.Marshal(state)
	return stateJson
}

//sign
func (c *CrossTrainSysCC) signJson(balanceType string, str string, platform string) (string, error) {
	publicNode, ok := wallet.PublicInfos[platform]
	if !ok {
		return "", errors.New("platform not support")
	}
	privateKey := publicNode.PrivateKey
	if privateKey == "" {
		return "", errors.New("platform info not exist...")
	}
	contract, ok := publicNode.ContractList[balanceType]
	if !ok {
		return "", errors.New("balance type not support..")
	}
	str += contract.Version
	result, err := wallet.SignJson([]byte(str), privateKey)
	if result == nil {
		err = errors.New("signature failed")
		return "", err
	} else {
		return wallet.SignatureBytesToString(result), nil
	}
}
