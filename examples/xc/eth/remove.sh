#! /bin/bash

kill -9 $(lsof -i:8545 | awk '{print $2}')
rm -rvf solar.development.json
rm -rvf .eth_node
rm -rvf nohup.out