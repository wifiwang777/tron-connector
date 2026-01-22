package common

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestDecode(t *testing.T) {
	address := "TBWa43fqDWBnjAkqPqM3SkP9Aoo5spMZgb"
	decoded, err := DecodeAddress(address)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%x", []byte(decoded))
	// 4110e697f0d602ca3caadb50f571854f23872ff07f
}

func TestEncode(t *testing.T) {
	b, _ := hex.DecodeString("4110e697f0d602ca3caadb50f571854f23872ff07f")
	address := Address(b)
	t.Log(address.String())
	// TBWa43fqDWBnjAkqPqM3SkP9Aoo5spMZgb
}

func TestAddress(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
		return
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	t.Logf("%x", privateKeyBytes)

	publicKey := privateKey.PublicKey
	address := PublicKeyToAddress(publicKey)
	t.Log(address)
}

func TestPrivateKeyToAddress(t *testing.T) {
	privateKeyBytes, _ := hex.DecodeString("xxx")

	address := PrivateKeyToAddress(privateKeyBytes)
	t.Log(address)
}
