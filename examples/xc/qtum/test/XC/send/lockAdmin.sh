#! /bin/bash

source p_init.sh

qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol lockAdmin '["'$1'","'$2'",'$3']'` 0 6000000 $gasPrice $XCOwner
