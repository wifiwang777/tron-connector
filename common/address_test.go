package common

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDecodeAddress(t *testing.T) {
	address := "TDksuv3bYk1TBWPXM4wRRFthb4X1w2J6cD"
	decoded, err := DecodeAddress(address)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%x", []byte(decoded))
	// 41298b93c0fc1e70e68b267ae93800caf0f502caf1
}

func TestBytesToAddress(t *testing.T) {
	bytes, _ := hex.DecodeString("41298b93c0fc1e70e68b267ae93800caf0f502caf1")
	address := Address(bytes)
	t.Logf("%s", address.String())
	// TDksuv3bYk1TBWPXM4wRRFthb4X1w2J6cD
}

func TestPublicKeyToAddress(t *testing.T) {
	publicKeyBytes := common.FromHex("029f23e018f77aec5a0b9a7795c3ae7da87baa00814ff0a375df30d1f6a8a75058")
	publicKey, err := crypto.DecompressPubkey(publicKeyBytes)
	if err != nil {
		t.Error(err)
		return
	}
	address := PublicKeyToAddress(*publicKey)
	t.Log(address)

}

func TestPrivateKeyToAddress(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("xxx")
	if err != nil {
		t.Error(err)
		return
	}
	address := PublicKeyToAddress(privateKey.PublicKey)
	t.Log(address.String())
}

func TestGeneratePrivateKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("private key: %x", crypto.FromECDSA(privateKey))
}
