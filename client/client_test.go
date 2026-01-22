package client

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wifiwang777/tron-connector/common"
	"github.com/wifiwang777/tron-protocol/protos/api"
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
	client := &Client{
		Conn: api.NewWalletClient(conn),
	}
	address, _ := common.DecodeAddress("TF9J64qtrNfSChQ217zcmd8HPMsHbUhK6Y")
	account, err := client.GetAccount(address)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(account.Balance)
}

func TestGetAccountResource(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	client := &Client{
		Conn: api.NewWalletClient(conn),
	}
	address, _ := common.DecodeAddress("TF9J64qtrNfSChQ217zcmd8HPMsHbUhK6Y")

	resource, err := client.GetAccountResource(address)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resource.EnergyUsed)
	t.Log(resource.EnergyLimit)

	currentEnergy := resource.EnergyLimit - resource.EnergyUsed
	t.Log("current energy:", currentEnergy)
}

func TestGetEnergyPrice(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	client := &Client{
		Conn: api.NewWalletClient(conn),
	}
	res, err := client.Conn.GetEnergyPrices(context.Background(), &api.EmptyMessage{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Prices)
}

func TestGetContractInfo(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	client := &Client{
		Conn: api.NewWalletClient(conn),
	}
	contract, _ := common.DecodeAddress(USDTContractAddressNile)

	in := &api.BytesMessage{
		Value: contract,
	}
	res, err := client.Conn.GetContractInfo(context.Background(), in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.SmartContract)
	t.Log(res.ContractState)
}

func TestTrc20Balance(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	client := &Client{
		Conn: api.NewWalletClient(conn),
	}
	contractAddress, err := common.DecodeAddress(USDTContractAddressNile)
	if err != nil {
		t.Error(err)
		return
	}

	address, err := common.DecodeAddress("TF9J64qtrNfSChQ217zcmd8HPMsHbUhK6Y")
	if err != nil {
		t.Error(err)
		return
	}

	balance, err := client.Trc20Balance(contractAddress, address)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Balance: %s", balance)
}

func TestTrc20Transfer(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	client := &Client{
		Conn: api.NewWalletClient(conn),
	}

	contractAddress, err := common.DecodeAddress(USDTContractAddressNile)
	if err != nil {
		t.Error(err)
		return
	}

	hexKey := "YOUR_PRIVATE_KEY"
	toAddress := "YOUR_RECEIVER_ADDRESS"

	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		t.Error(err)
		return
	}

	from := common.PublicKeyToAddress(privateKey.PublicKey)
	to, err := common.DecodeAddress(toAddress)
	if err != nil {
		t.Error(err)
	}

	amount := new(big.Int).Mul(big.NewInt(100), big.NewInt(1000000))

	tx, err := client.Trc20Transfer(contractAddress, from, to, amount)
	if err != nil {
		t.Error(err)
		return
	}

	err = common.SignTx(tx.Transaction, privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = client.Conn.BroadcastTransaction(context.Background(), tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Transaction ID: %x", tx.Txid)
}
