#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


if [ "$2" != "inkchain" ]; then

        echo " Pull Request number is $1 "
        echo " User Name is $2 "
	echo " Repository Name is $3 "

mkdir -p $HOME/gopath/src/github.com/inkchain

	echo "inklabsfoundation/inkchain folder created"

git clone -ql $HOME/gopath/src/github.com/$2/$3 $HOME/gopath/src/github.com/inklabsfoundation/inkchain

	echo "linked $2 user repo into inklabsfoundation/inkchain folder"

fi
