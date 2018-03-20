#! /bin/bash

solar deploy contracts/INK.sol

solar deploy contracts/XCPlugin.sol '["7174756d00000000000000000000000000000000000000000000000000000000"]'

solar deploy contracts/XC.sol '["7174756d00000000000000000000000000000000000000000000000000000000"]'

#
#echo '      ############################'
#echo '     #     1: INK               #'
#echo '    #     2: XCPlugin          #'
#echo '   #     3: XC                #'
#echo '  ############################'
#echo ' # (default:All) Enter [1~3]:'
#read aNum
#
#case ${aNum} in
#    '')
#    echo "###################################################"
#    echo "#                      START                      #"
#    echo "###################################################"
#    echo "### 1）Transfer 3 QTUM to the contract account. ###"
#    echo "Please enter a release contract account:"
#    read deployAccount
#    sh qtum-cli.sh sendtoaddress ${deployAccount} 3
#    echo "### 2）Preparation of INK contract bytecode. ###"
#
#    INKBytecode=`cat ../contracts/src/ink-token/INK.bin`
#    echo "INK Bytecode:${INKBytecode:0:16}..."
#
#    echo "### 3）Deployment of INK contract. ###"
#    sh qtum-cli.sh createcontract ${INKBytecode} 2500000  0.00000049 ${deployAccount}
#
#    echo "###################################################"
#    echo "#                       END                       #"
#    echo "###################################################"
#
#    echo "###################################################"
#    echo "#                      START                      #"
#    echo "###################################################"
#    echo "### 1）Transfer 3 QTUM to the contract account. ###"
#    echo "Please enter the platform name:"
#    read platform
#    echo "Please enter a release contract account:"
#    read deployAccount
#    sh qtum-cli.sh sendtoaddress ${deployAccount} 3
#    echo "### 2）Preparation of XCPlugin contract bytecode. ###"
#
#    XCPluginBytecode=`cat ../contracts/src/xc-plugin/XCPlugin.bin && docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest encode params -v string ${platform} --lenient`
#    echo "XCPlugin Bytecode:${XCPluginBytecode:0:16}..."
#
#    echo "### 3）Deployment of XCPlugin contract. ###"
#    sh qtum-cli.sh createcontract ${XCPluginBytecode} 2500000  0.00000049 ${deployAccount}
#
#    echo "###################################################"
#    echo "#                       END                       #"
#    echo "###################################################"
#
#    echo "###################################################"
#    echo "#                      START                      #"
#    echo "###################################################"
#    echo "### 1）Transfer 3 QTUM to the contract account. ###"
#    echo "Please enter the platform name:"
#    read platform
#    echo "Please enter a release contract account:"
#    read deployAccount
#    sh qtum-cli.sh sendtoaddress ${deployAccount} 3
#    echo "### 2）Preparation of XC contract bytecode. ###"
#
#    XCBytecode=`cat ../contracts/src/ink-xc/XC.bin && docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest encode params -v string ${platform} --lenient`
#    echo "XC Bytecode:${XCBytecode:0:16}..."
#
#    echo "### 3）Deployment of XC contract. ###"
#    sh qtum-cli.sh createcontract ${XCBytecode} 2500000  0.00000049 ${deployAccount}
#
#    echo "###################################################"
#    echo "#                       END                       #"
#    echo "###################################################"
#    ;;
#    1)
#    echo "###################################################"
#    echo "#                      START                      #"
#    echo "###################################################"
#    echo "### 1）Transfer 3 QTUM to the contract account. ###"
#    echo "Please enter a release contract account:"
#    read deployAccount
#    sh qtum-cli.sh sendtoaddress ${deployAccount} 3
#    echo "### 2）Preparation of INK contract bytecode. ###"
#
#    INKBytecode=`cat ../contracts/src/ink-token/INK.bin`
#    echo "INK Bytecode:${INKBytecode:0:16}..."
#
#    echo "### 3）Deployment of INK contract. ###"
#    sh qtum-cli.sh createcontract ${INKBytecode} 2500000  0.00000049 ${deployAccount}
#
#    echo "###################################################"
#    echo "#                       END                       #"
#    echo "###################################################"
#    ;;
#    2)
#    echo "###################################################"
#    echo "#                      START                      #"
#    echo "###################################################"
#    echo "### 1）Transfer 3 QTUM to the contract account. ###"
#    echo "Please enter the platform name:"
#    read platform
#    echo "Please enter a release contract account:"
#    read deployAccount
#    sh qtum-cli.sh sendtoaddress ${deployAccount} 3
#    echo "### 2）Preparation of XCPlugin contract bytecode. ###"
#
#    XCPluginBytecode=`cat ../contracts/src/xc-plugin/XCPlugin.bin && docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest encode params -v string ${platform} --lenient`
#    echo "XCPlugin Bytecode:${XCPluginBytecode:0:16}..."
#
#    echo "### 3）Deployment of XCPlugin contract. ###"
#    sh qtum-cli.sh createcontract ${XCPluginBytecode} 2500000  0.00000049 ${deployAccount}
#
#    echo "###################################################"
#    echo "#                       END                       #"
#    echo "###################################################"
#    ;;
#    3)
#    echo "###################################################"
#    echo "#                      START                      #"
#    echo "###################################################"
#    echo "### 1）Transfer 3 QTUM to the contract account. ###"
#    echo "Please enter the platform name:"
#    read platform
#    echo "Please enter a release contract account:"
#    read deployAccount
#    sh qtum-cli.sh sendtoaddress ${deployAccount} 3
#    echo "### 2）Preparation of XC contract bytecode. ###"
#
#    XCBytecode=`cat ../contracts/src/ink-xc/XC.bin && docker run --rm -v ${PWD}:/ethabi cryptominder/ethabi:latest encode params -v string ${platform} --lenient`
#    echo "XC Bytecode:${XCBytecode:0:16}..."
#
#    echo "### 3）Deployment of XC contract. ###"
#    sh qtum-cli.sh createcontract ${XCBytecode} 2500000  0.00000049 ${deployAccount}
#
#    echo "###################################################"
#    echo "#                       END                       #"
#    echo "###################################################"
#    ;;
#    *)
#    echo "###################################################"
#    echo "#                        ~^@^~                    #"
#    echo "###################################################"
#    ;;
#esac
