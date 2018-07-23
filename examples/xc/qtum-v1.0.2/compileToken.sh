#! /bin/bash

docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
        --optimize --bin --abi --hashes -o /solidity/src/Token  \
        --overwrite /solidity/Token.sol --evm-version homestead