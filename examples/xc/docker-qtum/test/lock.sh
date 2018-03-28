#! /bin/bash

source p_init.sh

# balanceOf
qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$INKHexOwner'"]'`
qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$XCHexAddress'"]'`

# approve
qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol approve '["'$XCHexAddress'",1000]'` 0 6000000 0.0000004 $INKOwner
sleep 10

# allowance
qcli callcontract $INKHexAddress `solar encode contracts/INK.sol allowance '["'$INKHexOwner'","'$XCHexAddress'"]'`

# lock d6b39eb631df8ee60e46a576231ccf1fcd204a5e xc`s toAccount
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol lock '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",500]'` 0 6000000 0.0000004 $INKOwner
sleep 20

# balanceOf
qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$INKHexOwner'"]'`
qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$XCHexAddress'"]'`