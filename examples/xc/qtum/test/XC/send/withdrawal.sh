#! /bin/bash

source p_init.sh

qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol withdrawal '["'$1'",'$2']'` 0 6000000 $gasPrice $XCOwner
