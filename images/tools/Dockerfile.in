# Copyright Greg Haskins All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
FROM _BASE_NS_/inkchain-baseimage:_BASE_TAG_
ENV INKCHAIN_CFG_PATH /etc/inklabsfoundation/inkchain
VOLUME /etc/inklabsfoundation/inkchain
ADD  payload/sampleconfig.tar.bz2 $INKCHAIN_CFG_PATH
RUN  apt-get update \
         && apt-get install -y vim tree jq \
         && apt-get install -y unzip \
         && rm -rf /var/cache/apt
COPY payload/cryptogen /usr/local/bin
COPY payload/configtxgen /usr/local/bin
COPY payload/configtxlator /usr/local/bin
COPY payload/peer /usr/local/bin
