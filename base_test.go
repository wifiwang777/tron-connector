package tron_connector

import (
	"testing"

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
