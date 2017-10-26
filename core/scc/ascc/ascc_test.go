package ascc

import (
	"testing"
	"github.com/inkchain/inkchain/core/chaincode/shim"
	"fmt"
	"strconv"
)

const(
	MPriKey = "bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe"
	MAddress= "4230a12f5b0693dd88bb35c79d7e56a68614b199"

	testPriKey = "60ef69c3e7d5a8a8d6c25406ab321f5e51c475dd44ddef3a9a47e91d764edae1"
	testAddress= "3c97f146e8de9807ef723538521fcecd5f64c79a"
	//MPriKey = "bab0c1204b2e7f344f9d1fbe8ad978d5355e32b8fa45b10b600d64ca970e0dc9"
	//MAddress= "411b6f8f24F28CaAFE514c16E11800167f8EBd89"
	//testPriKey = "703339e58975be8f91981b9fc6f001ad8cedf936067284b80e6e9c57b32a2ddd"
	//testAddress= "4708A97Bf6F53c2Ca664BB003eF80fa6B997D656"
)

// check Init function
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

// check Invoke function
func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

// check GetBalance
func checkGetBalance(t *testing.T, stub *shim.MockStub, args [][]byte, cmp int) {
	res1 := stub.MockInvoke("1", args)
	if res1.Status != shim.OK {
		fmt.Println("GetBalance failed", string(res1.Message))
		t.FailNow()
	}
	amount,_ := strconv.Atoi(string(res1.Payload))
	if amount != cmp {
		fmt.Printf("Query result error! %v", amount )
		t.FailNow()
	}
}

// check queryInfo
func checkQueryInfo(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res1 := stub.MockInvoke("1", args)
	if res1.Status != shim.OK {
		fmt.Println("GetBalance failed", string(res1.Message))
		t.FailNow()
	}
	fmt.Printf("QueryInfo Result: %s \n", string(res1.Payload))
}

// Test 0: Init ascc
func TestAssetSysCC_Init(t *testing.T) {
	ascc := new(AssetSysCC)
	stub := shim.NewMockStub("ascc", ascc)
	checkInit(t, stub, [][]byte{[]byte("")})
}

// Test 1: Test Register Token
func TestAssetSysCC_RegisterToken(t *testing.T) {

	fmt.Println("-----------------------------------")
	fmt.Println("Test1: registerToken")

	ascc := new(AssetSysCC)
	stub := shim.NewMockStub("ascc", ascc)
	checkInit(t, stub, [][]byte{[]byte("")})

	res_test3 := stub.MockInvoke("1", [][]byte{[]byte("registerToken"), []byte("SSToken"), []byte("250"), []byte("18"), []byte(MAddress[:])})

	if res_test3.Status != shim.OK {
		fmt.Println("Register token failed", string(res_test3.Message))
		t.FailNow()
	}

	fmt.Println("Test registerToken Success!")

}

// Test 2: Test Issue Token
func TestAssetSysCC_IssueToken(t *testing.T) {

	fmt.Println("-----------------------------------")
	fmt.Println("Test2: issueToken")

	ascc := new(AssetSysCC)
	stub := shim.NewMockStub("ascc", ascc)
	checkInit(t, stub, [][]byte{[]byte("")})


	res_test2 := stub.MockInvoke("1", [][]byte{[]byte("registerToken"), []byte("SSToken"), []byte("250"), []byte("18"), []byte(MAddress[:])})

	if res_test2.Status != shim.OK {
		fmt.Println("Register token failed", string(res_test2.Message))
		t.FailNow()
	}

	res_test3 := stub.MockInvoke("1", [][]byte{[]byte("issueToken"), []byte("SSToken"), []byte("250"), []byte("18"), []byte(MAddress[:])})

	if res_test3.Status != shim.OK {
		fmt.Println("Register token failed", string(res_test3.Message))
		t.FailNow()
	}
	checkQueryInfo(t, stub, [][]byte{[]byte("getTokenInfo"), []byte("SSToken")})

	////query token quantity
	//	res1 := stub.MockInvoke("2", [][]byte{[]byte("getBalance"), []byte(MAddress[:]), []byte("SSToken")});
	//	if res1.Status != shim.OK {
	//		fmt.Println("Query failed", string(res1.Message))
	//		t.FailNow()
	//	}
	//	amount,_ := strconv.Atoi(string(res1.Payload))
	//	if amount != 250 {
	//		fmt.Printf("Query result error! %v", amount )
	//		t.FailNow()
	//	}

	fmt.Println("Test issueToken for a registered one Success!")

	res_test4 := stub.MockInvoke("2", [][]byte{[]byte("issueToken"), []byte("MToken"), []byte("888"), []byte("20"), []byte(testAddress[:])})
	if res_test4.Status != shim.OK {
		fmt.Println("Register token failed", string(res_test3.Message))
		t.FailNow()
	}
	checkQueryInfo(t, stub, [][]byte{[]byte("getTokenInfo"), []byte("MToken")})

	////query token quantity
	//res2 := stub.MockInvoke("2", [][]byte{[]byte("getBalance"), []byte(testAddress[:]), []byte("CMBToken")});
	//if res1.Status != shim.OK {
	//	fmt.Println("Query failed", string(res2.Message))
	//	t.FailNow()
	//}
	//amount2,_ := strconv.Atoi(string(res2.Payload))
	//if amount2 != 888 {
	//	fmt.Printf("Query result error! %v", amount2 )
	//	t.FailNow()
	//}

	fmt.Println("Test issueToken for an un registered one Success!")
}

