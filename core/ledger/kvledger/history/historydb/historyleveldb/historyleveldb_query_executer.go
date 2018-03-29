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

package historyleveldb

import (
	"errors"

	"encoding/json"

	commonledger "github.com/inklabsfoundation/inkchain/common/ledger"
	"github.com/inklabsfoundation/inkchain/common/ledger/blkstorage"
	"github.com/inklabsfoundation/inkchain/common/ledger/util"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/history/historydb"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/ledgerutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/ledgerconfig"
	"github.com/inklabsfoundation/inkchain/core/wallet"
	"github.com/inklabsfoundation/inkchain/protos/common"
	"github.com/inklabsfoundation/inkchain/protos/ledger/queryresult"
	putils "github.com/inklabsfoundation/inkchain/protos/utils"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// LevelHistoryDBQueryExecutor is a query executor against the LevelDB history DB
type LevelHistoryDBQueryExecutor struct {
	historyDB  *historyDB
	blockStore blkstorage.BlockStore
}

// GetHistoryForKey implements method in interface `ledger.HistoryQueryExecutor`
func (q *LevelHistoryDBQueryExecutor) GetHistoryForKey(namespace string, key string) (commonledger.ResultsIterator, error) {

	if ledgerconfig.IsHistoryDBEnabled() == false {
		return nil, errors.New("History tracking not enabled - historyDatabase is false")
	}

	var compositeStartKey []byte
	var compositeEndKey []byte
	compositeStartKey = historydb.ConstructPartialCompositeHistoryKey(namespace, key, false)
	compositeEndKey = historydb.ConstructPartialCompositeHistoryKey(namespace, key, true)

	// range scan to find any history records starting with namespace~key
	dbItr := q.historyDB.db.GetIterator(compositeStartKey, compositeEndKey)
	return newHistoryScanner(compositeStartKey, namespace, key, dbItr, q.blockStore), nil
}

//historyScanner implements ResultsIterator for iterating through history results
type historyScanner struct {
	compositePartialKey []byte //compositePartialKey includes namespace~key
	namespace           string
	key                 string
	dbItr               iterator.Iterator
	blockStore          blkstorage.BlockStore
}

func newHistoryScanner(compositePartialKey []byte, namespace string, key string,
	dbItr iterator.Iterator, blockStore blkstorage.BlockStore) *historyScanner {
	return &historyScanner{compositePartialKey, namespace, key, dbItr, blockStore}
}

func (scanner *historyScanner) Next() (commonledger.QueryResult, error) {
	if !scanner.dbItr.Next() {
		return nil, nil
	}
	historyKey := scanner.dbItr.Key() // history key is in the form namespace~key~blocknum~trannum

	// SplitCompositeKey(namespace~key~blocknum~trannum, namespace~key~) will return the blocknum~trannum in second position
	_, blockNumTranNumBytes := historydb.SplitCompositeHistoryKey(historyKey, scanner.compositePartialKey)
	blockNum, bytesConsumed := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes[0:])
	tranNum, _ := util.DecodeOrderPreservingVarUint64(blockNumTranNumBytes[bytesConsumed:])
	logger.Debugf("Found history record for namespace:%s key:%s at blockNumTranNum %v:%v\n",
		scanner.namespace, scanner.key, blockNum, tranNum)

	// Get the transaction from block storage that is associated with this history record
	tranEnvelope, err := scanner.blockStore.RetrieveTxByBlockNumTranNum(blockNum, tranNum)
	if err != nil {
		return nil, err
	}

	// Get the txid, key write value, timestamp, and delete indicator associated with this transaction
	queryResult, err := getKeyModificationFromTran(tranEnvelope, scanner.namespace, scanner.key)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Found historic key value for namespace:%s key:%s from transaction %s\n",
		scanner.namespace, scanner.key, queryResult.(*queryresult.KeyModification).TxId)
	return queryResult, nil
}

func (scanner *historyScanner) Close() {
	scanner.dbItr.Release()
}

// getTxIDandKeyWriteValueFromTran inspects a transaction for writes to a given key
func getKeyModificationFromTran(tranEnvelope *common.Envelope, namespace string, key string) (commonledger.QueryResult, error) {
	logger.Debugf("Entering getKeyModificationFromTran()\n", namespace, key)

	// extract action from the envelope
	payload, err := putils.GetPayload(tranEnvelope)
	if err != nil {
		return nil, err
	}

	tx, err := putils.GetTransaction(payload.Data)
	if err != nil {
		return nil, err
	}

	_, respPayload, err := putils.GetPayloads(tx.Actions[0])
	if err != nil {
		return nil, err
	}

	chdr, err := putils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, err
	}

	txID := chdr.TxId
	timestamp := chdr.Timestamp

	ledgerSet := &ledgerutil.LedgerSet{}
	if err = ledgerSet.FromProtoBytes(respPayload.Results); err != nil {
		return nil, err
	}
	txRWSet := ledgerSet.TxRwSet

	// Get the Result from the Action and then Unmarshal
	// it into a TxReadWriteSet using custom unmarshalling

	if namespace == wallet.WALLET_NAMESPACE {
		if ledgerSet.TranSet.From == key {
			account := &wallet.Account{}
			//***************** [issue] multiple key
			accountBytes, jsonErr := json.Marshal(account)
			if jsonErr != nil {
				return nil, errors.New("error marshaling account")
			}
			return &queryresult.KeyModification{TxId: txID, Value: accountBytes,
				Timestamp: timestamp, IsDelete: false}, nil
		}
		for _, senderTranSet := range ledgerSet.TranSet.KvTranSet.Trans {
			if senderTranSet.To == key {
				account := &wallet.Account{}
				//***************** [issue] multiple key
				accountBytes, jsonErr := json.Marshal(account)
				if jsonErr != nil {
					return nil, errors.New("error marshaling account")
				}
				return &queryresult.KeyModification{TxId: txID, Value: accountBytes,
					Timestamp: timestamp, IsDelete: false}, nil
			}
		}
	}
	// look for the namespace and key by looping through the transaction's ReadWriteSets
	for _, nsRWSet := range txRWSet.NsRwSets {
		if nsRWSet.NameSpace == namespace {
			// got the correct namespace, now find the key write
			for _, kvWrite := range nsRWSet.KvRwSet.Writes {
				if kvWrite.Key == key {
					return &queryresult.KeyModification{TxId: txID, Value: kvWrite.Value,
						Timestamp: timestamp, IsDelete: kvWrite.IsDelete}, nil
				}
			} // end keys loop
			return nil, errors.New("Key not found in namespace's writeset")
		} // end if
	} //end namespaces loop
	return nil, errors.New("Namespace not found in transaction's ReadWriteSets")

}
