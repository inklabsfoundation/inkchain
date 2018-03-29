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

# addUser $USER_TOKEN_01 hanmeimei hanmeimei@qq.com
# addUser $USER_TOKEN_02 lilei lilei@qq.com
addUser(){
    tag="add user"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["AddUser","'$2'","'$3'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative $tag successfully======================= "
    echo
}

# deleteUser $USER_TOKEN_01 hanmeimei
# deleteUser $USER_TOKEN_02 lilei
deleteUser () {
    tag="delete user"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["DeleteUser","'$2'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative $tag successfully======================= "
    echo
}

# modifyUser $USER_TOKEN_01 hanmeimei Email meimei@qq.com
# modifyUser $USER_TOKEN_02 lilei Email leilei@qq.com
modifyUser () {
    tag="modify user"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["ModifyUser","'$2'","'$3'","'$4'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative $tag successfully======================= "
    echo
}

# queryUser hanmeimei
# queryUser lilei
queryUser () {
    tag="query user"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["QueryUser","'$1'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative $tag successfully======================= "
    echo
}

# listOfUser
listOfUser () {
    tag="list of user"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["ListOfUser"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}