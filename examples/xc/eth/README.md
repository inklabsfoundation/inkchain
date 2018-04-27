### Cross Chain

#### 1） ETH Environment Setup.

##### contracts
```
$ ln -s ../docker-qtum/contracts contracts
```

##### init eth network
```
$ sh init.sh
```

##### start eth network
```
$ sh start.sh
```

##### stop eth network
```
$ sh stop.sh
```

##### remove eth network
```
$ sh remove.sh
```

##### eth console
```
$ sh console.sh
```

##### new Account
```
$ sh console.sh
> personal.newAccount("123456")
> miner.start();admin.sleepBlocks(20);miner.stop();
> web3.fromWei(eth.getBalance(eth.accounts[0]),'ether')
> 
```

##### unlock Account
````
$ personal.unlockAccount(eth.accounts[0],"ma123456")
````

