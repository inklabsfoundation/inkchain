#! /bin/bash

docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
        --optimize --bin --abi --hashes -o /solidity/src/XCPlugin  \
        --overwrite /solidity/XCPlugin.sol --evm-version homestead