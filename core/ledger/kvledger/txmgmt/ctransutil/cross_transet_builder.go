
/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package ctransutil

import (
	"fmt"

	"github.com/inklabsfoundation/inkchain/common/flogging"
	"github.com/inklabsfoundation/inkchain/core/ledger/util"
	"github.com/inklabsfoundation/inkchain/protos/ledger/crosstranset"
	"github.com/inklabsfoundation/inkchain/protos/ledger/crosstranset/kvcrosstranset"
)

type CrossTranSetBuilder struct {
	tokenAccount string
	tokenType    string
	fromVer      *crosstranset.Version
	toSet        map[string]*kvcrosstranset.KVCrossTrans
}

func NewCrossTranSetBuilder() *CrossTranSetBuilder {
	return &CrossTranSetBuilder{"", "", nil, make(map[string]*kvcrosstranset.KVCrossTrans)}
}

var logger = flogging.MustGetLogger("crosstransutil")

func (builder *CrossTranSetBuilder) AddToTranSet(to string,  amount []byte) {
	builder.toSet[to] = NewCrossTranSet(to, amount)
}

func (builder *CrossTranSetBuilder) SetTokenType(tokenType string) {
	if tokenType == "" {
		fmt.Println("fatal: empty token type")
		return
	}
	builder.tokenType = tokenType
}

func (builder *CrossTranSetBuilder) GetTranSet() *CrossTranSet {
	tranSet := &CrossTranSet{}
	//from
	tranSet.TokenType = builder.tokenType
	//transferSet
	var trans []*kvcrosstranset.KVCrossTrans
	sortedTransKeys := util.GetSortedKeys(builder.toSet)
	for _, key := range sortedTransKeys {
		trans = append(trans, builder.toSet[key])
	}
	kvTrans := &kvcrosstranset.KVCrossTranSet{Trans: trans}
	tranSet.KvTranSet = kvTrans
	return tranSet
}