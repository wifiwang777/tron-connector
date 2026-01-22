package tron_connector

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/wifiwang777/tron-connector/common"
	"github.com/wifiwang777/tron-protocol/protos/api"
	"github.com/wifiwang777/tron-protocol/protos/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	GrpcEndpointMainnet = "grpc.trongrid.io:50051"
	GrpcEndpointNile    = "grpc.nile.trongrid.io:50051"

	USDTContractAddressMainnet = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	USDTContractAddressNile    = "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj"
)

func TestGetAccount(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)
	account, err := tron.GetAccount("TK3C8W8Ei6xk6EiRW4nMknPNoR7viQDC24")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(account.Balance)
}

func TestGetTransactionInfo(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)

	in := new(api.BytesMessage)
	in.Value, err = hex.DecodeString("20f59b5b263e63e0ff0110052c22bb9f3b6700ac8fe86bab773111b5bc0230b0")
	if err != nil {
		t.Error(err)
		return
	}

	transactionInfo, err := tron.client.GetTransactionInfoById(context.Background(), in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(transactionInfo.BlockNumber)
	t.Log(transactionInfo.BlockTimeStamp)
}

func TestGetTransactionInfoByBlockNum(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)

	numMessage := new(api.NumberMessage)
	numMessage.Num = 62109747

	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 10e6)
	block, err := tron.client.GetTransactionInfoByBlockNum(context.Background(), numMessage, maxSizeOption)
	if err != nil {
		t.Error(err)
		return
	}
	for _, transactionInfo := range block.TransactionInfo {
		t.Log(transactionInfo.BlockNumber)
		t.Log(transactionInfo.BlockTimeStamp)
		return
	}
}

func TestGetAccountResource(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointMainnet, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)
	account := new(core.Account)
	account.Address, err = common.DecodeAddress("TVt7oQuLnHZz252eaDFLbh66zHDGgksoSY")
	if err != nil {
		t.Error(err)
		return
	}
	resource, err := tron.client.GetAccountResource(context.Background(), account)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resource.EnergyUsed)
	t.Log(resource.EnergyLimit)
	t.Log(resource.TotalEnergyLimit)

	currentEnergy := resource.EnergyLimit - resource.EnergyUsed
	t.Log("current energy:", currentEnergy)
}
