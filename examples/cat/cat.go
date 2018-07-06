/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.
*/

// cat: a demo chaincode for inkchain

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
)

const (
	// function
	InitSystemCat = "initSystemCat"
	QueryCat      = "query"
	QueryAll      = "queryAll"
	QuerySale     = "querySale"
	DelCat        = "delete"
	BreedCat      = "breed"
	BuyCat        = "buy"
	SetState      = "setState"
	CreateAuction = "createAuction"
	Bid           = "bid"
	EndAuction    = "endAuction"
	PayAuction    = "payAuction"
	QueryAuction  = "queryAuction"
	// config
	BuyHistory         = "buyHistory" // cat's record, ps: time ~ price ~ catGene
	BidHistory         = "bidHistory" // cat's record, ps: catGene ~ bidder ~ bidTime ~ bidPrice
	BidOrder           = "bidOrder"   // user confirm bid, ps: bidder ~ bidTime ~ catGene ~ bidPrice ~ orderTime
	BidOrderExpiryDate = 86400        // 24 hour
)

type CatChainCode struct {
	SysCatTime     int    // system generate cat max time, if timeout return failed
	LatestBuyPrice int    // count latest buy record number
	InitPrice      int    // init cat price
	NextMateTime   int    // cat can mate time since last mate
	TokenType      string // tokenType as "INK"
	SystemAccount  string // system account, user buy gen0 pay money to system
}

// setState type
const (
	SaleState = iota
	SalePrice
	MateState
	MatePrice
)

type buyHistory struct {
	gene  string
	time  int64
	price int
}

type cat struct {
	Gene      string       `json:"gene"` // cat unique key
	Name      string       `json:"name"`
	SaleState int          `json:"sale_state"` // 0 not sale, 1 on sale
	SalePrice int          `json:"sale_price"`
	GenNum    int          `json:"gen"` // Generations number
	Birth     int64       `json:"birth"`
	Parents   [2]string    `json:"parents"`
	Children  []string     `json:"children"`
	Owner     string       `json:"owner"`
	MateState int          `json:"mate_state"` // 0 forbid mate, 1 allow mate
	MatePrice int          `json:"mate_price"`
	MateTime  int64        `json:"mate_time"`
	Auction   *saleAuction `json:"auction"`
}

type saleAuction struct {
	AuctionState    int // 0not auction, 1 auction
	LastAuctionTime int64

	StartingTime  int64
	StartingPrice int
	Duration      int

	Bidder   string
	BidTime  int64
	BidPrice int
}

func main() {
	err := shim.Start(new(CatChainCode))
	if err == nil {
		fmt.Printf("Error starting CatChaincode: %s", err)
	}
}

// Init initializes chaincode
func (c *CatChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	fmt.Println(args[0], args[1])
	if len(args) != 6 {
		return shim.Error("initConfig, Incorrect number of arguments. Expecting 6")
	}

	return c.initConfig(stub, args)
}

// Invoke func
func (c *CatChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case InitSystemCat:
		if len(args) != 1 {
			return shim.Error("initSystemCat, Incorrect number of arguments. Expecting 1")
		}
		return c.newGenZeroCat(stub, args)
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
		if len(args) != 3 {
			return shim.Error("breedCat, Incorrect number of arguments. Expecting 3")
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
	case CreateAuction:
		if len(args) != 3 {
			return shim.Error("CreateAuction, Incorrect number of arguments. Expecting 3")
		}
		return c.createAuction(stub, args)
	case Bid:
		if len(args) != 2 {
			return shim.Error("bid, Incorrect number of arguments. Expecting 2")
		}
		return c.bid(stub, args)
	case EndAuction:
		if len(args) != 1 {
			return shim.Error("EndAuction, Incorrect number of arguments. Expecting 1")
		}
		return c.endAuction(stub, args)
	case PayAuction:
		if len(args) != 3 {
			return shim.Error("PayAuction, Incorrect number of arguments. Expecting 3")
		}
		return c.payAuction(stub, args)
	case QueryAuction:
		if len(args) != 1 {
			return shim.Error("PayAuction, Incorrect number of arguments. Expecting 1")
		}
		return c.queryAuction(stub, args)

	}

	return shim.Error("function is invalid: " + function)
}

