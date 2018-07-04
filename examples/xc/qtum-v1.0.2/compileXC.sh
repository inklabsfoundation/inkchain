#! /bin/bash

docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
        --optimize --bin --abi --hashes -o /solidity/src/XC  \
        --overwrite /solidity/XC.sol --evm-version homestead