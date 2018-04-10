#! /bin/bash

source p_init.sh

solar deploy contracts/INK.sol

solar deploy contracts/XCPlugin.sol '["'${myPlatform}'"]'

solar deploy contracts/XC.sol '["'${myPlatform}'"]'
