#! /bin/bash

source p_init.sh

# === XCPlugin ===

# .set Owner PlatformName  default "QTUM"
#qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol setPlatformName '["'${myPlatform}'"]'` 0 6000000 0.000001 $XCPluginOwner

# .add credible xc Platform （"INK"）
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addPlatform '["'${otherPlatform}'"]'` 0 300000 0.000001 $XCPluginOwner
sleep 180

# .add credible xc Platform's PublicKey
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addPublicKey '["'${otherPlatform}'","'$otherPlatformPubKey'"]'` 0 300000 0.000001 $XCPluginOwner
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addPublicKey '["'${otherPlatform}'","'$otherPlatformPubKey2'"]'` 0 300000 0.000001 $XCPluginOwner
sleep 180

#
## .set  credible xc Platform's weight, default 1
##qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol setWeight '["'${otherPlatform}'",1]'` 0 6000000 0.000001 $XCPluginOwner
#
## set XCPlugin contract caller.
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addCaller '["'$XCHexAddress'"]'` 0 300000 0.000001 $XCPluginOwner
#
## === XC ===
#
## .set INK contract address
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setINK '["'$INKHexAddress'"]'` 0 300000 0.000001 $XCOwner
## .set XCPlugin contract address
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setXCPlugin '["'$XCPluginHexAddress'"]'` 0 300000 0.000001 $XCOwner
#
## .set Owner PlatformName  default "QTUM"
##qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setPlatformName '["'${myPlatform}'"]'` 0 6000000 0.000001 $XCOwner
#
## === start ===
#
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol start` 0 300000 0.000001 $XCPluginOwner
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setStatus [3]` 0 300000 0.000001 $XCOwner
