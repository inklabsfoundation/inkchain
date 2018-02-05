#!/usr/bin/env bash

#
#Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

#system account
#07caf88941eafcaaa3370657fccc261acb75dfba
#70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4

#user account
#a5ff00eb44bf19d5dfbde501c90e286badb58df4
#344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5

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

catSystemInit(){
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["initSystemCat"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat invoke successfully======================= "
    echo
}

createAuction(){
    echo_b "Attempting to set cat createuction"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["createAuction","8415","10","300"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat createAuction successfully======================= "
    echo
}

bid(){
    echo_b "Attempting to bid cat "
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["bid","8415","310"]}' -i "10" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >log.txt
    res=$?
    cat log.txt
    verifyResult $res "bid Cat has Failed."
    echo_g "===================== bid cat success ======================= "
    echo
}

endAuction(){
    echo_b "Attempting to endAuction"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["endAuction","8415"]}' -i "10" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat endAuction successfully======================= "
    echo
}

payAuction(){
    echo_b "Attempting to payAuction"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["payAuction","8415","1","310",""]}' -i "10" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat payAuction successfully======================= "
    echo
}

echo_b "=====================1.cat system init======================="
catSystemInit

echo_b "=====================2.cat auction======================="
createAuction

echo_b "=====================3.cat bid======================="
bid

echo_b "=====================4.cat end auction====================="
endAuction

echo_b "=====================5.cat pay auction====================="
#payAuction

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

