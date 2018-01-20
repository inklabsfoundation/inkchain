/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.
*/

// cat: a demo chaincode for inkchain

// ====CHAINCODE EXECUTION SAMPLES (CLI) =================

// ==== Instantiate chaincode ====
// peer chaincode instantiate -o orderer.example.com:7050 -c '{"Args":["init","5","5","6","60","INK","07caf88941eafcaaa3370657fccc261acb75dfba"]}'

// ==== Invoke cat ====
// peer chaincode invoke -C mychannel -n cat -c '{"Args":["initSystemCat"]}'
// peer chaincode invoke -C mychannel -n cat -c '{"Args":["delete","7918"]}'
// peer chaincode invoke -C mychannel -n cat -c '{"Args":["setState","7918","0","1"]}'
// peer chaincode invoke -C mychannel -n cat -c '{"Args":["breed","5060","7918"]}'
// peer chaincode invoke -C mychannel -n cat -c '{"Args":["buy","7918"]}'

// ==== Query cat ====
// peer chaincode query -C mychannel -n cat -c '{"Args":["query","7918"]}'
// peer chaincode query -C mychannel -n cat -c '{"Args":["queryAll"]}'
// peer chaincode query -C mychannel -n cat -c '{"Args":["querySale"]}'

package main

import (
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"fmt"
	"encoding/json"
	"bytes"
	"math/rand"
	"strconv"
	"strings"
	"math/big"
	"time"
	"sort"
)

const (
	// function
	InitSystemCat  	string = "initSystemCat"
	QueryCat 		string = "query"
	QueryAll		string = "queryAll"
	QuerySale		string = "querySale"
	DelCat 	 		string = "delete"
	BreedCat 		string = "breed"
	BuyCat			string = "buy"
	SetState		string = "setState"
	// config
	HistoryIndex 	string = "buyhistory"
)

type CatChainCode struct {
	SysCatTime 		int			// system generate cat max time, if timeout return failed
	LatestBuyPrice 	int			// count latest buy record number
	InitPrice 		int			// init cat price
	NextMateTime 	int			// cat can mate time since last mate
	TokenType 		string		// tokenType as "INK"
	SystemAccount	string		// system account, user buy gen0 pay money to system
}

// Gene type：per type 1~9
const (
	Color 	= iota
	Eye
	Hair
	Tail
	MaxGene
)

// setState type
const (
	SaleState = iota
	SalePrice
	MateState
	MatePrice
)

type buyHistory struct {
	gene 	string
	time	int64
	price 	int
}
type BuySlice []*buyHistory

func (b BuySlice) Len() int {
	return len(b)
}
func (b BuySlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b BuySlice) Less(i, j int) bool {
	return b[i].time > b[j].time
}

type cat struct {
	Gene 		string   	`json:"gene"`		// cat unique key
	Name 	   	string 	 	`json:"name"`
	SaleState 	int			`json:"sale_state"`	// 0 not sale, 1 on sale
	SalePrice 	int    	 	`json:"sale_price"`
	GenNum		int		 	`json:"gen"`  		// Generations number
	Birth 		int64 	 	`json:"birth"`
	Parents		[2]string 	`json:"parents"`
	Children	[]string	`json:"children"`
	Owner		string 	 	`json:"owner"`
	MateState	int			`json:"mate_state"`	// 0 forbid mate, 1 allow mate
	MatePrice	int			`json:"mate_price"`
	MateTime 	int64		`json:"mate_time"`
}

func main() {
	err := shim.Start(new (CatChainCode))
	if err == nil {
		fmt.Printf("Error starting CatChaincode: %s", err)
	}
}

// Init initializes chaincode
func (c *CatChainCode)Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	fmt.Println(args[0], args[1])
	if len(args) != 6 {
		return shim.Error("initConfig, Incorrect number of arguments. Expecting 6")
	}

	return c.initConfig(stub, args)
}

