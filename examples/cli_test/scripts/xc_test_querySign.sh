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

queryTxInfo () {
    echo_b "Attempting to Query Tx Info "
    sleep 3
    peer chaincode query  -C ${CHANNEL_NAME} -n xscc -c '{"Args":["querySign","eafccedd1d57e899d6e17c9fc9e6486ef8ee7bfa61ab55fbe50221ca7ec8c01b"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query Tx Info Failed."
}


echo_b "=====================1.QueryTxInfo====================="
queryTxInfo

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

