package util

import (
	. "github.com/inklabsfoundation/inkchain/examples/creative/conf"
	. "github.com/inklabsfoundation/inkchain/examples/creative/model"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	"bytes"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"
	"errors"
)

func GetUserKey(username string) string {
	return GetStateKey(UserPrefix, username)
}

func GetArtistKey(username string) string {
	return GetStateKey(ArtistPrefix, username)
}

func GetProductionKey(username, production_type, serial string) string {
	return GetStateKey(ProductionPrefix, username, production_type, serial)
}

func GetStateKey(prefix string, args ...string) string {
	for _, value := range args {
		if value == "" {
			break
		} else {
			prefix = prefix + value + StateSplitSymbol
		}
	}
	return prefix
}

func Add(num_str_1, num_str_2 string) string {
	num1, _ := strconv.Atoi(num_str_1)
	num2, _ := strconv.Atoi(num_str_2)
	return strconv.Itoa(num1 + num2)
}

func Mul(num_str_1, num_str_2 string) string {
	num1, _ := strconv.Atoi(num_str_1)
	num2, _ := strconv.Atoi(num_str_2)
	return strconv.Itoa(num1 * num2)
}

func GetListResult(resultsIterator shim.StateQueryIteratorInterface) ([]byte, error) {
	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func GetListProduction(resultsIterator shim.StateQueryIteratorInterface) ([]Production, error) {
	defer resultsIterator.Close()
	var products []Production
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var pro Production
		json.Unmarshal(queryResponse.Value, &pro)
		products = append(products, pro)
	}
	return products, nil
}

func GetModifyProduction(productionJSON *Production, params []string) error {
	fmt.Println("params:", params)
	for _, param := range params {
		index := strings.Index(param, ",")
		if index == -1 {
			return errors.New("Wrong parameter format:" + param)
		}
		name := param[:index]
		value := param[index+1:]
		switch name {
		case "Name":
			productionJSON.Name = value
		case "Desc":
			productionJSON.Desc = value
		case "CopyrightPriceType":
			productionJSON.CopyrightPriceType = value
		case "CopyrightPrice":
			productionJSON.CopyrightPrice = value
		case "CopyrightNum":
			productionJSON.CopyrightNum = value
		case "CopyrightTransferPart":
			productionJSON.CopyrightTransferPart = value
		default:
			return errors.New("This field_name doesn't exist:" + value)
		}
	}
	return nil
}

func GetModifyArtist(artistJSON *Artist, params []string) error {
	fmt.Println("params:", params)
	for _, param := range params {
		index := strings.Index(param, ",")
		if index == -1 {
			return errors.New("Wrong parameter format:" + param)
		}
		name := param[:index]
		value := param[index+1:]
		switch name {
		case "Name":
			artistJSON.Name = value
		case "Desc":
			artistJSON.Desc = value
		default:
			return errors.New("This field_name doesn't exist:" + value)
		}
	}
	return nil
}

func GetModifyUser(userJSON *User, params []string) error {
	fmt.Println("params:", params)
	for _, param := range params {
		if param != "" {
			index := strings.Index(param, ",")
			if index == -1 {
				return errors.New("Wrong parameter format:" + param)
			}
			name := param[:index]
			value := param[index+1:]
			switch name {
			case "Email":
				userJSON.Email = value
			default:
				return errors.New("This field_name doesn't exist:" + value)
			}
		}
	}
	return nil
}
