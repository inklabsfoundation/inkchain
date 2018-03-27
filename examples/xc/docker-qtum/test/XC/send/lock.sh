#! /bin/bash

source p_init.sh

qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol lock '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'` 0 6000000 0.0000004 $XCOwner
