#! /bin/bash

source p_init.sh

solar deploy contracts/INK.sol --gasLimit=300000

solar deploy contracts/XCPlugin.sol '["7174756d00000000000000000000000000000000000000000000000000000000"]' --gasLimit=500000

solar deploy contracts/XC.sol '["7174756d00000000000000000000000000000000000000000000000000000000"]' --gasLimit=500000
