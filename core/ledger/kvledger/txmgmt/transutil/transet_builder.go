package transutil

import (
	"fmt"

	"github.com/inklabsfoundation/inkchain/common/flogging"
	"github.com/inklabsfoundation/inkchain/core/ledger/util"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet/kvtranset"
)

type TranSetBuilder struct {
	from    string
	fromVer *transet.Version
	toSet   map[string]*kvtranset.KVTrans
}

func NewTranSetBuilder() *TranSetBuilder {
	return &TranSetBuilder{"", nil, make(map[string]*kvtranset.KVTrans)}
}

var logger = flogging.MustGetLogger("transutil")

func (builder *TranSetBuilder) AddToTranSet(to string, balanceType string, amount []byte) {
	builder.toSet[to] = NewTranSet(to, balanceType, amount)
}

func (builder *TranSetBuilder) SetSender(sender string, version *transet.Version) {
	if sender == "" {
		fmt.Println("fatal: empty sender")
		return
	}
	if builder.from == "" {
		builder.from = sender
	} else if builder.from != sender {
		panic("fatal: multiple sender in one transaction!")
		return
	}
	builder.fromVer = version
}

func (builder *TranSetBuilder) GetTranSet() *TranSet {
	tranSet := &TranSet{}
	//from
	tranSet.From = builder.from
	tranSet.FromVer = builder.fromVer
	//transferSet
	var trans []*kvtranset.KVTrans
	sortedTransKeys := util.GetSortedKeys(builder.toSet)
	for _, key := range sortedTransKeys {
		trans = append(trans, builder.toSet[key])
	}
	kvTrans := &kvtranset.KVTranSet{Trans: trans}
	tranSet.KvTranSet = kvTrans
	return tranSet
}
