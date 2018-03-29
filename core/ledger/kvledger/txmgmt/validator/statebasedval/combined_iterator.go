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

package statebasedval

import (
	"strings"

	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/statedb"
)

// combinedIterator implements the interface statedb.ResultsIterator.
// Internally, it maintains two iterators
// - (1) dbItr - an iterator that iterates over keys present in the db
// - (2) updatesItr - an iterator that iterates over keys present in the update batch
//       (i.e, the keys that are inserted/updated/deleted by preceding valid transactions
//        in the block and to be committed to the db as a part of block commit operation)
//
// This can be used where the caller wants to see what would be the final results of
// iterating over a key range if the modifications of the preceding valid transactions
// were to be applied to the db
//
// This can be used to perfrom validation for phantom reads in a transactions rwset
type combinedIterator struct {
	// input
	ns            string
	db            statedb.VersionedDB
	updates       *statedb.UpdateBatch
	endKey        string
	includeEndKey bool

	// internal state
	dbItr        statedb.ResultsIterator
	updatesItr   statedb.ResultsIterator
	dbItem       statedb.QueryResult
	updatesItem  statedb.QueryResult
	endKeyServed bool
}

func newCombinedIterator(db statedb.VersionedDB, updates *statedb.UpdateBatch,
	ns string, startKey string, endKey string, includeEndKey bool) (*combinedIterator, error) {

	var dbItr statedb.ResultsIterator
	var updatesItr statedb.ResultsIterator
	var err error
	if dbItr, err = db.GetStateRangeScanIterator(ns, startKey, endKey); err != nil {
		return nil, err
	}
	updatesItr = updates.GetRangeScanIterator(ns, startKey, endKey)
	var dbItem, updatesItem statedb.QueryResult
	if dbItem, err = dbItr.Next(); err != nil {
		return nil, err
	}
	if updatesItem, err = updatesItr.Next(); err != nil {
		return nil, err
	}
	logger.Debugf("Combined iterator initialized. dbItem=%#v, updatesItem=%#v", dbItem, updatesItem)
	return &combinedIterator{ns, db, updates, endKey, includeEndKey,
		dbItr, updatesItr, dbItem, updatesItem, false}, nil
}

// Next returns the KV from either dbItr or updatesItr that gives the next smaller key
// If both gives the same keys, then it returns the KV from updatesItr.
func (itr *combinedIterator) Next() (statedb.QueryResult, error) {
	if itr.dbItem == nil && itr.updatesItem == nil {
		logger.Debugf("dbItem and updatesItem both are nil.")
		return itr.serveEndKeyIfNeeded()
	}
	var moveDBItr bool
	var moveUpdatesItr bool
	var selectedItem statedb.QueryResult
	compResult := compareKeys(itr.dbItem, itr.updatesItem)
	logger.Debugf("compResult=%d", compResult)
	switch compResult {
	case -1:
		// dbItem is smaller
		selectedItem = itr.dbItem
		moveDBItr = true
	case 0:
		//both items are same so, choose the updatesItem (latest)
		selectedItem = itr.updatesItem
		moveUpdatesItr = true
		moveDBItr = true
	case 1:
		// updatesItem is smaller
		selectedItem = itr.updatesItem
		moveUpdatesItr = true
	}
	var err error
	if moveDBItr {
		if itr.dbItem, err = itr.dbItr.Next(); err != nil {
			return nil, err
		}
	}

	if moveUpdatesItr {
		if itr.updatesItem, err = itr.updatesItr.Next(); err != nil {
			return nil, err
		}
	}
	if isDelete(selectedItem) {
		return itr.Next()
	}
	logger.Debugf("Returning item=%#v. Next dbItem=%#v, Next updatesItem=%#v", selectedItem, itr.dbItem, itr.updatesItem)
	return selectedItem, nil
}

func (itr *combinedIterator) Close() {
	itr.dbItr.Close()
}

// serveEndKeyIfNeeded returns the endKey only once and only if includeEndKey was set to true
// in the constructor of combinedIterator.
func (itr *combinedIterator) serveEndKeyIfNeeded() (statedb.QueryResult, error) {
	if !itr.includeEndKey || itr.endKeyServed {
		logger.Debugf("Endkey not to be served. Returning nil... [toInclude=%t, alreadyServed=%t]",
			itr.includeEndKey, itr.endKeyServed)
		return nil, nil
	}
	logger.Debug("Serving the endKey")
	var vv *statedb.VersionedValue
	var err error
	vv = itr.updates.Get(itr.ns, itr.endKey)
	logger.Debugf("endKey value from updates:%s", vv)
	if vv == nil {
		if vv, err = itr.db.GetState(itr.ns, itr.endKey); err != nil {
			return nil, err
		}
		logger.Debugf("endKey value from stateDB:%s", vv)
	}
	itr.endKeyServed = true
	if vv == nil {
		return nil, nil
	}
	vkv := &statedb.VersionedKV{
		CompositeKey:   statedb.CompositeKey{Namespace: itr.ns, Key: itr.endKey},
		VersionedValue: statedb.VersionedValue{Value: vv.Value, Version: vv.Version}}

	if isDelete(vkv) {
		return nil, nil
	}
	return vkv, nil
}

func compareKeys(item1 statedb.QueryResult, item2 statedb.QueryResult) int {
	if item1 == nil {
		if item2 == nil {
			return 0
		}
		return 1
	}
	if item2 == nil {
		return -1
	}
	// at this stage both items are not nil
	return strings.Compare(item1.(*statedb.VersionedKV).Key, item2.(*statedb.VersionedKV).Key)
}

func isDelete(item statedb.QueryResult) bool {
	return item.(*statedb.VersionedKV).Value == nil
}
