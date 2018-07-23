#! /bin/bash

source p_init.sh

qcli callcontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol getPlatformName` $XCPluginOwner