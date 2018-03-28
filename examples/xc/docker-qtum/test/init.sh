#! /bin/bash

source p_init.sh

# ===  init amount ===
qcli sendtoaddress $XCPluginOwner 20000
qcli sendtoaddress $XCOwner 20000
qcli sendtoaddress $XCAddress 20000

# === XCPlugin ===

# .set Owner PlatformName  default "qtum"
#qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol setPlatformName '["7174756d00000000000000000000000000000000000000000000000000000000"]'` 0 6000000 0.0000004 $XCPluginOwner

# .add credible xc Platform （"A","B"）
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addPlatform '["4100000000000000000000000000000000000000000000000000000000000000"]'` 0 6000000 0.0000004 $XCPluginOwner
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addPlatform '["4200000000000000000000000000000000000000000000000000000000000000"]'` 0 6000000 0.0000004 $XCPluginOwner

# .add credible xc Platform's PublicKey
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addPublicKey '["4100000000000000000000000000000000000000000000000000000000000000","'$XCHexAddress'"]'` 0 6000000 0.0000004 $XCPluginOwner
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addPublicKey '["4200000000000000000000000000000000000000000000000000000000000000","'$XCHexAddress'"]'` 0 6000000 0.0000004 $XCPluginOwner

# .set  credible xc Platform's weight, default 1
#qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol setWeight '["4100000000000000000000000000000000000000000000000000000000000000",1]'` 0 6000000 0.0000004 $XCPluginOwner

# set XCPlugin contract caller.
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol addCaller '["'$XCHexAddress'"]'` 0 6000000 0.0000004 $XCPluginOwner

# === XC ===

# .set INK contract address
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setINK '["'$INKHexAddress'"]'` 0 6000000 0.0000004 $XCOwner
# .set XCPlugin contract address
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setXCPlugin '["'$XCPluginHexAddress'"]'` 0 6000000 0.0000004 $XCOwner

# .set Owner PlatformName  default "qtum"
#qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setPlatformName '["7174756d00000000000000000000000000000000000000000000000000000000"]'` 0 6000000 0.0000004 $XCOwner

# === start ===

qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol start` 0 6000000 0.0000004 $XCPluginOwner
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol start` 0 6000000 0.0000004 $XCOwner
