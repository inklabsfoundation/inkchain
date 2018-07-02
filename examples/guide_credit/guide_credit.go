package main

import (
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	"fmt"
	"encoding/json"
	"strconv"
	"errors"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

type GuideCreditChainCode struct {
}

const (
	COMPANY_GUIDES = "CompanyGuides" //sender~guideAddress~timestamp
	BLACK_GUIDE    = "BlackGuide"    //companyAddress~guideAddress~timestamp
	OPERATE_LOG    = "OperateLog"    //sender~operateType~timestamp~md5key
	LEAVE_LOG      = "LeaveLog"      //guideAddress~companyAddress~timestamp
	WORK_LOG       = "WorkLog"       //guideAddress～companyAddress~timestamp
)

const (
	GUIDE_PREFIX   = "guide"   //prefix of guide key
	COMPANY_PREFIX = "company" //prefix of company key
)

const (
	OPERATE_TYPE_GUIDE   uint8 = 0
	OPERATE_TYPE_COMPANY uint8 = 1
)

//tourist guide info
type GuideInfo struct {
	Number          string `json:"number"`          //certificate of guide code
	Name            string `json:"name"`            // tourist guide name
	Sex             bool   `json:"sex"`             //tourist guide sex male= true ; female = false
	Age             int    `json:"age"`             //tourist guide age
	CompanyKey      string `json:"companyKey"`      //the company's key to which the tourist guide belongs
	JoinCompanyTime string `json:"joinCompanyTime"` //the time of tourist guide join company
	RegisterTime    string `json:"registerTime"`    //tourist guide register time
	BlackNumber     int    `json:"blackNumber"`     //count of company set guide to black list
}

//guide leave log
type LeaveLog struct {
	GuideKey   string `json:"guideKey"`   //key for guide state
	CompanyKey string `json:"companyKey"` //key for company
	Reason     string `json:"reason"`     //leave reason
	DateTime   string `json:"dateTime"`   //leave time
}

//work log
type WorkLog struct {
	GuideKey   string `json:"guideKey"`   //key for guide state
	CompanyKey string `json:"companyKey"` //key for company state
	DateTime   string `json:"dateTime"`   //join date
}

//company info
type CompanyInfo struct {
	Name         string `json:"name"`         //company name
	Address      string `json:"address"`      //company address
	Code         string `json:"code"`         // certificate of organization code
	TotalGuides  int    `json:"totalGuides"`  // total number of staff guides
	RegisterTime string `json:"registerTime"` //company register time
}

//black list info
type BlackInfo struct {
	CompanyCode string `json:"companyCode"` // company key
	GuideNumber string `json:"guideNumber"` //guide key
	CompanyName string `json:"companyName"` //company name
	Reason      string `json:"reason"`      //set to black list reason
	OperateTime string `json:"operateTime"` //operate time
}

//operate log
type OperateInfo struct {
	OperatorKey  string `json:"operatorKey"`  //operator key
	OperatorType uint8  `json:"operatorType"` //operator type : 0 Guide  1 Company
	OperateTime  string `json:"operateTime"`  //operate Time
	OperateMsg   string `json:"operatorMsg"`  //operate msg
}

const (
	RegisterCompany     = "registerCompany"
	RegisterGuide       = "registerGuide"
	AddGuide            = "addGuide"
	SetGuideToBlackList = "setGuideToBlackList"
	QueryGuideInfo      = "queryGuideInfo"
	QueryCompanyInfo    = "queryCompanyInfo"
	QueryOperateLog     = "queryOperateLog"
	QueryBlackList      = "queryBlackList"
	RemoveFromCompany   = "removeFromCompany"
	QueryGuideWorkList  = "queryGuideWorkList"
	QueryLeaveLogs      = "queryLeaveLogs"
)

func main() {
	err := shim.Start(new(GuideCreditChainCode))
	if err == nil {
		fmt.Printf("Error starting WorkChaincode: %s", err)
	}
}

func (g *GuideCreditChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (g *GuideCreditChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case RegisterCompany:
		if len(args) < 3 {
			return shim.Error("RegisterCompany, Incorrect number of arguments. Expecting 3")
		}
		return g.registerCompany(stub, args)
	case RegisterGuide:
		if len(args) < 4 {
			return shim.Error("RegisterGuide, Incorrect number of arguments. Expecting 4")
		}
		return g.registerGuide(stub, args)
	case AddGuide:
		if len(args) < 1 {
			return shim.Error("AddGuide, Incorrect number of arguments. Expecting 2")
		}
		return g.addGuide(stub, args)
	case SetGuideToBlackList:
		if len(args) < 2 {
			return shim.Error("SetGuideToBlackList, Incorrect number of arguments. Expecting 2")
		}
		return g.setGuideToBlackList(stub, args)
	case QueryGuideInfo:
		if len(args) < 1 {
			return shim.Error("QueryGuideInfo , Incorrect number of arguments. Excepting 1")
		}
		return g.queryGuideInfo(stub, args)
	case QueryCompanyInfo:
		if len(args) < 1 {
			return shim.Error("QueryCompany , Incorrect number of arguments. Excepting 1")
		}
		return g.queryCompanyInfo(stub, args)
	case QueryOperateLog:
		if len(args) < 2 {
			return shim.Error("QueryOperateLog , Incorrect number of arguments. Excepting 2")
		}
		return g.queryOperateLog(stub, args)
	case QueryBlackList:
		if len(args) < 1 {
			return shim.Error("QueryBlackList , Incorrect number of arguments, Excepting 1")
		}
		return g.queryBlackList(stub, args)
	case RemoveFromCompany:
		if len(args) < 2 {
			return shim.Error("RemoveFromCompany , Incorrect number of arguments, Excepting 2")
		}
		return g.removeFromCompany(stub, args)
	case QueryGuideWorkList:
		if len(args) < 1 {
			return shim.Error("QueryGuideWorkList , Incorrect number of arguments, Excepting 1")
		}
		return g.queryGuideWorkList(stub, args)
	case QueryLeaveLogs:
		if len(args) < 1 {
			return shim.Error("QueryLeaveLogs , Incorrect number of arguments, Excepting 1")
		}
		return g.queryGuideLeaveLogs(stub, args)
	}
	return shim.Error("Invalid call function name. Expecting \"registerCompany\" , \"RegisterGuide\" , \"AddGuide\" , \"SetGuideToBlackList\" , " +
		"\"QueryGuideInfo\" , \"QueryCompanyInfo\" , \"QueryOperateLog\" , \"QueryBlackList\" , \"RemoveFromCompany\" , \"QueryGuideWorkList\" , \"QueryLeaveLogs\".")
}

//register company info
//use sender address for key of company info
func (g *GuideCreditChainCode) registerCompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	registerTime := fmt.Sprintf("%d", timeStamp.GetSeconds())
	name := strings.TrimSpace(strings.ToLower(args[0]))
	if len(name) <= 0 {
		return shim.Error("1st arg must be a non-empty string")
	}
	address := strings.TrimSpace(strings.ToLower(args[1]))
	if len(address) <= 0 {
		return shim.Error("2st arg must be a non-empty string")
	}
	code := strings.TrimSpace(strings.ToLower(args[2]))
	if len(code) <= 0 {
		return shim.Error("3st arg must be a non-empty string")
	}
	company := &CompanyInfo{
		Name:         name,
		Address:      address,
		Code:         code,
		RegisterTime: registerTime,
	}
	sender, err := g.getSender(stub)
	if err != nil {
		return shim.Error("Failed to get sender info : " + err.Error())
	}
	key := COMPANY_PREFIX + sender
	//check sender had been registered
	old, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to check sender registration statue : " + err.Error())
	} else if old != nil {
		return shim.Error("Sender had been registered company")
	}
	//marshal company info
	companyJson, err := json.Marshal(company)
	if err != nil {
		return shim.Error("Failed to marshal company info to json : " + err.Error())
	}
	err = stub.PutState(key, companyJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	//save operate log
	logInfo := fmt.Sprintf("%s registered at %s", sender, registerTime)
	err = g.saveOperateLog(stub, OPERATE_TYPE_COMPANY, sender, logInfo, registerTime)
	if err != nil {
		return shim.Error("Failed to save operate log : " + err.Error())
	}
	return shim.Success(nil)
}

