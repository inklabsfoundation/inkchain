#! /bin/bash

source p_init.sh

#qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$INKHexOwner'"]'`

#qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$XCHexAddress'"]'`

qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$XCHexOwner'"]'`

#qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]' `