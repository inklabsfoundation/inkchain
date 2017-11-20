#!/usr/bin/env bash

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
CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/token

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
    peer chaincode invoke -o orderer.example.com:7050 -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","100","18","4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Issue a new token using ascc has Failed."
    echo_g "===================== A new token has been successfully issued======================= "
    echo
}

makeTransfer(){
    echo_b "pls wait 5 secs..."
    sleep 5
    peer chaincode invoke -o orderer.example.com:7050 -C ${CHANNEL_NAME} -n token -c '{"Args":["transfer","3c97f146e8de9807ef723538521fcecd5f64c79a","INK","10"]}' -i "1" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Make transfer has Failed."
    echo_g "===================== Make transfer success ======================= "
    echo
}

chaincodeQueryA () {
    echo_b "Attempting to Query account A's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","4230a12f5b0693dd88bb35c79d7e56a68614b199","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account A Failed."
}

chaincodeQueryB () {
    echo_b "Attempting to  query account B's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","3c97f146e8de9807ef723538521fcecd5f64c79a","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account B Failed."
   
}

echo_b "=====================6.Issue a token using ascc========================"
issueToken INK

echo_b "=====================7.Transfer 100 amount of INK====================="
makeTransfer

echo_b "=====================8.Query transfer result of From account====================="
#checkTransferRes1
chaincodeQueryA

echo_b "=====================9.Query transfer result of To account====================="
#checkTransferRes2
chaincodeQueryB

echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

