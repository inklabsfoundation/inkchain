package ledgerutil

import (
	"github.com/inkchain/inkchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/inkchain/inkchain/core/ledger/kvledger/txmgmt/transutil"
)

type LedgerSetBuilder struct {
	TranSetBuilder *transutil.TranSetBuilder
	RwSetBuilder   *rwsetutil.RWSetBuilder
}

func (builder *LedgerSetBuilder) GetLedgerSet() *LedgerSet {
	ledgerSet := &LedgerSet{}
	ledgerSet.TranSet = builder.TranSetBuilder.GetTranSet()
	ledgerSet.TxRwSet = builder.RwSetBuilder.GetTxReadWriteSet()
	return ledgerSet
}

func NewLedgerBuilder() *LedgerSetBuilder {
	ledgerSetBuilder := &LedgerSetBuilder{}
	ledgerSetBuilder.RwSetBuilder = rwsetutil.NewRWSetBuilder()
	ledgerSetBuilder.TranSetBuilder = transutil.NewTranSetBuilder()
	return ledgerSetBuilder

}
