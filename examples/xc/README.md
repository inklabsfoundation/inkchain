### Cross Chain

#### 1） Qtum Environment Setup.

##### 1.1）Installing Docker

> https://docs.docker.com/engine/installation/
>

##### 1.2）Pull the Qtum Docker Image

Download qtum,ethabi,solc Images.
```
$ cd docker-qtum
$ sh download_images.sh
```
##### 1.3）Using the Qtum Docker

Starting up
```
$ cd docker-qtum
$ sh start.sh
```
 
Stop Running 
```
$ cd docker-qtum
$ sh stop.sh
```
Remove Containers 
```
$ cd docker-qtum
$ sh remove.sh
```

#### 2）Script Commands.

##### 2.1）Help
```
$ cd docker-qtum
$ sh qtum-cli.sh help
```

##### 2.2）Common commands

Get the account address
```
$ cd docker-qtum
$ sh qtum-cli.sh getaccountaddress greeter_owner
```

Check accounts balance
```
$ cd docker-qtum
$ sh qtum-cli.sh listaccounts
```

Send to address
```
$ cd docker-qtum
$ sh qtum-cli.sh sendtoaddress greeter_owner 3
```

List unspent
```
$ cd docker-qtum
$ sh qtum-cli.sh listunspent 1 9999999 [\"greeter_owner\"]
```

List contracts
```
$ cd docker-qtum
$ sh qtum-cli.sh listcontracts
```

Get Account Info
```
$ cd docker-qtum
$ sh qtum-cli.sh  getaccountinfo contractaddress
```

Create contract 
```
$ cd docker-qtum
$ sh qtum-cli.sh createcontract bytecode gasLimit gasPrice senderaddress
```

Estimate smart fee
```
$ cd docker-qtum
$ sh qtum-cli.sh estimatesmartfee bytecode
```

Estimate smart priority
```
$ cd docker-qtum
$ sh qtum-cli.sh estimatesmartpriority bytecode gasLimit gasPrice senderaddress
```

Call contract
```
$ cd docker-qtum
$ sh qtum-cli.sh  callcontract "contractaddress" "data"
```

Send to contract
```
$ cd docker-qtum
$ sh qtum-cli.sh  sendtocontract "contractaddress" "data"
```

Estimate fee
```
$ cd docker-qtum
$ sh qtum-cli.sh  estimatefee "contractaddress" "data"
```

Estimate priority
```
$ cd docker-qtum
$ sh qtum-cli.sh  estimatepriority "contractaddress" "data"
```

#### 3）Compile The Smart Contracts.

```
$ cd contracts
$ sh compile.sh
$ ls -al src
$ *.abi、*.bin、*.signatures
```

#### 4）Deploy The Smart Contracts.

```
$ cd docker-qtum
$ sh deploy.sh
```
#### 5）Using The Smart Contracts.

Contract address
```
{
    INK:"5563547cf778cd54db7b3e31cffbf82f39d6cf7e"
    XCPlugin:"c77c7799c2f008f381dedef20d88878a1c02daa9"
    XC:"4fb66da2904e84a5726f0ec549aafcc154be0dbf"
}
```

callcontract
```
$ cd docker-qtum
$ sh qtum-cli.sh callcontract 5563547cf778cd54db7b3e31cffbf82f39d6cf7e 06fdde03
```

sendtocontract
```
$ cd docker-qtum
$ sh qtum-cli.sh sendtocontract 5563547cf778cd54db7b3e31cffbf82f39d6cf7e 06fdde03
```

#### 6）Ethabi util

```
$ cd docker-qtum
$ sh ethabi-util.sh
```




