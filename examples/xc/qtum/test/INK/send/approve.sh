#! /bin/bash

source p_init.sh

qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol approve '["'$1'",'$2']'` 0 6000000 $gasPrice $INKOwner
