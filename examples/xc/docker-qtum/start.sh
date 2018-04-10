#! /bin/bash

# -e "QTUM_NETWORK=testnet" \
# -e "QTUM_NETWORK=mainnet" \

docker run -d --rm --name qtumd_node \
    -v ${PWD}:/dapp \
    -p 9899:9899 \
    -p 9888:9888 \
    -p 3889:3889 \
    -p 3888:3888 \
    -p 13888:13888 \
    hayeah/qtumportal:latest
