package wallet

import (
	"encoding/hex"
	"math/big"
)

const (
	HashLength          = 32
	AddressLength       = 20
	HashStringLength    = 64
	AddressStringLength = 40
	PriKeyLength        = 32
	PriKeyStringLength  = 64
	WALLET_NAMESPACE    = "mtc"
	MTC_BALANCE_NAME    = "mtc"
)

type Hash [HashLength]byte
type Address [AddressLength]byte

type Account struct {
	Address *Address            `json:"address"`
	Balance map[string]*big.Int `json:"balance"`
	Counter uint64              `json:"counter"`
}

type TxData struct {
	Sender      *Address `json:"from"`
	Recipient   *Address `json:"to"`
	BalanceType string   `json:"balanceType"`
	Amount      *big.Int `json:"amount"`
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > AddressLength {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func BytesToAddress(b []byte) *Address {
	a := Address{}
	a.SetBytes(b)
	return &a
}

func (a *Address) ToBytes() []byte {
	return a[:]
}

func StringToAddress(b string) *Address {
	a := Address{}
	bytes, err := hex.DecodeString(b)
	if err != nil {
		return nil
	}
	a.SetBytes(bytes)
	return &a
}

func (a *Address) ToString() string {
	return string(hex.EncodeToString(a[:]))
}

func (a *Hash) SetBytes(b []byte) {
	if len(b) > HashLength {
		b = b[len(b)-HashLength:]
	}
	copy(a[HashLength-len(b):], b)
}

func BytesToHash(b []byte) *Hash {
	a := Hash{}
	a.SetBytes(b)
	return &a
}

func (a *Hash) ToBytes() []byte {
	return a[:]
}

func SignatureStringToBytes(sig string) ([]byte, error) {
	return hex.DecodeString(sig)
}

func SignatureBytesToString(sig []byte) string {
	return hex.EncodeToString(sig)
}
