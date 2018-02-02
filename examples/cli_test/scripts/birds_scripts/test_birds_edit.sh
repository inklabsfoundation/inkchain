#!/usr/bin/env bash

#
#Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

# Detecting whether can import the header file to render colorful cli output
CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="mychannel"}
: ${TIMEOUT:="60"}
COUNTER=0
MAX_RETRY=5

ORDERER_CA=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "Chaincode Path : "$CC_PATH
echo "Channel name : "$CHANNEL_NAME

verifyResult () {
    if [ $1 -ne 0 ] ; then
        echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
        echo "================== ERROR !!! FAILED to execute MVE =================="
        echo
        exit 1
    fi
}

birdDel () {
    echo "Attempting to Delete bird "
    peer chaincode invoke --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C mychannel -n birds -c '{"Args":["editBird","birdqt","parrot","blue"]}' -i "10" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account A Failed."
}

echo "=====================7.del bird====================="
birdDel

echo
echo "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0
