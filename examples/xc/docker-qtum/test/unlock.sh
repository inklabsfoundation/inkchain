#! /bin/bash

source p_init.sh

# bytes32 txId, bytes32 fromPlatform, address fromAccount, address toAccount, uint amount
qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol unlock '["0xedebdb5c8feffb95850628ea17816aa8579e1070e2106f619dd733c1d0e89b7c","4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",900]'` 0 6000000 0.0000004 $XCOwner