// Test 3: Test Invalidate Token
func TestAssetSysCC_InvalidateToken(t *testing.T) {

	fmt.Println("-----------------------------------")
	fmt.Println("Test3: invalidateToken")

	//fmt.Println("******test string to big.newInt")
	//str := "12321"
	//strInt := big.NewInt(0)
	//strInt.SetString(str,10)
	//fmt.Println(strInt.String())
	//fmt.Println("*******************************")

	ascc := new(AssetSysCC)
	stub := shim.NewMockStub("ascc", ascc)
	checkInit(t, stub, [][]byte{[]byte("")})

	res_test3 := stub.MockInvoke("1", [][]byte{[]byte("issueToken"), []byte("SSToken"), []byte("250"), []byte("18"), []byte(MAddress[:])})

	if res_test3.Status != shim.OK {
		fmt.Println("Register token failed", string(res_test3.Message))
		t.FailNow()
	}

	////query token quantity
	//res1 := stub.MockInvoke("2", [][]byte{[]byte("getBalance"), []byte(MAddress[:]), []byte("SSToken")});
	//if res1.Status != shim.OK {
	//	fmt.Println("Query failed", string(res1.Message))
	//	t.FailNow()
	//}
	//amount,_ := strconv.Atoi(string(res1.Payload))
	//if amount != 250 {
	//	fmt.Printf("Query result error! %v", amount )
	//	t.FailNow()
	//}

	//beging to invalidate this token
	checkQueryInfo(t, stub, [][]byte{[]byte("getTokenInfo"), []byte("SSToken")})

	testInvalidate := stub.MockInvoke("4", [][]byte{[]byte("invalidateToken"), []byte("SSToken")});
	if testInvalidate.Status != shim.OK {
		fmt.Println("Query failed", string(testInvalidate.Message))
		t.FailNow()
	}

	checkQueryInfo(t, stub, [][]byte{[]byte("getTokenInfo"), []byte("SSToken")})
}

//// Test 4: Transfer Token & queryTokenInfo
//func TestAssetSysCC_TransferToken(t *testing.T) {
//
//	fmt.Println("-----------------------------------")
//	fmt.Println("Test4: transferToken")
//
//	ascc := new(AssetSysCC)
//	stub := shim.NewMockStub("ascc", ascc)
//	checkInit(t, stub, [][]byte{[]byte("")})
//
//	//register token
//	res_test2 := stub.MockInvoke("1", [][]byte{[]byte("registerToken"), []byte("SSToken"), []byte("1000000000"), []byte(MAddress[:])})
//
//	if res_test2.Status != shim.OK {
//		fmt.Println("Register token failed", string(res_test2.Message))
//		t.FailNow()
//	}
//	checkQueryInfo(t, stub, [][]byte{[]byte("getTokenInfo"), []byte("SSToken")})
//
//	//issue token
//	res_test3 := stub.MockInvoke("2", [][]byte{[]byte("issueToken"), []byte("SSToken"), []byte("1000000000"), []byte(MAddress[:])})
//
//	if res_test3.Status != shim.OK {
//		fmt.Println("Issue token failed", string(res_test3.Message))
//		t.FailNow()
//	}
//	checkQueryInfo(t, stub, [][]byte{[]byte("getTokenInfo"), []byte("SSToken")})
//
//	//query token quantity
//	checkGetBalance(t, stub, [][]byte{[]byte("getBalance"), []byte(MAddress[:]), []byte("SSToken")}, 1000000000)
//
//	//make transfer
//	res_test4 := stub.MockInvoke("3", [][]byte{[]byte("transfer"), []byte(MAddress[:]), []byte(testAddress[:]), []byte("SSToken"), []byte("100000000"), })
//
//	if res_test4.Status != shim.OK {
//		fmt.Println("Transfer token failed", string(res_test3.Message))
//		t.FailNow()
//	}
//
//	//Check the transfer result
//	checkGetBalance(t, stub, [][]byte{[]byte("getBalance"), []byte(MAddress[:]), []byte("SSToken")}, 900000000)
//	checkGetBalance(t, stub, [][]byte{[]byte("getBalance"), []byte(testAddress[:]), []byte("SSToken")}, 100000000)
//
//	//output info of SSToken
//	checkQueryInfo(t, stub, [][]byte{[]byte("getTokenInfo"), []byte("SSToken")})
//}