//register guide info
//use sender address for key of guide info
func (g *GuideCreditChainCode) registerGuide(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	registerTime := fmt.Sprintf("%d", timeStamp.GetSeconds())
	number := strings.TrimSpace(strings.ToLower(args[0]))
	if len(number) <= 0 {
		return shim.Error("1st arg must be a non-empty string")
	}
	name := strings.TrimSpace(strings.ToLower(args[1]))
	if len(name) <= 0 {
		return shim.Error("2st arg must be a non-empty string")
	}
	sex, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error("Failed to parse 2st arg to bool : " + err.Error())
	}
	age, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Failed to parse 3st arg to int : " + err.Error())
	}
	if age <= 0 {
		return shim.Error("4st arg must be greater than zero")
	}
	guide := &GuideInfo{
		Number:       number,
		Name:         name,
		Sex:          sex,
		Age:          age,
		RegisterTime: registerTime,
	}
	sender, err := g.getSender(stub)
	if err != nil {
		return shim.Error("Failed to get sender info : " + err.Error())
	}
	key := GUIDE_PREFIX + sender
	//check sender has been used
	old, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to check sender registration statue : " + err.Error())
	} else if old != nil {
		return shim.Error("sender had been registered guide")
	}
	//marshal guide info
	guideJson, err := json.Marshal(guide)
	if err != nil {
		return shim.Error("Failed to marshal guide info to json : " + err.Error())
	}
	err = stub.PutState(key, guideJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	//save operate log
	logInfo := fmt.Sprintf("%s registered at %s", sender, registerTime)
	err = g.saveOperateLog(stub, OPERATE_TYPE_GUIDE, sender, logInfo, registerTime)
	if err != nil {
		return shim.Error("Failed to save operate log : " + err.Error())
	}
	return shim.Success(nil)
}

