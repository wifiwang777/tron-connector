package client

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/wifiwang777/tron-connector/common"
	"github.com/wifiwang777/tron-connector/protos/api"
	"github.com/wifiwang777/tron-connector/protos/core"
)

const (
	GrpcEndpointMainnet = "grpc.trongrid.io:50051"
	GrpcEndpointNile    = "grpc.nile.trongrid.io:50051"

	USDTContractAddressMainnet = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	USDTContractAddressNile    = "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj"
)

var (
	Endpoint                  = GrpcEndpointNile
	TokenAddress              = USDTContractAddressNile
	ApiHeaderName             = "TRON-PRO-API-KEY"
	ApiKey                    = "your api key"
	TokenDecimals             = 6
	Trc20TransferFeeLimit     = int64(2686000)
	Trc20ApproveFeeLimit      = int64(2250600)
	Trc20TransferFromFeeLimit = int64(2076000)

	MainAddress     = "your main address"
	SlaveAddress    = "your slave address"
	ReceiverAddress = "your receiver address"

	MainPrivateKey  = "your hex private key"
	SlavePrivateKey = "your hex private key"
)

func getClient() (*Client, error) {
	return NewClient(
		Endpoint,
		//WithAPIKey(ApiHeaderName, ApiKey),
	)
}

func TestGetNodeInfo(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	nodeInfo, err := client.Conn.GetNodeInfo(context.Background(), new(api.EmptyMessage))
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("nodeInfo: %v", nodeInfo)
}

func TestGetAccount(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
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
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	address, _ := common.DecodeAddress("TF9J64qtrNfSChQ217zcmd8HPMsHbUhK6Y")

	resource, err := client.GetAccountResource(address)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("EnergyUsed", resource.EnergyUsed)
	t.Log("EnergyLimit", resource.EnergyLimit)

	currentEnergy := resource.EnergyLimit - resource.EnergyUsed
	t.Log("current energy:", currentEnergy)

	t.Log("FreeNetLimit", resource.FreeNetLimit)
	t.Log("FreeNetUsed", resource.FreeNetUsed)
	t.Log("current bandwidth: ", resource.FreeNetLimit-resource.FreeNetUsed)
}

func TestCreateAccount(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	ownerAddress, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}
	accountAddress, err := common.DecodeAddress("THFvKmA3jRDb6xpKbup37d7tESQ3zjQFBU")
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := crypto.HexToECDSA(MainPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}

	tx, err := client.CreateAccount(ownerAddress, accountAddress, core.AccountType_Normal)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	txId := tx.Txid
	t.Logf("txId: %x", txId)
}

func TestDelegateResource(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	from, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}
	to, err := common.DecodeAddress("THFvKmA3jRDb6xpKbup37d7tESQ3zjQFBU")
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := crypto.HexToECDSA(MainPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}

	tx, err := client.DelegateResource(from, to, core.ResourceCode_ENERGY, 100000000, false, 0)
	if err != nil {
		t.Error(err)
		return
	}

	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	txId := tx.Txid
	t.Logf("txId: %x", txId)
}

func TestUnDelegateResource(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	from, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}
	to, err := common.DecodeAddress("THFvKmA3jRDb6xpKbup37d7tESQ3zjQFBU")
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := crypto.HexToECDSA(MainPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}

	tx, err := client.UnDelegateResource(from, to, core.ResourceCode_ENERGY, 100000000)
	if err != nil {
		t.Error(err)
		return
	}

	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	txId := tx.Txid
	t.Logf("txId: %x", txId)
}

func TestTransferTRX(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	hexKey := MainPrivateKey
	toAddress := SlaveAddress

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

	amountStr := "100.10000000"
	amount, _ := decimal.NewFromString(amountStr)
	tx, err := client.Transfer(from, to, amount)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	txId := tx.Txid
	t.Logf("txId: %x", txId)
}

