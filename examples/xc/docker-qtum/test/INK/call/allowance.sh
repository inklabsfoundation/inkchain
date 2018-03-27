#! /bin/bash

source p_init.sh

qcli callcontract $INKHexAddress `solar encode contracts/INK.sol allowance '["'$INKHexOwner'","'$XCHexAddress'"]'`
#qcli callcontract $INKHexAddress `solar encode contracts/INK.sol allowance '["'$INKHexOwner'","'d2e12af1eda54d3a7b29fb1f53813819733697e4'"]'`