//company add guide
func (g *GuideCreditChainCode) addGuide(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	//get sender's company info
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	}
	companyKey := COMPANY_PREFIX + sender
	company, err := g.getCompanyInfo(stub, companyKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	//get guide info
	guideAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(guideAddress) <= 0 {
		return shim.Error("1st arg must be a non-empty string")
	}
	guideKey := GUIDE_PREFIX + guideAddress
	guideInfo, err := g.getGuideInfo(stub, guideKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	//validate guide work info
	if len(guideInfo.CompanyKey) > 0 {
		return shim.Error("Guide has worked for another company")
	}
	//set guideInfo data
	operateTime := fmt.Sprintf("%d", timeStamp.GetSeconds())
	guideInfo.CompanyKey = sender
	guideInfo.JoinCompanyTime = operateTime
	guideJson, err := json.Marshal(guideInfo)
	if err != nil {
		return shim.Error("Failed to marshal guide info to json : " + err.Error())
	}
	company.TotalGuides += 1
	companyJson, err := json.Marshal(company)
	if err != nil {
		return shim.Error("Failed to marshal company info to json : " + err.Error())
	}
	err = stub.PutState(guideKey, guideJson)
	if err != nil {
		return shim.Error("Failed to save guide info : " + err.Error())
	}
	//update company info
	err = stub.PutState(companyKey, companyJson)
	if err != nil {
		return shim.Error("Failed to save company info : " + err.Error())
	}
	err = g.companyGuideByComposite(stub, []string{sender, guideAddress, operateTime})
	if err != nil {
		return shim.Error("Failed to save guide company info : " + err.Error())
	}
	//save work log
	err = g.WorkLogByComposite(stub, sender, guideAddress, operateTime)
	if err != nil {
		return shim.Error("Failed to save guide work log : " + err.Error())
	}
	//save operate log
	logInfo := fmt.Sprintf("%s add guide %s to company at %s", sender, guideAddress, operateTime)
	err = g.saveOperateLog(stub, OPERATE_TYPE_COMPANY, sender, logInfo, operateTime)
	if err != nil {
		return shim.Error("Failed to save operate log : " + err.Error())
	}
	return shim.Success(nil)
}

//set guide worked for sender's company to black list
func (g *GuideCreditChainCode) setGuideToBlackList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	operateTime := fmt.Sprintf("%d", timeStamp.Seconds)
	//get sender's company info
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error(err.Error())
	}
	companyKey := COMPANY_PREFIX + sender
	company, err := g.getCompanyInfo(stub, companyKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	//get guide info
	guideAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(guideAddress) <= 0 {
		return shim.Error("1st arg must be a non-empty string")
	}
	guideKey := GUIDE_PREFIX + guideAddress
	guideInfo, err := g.getGuideInfo(stub, guideKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	//validate guide work info
	if len(guideInfo.CompanyKey) <= 0 {
		return shim.Error("Guide has not worked for your company")
	}
	reason := strings.TrimSpace(strings.ToLower(args[1]))
	if len(reason) <= 0 {
		return shim.Error("1st arg must be a non-empty string")
	}
	blackInfo := &BlackInfo{
		CompanyCode: sender,
		CompanyName: company.Name,
		GuideNumber: guideAddress,
		Reason:      reason,
		OperateTime: operateTime,
	}
	blackInfoJson, err := json.Marshal(blackInfo)
	if err != nil {
		return shim.Error("Failed marshal black info to json : " + err.Error())
	}
	guideInfo.BlackNumber += 1
	guideInfoJson, err := json.Marshal(guideInfo)
	if err != nil {
		return shim.Error("Failed marshal guide info to json : " + err.Error())
	}
	err = stub.PutState(guideKey, guideInfoJson)
	err = g.blackListByComposite(stub, []string{guideAddress, sender, operateTime}, blackInfoJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	//save operate log
	logInfo := fmt.Sprintf("%s set guide %s to black list at %s", sender, guideAddress, operateTime)
	err = g.saveOperateLog(stub, OPERATE_TYPE_COMPANY, sender, logInfo, operateTime)
	if err != nil {
		return shim.Error("Failed to save operate log : " + err.Error())
	}
	return shim.Success(nil)
}

//remove guide from company
func (g *GuideCreditChainCode) removeFromCompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	operateTime := fmt.Sprintf("%d", timeStamp.Seconds)
	//check sender company info
	sender, err := g.getSender(stub)
	if err != nil {
		return shim.Error("Failed to get sender : " + err.Error())
	}
	companyKey := COMPANY_PREFIX + sender
	company, err := g.getCompanyInfo(stub, companyKey)
	if err != nil {
		return shim.Error("Failed to get company info : " + err.Error())
	}
	//check guide info
	guideAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(guideAddress) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	guideKey := GUIDE_PREFIX + guideAddress
	guide, err := g.getGuideInfo(stub, guideKey)
	if err != nil {
		return shim.Error("Failed to get guide info : " + err.Error())
	}
	if guide.CompanyKey != sender {
		return shim.Error("Guide not worked for you")
	}
	//check 2st arg
	reason := strings.TrimSpace(strings.ToLower(args[1]))
	if len(reason) <= 0 {
		return shim.Error("2st arg must be non-empty string")
	}
	//make write info
	company.TotalGuides -= 1
	guide.CompanyKey = ""
	leaveLog := LeaveLog{
		CompanyKey: sender,
		GuideKey:   guideAddress,
		Reason:     reason,
		DateTime:   operateTime,
	}
	fmt.Println(guide)
	//start to write info
	companyJson, err := json.Marshal(company)
	if err != nil {
		return shim.Error("Failed to marshal company info to json : " + err.Error())
	}
	err = stub.PutState(companyKey, companyJson)
	if err != nil {
		return shim.Error("Failed to update company info : " + err.Error())
	}
	guideJson, err := json.Marshal(guide)
	if err != nil {
		return shim.Error("Failed to marshal guide info to json : " + err.Error())
	}
	err = stub.PutState(guideKey, guideJson)
	if err != nil {
		return shim.Error("Failed to update guide info : " + err.Error())
	}
	leaveLogJson, err := json.Marshal(leaveLog)
	indexKey, err := stub.CreateCompositeKey(LEAVE_LOG, []string{sender, guideAddress, operateTime})
	if err != nil {
		return shim.Error("Failed to create key for leave log : " + err.Error())
	}
	err = stub.PutState(indexKey, leaveLogJson)
	if err != nil {
		return shim.Error("Failed to create leave log : " + err.Error())
	}
	//save operate log
	logInfo := fmt.Sprintf("%s remove the %s from company at %s beacuse of  %s", sender, guideAddress, operateTime, reason)
	g.saveOperateLog(stub, OPERATE_TYPE_COMPANY, sender, logInfo, operateTime)
	if err != nil {
		return shim.Error("Failed to save operate log : " + err.Error())
	}
	return shim.Success(nil)
}

//save user operate log
func (g *GuideCreditChainCode) saveOperateLog(stub shim.ChaincodeStubInterface, operateType uint8, sender, msg, timestamp string) error {
	log := &OperateInfo{
		OperatorKey:  sender,
		OperatorType: operateType,
		OperateTime:  timestamp,
		OperateMsg:   msg,
	}
	logJson, err := json.Marshal(log)
	if err != nil {
		return errors.New("Can not marshal log to json : " + err.Error())
	}
	key := g.calcMd5Key(logJson)
	indexKey, err := stub.CreateCompositeKey(OPERATE_LOG, []string{sender, fmt.Sprintf("%d", operateType), timestamp, key})
	if err != nil {
		return err
	}
	err = stub.PutState(indexKey, logJson)
	if err != nil {
		return err
	}
	return nil
}

//query guide info
func (g *GuideCreditChainCode) queryGuideInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	guideAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(guideAddress) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	guide, err := stub.GetState(GUIDE_PREFIX + guideAddress)
	if err != nil {
		return shim.Error("Failed to get guide info : " + err.Error())
	} else if guide == nil {
		return shim.Error("Failed to get guide info : Can not find the guide info")
	}
	blackList, err := g.getBlackList(stub, guideAddress)
	if err != nil {
		return shim.Error("Failed to get guide's black list : " + err.Error())
	}
	guideInfo := &GuideInfo{}
	err = json.Unmarshal(guide, guideInfo)
	if err != nil {
		return shim.Error("Failed to unmarshal guide info from json to struct : " + err.Error())
	}
	result := map[string]interface{}{
		"guide":     guideInfo,
		"blackList": blackList,
	}
	resJson, err := json.Marshal(result)
	if err != nil {
		return shim.Error("Failed to marshal result to json : " + err.Error())
	}
	return shim.Success(resJson)
}

