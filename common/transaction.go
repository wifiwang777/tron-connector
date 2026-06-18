package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wifiwang777/tron-connector/protos/api"
	"google.golang.org/protobuf/proto"
)

func NewTransaction(tx *api.TransactionExtention) (*Transaction, error) {
	if tx == nil {
		return nil, fmt.Errorf("tx is nil")
	}
	if tx.Result == nil {
		return nil, fmt.Errorf("tx result is nil")
	}
	if tx.Result.Code != api.Return_SUCCESS {
		return nil, fmt.Errorf("tx failed with code %s, message: %s", tx.Result.Code.String(), tx.Result.Message)
	}
	return &Transaction{TransactionExtention: tx}, nil
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
