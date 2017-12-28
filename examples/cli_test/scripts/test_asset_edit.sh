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

assetQuery_Asset () {
    echo_b "Attempting to Query asset "
    sleep 3
    peer chaincode query -C mychannel -n asset -c '{"Args":["readAsset","Blockchain Guide"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query asset Failed."
}

assetInvoke_Edit_TYPE(){

    sleep 3
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["editAsset","Blockchain Guide","Type","E-Book"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: editAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetInvoke_Edit_CONTENT(){

    sleep 3
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["editAsset","Blockchain Guide","Content","a e-book version."]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: editAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetInvoke_Edit_PRICE_TYPE(){

    sleep 3
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["editAsset","Blockchain Guide","PriceType","INK"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: editAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

assetInvoke_Edit_PRICE(){

    sleep 3
    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["editAsset","Blockchain Guide","Price","8"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "asset invoke: editAsset has Failed."
    echo_g "===================== asset invoke successfully======================= "
    echo
}

chaincodeQueryB () {

    sleep 3
    echo_b "Attempting to  query account B's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","a5ff00eb44bf19d5dfbde501c90e286badb58df4","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account B Failed."

}


echo_b "=====================Test Asset's edit invoke======================="

echo_b "=====================1.query asset====================="
assetQuery_Asset

echo_b "=====================2.edit asset====================="
assetInvoke_Edit_TYPE
assetInvoke_Edit_CONTENT
assetInvoke_Edit_PRICE_TYPE
assetInvoke_Edit_PRICE

echo_b "=====================3.query asset again====================="
assetQuery_Asset

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

