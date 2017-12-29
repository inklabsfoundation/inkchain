#!/usr/bin/env bash

#
#Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

# Detecting whether can import the header file to render colorful cli output
if [ -f ./header.sh ]; then
 source ./header.sh
elif [ -f scripts/header.sh ]; then
 source scripts/header.sh
else
 alias echo_r="echo"
 alias echo_g="echo"
 alias echo_b="echo"
fi

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="mychannel"}
: ${TIMEOUT:="60"}
COUNTER=0
MAX_RETRY=5
CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/token
ORDERER_CA=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo_b "Chaincode Path : "$CC_PATH
echo_b "Channel name : "$CHANNEL_NAME

verifyResult () {
    if [ $1 -ne 0 ] ; then
        echo_b "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
        echo_r "================== ERROR !!! FAILED to execute MVE =================="
        echo
        exit 1
    fi
}

issueToken(){
#    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","100","18","4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","100","18","07caf88941eafcaaa3370657fccc261acb75dfba"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Issue a new token using ascc has Failed."
    echo_g "===================== A new token has been successfully issued======================= "
    echo
}

makeTransfer(){
    echo_b "pls wait 5 secs..."
    sleep 5
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n token -c '{"Args":["transfer","a5ff00eb44bf19d5dfbde501c90e286badb58df4","INK","30"]}' -i "1" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >log.txt
#    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n token -c '{"Args":["transfer","3c97f146e8de9807ef723538521fcecd5f64c79a","INK","10"]}' -i "1" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt

    res=$?
    cat log.txt
    verifyResult $res "Make transfer has Failed."
    echo_g "===================== Make transfer success ======================= "
    echo
}

chaincodeQueryA () {
    echo_b "Attempting to Query account A's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","07caf88941eafcaaa3370657fccc261acb75dfba","INK"]}' >log.txt
#    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","4230a12f5b0693dd88bb35c79d7e56a68614b199","INK"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query account A Failed."
}

chaincodeQueryB () {
    echo_b "Attempting to  query account B's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","a5ff00eb44bf19d5dfbde501c90e286badb58df4","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account B Failed."
   
}

# for init asset
assetInvoke_AddUser01(){
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["addUser","Daniel","18"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: addUser has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetInvoke_AddUser02(){
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["addUser","Rose","16"]}' -i "10" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: addUser has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetQuery_User01 () {
    echo_b "Attempting to Query user "
    sleep 3
    peer chaincode query -C mychannel -n asset -c '{"Args":["queryUser","Daniel"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query user: Dainel Failed."
}

assetQuery_User02 () {
    echo_b "Attempting to Query user "
    sleep 3
    peer chaincode query -C mychannel -n asset -c '{"Args":["queryUser","Rose"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query user: Dainel Failed."
}

assetInvoke_AddAsset(){
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["addAsset","Blockchain Guide","BOOK","a book on blockchain.","CCToken","10","Daniel"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: addAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetInvoke_AddAsset2(){
    sleep 3

    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["addAsset","Blockchain Guide2","BOOK","a book on blockchain 2.","INK","100","Daniel"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: addAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetInvoke_AddAsset3(){
    sleep 3

    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["addAsset","Z Blockchain Guide3","BOOK","a book on blockchain 3.","INK","1000","Daniel"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: addAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetQuery_Asset () {
    echo_b "Attempting to Query asset "
    sleep 3
    peer chaincode query -C mychannel -n asset -c '{"Args":["readAsset","Blockchain Guide"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query asset: To The Moon Failed."
}

echo_b "=====================6.Issue a token using ascc========================"
issueToken INK

echo_b "=====================7.Transfer 10 amount of INK====================="
makeTransfer

echo_b "=====================8.Query transfer result of From account====================="
chaincodeQueryA

echo_b "=====================9.Query transfer result of To account====================="
chaincodeQueryB

echo_b "=====================0. Init for test Asset====================="
echo_b "=====================0.1 add 2 user======================="
assetInvoke_AddUser01
assetInvoke_AddUser02

echo_b "=====================0.2 query 2 user======================="
assetQuery_User01
assetQuery_User02


echo_b "=====================0.3 add asset======================="
assetInvoke_AddAsset
assetInvoke_AddAsset2
assetInvoke_AddAsset3

echo_b "=====================0.4 query asset====================="
assetQuery_Asset


echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

