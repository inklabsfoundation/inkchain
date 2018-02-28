#!/bin/bash
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
XC_CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/xcDemo
ORDERER_CA=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
XC_V=1.6

echo_b "Chaincode Path : " $CC_PATH
echo_b "Channel name : " $CHANNEL_NAME

verifyResult () {
    if [ $1 -ne 0 ] ; then
        echo_b "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
        echo_r "================== ERROR !!! FAILED to execute MVE =================="
        echo
        exit 1
    fi
}

installChaincode () {
    peer chaincode install -n xc -v $XC_V -p ${XC_CC_PATH} -o orderer.example.com:7050 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode token installation on remote peer0 has Failed"
    echo_g "===================== Chaincode is installed success on remote peer0===================== "
    echo
}

instantiateChaincode () {
    local starttime=$(date +%s)
    peer chaincode upgrade -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n xc -v $XC_V -c '{"Args":["init","INK","3c97f146e8de9807ef723538521fcecd5f64c79a"]}' -P "OR ('Org1MSP.member')" >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode instantiation on pee0.org1 on channel '$CHANNEL_NAME' failed"
    echo_g "=========== Chaincode token Instantiation on peer0.org1 on channel '$CHANNEL_NAME' is successful ========== "
    echo_b "Instantiate spent $(($(date +%s)-starttime)) secs"
    echo
}

echo_b "=====================1.Install chaincode token on Peer0/Org0========================"
installChaincode

echo_b "=====================2.Instantiate chaincode token, this will take a while, pls waiting...==="
instantiateChaincode

echo
echo_g "=====================All GOOD, MVE initialization completed ===================== "
echo
exit 0
