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
	"math/rand"
)

type WorkChainCode struct {
	sysAddress		string
	TokenType		string		// as "INK"
	MaxWorkTime		int 		// max auto generate work time, sec
}

func main() {
	err := shim.Start(new(WorkChainCode))
	if err == nil {
		fmt.Printf("Error starting WorkChaincode: %s", err)
	}
}

func (w *WorkChainCode)Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 3 {
		return shim.Error("init, Incorrect number of arguments. Expecting 3")
	}
	w.sysAddress = args[0]
	w.TokenType = args[1]
	nTime, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	if nTime <= 0 {
		return shim.Error("init, auto generate work time invalid")
	}
	w.MaxWorkTime = nTime

	return shim.Success(nil)
}

func (w *WorkChainCode)Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case GetWork:
		if len(args) != 0 {
			return shim.Error("getWork, Incorrect number of arguments. Expecting 0")
		}
		return w.getWork(stub)
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
	}

	return shim.Error("function is invalid: " + function)
}

func (w *WorkChainCode) workState(stub shim.ChaincodeStubInterface, workId string) (work WorkDef, err error) {
	var valAsbytes []byte
	work = WorkDef{}
	valAsbytes, err = stub.GetState(workId)
	if err != nil {
		// Failed to get state
		return
	} else if valAsbytes == nil {
		// work does not exist
		fmt.Printf("%s work does not exist", workId)
		return
	}

	err = json.Unmarshal([]byte(valAsbytes), &work)
	return
}

func (w *WorkChainCode) getWork(stub shim.ChaincodeStubInterface) pb.Response {
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get sender's address.")
	}
	sender = strings.ToLower(sender)

	bGet, err := w.getFreeState(stub, sender)
	if err != nil {
		return shim.Error("Fail to getFreeState")
	}
	if bGet == true {
		return shim.Error("already get")
	}

	workId, level := w.autoWork(stub)
	if workId == "" {
		return shim.Error("getWork, auto generate work failed.")
	}

	work := &WorkDef {
		WorkId:		workId,
		Level: 		level,
		Birth:		time.Now().Unix(),
		Owner:		sender,
		Sale:		0,
		Price:		0,
		SaleTime:	0,
	}

	workJSONasByte, err := json.Marshal(work)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(workId, workJSONasByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = w.historyByComposite(stub, FreeHistory, []string{sender, workId})
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("get work:%v Succeed", workId)
	return shim.Success([]byte(result))
}

