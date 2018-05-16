package transutil

import (
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/core/ledger/kvledger/txmgmt/version"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet"
	"github.com/inklabsfoundation/inkchain/protos/ledger/transet/kvtranset"
)

type TranSet struct {
	From      string
	FromVer   *transet.Version
	KvTranSet *kvtranset.KVTranSet
}

type SenderCounter struct {
	Sender string
	Ink    *big.Int
}

// NsRwSet encapsulates 'kvrwset.KVRWSet' proto message for a specific name space (chaincode)

// ToProtoBytes constructs TxReadWriteSet proto message and serializes using protobuf Marshal
func (tranSet *TranSet) ToProtoBytes() ([]byte, error) {
	protoTranSet := &transet.TranSet{}
	protoTranSetBytes, err := proto.Marshal(tranSet.KvTranSet)
	if err != nil {
		return nil, err
	}
	protoTranSet.From = tranSet.From
	protoTranSet.FromVer = tranSet.FromVer
	protoTranSet.Transet = protoTranSetBytes
	protoTxBytes, err := proto.Marshal(protoTranSet)
	if err != nil {
		return nil, err
	}
	return protoTxBytes, nil
}

// FromProtoBytes deserializes protobytes into TxReadWriteSet proto message and populates 'TxRwSet'
func (tranSet *TranSet) FromProtoBytes(protoBytes []byte) error {
	protoTranSet := &transet.TranSet{}
	if err := proto.Unmarshal(protoBytes, protoTranSet); err != nil {
		return err
	}
	tranSet.From = protoTranSet.From
	tranSet.FromVer = protoTranSet.FromVer
	protoTranSetBytes := protoTranSet.Transet
	protoKvTranSet := &kvtranset.KVTranSet{}
	if err := proto.Unmarshal(protoTranSetBytes, protoKvTranSet); err != nil {
		return err
	}
	tranSet.KvTranSet = protoKvTranSet
	return nil
}

func NewTranSet(to string, balanceType string, amount []byte) *kvtranset.KVTrans {
	return &kvtranset.KVTrans{To: to, BalanceType: balanceType, Amount: amount}
}

func NewVersion(protoVersion *transet.Version) *version.Height {
	if protoVersion == nil {
		return nil
	}
	return version.NewHeight(protoVersion.BlockNum, protoVersion.TxNum)
}
