package main

import (
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"fmt"
	"encoding/json"
	"strings"
	"math/big"
	"strconv"
	"time"
	"bytes"
)

type WorkChainCode struct {
	sysAddress		string
	TokenType		string		// as "INK"
}

func main() {
	err := shim.Start(new(WorkChainCode))
	if err == nil {
		fmt.Printf("Error starting WorkChaincode: %s", err)
	}
}

func (w *WorkChainCode)Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 2 {
		return shim.Error("init, Incorrect number of arguments. Expecting 2")
	}
	w.sysAddress = args[0]
	w.TokenType = args[1]

	return shim.Success(nil)
}

func (w *WorkChainCode)Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case RegisterWork:
		if len(args) != 3 {
			return shim.Error("getWork, Incorrect number of arguments. Expecting 3")
		}
		return w.registerWork(stub, args)
	case Purchase:
		if len(args) != 1 {
			return shim.Error("purchase, Incorrect number of arguments. Expecting 1")
		}
		return w.purchase(stub, args)
	case Sell:
		if len(args) != 3 {
			return shim.Error("sell, Incorrect number of arguments. Expecting 3")
		}
		return w.sell(stub, args)
	case Query:
		if len(args) != 1 {
			return shim.Error("query, Incorrect number of arguments. Expecting 1")
		}
		return w.query(stub, args)
	case QueryInkwork:
		if len(args) != 1 {
			return shim.Error("queryInkwork, Incorrect number of arguments. Expecting 1")
		}
		return w.queryInkwork(stub, args)
	}

	return shim.Error("function is invalid: " + function)
}

func (w *WorkChainCode) workState(stub shim.ChaincodeStubInterface, workId string) (work *WorkDef, err error) {
	var valAsbytes []byte
	if valAsbytes, err = stub.GetState(workId); err != nil {
		// Failed to get state
		return nil, err
	} else if valAsbytes == nil {
		// work does not exist
		fmt.Printf("%s work does not exist", workId)
		return nil, nil
	}
	work = &WorkDef{}
	err = json.Unmarshal([]byte(valAsbytes), work)
	return
}

//@param workId
//@param level
//@param price
func (w *WorkChainCode) registerWork(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	workId := args[0]
	var level, price int
	var err error
	if level, err = strconv.Atoi(args[1]); err != nil{
		return shim.Error(err.Error())
	}
	if price, err = strconv.Atoi(args[2]); err != nil {
		return shim.Error(err.Error())
	}

	var sender string
	if sender, err = stub.GetSender(); err != nil {
		return shim.Error("Fail to get sender's address.")
	}
	sender = strings.ToLower(sender)

	if bGet, err := w.getFreeState(stub, sender); err != nil {
		return shim.Error("Fail to getFreeState")
	} else if bGet == true {
		return shim.Error("already get")
	}

	if data, err := w.workState(stub, workId); err != nil {
		return shim.Error(err.Error())
	} else if data != nil {
		return shim.Error("this work already exist")
	}

	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("GetTxTimestamp failed")
	}
	work := &WorkDef {
		WorkId:		workId,
		Level: 		level,
		Birth:		txTimestamp.String(),
		Owner:		sender,
		Sale:		0,
		Price:		price,
		SaleTime:	0,
	}

	if workJSONasByte, err := json.Marshal(work); err == nil {
		if err = stub.PutState(workId, workJSONasByte); err != nil {
			return shim.Error(err.Error())
		}
	} else {
		return shim.Error(err.Error())
	}

	err = w.historyByComposite(stub, FreeHistory, []string{sender, workId})
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("register work:%v Succeed", workId)
	return shim.Success([]byte(result))
}

func (w *WorkChainCode) historyByComposite(stub shim.ChaincodeStubInterface, recordType string, args []string) (err error) {
	var indexKey string
	if indexKey, err = stub.CreateCompositeKey(recordType, args); err != nil {
		return
	}

	value := []byte{0x00}
	if err = stub.PutState(indexKey, value); err != nil {
		return
	}
	return
}

func (w *WorkChainCode) getFreeState(stub shim.ChaincodeStubInterface, user string) (bool, error) {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(FreeHistory, []string{})
	if err != nil {
		fmt.Println("getFreeState, stub.GetStateByPartialCompositeKey failed: " + err.Error())
		return false, err
	}
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return false, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return false, err
		}

		userAddr := compositeKeyParts[0]
		if userAddr == user {
			return true, nil
		}
	}
	return false, nil
}