// init constant
func (c *CatChainCode) initConfig(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	c.SystemAccount = strings.ToLower(args[5])
	fmt.Println("cat Init Config.")
	fmt.Println("cat chaincode Init.")
	return shim.Success([]byte("initConfig succeed"))
}

// generate Gen0 cat only by system address
// @Param Gene
func (c *CatChainCode) newGenZeroCat(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	gene := args[0]
	if gene == "" {
		return shim.Error("{\"Error\":\"invalid cat gene!\"}")
	}

	// check sender is cat owner or not
	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)
	if new_add != c.SystemAccount {
		return shim.Error("init system cat failed,sysaccount:" + c.SystemAccount + ", sender:" + new_add)
	}

	if catAsBytes, err := stub.GetState(gene); err != nil {
		return shim.Error(err.Error())
	} else if catAsBytes != nil {
		return shim.Error("{\"Error\":\" cat already exist!\"}")
	}

	// according to buy history,decide cat price
	price := c.InitPrice
	historySlice := c.getBuyHistory(stub)
	if len(historySlice) > 0 {
		priceCount := 0
		i := 0
		for _, v := range historySlice {
			priceCount += v.price
			i++
		}
		price = priceCount / i
	}
	timstamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	// get total cat number, give cat a name
	nCount := c.getStateByRange(stub)
	name := fmt.Sprintf("Kitty%v", nCount+1)
	myCat := &cat{
		Gene:      gene,
		Name:      name,
		SaleState: 1,
		SalePrice: price,
		GenNum:    0,
		Birth:     timstamp.GetSeconds(),
		Parents:   [2]string{"", ""},
		Children:  []string{},
		Owner:     c.SystemAccount,
		MateState: 0,
		MatePrice: 0,
		MateTime:  0,
		Auction:   &saleAuction{},
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
func (c *CatChainCode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: catGene
	gene := args[0]
	err := stub.DelState(gene)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to delete cat for " + gene + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	}

	result := fmt.Sprintf("delete %v Cat Succeed", gene)
	return shim.Success([]byte(result))
}

// query specify cat
func (c *CatChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
func (c *CatChainCode) queryAll(stub shim.ChaincodeStubInterface) pb.Response {
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
		bArrayIndex++
	}
	buffer.WriteString("]")

	return shim.Success([]byte(buffer.String()))
}

// query for sale cat
func (c *CatChainCode) querySale(stub shim.ChaincodeStubInterface) pb.Response {
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
		bArrayIndex++
	}
	buffer.WriteString("]")

	return shim.Success([]byte(buffer.String()))
}