// Invoke func
func (c *CatChainCode)Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case InitSystemCat:
		if len(args) != 0 {
			return shim.Error("initSystemCat, Incorrect number of arguments. Expecting 0")
		}
		return c.initSystemCat(stub)
	case QueryCat:
		if len(args) != 1 {
			return shim.Error("cat query, Incorrect number of arguments. Expecting 1")
		}
		return c.query(stub, args)
	case QueryAll:
		if len(args) != 0 {
			return shim.Error("cat queryAll, Incorrect number of arguments. Expecting 0")
		}
		return c.queryAll(stub)
	case QuerySale:
		if len(args) != 0 {
			return shim.Error("cat querySale, Incorrect number of arguments. Expecting 0")
		}
		return c.querySale(stub)
	case DelCat:
		if len(args) != 1 {
			return shim.Error("cat delete, Incorrect number of arguments. Expecting 1")
		}
		return c.delete(stub, args)
	case BreedCat:
		if len(args) != 2 {
			return shim.Error("breedCat, Incorrect number of arguments. Expecting 2")
		}
		return c.breed(stub, args)
	case BuyCat:
		if len(args) != 1 {
			return shim.Error("buy cat, Incorrect number of arguments. Expecting 1")
		}
		return c.buy(stub, args)
	case SetState:
		if len(args) != 3 {
			return shim.Error("setState, Incorrect number of arguments. Expecting 3")
		}
		return c.setState(stub, args)
	}

	return shim.Error("function is invalid: " + function)
}

// init constant
func (c *CatChainCode)initConfig(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	c.SysCatTime, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("initConfig, strconv.Atoi(args[0]) failed, " + err.Error())
	}
	c.LatestBuyPrice, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("initConfig, strconv.Atoi(args[1]) failed, " + err.Error())
	}
	c.InitPrice, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("initConfig, strconv.Atoi(args[2]) failed, " + err.Error())
	}
	c.NextMateTime, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("initConfig, strconv.Atoi(args[3]) failed, " + err.Error())
	}
	c.TokenType = args[4]
	c.SystemAccount = args[5]
	fmt.Println("cat Init Config.")
	fmt.Println("cat chaincode Init.")
	return shim.Success([]byte("initConfig succeed"))
}

// generate Gen0 cat
func (c *CatChainCode)initSystemCat(stub shim.ChaincodeStubInterface) pb.Response {
	// check sender is cat owner or not
	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)
	if new_add != c.SystemAccount {
		return shim.Error("init system cat failed,sysaccount:" + c.SystemAccount + ", sender:" + new_add)
	}

	// random a cat
	var gene string
	curTime := time.Now().Unix()
	for {
		if c.timeDuration(curTime) > c.SysCatTime {
			break
		}
		gene = c.newGenome()
		catAsbytes, err := stub.GetState(gene)
		if err == nil && catAsbytes == nil {
			break
		}
	}
	if gene == "" {
		jsonResp := "{\"Error\":\"initSystemCat timeout!\"}"
		return shim.Error(jsonResp)
	}

	// according to buy history,decide cat price
	price := c.InitPrice
	historySlice := c.getBuyHistory(stub)
	if len(historySlice) > 0 {
		priceCount := 0
		sort.Sort(historySlice)
		i := 0
		for _, v := range historySlice {
			priceCount += v.price
			i++
			if i >= c.LatestBuyPrice {
				break
			}
		}
		price = priceCount/i
	}

	// get total cat number, give cat a name
	nCount := c.getStateByRange(stub)
	name := fmt.Sprintf("Kitty%v", nCount+1)
	myCat := &cat {
		Gene:		gene,
		Name:		name,
		SaleState:	1,
		SalePrice:	price,
		GenNum:		0,
		Birth:		time.Now().Unix(),
		Parents:	[2]string{"",""},
		Children:	[]string{},
		Owner:		c.SystemAccount,
		MateState: 	0,
		MatePrice: 	0,
		MateTime: 	0,
	}

	catJSONasByte, err := json.Marshal(myCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(gene, catJSONasByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("init %v Cat Succeed", gene)
	return shim.Success([]byte(result))
}

// delete specify cat
func (c *CatChainCode)delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: catGene
	gene := args[0]
	err := stub.DelState(gene)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to delete cat for " + gene + "," + err.Error() +"\"}"
		return shim.Error(jsonResp)
	}

	result := fmt.Sprintf("delete %v Cat Succeed", gene)
	return shim.Success([]byte(result))
}

