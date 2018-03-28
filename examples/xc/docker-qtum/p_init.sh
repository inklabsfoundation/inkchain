#! /bin/bash

# INK
INKHexAddress="bd557ae7ea1030262af6cfaa3d0e1b556dc61d6d"
INKAddress=`qcli fromhexaddress $INKHexAddress`
INKOwner="qKSb9J7F5RHGgtnsgZTJ3rvM9GF29oE9Dz"
INKHexOwner=`qcli gethexaddress $INKOwner`

# XC
XCHexAddress="b3ba697c5ff1e0bff5bf17f3bebbb0ccf9234ca9"
XCAddress=`qcli fromhexaddress ${XCHexAddress}`
XCOwner="qVuZrmVB4mALHEbYMSBH63Vm4ysfDBxR5B"
XCHexOwner=`qcli gethexaddress ${XCOwner}`

# XCPlugin
XCPluginHexAddress="dbed340e3a3d8bfe2aacbc2289ac3dfa2b27f93b"
XCPluginAddress=`qcli fromhexaddress ${XCPluginHexAddress}`
XCPluginOwner="qUgv6rtYnhPZmhsb1U3pSAFNsX4d82eXf6"
XCPluginHexOwner=`qcli gethexaddress ${XCPluginOwner}`
