package common

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wifiwang777/tron-protocol/protos/core"
	"google.golang.org/protobuf/proto"
)

func SignTx(tx *core.Transaction, privateKey *ecdsa.PrivateKey) error {
	bytes, err := proto.Marshal(tx.RawData)
	if err != nil {
		return err
	}
	h256h := sha256.New()
	h256h.Write(bytes)
	hash := h256h.Sum(nil)

	sign, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return err
	}

	tx.Signature = append(tx.Signature, sign)
	return nil
}
