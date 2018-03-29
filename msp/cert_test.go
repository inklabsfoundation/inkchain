/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package msp

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/inklabsfoundation/inkchain/bccsp/sw"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeCertWithRSA(t *testing.T) {
	cert := &x509.Certificate{}
	cert.SignatureAlgorithm = x509.MD2WithRSA
	result := isECDSASignedCert(cert)
	assert.False(t, result)

	cert.SignatureAlgorithm = x509.ECDSAWithSHA512
	result = isECDSASignedCert(cert)
	assert.True(t, result)
}

func TestSanitizeCertInvalidInput(t *testing.T) {
	_, err := sanitizeECDSASignedCert(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Certificate must be different from nil.")

	_, err = sanitizeECDSASignedCert(&x509.Certificate{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Parent certificate must be different from nil")

	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	cert := &x509.Certificate{}
	cert.PublicKey = &k.PublicKey
	sigma, err := sw.MarshalECDSASignature(big.NewInt(1), elliptic.P256().Params().N)
	assert.NoError(t, err)
	cert.Signature = sigma
	cert.PublicKeyAlgorithm = x509.ECDSA
	cert.Raw = []byte{0, 1}
	_, err = sanitizeECDSASignedCert(cert, cert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "asn1: structure error: tags don't match")
}

func TestSanitizeCert(t *testing.T) {
	var k *ecdsa.PrivateKey
	var cert *x509.Certificate
	for {
		k, cert = generateSelfSignedCert(t, time.Now())

		_, s, err := sw.UnmarshalECDSASignature(cert.Signature)
		assert.NoError(t, err)

		lowS, err := sw.IsLowS(&k.PublicKey, s)
		assert.NoError(t, err)

		if !lowS {
			break
		}
	}

	sanitizedCert, err := sanitizeECDSASignedCert(cert, cert)
	assert.NoError(t, err)
	assert.NotEqual(t, cert.Signature, sanitizedCert.Signature)

	_, s, err := sw.UnmarshalECDSASignature(sanitizedCert.Signature)
	assert.NoError(t, err)

	lowS, err := sw.IsLowS(&k.PublicKey, s)
	assert.NoError(t, err)
	assert.True(t, lowS)
}

func TestCertExpiration(t *testing.T) {
	msp := &bccspmsp{}
	msp.opts = &x509.VerifyOptions{}
	msp.opts.DNSName = "test.example.com"

	// Certificate is in the future
	_, cert := generateSelfSignedCert(t, time.Now().Add(24*time.Hour))
	msp.opts.Roots = x509.NewCertPool()
	msp.opts.Roots.AddCert(cert)
	_, err := msp.getUniqueValidationChain(cert, msp.getValidityOptsForCert(cert))
	assert.NoError(t, err)

	// Certificate is in the past
	_, cert = generateSelfSignedCert(t, time.Now().Add(-24*time.Hour))
	msp.opts.Roots = x509.NewCertPool()
	msp.opts.Roots.AddCert(cert)
	_, err = msp.getUniqueValidationChain(cert, msp.getValidityOptsForCert(cert))
	assert.NoError(t, err)

	// Certificate is in the middle
	_, cert = generateSelfSignedCert(t, time.Now())
	msp.opts.Roots = x509.NewCertPool()
	msp.opts.Roots.AddCert(cert)
	_, err = msp.getUniqueValidationChain(cert, msp.getValidityOptsForCert(cert))
	assert.NoError(t, err)
}

func generateSelfSignedCert(t *testing.T, now time.Time) (*ecdsa.PrivateKey, *x509.Certificate) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	// Generate a self-signed certificate
	testExtKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	testUnknownExtKeyUsage := []asn1.ObjectIdentifier{[]int{1, 2, 3}, []int{2, 59, 1}}
	extraExtensionData := []byte("extra extension")
	commonName := "test.example.com"
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Σ Acme Co"},
			Country:      []string{"US"},
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  []int{2, 5, 4, 42},
					Value: "Gopher",
				},
				// This should override the Country, above.
				{
					Type:  []int{2, 5, 4, 6},
					Value: "NL",
				},
			},
		},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(1 * time.Hour),
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		SubjectKeyId:          []byte{1, 2, 3, 4},
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           testExtKeyUsage,
		UnknownExtKeyUsage:    testUnknownExtKeyUsage,
		BasicConstraintsValid: true,
		IsCA:                  true,
		OCSPServer:            []string{"http://ocurrentBCCSP.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},
		DNSNames:              []string{"test.example.com"},
		EmailAddresses:        []string{"gopher@golang.org"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1).To4(), net.ParseIP("2001:4860:0:2001::68")},
		PolicyIdentifiers:     []asn1.ObjectIdentifier{[]int{1, 2, 3}},
		PermittedDNSDomains:   []string{".example.com", "example.com"},
		CRLDistributionPoints: []string{"http://crl1.example.com/ca1.crl", "http://crl2.example.com/ca1.crl"},
		ExtraExtensions: []pkix.Extension{
			{
				Id:    []int{1, 2, 3, 4},
				Value: extraExtensionData,
			},
		},
	}
	certRaw, err := x509.CreateCertificate(rand.Reader, &template, &template, &k.PublicKey, k)
	assert.NoError(t, err)

	cert, err := x509.ParseCertificate(certRaw)
	assert.NoError(t, err)

	return k, cert
}
