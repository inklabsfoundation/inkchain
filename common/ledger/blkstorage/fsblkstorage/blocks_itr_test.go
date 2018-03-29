/*
Copyright Ziggurat Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fsblkstorage

import (
	"testing"
	"time"

	"github.com/inklabsfoundation/inkchain/common/ledger/testutil"
	"github.com/inklabsfoundation/inkchain/protos/common"
)

func TestBlocksItrBlockingNext(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()
	blkfileMgr := blkfileMgrWrapper.blockfileMgr

	blocks := testutil.ConstructTestBlocks(t, 10)
	blkfileMgrWrapper.addBlocks(blocks[:5])

	itr, err := blkfileMgr.retrieveBlocks(1)
	defer itr.Close()
	testutil.AssertNoError(t, err, "")
	doneChan := make(chan bool)
	go testIterateAndVerify(t, itr, blocks[1:], doneChan)
	for {
		if itr.blockNumToRetrieve == 5 {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	testAppendBlocks(blkfileMgrWrapper, blocks[5:7])
	blkfileMgr.moveToNextFile()
	time.Sleep(time.Millisecond * 10)
	testAppendBlocks(blkfileMgrWrapper, blocks[7:])
	<-doneChan
}

func TestBlockItrClose(t *testing.T) {
	env := newTestEnv(t, NewConf(testPath(), 0))
	defer env.Cleanup()
	blkfileMgrWrapper := newTestBlockfileWrapper(env, "testLedger")
	defer blkfileMgrWrapper.close()
	blkfileMgr := blkfileMgrWrapper.blockfileMgr

	blocks := testutil.ConstructTestBlocks(t, 5)
	blkfileMgrWrapper.addBlocks(blocks)

	itr, err := blkfileMgr.retrieveBlocks(1)
	testutil.AssertNoError(t, err, "")

	bh, _ := itr.Next()
	testutil.AssertNotNil(t, bh)
	itr.Close()

	bh, err = itr.Next()
	testutil.AssertNoError(t, err, "")
	testutil.AssertNil(t, bh)
}

func testIterateAndVerify(t *testing.T, itr *blocksItr, blocks []*common.Block, doneChan chan bool) {
	blocksIterated := 0
	for {
		t.Logf("blocksIterated: %v", blocksIterated)
		block, err := itr.Next()
		testutil.AssertNoError(t, err, "")
		testutil.AssertEquals(t, block, blocks[blocksIterated])
		blocksIterated++
		if blocksIterated == len(blocks) {
			break
		}
	}
	doneChan <- true
}

func testAppendBlocks(blkfileMgrWrapper *testBlockfileMgrWrapper, blocks []*common.Block) {
	blkfileMgrWrapper.addBlocks(blocks)
}