//query guide work logs
func (g *GuideCreditChainCode) queryGuideWorkList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	guideAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(guideAddress) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	//get blackList
	resultsIterator, err := stub.GetStateByPartialCompositeKey(WORK_LOG, []string{guideAddress})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	resultList := make([]WorkLog, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		value := responseRange.Value
		tmp := WorkLog{}
		err = json.Unmarshal(value, &tmp)
		if err == nil {
			resultList = append(resultList, tmp)
		}
	}
	resJson, err := json.Marshal(resultList)
	if err != nil {
		return shim.Error("Failed to marshal result to json : " + err.Error())
	}
	return shim.Success(resJson)
}

//query guide leave logs
func (g *GuideCreditChainCode) queryGuideLeaveLogs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	guideAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(guideAddress) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	//get blackList
	resultsIterator, err := stub.GetStateByPartialCompositeKey(LEAVE_LOG, []string{guideAddress})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	resultList := make([]LeaveLog, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		value := responseRange.Value
		tmp := LeaveLog{}
		err = json.Unmarshal(value, &tmp)
		if err == nil {
			resultList = append(resultList, tmp)
		}
	}
	resJson, err := json.Marshal(resultList)
	if err != nil {
		return shim.Error("Failed to marshal result to json : " + err.Error())
	}
	return shim.Success(resJson)
}

