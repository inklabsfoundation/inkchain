# Copyright Greg Haskins All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
FROM _NS_/inkchain-buildenv:_TAG_

# inkchain configuration locations
ENV INKCHAIN_CFG_PATH /etc/inklabsfoundation/inkchain

# create needed directories
RUN mkdir -p \
  $INKCHAIN_CFG_PATH \
  /var/inkchain/production

# inkchain configuration files
ADD payload/sampleconfig.tar.bz2 $INKCHAIN_CFG_PATH

# inkchain binaries
COPY payload/orderer /usr/local/bin
COPY payload/peer /usr/local/bin

# softhsm2
COPY payload/install-softhsm2.sh /tmp
RUN bash /tmp/install-softhsm2.sh && rm -f install-softhsm2.sh

# typically, this is mapped to a developer's dev environment
WORKDIR /opt/gopath/src/github.com/inklabsfoundation/inkchain
