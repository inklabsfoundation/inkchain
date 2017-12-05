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
IMG_TAG_LATEST=0.10.2

echo_b "===Pulling inkchain images... with tag = ${IMG_TAG_LATEST}"
docker pull inklabsfoundation/inkchain-peer:$ARCH-$IMG_TAG_LATEST
docker pull inklabsfoundation/inkchain-tools:$ARCH-$IMG_TAG_LATEST
docker pull inklabsfoundation/inkchain-orderer:$ARCH-$IMG_TAG_LATEST
docker pull inklabsfoundation/inkchain-ccenv:$ARCH-$IMG_TAG_LATEST

docker pull inklabsfoundation/inkchain-baseimage:$ARCH-$BASEIMAGE_RELEASE
docker pull inklabsfoundation/inkchain-baseos:$ARCH-$BASEIMAGE_RELEASE
docker pull inklabsfoundation/inkchain-ca:$ARCH-$IMG_TAG_LATEST

docker pull inklabsfoundation/inkchain-couchdb:$ARCH-$IMG_TAG_LATEST
docker pull inklabsfoundation/inkchain-kafka:$ARCH-$IMG_TAG_LATEST
docker pull inklabsfoundation/inkchain-zookeeper:$ARCH-$IMG_TAG_LATEST
docker pull inklabsfoundation/inkchain-javaenv:$ARCH-$IMG_TAG_LATEST

echo_b "===Re-tagging images to *latest* tag"
docker tag inklabsfoundation/inkchain-peer:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-peer
docker tag inklabsfoundation/inkchain-tools:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-tools
docker tag inklabsfoundation/inkchain-orderer:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-orderer
docker tag inklabsfoundation/inkchain-zookeeper:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-zookeeper
docker tag inklabsfoundation/inkchain-kafka:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-kafka
docker tag inklabsfoundation/inkchain-couchdb:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-couchdb
docker tag inklabsfoundation/inkchain-javaenv:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-javaenv
docker tag inklabsfoundation/inkchain-ccenv:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-ccenv
docker tag inklabsfoundation/inkchain-ca:$ARCH-$IMG_TAG_LATEST inklabsfoundation/inkchain-ca

echo_b "Done, now can startup the network using docker-compose..."
