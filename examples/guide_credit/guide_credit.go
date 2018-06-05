package main

import (
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"inkchain/core/chaincode/shim"
	"fmt"
)

type GuideCreditChainCode struct {
}

//tourist guide info
type GuideInfo struct {
	Id              string `json:"id"`              //tourist guide number
	Name            string `json:"name"`            // tourist guide name
	Sex             string `json:"sex"`             //tourist guide sex
	Age             string `json:"age"`             //tourist guide age
	CompanyId       string `json:"companyId"`       //the company's id to which the tourist guide belongs
	JoinCompanyTime string `json:"joinCompanyTime"` //the time of tourist guide join company
	RegisterTime    string `json:"registerTime"`    //tourist guide register time
}

//company info
type CompanyInfo struct {
	Id           int    `json:"id"`           //company id
	Name         string `json:"name"`         //company name
	Address      string `json:"address"`      //company address
	Code         string `json:"code"`         // certificate of organization code
	TotalGuides  int    `json:"totalGuides"`  // total number of staff guides
	RegisterTime string `json:"registerTime"` //company register time
}

type CreditInfo struct {
	GuideId     int     `json:"guideId"`
	Number      float64 `json:"number"`
	BlackNumber int     `json:"blackNumber"`
}

type BlackInfo struct {
	CompanyId int `json:"companyId"`
}

type OperateInfo struct {
	OperatorId   int    `json:"operatorId"`   //operator Id
	OperatorType int    `json:"operatorType"` //operator type : 0 Guide  1 Company
	Address      string `json:"address"`      //operator address
	OperateTime  string `json:"operateTime"`  //operate Time
	OperateMsg   string `json:"operatorMsg"`  //operate msg
}

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
	return shim.Error("")
}
