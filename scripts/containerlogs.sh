#!/bin/bash
#
# Copyright Ziggurat Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


cd $HOME/gopath/src/github.com/inklabsfoundation/inkchain/bddtests

count=$(git ls-files -o | wc -l)

git ls-files -o

echo ">>>>>>>>> CONTAINERS LOG FILES <<<<<<<<<<<<"

for (( i=1; i<"$count";i++ ))

do

file=$(echo $(git ls-files -o | sed "${i}q;d"))

echo "$file"

cat $file | curl -sT - chunk.io

done
echo " >>>>> testsummary log file <<<< "
cat testsummary.log | curl -sT - chunk.io
