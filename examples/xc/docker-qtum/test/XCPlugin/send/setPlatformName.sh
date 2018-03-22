#! /bin/bash

source p_init.sh

qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol setPlatformName '["7174756d00000000000000000000000000000000000000000000000000000000"]'` 0 6000000 0.0000004 $XCPluginOwner