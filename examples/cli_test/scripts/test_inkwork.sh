#!/usr/bin/env bash

#
#Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

#system account
#i07caf88941eafcaaa3370657fccc261acb75dfba
#70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4

#user account
#ia5ff00eb44bf19d5dfbde501c90e286badb58df4
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

issueToken(){
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","1000000000000000","18","i07caf88941eafcaaa3370657fccc261acb75dfba"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Issue a new token using ascc has Failed."
    echo_g "===================== A new token has been successfully issued======================= "
    echo
}

makeTransfer(){
    echo_b "pls wait 5 secs..."
    sleep 5
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n token -c '{"Args":["transfer","ia5ff00eb44bf19d5dfbde501c90e286badb58df4","INK","50000000000000"]}' -i "1000000000" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Make transfer has Failed."
    echo_g "===================== Make transfer success ======================= "
    echo
}

registerWork(){
    sleep 3
    peer chaincode invoke -C mychannel -n inkwork --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["registerWork","2018","2","4"]}' -i "100000000" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "registerWork invoke has Failed."
    echo_g "===================== registerWork invoke successfully======================= "
    echo
}

sell(){
    sleep 3
    peer chaincode invoke -C mychannel -n inkwork --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["sell","2018","10","1"]}' -i "100000000" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "inkwork invoke has Failed."
    echo_g "===================== inkwork sell successfully======================= "
    echo
}

purchase(){
    sleep 3
    peer chaincode invoke -C mychannel -n inkwork --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["purchase","2018"]}' -i "100000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "inkwork invoke has Failed."
    echo_g "===================== inkwor purchase successfully======================= "
    echo
}

inkworkQuery () {
    sleep 3
    peer chaincode query -C mychannel -n inkwork -c '{"Args":["queryInkwork","2018"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query inkwork Failed."
    echo_g "===================== inkwork query successfully======================= "
    echo
}

listQuery () {
    sleep 3
    peer chaincode query -C mychannel -n inkwork -c '{"Args":["query","1"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query inkwork list Failed."
    echo_g "===================== list query successfully======================= "
    echo
}


echo_b "=====================register token======================="
issueToken INK
makeTransfer

echo_b "=====================1.registerWork======================="
registerWork

echo_b "=====================2.sell work======================="
sell

echo_b "=====================3.purchase purchase======================="
purchase

echo_b "=====================4.inkworkQuery======================="
inkworkQuery

echo_b "=====================5.listQuery======================="
listQuery

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