//@param workId
func (w *WorkChainCode) purchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	workId := args[0]
	var data *WorkDef
	var err error
	if data, err = w.workState(stub, workId); err != nil {
		return shim.Error(err.Error())
	} else if data == nil {
		return shim.Error("this work not exist")
	}

	if 0 == data.Sale {
		jsonResp := "{\"Error\":\"buy failed, work does not sale: " + workId + "\"}"
		return shim.Error(jsonResp)
	}

	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get sender's address.")
	}
	sender = strings.ToLower(sender)
	if sender == data.Owner {
		jsonResp := "{\"Error\":\"buy failed, already is work owner\"}"
		return shim.Error(jsonResp)
	}
	balance := w.getBalance(stub, sender)
	if balance < int64(data.Price) {
		return shim.Error("bid failed, balance not enough")
	}

	// pay cost
	amount := big.NewInt(0)
	_, good := amount.SetString(strconv.Itoa(data.Price), 10)
	if !good  {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(data.Owner, w.TokenType, amount)
	if err != nil {
		return shim.Error("transfer error: " + err.Error())
	}

	data.Owner = sender
	data.Sale = 0
	data.SaleTime = 0
	if workAsbytes, err := json.Marshal(data); err != nil {
		return shim.Error(err.Error())
	} else if err = stub.PutState(workId, workAsbytes); err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("buy Work:%v Succeed", workId)
	return shim.Success([]byte(result))
}

func (w *WorkChainCode) getBalance(stub shim.ChaincodeStubInterface, userAddr string) int64 {
	account, err := stub.GetAccount(userAddr)
	if err != nil {
		// account not exists
		return 0
	}
	if account == nil || account.Balance[w.TokenType] == nil {
		return 0
	}

	balance := account.Balance[w.TokenType].Int64()
	return balance
}

//@param workId
//@param price
//@param sellState
func (w *WorkChainCode) sell(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	workId := args[0]
	var work *WorkDef
	var err error
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get sender's address.")
	}
	sender = strings.ToLower(sender)

	if work, err = w.workState(stub, workId); err != nil {
		return shim.Error(err.Error())
	} else if work == nil {
		return shim.Error("this work not exist")
	}

	if sender != work.Owner {
		jsonResp := "{\"Error\":\"sell failed, not work owner\"}"
		return shim.Error(jsonResp)
	}

	priceStr := args[1]
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if 1 == work.Sale && price == work.Price  {
		jsonResp := "{\"Error\":\"invalid operate, nothing change\"}"
		return shim.Error(jsonResp)
	}

	sale, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	var saleTime int64 = 0
	if 1 == sale {
		saleTime = time.Now().Unix()
	}
	work.Sale = sale
	work.Price = price
	work.SaleTime = saleTime
	if workJSONasByte, err := json.Marshal(work); err != nil {
		return shim.Error(err.Error())
	} else if err = stub.PutState(workId, workJSONasByte); err != nil {
			return shim.Error(err.Error())
	}

	result := fmt.Sprintf("sell work:%v, sale:%v, price:%v", workId, sale, price)
	return shim.Success([]byte(result))
}

//@param query_type
func (w *WorkChainCode)query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryType := args[0]
	nType, err := strconv.Atoi(queryType)
	if err != nil {
		return shim.Error(err.Error())
	}
	if nType <= QueryStart || nType >= QueryEnd {
		return shim.Error("{\"Error\":\"invalid query\"}")
	}

	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get sender's address.")
	}
	sender = strings.ToLower(sender)

	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	bArrayIndex := 1
	work := WorkDef{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(queryResponse.Value, &work)
		if err != nil {
			continue
		}
		switch nType {
		case All:
		case Sale:
			if 0 == work.Sale {
				continue
			}
		case Self:
			if sender != work.Owner {
				continue
			}
		default:
			return shim.Error("{\"Error\":\"invalid query\"}")
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Number\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(bArrayIndex))
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
		bArrayIndex ++
	}
	buffer.WriteString("]")

	return shim.Success([]byte(buffer.String()))
}

//@param workId
func (w *WorkChainCode) queryInkwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	workId := args[0]
	valAsbytes, err := stub.GetState(workId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + workId + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"cat does not exist: " + workId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}
