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

marbleInvoke(){
    peer chaincode invoke -C mychannel -n marbles --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["initMarble","marble1","blue","35","tom"]}' -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "marble invoke has Failed."
    echo_g "===================== marble invoke successfully======================= "
    echo
}

marbleQuery () {
    echo_b "Attempting to Query marble1 "
    sleep 3
    peer chaincode query -C mychannel -n marbles -c '{"Args":["readMarble","marble1"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account A Failed."
}

echo_b "=====================6.marble invoke======================="
marbleInvoke

echo_b "=====================7.query marble1====================="
marbleQuery

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

