package invoke

import (
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	. "github.com/inklabsfoundation/inkchain/examples/creative/model"
	. "github.com/inklabsfoundation/inkchain/examples/creative/conf"
	. "github.com/inklabsfoundation/inkchain/examples/creative/util"
	"fmt"
	"encoding/json"
	"strings"
	"math/big"
	"strconv"
)

type ProductionInvoke struct{}

func (*ProductionInvoke) AddProduction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("add production start.")
	username := args[0]
	production_type := args[1]
	production_serial := args[2]
	production_name := args[3]
	production_desc := args[4]
	production_price_type := args[5]
	production_price := args[6]
	production_num := args[7]
	production_transfer_part := args[8]

	// verify weather the user exists
	user_key := GetUserKey(username)
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}

	// verify weather the artist exists
	artist_key := GetArtistKey(username)
	artistAsBytes, err := stub.GetState(artist_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}

	if artistAsBytes == nil {
		fmt.Println("This artist doesn't exist: " + artist_key)
		return shim.Error("This artist doesn't exist: " + artist_key)
	}

	product_key := GetProductionKey(username, production_type, production_serial)

	// check if user exists
	productionAsBtyes, err := stub.GetState(product_key)
	if err != nil {
		return shim.Error("Fail to get production: " + err.Error())
	} else if productionAsBtyes != nil {
		fmt.Println("This production already exists: " + product_key)
		return shim.Error("This production already exists: " + product_key)
	}

	// add production
	production := Production{
		Type:                  production_type,
		Serial:                production_serial,
		Name:                  production_name,
		Desc:                  production_desc,
		CopyrightPriceType:    production_price_type,
		CopyrightPrice:        production_price,
		CopyrightNum:          production_num,
		CopyrightTransferPart: production_transfer_part,
		Username:              username,
	}

	productionAsBtyes, err = json.Marshal(production)
	err = stub.PutState(product_key, productionAsBtyes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("add production success."))
}

func (*ProductionInvoke) ModifyProduction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("modify production start.")
	username := args[0]
	production_type := args[1]
	production_serial := args[2]
	user_key := GetUserKey(username)
	// verify weather the user exists
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}

	// get user's address
	address, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	address = strings.ToLower(address)

	var userJSON User
	err = json.Unmarshal([]byte(userAsBytes), &userJSON)
	if userJSON.Address != address {
		return shim.Error("The sender's address doesn't correspond with the user's.")
	}

	production_key := GetProductionKey(username, production_type, production_serial)
	// verify weather the production exists
	productionAsBytes, err := stub.GetState(production_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}
	if productionAsBytes == nil {
		fmt.Println("This production doesn't exist: " + production_key)
		return shim.Error("This production doesn't exist: " + production_key)
	}

	var productionJSON Production
	err = json.Unmarshal([]byte(productionAsBytes), &productionJSON)

	if productionJSON.Username != userJSON.Username {
		return shim.Error("The artist's username doesn't correspond with the user's.")
	}

	err = GetModifyProduction(&productionJSON, args[3:])
	if err != nil {
		return shim.Error(err.Error())
	}
	artistJSONasBytes, err := json.Marshal(productionJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(production_key, artistJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("modify production success."))
}

func (*ProductionInvoke) DeleteProduction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("delete production start.")
	username := args[0]
	production_type := args[1]
	production_serial := args[2]
	user_key := GetUserKey(username)
	// verify weather the user exists
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}
	// get user's address
	address, err := stub.GetSender()
	if err != nil {
		return shim.Error("Fail to reveal user's address.")
	}
	address = strings.ToLower(address)
	var userJSON User
	err = json.Unmarshal([]byte(userAsBytes), &userJSON)
	if userJSON.Address != address {
		return shim.Error("The sender's address doesn't correspond with the user's.")
	}
	production_key := GetProductionKey(username, production_type, production_serial)
	// verify weather the production exists
	productionAsBytes, err := stub.GetState(production_key)
	if err != nil {
		return shim.Error("Fail to get production: " + err.Error())
	}
	if productionAsBytes == nil {
		fmt.Println("This production doesn't exist: " + production_key)
		return shim.Error("This production doesn't exist: " + production_key)
	}
	var productionJSON Production
	err = json.Unmarshal([]byte(productionAsBytes), &productionJSON)
	if productionJSON.Username != userJSON.Username {
		return shim.Error("The artist's username doesn't correspond with the user's.")
	}
	// delete production's info
	err = stub.DelState(production_key)
	if err != nil {
		fmt.Println("Fail to delete: " + production_key)
		return shim.Error("Fail to delete" + production_key)
	}
	return shim.Success([]byte("delete production success."))
}

