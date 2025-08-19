package common

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestDecode(t *testing.T) {
	address := "TK3C8W8Ei6xk6EiRW4nMknPNoR7viQDC24"
	decoded, err := DecodeAddress(address)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%x", decoded)
	// 41637a117ac5aff17dc83e46b34d6c73f54d59416d
}

func TestEncode(t *testing.T) {
	b, _ := hex.DecodeString("41637a117ac5aff17dc83e46b34d6c73f54d59416d")
	encoded := EncodeAddress(b)
	t.Log(encoded)
	// TK3C8W8Ei6xk6EiRW4nMknPNoR7viQDC24
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
	address := PubkeyToAddress(publicKey)
	t.Log(address)
}
