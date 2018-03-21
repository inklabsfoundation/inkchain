#! /bin/bash

INKHexAddress="223a30419a3dc234ec659e390d7bc65e70dcea0b"
INKAddress=`qcli fromhexaddress $INKHexAddress`
INKOwner="qP2RGcwGunqiZvSUidSQyEW64G1XuyGUyM"
INKHexOwner=`qcli gethexaddress $INKOwner`


XCHexAddress="54acd8f45dad524584003ee4d4f4d4f8b00ef4c1"
XCAddress=`qcli fromhexaddress ${XCHexAddress}`
XCOwner="qMqd6usqYmUXpTB3FShUiKP59TfMrXFQJf"
XCHexOwner=`qcli gethexaddress ${XCOwner}`


XCPluginHexAddress="fdfe183ff6d196a2871d3b856629222f1dd1264b"
XCPluginAddress=`qcli fromhexaddress ${XCPluginHexAddress}`
XCPluginOwner="qf2FFX6c5qohyqPvbUoyy8Xo4Q6QrjsaKX"
XCPluginHexOwner=`qcli gethexaddress ${XCPluginOwner}`