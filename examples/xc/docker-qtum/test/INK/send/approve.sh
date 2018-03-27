#! /bin/bash

source p_init.sh

qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol approve '["'$XCHexAddress'",1000]'` 0 6000000 0.0000004 $INKOwner
#qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol approve '["'d2e12af1eda54d3a7b29fb1f53813819733697e4'",1000]'` 0 6000000 0.0000004 $INKOwner