func (*ProductionInvoke) QueryProduction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("query production start.")
	username := args[0]
	production_type := args[1]
	production_serial := args[2]
	production_key := GetProductionKey(username, production_type, production_serial)
	productionAsBytes, err := stub.GetState(production_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}
	if productionAsBytes == nil {
		fmt.Println("This production doesn't exist: " + production_key)
		return shim.Error("This production doesn't exist: " + production_key)
	}
	return shim.Success(productionAsBytes)
}

func (*ProductionInvoke) ListOfProduction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("list of production start.")
	username := args[0]
	product_type := args[1]
	state_key := GetStateKey(ProductionPrefix, username, product_type)
	resultsIterator, err := stub.GetStateByRange(state_key+StateStartSymbol, state_key+StateEndSymbol)
	if err != nil {
		return shim.Error(err.Error())
	}
	list, err := GetListResult(resultsIterator)
	if err != nil {
		return shim.Error("getListResult failed")
	}
	return shim.Success(list)
}

func (*ProductionInvoke) ListOfSupporter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("list of supporter start.")
	username := args[0]
	product_type := args[1]
	product_serial := args[2]
	state_key := GetStateKey(ProductionPrefix, username, product_type, product_serial)
	var list []Production
	fmt.Println("state_key:", state_key)
	if product_serial != "" {
		product, err := stub.GetState(state_key)
		if err != nil {
			return shim.Error(err.Error())
		}
		var pro Production
		json.Unmarshal(product, &pro)
		list = append(list, pro)
	} else {
		resultsIterator, err := stub.GetStateByRange(state_key+StateStartSymbol, state_key+StateEndSymbol)
		if err != nil {
			return shim.Error(err.Error())
		}
		list, err = GetListProduction(resultsIterator)
		fmt.Println("list:", list)
		if err != nil {
			return shim.Error("getListResult failed")
		}
	}
	result := make(map[string]string)
	for _, production := range list {
		supporters := production.Supporters
		for k, v := range supporters {
			if _, ok := result[k]; ok {
				result[k] = Add(result[k], v)
			} else {
				result[k] = v
			}
		}
	}
	resultAsBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultAsBytes == nil {
		fmt.Println("This supporters doesn't exist: ")
		return shim.Error("This supporters doesn't exist: ")
	}
	return shim.Success(resultAsBytes)
}

func (*ProductionInvoke) AddSupporter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("add buyer start.")
	username := args[0]
	production_type := args[1]
	production_serial := args[2]
	price_type := args[3]
	price := args[4]
	supporter := args[5]
	user_key := GetUserKey(username)
	// verify weather the user exists
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}

	production_key := GetProductionKey(username, production_type, production_serial)
	// verify weather the production exists
	productionAsBytes, err := stub.GetState(production_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}
	if productionAsBytes == nil {
		fmt.Println("This production doesn't exist: " + production_key)
		return shim.Error("This production doesn't exist: " + production_key)
	}

	var productionJSON Production
	err = json.Unmarshal([]byte(productionAsBytes), &productionJSON)

	var userJSON User
	err = json.Unmarshal([]byte(userAsBytes), &userJSON)
	if userJSON.Username != productionJSON.Username {
		return shim.Error("The production's username doesn't correspond with the user's.")
	}

	if price_type != productionJSON.CopyrightPriceType {
		return shim.Error("The production's priceType doesn't correspond with the supporter's.")
	}

	// step 4: make transfer
	toAddress := userJSON.Address
	amount := big.NewInt(0)
	_, good := amount.SetString(price, 10)
	if !good {
		return shim.Error("Expecting integer value for amount")
	}

	fmt.Println(toAddress, ", ", price_type, "，", amount)
	err = stub.Transfer(toAddress, price_type, amount)

	if err != nil {
		return shim.Error("Error when making transfer。")
	}
	if _, ok := productionJSON.Supporters[supporter]; ok {

		productionJSON.Supporters[supporter] = Add(productionJSON.Supporters[supporter], price)
	} else {
		productionJSON.Supporters = make(map[string]string)
		productionJSON.Supporters[supporter] = price
	}
	productionJSONasBytes, err := json.Marshal(productionJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(production_key, productionJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("add supporter success."))
}

