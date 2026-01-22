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

func (c *Client) GetAccount(address string) (*core.Account, error) {
	account := new(core.Account)
	var err error
	account.Address, err = common.DecodeAddress(address)
	if err != nil {
		return nil, err
	}
	account, err = c.Conn.GetAccount(context.Background(), account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (c *Client) BalanceOf(contractAddress, account common.Address) (*big.Int, error) {
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
