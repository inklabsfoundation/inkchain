### Cross Chain

#### 1） Qtum Environment Setup.

##### 1.1）Installing Docker

> https://docs.docker.com/engine/installation/
>
##### 1.2）Pull the Qtum Docker Image

```
$ cd docker-qtum
$ sh download_images.sh

```
##### 1.3）Installing solar

```
$ go get -u github.com/qtumproject/solar/cli/solar
```
##### 1.4）Using the Qtum Docker



Starting up
```
$ cd docker-qtum
$ sh start.sh
```
Stopping 
```
$ cd docker-qtum
$ sh stop.sh
```
Remove Containers 
```
$ cd docker-qtum
$ sh remove.sh
```
login
```
$ cd docker-qtum
$ sh login.sh
```

#### 2）Script Commands.

##### 2.1）Help
```
$ cd docker-qtum
$ sh qtum-cli.sh help
```
##### 2.2）Common commands

Get New Address
```
$ cd docker-qtum
$ sh qtum-cli.sh getnewaddress
```

Get the account address
```
$ cd docker-qtum
$ sh qtum-cli.sh getaccountaddress ${account} // lilei
```

Get Hex Address
```
$ cd docker-qtum
$ sh qtum-cli.sh gethexaddress ${account_address} // qHhn8KTpfootyYy6Eo51dfR8kXn8s1Wv8i
```

fromhexaddress
```
$ cd docker-qtum
$ sh qtum-cli.sh fromhexaddress ${account_hexaddress} // 019680a2de9e2d8753dfb86437ed52865f0a135c
```

Check accounts balance
```
$ cd docker-qtum
$ sh qtum-cli.sh listaccounts 
```

Send to address
```
$ cd docker-qtum
$ sh qtum-cli.sh sendtoaddress ${toAccount} ${amount} // qHhn8KTpfootyYy6Eo51dfR8kXn8s1Wv8i 1000
```

List unspent
```
$ cd docker-qtum
$ sh qtum-cli.sh listunspent 1 9999999 [\"${account_address}\"] // [\"qHhn8KTpfootyYy6Eo51dfR8kXn8s1Wv8i\\"] 
```

List contracts
```
$ cd docker-qtum
$ sh qtum-cli.sh listcontracts
```

Get Account Info
```
$ cd docker-qtum
$ sh qtum-cli.sh  getaccountinfo ${account_address} // contract address
```

Create contract 
```
$ cd docker-qtum
$ sh qtum-cli.sh createcontract ${bytecode} ${gasLimit} ${gasPrice} ${senderaddress} 
```

Call contract
```
$ cd docker-qtum
$ sh qtum-cli.sh  callcontract ${contractaddress} ${bytecode}
```

Send to contract
```
$ cd docker-qtum
$ sh qtum-cli.sh  sendtocontract ${contractaddress} ${bytecode}
```

Estimate fee
```
$ cd docker-qtum
$ sh qtum-cli.sh  estimatefee ${contractaddress} ${bytecode}
```

#### 3）Deploy The Smart Contracts.

##### 3.1）deploy
```
$ cd docker-qtum/contracts
$ sh solar_deploy.sh
```
##### 3.2）deploy status
```
$ cd docker-qtum
$ sh solar_deploy_status.sh
```

#### 4) method signature

login
```
$ sh login.sh
$ qcli generate 600
```

#### 4.1) INK
##### 4.1.1) allowance(address,address)

```
$ solar encode contracts/INK.sol allowance '["18ac89acd86ea466cfb964dd9847d72e427c7886","d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'

dd62ed3e00000000000000000000000018ac89acd86ea466cfb964dd9847d72e427c7886000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```

##### 4.1.2) approve(address,uint256)

```
$ solar encode contracts/INK.sol approve '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'

095ea7b3000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```

##### 4.1.3) balanceOf(address)

```
$ solar encode contracts/INK.sol balanceOf '["18ac89acd86ea466cfb964dd9847d72e427c7886"]'

70a0823100000000000000000000000018ac89acd86ea466cfb964dd9847d72e427c7886

$ solar encode contracts/INK.sol balanceOf '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'

70a08231000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```

##### 4.1.4) decimals()

```
$ solar encode contracts/INK.sol decimals
313ce567
```

##### 4.1.5) totalSupply()
```
$ solar encode contracts/INK.sol totalSupply
18160ddd
```

##### 4.1.6) transfer(address,uint256)
```
$ solar encode contracts/INK.sol transfer '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
a9059cbb000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```

##### 4.1.7) transferFrom(address,address,uint256)

```
$ solar encode contracts/INK.sol transferFrom '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
23b872dd000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```

##### 4.1.8) name()
```
$ solar encode contracts/INK.sol name 
06fdde03
```

##### 4.1.9) symbol()

```
$ solar encode contracts/INK.sol symbol 
95d89b41
```

##### 4.1.10) decreaseApproval(address,uint256)
```
$ solar encode contracts/INK.sol decreaseApproval  '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
66188463000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```

