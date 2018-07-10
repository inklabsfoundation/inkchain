/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package wallet

import (
	"encoding/hex"
	"math/big"
	"strings"
)

const (
	HashLength          = 32
	AddressLength       = 20
	HashStringLength    = 64
	AddressStringLength = 41
	PriKeyLength        = 32
	PriKeyStringLength  = 64
	ADDRESS_PREFIX      = "i"
	WALLET_NAMESPACE    = "ink"
	MAIN_BALANCE_NAME   = "INK"
)

type Hash [HashLength]byte
type Address [AddressLength]byte

var InkMinimumFee *big.Int
var SignPrivateKey string

type Account struct {
	Balance map[string]*big.Int `json:"balance"`
	Counter uint64              `json:"counter"`
}

type TxData struct {
	Sender      *Address `json:"from"`
	Recipient   *Address `json:"to"`
	BalanceType string   `json:"balanceType"`
	Amount      *big.Int `json:"amount"`
}

// type GenAccount
type Token struct {
	// token name
	Name string `json:"tokenName"`
	// total supply of the token
	totalSupply *big.Int `json:"totalSupply"`
	// initial address to issue
	Address string `json:"address"`
	// token status : Created, Delivered, Invalidate
	Status string `json:"status"`
	// token decimals
	Decimals int `json:"decimals"`
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
	if !strings.HasPrefix(b, ADDRESS_PREFIX) {
		return nil
	}
	c := strings.TrimLeft(b, ADDRESS_PREFIX)
	a := Address{}
	bytes, err := hex.DecodeString(strings.ToLower(c))
	if err != nil {
		return nil
	}
	a.SetBytes(bytes)
	return &a
}

func (a *Address) ToString() string {
	return string(ADDRESS_PREFIX + hex.EncodeToString(a[:]))
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
