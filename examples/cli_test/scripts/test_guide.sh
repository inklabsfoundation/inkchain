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
CC_ID=guide
CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/guide_credit
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
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","1000000000000000000","18","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Issue a new token using ascc has Failed."
    echo_g "===================== A new token has been successfully issued======================= "
    echo
}

registerGuide(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["registerGuide","100099999","evans","true","23"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Issue a new token using ascc has Failed."
    echo_g "=====================A new guide been successfully registered======================= "
    echo
}

registerCompany(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["registerCompany","test","china","999999"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Register a new guide has Failed."
    echo_g "===================== A new company has been successfully registered======================= "
    echo
}

addGuide(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["addGuide","i4230a12f5b0693dd88bb35c79d7e56a68614b199",""]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Register a new guide has Failed."
    echo_g "===================== A new company has been successfully registered======================= "
    echo
}

queryGuideInfo(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryGuideInfo","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query guide info has Failed."
    echo_g "===================== Query guide info has been successfully ======================= "
    echo
}

queryCompanyInfo(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryCompanyInfo","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query guide info has Failed."
    echo_g "===================== Query guide info has been successfully ======================= "
    echo
}

setGuideToBlackList(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["setGuideToBlackList","i4230a12f5b0693dd88bb35c79d7e56a68614b199","nothing"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Register a new guide has Failed."
    echo_g "===================== A new company has been successfully registered======================= "
    echo
}

queryOperateLogs(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryOperateLog","i4230a12f5b0693dd88bb35c79d7e56a68614b199","0"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query operate logs has Failed."
    echo_g "===================== Query operate logs has been successfully ======================= "
    echo
}

queryBlackList(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryBlackList","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query black list has Failed."
    echo_g "===================== Query black list has been successfully ======================= "
    echo
}

removeFromCompany(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["removeFromCompany","i4230a12f5b0693dd88bb35c79d7e56a68614b199","just remove"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Remove a guide from company has Failed."
    echo_g "===================== Remove a guide from company has been successfully registered======================= "
    echo
}

queryGuideWorkList(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryGuideWorkList","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query guide work list has Failed."
    echo_g "===================== Query guide work list has been successfully ======================= "
    echo
}

queryGuideLeaveLogs(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryLeaveLogs","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query guide leave logs has Failed."
    echo_g "===================== Query guide leave logs has been successfully ======================= "
    echo
}

echo_b "=====================5.Issue a token using ascc========================"
issueToken  INK

echo_b "=====================6.Register a guide================================"
registerGuide

echo_b "=====================7.Register a company=============================="
registerCompany

echo_b "=====================8.Add a guide to company=========================="
addGuide

echo_b "=====================9.Query a guide info=============================="
queryGuideInfo

echo_b "=====================10.Query a company info==========================="
queryCompanyInfo

echo_b "=====================11.Add a guide into black list========================"
setGuideToBlackList

echo_b "=====================12.Query a guide info============================="
queryGuideInfo

echo_b "=====================13.Query operate log info========================="
queryOperateLogs

echo_b "=====================14.Query black list info========================="
queryBlackList

echo_b "=====================15.Remove guide from company========================="
removeFromCompany

echo_b "=====================16.Query a guide info============================="
queryGuideInfo

echo_b "=====================17.Query guide work list============================="
queryGuideWorkList

echo_b "=====================18.Query guide leave logs============================="
queryGuideLeaveLogs

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0