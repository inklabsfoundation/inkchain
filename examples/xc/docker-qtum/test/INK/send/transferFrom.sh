#! /bin/bash

source p_init.sh

qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol transferFrom '["'$1'","'$2'",'$3']'` 0 6000000 $gasPrice $4
