#! /bin/bash

source p_init.sh

qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol setWeight '["4100000000000000000000000000000000000000000000000000000000000000",3]'` 0 6000000 0.0000004 $XCPluginOwner