#! /bin/bash

INKHexAddress="3a1ee4aee3fb23ac65874cba182c98d18ccbae9e"
INKAddress=`qcli fromhexaddress $INKHexAddress`
INKOwner="qejjheVXroaEcyqEDrdF7HSdAGCR6EveeX"
INKHexOwner=`qcli gethexaddress $INKOwner`


XCHexAddress="4efd25089b49f3ddefdbebf7c7c4d8e7fa9a86d5"
XCAddress=`qcli fromhexaddress ${XCHexAddress}`
XCOwner="qbhM5imeTbZc3GavqJLDebbXNpzzHu54y5"
XCHexOwner=`qcli gethexaddress ${XCOwner}`


XCPluginHexAddress="98553eeaf215c1c586ffd2453ac3a6eff78cb541"
XCPluginAddress=`qcli fromhexaddress ${XCPluginHexAddress}`
XCPluginOwner="qLRqiZWbnj1s667W6M67iPeamWVMHoBiay"
XCPluginHexOwner=`qcli gethexaddress ${XCPluginOwner}`
