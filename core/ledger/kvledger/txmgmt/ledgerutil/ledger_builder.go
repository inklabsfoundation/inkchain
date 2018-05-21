package ledgerutil

import (
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/transutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/ctransutil"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/eftransutil"
)

type LedgerSetBuilder struct {
	TranSetBuilder      *transutil.TranSetBuilder
	RwSetBuilder        *rwsetutil.RWSetBuilder
	CrossTranSetBuilder *ctransutil.CrossTranSetBuilder
	EfTranSetBuilder    *eftransutil.EfTranSetBuilder
}

func (builder *LedgerSetBuilder) GetLedgerSet() *LedgerSet {
	ledgerSet := &LedgerSet{}
	ledgerSet.TranSet = builder.TranSetBuilder.GetTranSet()
	ledgerSet.TxRwSet = builder.RwSetBuilder.GetTxReadWriteSet()
	ledgerSet.CrossTranSet = builder.CrossTranSetBuilder.GetTranSet()
	ledgerSet.EfTranSet = builder.EfTranSetBuilder.GetEfTranSet()
	return ledgerSet
}

func NewLedgerBuilder() *LedgerSetBuilder {
	ledgerSetBuilder := &LedgerSetBuilder{}
	ledgerSetBuilder.RwSetBuilder = rwsetutil.NewRWSetBuilder()
	ledgerSetBuilder.TranSetBuilder = transutil.NewTranSetBuilder()
	ledgerSetBuilder.CrossTranSetBuilder = ctransutil.NewCrossTranSetBuilder()
	ledgerSetBuilder.EfTranSetBuilder = eftransutil.NewEfTranSetBuilder()
	return ledgerSetBuilder

}
