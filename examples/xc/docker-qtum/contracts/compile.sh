#! /bin/bash

echo '      ############################'
echo '     #     1: INK.sol     #'
echo '    #     2: XCPlugin.sol     #'
echo '   #     3: XC.sol        #'
echo '  ############################'
echo ' # (default:All) Enter [1~3]:'
read aNum

case $aNum in
    '') docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/INK  \
                    --overwrite /solidity/INK.sol

        docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/XCPlugin  \
                    --overwrite /solidity/XCPlugin.sol

        docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/XC  \
                    --overwrite /solidity/XC.sol
    ;;
    1) docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/INK  \
                    --overwrite /solidity/INK.sol
    ;;
    2) docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/XCPlugin  \
                    --overwrite /solidity/XCPlugin.sol
    ;;
    3) docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/XC  \
                    --overwrite /solidity/XC.sol
    ;;
esac
