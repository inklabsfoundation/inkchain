#! /bin/bash

source p_init.sh

qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol unlockAdmin '["'$1'","'$2'","'$3'","'$4'",'$5']'` 0 6000000 $gasPrice $XCOwner
