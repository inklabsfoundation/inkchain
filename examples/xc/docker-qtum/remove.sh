#! /bin/bash

#docker stop qtumd_node1 qtumd_node2 qtumd_node3
#rm -rf ./node1_data
#rm -rf ./node2_data
#rm -rf ./node3_data

docker stop qtumd_node
rm -rf ./.qtum
rm solar.development.json
sleep 5
docker rm `docker ps -aq`
