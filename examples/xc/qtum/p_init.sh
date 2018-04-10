#! /bin/bash

#QTUM
myPlatform="5154554d00000000000000000000000000000000000000000000000000000000"

#INK
otherPlatform="494e4b0000000000000000000000000000000000000000000000000000000000"
otherPlatformPubKey="4230a12f5b0693dd88bb35c79d7e56a68614b199"
otherPlatformPubKey2="07caf88941eafcaaa3370657fccc261acb75dfba"

gasPrice=0.0000004

INKHexAddress="13ea393938e637747960a769dd2883d278dd8d64"
INKAddress=`qcli fromhexaddress $INKHexAddress`
INKOwner="qf4FJzETksomRxgXvTLqvZwAVihjAe1ifE"
INKHexOwner=`qcli gethexaddress $INKOwner`

XCHexAddress="7456a68540ab5a0eaa82cd00fc2cea64f412cd7c"
XCAddress=`qcli fromhexaddress ${XCHexAddress}`
XCOwner="qJsUXqnd71btLaKLy2Cnf4p5up1GiyfCwX"
XCHexOwner=`qcli gethexaddress ${XCOwner}`

XCPluginHexAddress="32fbd6b99a0d7b733a826ec738d88618340a3308"
XCPluginAddress=`qcli fromhexaddress ${XCPluginHexAddress}`
XCPluginOwner="qdkVBzhos1N9KN7wmfGdRmr9G8PSnTCtfm"
XCPluginHexOwner=`qcli gethexaddress ${XCPluginOwner}`

#✅  contracts/INK.sol
#        txid: 544b8d54b8d39c4f3c4073b327548536fd0172416ca37d50b182e3d97eff5ddf
#     address: 13ea393938e637747960a769dd2883d278dd8d64
#   confirmed: true
#       owner: qf4FJzETksomRxgXvTLqvZwAVihjAe1ifE
#
#✅  contracts/XCPlugin.sol
#        txid: 60471cb61f7b91c94dfcc76a65a2dde9f74cf7b4e41e3885b3c2739e3bab8be0
#     address: 32fbd6b99a0d7b733a826ec738d88618340a3308
#   confirmed: true
#       owner: qdkVBzhos1N9KN7wmfGdRmr9G8PSnTCtfm
#
#✅  contracts/XC.sol
#        txid: 009b7f44a710865e3e7dde7978b929c3918cd9e8e45c830933d3111e7370ba9e
#     address: 7456a68540ab5a0eaa82cd00fc2cea64f412cd7c
#   confirmed: true
#       owner: qJsUXqnd71btLaKLy2Cnf4p5up1GiyfCwX


