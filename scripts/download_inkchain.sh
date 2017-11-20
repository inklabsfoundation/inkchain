#!/bin/bash

if [ -f ./colour.sh ]; then
 source ./colour.sh
elif [ -f scripts/colour.sh ]; then
 source scripts/colour.sh
else
 alias echo_r="echo"
 alias echo_g="echo"
 alias echo_b="echo"
fi

echo_b "===Download official images from https://hub.docker.com/u/inkchain/"

# pull inkchain images
ARCH=x86_64
BASEIMAGE_RELEASE=0.3.1
IMG_TAG_LATEST=0.10

echo_b "===Pulling inkchain images... with tag = ${IMG_TAG_LATEST}"
docker pull inkchain/inkchain-peer:$ARCH-$IMG_TAG_LATEST
docker pull inkchain/inkchain-tools:$ARCH-$IMG_TAG_LATEST
docker pull inkchain/inkchain-orderer:$ARCH-$IMG_TAG_LATEST
docker pull inkchain/inkchain-ccenv:$ARCH-$IMG_TAG_LATEST

docker pull inkchain/inkchain-baseimage:$ARCH-$BASEIMAGE_RELEASE
docker pull inkchain/inkchain-baseos:$ARCH-$BASEIMAGE_RELEASE
docker pull inkchain/inkchain-ca:x86_64-1.0.5

docker pull inkchain/inkchain-couchdb:$ARCH-$IMG_TAG_LATEST
docker pull inkchain/inkchain-kafka:$ARCH-$IMG_TAG_LATEST
docker pull inkchain/inkchain-zookeeper:$ARCH-$IMG_TAG_LATEST
docker pull inkchain/inkchain-javaenv:$ARCH-$IMG_TAG_LATEST

echo_b "===Re-tagging images to *latest* tag"
docker tag inkchain/inkchain-peer:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-peer
docker tag inkchain/inkchain-tools:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-tools
docker tag inkchain/inkchain-orderer:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-orderer
docker tag inkchain/inkchain-zookeeper:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-zookeeper
docker tag inkchain/inkchain-kafka:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-kafka
docker tag inkchain/inkchain-couchdb:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-couchdb
docker tag inkchain/inkchain-javaenv:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-javaenv
docker tag inkchain/inkchain-ccenv:$ARCH-$IMG_TAG_LATEST inkchain/inkchain-ccenv
docker tag inkchain/inkchain-ca:x86_64-1.0.5 inkchain/inkchain-ca

echo_b "Done, now can startup the network using docker-compose..."
