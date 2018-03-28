#! /bin/bash

solc -o ${PWD}/src/INK --bin --abi ${PWD}/INK.sol
solc -o ${PWD}/src/XCPlugin --bin --abi ${PWD}/XCPlugin.sol
solc -o ${PWD}/src/XC --bin --abi ${PWD}/XC.sol
