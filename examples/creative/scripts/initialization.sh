#!/usr/bin/env bash
#
#Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

# Detecting whether can import the header file to render colorful cli output
if [ -f ./func_user.sh ]; then
 source ./func_user.sh
elif [ -f scripts/func_user.sh ]; then
 source scripts/func_user.sh
fi

if [ -f ./func_artist.sh ]; then
 source ./func_artist.sh
elif [ -f scripts/func_artist.sh ]; then
 source scripts/func_artist.sh
fi

if [ -f ./func_production.sh ]; then
 source ./func_production.sh
elif [ -f scripts/func_production.sh ]; then
 source scripts/func_production.sh
fi

if [ -f ./func_token.sh ]; then
 source ./func_token.sh
elif [ -f scripts/func_token.sh ]; then
 source scripts/func_token.sh
fi


#echo_b "====================1.Create channel(default newchannel) ============================="
#createChannel
#
#echo_b "====================2.Join pee0 to the channel ======================================"
#joinChannel
#
#echo_b "====================3.set anchor peers for org1 in the channel==========================="
##updateAnchorPeers
#
#echo_b "=====================4.Install chaincode token on Peer0/Org0========================"
#installChaincode token 1.0
#installChaincode creative 1.0
#
#echo_b "=====================5.Instantiate chaincode token, this will take a while, pls waiting...==="
#instantiateChaincode token 1.0
#instantiateChaincode creative 1.0


############################################################
# installChaincode creative 1.6
# upgradeChaincode creative 1.6