// query specify cat
func (c *CatChainCode)query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: catGene
	gene := args[0]
	valAsbytes, err := stub.GetState(gene)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + gene + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"cat does not exist: " + gene + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// query all cat info
func (c *CatChainCode)queryAll(stub shim.ChaincodeStubInterface,) pb.Response {
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	bArrayIndex := 1
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

// query for sale cat
func (c *CatChainCode)querySale(stub shim.ChaincodeStubInterface,) pb.Response {
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	bArrayIndex := 1
	var tmpCat cat
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(queryResponse.Value, &tmpCat)
		if err != nil {
			continue
		}
		if tmpCat.SaleState == 0 {
			continue
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

// set cat state
func (c *CatChainCode)setState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: catGene, type, value
	// 1. check cat exist or not
	gene := args[0]
	stateType, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	valAsbytes, err := stub.GetState(gene)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + gene + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"cat does not exist: " + gene + "\"}"
		return shim.Error(jsonResp)
	}

	var  myCat cat
	err = json.Unmarshal([]byte(valAsbytes), &myCat)
	if err != nil {
		jsonResp := "{\"Error\":\"Json Unmarshal failed, " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	}

	// 2. check sender is cat owner or not
	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)
	if new_add != myCat.Owner {
		return shim.Error("set cat sale state failed,cat owner:" + myCat.Owner + ", sender:" + new_add)
	}

	// 3. edit cat state according to stateType
	switch stateType {
	case SaleState:
		// 0 not for sale, 1 for sale
		myCat.SaleState = value
	case SalePrice:
		myCat.SalePrice = value
	case MateState:
		// 0 forbid mate, 1 allow mate
		myCat.MateState = value
	case MatePrice:
		myCat.MatePrice = value
	default:
		return shim.Error("set cat state failed, stateType err: " + args[1])
	}

	// 4. put state to save
	catJSONasByte, err := json.Marshal(myCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(gene, catJSONasByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("set %v Cat sale stateType:%v, value:%v", args[0], args[1], args[2])
	return shim.Success([]byte(result))
}

// user's cat request breed with another cat
func (c *CatChainCode)breed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: catGene1, catGene2
	// sender is father,  another is mother
	// 1. check two cat is exist or not
	geneFather := args[0]
	fatherAsbytes, err := stub.GetState(geneFather)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get father state for " + geneFather + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if fatherAsbytes == nil {
		jsonResp := "{\"Error\":\"cat father does not exist: " + geneFather + "\"}"
		return shim.Error(jsonResp)
	}
	geneMother := args[1]
	motherAsbytes, err := stub.GetState(geneMother)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get mother state for " + geneMother + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if motherAsbytes == nil {
		jsonResp := "{\"Error\":\"cat mother does not exist: " + geneMother + "\"}"
		return shim.Error(jsonResp)
	}

	// 2. check father cat's owner is right or not
	var catFather cat
	err = json.Unmarshal(fatherAsbytes, &catFather)
	if err != nil {
		return shim.Error(err.Error())
	}
	var catMother cat
	err = json.Unmarshal(motherAsbytes, &catMother)
	if err != nil {
		return shim.Error(err.Error())
	}

	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)
	if new_add != catFather.Owner {
		return shim.Error("breed cat failed, cat owner:" + catFather.Owner + ", sender:" + new_add)
	}

	// 3. check mother cat is allow mate or not
	if catMother.MateState != 1 {
		jsonResp := "{\"Error\":\"cat mother not allow mate: " + geneMother + "\"}"
		return shim.Error(jsonResp)
	}else if c.timeDuration(catMother.MateTime) < c.NextMateTime {
		jsonResp := fmt.Sprintf("{\"Error\":\"cat mother after %v seconds can mate \"}", c.NextMateTime - c.timeDuration(catMother.MateTime))
		return shim.Error(jsonResp)
	}

	// 4. check account balance is enough or not
	account, err := stub.GetAccount(new_add)
	if err != nil {
		jsonResp := "{\"Error\":\"account not exists\"}"
		return shim.Error(jsonResp)
	}
	if account == nil || account.Balance[c.TokenType] == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + new_add + "\"}"
		return shim.Error(jsonResp)
	}

	balance := account.Balance[c.TokenType].Int64()
	if balance < int64(catMother.MatePrice) {
		jsonResp := "{\"Error\":\"breed failed, balance not enough, balance: " + account.Balance[c.TokenType].String() + ",price: " + strconv.Itoa(catMother.MatePrice) + "\"}"
		return shim.Error(jsonResp)
	}

	// 5. pay money
	amount := big.NewInt(0)
	_, good := amount.SetString(strconv.Itoa(catMother.MatePrice), 10)
	if !good  {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(catMother.Owner, c.TokenType, amount)
	if err != nil {
		return shim.Error("transfer error" + err.Error())
	}
	// 6.create new cat
	var newGene string
	for {
		newGene = c.breedGenome(geneFather, geneMother)
		if newGene == "" {
			jsonResp := "{\"Error\":\"createGenome failed: " + geneFather + "," + geneMother + "\"}"
			return shim.Error(jsonResp)
		}
		newAsbytes, err := stub.GetState(newGene)
		if err == nil && newAsbytes == nil {
			break
		}
	}

	// 7. set cat data and put
	catFather.Children = append(catFather.Children, newGene)
	catMother.Children = append(catMother.Children, newGene)
	catFather.MateTime = time.Now().Unix()
	catMother.MateTime = time.Now().Unix()

	genNum := catFather.GenNum
	if genNum < catMother.GenNum {
		genNum = catMother.GenNum
	}

	nCount := c.getStateByRange(stub)
	name := fmt.Sprintf("Kitty%v", nCount+1)

	myCat := &cat {
		Gene:		newGene,
		Name:		name,
		SaleState:	0,
		SalePrice:	0,
		GenNum:		genNum+1,
		Birth:		time.Now().Unix(),
		Parents:	[2]string{geneFather,geneMother},
		Children:	[]string{},
		Owner:		catFather.Owner,
		MateState: 	0,
		MatePrice:  0,
		MateTime: 	0,
	}

	catJSONasByte, err := json.Marshal(myCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(newGene, catJSONasByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	fatherBytes, err := json.Marshal(catFather)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(geneFather, fatherBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	motherBytes, err := json.Marshal(catMother)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(geneMother, motherBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("createGenome new Cat:%v", newGene)
	return shim.Success([]byte(result))
}

// create gen0 cat gene by random num
func (c *CatChainCode)newGenome() string {
	var buffer bytes.Buffer
	for i := 0; i < MaxGene; i++ {
		buffer.WriteString(strconv.Itoa(rand.Intn(10)))
	}

	return buffer.String()
}

func (c *CatChainCode)timeDuration(srcTime int64) int {
	return int(time.Now().Unix() - srcTime)
}

// breed cat inherit gene(0 father gene, 1 mother gene, 2 average)
func (c *CatChainCode)breedGenome(geneFather, geneMother string) string {
	if len(geneFather) != MaxGene {
		return ""
	}
	if len(geneMother) != MaxGene {
		return ""
	}

	var buffer bytes.Buffer
	for i := 0; i < MaxGene; i++ {
		switch rand.Intn(3) {
		case 0:
			buffer.WriteByte(geneFather[i])
		case 1:
			buffer.WriteByte(geneMother[i])
		case 2:
			j := geneFather[i] +geneMother[i]
			buffer.WriteByte(j/2)
		default:
			buffer.WriteByte(geneFather[i])
		}
	}

	return buffer.String()
}

func (c *CatChainCode) getStateByRange(stub shim.ChaincodeStubInterface) (count int) {
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		fmt.Println("getCountByRange: stub.GetStateByRange" + err.Error())
		return
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return
		}
		count++
	}

	return
}

// user buy cat
func (c *CatChainCode) buy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: caetGene
	// 1. check cat is exist or not
	catGene := args[0]
	catAsbytes, err := stub.GetState(catGene)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + catGene + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if catAsbytes == nil {
		jsonResp := "{\"Error\":\"cat does not exist: " + catGene + "\"}"
		return shim.Error(jsonResp)
	}
	// 2. check cat is for sale or not
	var buyCat cat
	err = json.Unmarshal(catAsbytes, &buyCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	if buyCat.SaleState == 0 {
		jsonResp := "{\"Error\":\"buy failed, cat does not sale: " + catGene + "\"}"
		return shim.Error(jsonResp)
	}

	// 3. check cat owner is mine or not
	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)
	if new_add == buyCat.Owner {
		jsonResp := "{\"Error\":\"buy failed, already is cat owner: " + catGene + "\"}"
		return shim.Error(jsonResp)
	}

	// 4. check account balance is enough or not
	account, err := stub.GetAccount(new_add)
	if err != nil {
		jsonResp := "{\"Error\":\"account not exists\"}"
		return shim.Error(jsonResp)
	}
	if account == nil || account.Balance[c.TokenType] == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + new_add + "\"}"
		return shim.Error(jsonResp)
	}

	balance := account.Balance[c.TokenType].Int64()
	if balance < int64(buyCat.SalePrice) {
		jsonResp := "{\"Error\":\"buy failed, balance not enough, balance: " + account.Balance[c.TokenType].String() + ",price: " + strconv.Itoa(buyCat.SalePrice) + "\"}"
		return shim.Error(jsonResp)
	}

	// 5. pay cat owner cost
	amount := big.NewInt(0)
	_, good := amount.SetString(strconv.Itoa(buyCat.SalePrice), 10)
	if !good  {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(buyCat.Owner, c.TokenType, amount)
	if err != nil {
		return shim.Error("transfer error" + err.Error())
	}

	// 6. edit cat date and put
	buyCat.Owner = new_add
	buyCat.SaleState = 0
	buyCat.MateState = 0

	buyAsbytes, err := json.Marshal(buyCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(catGene, buyAsbytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 7. buy record CreateCompositeKey and save
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	err = c.addBuyHistory(stub, []string{curTime, strconv.Itoa(buyCat.SalePrice), catGene})
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("buy %v Cat Succeed", catGene)
	return shim.Success([]byte(result))
}

// save buy record
func (c *CatChainCode) addBuyHistory(stub shim.ChaincodeStubInterface, args []string) error {
	// time ~ price ~ catGene
	var indexKey string
	var err error
	indexKey, err = stub.CreateCompositeKey(HistoryIndex, args)
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

func (c *CatChainCode) getBuyHistory(stub shim.ChaincodeStubInterface) (buySlice BuySlice) {
	buySlice = BuySlice{}
	resultsIterator, err := stub.GetStateByPartialCompositeKey(HistoryIndex, []string{})
	if err != nil {
		fmt.Println("getBuyHistory, stub.GetStateByPartialCompositeKey failed: " + err.Error())
		return
	}

	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return
		}
		// time ~ price ~ catGene
		buyTime := compositeKeyParts[0]
		buyPrice := compositeKeyParts[1]
		gene := compositeKeyParts[2]
		tm, _ := strconv.ParseInt(buyTime, 10,64)
		price, _ := strconv.ParseInt(buyPrice, 10,32)
		buySlice = append(buySlice, &buyHistory{gene, tm, int(price)})
		fmt.Println("getBuyHistory, stub.GetStateByPartialCompositeKey: " +  buyTime + "," + buyPrice + "," +gene)
	}
	return
}