func (*ProductionInvoke) ListOfBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("list of buyer start.")
	username := args[0]
	product_type := args[1]
	product_serial := args[2]
	state_key := GetStateKey(ProductionPrefix, username, product_type, product_serial)
	resultsIterator, err := stub.GetStateByRange(state_key+StateStartSymbol, state_key+StateEndSymbol)
	if err != nil {
		return shim.Error(err.Error())
	}
	list, err := GetListResult(resultsIterator)
	if err != nil {
		return shim.Error("getListResult failed")
	}
	var listOfProductionJson []Production
	err = json.Unmarshal(list, &listOfProductionJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	result := make(map[string]string)
	for _, production := range listOfProductionJson {
		buyers := production.Buyers
		for k, v := range buyers {
			if _, ok := result[k]; ok {
				result[k] = Add(result[k], v)
			} else {
				result[k] = v
			}
		}
	}
	resultAsBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultAsBytes == nil {
		fmt.Println("This buyers doesn't exist: ")
		return shim.Error("This buyers doesn't exist: ")
	}
	return shim.Success(resultAsBytes)
}

func (w *ProductionInvoke) AddBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("add buyer start.")
	username := args[0]
	production_type := args[1]
	production_serial := args[2]
	price_type := args[3]
	buy_part := args[4]
	buyer := args[5]
	user_key := GetUserKey(username)
	// verify weather the user exists
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	}
	if userAsBytes == nil {
		fmt.Println("This user doesn't exist: " + username)
		return shim.Error("This user doesn't exist: " + username)
	}

	production_key := GetProductionKey(username, production_type, production_serial)
	// verify weather the production exists
	productionAsBytes, err := stub.GetState(production_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}
	if productionAsBytes == nil {
		fmt.Println("This production doesn't exist: " + production_key)
		return shim.Error("This production doesn't exist: " + production_key)
	}

	var productionJSON Production
	err = json.Unmarshal([]byte(productionAsBytes), &productionJSON)

	var userJSON User
	err = json.Unmarshal([]byte(userAsBytes), &userJSON)
	if userJSON.Username != productionJSON.Username {
		return shim.Error("The production's username doesn't correspond with the user's.")
	}

	if price_type != productionJSON.CopyrightPriceType {
		return shim.Error("The production's priceType doesn't correspond with the supporter's.")
	}

	// 是否满足购买条件
	total := "0"
	buyers := productionJSON.Buyers
	if buyers != nil {
		for _, value := range buyers {
			total = Add(total, value)
		}
	}
	num1, _ := strconv.Atoi(Add(total, buy_part))
	num2, _ := strconv.Atoi(productionJSON.CopyrightTransferPart)
	if num1 > num2 {
		return shim.Error("buy part too much")
	}

	// step 4: make transfer
	toAddress := userJSON.Address
	amount := big.NewInt(0)
	_, good := amount.SetString(Mul(productionJSON.CopyrightPrice, buy_part), 10)
	if !good {
		return shim.Error("Expecting integer value for amount")
	}

	fmt.Println(toAddress, ", ", price_type, "，", amount)
	err = stub.Transfer(toAddress, price_type, amount)
	if err != nil {
		return shim.Error("Error when making transfer。")
	}
	if _, ok := productionJSON.Buyers[buyer]; ok {
		productionJSON.Buyers[buyer] = Add(productionJSON.Buyers[buyer], buy_part)
	} else {
		productionJSON.Buyers = make(map[string]string)
		productionJSON.Buyers[buyer] = buy_part
	}
	productionJSONasBytes, err := json.Marshal(productionJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(production_key, productionJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("add supporter success."))
}

func (*ProductionInvoke) ModifyBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("modify buyer start.")
	// TODO
	return shim.Success([]byte("modify buyer success."))
}

func (*ProductionInvoke) DeleteBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("delete buyer start.")
	// TODO
	return shim.Success([]byte("delete buyer success."))
}
