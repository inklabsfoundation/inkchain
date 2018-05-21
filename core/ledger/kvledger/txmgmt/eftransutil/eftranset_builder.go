package eftransutil

import (
	"github.com/inklabsfoundation/inkchain/protos/ledger/eftranset/kveftranset"
	"github.com/inklabsfoundation/inkchain/core/ledger/util"
	"github.com/inklabsfoundation/inkchain/common/flogging"
	"fmt"
)

type EfTranSetBuilder struct {
	efFrom string
	toSet  map[string]*kveftranset.KVEfTrans
}

func NewEfTranSetBuilder() *EfTranSetBuilder {
	return &EfTranSetBuilder{"", make(map[string]*kveftranset.KVEfTrans)}
}

var logger = flogging.MustGetLogger("eftransutil")

func (builder *EfTranSetBuilder) AddToTranSet(to string, amount []byte) {
	builder.toSet[to] = NewEfTranSet(to, amount)
}

func (builder *EfTranSetBuilder) SetSender(sender string) {
	if sender == "" {
		fmt.Println("fatal: empty sender")
		return
	}
	if builder.efFrom == "" {
		builder.efFrom = sender
	} else if builder.efFrom != sender {
		panic("fatal: multiple sender in one transaction!")
		return
	}
}


func (builder *EfTranSetBuilder) GetEfTranSet() *EfTranSet {
	efTranSet := &EfTranSet{}
	efTranSet.EfFrom = builder.efFrom
	//transferSet
	var efTrans []*kveftranset.KVEfTrans
	sortedTransKeys := util.GetSortedKeys(builder.toSet)
	for _, key := range sortedTransKeys {
		efTrans = append(efTrans, builder.toSet[key])
	}
	kvefTrans := &kveftranset.KVEfTranSet{efTrans}
	efTranSet.KvEfTranSet = kvefTrans
	return efTranSet
}
