package main

import (
	"fmt"
	"strconv"
	. "github.com/inklabsfoundation/inkchain/examples/creative/conf"
	. "github.com/inklabsfoundation/inkchain/examples/creative/invoke"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
)

type mainChaincode struct {
	UserInvoke
	ArtistInvoke
	ProductionInvoke
}

func main() {
	err := shim.Start(&mainChaincode{})
	if err != nil {
		fmt.Printf("Error starting assetChaincode: %s", err)
	}
}

func (c *mainChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mainChaincode Init.")
	return shim.Success([]byte("Init success."))
}

func (c *mainChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mainChaincode Invoke.")
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case AddUser:
		parameter_length := 2
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length) + ":" + strconv.Itoa(len(args)))
		}
		return c.AddUser(stub, args)
	case ModifyUser:
		parameter_length := 2
		if len(args) <= parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ModifyUser(stub, args)
	case DeleteUser:
		parameter_length := 1
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.DeleteUser(stub, args)
	case QueryUser:
		parameter_length := 1
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.QueryUser(stub, args)
	case ListOfUser:
		parameter_length := 0
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ListOfUser(stub, args)
	case AddArtist:
		parameter_length := 3
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.AddArtist(stub, args)
	case DeleteArtist:
		parameter_length := 1
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.DeleteArtist(stub, args)
	case ModifyArtist:
		parameter_length := 1
		if len(args) <= parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ModifyArtist(stub, args)
	case QueryArtist:
		parameter_length := 1
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.QueryArtist(stub, args)
	case ListOfArtist:
		parameter_length := 0
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ListOfArtist(stub, args)
	case AddProduction:
		parameter_length := 9
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.AddProduction(stub, args)
	case DeleteProduction:
		parameter_length := 3
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.DeleteProduction(stub, args)
	case ModifyProduction:
		parameter_length := 3
		if len(args) <= parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ModifyProduction(stub, args)
	case QueryProduction:
		parameter_length := 3
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.QueryProduction(stub, args)
	case ListOfProduction:
		parameter_length := 2
		if len(args) > parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ListOfProduction(stub, args)
	case ListOfSupporter:
		parameter_length := 3
		if len(args) <= 1 && len(args) > parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ListOfSupporter(stub, args)
	case ListOfBuyer:
		parameter_length := 3
		if len(args) <= 1 && len(args) > parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.ListOfBuyer(stub, args)
	case AddSupporter:
		parameter_length := 6
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.AddSupporter(stub, args)
	case AddBuyer:
		parameter_length := 6
		if len(args) != parameter_length {
			return shim.Error("Incorrect number of arguments. Expecting " + strconv.Itoa(parameter_length))
		}
		return c.AddBuyer(stub, args)
	case ModifyBuyer:
		// TODO
		return c.ModifyBuyer(stub, args)
	case DeleteBuyer:
		// TODO
		return c.DeleteBuyer(stub, args)
	}
	return shim.Error("Invalid invoke function name.")
}
