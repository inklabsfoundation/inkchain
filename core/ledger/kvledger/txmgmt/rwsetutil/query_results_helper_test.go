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

package rwsetutil

import (
	"testing"

	"fmt"

	"github.com/golang/protobuf/proto"
	bccspfactory "github.com/inklabsfoundation/inkchain/bccsp/factory"
	"github.com/inklabsfoundation/inkchain/common/ledger/testutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/version"
	"github.com/inklabsfoundation/inkchain/protos/ledger/rwset/kvrwset"
)

func TestQueryResultHelper_NoResults(t *testing.T) {
	helper, _ := NewRangeQueryResultsHelper(true, 3)
	r, h, err := helper.Done()
	testutil.AssertNoError(t, err, "")
	testutil.AssertNil(t, h)
	testutil.AssertNil(t, r)
}

func TestQueryResultHelper_HashNotEnabled(t *testing.T) {
	kvReads := buildTestKVReads(t, 5)
	r, h := buildTestResults(t, false, 3, kvReads)
	testutil.AssertNil(t, h)
	testutil.AssertEquals(t, r, kvReads)
}

func TestQueryResultHelper_ResultsNoMoreThanMaxDegree(t *testing.T) {
	maxDegree := 3
	numResults := 3
	kvReads := buildTestKVReads(t, numResults)
	r, h := buildTestResults(t, true, maxDegree, kvReads)
	testutil.AssertNil(t, h)
	testutil.AssertEquals(t, r, kvReads)
}

func TestQueryResultHelper_Hash_OneLevel(t *testing.T) {
	maxDegree := 3
	kvReads := buildTestKVReads(t, 9)
	r, h := buildTestResults(t, true, maxDegree, kvReads)
	level1_1 := computeTestHashKVReads(t, kvReads[0:4])
	level1_2 := computeTestHashKVReads(t, kvReads[4:8])
	level1_3 := computeTestHashKVReads(t, kvReads[8:])
	testutil.AssertNil(t, r)
	testutil.AssertEquals(t, h, &kvrwset.QueryReadsMerkleSummary{
		MaxDegree:      uint32(maxDegree),
		MaxLevel:       1,
		MaxLevelHashes: hashesToBytes([]Hash{level1_1, level1_2, level1_3})})
}

func TestQueryResultHelper_Hash_TwoLevel(t *testing.T) {
	maxDegree := 3
	kvReads := buildTestKVReads(t, 25)
	r, h := buildTestResults(t, true, maxDegree, kvReads)
	level1_1 := computeTestHashKVReads(t, kvReads[0:4])
	level1_2 := computeTestHashKVReads(t, kvReads[4:8])
	level1_3 := computeTestHashKVReads(t, kvReads[8:12])
	level1_4 := computeTestHashKVReads(t, kvReads[12:16])
	level1_5 := computeTestHashKVReads(t, kvReads[16:20])
	level1_6 := computeTestHashKVReads(t, kvReads[20:24])
	level1_7 := computeTestHashKVReads(t, kvReads[24:])

	level2_1 := computeTestCombinedHash(t, level1_1, level1_2, level1_3, level1_4)
	level2_2 := computeTestCombinedHash(t, level1_5, level1_6, level1_7)
	testutil.AssertNil(t, r)
	testutil.AssertEquals(t, h, &kvrwset.QueryReadsMerkleSummary{
		MaxDegree:      uint32(maxDegree),
		MaxLevel:       2,
		MaxLevelHashes: hashesToBytes([]Hash{level2_1, level2_2})})
}

func TestQueryResultHelper_Hash_ThreeLevel(t *testing.T) {
	maxDegree := 3
	kvReads := buildTestKVReads(t, 65)
	r, h := buildTestResults(t, true, maxDegree, kvReads)
	level1_1 := computeTestHashKVReads(t, kvReads[0:4])
	level1_2 := computeTestHashKVReads(t, kvReads[4:8])
	level1_3 := computeTestHashKVReads(t, kvReads[8:12])
	level1_4 := computeTestHashKVReads(t, kvReads[12:16])
	level1_5 := computeTestHashKVReads(t, kvReads[16:20])
	level1_6 := computeTestHashKVReads(t, kvReads[20:24])
	level1_7 := computeTestHashKVReads(t, kvReads[24:28])
	level1_8 := computeTestHashKVReads(t, kvReads[28:32])
	level1_9 := computeTestHashKVReads(t, kvReads[32:36])
	level1_10 := computeTestHashKVReads(t, kvReads[36:40])
	level1_11 := computeTestHashKVReads(t, kvReads[40:44])
	level1_12 := computeTestHashKVReads(t, kvReads[44:48])
	level1_13 := computeTestHashKVReads(t, kvReads[48:52])
	level1_14 := computeTestHashKVReads(t, kvReads[52:56])
	level1_15 := computeTestHashKVReads(t, kvReads[56:60])
	level1_16 := computeTestHashKVReads(t, kvReads[60:64])
	level1_17 := computeTestHashKVReads(t, kvReads[64:])

	level2_1 := computeTestCombinedHash(t, level1_1, level1_2, level1_3, level1_4)
	level2_2 := computeTestCombinedHash(t, level1_5, level1_6, level1_7, level1_8)
	level2_3 := computeTestCombinedHash(t, level1_9, level1_10, level1_11, level1_12)
	level2_4 := computeTestCombinedHash(t, level1_13, level1_14, level1_15, level1_16)

	level3_1 := computeTestCombinedHash(t, level2_1, level2_2, level2_3, level2_4)
	level3_2 := level1_17
	testutil.AssertNil(t, r)
	testutil.AssertEquals(t, h, &kvrwset.QueryReadsMerkleSummary{
		MaxDegree:      uint32(maxDegree),
		MaxLevel:       3,
		MaxLevelHashes: hashesToBytes([]Hash{level3_1, level3_2})})
}

