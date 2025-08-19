package tron_connector

import (
	"context"
	"fmt"
	"github.com/wifiwang777/tron-protocol/protos/api"
	"github.com/wifiwang777/tron-protocol/protos/core/contract"
	"github/wifiwang777/tron-connector/common"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type Trc20 struct {
	*Tron
}

func NewTrc20(t *Tron) *Trc20 {
	return &Trc20{
		Tron: t,
	}
}

func (t *Trc20) BalanceOf(contractAddress, account common.Address) (*big.Int, error) {
	methodSignature := []byte("balanceOf(address)")
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(methodSignature)
	methodId := keccak256.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(account, 32)...)

	ct := &contract.TriggerSmartContract{
		OwnerAddress:    account,
		ContractAddress: contractAddress,
		Data:            data,
	}

	result, err := t.client.TriggerConstantContract(context.Background(), ct)
	if err != nil {
		return nil, err
	}
	if len(result.ConstantResult) == 0 {
		return nil, fmt.Errorf("invalid result from contract")
	}

	return new(big.Int).SetBytes(result.ConstantResult[0]), nil
}

func (t *Trc20) Transfer(contractAddress, from, to common.Address, amount *big.Int) (*api.TransactionExtention, error) {
	methodSignature := []byte("transfer(address,uint256)")
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(methodSignature)
	methodId := keccak256.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(to, 32)...)
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)
	ct := &contract.TriggerSmartContract{
		OwnerAddress:    from,
		ContractAddress: contractAddress,
		Data:            data,
	}

	result, err := t.client.TriggerContract(context.Background(), ct)
	return result, err
}

func (t *Trc20) Allowance(contractAddress, owner, spender common.Address) (*big.Int, error) {
	transferFnSignature := []byte("allowance(address,address)")
	erc20hash := sha3.NewLegacyKeccak256()
	erc20hash.Write(transferFnSignature)
	methodId := erc20hash.Sum(nil)[:4]

	paddedOwner := common.LeftPadBytes(owner, 32)
	paddedSpender := common.LeftPadBytes(spender, 32)
	var data []byte
	data = append(data, methodId...)
	data = append(data, paddedOwner...)
	data = append(data, paddedSpender...)

	ct := &contract.TriggerSmartContract{
		OwnerAddress:    owner,
		ContractAddress: contractAddress,
		Data:            data,
	}

	result, err := t.client.TriggerConstantContract(context.Background(), ct)
	if err != nil {
		return nil, err
	}
	if len(result.ConstantResult) == 0 {
		return nil, fmt.Errorf("invalid result from contract")
	}

	return new(big.Int).SetBytes(result.ConstantResult[0]), nil
}