// set cat state
func (c *CatChainCode) setState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	valAsBytes, err := stub.GetState(gene)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + gene + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsBytes == nil {
		jsonResp := "{\"Error\":\"cat does not exist: " + gene + "\"}"
		return shim.Error(jsonResp)
	}

	var myCat cat
	err = json.Unmarshal([]byte(valAsBytes), &myCat)
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

	timStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	if 1 == myCat.Auction.AuctionState || (timStamp.GetSeconds()-myCat.Auction.LastAuctionTime) < BidOrderExpiryDate {
		return shim.Error("cat is on Auction state, cant sale or mate")
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
// @Param fatherCat gene
// @Param motherCat gene
// @Param newCat gene
func (c *CatChainCode) breed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	newGene := args[2]
	if newAsBytes, err := stub.GetState(newGene); err != nil {
		return shim.Error(err.Error())
	} else if newAsBytes != nil {
		return shim.Error("{\"Error\":\" cat already exist!\"}")
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

	timStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	// 3. check mother cat is allow mate or not
	if catMother.MateState != 1 {
		jsonResp := "{\"Error\":\"cat mother not allow mate: " + geneMother + "\"}"
		return shim.Error(jsonResp)
	} else if (timStamp.GetSeconds()-catMother.MateTime) < int64(c.NextMateTime) {
		jsonResp := fmt.Sprintf("{\"Error\":\"cat mother after %v seconds can mate \"}", int64(c.NextMateTime)+catMother.MateTime-timStamp.GetSeconds())
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
	if !good {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(catMother.Owner, c.TokenType, amount)
	if err != nil {
		return shim.Error("transfer error" + err.Error())
	}

	// 7. set cat data and put
	catFather.Children = append(catFather.Children, newGene)
	catMother.Children = append(catMother.Children, newGene)
	catFather.MateTime = timStamp.GetSeconds()
	catMother.MateTime = timStamp.GetSeconds()

	genNum := catFather.GenNum
	if genNum < catMother.GenNum {
		genNum = catMother.GenNum
	}

	nCount := c.getStateByRange(stub)
	name := fmt.Sprintf("Kitty%v", nCount+1)

	myCat := &cat{
		Gene:      newGene,
		Name:      name,
		SaleState: 0,
		SalePrice: 0,
		GenNum:    genNum + 1,
		Birth:     timStamp.GetSeconds(),
		Parents:   [2]string{geneFather, geneMother},
		Children:  []string{},
		Owner:     catFather.Owner,
		MateState: 0,
		MatePrice: 0,
		MateTime:  0,
		Auction:   &saleAuction{},
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
	// args: catGene
	// 1. check cat is exist or not
	catGene := args[0]
	catAsBytes, err := stub.GetState(catGene)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + catGene + "," + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if catAsBytes == nil {
		jsonResp := "{\"Error\":\"cat does not exist: " + catGene + "\"}"
		return shim.Error(jsonResp)
	}
	// 2. check cat is for sale or not
	var buyCat cat
	err = json.Unmarshal(catAsBytes, &buyCat)
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
	bEnough := c.checkBalance(stub, new_add, buyCat.SalePrice)
	if bEnough == false {
		return shim.Error("bid failed, balance not enough")
	}

	// 5. pay cat owner cost
	amount := big.NewInt(0)
	_, good := amount.SetString(strconv.Itoa(buyCat.SalePrice), 10)
	if !good {
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

	buyAsBytes, err := json.Marshal(buyCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(catGene, buyAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 7. buy record CreateCompositeKey and save
	timStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	err = c.historyByComposite(stub, BuyHistory, []string{timStamp.String(), strconv.Itoa(buyCat.SalePrice), catGene})
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("buy %v Cat Succeed", catGene)
	return shim.Success([]byte(result))
}

// save record
func (c *CatChainCode) historyByComposite(stub shim.ChaincodeStubInterface, sType string, args []string) error {
	var indexKey string
	var err error
	indexKey, err = stub.CreateCompositeKey(sType, args)
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

func (c *CatChainCode) getBuyHistory(stub shim.ChaincodeStubInterface) (buySlice []*buyHistory) {
	buySlice = []*buyHistory{}
	resultsIterator, err := stub.GetStateByPartialCompositeKey(BuyHistory, []string{})
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
		tm, _ := strconv.ParseInt(buyTime, 10, 64)
		price, _ := strconv.ParseInt(buyPrice, 10, 32)
		buySlice = append(buySlice, &buyHistory{gene, tm, int(price)})
		fmt.Println("getBuyHistory, stub.GetStateByPartialCompositeKey: " + buyTime + "," + buyPrice + "," + gene)
		if len(buySlice) >= c.LatestBuyPrice {
			break
		}
	}
	return
}
func (c *CatChainCode) createAuction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Param1:catGene
	// Param2:startingPrice
	// Param3:duration(seconds)
	// 1. check cat exist or not
	gene := args[0]
	startingPrice, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	duration, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	myCat, err := c.getCat(stub, gene)
	if err != nil {
		return shim.Error(err.Error())
	} else if myCat.Gene == "" {
		jsonResp := "{\"Error\":\"cat does not exist: " + gene + "\"}"
		return shim.Error(jsonResp)
	}

	// 2. check sender is cat owner or not
	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)
	if new_add != myCat.Owner {
		return shim.Error("auction failed,cat owner:" + myCat.Owner + ", sender:" + new_add)
	}

	// 3.check cat is on auction
	if nil == myCat.Auction {
		return shim.Error("Cat no Auction info")
	}
	if 1 == myCat.Auction.AuctionState {
		return shim.Error("Cat is on Auction")
	}
	timStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	myCat.MateState = 0
	myCat.SaleState = 0
	myCat.Auction.AuctionState = 1
	myCat.Auction.StartingTime = timStamp.GetSeconds()
	myCat.Auction.StartingPrice = startingPrice
	myCat.Auction.Duration = duration
	myCat.Auction.BidPrice = startingPrice

	// 4. put state to save
	catJSONasByte, err := json.Marshal(myCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(gene, catJSONasByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("Cat:%v auctuon publish succeed", args[0])
	return shim.Success([]byte(result))
}

func (c *CatChainCode) bid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Param1:bidCatGene
	// Param2:bidPrice
	// 1. check cat exist or not
	gene := args[0]
	bidPrice, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	myCat, err := c.getCat(stub, gene)
	if err != nil {
		return shim.Error(err.Error())
	} else if myCat.Gene == "" {
		jsonResp := "{\"Error\":\"cat does not exist: " + gene + "\"}"
		return shim.Error(jsonResp)
	}

	// 2. check cat auction state
	if 0 == myCat.Auction.AuctionState {
		result := fmt.Sprintf("Cat:%v is not Auction", myCat.Gene)
		return shim.Error(result)
	}

	bidder, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	bidder = strings.ToLower(bidder)
	if bidder == myCat.Owner {
		return shim.Error("cat owner cant join auction")
	}

	// 3.check auction Expiration Date

	timstamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	if (timstamp.GetSeconds()-myCat.Auction.StartingTime) >= int64(myCat.Auction.Duration) {
		return shim.Error("auction already end!")
	}

	// 4.check bidPrice
	if nil == myCat.Auction {
		return shim.Error("Cat on Auction info")
	}
	if bidPrice <= myCat.Auction.BidPrice {
		return shim.Error("bidPrice lower than current price")
	}

	// 5. check account balance is enough or not
	bEnough := c.checkBalance(stub, bidder, bidPrice)
	if bEnough == false {
		return shim.Error("bid failed, balance not enough")
	}

	err = c.historyByComposite(stub, BidHistory, []string{gene, bidder,timstamp.String(), strconv.Itoa(bidPrice)})
	if err != nil {
		return shim.Error(err.Error())
	}

	myCat.Auction.Bidder = bidder
	myCat.Auction.BidTime = timstamp.GetSeconds()
	myCat.Auction.BidPrice = bidPrice
	// 4. put state to save
	catJSONasByte, err := json.Marshal(myCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(gene, catJSONasByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("bidder:%v bid cat:%v succeed", bidder, gene)
	return shim.Success([]byte(result))
}

func (c *CatChainCode) checkBalance(stub shim.ChaincodeStubInterface, add string, amount int) bool {
	account, err := stub.GetAccount(add)
	if err != nil {
		// account not exists
		return false
	}
	if account == nil || account.Balance[c.TokenType] == nil {
		return false
	}

	balance := account.Balance[c.TokenType].Int64()
	if balance < int64(amount) {
		// balance not enough
		return false
	}

	return true
}

func (c *CatChainCode) getCat(stub shim.ChaincodeStubInterface, gene string) (myCat cat, err error) {
	var valAsbytes []byte
	myCat = cat{}
	valAsbytes, err = stub.GetState(gene)
	if err != nil {
		// Failed to get state
		return
	} else if valAsbytes == nil {
		// cat does not exist
		return
	}

	err = json.Unmarshal([]byte(valAsbytes), &myCat)
	return
}

func (c *CatChainCode) endAuction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Param1:catGene
	// 1. check cat exist or not
	gene := args[0]
	myCat, err := c.getCat(stub, gene)
	if err != nil {
		return shim.Error(err.Error())
	} else if myCat.Gene == "" {
		jsonResp := "{\"Error\":\"cat does not exist: " + gene + "\"}"
		return shim.Error(jsonResp)
	}

	// 2. check sender is cat owner or not
	new_add, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	new_add = strings.ToLower(new_add)
	if new_add != myCat.Owner {
		return shim.Error("auction failed,cat owner:" + myCat.Owner + ", sender:" + new_add)
	}

	// 3.check cat is on auction
	if 0 == myCat.Auction.AuctionState {
		return shim.Error("Cat is not Auction")
	}
	if nil == myCat.Auction {
		return shim.Error("Cat no Auction info")
	}
	bidder := myCat.Auction.Bidder
	bidTime := myCat.Auction.BidTime
	bidPrice := myCat.Auction.BidPrice
	// 4.check auction Expiration Date

	timStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	if (timStamp.GetSeconds()-myCat.Auction.StartingTime) < int64(myCat.Auction.Duration) && bidder == "" {
		return shim.Error("auction is not end yet!")
	}
	if bidder != "" {
		err = c.historyByComposite(stub, BidOrder, []string{bidder, strconv.FormatInt(bidTime, 10), gene, strconv.Itoa(bidPrice),
			timStamp.String()})
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	myCat.Auction.AuctionState = 0
	myCat.Auction.LastAuctionTime = timStamp.GetSeconds()
	myCat.Auction.StartingTime = 0
	myCat.Auction.StartingPrice = 0
	myCat.Auction.Duration = 0
	myCat.Auction.Bidder = ""
	myCat.Auction.BidTime = 0
	myCat.Auction.BidPrice = 0

	// 5. put state to save
	catJSONasByte, err := json.Marshal(myCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(gene, catJSONasByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("Cat:%v auctuon end", args[0])
	return shim.Success([]byte(result))
}

func (c *CatChainCode) payAuction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Param1:catGene
	// Param2:bidTime
	// Param3:bidPrice
	// Param4:orderTime
	// 1. check cat exist or not
	gene := args[0]
	bidPrice, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}

	myCat, err := c.getCat(stub, gene)
	if err != nil {
		return shim.Error(err.Error())
	} else if myCat.Gene == "" {
		jsonResp := "{\"Error\":\"cat does not exist: " + gene + "\"}"
		return shim.Error(jsonResp)
	}

	// 2. check sender is cat owner or not
	bidder, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	bidder = strings.ToLower(bidder)
	if bidder == myCat.Owner {
		return shim.Error("pay auction failed,already is cat owner")
	}

	// 4. check account balance is enough or not
	bEnough := c.checkBalance(stub, bidder, int(bidPrice))
	if bEnough == false {
		return shim.Error("bid failed, balance not enough")
	}

	// 5. pay cat owner cost
	amount := big.NewInt(0)
	_, good := amount.SetString(args[2], 10)
	if !good {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(myCat.Owner, c.TokenType, amount)
	if err != nil {
		return shim.Error("transfer error" + err.Error())
	}

	// 6. edit cat date and put
	myCat.Owner = bidder
	myCat.SaleState = 0
	myCat.MateState = 0

	buyAsBytes, err := json.Marshal(myCat)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(gene, buyAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	result := fmt.Sprintf("user:%v pay auctuon succeed", bidder)
	return shim.Success([]byte(result))
}

func (c *CatChainCode) checkBidOrder(stub shim.ChaincodeStubInterface, args []string) bool {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(BidOrder, args)
	if err != nil {
		fmt.Println("getBuyHistory, stub.GetStateByPartialCompositeKey failed: " + err.Error())
		return false
	}

	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return false
		}
		_, _, err = stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return false
		}

		return true
	}
	return false
}

func (c *CatChainCode) queryAuction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(BidHistory, args)
	if err != nil {
		result := fmt.Sprintf("queryAuction, stub.GetStateByPartialCompositeKey failed: " + err.Error())
		return shim.Error(result)
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		// catGene ~ bider ~ bidTime ~ bidPrice
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Number\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString("\"")
		buffer.WriteString(", \"catGene\":")
		buffer.WriteString(compositeKeyParts[0])
		buffer.WriteString("\"")
		buffer.WriteString(", \"bidder\":")
		buffer.WriteString(compositeKeyParts[1])
		buffer.WriteString("\"")
		buffer.WriteString(", \"bidTime\":")
		buffer.WriteString(compositeKeyParts[2])
		buffer.WriteString(", \"bidPrice\":")
		buffer.WriteString(compositeKeyParts[3])
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success([]byte(buffer.String()))
}
