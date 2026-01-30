package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/wifiwang777/tron-connector/common"
	"github.com/wifiwang777/tron-protocol/protos/api"
	"github.com/wifiwang777/tron-protocol/protos/core"
	"github.com/wifiwang777/tron-protocol/protos/core/contract"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn api.WalletClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn: api.NewWalletClient(conn),
	}, nil
}

// Base

func (c *Client) GetAccount(address common.Address) (*core.Account, error) {
	account := &core.Account{
		Address: address,
	}
	account, err := c.Conn.GetAccount(context.Background(), account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (c *Client) GetAccountResource(address common.Address) (*api.AccountResourceMessage, error) {
	account := &core.Account{
		Address: address,
	}
	return c.Conn.GetAccountResource(context.Background(), account)
}

func (c *Client) TransferTRX(from, to common.Address, amount *big.Int) (*api.TransactionExtention, error) {
	amount = common.MultiplyBy10Power(amount, 6)

	in := &contract.TransferContract{
		OwnerAddress: from,
		ToAddress:    to,
		Amount:       amount.Int64(),
	}
	tx, err := c.Conn.CreateTransaction2(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *Client) BroadcastTransaction(tx *core.Transaction) error {
	_, err := c.Conn.BroadcastTransaction(context.Background(), tx)
	return err
}

// TRC20

func (c *Client) Trc20Balance(contractAddress, account common.Address) (*big.Int, error) {
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

	result, err := c.Conn.TriggerConstantContract(context.Background(), ct)
	if err != nil {
		return nil, err
	}
	if len(result.ConstantResult) == 0 {
		return nil, fmt.Errorf("invalid result from contract")
	}

	return new(big.Int).SetBytes(result.ConstantResult[0]), nil
}

func (c *Client) Trc20Transfer(contractAddress, from, to common.Address, amount *big.Int, decimal int) (*api.TransactionExtention, error) {
	amount = common.MultiplyBy10Power(amount, decimal)
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

	result, err := c.Conn.TriggerContract(context.Background(), ct)
	return result, err
}

func (c *Client) Trc20Allowance(contractAddress, owner, spender common.Address) (*big.Int, error) {
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

	result, err := c.Conn.TriggerConstantContract(context.Background(), ct)
	if err != nil {
		return nil, err
	}
	if len(result.ConstantResult) == 0 {
		return nil, fmt.Errorf("invalid result from contract")
	}

	return new(big.Int).SetBytes(result.ConstantResult[0]), nil
}
