package eftransutil

import (
	"github.com/golang/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/protos/ledger/eftranset/kveftranset"
	"github.com/inklabsfoundation/inkchain/protos/ledger/eftranset"
)

type EfTranSet struct {
	EfFrom      string
	KvEfTranSet *kveftranset.KVEfTranSet
}

func NewEfTranSet(to string, amount []byte) *kveftranset.KVEfTrans {
	return &kveftranset.KVEfTrans{To: to, Amount: amount}
}

// NsRwSet encapsulates 'kvrwset.KVRWSet' proto message for a specific name space (chaincode)

// ToProtoBytes constructs TxReadWriteSet proto message and serializes using protobuf Marshal
func (efTranSet *EfTranSet) ToProtoBytes() ([]byte, error) {
	protoTranSet := &eftranset.EfTranset{}
	protoTranSetBytes, err := proto.Marshal(efTranSet.KvEfTranSet)
	if err != nil {
		return nil, err
	}
	protoTranSet.From = efTranSet.EfFrom
	protoTranSet.Eftranset = protoTranSetBytes
	protoTxBytes, err := proto.Marshal(protoTranSet)
	if err != nil {
		return nil, err
	}
	return protoTxBytes, nil
}

// FromProtoBytes deserializes protobytes into TxReadWriteSet proto message and populates 'TxRwSet'
func (efTranSet *EfTranSet) FromProtoBytes(protoBytes []byte) error {
	protoTranSet := &eftranset.EfTranset{}
	if err := proto.Unmarshal(protoBytes, protoTranSet); err != nil {
		return err
	}
	efTranSet.EfFrom = protoTranSet.From
	protoTranSetBytes := protoTranSet.Eftranset
	protoKvTranSet := &kveftranset.KVEfTranSet{}
	if err := proto.Unmarshal(protoTranSetBytes, protoKvTranSet); err != nil {
		return err
	}
	efTranSet.KvEfTranSet = protoKvTranSet
	return nil
}
