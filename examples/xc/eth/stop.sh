#! /bin/bash

kill -9 $(lsof -i:8545 | awk '{print $2}')