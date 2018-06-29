#! /bin/bash
# db,eth,net,web3,personal,admin,miner
nohup geth --datadir .eth_node --networkid 10086 --rpc --rpcapi "db,eth,net,web3,personal,admin,miner" --nodiscover &

