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

package lockbasedtxmgr

import (
	"errors"

	"github.com/inklabsfoundation/inkchain/common/util"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/ledgerutil"
	"github.com/inklabsfoundation/inkchain/core/wallet"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet/kvtranset"
	"github.com/inklabsfoundation/inkchain/protos/ledger/crosstranset/kvcrosstranset"
	"github.com/inklabsfoundation/inkchain/protos/ledger/eftranset/kveftranset"
)

// LockBasedTxSimulator is a transaction simulator used in `LockBasedTxMgr`
type lockBasedTxSimulator struct {
	lockBasedQueryExecutor
	ledgerBuilder *ledgerutil.LedgerSetBuilder
}

func newLockBasedTxSimulator(txmgr *LockBasedTxMgr) *lockBasedTxSimulator {
	ledgerBuilder := ledgerutil.NewLedgerBuilder()
	helper := &queryHelper{txmgr: txmgr, ledgerSetBuilder: ledgerBuilder}
	id := util.GenerateUUID()
	logger.Debugf("constructing new tx simulator [%s]", id)
	return &lockBasedTxSimulator{lockBasedQueryExecutor{helper, id}, ledgerBuilder}
}

// GetState implements method in interface `ledger.TxSimulator`
func (s *lockBasedTxSimulator) GetState(ns string, key string) ([]byte, error) {
	return s.helper.getState(ns, key)
}

// SetState implements method in interface `ledger.TxSimulator`
func (s *lockBasedTxSimulator) SetState(ns string, key string, value []byte) error {
	s.helper.checkDone()
	if err := s.helper.txmgr.db.ValidateKey(key); err != nil {
		return err
	}
	s.ledgerBuilder.RwSetBuilder.AddToWriteSet(ns, key, value)
	return nil
}

// Implementation for inkchain TxSimulator interface
func (s *lockBasedTxSimulator) Transfer(transet *kvtranset.KVTranSet) error {
	s.helper.checkDone()
	for _, tran := range transet.Trans {
		if err := s.helper.txmgr.db.ValidateKey(tran.To); err != nil {
			return err
		}
		s.ledgerBuilder.TranSetBuilder.AddToTranSet(tran.To, tran.BalanceType, tran.Amount)
	}
	return nil
}

// Implementation for inkchain TxSimulator interface
func (s *lockBasedTxSimulator) SetSender(sender string) error {
	if err := s.helper.txmgr.db.ValidateKey(sender); err != nil {
		return err
	}
	versionedValue, err := s.helper.txmgr.db.GetState(wallet.WALLET_NAMESPACE, sender)
	if err != nil {
		return err
	}
	_, verHeight := decomposeVersionedValue(versionedValue)
	var ver *transet.Version
	if verHeight != nil {
		ver = &transet.Version{verHeight.BlockNum, verHeight.TxNum}
	} else {
		ver = &transet.Version{0, 0}
	}
	s.ledgerBuilder.TranSetBuilder.SetSender(sender, ver)
	return nil
}

// Implementation for inkchain TxSimulator interface
func (s *lockBasedTxSimulator) CrossTransfer(transet *kvcrosstranset.KVCrossTranSet) error {
	s.helper.checkDone()
	for _, tran := range transet.Trans {
		if err := s.helper.txmgr.db.ValidateKey(tran.To); err != nil {
			return err
		}
		s.ledgerBuilder.CrossTranSetBuilder.AddToTranSet(tran.To, tran.Amount)
	}
	s.ledgerBuilder.CrossTranSetBuilder.SetTokenType(transet.BalanceType)
	return nil
}

// Implementation for inkchain TxSimulator interface
func (s *lockBasedTxSimulator) TransferExtractFee(efTranset *kveftranset.KVEfTranSet) error {
	s.helper.checkDone()
	for _, efTran := range efTranset.Eftrans {
		if err := s.helper.txmgr.db.ValidateKey(efTran.To); err != nil {
			return err
		}
		s.ledgerBuilder.EfTranSetBuilder.AddToTranSet(efTran.To, efTran.Amount)
	}
	return nil
}


// DeleteState implements method in interface `ledger.TxSimulator`
func (s *lockBasedTxSimulator) DeleteState(ns string, key string) error {
	return s.SetState(ns, key, nil)
}

// SetStateMultipleKeys implements method in interface `ledger.TxSimulator`
func (s *lockBasedTxSimulator) SetStateMultipleKeys(namespace string, kvs map[string][]byte) error {
	for k, v := range kvs {
		if err := s.SetState(namespace, k, v); err != nil {
			return err
		}
	}
	return nil
}

// GetTxSimulationResults implements method in interface `ledger.TxSimulator`
func (s *lockBasedTxSimulator) GetTxSimulationResults() ([]byte, error) {
	logger.Debugf("Simulation completed, getting simulation results")
	s.Done()
	if s.helper.err != nil {
		return nil, s.helper.err
	}
	return s.ledgerBuilder.GetLedgerSet().ToProtoBytes()
}

// ExecuteUpdate implements method in interface `ledger.TxSimulator`
func (s *lockBasedTxSimulator) ExecuteUpdate(query string) error {
	return errors.New("Not supported")
}
