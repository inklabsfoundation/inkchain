#!/bin/bash
#
# Copyright Ziggurat Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

# Shut down the Docker containers for the system tests.
docker-compose -f docker-compose.yml stop && docker-compose -f docker-compose.yml down

# remove chaincode docker images
docker rmi $(docker images dev-* -q)

# Your system is now clean
