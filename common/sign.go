package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wifiwang777/tron-protocol/protos/api"
	"google.golang.org/protobuf/proto"
)

func SignTx(tx *api.TransactionExtention, privateKey *ecdsa.PrivateKey) (*api.TransactionExtention, error) {
	bytes, err := proto.Marshal(tx.Transaction.RawData)
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

	tx.Transaction.Signature = append(tx.Transaction.Signature, sign)
	return tx, nil
}
