/*
Copyright IBM Corp. 2016 All Rights Reserved.

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
	"encoding/json"

	"math/big"

	"fmt"

	"github.com/inklabsfoundation/inkchain/common/flogging"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/ledgerutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/transutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/version"
	"github.com/inklabsfoundation/inkchain/core/ledger/util"
	"github.com/inklabsfoundation/inkchain/core/wallet"
	"github.com/inklabsfoundation/inkchain/core/wallet/ink"
	"github.com/inklabsfoundation/inkchain/core/wallet/ink/impl"
	"github.com/inklabsfoundation/inkchain/protos/common"
	"github.com/inklabsfoundation/inkchain/protos/ledger/rwset/kvrwset"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet/kvtranset"
	"github.com/inklabsfoundation/inkchain/protos/peer"
	putils "github.com/inklabsfoundation/inkchain/protos/utils"
)

var logger = flogging.MustGetLogger("statevalidator")

// Validator validates a tx against the latest committed state
// and preceding valid transactions with in the same block
type Validator struct {
	db            statedb.VersionedDB
	inkCalculator ink.InkAlg
}

// NewValidator constructs StateValidator
func NewValidator(db statedb.VersionedDB) *Validator {
	return &Validator{db, impl.NewSimpleInkAlg()}
}

func (v *Validator) validateCounterAndInk(sender string, cis *peer.ChaincodeInvocationSpec, batch *statedb.TransferBatch, ledgerByteCount int) (*big.Int, error) {
	//validate counter
	counterValidated := false
	senderCounter, ok := batch.GetSenderCounter(sender)
	if ok && senderCounter != 0 {
		if senderCounter != cis.SenderSpec.Counter {
			return nil, fmt.Errorf("invalid counter")
		}
		counterValidated = true
	}

	versionedValue, err := v.db.GetState(wallet.WALLET_NAMESPACE, sender)
	if err != nil {
		return nil, err
	}
	if versionedValue != nil {
		account := &wallet.Account{}
		jsonErr := json.Unmarshal(versionedValue.Value, account)
		if jsonErr != nil {
			return nil, jsonErr
		}
		if !counterValidated && account.Counter != cis.SenderSpec.Counter {
			return nil, fmt.Errorf("counter mismatch")
		}
		inkFee, err := v.inkCalculator.CalcInk(ledgerByteCount)
		if inkFee.Cmp(big.NewInt(0)) > 0 {
			if err != nil {
				return nil, fmt.Errorf("commiter: error when calculating ink.")
			}
			mtcBalance, ok := account.Balance[wallet.MAIN_BALANCE_NAME]
			if !ok {
				return nil, fmt.Errorf("commiter: insuffient mtc balance for ink consumption.")
			}
			if batch.ExistsFrom(sender) {
				balanceUpdate := batch.GetBalanceUpdate(sender, wallet.MAIN_BALANCE_NAME)
				if balanceUpdate != nil {
					mtcBalance = mtcBalance.Add(mtcBalance, balanceUpdate)
				}
			}
			inkLimit, ok := new(big.Int).SetString(string(cis.SenderSpec.InkLimit), 10)
			if !ok {
				return nil, fmt.Errorf("commiter: invalid inklimit.")
			}
			if inkFee.Cmp(inkLimit) > 0 {
				return nil, fmt.Errorf("commiter: ink exceeds inkLimit.")
			}
			if !ok || mtcBalance.Cmp(inkFee) < 0 {
				return nil, fmt.Errorf("commiter: insuffient balance for ink consumption.")
			}
		}
		return inkFee, nil
	}
	return nil, fmt.Errorf("sender not exists")
}

//validate endorser transaction
func (v *Validator) validateEndorserTX(envBytes []byte, doMVCCValidation bool, updates *statedb.UpdateBatch, transferUpdates *statedb.TransferBatch) (*rwsetutil.TxRwSet, *transutil.TranSet, *transutil.SenderCounter, peer.TxValidationCode, error) {
	// extract actions from the envelope message
	cis, respPayload, err := putils.GetActionFromEnvelope(envBytes)
	if err != nil {
		return nil, nil, nil, peer.TxValidationCode_NIL_TXACTION, nil
	}
	//preparation for extracting RWSet from transaction
	// Get the Result from the Action
	// and then Unmarshal it into a TxReadWriteSet using custom unmarshalling
	ledgerSet := &ledgerutil.LedgerSet{}
	if err = ledgerSet.FromProtoBytes(respPayload.Results); err != nil {
		return nil, nil, nil, peer.TxValidationCode_INVALID_OTHER_REASON, nil
	}

	txResult := peer.TxValidationCode_VALID

	senderCounter := &transutil.SenderCounter{}
	// check signature
	if cis.ChaincodeSpec.ChaincodeId.Name != "lscc" && cis.ChaincodeSpec.ChaincodeId.Name != "ascc" {
		if cis.Sig == nil || cis.SenderSpec == nil {
			return nil, nil, nil, peer.TxValidationCode_BAD_SIGNATURE, nil
		}
		hash, err := wallet.GetInvokeHash(cis.ChaincodeSpec, cis.IdGenerationAlg, cis.SenderSpec)
		if err != nil {
			return nil, nil, nil, peer.TxValidationCode_BAD_SIGNATURE, nil
		}
		sender, err := wallet.GetSenderFromSignature(hash, cis.Sig)
		senderStr := sender.ToString()
		if senderStr != string(cis.SenderSpec.Sender) || senderStr != ledgerSet.TranSet.From {
			return nil, nil, nil, peer.TxValidationCode_BAD_SIGNATURE, nil
		}

		inkFee, err := v.validateCounterAndInk(senderStr, cis, transferUpdates, len(respPayload.Results))
		if err != nil {
			fmt.Println(err)
			return nil, nil, nil, peer.TxValidationCode_BAD_COUNTER, nil
		}
		senderCounter.Sender = sender.ToString()
		senderCounter.Counter = cis.SenderSpec.Counter
		senderCounter.Ink = inkFee
	}
	//mvccvalidation, may invalidate transaction
	if doMVCCValidation {
		// validate signature
		if ledgerSet.TxRwSet != nil {
			if txResult, err = v.validateTx(ledgerSet.TxRwSet, updates); err != nil || txResult != peer.TxValidationCode_VALID {
				ledgerSet.TxRwSet = nil
				ledgerSet.TranSet = nil
				return nil, nil, nil, txResult, err
			}
		}
		//validate transfer
		if ledgerSet.TranSet != nil {
			if txResult, err = v.validateTrans(ledgerSet.TranSet, transferUpdates); err != nil || txResult != peer.TxValidationCode_VALID {
				ledgerSet.TxRwSet = nil
				ledgerSet.TranSet = nil
				return nil, nil, nil, txResult, err
			}
		}
	}
	return ledgerSet.TxRwSet, ledgerSet.TranSet, senderCounter, txResult, err
}

// ValidateAndPrepareBatch implements method in Validator interface
func (v *Validator) ValidateAndPrepareBatch(block *common.Block, doMVCCValidation bool) (*statedb.UpdateBatch, error) {
	logger.Debugf("New block arrived for validation:%#v, doMVCCValidation=%t", block, doMVCCValidation)
	updates := statedb.NewUpdateBatch()
	transferUpdates := statedb.NewTransferBatch()
	logger.Debugf("Validating a block with [%d] transactions", len(block.Data.Data))

	// Committer validator has already set validation flags based on well formed tran checks
	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	// Precaution in case committer validator has not added validation flags yet
	if len(txsFilter) == 0 {
		txsFilter = util.NewTxValidationFlags(len(block.Data.Data))
		block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
	}
	lastTxIndex := 0
	for txIndex, envBytes := range block.Data.Data {
		lastTxIndex = txIndex
		if txsFilter.IsInvalid(txIndex) {
			// Skiping invalid transaction
			logger.Warningf("Block [%d] Transaction index [%d] marked as invalid by committer. Reason code [%d]",
				block.Header.Number, txIndex, txsFilter.Flag(txIndex))
			continue
		}

		env, err := putils.GetEnvelopeFromBlock(envBytes)
		if err != nil {
			return nil, err
		}

		payload, err := putils.GetPayload(env)
		if err != nil {
			return nil, err
		}

		chdr, err := putils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return nil, err
		}

		txType := common.HeaderType(chdr.Type)

		if txType != common.HeaderType_ENDORSER_TRANSACTION {
			logger.Debugf("Skipping mvcc validation for Block [%d] Transaction index [%d] because, the transaction type is [%s]",
				block.Header.Number, txIndex, txType)
			continue
		}

		txRWSet, txTranSet, senderCounter, txResult, err := v.validateEndorserTX(envBytes, doMVCCValidation, updates, transferUpdates)

		if err != nil {
			return nil, err
		}
		txsFilter.SetFlag(txIndex, txResult)
		//txRWSet != nil => t is valid
		committingTxHeight := version.NewHeight(block.Header.Number, uint64(txIndex))
		if txResult == peer.TxValidationCode_VALID {
			if senderCounter != nil {
				transferUpdates.UpdateSender(senderCounter.Sender, senderCounter.Counter, senderCounter.Ink, committingTxHeight)
			}
			if txRWSet != nil {
				addWriteSetToBatch(txRWSet, committingTxHeight, updates)
			}
			if txTranSet != nil {
				addTranSetToBatch(txTranSet, committingTxHeight, transferUpdates)
			}
		}
		if txsFilter.IsValid(txIndex) {
			logger.Debugf("Block [%d] Transaction index [%d] TxId [%s] marked as valid by state validator",
				block.Header.Number, txIndex, chdr.TxId)
		} else {
			logger.Warningf("Block [%d] Transaction index [%d] TxId [%s] marked as invalid by state validator. Reason code [%d]",
				block.Header.Number, txIndex, chdr.TxId, txsFilter.Flag(txIndex))
		}

	}
	v.addTransferToRWSet(transferUpdates, updates, block.Header.FeeAddress, version.NewHeight(block.Header.Number, uint64(lastTxIndex)))
	block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
	return updates, nil
}

func addTranSetToBatch(tranSet *transutil.TranSet, txHeight *version.Height, transferBatch *statedb.TransferBatch) {
	from := tranSet.From
	if wallet.StringToAddress(from) == nil {
		return
	}
	for _, tran := range tranSet.KvTranSet.Trans {
		if wallet.StringToAddress(tran.To) == nil || tran.Amount == nil {
			continue
		}
		transferBatch.Put(from, tran.To, tran.BalanceType, tran.Amount, txHeight)
	}
}

func (v *Validator) validateTrans(tranSet *transutil.TranSet, updates *statedb.TransferBatch) (peer.TxValidationCode, error) {
	from := tranSet.From
	fromVer := tranSet.FromVer
	var accountBalance map[string]*big.Int
	versionedValue, err := v.db.GetState(wallet.WALLET_NAMESPACE, from)
	if err != nil {
		return peer.TxValidationCode_TRANSFER_CONFLICT, nil
	}
	var committedVersion *version.Height
	if versionedValue != nil {
		committedVersion = versionedValue.Version
		// check version
		if !version.AreSame(committedVersion, transutil.NewVersion(fromVer)) {
			logger.Debugf("Version mismatch for (sender) (%s). Committed version = [%s], Version in transferSet [%s]",
				from, committedVersion, fromVer)
			return peer.TxValidationCode_TRANSFER_CONFLICT, nil
		}
		// check sender balance
		account := &wallet.Account{}
		jsonErr := json.Unmarshal(versionedValue.Value, account)
		if jsonErr != nil {
			return peer.TxValidationCode_TRANSFER_CONFLICT, nil
		}
		accountBalance = account.Balance
		if accountBalance == nil {
			return peer.TxValidationCode_TRANSFER_CONFLICT, nil
		}
	}

	for _, kvTo := range tranSet.KvTranSet.Trans {
		if valid, err := v.validateKVTransfer(from, fromVer, kvTo, accountBalance, updates); !valid || err != nil {
			return peer.TxValidationCode_TRANSFER_CONFLICT, nil
		}
	}

	return peer.TxValidationCode_VALID, nil
}
func (v *Validator) validateKVTransfer(from string, fromVer *transet.Version, kvTo *kvtranset.KVTrans, accountBalance map[string]*big.Int, updates *statedb.TransferBatch) (bool, error) {
	// one (from,to) in one block
	if updates.ExistsTransfer(from, kvTo.To) {
		return false, nil
	}
	balance, ok := accountBalance[kvTo.BalanceType]
	if !ok {
		return false, nil
	}
	if updates.ExistsFrom(from) {
		balanceUpdate := updates.GetBalanceUpdate(from, kvTo.BalanceType)
		if balanceUpdate != nil {
			balance = balance.Add(balance, balanceUpdate)
		}
	}
	if balance.Cmp(new(big.Int).SetBytes(kvTo.Amount)) >= 0 {
		return true, nil
	}
	return false, nil
}
func (v *Validator) addTransferToRWSet(transferBatch *statedb.TransferBatch, batch *statedb.UpdateBatch, feeAddress []byte, txHeight *version.Height) {
	doInkDist := false
	if feeAddress != nil && len(feeAddress) == wallet.AddressLength {
		doInkDist = true
	}
	inkTotal := big.NewInt(0)
	for accountUpdate, _ := range transferBatch.Updates {
		balanceChange := transferBatch.GetAllBalanceUpdates(accountUpdate)
		inkFee, ok := transferBatch.GetSenderInk(accountUpdate)
		if !ok || inkFee == nil {
			inkFee = big.NewInt(0)
		}
		feeBalance, ok := balanceChange[wallet.MAIN_BALANCE_NAME]
		if !ok && inkFee.Cmp(big.NewInt(0)) > 0 {
			continue
		}
		versionedValue, err := v.db.GetState(wallet.WALLET_NAMESPACE, accountUpdate)
		if err != nil {
			continue
		}
		if doInkDist && feeBalance != nil {
			feeBalance = feeBalance.Sub(feeBalance, inkFee)
			inkTotal = inkTotal.Add(inkTotal, inkFee)
		}
		account := &wallet.Account{}
		if versionedValue != nil {
			jsonErr := json.Unmarshal(versionedValue.Value, account)
			if jsonErr != nil {
				continue
			}
			allBalance := account.Balance
			for key, value := range balanceChange {
				balance, ok := allBalance[key]
				if !ok || balance == nil {
					balance = big.NewInt(0)
					allBalance[key] = balance
				}
				balance = balance.Add(balance, value)
			}
		} else {
			account.Address = wallet.StringToAddress(accountUpdate)
			account.Balance = balanceChange
		}
		account.Counter, _ = transferBatch.GetSenderCounter(accountUpdate)
		accountBytes, err := json.Marshal(account)
		if err != nil {
			continue
		}
		batch.Put(wallet.WALLET_NAMESPACE, accountUpdate, accountBytes, transferBatch.GetBalanceVersion(accountUpdate))
	}

	if doInkDist {
		account := &wallet.Account{}
		var accountVersion *version.Height
		feeAccountName := wallet.BytesToAddress(feeAddress).ToString()
		if batch.Exists(wallet.WALLET_NAMESPACE, feeAccountName) {
			versionedValue := batch.Get(wallet.WALLET_NAMESPACE, feeAccountName)
			accountVersion = versionedValue.Version
			if versionedValue != nil {
				jsonErr := json.Unmarshal(versionedValue.Value, account)
				if jsonErr != nil {
					logger.Debugf("committer: fee account error")
					return
				}
			}
		} else {
			versionedValue, err := v.db.GetState(wallet.WALLET_NAMESPACE, feeAccountName)
			if err != nil {
				logger.Debugf("committer: fee account error")
				return
			}
			if versionedValue != nil {
				jsonErr := json.Unmarshal(versionedValue.Value, account)
				if jsonErr != nil {
					logger.Debugf("committer: fee account error")
					return
				}
				accountVersion = versionedValue.Version
			} else {
				account.Address = wallet.BytesToAddress(feeAddress)
				accountVersion = txHeight
			}
		}
		if account.Balance == nil {
			account.Balance = make(map[string]*big.Int)
		}
		feeBalance, ok := account.Balance[wallet.MAIN_BALANCE_NAME]
		if !ok {
			feeBalance = big.NewInt(0)
			account.Balance[wallet.MAIN_BALANCE_NAME] = feeBalance
		}
		feeBalance = feeBalance.Add(feeBalance, inkTotal)
		accountBytes, err := json.Marshal(account)
		if err != nil {
			logger.Debugf("committer: fee account error")
			return
		}
		batch.Put(wallet.WALLET_NAMESPACE, feeAccountName, accountBytes, accountVersion)
	}

}
func addWriteSetToBatch(txRWSet *rwsetutil.TxRwSet, txHeight *version.Height, batch *statedb.UpdateBatch) {
	for _, nsRWSet := range txRWSet.NsRwSets {
		ns := nsRWSet.NameSpace
		for _, kvWrite := range nsRWSet.KvRwSet.Writes {
			if kvWrite.IsDelete {
				batch.Delete(ns, kvWrite.Key, txHeight)
			} else {
				batch.Put(ns, kvWrite.Key, kvWrite.Value, txHeight)
			}
		}
	}
}

func (v *Validator) validateTx(txRWSet *rwsetutil.TxRwSet, updates *statedb.UpdateBatch) (peer.TxValidationCode, error) {
	for _, nsRWSet := range txRWSet.NsRwSets {
		ns := nsRWSet.NameSpace

		//filter out direct write to WALLET_NAMESPACE
		/*
			if ns == wallet.WALLET_NAMESPACE {
				return peer.TxValidationCode_TRANSFER_CONFLICT, nil
			}
		*/
		//*****
		if valid, err := v.validateReadSet(ns, nsRWSet.KvRwSet.Reads, updates); !valid || err != nil {
			if err != nil {
				return peer.TxValidationCode(-1), err
			}
			return peer.TxValidationCode_MVCC_READ_CONFLICT, nil
		}
		if valid, err := v.validateRangeQueries(ns, nsRWSet.KvRwSet.RangeQueriesInfo, updates); !valid || err != nil {
			if err != nil {
				return peer.TxValidationCode(-1), err
			}
			return peer.TxValidationCode_PHANTOM_READ_CONFLICT, nil
		}
	}
	return peer.TxValidationCode_VALID, nil
}

