#! /bin/bash
#
# Copyright INK Corp. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0
#

# At a new vm we shouled install following:
# docker
# docker-compose
# pyhont2 python3
# go
# grpc
# java

#if [ xroot != x$(whoami) ]
#then
#   echo "You must run as root (Hint: sudo su)"
#   exit
#fi

username=`whoami`
USERGROUP=`id -gn`

# update system
sudo apt-get update

# Install necessary pkg

apt-get install -y build-essential git make curl unzip g++ libtool

echo "Install Java begin"

# Install Java
# ----------------------------------------------------------------
apt-get install -y openjdk-8-jdk maven

wget https://services.gradle.org/distributions/gradle-2.12-bin.zip -P /tmp --quiet
unzip -q /tmp/gradle-2.12-bin.zip -d /opt && rm /tmp/gradle-2.12-bin.zip
ln -s /opt/gradle-2.12/bin/gradle /usr/bin

echo "Install java end"


# Install nodejs
echo "Install nodejs begin"

NODE_VER=6.9.5
NODE_URL=https://nodejs.org/dist/v$NODE_VER/node-v$NODE_VER-linux-x64.tar.gz

curl -sL $NODE_URL | (cd /usr/local && tar --strip-components 1 -xz )

echo "Install nodejs end"

# Install docker
echo "Install latest release docker for ubuntu-16.04"
# Prep apt-get for docker install
sudo apt-get install -y apt-transport-https ca-certificates
sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D


# add ubuntu-16.04 repo
sudo sh -c "echo deb https://apt.dockerproject.org/repo ubuntu-xenial main > /etc/apt/sources.list.d/docker.list"

# update system
sudo apt-get update -qq

# Install docker
apt-get install -y linux-image-extra-$(uname -r) apparmor docker-engine

# Configure docker ubuntu-16.04
DOCKER_STORAGE_BACKEND_STRING=auf

# restart and add user
usermod -aG docker ${username}

# config docker accelerate
sudo mkdir -p /etc/docker

# use your own docker accelerator
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["http://xxxxxx.m.daocloud.io"]
}
EOF

sudo systemctl daemon-reload
sudo systemctl restart docker

echo "Install docker end"

# some necessary lib
sudo apt-get install python-dev python-pip libssl-dev libffi-dev curl -y
sudo pip install --upgrade pip

# install Docker-compose
echo "Checking to install Docker-compose…"
sudo pip install docker-compose==1.8.1 

# ----------------------------------------------------------------
# Install Golang
# ----------------------------------------------------------------
GO_VER=1.8.3
GO_URL=https://storage.googleapis.com/golang/go${GO_VER}.linux-amd64.tar.gz

# Set Go environment variables needed by other scripts
export GOPATH="/opt/gopath"
export GOROOT="/opt/go"
PATH=$GOROOT/bin:$GOPATH/bin:$PATH

cat <<EOF >/etc/profile.d/goroot.sh
export GOROOT=$GOROOT
export GOPATH=$GOPATH
export PATH=\$PATH:$GOROOT/bin:$GOPATH/bin
EOF

mkdir -p $GOROOT

curl -sL $GO_URL | (cd $GOROOT && tar --strip-components 1 -xz)

# Ensure permissions are set for GOPATH
sudo chown -R ${USER}:${USERGROUP} $GOPATH
