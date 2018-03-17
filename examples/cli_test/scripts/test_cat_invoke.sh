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

catSystemInit(){
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["initSystemCat", "'$1'"]}' -i "1000000000" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat invoke successfully======================= "
    echo
}

catDel(){
    echo_b "Attempting to delete cat"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["delete","7918"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat delete successfully======================= "
    echo
}

saleState(){
    echo_b "Attempting to set cat saleState"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["setState","7918","0","1"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat setSaleStae successfully======================= "
    echo
}

mateState(){
    echo_b "Attempting to set cat mateState"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["setState","7918","2","1"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat setMateStae successfully======================= "
    echo
}

salePrice(){
    echo_b "Attempting to set cat salePrice"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["setState","7918","1","4"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat setSalePrice successfully======================= "
    echo
}

matePrice(){
    echo_b "Attempting to set cat matePrice"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["setState","7918","3","5"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat setMalePrice successfully======================= "
    echo
}

catBreed(){
    echo_b "Attempting to cat breed"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["breed","5060","7918","2108"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "cat invoke has Failed."
    echo_g "===================== cat breed successfully======================= "
    echo
}

catBuyA(){
    echo_b "Attempting to buy cat A"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["buy","7918"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >log.txt
    res=$?
    cat log.txt
    verifyResult $res "buyCat has Failed."
    echo_g "===================== buycat success ======================= "
    echo
}

catBuyB(){
    echo_b "Attempting to buy cat B"
    sleep 3
    peer chaincode invoke -C mychannel -n cat --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["buy","5060"]}' -i "1000000000" -z 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5 >log.txt
    res=$?
    cat log.txt
    verifyResult $res "buyCat has Failed."
    echo_g "===================== buycat success ======================= "
    echo
}

echo_b "=====================register token======================="
issueToken INK
makeTransfer

echo_b "=====================1.cat init======================="
catSystemInit 7918
catSystemInit 5060

echo_b "=====================2.cat delete======================="
#catDel

echo_b "=====================3.cat buy======================="
catBuyA

echo_b "=====================4.cat buy======================="
catBuyB

echo_b "=====================5.cat sale state====================="
saleState

echo_b "=====================6.cat mate state====================="
mateState

echo_b "=====================7.cat sale price====================="
salePrice

echo_b "=====================8.cat mate price====================="
matePrice

echo_b "=====================9.cat breed====================="
catBreed

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

