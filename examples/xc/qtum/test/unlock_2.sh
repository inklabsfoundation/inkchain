#! /bin/bash

source p_init.sh

# voter(bytes32,address,address,uint256,bytes32   ,bytes32,bytes32,uint8)
fromPlatform="494e4b0000000000000000000000000000000000000000000000000000000000"
fromAccount="4230a12f5b0693dd88bb35c79d7e56a68614b199"
toAccount="db39e7c2a4009e69804c1d9737db29a133c80e5d"
value=100000
txid="2ba06c917766ac7fe5711da6642593e1941eff10c35ebe0cae700cb4beda39d0"
#sign="d0dbd5c9f1fc491bfdf4b2c9c751608168388f1c15a90767ad88bc63504a5e8a7b347c21c8d537d18949d667d8ae6d2e8d7c61973342c13048722dccc692048001"
#d0dbd5c9f1fc491bfdf4b2c9c751608168388f1c15a90767ad88bc63504a5e8a
#7b347c21c8d537d18949d667d8ae6d2e8d7c61973342c13048722dccc6920480
#01
r="d0dbd5c9f1fc491bfdf4b2c9c751608168388f1c15a90767ad88bc63504a5e8a"
s="7b347c21c8d537d18949d667d8ae6d2e8d7c61973342c13048722dccc6920480"
v=28

# balanceOf
qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$toAccount'"]'`

qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol voteProposal '["'$fromPlatform'","'$fromAccount'","'$toAccount'",'$value',"'$txid'","'$r'","'$s'",'$v']'` 0 6000000 0.000001 $XCPluginOwner
sleep 180
#2488a4d8a028c56cc1bab54dfa2b186897f765757877ab638214a8719e06da777fbc930b59d48817eb3230f577b1a66076597a3d8c1370f53cd692d4a6c63d8800
r1="2488a4d8a028c56cc1bab54dfa2b186897f765757877ab638214a8719e06da77"
s1="7fbc930b59d48817eb3230f577b1a66076597a3d8c1370f53cd692d4a6c63d88"
v1=27
qcli sendtocontract $XCPluginHexAddress `solar encode contracts/XCPlugin.sol voteProposal '["'$fromPlatform'","'$fromAccount'","'$toAccount'",'$value',"'$txid'","'$r1'","'$s1'",'$v1']'` 0 6000000 0.000001 $XCPluginOwner
sleep 180

qcli sendtocontract $XCHexAddress `solar encode contracts/XC.sol unlock '["'$txid'","'$fromPlatform'","'$fromAccount'","'$toAccount'",'$value']'` 0 6000000 0.000001 $XCPluginOwner
sleep 180

# balanceOf
qcli callcontract $INKHexAddress `solar encode contracts/INK.sol balanceOf '["'$toAccount'"]'`
