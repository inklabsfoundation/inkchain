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

assetInvoke_AddUser(){
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["addUser","Daniel","18"]}' -i "10" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: addUser has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetQuery_User () {
    echo_b "Attempting to Query user "
    sleep 3
    peer chaincode query -C mychannel -n asset -c '{"Args":["queryUser","Daniel"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query user: Dainel Failed."
}

assetInvoke_AddAsset(){
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["addAsset","To The Moon","GAME","a game developed by Freebird Games","Daniel"]}' -i "10" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: addAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetQuery_Asset () {
    echo_b "Attempting to Query asset "
    sleep 3
    peer chaincode query -C mychannel -n asset -c '{"Args":["readAsset","To The Moon"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query asset: To The Moon Failed."
}

assetInvoke_Delete () {
    echo_b "Attempting to delete asset "
    sleep 3

    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["deleteAsset","To The Moon"]}' -i "10" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >&log.txt

    res=$?
    cat log.txt
    verifyResult $res "query asset: To The Moon Failed."
}

echo_b "=====================6.add user======================="
assetInvoke_AddUser

echo_b "=====================7.query user======================="
assetQuery_User

echo_b "=====================8.add product======================="
assetInvoke_AddAsset

echo_b "=====================7.query product====================="
assetQuery_Asset

echo_b "=====================8.delete product====================="
assetInvoke_Delete


echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

