#! /bin/bash

docker run -i --network container:qtumd_node1 \
    -v ${PWD}/node1_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node1_data:/data \
    cryptominder/qtum:latest \
    qtum-cli $@
