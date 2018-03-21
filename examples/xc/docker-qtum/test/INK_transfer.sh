#! /bin/bash

source p_init.sh

qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol transfer '["'$XCHexAddress'",1000]'`