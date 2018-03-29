#! /bin/bash

source p_init.sh

# ===  init amount ===
qcli sendtoaddress $XCPluginOwner 20
qcli sendtoaddress $XCOwner 20
qcli sendtoaddress $XCAddress 20
qcli sendtoaddress $INKOwner 20
