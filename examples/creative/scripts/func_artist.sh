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

# addArtist  $USER_TOKEN_01 "hanmeimei" "女作家" "韩梅梅，女，金牛座，1981年出生于云南。毕业于北京电影学院导演系，畅销书作家。"
# addArtist  $USER_TOKEN_02 "lilei" "程序员" "李雷，男，狮子座，1980年出生于山东。毕业于蓝翔，屌丝程序员。"
addArtist(){
    tag="add artist"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["AddArtist","'$2'","'$3'","'$4'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# deleteArtist $USER_TOKEN_01 "hanmeimei"
# deleteArtist $USER_TOKEN_02 "lilei"
deleteArtist () {
    tag="delete artist"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["DeleteArtist","'$2'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# TODO 可1次修改多值
# modifyArtist $USER_TOKEN_01 Name "女作家&女演员"
# modifyArtist $USER_TOKEN_01 Desc "韩梅梅，女，金牛座，1981年出生于云南。毕业于北京电影学院导演系，畅销书作家。TO DO +"
# modifyArtist $USER_TOKEN_02 Name "程序员&架构师"
# modifyArtist $USER_TOKEN_02 Desc "李雷，男，狮子座，1980年出生于山东。毕业于蓝翔，屌丝程序员。TO DO +"
modifyArtist () {
    tag="modify artist"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["ModifyArtist","'$2'","'$3'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# queryArtist "hanmeimei"
# queryArtist "lilei"
queryArtist () {
    tag="query artist"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["QueryArtist","'$1'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# listOfArtist
# listOfArtist
listOfArtist () {
    tag="list of artist"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["ListOfArtist"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}