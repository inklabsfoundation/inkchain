#! /bin/bash

source p_init.sh

# voter(bytes32,address,address,uint256,bytes32,bytes32,bytes32,uint8)
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol voteProposal '["'$1'","'$2'","'$3'",'$4',"'$5'","'$6'","'$7'",'$8']'` 0 6000000 $gasPrice $XCPluginOwner
