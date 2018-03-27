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

# addProduction $USER_TOKEN_01 "hanmeimei" "book" "00000001" "《李雷和韩梅梅》" "青春偶像" "INK" "100" "10000" "1000"
# addProduction $USER_TOKEN_02 "lilei" "code" "00000001" "《PHP从出门到放弃》" "抓狂日志" "INK" "1" "10000" "5000"
addProduction(){
    tag="add production"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["AddProduction","'$2'","'$3'","'$4'","'$5'","'$6'","'$7'","'$8'","'$9'","'${10}'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# deleteProduction $USER_TOKEN_01 "hanmeimei" "book" "00000001"
# deleteProduction $USER_TOKEN_02 "lilei" "code" "00000001"
deleteProduction(){
    tag="delete production"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["DeleteProduction","'$2'","'$3'","'$4'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# Name、Desc、CopyrightPriceType、CopyrightPrice、CopyrightNum、CopyrightTransferPart
# modifyProduction $USER_TOKEN_01 "hanmeimei" "book" "00000001" "Name" "《老王和韩梅梅》"
# modifyProduction $USER_TOKEN_02 "lilei" "code" "00000001" "Name" "《再见！coding》"
modifyProduction(){
    tag="modify production"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["ModifyProduction","'$2'","'$3'","'$4'","'$5'","'$6'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# queryProduction "hanmeimei" "book" "00000001"
# queryProduction "lilei" "code" "00000001"
queryProduction () {
    tag="query production"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["QueryProduction","'$1'","'$2'","'$3'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# listOfProduction username product_type
# listOfProduction "hanmeimei" "book"
# listOfProduction "lilei" "code"
listOfProduction () {
    tag="list of production"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["ListOfProduction","'$1'","'$2'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# addSupporter $USER_TOKEN_01 "lilei" "code" "00000001" "INK" "99" "hanmeimei"
# addSupporter $USER_TOKEN_02 "hanmeimei" "book" "00000001" "INK" "199" "lilei"
addSupporter(){
    tag="add supporter"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["AddSupporter","'$2'","'$3'","'$4'","'$5'","'$6'","'$7'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# listOfSupporter username product_type  product_serial
# listOfSupporter "hanmeimei" "book" "00000001"
# listOfSupporter "lilei" "code" "00000001"
listOfSupporter(){
    tag="list of supporter"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["ListOfSupporter","'$1'","'$2'","'$3'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# addBuyer $USER_TOKEN_01 "lilei" "code" "00000001" "INK" "5" "hanmeimei"
# addBuyer $USER_TOKEN_01 "lilei" "code" "00000001" "INK" "5" "hanmeimei"
addBuyer(){
    tag="add buyer"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["AddBuyer","'$2'","'$3'","'$4'","'$5'","'$6'","'$7'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# modifyBuyer
modifyBuyer(){
    tag="modify buyer"
    echo_b "Attempting to $tag"
#    sleep 3
#    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["ModifyBuyer","'$2'","'$3'","'$4'","'$5'","'$6'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
#    res=$?
#    cat log.txt
#    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# deleteBuyer
deleteBuyer(){
    tag="delete buyer"
    echo_b "Attempting to $tag"
#    sleep 3
#    peer chaincode invoke -C $CHANNEL_NAME -n creative --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["DeleteBuyer","'$2'","'$3'","'$4'","'$5'","'$6'"]}' -i $SERVICE_CHARGE -z $1 >&log.txt
#    res=$?
#    cat log.txt
#    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}

# listOfBuyer username product_type  product_serial
# listOfBuyer "hanmeimei" "book" "00000001"
# listOfBuyer "lilei" "code" "00000001"
listOfBuyer(){
    tag="list of buyer"
    echo_b "Attempting to $tag"
    sleep 3
    peer chaincode query -C $CHANNEL_NAME -n creative -c '{"Args":["ListOfBuyer","'$1'","'$2'","'$3'"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "$tag: Dainel Failed."
    echo_g "===================== creative  $tag successfully======================= "
    echo
}