func (v *Validator) validateReadSet(ns string, kvReads []*kvrwset.KVRead, updates *statedb.UpdateBatch) (bool, error) {
	for _, kvRead := range kvReads {
		if valid, err := v.validateKVRead(ns, kvRead, updates); !valid || err != nil {
			return valid, err
		}
	}
	return true, nil
}

// validateKVRead performs mvcc check for a key read during transaction simulation.
// i.e., it checks whether a key/version combination is already updated in the statedb (by an already committed block)
// or in the updates (by a preceding valid transaction in the current block)
func (v *Validator) validateKVRead(ns string, kvRead *kvrwset.KVRead, updates *statedb.UpdateBatch) (bool, error) {
	if updates.Exists(ns, kvRead.Key) {
		return false, nil
	}
	versionedValue, err := v.db.GetState(ns, kvRead.Key)
	if err != nil {
		return false, nil
	}
	var committedVersion *version.Height
	if versionedValue != nil {
		committedVersion = versionedValue.Version
	}

	if !version.AreSame(committedVersion, rwsetutil.NewVersion(kvRead.Version)) {
		logger.Debugf("Version mismatch for key [%s:%s]. Committed version = [%s], Version in readSet [%s]",
			ns, kvRead.Key, committedVersion, kvRead.Version)
		return false, nil
	}
	return true, nil
}

