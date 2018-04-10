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
    echo "encode data type enter "
    read type
    echo "encode data enter :"
    read data
    docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest encode params -v ${type} ${data} --lenient
    ;;
    2)
    echo "decode data type enter "
    read type
    echo "decode data enter :"
    read data
    docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest decode params -t ${type} ${data}
    ;;
    3)
    echo "encode data enter (1:INK 2:XCPlugin 3:XC):"
    read contract
    echo "encode method data enter:"
    read method

    case ${contract} in
    1)
    docker run --rm -v ${PWD}/contracts/src/INK:/ethabi cryptominder/ethabi:latest encode function /ethabi/INK.abi ${method} --lenient
    ;;
    2)
    docker run --rm -v ${PWD}/contracts/src/XCPlugin:/ethabi cryptominder/ethabi:latest encode function /ethabi/XCPlugin.abi ${method} --lenient
    ;;
    3)
    docker run --rm -v ${PWD}/contracts/src/XC:/ethabi cryptominder/ethabi:latest encode function /ethabi/XC.abi ${method} --lenient
    ;;
    esac
esac
