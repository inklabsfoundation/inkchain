package ledgerutil

import (
	"github.com/golang/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/transutil"
	"github.com/inklabsfoundation/inkchain/protos/ledger"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/ctransutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/eftransutil"
)

type LedgerSet struct {
	TranSet      *transutil.TranSet
	TxRwSet      *rwsetutil.TxRwSet
	CrossTranSet *ctransutil.CrossTranSet
	EfTranSet    *eftransutil.EfTranSet
}

func (ledgerSet *LedgerSet) ToProtoBytes() ([]byte, error) {
	protoLedgerSet := &ledger.LedgerSet{}
	var err error
	protoLedgerSet.Transet, err = ledgerSet.TranSet.ToProtoBytes()
	if err != nil {
		return nil, err
	}

	protoLedgerSet.Txrwset, err = ledgerSet.TxRwSet.ToProtoBytes()
	if err != nil {
		return nil, err
	}
	protoLedgerSet.Crosstranset, err = ledgerSet.CrossTranSet.ToProtoBytes()
	if err != nil {
		return nil, err
	}
	protoLedgerSet.Eftranset, err = ledgerSet.EfTranSet.ToProtoBytes()
	if err != nil {
		return nil, err
	}
	protoLedgerSetBytes, err := proto.Marshal(protoLedgerSet)
	return protoLedgerSetBytes, nil
}

func (ledgerSet *LedgerSet) FromProtoBytes(protoBytes []byte) error {
	var err error
	protoLedgerSet := &ledger.LedgerSet{}
	if err = proto.Unmarshal(protoBytes, protoLedgerSet); err != nil {
		return err
	}
	if protoLedgerSet.Transet != nil {
		ledgerSet.TranSet = &transutil.TranSet{}
		err = ledgerSet.TranSet.FromProtoBytes(protoLedgerSet.Transet)
		if err != nil {
			return err
		}
		if ledgerSet.TranSet.From == "" {
			ledgerSet.TranSet = nil
		}
	} else {
		ledgerSet.TranSet = nil
	}

	if protoLedgerSet.Crosstranset != nil {
		ledgerSet.CrossTranSet = &ctransutil.CrossTranSet{}
		err = ledgerSet.CrossTranSet.FromProtoBytes(protoLedgerSet.Crosstranset)
		if err != nil {
			return err
		}
		if ledgerSet.CrossTranSet.TokenType == "" {
			ledgerSet.CrossTranSet = nil
		}
	} else {
		ledgerSet.CrossTranSet = nil
	}

	if protoLedgerSet.Eftranset != nil {
		ledgerSet.EfTranSet = &eftransutil.EfTranSet{}
		err = ledgerSet.EfTranSet.FromProtoBytes(protoLedgerSet.Eftranset)
		if err != nil {
			return err
		}
		if ledgerSet.EfTranSet.EfFrom == "" {
			ledgerSet.EfTranSet = nil
		}
	} else {
		ledgerSet.EfTranSet = nil
	}
	ledgerSet.TxRwSet = &rwsetutil.TxRwSet{}
	err = ledgerSet.TxRwSet.FromProtoBytes(protoLedgerSet.Txrwset)
	if err != nil {
		return err
	}
	return nil
}
