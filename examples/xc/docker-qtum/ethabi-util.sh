#! /bin/bash

echo '      ############################'
echo '     #     1:encode             #'
echo '    #     2:decode             #'
echo '   #     3: encode function   #'
echo '  ############################'
echo ' # enter [1~3]:'

read aNum

case ${aNum} in
    1)
    echo "encode data enter :"
    read data
    docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest encode params -v string ${data} --lenient
    ;;
    2)
    echo "decode data enter :"
    read data
    docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest decode params -t string ${data}
    ;;
    3)
    echo "encode data enter (1:INK 2:XCPlugin 3:XC):"
    read contract
    read method
    case ${contract} in
    1)
    docker run --rm -v ${PWD}/../contracts/src/ink-token:/ethabi cryptominder/ethabi:latest encode function /ethabi/INK.abi ${method}
    ;;
    2)
    docker run --rm -v ${PWD}/../contracts/src/xc-plugin:/ethabi cryptominder/ethabi:latest encode function /ethabi/XCPlugin.abi ${method}
    ;;
    3)
    docker run --rm -v ${PWD}/../contracts/src/ink-xc:/ethabi cryptominder/ethabi:latest encode function /ethabi/XC.abi ${method}
    ;;
    esac
esac
