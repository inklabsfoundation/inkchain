/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package ctransutil

import (
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/version"
	"github.com/inklabsfoundation/inkchain/protos/ledger/crosstranset"
	"github.com/inklabsfoundation/inkchain/protos/ledger/crosstranset/kvcrosstranset"
)

type CrossTranSet struct {
	TokenAddr string
	TokenType string
	KvTranSet *kvcrosstranset.KVCrossTranSet
}

type SenderCounter struct {
	Sender  string
	Counter uint64
	Ink     *big.Int
}

// NsRwSet encapsulates 'kvrwset.KVRWSet' proto message for a specific name space (chaincode)

// ToProtoBytes constructs TxReadWriteSet proto message and serializes using protobuf Marshal
func (tranSet *CrossTranSet) ToProtoBytes() ([]byte, error) {
	protoTranSet := &crosstranset.CrossTranSet{}
	protoTranSetBytes, err := proto.Marshal(tranSet.KvTranSet)
	if err != nil {
		return nil, err
	}
	protoTranSet.TokenType = tranSet.TokenType
	protoTranSet.Ctranset = protoTranSetBytes
	protoTxBytes, err := proto.Marshal(protoTranSet)
	if err != nil {
		return nil, err
	}
	return protoTxBytes, nil
}

// FromProtoBytes deserializes protobytes into TxReadWriteSet proto message and populates 'TxRwSet'
func (tranSet *CrossTranSet) FromProtoBytes(protoBytes []byte) error {
	protoTranSet := &crosstranset.CrossTranSet{}
	if err := proto.Unmarshal(protoBytes, protoTranSet); err != nil {
		return err
	}
	tranSet.TokenAddr = protoTranSet.TokenAddr
	tranSet.TokenType = protoTranSet.TokenType
	protoTranSetBytes := protoTranSet.Ctranset
	protoKvTranSet := &kvcrosstranset.KVCrossTranSet{}
	if err := proto.Unmarshal(protoTranSetBytes, protoKvTranSet); err != nil {
		return err
	}
	tranSet.KvTranSet = protoKvTranSet
	return nil
}

func NewCrossTranSet(to string, amount []byte) *kvcrosstranset.KVCrossTrans {
	return &kvcrosstranset.KVCrossTrans{To: to, Amount: amount}
}

func NewVersion(protoVersion *crosstranset.Version) *version.Height {
	if protoVersion == nil {
		return nil
	}
	return version.NewHeight(protoVersion.BlockNum, protoVersion.TxNum)
}
