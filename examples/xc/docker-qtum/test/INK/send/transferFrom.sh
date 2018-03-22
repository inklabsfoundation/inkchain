#! /bin/bash

source p_init.sh

#qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol transferFrom '["'$INKHexOwner'","'$XCHexAddress'",500]'` 0 6000000 0.0000004 $XCAddress
qcli sendtocontract $INKHexAddress `solar encode contracts/INK.sol transferFrom '["'$INKHexOwner'","d2e12af1eda54d3a7b29fb1f53813819733697e4",500]'` 0 6000000 0.0000004 qcnQoFCBc9xgaHhygBfSvJoZghSTdSxECQ

