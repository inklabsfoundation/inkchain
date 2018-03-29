#! /bin/bash

source p_init.sh

qcli callcontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol verifyProposal '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000,"4200000000000000000000000000000000000000000000000000000000000000"]'`
# verifyProposal(bytes32,address,address,uint256,bytes32)