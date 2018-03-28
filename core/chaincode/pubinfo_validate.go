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

//struct for response which from eth JSON_RPC
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
	TransactionHash  string                 `json:"transactionHash"`
	TransactionIndex string                 `json:"transactionIndex"`
	BlockNumber      string                 `json:"blockNumber"`
	BlockHash        string                 `json:"blockHash"`
	GasUsed          string                 `json:"gasUsed"`
	Logs             map[string]interface{} `json:"logs"`
}

//struct for response which from qtum-insight-api
type qtumTransInfo struct {
	TxId      string `json:"txid"`
	BlockHash string `json:"blockhash"`
	//Vin         map[string]interface{} `json:"vin"`
	//Vout        map[string]interface{} `json:"vout"`
	BlockHeight int `json:"blockheight"`
	Time        int `json:"time"`
}

type qtumBlockInfo struct {
	Hash    string `json:"hash"`
	Height  int    `json:"height"`
	Size    int    `json:"size"`
	Version int    `json:"version"`
	//Tx      map[string]string `json:"tx"`
}

//validate pubTxId from eth
func (handler *Handler) validateEthTrans(pubTxId string, amount *big.Int) (result bool, err error) {
	url := wallet.FullNodeIps["eth"]
	if url == "" {
		err = errors.New("not support this public chain")
		return
	}
	//get eth transaction detail
	transInfo, err := getEthTransInfo(url, pubTxId)
	if err != nil {
		return
	}
	//calculate block index
	blockNum, err := strconv.ParseInt(transInfo.BlockNumber, 16, 64)
	blockNum = 6 + blockNum
	if err != nil {
		return
	}
	//get block info to validate transaction confirmed
	_, err = getEthBlockInfo(url, blockNum)
	if err != nil {
		return
	} else {
		result = true
	}
	return
}

//validate pubTxId from qtum
func (handler *Handler) validateQtumPubTxId(pubTxId string, amount *big.Int) (result bool, err error) {
	url := wallet.FullNodeIps["qtum"]
	if url == "" {
		err = errors.New("not support this public chain")
		return
	}
	//get transaction detail
	transInfo, err := getQtumTransInfo(url, pubTxId)
	if err != nil {
		return
	}
	//get block detail by block hash
	blockInfo, err := getQtumBlockInfo(url, transInfo.BlockHash)
	if err != nil {
		return
	}
	//get block hash by index
	confirmBlock, err := getQtumBlockHashByHeight(url, blockInfo.Height+6)
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
		err = errors.New("txId not found")
		return nil, err
	}
	return data.Result, nil
}

//get block detail from eth
func getEthBlockInfo(url string, number int64) (*ethBlockInfo, error) {
	numHex := fmt.Sprintf("0x%x", number)
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
	url = url + "/qtum-insight-api/tx/" + pubTxId
	res, err := quest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	data := qtumTransInfo{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}
	if data.BlockHash == "" {
		err = errors.New("transaction not confirmed")
		return nil, err
	}
	return &data, nil
}

//get block detail from qtum
func getQtumBlockInfo(url string, blockHash string) (*qtumBlockInfo, error) {
	url = url + "/qtum-insight-api/block/" + blockHash
	res, err := quest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	data := qtumBlockInfo{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

//get block hash from qtum
func getQtumBlockHashByHeight(url string, height int) (map[string]interface{}, error) {
	url = fmt.Sprintf("%s/qtum-insight-api/block-index/%d", url, height)
	res, err := quest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
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
