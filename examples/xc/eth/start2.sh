#! /bin/bash

0xacfee5cc1c4b41e4928291dc5bac5e3fbcb9f21e

personal.newAccount("123456")

miner.start();admin.sleepBlocks(10);miner.stop();

personal.unlockAccount(eth.accounts[0],"123456")
