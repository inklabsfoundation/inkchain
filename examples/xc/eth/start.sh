#! /bin/bash

nohup geth --datadir .eth_node --networkid 10086 --rpc --rpcapi "eth,miner,personal" --nodiscover &