//query company info
func (g *GuideCreditChainCode) queryCompanyInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	companyAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(companyAddress) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	company, err := stub.GetState(COMPANY_PREFIX + companyAddress)
	if err != nil {
		return shim.Error("Failed to get company info : " + err.Error())
	} else if company == nil {
		return shim.Error("Failed to get company info : Can not find the company info")
	}
	return shim.Success(company)
}

//query black list of guide
func (g *GuideCreditChainCode) queryBlackList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get operator address
	guideAddress := strings.TrimSpace(strings.ToLower(args[0]))
	if len(guideAddress) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	//get blackList
	resultsIterator, err := stub.GetStateByPartialCompositeKey(BLACK_GUIDE, []string{guideAddress})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	resultList := make([]BlackInfo, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		value := responseRange.Value
		tmp := BlackInfo{}
		err = json.Unmarshal(value, &tmp)
		if err == nil {
			resultList = append(resultList, tmp)
		}
	}
	resJson, err := json.Marshal(resultList)
	if err != nil {
		return shim.Error("Failed to marshal result to json : " + err.Error())
	}
	return shim.Success(resJson)
}

//query operate log
func (g *GuideCreditChainCode) queryOperateLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get operator address
	operator := strings.TrimSpace(strings.ToLower(args[0]))
	if len(operator) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	operatorType, err := strconv.ParseUint(args[1], 10, 8)
	if err != nil {
		return shim.Error("Failed to paras 2st arg to int : " + err.Error())
	}
	if uint8(operatorType) != OPERATE_TYPE_GUIDE && uint8(operatorType) != OPERATE_TYPE_COMPANY {
		return shim.Error("2st must be 0 or 1")
	}
	resultsIterator, err := stub.GetStateByPartialCompositeKey(OPERATE_LOG, []string{operator, fmt.Sprintf("%d", operatorType)})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	resultList := make([]OperateInfo, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		value := responseRange.Value
		tmp := OperateInfo{}
		err = json.Unmarshal(value, &tmp)
		if err == nil {
			resultList = append(resultList, tmp)
		}
	}
	resJson, err := json.Marshal(resultList)
	if err != nil {
		return shim.Error("Failed to marshal result to json : " + err.Error())
	}
	return shim.Success(resJson)
}

