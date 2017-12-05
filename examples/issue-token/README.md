## Token-issue Flow

-----------------------------

This example provisions a sample Inkchain network consisting of 1 CA, 1 orderer, 1 peer, and a cli client.

### Download Images

Prerequisites:

- [Git client](https://git-scm.com/downloads)

- [Docker](https://www.docker.com/products/overview) v1.12 or higher

- [Docker Compose](https://docs.docker.com/compose/overview/) v1.8 or higher

Download images from `https://hub.docker.com/u/inkchain/`, here we use images version 0.9.1.

### Start Inkchain network

We provide the basic yaml file to start the network, and we can use following command to launch.

```bash
$ docker-compose up (-d)
```

wait,10 seconds or so, without error log output. execute a `docker ps` to view your active
containers. you shouled see an output identical to the following:

```bash
CONTAINER ID        IMAGE                                                  COMMAND                  CREATED             STATUS              PORTS                                            NAMES
1c3c57cb4fda        inklabsfoundation/inkchain-tools                       "bash -c 'while tr..."   About an hour ago   Up About an hour                                                     cli
cc5376420976        inklabsfoundation/inkchain-peer                        "peer node start"        About an hour ago   Up About an hour    0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp   peer0.org1.example.com
06a7356b3b93        inklabsfoundation/inkchain-orderer                     "orderer"                About an hour ago   Up About an hour    0.0.0.0:7050->7050/tcp                           orderer.example.com
094d74de2ee6        inklabsfoundation/inkchain-ca                          "sh -c 'inkchain-c..."   About an hour ago   Up About an hour    0.0.0.0:7054->7054/tcp                           ca_peerOrg1

```

### Execute test sctipts

At a new terminal, we will go into the `cli` container, and execute test command:

```bash
$ docker exec -it cli bash

$ root@:/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer# bash ./scripts/initialization.sh
```

You should see result like the following if the initialization is successful.

```bash

...
=====================5.Instantiate chaincode token, this will take a while, pls waiting...===
2017-11-17 03:08:00.973 UTC [msp] GetLocalMSP -> DEBU 001 Returning existing local MSP
2017-11-17 03:08:00.973 UTC [msp] GetDefaultSigningIdentity -> DEBU 002 Obtaining default signing identity
2017-11-17 03:08:00.974 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 003 Using default escc
2017-11-17 03:08:00.974 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 004 Using default vscc
2017-11-17 03:08:00.975 UTC [msp/identity] Sign -> DEBU 005 Sign: plaintext: 0A95070A6708031A0C0890A5B9D00510...314D53500A04657363630A0476736363
2017-11-17 03:08:00.975 UTC [msp/identity] Sign -> DEBU 006 Sign: digest: C4C272C304B7F94526BB9E3B1EB548CEF22536361FC7B9CA78B1D2B82F2E4F56
2017-11-17 03:08:10.722 UTC [msp/identity] Sign -> DEBU 007 Sign: plaintext: 0A95070A6708031A0C0890A5B9D00510...1A27A866A689AB9E01AA12021DF03C45
2017-11-17 03:08:10.722 UTC [msp/identity] Sign -> DEBU 008 Sign: digest: 1CAF6F8DAAE1D1A294893BB05209BFE9F6A49CA9E5791C5E59C984558E4A5BD4
2017-11-17 03:08:10.724 UTC [main] main -> INFO 009 Exiting.....
=========== Chaincode token Instantiation on peer0.org1 on channel 'mychannel' is successful ==========
Instantiate spent 10 secs

=====================All GOOD, MVE Initialization completed =====================

```

Test issue-token, transfer and query using following script:

```bash
$ root@:/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer# bash ./scripts/test_token.sh

```

You should see result like the following if the test-issue-token successful.

```bash
...
=====================9.Query transfer result of To account=====================
Attempting to  query account B's balance on peer
2017-11-17 03:08:31.017 UTC [msp] GetLocalMSP -> DEBU 001 Returning existing local MSP
2017-11-17 03:08:31.017 UTC [msp] GetDefaultSigningIdentity -> DEBU 002 Obtaining default signing identity
2017-11-17 03:08:31.017 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 003 Using default escc
2017-11-17 03:08:31.017 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 004 Using default vscc
2017-11-17 03:08:31.017 UTC [msp/identity] Sign -> DEBU 005 Sign: plaintext: 0A95070A6708031A0B08AFA5B9D00510...663634633739610A074343546F6B656E
2017-11-17 03:08:31.018 UTC [msp/identity] Sign -> DEBU 006 Sign: digest: C4AB537F895CE1D1375C09F2E83484EF9028F5342A25710294FAF8F66F532B5B
2017-11-17 03:08:31.024 UTC [main] main -> INFO 007 Exiting.....
Query Result: {"CCToken":"10"}

=====================All GOOD, MVE Test completed =====================

```

