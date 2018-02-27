#! /bin/bash

docker run -d --rm --name qtumd_node1 --network=qtum_network \
    -v ${PWD}/node1_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node1_data:/data cryptominder/qtum:latest \
    qtumd

docker run -d --rm --name qtumd_node2 --network=qtum_network \
    -v ${PWD}/node2_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node2_data:/data cryptominder/qtum:latest \
    qtumd

docker run -d --rm --name qtumd_node3 --network=qtum_network \
    -v ${PWD}/node3_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node3_data:/data cryptominder/qtum:latest \
    qtumd
