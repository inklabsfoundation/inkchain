#! /bin/bash

source p_init.sh

qcli callcontract $XCHexAddress `solar encode contracts/XC.sol getINK`