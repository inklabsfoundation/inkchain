#! /bin/bash

source p_init.sh

qcli callcontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol getWeight '["4100000000000000000000000000000000000000000000000000000000000000"]'`