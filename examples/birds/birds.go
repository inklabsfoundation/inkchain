package main

import (
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"fmt"
	"strings"
	"encoding/json"
)

/**
private_key 07caf88941eafcaaa3370657fccc261acb75dfba
token 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4


private_key a5ff00eb44bf19d5dfbde501c90e286badb58df4
token 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5
**/
type BirdsChainCode struct{}

const OBJECT_TYPE_BIRD = "bird"
const (
	InitBird = "initBird"
	GetBird  = "getBird"
	DelBird  = "deleteBird"
	EditBird = "editBird"
)

type birds struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	Color      string `json:"color"`
	Owner      string `json:"owner"`
}

func (b *BirdsChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (b *BirdsChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	switch function {
	case InitBird:
		return b.initBird(stub, args)
	case GetBird:
		return b.getBird(stub, args)
	case DelBird:
		return b.deleteBird(stub, args)
	case EditBird:
		return b.editBird(stub, args)
	default:
		return shim.Success(nil)
	}
}

//init a new birds
func (b *BirdsChainCode) initBird(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- Start init bird")
	if len(args) < 4 {
		return shim.Error("Incorrect of arguments,there must be 4 arguments,you may need send one of these arguments(name/kind/color/owner)")
	}

	if len(args[0]) <= 0 || len(args[1]) <= 0 || len(args[2]) <= 0 || len(args[3]) <= 0 {
		return shim.Error("Arguments error,every arguments can't be empty")
	}

	name := strings.ToLower(args[0])
	kind := strings.ToLower(args[1])
	color := strings.ToLower(args[2])
	owner := strings.ToLower(args[3])
	birdsAsByte, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get bird named " + name)
	} else if len(birdsAsByte) > 0 {
		return shim.Error("This bird has exists : " + name)
	}

	bird := &birds{OBJECT_TYPE_BIRD, name, kind, color, owner}
	birdsAsJson, err := json.Marshal(bird)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(name, birdsAsJson)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{bird.Color, bird.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	err = stub.PutState(colorNameIndexKey, value)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- End init bird")
	return shim.Success(nil)
}

//get a bird
func (b *BirdsChainCode) getBird(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- Start get bird")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the bird to query")
	}
	name := strings.ToLower(args[0])

	valAsBytes, err := stub.GetState(name)
	fmt.Println(valAsBytes)
	if err != nil {
		return shim.Error("Failed to get bird named " + name)
	} else if valAsBytes == nil {
		return shim.Error("Bird does not exists:" + name)
	}
	fmt.Println("- End get bird")
	return shim.Success(valAsBytes)
}

//delete bird
func (b *BirdsChainCode) deleteBird(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("- Start delete bird")
	var birdMarshal birds
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name")
	}
	name := args[0]

	valAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get state for " + name)
	} else if valAsBytes == nil {
		return shim.Error("Bird does not exist:" + name)
	}

	err = json.Unmarshal([]byte(valAsBytes), &birdMarshal)
	if err != nil {
		return shim.Error("Failed to decode JSON of " + name)
	}

	//remove the bird from chaincode state
	err = stub.DelState(name)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// bird the index
	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{birdMarshal.Color, birdMarshal.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(colorNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	fmt.Printf("- End delete birds")
	return shim.Success(nil)
}

//edit bird
func (b *BirdsChainCode) editBird(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- Start delete bird")

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments.")
	}
	name := strings.ToLower(args[0])
	kind := strings.ToLower(args[1])
	owner := strings.ToLower(args[2])
	valAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get state for " + name)
	} else if valAsBytes == nil {
		return shim.Error("Bird does not exist:" + name)
	}
	var birdMarshal birds

	err = json.Unmarshal([]byte(valAsBytes), &birdMarshal)
	if err != nil {
		return shim.Error("Failed to decode JSON of " + name)
	}

	if len(kind) > 0 && kind != birdMarshal.Kind {
		birdMarshal.Kind = kind
	}

	if len(owner) > 0 && owner != birdMarshal.Owner {
		birdMarshal.Owner = owner
	}

	birdsAsJson, err := json.Marshal(birdMarshal)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(name, birdsAsJson)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- End delete bird")
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(BirdsChainCode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err)
	}
}
