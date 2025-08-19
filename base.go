package tron_connector

import (
	"context"
	"github.com/wifiwang777/tron-protocol/protos/api"
	"github.com/wifiwang777/tron-protocol/protos/core"
	"github/wifiwang777/tron-connector/common"
	"google.golang.org/grpc"
)

type Tron struct {
	client api.WalletClient
}

func NewTron(conn *grpc.ClientConn) *Tron {
	client := api.NewWalletClient(conn)
	return &Tron{
		client: client,
	}
}

func (t *Tron) GetAccount(address string) (*core.Account, error) {
	account := new(core.Account)
	var err error
	account.Address, err = common.DecodeAddress(address)
	if err != nil {
		return nil, err
	}
	account, err = t.client.GetAccount(context.Background(), account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
