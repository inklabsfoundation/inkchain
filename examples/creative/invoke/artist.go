package invoke

import (
	. "github.com/inklabsfoundation/inkchain/examples/creative/conf"
	. "github.com/inklabsfoundation/inkchain/examples/creative/util"
	. "github.com/inklabsfoundation/inkchain/examples/creative/model"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"encoding/json"
	"fmt"
	"strings"
)

type ArtistInvoke struct{}

func (*ArtistInvoke) AddArtist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("add artist start.")
	username := args[0]
	artist_name := args[1]
	artist_desc := args[2]
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
	if artistAsBytes != nil {
		fmt.Println("This artist already exist: " + artist_name)
		return shim.Error("This artist already exist: " + artist_name)
	}
	// add artist
	artist := Artist{artist_name, artist_desc, username}
	artistAsBtyes, err := json.Marshal(artist)
	err = stub.PutState(artist_key, artistAsBtyes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("add artist success."))
}

func (c *ArtistInvoke) ModifyArtist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("modify artist start.")
	username := args[0]
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
	artist_key := GetArtistKey(username)
	// verify weather the artist exists
	artistAsBytes, err := stub.GetState(artist_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}
	if artistAsBytes == nil {
		fmt.Println("This artist doesn't exist: " + artist_key)
		return shim.Error("This artist doesn't exist: " + artist_key)
	}
	var artistJSON Artist
	err = json.Unmarshal([]byte(artistAsBytes), &artistJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	if artistJSON.Username != userJSON.Username {
		return shim.Error("The artist's username doesn't correspond with the user's.")
	}
	err = GetModifyArtist(&artistJSON, args[1:])
	if err != nil {
		return shim.Error(err.Error())
	}
	artistJSONasBytes, err := json.Marshal(artistJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(artist_key, artistJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("modify artist success."))
}

func (c *ArtistInvoke) DeleteArtist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("delete artist start.")
	username := args[0]
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
	artist_key := GetArtistKey(username)
	// verify weather the artist exists
	artistAsBytes, err := stub.GetState(artist_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}
	if artistAsBytes == nil {
		fmt.Println("This artist doesn't exist: " + artist_key)
		return shim.Error("This artist doesn't exist: " + artist_key)
	}
	var artistJSON Artist
	err = json.Unmarshal([]byte(artistAsBytes), &artistJSON)
	if artistJSON.Username != userJSON.Username {
		return shim.Error("The artist's username doesn't correspond with the user's.")
	}
	// delete artist's info
	err = stub.DelState(artist_key)
	if err != nil {
		fmt.Println("Fail to delete: " + artist_key)
		return shim.Error("Fail to delete" + artist_key)
	}
	return shim.Success([]byte("delete artist success."))
}

func (c *ArtistInvoke) QueryArtist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("query artist start.")
	username := args[0]
	artist_key := GetArtistKey(username)
	artistAsBytes, err := stub.GetState(artist_key)
	if err != nil {
		return shim.Error("Fail to get artist: " + err.Error())
	}
	if artistAsBytes == nil {
		fmt.Println("This artist doesn't exist: " + username)
		return shim.Error("This artist doesn't exist: " + username)
	}
	return shim.Success(artistAsBytes)
}

func (c *ArtistInvoke) ListOfArtist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("list of artist start.")
	resultsIterator, err := stub.GetStateByRange(ArtistPrefix+StateStartSymbol, ArtistPrefix+StateEndSymbol)
	if err != nil {
		return shim.Error(err.Error())
	}
	list, err := GetListResult(resultsIterator)
	if err != nil {
		return shim.Error("getListResult failed")
	}
	return shim.Success(list)
}
