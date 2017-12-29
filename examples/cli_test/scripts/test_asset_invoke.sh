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
    verifyResult $res "query asset: Failed."
}

chaincodeQueryBuyer () {
    echo_b "Attempting to  query account B's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","07caf88941eafcaaa3370657fccc261acb75dfba","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account B Failed."

}

assetInvoke_Buy () {
    echo_b "Attempting to buy asset "
    sleep 3

    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["buyAsset","Blockchain Guide","Rose"]}' -i "10" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt

    res=$?
    cat log.txt
    verifyResult $res "query asset: To The Moon Failed."
}

assetInvoke_Delete () {
    echo_b "Attempting to delete asset "
    sleep 3

    peer chaincode invoke -C mychannel -n asset --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["deleteAsset","Blockchain Guide"]}' -i "10" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt

    res=$?
    cat log.txt
    verifyResult $res "query asset: To The Moon Failed."
}

assetQueryRange_Asset () {
    echo_b "Attempting to Query asset "
    sleep 3
    peer chaincode query -C mychannel -n asset -c '{"Args":["readAssetByRange","",""]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "range query asset failed."
}

echo_b "=====================Test Asset's buyAsset invoke====================="
echo_b "=====================1.query asset====================="
assetQuery_Asset

echo_b "=====================2.0 query balance before transfer asset====================="
chaincodeQueryBuyer

echo_b "=====================2.1 transfer asset====================="
assetInvoke_Buy

echo_b "=====================2.2 query balance after transfer asset====================="
chaincodeQueryBuyer

echo_b "=====================3. query asset====================="
assetQuery_Asset

echo_b "=====================Test Asset's delete invoke====================="
echo_b "=====================4.delete asset====================="
assetInvoke_Delete


echo_b "=====================5. query all the assets====================="
assetQueryRange_Asset

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

