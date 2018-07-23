#!/usr/bin/env bash
#
#Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

# Detecting whether can import the header file to render colorful cli output
if [ -f ./func_init.sh ]; then
 source ./func_init.sh
elif [ -f scripts/func_init.sh ]; then
 source scripts/func_init.sh
else
 alias echo_r="echo"
 alias echo_g="echo"
 alias echo_b="echo"
fi

# issueToken INK 1000 18 ${USER_ADDRESS_01}
issueToken(){
    sleep 5
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","'$2'","'$3'","'$4'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Issue a new token using ascc has Failed."
    echo_g "===================== A new token has been successfully issued======================= "
    echo
}

# makeTransfer ${USER_ADDRESS_02} INK  500 ${USER_TOKEN_01}
makeTransfer(){
    echo_b "pls wait 5 secs..."
    sleep 5
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n token -c '{"Args":["transfer","'$1'","'$2'","'$3'"]}' -i "10" -z "$4" >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Make transfer has Failed."
    echo_g "===================== Make transfer success ======================= "
    echo
}

# chaincodeQuery ${USER_ADDRESS_01} INK
# chaincodeQuery ${USER_ADDRESS_02} INK
chaincodeQuery () {
    echo_b "Attempting to Query account A's balance on peer "
    sleep 3
    peer chaincode query -C ${CHANNEL_NAME} -n token -c '{"Args":["getBalance","'$1'","'$2'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account A Failed."
}
