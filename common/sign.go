package common

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wifiwang777/tron-protocol/protos/core"
	"google.golang.org/protobuf/proto"
)

func GenerateSignature(tx *core.Transaction, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	bytes, err := proto.Marshal(tx.RawData)
	if err != nil {
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(bytes)
	hash := h256h.Sum(nil)

	sign, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}

	return sign, nil
}
