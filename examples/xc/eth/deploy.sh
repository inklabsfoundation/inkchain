#! /bin/bash

source p_init.sh

solar deploy ../docker-qtum/contracts/INK.sol

solar deploy ../docker-qtum/contracts/XCPlugin.sol '["7174756d00000000000000000000000000000000000000000000000000000000"]'

solar deploy ../docker-qtum/contracts/XC.sol '["7174756d00000000000000000000000000000000000000000000000000000000"]'
