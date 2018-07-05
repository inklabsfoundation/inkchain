#! /bin/bash

docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
        --optimize --bin --abi --hashes -o /solidity/src/INK  \
        --overwrite /solidity/INK.sol --evm-version homestead