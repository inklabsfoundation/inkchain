#!/bin/bash
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
TOKEN_CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/token
MARBLES_CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/marbles
ASSET_CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/asset
CAT_CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/cat
INKWORK_CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/inkwork
ORDERER_CA=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo_b "Chaincode Path : " $CC_PATH
echo_b "Channel name : " $CHANNEL_NAME

verifyResult () {
    if [ $1 -ne 0 ] ; then
        echo_b "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
        echo_r "================== ERROR !!! FAILED to execute MVE =================="
        echo
        exit 1
    fi
}

createChannel() {
    peer channel create -o orderer.example.com:7050 -c ${CHANNEL_NAME} -f ./channel-artifacts/mychannel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
    res=$?
    cat log.txt

    verifyResult $res "Channel creation failed"
    echo
    # verify file mychannel.block exist
    if [ -s mychannel.block ]; then
        res=$?
        verifyResult $res "Channel created failed"
    fi
        echo_g "================channel \"$CHANNEL_NAME\" is created successfully ==============="
}

## Sometimes Join takes time hence RETRY atleast for 5 times

joinChannel () {
    echo_b "===================== PEER0 joined on the channel \"$CHANNEL_NAME\" ===================== "
    peer channel join -b ${CHANNEL_NAME}.block -o orderer.example.com:7050 >&log.txt
    res=$?
    cat log.txt
    if [ $res -ne 0 -a $COUNTER -lt $MAX_RETRY ]; then
        COUNTER=` expr $COUNTER + 1`
        echo_r "PEER0 failed to join the channel, Retry after 2 seconds"
        sleep 2
        joinWithRetry
    else
        COUNTER=0
    fi
        verifyResult $res "After $MAX_RETRY attempts, PEER0 has failed to Join the Channel"
}

updateAnchorPeers() {
    peer channel create -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -f ./channel-artifacts/Org1MSPanchors.tx >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Anchor peer update failed"
    echo_g "==== Anchor peers for org1 on mychannel is updated successfully======"
    echo
}

installChaincode () {
    peer chaincode install -n token -v 1.0 -p ${TOKEN_CC_PATH} -o orderer.example.com:7050 >&log.txt

    #peer chaincode install -n marbles -v 1.0 -p ${MARBLES_CC_PATH} -o orderer.example.com:7050 >&log.txt

    peer chaincode install -n asset -v 1.0 -p ${ASSET_CC_PATH} -o orderer.example.com:7050 >&log.txt

    #peer chaincode install -n cat -v 1.0 -p ${CAT_CC_PATH} -o orderer.example.com:7050 >&log.txt

    #peer chaincode install -n inkwork -v 1.0 -p ${INKWORK_CC_PATH} -o orderer.example.com:7050 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode token installation on remote peer0 has Failed"
    echo_g "===================== Chaincode is installed success on remote peer0===================== "
    echo
}

instantiateChaincode () {
    local starttime=$(date +%s)
    peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n token -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member')" >&log.txt
    #peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n marbles -v 1.0 -c '{"Args":["initMarble","marble1","blue","35","tom"]}' -P "OR ('Org1MSP.member')" >&log.txt
    #peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n marbles -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member')" >&log.txt
    peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n asset -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member')" >&log.txt
    #peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n cat -v 1.0 -c '{"Args":["init","5","5","6","60","INK","i07caf88941eafcaaa3370657fccc261acb75dfba"]}' -P "OR ('Org1MSP.member')" >&log.txt
    #peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n inkwork -v 1.0 -c '{"Args":["init","i07caf88941eafcaaa3370657fccc261acb75dfba","INK"]}' -P "OR ('Org1MSP.member')" >&log.txt

    res=$?
    cat log.txt
    verifyResult $res "Chaincode instantiation on pee0.org1 on channel '$CHANNEL_NAME' failed"
    echo_g "=========== Chaincode token Instantiation on peer0.org1 on channel '$CHANNEL_NAME' is successful ========== "
    echo_b "Instantiate spent $(($(date +%s)-starttime)) secs"
    echo
}

echo_b "====================1.Create channel(default newchannel) ============================="
createChannel

echo_b "====================2.Join pee0 to the channel ======================================"
joinChannel

echo_b "====================3.set anchor peers for org1 in the channel==========================="
#updateAnchorPeers

echo_b "=====================4.Install chaincode token on Peer0/Org0========================"
installChaincode

echo_b "=====================5.Instantiate chaincode token, this will take a while, pls waiting...==="
instantiateChaincode

echo
echo_g "=====================All GOOD, MVE initialization completed ===================== "
echo
exit 0
