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

# example

#echo_b "===================== issueToken ========================"
#issueToken INK "1000" "18" ${USER_ADDRESS_01}
#echo_b "===================== makeTransfer ========================"
#makeTransfer ${USER_ADDRESS_02} INK  500 ${USER_TOKEN_01}
#echo_b "===================== chaincodeQuery ========================"
chaincodeQuery $USER_ADDRESS_01 INK
chaincodeQuery $USER_ADDRESS_02 INK

#echo_b "===================== add user ========================"
# addUser $USER_TOKEN_01 hanmeimei hanmeimei@qq.com
# addUser $USER_TOKEN_02 lilei lilei@qq.com

#echo_b "===================== delete user ========================"
# deleteUser $USER_TOKEN_01 hanmeimei
# deleteUser $USER_TOKEN_02 lilei

#echo_b "===================== modify user 可1次修改多值 "Name,xxx" "Desc,xxx" 用 ',' 分割键值  ========================"
# modifyUser $USER_TOKEN_01 "hanmeimei" "Email,meimei@qq.com"
# modifyUser $USER_TOKEN_02 "lilei" "Email,leilei@qq.com"
#
#echo_b "===================== query user ========================"
# queryUser hanmeimei
# queryUser lilei

#echo_b "===================== list user ========================"
# listOfUser

#echo_b "===================== addArtist ========================"
# addArtist  $USER_TOKEN_01 "hanmeimei" "女作家" "韩梅梅，女，金牛座，1981年出生于云南。毕业于北京电影学院导演系，畅销书作家。"
# addArtist  $USER_TOKEN_02 "lilei" "程序员" "李雷，男，狮子座，1980年出生于山东。毕业于蓝翔，屌丝程序员。"

# deleteArtist $USER_TOKEN_01 "hanmeimei"
# deleteArtist $USER_TOKEN_02 "lilei"

# 可1次修改多值 "Name,xxx" "Desc,xxx" 用 ',' 分割键值
# modifyArtist $USER_TOKEN_01 "hanmeimei" "Name,女作家、女演员" "Desc,韩梅梅"
# modifyArtist $USER_TOKEN_02 "lilei"  "Name,程序员&架构师" "Desc,李雷，男，狮子座，1980年出生于山东。毕业于蓝翔，屌丝程序员。TO DO +"
#
# queryArtist "hanmeimei"
# queryArtist "lilei"

# listOfArtist

# addProduction $USER_TOKEN_01 "hanmeimei" "book" "00000001" "《李雷和韩梅梅》" "青春偶像" "INK" "100" "10000" "1000"
# addProduction $USER_TOKEN_02 "lilei" "code" "00000001" "《PHP从出门到放弃》" "抓狂日志" "INK" "1" "10000" "5000"

# deleteProduction $USER_TOKEN_01 "hanmeimei" "book" "00000001"
# deleteProduction $USER_TOKEN_02 "lilei" "code" "00000001"

# Name、Desc、CopyrightPriceType、CopyrightPrice、CopyrightNum、CopyrightTransferPart 可1次修改多值 "Name,xxx" "Desc,xxx" 用 ',' 分割键值
# modifyProduction $USER_TOKEN_01 "hanmeimei" "book" "00000001" "Name,《老王,韩梅梅》" "Desc,xxxxx"
# modifyProduction $USER_TOKEN_02 "lilei" "code" "00000001" "Name,《再见！coding》"
##
# queryProduction "hanmeimei" "book" "00000001"
# queryProduction "lilei" "code" "00000001"


# listOfProduction "username" "product_type" 参数可变
# listOfProduction "hanmeimei" "book"
# listOfProduction "lilei" "code"
# listOfProduction "hanmeimei"
# listOfProduction "lilei"
# listOfProduction

# addSupporter $USER_TOKEN_01 "lilei" "code" "00000001" "INK" "99" "hanmeimei"
# addSupporter $USER_TOKEN_02 "hanmeimei" "book" "00000001" "INK" "199" "lilei"

# listOfSupporter "username" "product_type"  "product_serial" 参数可变
# listOfSupporter "hanmeimei" "book" "00000001"
# listOfSupporter "lilei" "code" "00000001"
#
# addBuyer $USER_TOKEN_02 "hanmeimei" "book" "00000001" "INK" "1" "lilei"
# addBuyer $USER_TOKEN_01 "lilei" "code" "00000001" "INK" "5" "hanmeimei"
#
# queryProduction "hanmeimei" "book" "00000001"
# queryProduction "lilei" "code" "00000001"

# listOfBuyer username product_type  product_serial
# listOfBuyer "hanmeimei" "book" "00000001"
# listOfBuyer "lilei" "code" "00000001"