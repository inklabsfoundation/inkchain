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

catQuery () {
    echo_b "Attempting to Query cat"
    sleep 3
    peer chaincode query -C mychannel -n cat -c '{"Args":["query","7918"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query cat Failed."
    echo_g "===================== cat query successfully======================= "
    echo
}

catQueryAll () {
    echo_b "Attempting to Query all cat "
    sleep 3
    peer chaincode query -C mychannel -n cat -c '{"Args":["queryAll"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query cat Failed."
    echo_g "===================== cat query all successfully======================= "
    echo
}

catQuerySale () {
    echo_b "Attempting to Query sale cat "
    sleep 3
    peer chaincode query -C mychannel -n cat -c '{"Args":["querySale"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query cat Failed."
    echo_g "===================== cat query sale successfully======================= "
    echo
}

echo_b "=====================1.cat query all====================="
catQueryAll

echo_b "=====================2.cat query sale====================="
catQuerySale

echo_b "=====================3.cat query====================="
catQuery

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