func TestGetCurrentEnergyPrice(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	res, err := client.GetCurrentEnergyPrice()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestGetCurrentBandwidthPrice(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	res, err := client.GetCurrentBandwidthPrice()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestGetContractInfo(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	contract, _ := common.DecodeAddress(TokenAddress)

	in := &api.BytesMessage{
		Value: contract,
	}
	res, err := client.Conn.GetContractInfo(context.Background(), in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(common.Address(res.SmartContract.OriginAddress).String())
	t.Log(res.SmartContract.OriginEnergyLimit)
	t.Log(res.ContractState)
}

func TestGetChainParameters(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	chainParameters, err := client.GetChainParameterMap()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("chain parameters: %+v", chainParameters)
}

func TestGetTrc20Balance(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	contractAddress, err := common.DecodeAddress(TokenAddress)
	if err != nil {
		t.Error(err)
		return
	}

	account, err := common.DecodeAddress("TF9J64qtrNfSChQ217zcmd8HPMsHbUhK6Y")
	if err != nil {
		t.Error(err)
		return
	}

	balance, err := client.GetTrc20Balance(account, contractAddress, int32(TokenDecimals))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Balance: %s", balance)
}

func TestGetTrc20TransferEnergyCost(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	contractAddress, err := common.DecodeAddress(TokenAddress)
	if err != nil {
		t.Error(err)
		return
	}

	from, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}
	to, err := common.DecodeAddress(ReceiverAddress)
	if err != nil {
		t.Error(err)
		return
	}

	amountStr := "101"
	//amountStr := "1010000000000"

	amount, _ := decimal.NewFromString(amountStr)
	decimals := int32(TokenDecimals)

	tcs, err := client.GenerateTrc20TransferTrigger(contractAddress, from, to, amount, decimals)
	if err != nil {
		t.Error(err)
		return
	}
	energyCost, err := client.GetEnergyCost(tcs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Energy Cost: %d", energyCost)
}

func TestTrc20Transfer(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	contractAddress, err := common.DecodeAddress(TokenAddress)
	if err != nil {
		t.Error(err)
		return
	}

	hexKey := MainPrivateKey
	toAddress := ReceiverAddress

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

	amountStr := "100.123456"

	amount, _ := decimal.NewFromString(amountStr)
	decimals := int32(TokenDecimals)
	tx, err := client.Trc20Transfer(contractAddress, from, to, amount, decimals, Trc20TransferFeeLimit)
	if err != nil {
		t.Error(err)
		return
	}

	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Transaction ID: %x", tx.Txid)
}

func TestGetTrc20ApproveEnergyCost(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	contractAddress, err := common.DecodeAddress(TokenAddress)
	if err != nil {
		t.Error(err)
		return
	}

	from, err := common.DecodeAddress(SlaveAddress)
	if err != nil {
		t.Error(err)
		return
	}

	to, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}

	tsc, err := client.GenerateTrc20ApproveTrigger(contractAddress, from, to, common.UnlimitedApproveAmount, 0)
	if err != nil {
		t.Error(err)
		return
	}
	energyCost, err := client.GetEnergyCost(tsc)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Energy Cost: %d", energyCost)
}

func TestTrc20Approve(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	contractAddress, err := common.DecodeAddress(TokenAddress)
	if err != nil {
		t.Error(err)
		return
	}

	from, err := common.DecodeAddress(SlaveAddress)
	if err != nil {
		t.Error(err)
		return
	}

	spender, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}

	tsc, err := client.GenerateTrc20ApproveTrigger(contractAddress, from, spender, common.UnlimitedApproveAmount, int32(TokenDecimals))
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := client.BuildContractTransaction(tsc, Trc20ApproveFeeLimit)
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := crypto.HexToECDSA(SlavePrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Transaction ID: %x", tx.Txid)
}

func TestGetTrc20TransferFromEnergyCost(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	contractAddress, err := common.DecodeAddress(TokenAddress)
	if err != nil {
		t.Error(err)
		return
	}

	spender, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}

	sender, err := common.DecodeAddress(SlaveAddress)
	if err != nil {
		t.Error(err)
		return
	}

	receiver, err := common.DecodeAddress(ReceiverAddress)
	if err != nil {
		t.Error(err)
		return
	}

	amount, _ := decimal.NewFromString("100")

	tsc, err := client.GenerateTrc20TransferFromTrigger(contractAddress, spender, sender, receiver, amount, int32(TokenDecimals))
	if err != nil {
		t.Error(err)
		return
	}
	energyCost, err := client.GetEnergyCost(tsc)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Energy Cost: %d", energyCost)
}

func TestTrc20TransferFrom(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	contractAddress, err := common.DecodeAddress(TokenAddress)
	if err != nil {
		t.Error(err)
		return
	}

	spender, err := common.DecodeAddress(MainAddress)
	if err != nil {
		t.Error(err)
		return
	}

	sender, err := common.DecodeAddress(SlaveAddress)
	if err != nil {
		t.Error(err)
		return
	}

	receiver, err := common.DecodeAddress(ReceiverAddress)
	if err != nil {
		t.Error(err)
		return
	}

	amount, _ := decimal.NewFromString("100")

	tsc, err := client.GenerateTrc20TransferFromTrigger(contractAddress, spender, sender, receiver, amount, int32(TokenDecimals))
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := client.BuildContractTransaction(tsc, Trc20TransferFromFeeLimit)
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := crypto.HexToECDSA(MainPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Transaction ID: %x", tx.Txid)
}
