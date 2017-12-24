#!/bin/bash
#
# Copyright Ziggurat Corp. 2017 All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

export VERSION=0.10.3
export BASEIMAGE_RELEASE=0.3.1

#Set MARCH variable i.e x86_64
MARCH=`uname -m`

InkchainDockerPull() {
  local INKCHAIN_TAG=$1
  for IMAGES in peer orderer tools ccenv javaenv couchdb kafka zookeeper; do
      echo "==> INKCHAIN IMAGE: $IMAGES"
      echo
      docker pull inklabsfoundation/inkchain-$IMAGES:$INKCHAIN_TAG
      docker tag inklabsfoundation/inkchain-$IMAGES:$INKCHAIN_TAG inklabsfoundation/inkchain-$IMAGES
  done
}

CaDockerPull() {
      local CA_TAG=$1
      echo "==> INKCHAIN CA IMAGE"
      echo
      docker pull inklabsfoundation/inkchain-ca:$CA_TAG
      docker tag inklabsfoundation/inkchain-ca:$CA_TAG inklabsfoundation/inkchain-ca
}

BaseImagesPull() {
      docker pull inklabsfoundation/inkchain-baseimage:$MARCH-$BASEIMAGE_RELEASE
      docker pull inklabsfoundation/inkchain-baseos:$MARCH-$BASEIMAGE_RELEASE

}

: ${CA_TAG:="$MARCH-$VERSION"}
: ${INKCHAIN_TAG:="$MARCH-$VERSION"}

echo "===> Pulling Base Images"
BaseImagesPull

echo "===> Pulling inkchain Images"
InkchainDockerPull ${INKCHAIN_TAG}

echo "===> Pulling inkchain ca Image"
CaDockerPull ${CA_TAG}
echo
echo "===> List out inkchain docker images"
docker images | grep inklabsfoundation*

