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

invokeInke () {
    echo_b "Attempting to Query account "
    sleep 3
    # args public platform,  public account ,amount ,ink address,txId
    peer chaincode invoke -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n xscc -c '{"Args":["unlock","QTUM","41f14eebf44f921beb0674d53b080afa33db2ff1","1000","id6b39eb631df8ee60e46a576231ccf1fcd204a5e","6df6dc30539d968a71e53213855e680fc50bd43885e41f5ee4ddd1ae19e68cb9"]}' -i "100000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
#    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n xc -c '{"Args":["lock","QTUM","HelloWorld","10","3c97f146e8de9807ef723538521fcecd5f64c79a"]}' -i "5" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query unlock A Failed."
}

chaincodeQueryB () {
    echo_b "Attempting to  query account B's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","id6b39eb631df8ee60e46a576231ccf1fcd204a5e","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account B Failed."

}

chaincodeQueryA () {
    echo_b "Attempting to Query account A's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","i4230a12f5b0693dd88bb35c79d7e56a68614b199","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account A Failed."
}


echo_b "=====================1.invoke account====================="
invokeInke
echo_b "=====================2.query account A===================="
chaincodeQueryA
echo_b "=====================3.query account B===================="
chaincodeQueryB
echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