##### 4.1.11) increaseApproval(address,uint256)
```
$ solar encode contracts/INK.sol increaseApproval  '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
d73dd623000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```

##### 4.1.12) approveAndCall(address,uint256,bytes)
```
$ solar encode contracts/INK.sol approveAndCall  '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000,"41"]'
cae9ca51000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000014100000000000000000000000000000000000000000000000000000000000000
```
#### 4.2) XCPlugin

##### 4.2.1) addCaller(address)
```
$ solar encode contracts/XCPlugin.sol addCaller  '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
747293fb000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```

##### 4.2.2) addPlatform(bytes32)
```
$ solar encode contracts/XCPlugin.sol addPlatform  '["7174756d00000000000000000000000000000000000000000000000000000000"]'
df8d66e47174756d00000000000000000000000000000000000000000000000000000000
```
##### 4.2.3) addPublicKey(bytes32,address)
```
$ solar encode contracts/XCPlugin.sol addPublicKey  '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
1de1ee7b4100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```
##### 4.2.4) callers(uint256)
```
$ solar encode contracts/XCPlugin.sol callers  '[0]'
aa03fa3d0000000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.5) countOfPublicKey(bytes32)
```
$ solar encode contracts/XCPlugin.sol countOfPublicKey  '["4100000000000000000000000000000000000000000000000000000000000000"]'
c79c91054100000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.6) deleteCaller(address)
```
$ solar encode contracts/XCPlugin.sol deleteCaller  '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
258c0f1e000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```
##### 4.2.7) deletePlatform(bytes32)
```
$ solar encode contracts/XCPlugin.sol deletePlatform  '["4100000000000000000000000000000000000000000000000000000000000000"]'
8147bcf74100000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.8) deleteProposal(bytes32,bytes32)
```
$ solar encode contracts/XCPlugin.sol deleteProposal  '["4100000000000000000000000000000000000000000000000000000000000000","4100000000000000000000000000000000000000000000000000000000000000"]'
b56e0e3241000000000000000000000000000000000000000000000000000000000000004100000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.9) deletePublicKey(bytes32,address)
```
$ solar encode contracts/XCPlugin.sol deletePublicKey  '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
fa31b3714100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```
##### 4.2.10) existPlatform(bytes32)
```
$ solar encode contracts/XCPlugin.sol existPlatform  '["4100000000000000000000000000000000000000000000000000000000000000"]'
5de047d34100000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.11) getAdmin()
```
$ solar encode contracts/XCPlugin.sol getAdmin 
6e9960c3
```
##### 4.2.12) getPlatform(bytes32)
```
$ solar encode contracts/XCPlugin.sol getPlatform  '["4100000000000000000000000000000000000000000000000000000000000000"]'
569cd4cb4100000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.13) getPlatformName()
```
$ solar encode contracts/XCPlugin.sol getPlatformName 
db3b73db
```
##### 4.2.14) getStatus()
```
$ solar encode contracts/XCPlugin.sol getStatus
4e69d560
```
##### 4.2.15) getWeight(bytes32)
```
$ solar encode contracts/XCPlugin.sol getWeight  '["4100000000000000000000000000000000000000000000000000000000000000"]'
66435abf4100000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.16) setAdmin(address)
```
$ solar encode contracts/XCPlugin.sol setAdmin  '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
704b6c02000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```
##### 4.2.17) setPlatformName(bytes32)
```
$ solar encode contracts/XCPlugin.sol setPlatformName  '["4100000000000000000000000000000000000000000000000000000000000000"]'
8d8992d34100000000000000000000000000000000000000000000000000000000000000
```
##### 4.2.18) setWeight(bytes32,uint256)
```
$ solar encode contracts/XCPlugin.sol setWeight  '["4100000000000000000000000000000000000000000000000000000000000000",3]'
71c4275a41000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003
```
##### 4.2.19) start()
```
$ solar encode contracts/XCPlugin.sol start
be9a6555
```
##### 4.2.20) stop()
```
$ solar encode contracts/XCPlugin.sol stop
07da68f5
```
##### 4.2.21) verifyProposal(bytes32,address,address,uint256,bytes32)
```
$ solar encode contracts/XCPlugin.sol verifyProposal  '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000,"4200000000000000000000000000000000000000000000000000000000000000"]'
623aba4d4100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e84200000000000000000000000000000000000000000000000000000000000000
```

##### 4.2.22) voter(bytes32,address,address,uint256,bytes32,bytes32,bytes32,uint8)
```
$ solar encode contracts/XCPlugin.sol voter  '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000,"4200000000000000000000000000000000000000000000000000000000000000","4100000000000000000000000000000000000000000000000000000000000000","4100000000000000000000000000000000000000000000000000000000000000",27]'
ffe6bf064100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8420000000000000000000000000000000000000000000000000000000000000041000000000000000000000000000000000000000000000000000000000000004100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001b
```

#### 4.3) XC

##### 4.3.1) getAdmin()
```
$ solar encode contracts/XC.sol getAdmin
6e9960c3
```
##### 4.3.2) getINK()
```
$ solar encode contracts/XC.sol getINK
e5177ee4
```
##### 4.3.3) getPlatformName()
```
$ solar encode contracts/XC.sol getPlatformName
db3b73db
```
##### 4.3.4) getStatus()
```
$ solar encode contracts/XC.sol getStatus
4e69d560
```
##### 4.3.5) getXCPlugin()
```
$ solar encode contracts/XC.sol getXCPlugin
4567a6a7
```
##### 4.3.6) lock(bytes32,address,uint256)
```
$ solar encode contracts/XC.sol lock '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
0899146b4100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```
##### 4.3.7) lockAdmin(bytes32,address,uint256)
```
$ solar encode contracts/XC.sol lockAdmin '["4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
30b36e6a4100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```
##### 4.3.8) lockBalance()
```
$ solar encode contracts/XC.sol lockBalance
e228ecb4
```
##### 4.3.9) setAdmin(address)
```
$ solar encode contracts/XC.sol setAdmin '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
704b6c02000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```
##### 4.3.10) setINK(address)
```
$ solar encode contracts/XC.sol setINK '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
b8f23f52000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```
##### 4.3.11) setPlatformName(bytes32)
```
$ solar encode contracts/XC.sol setPlatformName '["4100000000000000000000000000000000000000000000000000000000000000"]'
8d8992d34100000000000000000000000000000000000000000000000000000000000000
```
##### 4.3.12) setXCPlugin(address)
```
$ solar encode contracts/XC.sol setXCPlugin '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e"]'
8aa5440e000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e
```
##### 4.3.13) start()
```
$ solar encode contracts/XC.sol start
be9a6555
```
##### 4.3.14) stop()
```
$ solar encode contracts/XC.sol stop
07da68f5
```
##### 4.3.15) unlock(bytes32,bytes32,address,address,uint256)
```
$ solar encode contracts/XC.sol unlock '["4100000000000000000000000000000000000000000000000000000000000000","4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
1a11d6d341000000000000000000000000000000000000000000000000000000000000004100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```
##### 4.3.16) unlockAdmin(bytes32,bytes32,address,address,uint256)
```
$ solar encode contracts/XC.sol unlockAdmin '["4100000000000000000000000000000000000000000000000000000000000000","4100000000000000000000000000000000000000000000000000000000000000","d6b39eb631df8ee60e46a576231ccf1fcd204a5e","d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
85ea4fd741000000000000000000000000000000000000000000000000000000000000004100000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```
##### 4.3.17) withdrawal(address,uint256)
```
$ solar encode contracts/XC.sol withdrawal '["d6b39eb631df8ee60e46a576231ccf1fcd204a5e",1000]'
5a6b26ba000000000000000000000000d6b39eb631df8ee60e46a576231ccf1fcd204a5e00000000000000000000000000000000000000000000000000000000000003e8
```

#### 5）Testing The Smart Contracts.

##### 5.1）prepare account 
```
$ qcli getaccountaddress ma
$ qcli getaccountaddress sender
$ qcli listaccounts