//create compositeKey of CompanyGuide
func (g *GuideCreditChainCode) companyGuideByComposite(stub shim.ChaincodeStubInterface, args []string) error {
	var indexKey string
	var err error
	indexKey, err = stub.CreateCompositeKey(COMPANY_GUIDES, args)
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

//create compositeKey of BlackInfo
func (g *GuideCreditChainCode) blackListByComposite(stub shim.ChaincodeStubInterface, args []string, blackInfo []byte) error {
	var indexKey string
	var err error
	indexKey, err = stub.CreateCompositeKey(BLACK_GUIDE, args)
	if err != nil {
		return err
	}
	err = stub.PutState(indexKey, blackInfo)
	if err != nil {
		return err
	}
	return nil
}

//save work log by composite key
func (g *GuideCreditChainCode) WorkLogByComposite(stub shim.ChaincodeStubInterface, sender, guide, operateTime string) error {
	var indexKey string
	var err error
	indexKey, err = stub.CreateCompositeKey(WORK_LOG, []string{sender, guide, operateTime})
	if err != nil {
		return err
	}
	workLog := WorkLog{
		CompanyKey: sender,
		GuideKey:   guide,
		DateTime:   operateTime,
	}
	workLogJson, err := json.Marshal(workLog)
	if err != nil {
		return err
	}
	err = stub.PutState(indexKey, workLogJson)
	if err != nil {
		return err
	}
	return nil
}

//get black list by guide address
func (g *GuideCreditChainCode) getBlackList(stub shim.ChaincodeStubInterface, guideAddress string) ([]BlackInfo, error) {
	//get operator address
	if len(guideAddress) <= 0 {
		return nil, errors.New("address for guide arg must be non-empty string")
	}
	//get blackList
	resultsIterator, err := stub.GetStateByPartialCompositeKey(BLACK_GUIDE, []string{guideAddress})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	resultList := make([]BlackInfo, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		value := responseRange.Value
		tmp := BlackInfo{}
		err = json.Unmarshal(value, &tmp)
		if err == nil {
			resultList = append(resultList, tmp)
		}
	}
	return resultList, nil
}

//get guide info
//args with ChaincodeStubInterface and guide address
func (g *GuideCreditChainCode) getGuideInfo(stub shim.ChaincodeStubInterface, key string) (*GuideInfo, error) {
	guide, err := stub.GetState(key)
	if err != nil {
		return nil, err
	} else if guide == nil {
		return nil, errors.New("Can not find the guide info")
	}
	guideInfo := &GuideInfo{}
	err = json.Unmarshal(guide, guideInfo)
	if err != nil {
		return nil, err
	}
	return guideInfo, nil
}

//get company info
//args with ChaincodeStubInterface and company address
func (g *GuideCreditChainCode) getCompanyInfo(stub shim.ChaincodeStubInterface, key string) (*CompanyInfo, error) {
	company, err := stub.GetState(key)
	if err != nil {
		return nil, err
	} else if company == nil {
		return nil, errors.New("Can not find the company info")
	}
	companyInfo := &CompanyInfo{}
	err = json.Unmarshal(company, companyInfo)
	if err != nil {
		return nil, err
	}
	return companyInfo, nil
}

//get sender
func (g *GuideCreditChainCode) getSender(stub shim.ChaincodeStubInterface) (string, error) {
	sender, err := stub.GetSender()
	if err != nil {
		return "", err
	} else if sender == "" {
		return "", errors.New("Can not get sender")
	}
	return sender, nil
}

//calc key of BlackInfo
func (g *GuideCreditChainCode) calcMd5Key(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}