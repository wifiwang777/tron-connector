package common

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wifiwang777/tron-connector/protos/api"
	"google.golang.org/protobuf/proto"
)

func NewTransaction(tx *api.TransactionExtention) *Transaction {
	return &Transaction{TransactionExtention: tx}
}

type Transaction struct {
	*api.TransactionExtention
}

func (tx *Transaction) SignWithPrivateKey(privateKey *ecdsa.PrivateKey) error {
	digest, err := tx.CalculateDigestHash()
	if err != nil {
		return err
	}

	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return err
	}

	tx.WithdrawSignature(signature)
	return nil
}

func (tx *Transaction) CalculateDigestHash() ([]byte, error) {
	bytes, err := proto.Marshal(tx.Transaction.RawData)
	if err != nil {
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(bytes)
	hash := h256h.Sum(nil)
	return hash, nil
}

func (tx *Transaction) WithdrawSignature(signature []byte) {
	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)
}
