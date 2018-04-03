#! /bin/bash

source p_init.sh

# voter(bytes32,address,address,uint256,bytes32   ,bytes32,bytes32,uint8)
r=""
s=""
v=27
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol voter '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e","'$XCHexOwner'",900,"0xedebdb5c8feffb95850628ea17816aa8579e1070e2106f619dd733c1d0e89b7c","'$r'","'$s'",'$v']'` 0 6000000 0.0000004 $XCPluginOwner