func (w *WorkChainCode) historyByComposite(stub shim.ChaincodeStubInterface, recordType string, args []string) error {
	var indexKey string
	var err error
	indexKey, err = stub.CreateCompositeKey(recordType, args)
	if err != nil {
		return err
	}
	value := []byte{0x00}
	err = stub.PutState(indexKey, value)
	if err != nil {
		return err
	}
	return nil
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
	work, err := w.workState(stub, workId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if 0 == work.Sale {
		jsonResp := "{\"Error\":\"buy failed, work does not sale: " + workId + "\"}"
		return shim.Error(jsonResp)
	}

	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get sender's address.")
	}
	sender = strings.ToLower(sender)
	if sender == work.Owner {
		jsonResp := "{\"Error\":\"buy failed, already is work owner\"}"
		return shim.Error(jsonResp)
	}
	balance := w.getBalance(stub, sender)
	if balance < int64(work.Price) {
		return shim.Error("bid failed, balance not enough")
	}

	// pay cost
	amount := big.NewInt(0)
	_, good := amount.SetString(strconv.Itoa(work.Price), 10)
	if !good  {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(work.Owner, w.TokenType, amount)
	if err != nil {
		return shim.Error("transfer error: " + err.Error())
	}

	work.Owner = sender
	work.Sale = 0
	work.SaleTime = 0
	workAsbytes, err := json.Marshal(work)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(workId, workAsbytes)
	if err != nil {
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
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get sender's address.")
	}
	sender = strings.ToLower(sender)

	work, err := w.workState(stub, workId)
	if err != nil {
		return shim.Error(err.Error())
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
	workJSONasByte, err := json.Marshal(work)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(workId, workJSONasByte)
	if err != nil {
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

func (w *WorkChainCode)timeDuration(srcTime int64) int {
	return int(time.Now().Unix() - srcTime)
}

func (w *WorkChainCode) autoWork(stub shim.ChaincodeStubInterface)(gene string, level int) {
	curTime := time.Now().Unix()
	for {
		if w.timeDuration(curTime) > w.MaxWorkTime {
			break
		}
		gene, level  = w.newGenome()
		// gene min length is 28
		if len(gene) < 28 {
			continue
		}
		catAsbytes, err := stub.GetState(gene)
		if err == nil && catAsbytes == nil {
			break
		}
	}

	return
}

func (w *WorkChainCode)newGenome() (string, int) {
	var buffer bytes.Buffer
	level := 0
	shape, index := w.rateRand(2,[]int{30,20,20,30})
	if shape == "" {
		return "", 0
	}
	if index == 1 || index == 2 {
		level += 2
	}
	buffer.WriteString(shape)

	line, index := w.rateRand(1,[]int{60,40})
	if line == "" {
		return "", 0
	}
	if index == 1 {
		level += 2
	}
	buffer.WriteString(line)
	shade, nColor := w.rateRand(1,[]int{60,30,10})
	if shade == "" {
		return "", 0
	}
	if nColor == 1 {
		level += 3
	}else if nColor == 2 {
		level += 5
	}
	buffer.WriteString(shape)
	//angle
	buffer.WriteString(w.limitRand(36000,5,1,35999))
	//draw
	buffer.WriteString(w.limitRand(10000,4,1,9999))
	//with
	buffer.WriteString(w.limitRand(900,3,1,899))
	for i:=0; i<=nColor; i++ {
		color , n := w.colorRand()
		level += n
		buffer.WriteString(color)
	}

	return buffer.String(), level
}

//rate total 100
//@param byte num
//@param per type rate
//@return bytes
//@return level
func (w *WorkChainCode)rateRand(length int, rateSlice []int) (string, int) {
	var buffer bytes.Buffer
	num := len(rateSlice)
	if num <=0 {
		return "", 0
	}
	val := rand.Intn(100)
	nCount := 0
	for i:=0; i<num; i++ {
		nCount += rateSlice[i]
		if val < nCount {
			valStr := strconv.Itoa(i)
			if len(valStr) < length {
				for j:=0; j<length-len(valStr); j++ {
					buffer.WriteString("0")
				}
			}
			buffer.WriteString(valStr)
			return buffer.String(), i
		}
	}

	return "", 0
}

func (w *WorkChainCode)limitRand(num, length, min, max int) string {
	if num <= 0 || length <= 0{
		return ""
	}
	if min >= max {
		return ""
	}

	var buffer bytes.Buffer
	val := 0
	for {
		val = rand.Intn(num)
		if val < min || val > max {
			continue
		}
		break
	}
	valStr := strconv.Itoa(val)
	if len(valStr) < length {
		for i:=0; i<length-len(valStr); i++ {
			buffer.WriteString("0")
		}
	}
	buffer.WriteString(valStr)
	return buffer.String()
}

func (w *WorkChainCode)colorRand() (string, int) {
	var buffer bytes.Buffer
	var level int
	r := w.limitRand(255,3,1,254)
	g := w.limitRand(255,3,1,254)
	b := w.limitRand(255,3,1,254)
	a := w.limitRand(101,3,1,100)
	if r == g && r == b {
		level += 2
	}else if (r != g && r != b && g != b) {
		level = 0
	}else {
		level += 1
	}
	buffer.WriteString(r)
	buffer.WriteString(g)
	buffer.WriteString(b)
	buffer.WriteString(a)
	return buffer.String(), level
}

