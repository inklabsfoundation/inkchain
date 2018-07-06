/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package chaincode

import (
	"encoding/json"
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"
	"errors"
	"strconv"
	"math/big"
	"strings"
	"github.com/inklabsfoundation/inkchain/core/wallet"
)

//struct for response from eth JSON_RPC
type ethBlockRes struct {
	Id      int           `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Result  *ethBlockInfo `json:"result"`
}
type ethTransRes struct {
	Id      int             `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Result  *ethTranRecInfo `json:"result"`
}

type ethBlockInfo struct {
	Number     string `json:"number"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Sha3Uncles string `json:"sha3Uncles"`
	Size       string `json:"size"`
}

type ethTranRecInfo struct {
	TransactionHash  string       `json:"transactionHash"`
	TransactionIndex string       `json:"transactionIndex"`
	ContractAddress  string       `json:"contractAddress"`
	BlockNumber      string       `json:"blockNumber"`
	ToContract       string       `json:"to"`
	BlockHash        string       `json:"blockHash"`
	Logs             []ethTranLog `json:"logs"`
}

type ethTranLog struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

//struct for response from qtum-insight-api
type qtumTransInfo struct {
	BlockHash       string    `json:"blockHash"`
	BlockNumber     int       `json:"blockNumber"`
	ContractAddress string    `json:"contractAddress"`
	Log             []qtumLog `json:"log"`
}

type qtumLog struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

type qtumBlockInfo struct {
	Hash    string `json:"hash"`
	Height  int    `json:"height"`
	Size    int    `json:"size"`
	Version int    `json:"version"`
}

const (
	CONFIRM_BLOCK_NUM = 6
	ETH_FIELD_LENGTH  = 64 //log field's length in eth
	QTUM_FIELD_LENGTH = 64 //log field's length in qtum
	ETH_PLACEHOLDER = "0"//placeholder in eth log
	QTUM_PLACTHOLDER = "0"//placeholder in qtum log
	ETH_HEADER_STR = "0x"//binary sign in eth api
	QTUM_API_URL_PREFIX = "/qtum-insight-api"
)

//validate pubTxId from eth
func (handler *Handler) validateEthPubTxId(pubTxId string, toUser string, amount *big.Int) (result bool, balanceType string, err error) {
	publicNode, ok := wallet.PublicInfos["eth"]
	if !ok {
		err = errors.New("platform not support")
		return
	}
	url := publicNode.FullNodeIp
	localPlatform := wallet.LocalPlatform
	contractList := publicNode.ContractList
	if url == "" {
		err = errors.New("not support this coin or public chain")
		return
	}
	//get eth transaction detail
	transInfo, err := getEthTransInfo(url, pubTxId)
	if err != nil {
		return
	}
	if len(transInfo.Logs) < 2 || len(transInfo.Logs[0].Topics) < 3 {
		err = errors.New("transaction verified failed")
		return
	}
	if transInfo.BlockNumber == "" {
		err = errors.New("transaction not confirmed")
	}
	valueData := strings.TrimLeft(transInfo.Logs[1].Data[2+2*ETH_FIELD_LENGTH:2+3*ETH_FIELD_LENGTH], ETH_PLACEHOLDER)
	value, err := strconv.ParseInt(valueData, 16, 64)
	if err != nil {
		return
	}
	valueInt := big.NewInt(value)
	if amount.Cmp(valueInt) < 0 {
		err = errors.New("transaction amount error")
		return
	}
	balanceType = ""
	for coinType, contract := range contractList {
		if transInfo.ToContract == contract.Address {
			balanceType = coinType
			break
		}
	}
	if balanceType == "" {
		err = errors.New("transaction data verified failed")
		return
	}

	platformData := strings.TrimRight(transInfo.Logs[1].Data[:ETH_FIELD_LENGTH+2], ETH_PLACEHOLDER)
	if len(platformData)%2 != 0 {
		platformData = platformData + "0"
	}

	platform, err := asciiToString(platformData[2:])
	if err != nil {
		return
	}
	if platform != localPlatform {
		err = errors.New("transaction platform error")
		return
	}
	toUserData := transInfo.Logs[1].Data[2+ETH_FIELD_LENGTH : 2+2*ETH_FIELD_LENGTH]
	if !strings.Contains(toUserData, toUser[len(wallet.ADDRESS_PREFIX):]) {
		err = errors.New("transaction turn out account error")
		return
	}
	coinNameAsc := strings.TrimRight(transInfo.Logs[1].Data[2+3*ETH_FIELD_LENGTH:], ETH_PLACEHOLDER)
	if len(coinNameAsc)%2 != 0 {
		coinNameAsc = coinNameAsc + "0"
	}
	coinName, err := asciiToString(coinNameAsc)
	if err != nil {
		return
	}
	if coinName != balanceType {
		err = errors.New("transaction coin type validate failed")
		return
	}
	//calculate block index
	blockNum, err := strconv.ParseInt(transInfo.BlockNumber[2:], 16, 64)
	if err != nil {
		return
	}

	blockNum = CONFIRM_BLOCK_NUM + blockNum
	// get block info to validate transaction confirmed
	_, err = getEthBlockInfo(url, blockNum)
	if err != nil {
		return
	} else {
		result = true
	}
	return
}

//validate pubTxId from qtum
func (handler *Handler) validateQtumPubTxId(pubTxId string, toUser string, amount *big.Int) (result bool, balanceType string, err error) {
	publicNode, ok := wallet.PublicInfos["qtum"]
	if !ok {
		err = errors.New("platform not support")
		return
	}
	url := publicNode.FullNodeIp
	localPlatform := wallet.LocalPlatform
	contractList := publicNode.ContractList
	if url == "" {
		err = errors.New("not support this coin or public chain")
		return
	}
	//get transaction detail
	transInfo, err := getQtumTransInfo(url, pubTxId)
	if err != nil {
		return
	}
	if len(transInfo.Log) < 2 {
		err = errors.New("transaction not belong to our contract")
		return
	}
	valueData := strings.TrimLeft(transInfo.Log[1].Data[2*QTUM_FIELD_LENGTH:3*QTUM_FIELD_LENGTH], QTUM_PLACTHOLDER)
	value, err := strconv.ParseInt(valueData, 16, 64)
	if err != nil {
		return
	}
	valueInt := big.NewInt(value)
	if amount.Cmp(valueInt) < 0 {
		err = errors.New("transaction amount error")
		return
	}
	platformData := strings.TrimRight(transInfo.Log[1].Data[:QTUM_FIELD_LENGTH], QTUM_PLACTHOLDER)
	if len(platformData)%2 != 0 {
		platformData = platformData + "0"
	}
	platform, err := asciiToString(platformData)
	if err != nil {
		return
	}
	if platform != localPlatform {
		err = errors.New("transaction platform error")
		return
	}
	toUserData := transInfo.Log[1].Data[QTUM_FIELD_LENGTH : 2*QTUM_FIELD_LENGTH]
	if !strings.Contains(toUserData, toUser[len(wallet.ADDRESS_PREFIX):]) {
		err = errors.New("transaction turn out account error")
		return
	}
	balanceType = ""
	for coinType, contract := range contractList {
		if transInfo.ContractAddress == contract.Address {
			balanceType = coinType
			break
		}
	}
	if balanceType == "" {
		err = errors.New("transaction data verified failed")
		return
	}
	coinNameAsc := strings.TrimRight(transInfo.Log[1].Data[3*QTUM_FIELD_LENGTH:], QTUM_PLACTHOLDER)
	if len(coinNameAsc)%2 != 0 {
		coinNameAsc = coinNameAsc + "0"
	}
	coinName, err := asciiToString(coinNameAsc)
	if err != nil {
		return
	}
	if coinName != balanceType {
		err = errors.New("transaction coin type validate failed")
		return
	}
	//get block detail by block hash
	blockInfo, err := getQtumBlockInfo(url, transInfo.BlockHash)
	if err != nil {
		return
	}
	//get block hash by index
	confirmBlock, err := getQtumBlockHashByHeight(url, blockInfo.Height+CONFIRM_BLOCK_NUM)
	if err != nil {
		return
	}
	//determined whether the blockHash field exists
	if _, ok := confirmBlock["blockHash"]; !ok {
		err = errors.New("transaction not confirmed by public chain")
		return
	} else {
		result = true
	}
	return
}

//get transaction detail from eth
func getEthTransInfo(url string, pubTxId string) (*ethTranRecInfo, error) {
	reqParam := map[string]interface{}{"id": 67, "method": "eth_getTransactionReceipt", "params": []string{pubTxId}}
	res, err := quest("POST", url, reqParam)
	if err != nil {
		return nil, err
	}
	data := ethTransRes{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}
	if data.Result == nil {
		return nil, errors.New("txId not found")
	}
	return data.Result, nil
}

//get block detail from eth
func getEthBlockInfo(url string, number int64) (*ethBlockInfo, error) {
	numHex := fmt.Sprintf(ETH_HEADER_STR+"%x", number)
	reqParam := map[string]interface{}{"id": 67, "method": "eth_getBlockByNumber", "params": []interface{}{numHex, true}}
	res, err := quest("POST", url, reqParam)
	if err != nil {
		return nil, err
	}
	data := ethBlockRes{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}
	if data.Result == nil {
		err = errors.New("block not exists")
		return nil, err
	}
	return data.Result, nil
}

//get transaction from qtum
func getQtumTransInfo(url string, pubTxId string) (*qtumTransInfo, error) {
	url = url + QTUM_API_URL_PREFIX+"/txs/" + pubTxId + "/receipt"
	res, err := quest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	var datas []qtumTransInfo
	err = json.Unmarshal(res, &datas)
	if err != nil {
		return nil, errors.New("transaction not exists")
	}
	if datas == nil {
		return nil, errors.New("transaction not exists")
	}
	data := datas[0]
	if data.BlockHash == "" {
		return nil, errors.New("transaction not confirmed")
	}
	return &data, nil
}

//get block detail from qtum
func getQtumBlockInfo(url string, blockHash string) (*qtumBlockInfo, error) {
	url = url + QTUM_API_URL_PREFIX+"/block/" + blockHash
	res, err := quest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	data := qtumBlockInfo{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, errors.New("block not existed")
	}
	return &data, nil
}

//get block hash from qtum
func getQtumBlockHashByHeight(url string, height int) (map[string]interface{}, error) {
	url = fmt.Sprintf("%s"+QTUM_API_URL_PREFIX+"/block-index/%d", url, height)
	res, err := quest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, errors.New("confirm block not existed")
	}
	return data, nil
}

//do request
func quest(method string, url string, params map[string]interface{}) ([]byte, error) {

	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "http://" + url
	}
	method = strings.ToUpper(method)

	var req *http.Request
	var err error
	if params == nil || len(params) <= 0 {
		req, err = http.NewRequest(method, url, nil)
	} else {
		reqJson, _ := json.Marshal(params)
		paramData := bytes.NewReader(reqJson)
		req, err = http.NewRequest(method, url, paramData)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		return body, nil
	}
	return nil, err
}

func asciiToString(str string) (s string, err error) {
	if len(str)%2 != 0 {
		err = errors.New("str validate error")
		return
	}
	for i := 0; i <= len(str)-2; i = i + 2 {
		tmp := ""
		if i == len(str)-2 {
			tmp = str[i:]
		} else {
			tmp = str[i : i+2]
		}
		tmpInt, err := strconv.ParseInt(tmp, 16, 64)
		if err != nil {
			s = ""
			break
		}
		str := string(rune(tmpInt))
		s = s + str
	}
	return
}