func TestQueryResultHelper_Hash_MaxLevelIncrementNeededInDone(t *testing.T) {
	maxDegree := 2
	kvReads := buildTestKVReads(t, 24)
	r, h := buildTestResults(t, true, maxDegree, kvReads)
	level1_1 := computeTestHashKVReads(t, kvReads[0:3])
	level1_2 := computeTestHashKVReads(t, kvReads[3:6])
	level1_3 := computeTestHashKVReads(t, kvReads[6:9])
	level1_4 := computeTestHashKVReads(t, kvReads[9:12])
	level1_5 := computeTestHashKVReads(t, kvReads[12:15])
	level1_6 := computeTestHashKVReads(t, kvReads[15:18])
	level1_7 := computeTestHashKVReads(t, kvReads[18:21])
	level1_8 := computeTestHashKVReads(t, kvReads[21:24])

	level2_1 := computeTestCombinedHash(t, level1_1, level1_2, level1_3)
	level2_2 := computeTestCombinedHash(t, level1_4, level1_5, level1_6)
	level2_3 := computeTestCombinedHash(t, level1_7, level1_8)

	level3_1 := computeTestCombinedHash(t, level2_1, level2_2, level2_3)

	testutil.AssertNil(t, r)
	testutil.AssertEquals(t, h, &kvrwset.QueryReadsMerkleSummary{
		MaxDegree:      uint32(maxDegree),
		MaxLevel:       3,
		MaxLevelHashes: hashesToBytes([]Hash{level3_1})})
}

func TestQueryResultHelper_Hash_FirstLevelSkipNeededInDone(t *testing.T) {
	maxDegree := 2
	kvReads := buildTestKVReads(t, 45)
	r, h := buildTestResults(t, true, maxDegree, kvReads)
	level1_1 := computeTestHashKVReads(t, kvReads[0:3])
	level1_2 := computeTestHashKVReads(t, kvReads[3:6])
	level1_3 := computeTestHashKVReads(t, kvReads[6:9])
	level1_4 := computeTestHashKVReads(t, kvReads[9:12])
	level1_5 := computeTestHashKVReads(t, kvReads[12:15])
	level1_6 := computeTestHashKVReads(t, kvReads[15:18])
	level1_7 := computeTestHashKVReads(t, kvReads[18:21])
	level1_8 := computeTestHashKVReads(t, kvReads[21:24])
	level1_9 := computeTestHashKVReads(t, kvReads[24:27])
	level1_10 := computeTestHashKVReads(t, kvReads[27:30])
	level1_11 := computeTestHashKVReads(t, kvReads[30:33])
	level1_12 := computeTestHashKVReads(t, kvReads[33:36])
	level1_13 := computeTestHashKVReads(t, kvReads[36:39])
	level1_14 := computeTestHashKVReads(t, kvReads[39:42])
	level1_15 := computeTestHashKVReads(t, kvReads[42:45])

	level2_1 := computeTestCombinedHash(t, level1_1, level1_2, level1_3)
	level2_2 := computeTestCombinedHash(t, level1_4, level1_5, level1_6)
	level2_3 := computeTestCombinedHash(t, level1_7, level1_8, level1_9)
	level2_4 := computeTestCombinedHash(t, level1_10, level1_11, level1_12)
	level2_5 := computeTestCombinedHash(t, level1_13, level1_14, level1_15)

	level3_1 := computeTestCombinedHash(t, level2_1, level2_2, level2_3)
	level3_2 := computeTestCombinedHash(t, level2_4, level2_5)

	testutil.AssertNil(t, r)
	testutil.AssertEquals(t, h, &kvrwset.QueryReadsMerkleSummary{
		MaxDegree:      uint32(maxDegree),
		MaxLevel:       3,
		MaxLevelHashes: hashesToBytes([]Hash{level3_1, level3_2})})
}

func buildTestResults(t *testing.T, enableHashing bool, maxDegree int, kvReads []*kvrwset.KVRead) ([]*kvrwset.KVRead, *kvrwset.QueryReadsMerkleSummary) {
	helper, _ := NewRangeQueryResultsHelper(enableHashing, uint32(maxDegree))
	for _, kvRead := range kvReads {
		helper.AddResult(kvRead)
	}
	r, h, err := helper.Done()
	testutil.AssertNoError(t, err, "")
	return r, h
}

func buildTestKVReads(t *testing.T, num int) []*kvrwset.KVRead {
	kvreads := []*kvrwset.KVRead{}
	for i := 0; i < num; i++ {
		kvreads = append(kvreads, NewKVRead(fmt.Sprintf("key_%d", i), version.NewHeight(1, uint64(i))))
	}
	return kvreads
}

func computeTestHashKVReads(t *testing.T, kvReads []*kvrwset.KVRead) Hash {
	queryReads := &kvrwset.QueryReads{}
	queryReads.KvReads = kvReads
	b, err := proto.Marshal(queryReads)
	testutil.AssertNoError(t, err, "")
	h, err := bccspfactory.GetDefault().Hash(b, hashOpts)
	testutil.AssertNoError(t, err, "")
	return h
}

func computeTestCombinedHash(t *testing.T, hashes ...Hash) Hash {
	h, err := computeCombinedHash(hashes)
	testutil.AssertNoError(t, err, "")
	return h
}
