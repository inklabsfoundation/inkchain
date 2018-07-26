/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"

	"crypto/sha256"

	"github.com/golang/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/common/crypto"
	"github.com/inklabsfoundation/inkchain/common/crypto/secp256k1"
	"github.com/inklabsfoundation/inkchain/common/crypto/sha3"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
)

//--------------------------------------
func ToECDSAPub(pub []byte) *ecdsa.PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(S256(), pub)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}
}

func NewECDSAPrivateKeyFromD(c elliptic.Curve, D *big.Int) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = D
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(D.Bytes())
	return priv
}

func SignJson(json []byte, priKey string) ([]byte, error) {
	hashT := sha256.Sum256(json)
	pri, err := HexToECDSA(priKey)
	if err != nil {
		return nil, nil
	}
	signature, err := crypto.Sign(hashT[:], pri)
	if err != nil {
		return nil, nil
	}
	return signature, nil
}

func GetInvokeHash(chaincodeSpec *pb.ChaincodeSpec, geneAlg string, senderSpec *pb.SenderSpec) ([]byte, error) {
	content := &pb.SignContent{ChaincodeSpec: chaincodeSpec, IdGenerationAlg: geneAlg, SenderSpec: senderSpec}
	protoBytes, err := proto.Marshal(content)
	if err != nil {
		return nil, err
	}
	hashT := sha256.Sum256(protoBytes)
	return hashT[:], nil
}

func SignInvoke(chaincodeSpec *pb.ChaincodeSpec, geneAlg string, senderSpec *pb.SenderSpec, priKey string) ([]byte, error) {
	content := &pb.SignContent{ChaincodeSpec: chaincodeSpec, IdGenerationAlg: geneAlg, SenderSpec: senderSpec}
	protoBytes, err := proto.Marshal(content)
	if err != nil {
		return nil, err
	}
	return SignJson(protoBytes, priKey)
}

func GetSenderFromSignature(hashT []byte, signature []byte) (*Address, error) {
	pub, err := crypto.Ecrecover(hashT, signature)
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %v", err)
	}
	return PubkeyToAddress(*ToECDSAPub(pub)), nil

}

func GetSenderPubKeyFromSignature(hashT []byte, signature []byte) (string, error) {
	pub, err := crypto.Ecrecover(hashT, signature)
	if err != nil {
		return "", fmt.Errorf("invalid signature: %v", err)
	}
	return hex.EncodeToString(pub[:]), nil

}

func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	if 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	return priv, nil
}

func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}
	return ToECDSA(b)
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

func PubkeyToAddress(p ecdsa.PublicKey) *Address {
	pubBytes := FromECDSAPub(&p)
	return BytesToAddress(Keccak256(pubBytes[1:])[12:])
}

func GetAddressFromPrikey(priKey string) (*Address, error) {
	ecdsa_prikey, err := HexToECDSA(priKey)
	if err != nil {
		return nil, err
	}
	return PubkeyToAddress(ecdsa_prikey.PublicKey), nil
}

func GetAddressHexFromPrikey(priKey string) (string, error) {
	ecdsa_prikey, err := HexToECDSA(priKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(PubkeyToAddress(ecdsa_prikey.PublicKey).ToBytes()), nil
}

func CheckAndGetSenderFromSignature(signature string, data []byte) (string, error) {
	signatureBytes, err := SignatureStringToBytes(signature)
	if err != nil {
		return "", err
	}
	hashT := sha256.Sum256(data)
	sender, err := GetSenderFromSignature(hashT[:], signatureBytes)
	if err != nil {
		return "", err
	}
	return sender.ToString(), nil
}

func HexToAddress(hexKey string) (*Address, error) {
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	return BytesToAddress(b), nil
}

func S256() elliptic.Curve {
	return secp256k1.S256()
}
