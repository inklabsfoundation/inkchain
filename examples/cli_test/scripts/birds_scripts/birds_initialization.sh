#!/bin/bash
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
BIRDS_CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/birds

echo "Chaincode Path : " $CC_PATH
echo "Channel name : " $CHANNEL_NAME

verifyResult () {
    if [ $1 -ne 0 ] ; then
        echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
        echo "================== ERROR !!! FAILED to execute MVE =================="
        echo
        exit 1
    fi
}

## Sometimes Join takes time hence RETRY atleast for 5 times
installChaincode () {
    peer chaincode install -n birds -v 1.0 -p ${BIRDS_CC_PATH} -o orderer.example.com:7050 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode token installation on remote peer0 has Failed"
    echo "===================== Chaincode is installed success on remote peer0===================== "
    echo
}

instantiateChaincode () {
    local starttime=$(date +%s)
    peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n birds -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member')" >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode instantiation on pee0.org1 on channel '$CHANNEL_NAME' failed"
    echo "=========== Chaincode token Instantiation on peer0.org1 on channel '$CHANNEL_NAME' is successful ========== "
    echo "Instantiate spent $(($(date +%s)-starttime)) secs"
    echo
}

echo "=====================1.Install chaincode token on Peer0/Org0========================"
installChaincode

echo "=====================2.Instantiate chaincode token, this will take a while, pls waiting...==="
instantiateChaincode

echo
echo "=====================All GOOD, MVE initialization completed ===================== "
echo
exit 0