{
  "": 2360000.00000000,
  "ma": 0.00000000,
  "sender": 0.00000000
}

```
##### 5.2）Contract
```
$ sh deploy.sh
... ...
$ sh deploy_status.sh

✅  contracts/XCPlugin.sol
      txid: d70dd064f8144cf6df82018fcf812a009304eed69a0953f4f13cd9a5b560e32d
   address: fdfe183ff6d196a2871d3b856629222f1dd1264b
 confirmed: true
     owner: qf2FFX6c5qohyqPvbUoyy8Xo4Q6QrjsaKX


✅  contracts/INK.sol
        txid: df3cda23e2b31f154bb31777b45ddbddff1eabbad636f9dea255528952261615
     address: 223a30419a3dc234ec659e390d7bc65e70dcea0b
   confirmed: true
       owner: qP2RGcwGunqiZvSUidSQyEW64G1XuyGUyM

✅  contracts/XC.sol
        txid: 74a98da3dc957d90bce9040d66ab0175f65f913acb07df7f857d3a2e716e516f
     address: 54acd8f45dad524584003ee4d4f4d4f8b00ef4c1
   confirmed: true
       owner: qMqd6usqYmUXpTB3FShUiKP59TfMrXFQJf

```
##### 5.3）prepare Contract address

```
$ vim p_init.sh
$ sh p_read.sh
```

#### 6） examples

##### 6.1）callcontract INK name
```
$ sh test/INK_name.sh
```
##### 6.2）callcontract INK symbol
```
$ sh test/INK_symbol.sh
```
##### 6.2）callcontract INK balanceOf(address) 
```
$ sh test/INK_balanceOf.sh
```
##### 6.3）sendtocontract INK transfer(address,uint)
```
$ sh test/INK_transfer.sh
```
... ...