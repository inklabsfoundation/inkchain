#! /bin/bash

echo '      ############################'
echo '     #     1: ink-token.sol     #'
echo '    #     2: xc-plugin.sol     #'
echo '   #     3: ink-xc.sol        #'
echo '  ############################'
echo ' # (default:All) Enter [1~3]:'
read aNum

case $aNum in
    '') docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/ink-token  \
                    --overwrite /solidity/ink-token.sol

        docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/xc-plugin  \
                    --overwrite /solidity/xc-plugin.sol

        docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/ink-xc  \
                    --overwrite /solidity/ink-xc.sol
    ;;
    1) docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/ink-token  \
                    --overwrite /solidity/ink-token.sol
    ;;
    2) docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/xc-plugin  \
                    --overwrite /solidity/xc-plugin.sol
    ;;
    3) docker run --rm -v ${PWD}:/solidity ethereum/solc:stable  \
                    --optimize --bin --abi --hashes -o /solidity/src/ink-xc  \
                    --overwrite /solidity/ink-xc.sol
    ;;
esac
