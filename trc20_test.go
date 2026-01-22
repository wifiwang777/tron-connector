package tron_connector

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

func TestBalanceOf(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)
	trc20 := NewTrc20(tron)

	contractAddress, err := common.DecodeAddress(USDTContractAddressNile)
	if err != nil {
		t.Error(err)
		return
	}

	account, err := common.DecodeAddress("TK3C8W8Ei6xk6EiRW4nMknPNoR7viQDC24")
	if err != nil {
		t.Error(err)
		return
	}

	balance, err := trc20.BalanceOf(contractAddress, account)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(balance.String())
}

func TestTransfer(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)
	trc20 := NewTrc20(tron)

	contractAddress, err := common.DecodeAddress(USDTContractAddressNile)
	if err != nil {
		t.Error(err)
		return
	}

	hexKey := "YOUR_PRIVATE_KEY"
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		t.Error(err)
		return
	}

	from := common.PublicKeyToAddress(privateKey.PublicKey)

	from, err := common.DecodeAddress("xx")
	if err != nil {
		t.Error(err)
		return
	}

	to, err := common.DecodeAddress("xx")

	amount := new(big.Int).Mul(big.NewInt(100), big.NewInt(1000000))

	tx, err := trc20.Transfer(contractAddress, from, to, amount)
	if err != nil {
		t.Error(err)
		return
	}

	err = common.SignTx(tx.Transaction, privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = trc20.client.BroadcastTransaction(context.Background(), tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Transaction ID: %x", tx.Txid)
}

func TestAllowance(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)
	trc20 := NewTrc20(tron)

	contractAddress, err := common.DecodeAddress(USDTContractAddressNile)
	if err != nil {
		t.Error(err)
		return
	}

	owner, err := common.DecodeAddress("TZ6CAfoc48NCAEvSkDD2kPou3dgKuhpnEf")
	if err != nil {
		t.Error(err)
		return
	}

	spender, err := common.DecodeAddress("TK8DFcRXsECeeN9fsHkNHT6wmkjLrnaDwi")
	if err != nil {
		t.Error(err)
		return
	}

	allowance, err := trc20.Allowance(contractAddress, owner, spender)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(allowance.String())
}

func TestGetContractInfo(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointMainnet, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)

	address, _ := common.DecodeAddress("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")

	in := &api.BytesMessage{
		Value: address,
	}
	res, err := tron.client.GetContractInfo(context.Background(), in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.SmartContract)
	t.Log(res.ContractState)
}

func TestGetEnergyPrice(t *testing.T) {
	conn, err := grpc.NewClient(GrpcEndpointNile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	tron := NewTron(conn)
	res, err := tron.client.GetEnergyPrices(context.Background(), &api.EmptyMessage{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Prices)
	//2030aa598e03cc4b3c88dfb03965211a44f5482ee40f92afa9993caf1237e5bd
}
