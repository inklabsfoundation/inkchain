#! /bin/bash

source p_init.sh

# verifyProposal(bytes32,address,address,uint256,bytes32)
qcli callcontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol verifyProposal '["'$1'","'$2'","'$3'",'$4',"'$5'"]'`
