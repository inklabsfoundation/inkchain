#!/usr/bin/env bash
#
#Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

## DO NOT MODIFY THE FOLLOWING PART, UNLESS YOU KNOW WHAT IT MEANS ##
echo_r () {
    [ $# -ne 1 ] && return 0
    echo -e "\033[31m$1\033[0m"
}
echo_g () {
    [ $# -ne 1 ] && return 0
    echo -e "\033[32m$1\033[0m"
}
echo_y () {
    [ $# -ne 1 ] && return 0
    echo -e "\033[33m$1\033[0m"
}
echo_b () {
    [ $# -ne 1 ] && return 0
    echo -e "\033[34m$1\033[0m"
}

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="mychannel"}
: ${TIMEOUT:="60"}
COUNTER=0
MAX_RETRY=5

CC_PATH=github.com/inklabsfoundation/inkchain/examples
ORDERER_CA=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
ARTIFACTS_PWD=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/channel-artifacts
SERVICE_CHARGE="10"

verifyResult () {
    if [ $1 -ne 0 ] ; then
        echo_b "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
        echo_r "================== ERROR !!! FAILED to execute MVE =================="
        echo
        exit 1
    fi
}

USER_TOKEN_01="70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4"
USER_ADDRESS_01="07caf88941eafcaaa3370657fccc261acb75dfba"

USER_TOKEN_02="bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe"
USER_ADDRESS_02="4230a12f5b0693dd88bb35c79d7e56a68614b199"

USER_TOKEN_03="344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5"
USER_ADDRESS_03="a5ff00eb44bf19d5dfbde501c90e286badb58df4"

createChannel() {
    peer channel create -o orderer.example.com:7050 -c ${CHANNEL_NAME} -f ${ARTIFACTS_PWD}/mychannel.tx --tls $CORE_PEER_TLS_ENABLED --cafile ${ORDERER_CA} >&log.txt
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
    peer channel create -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -f ${ARTIFACTS_PWD}/Org1MSPanchors.tx >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Anchor peer update failed"
    echo_g "==== Anchor peers for org1 on mychannel is updated successfully======"
    echo
}

# installChaincode token 1.0
# installChaincode creative 1.0
installChaincode () {
    peer chaincode install -n $1 -v $2 -p ${CC_PATH}/$1 -o orderer.example.com:7050 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode token installation on remote peer0 has Failed"
    echo_g "===================== Chaincode is installed success on remote peer0===================== "
    echo
}

# instantiateChaincode token 1.0
# instantiateChaincode creative 1.0
instantiateChaincode () {
    local starttime=$(date +%s)
    peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n $1 -v $2 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member')" >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode instantiation on pee0.org1 on channel "$CHANNEL_NAME" failed"
    echo_g "=========== Chaincode token Instantiation on peer0.org1 on channel "$CHANNEL_NAME" is successful ========== "
    echo_b "Instantiate spent $(($(date +%s)-starttime)) secs"
    echo
}

# upgradeChaincode creative 2.0
upgradeChaincode () {
    local starttime=$(date +%s)
    peer chaincode upgrade -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n $1 -v $2 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member')" >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode instantiation on pee0.org1 on channel "$CHANNEL_NAME" failed"
    echo_g "=========== Chaincode token Instantiation on peer0.org1 on channel "$CHANNEL_NAME" is successful ========== "
    echo_b "Instantiate spent $(($(date +%s)-starttime)) secs"
    echo
}