func (v *Validator) validateRangeQueries(ns string, rangeQueriesInfo []*kvrwset.RangeQueryInfo, updates *statedb.UpdateBatch) (bool, error) {
	for _, rqi := range rangeQueriesInfo {
		if valid, err := v.validateRangeQuery(ns, rqi, updates); !valid || err != nil {
			return valid, err
		}
	}
	return true, nil
}

// validateRangeQuery performs a phatom read check i.e., it
// checks whether the results of the range query are still the same when executed on the
// statedb (latest state as of last committed block) + updates (prepared by the writes of preceding valid transactions
// in the current block and yet to be committed as part of group commit at the end of the validation of the block)
func (v *Validator) validateRangeQuery(ns string, rangeQueryInfo *kvrwset.RangeQueryInfo, updates *statedb.UpdateBatch) (bool, error) {
	logger.Debugf("validateRangeQuery: ns=%s, rangeQueryInfo=%s", ns, rangeQueryInfo)

	// If during simulation, the caller had not exhausted the iterator so
	// rangeQueryInfo.EndKey is not actual endKey given by the caller in the range query
	// but rather it is the last key seen by the caller and hence the combinedItr should include the endKey in the results.
	includeEndKey := !rangeQueryInfo.ItrExhausted

	combinedItr, err := newCombinedIterator(v.db, updates,
		ns, rangeQueryInfo.StartKey, rangeQueryInfo.EndKey, includeEndKey)
	if err != nil {
		return false, err
	}
	defer combinedItr.Close()
	var validator rangeQueryValidator
	if rangeQueryInfo.GetReadsMerkleHashes() != nil {
		logger.Debug(`Hashing results are present in the range query info hence, initiating hashing based validation`)
		validator = &rangeQueryHashValidator{}
	} else {
		logger.Debug(`Hashing results are not present in the range query info hence, initiating raw KVReads based validation`)
		validator = &rangeQueryResultsValidator{}
	}
	validator.init(rangeQueryInfo, combinedItr)
	return validator.validate()
}
