#! /bin/bash

source p_init.sh

qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$INKHexOwner'"]'`

#qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$XCHexAddress'"]'`
