#! /bin/bash

source p_init.sh

qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol setAdmin '["'$1'"]'` 0 6000000 $gasPrice $XCOwner
