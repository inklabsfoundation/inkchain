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
CC_ID=student
CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/student
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
    echo_g "===================== A new token has success======================= "
    echo
}

registerSchool(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["registerSchool","100099999","Tsinghua University","北京","1","long times ago"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Register a school has Failed."
    echo_g "=====================A new school been success======================="
    echo
}

querySchoolInfo(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["querySchoolInfo","100099999"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query school info has failed."
    echo_g "=====================Query school info success======================="
    echo
}

registerStudent(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["registerStudent","330381199342105115115","TestStudent","26","1"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Register a student has Failed."
    echo_g "=====================A new student success======================="
    echo
}

queryStudentInfo(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryStudentInfo","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query student info has failed."
    echo_g "=====================Query student info success======================="
    echo
}

enrolment(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["enrolment","100099999","i4230a12f5b0693dd88bb35c79d7e56a68614b199","2010届","210711A","210711A054","1"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Student enrolment has Failed."
    echo_g "=====================Student enrolment success======================="
    echo
}

graduate(){
    sleep 3
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["graduate","100099999","i4230a12f5b0693dd88bb35c79d7e56a68614b199","表现良好，允许毕业","9999999999","40"]}' -i "1000000000" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Student graduate has Failed."
    echo_g "=====================Student graduate success======================="
    echo
}

queryStudentStudyLog(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryStudentStudyLog","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query student study log has failed."
    echo_g "=====================Query student study log success======================="
    echo
}

queryStudentGraduationLog(){
    sleep 3
    peer chaincode query -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ${CC_ID} -c '{"Args":["queryStudentGraduationLog","i4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Query student study log has failed."
    echo_g "=====================Query student graduation log success======================="
    echo
}

echo_b "=====================5.Issue a token using ascc========================"
issueToken  INK

echo_b "=====================6.Register a school================================"
registerSchool

echo_b "=====================7.Query school info================================"
querySchoolInfo

echo_b "=====================8.Register a student================================"
registerStudent

echo_b "=====================7.Query student info================================"
queryStudentInfo

echo_b "=====================8.Student enrolment================================"
enrolment

echo_b "=====================9.Query student info================================"
queryStudentInfo

echo_b "=====================10.Student graduate================================"
graduate

echo_b "=====================11.Query student info================================"
queryStudentInfo

echo_b "=====================12.Query student study log================================"
queryStudentStudyLog

echo_b "=====================13.Query student graduate log================================"
queryStudentGraduationLog

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0