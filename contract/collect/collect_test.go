package collect

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wifiwang777/tron-connector/client"
	"github.com/wifiwang777/tron-connector/common"
	"github.com/wifiwang777/tron-connector/contract"
)

const (
	GrpcEndpointNile           = "grpc.nile.trongrid.io:50051"
	CollectContractAddressNile = "TLDFVwxT76cciuQqDkNpieUEn5qDaRfaKB"
	USDTContractAddressNile    = "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj"

	worker1PrivateKey = "4be2df367fe87cf69cb6b45de780cf03eb414238173fec895517bb94873766ed"
	worker1Address    = "TYVWsM71CUYUQLRpuvbNSPYbuosWVJBpR3"

	worker2PrivateKey = "8112fb602e167ee738d01ec79fab1042151dc8387408d3bcba6b2e6f98a7c92e"
	worker2Address    = "TAgAwtYRL4ZPwo3CMbV88Pa2PCio4v9zZZ"

	user1PrivateKey = "292f21e0ace817d5d4f00975352221f4b954cebc92ce65dc9ee6c19fe3e10882"
	user1Address    = "THBn2TqubxehwZCZyP18EmVRPvpzqNWbVU"

	user2PrivateKey = "97034ab42819200c800be0a56872dea227bb7df3e03a66a2603d1a61cd49856b"
	user2Address    = "TK7sStvytqL4BjaQhQWb8333wWYJCrjq7s"

	user3PrivateKey = "cb8ff038f7d43689b423e876192a682c01fbc276a12eaab2fd72cc5a651f1440"
	user3Address    = "TRo5jMTGjAcP7NkjzozxX6jwPzuPiyNssr"
)

var (
	collectContractAddress = CollectContractAddressNile
	grpcEndpoint           = GrpcEndpointNile

	ownerPrivateKey = "5e09c723aa1010ac7b47c63f333a38e06727f67c897f5c338e3e992119de1cf3"
	ownerAddress    = "TSuiX2SRjsLQKSGEgA3VRPi9hBwse1738T"

	workerPrivateKey = worker2PrivateKey
	workerAddress    = worker2Address

	userPrivateKey = user3PrivateKey
	userAddress    = user3Address

	userAddresses = []string{
		user1Address,
		user2Address,
		user3Address,
	}
)

func newClient() *client.Client {
	cli, _ := client.NewClient(grpcEndpoint)
	return cli
}

func getContract() contract.Contract {
	contractAddress, _ := common.DecodeAddress(collectContractAddress)
	return contract.Contract{
		Address: contractAddress,
	}
}

func TestAddWorker(t *testing.T) {
	collector := NewCollector(getContract())
	cli := newClient()

	owner, err := common.DecodeAddress(ownerAddress)
	if err != nil {
		t.Error(err)
		return
	}

	worker, err := common.DecodeAddress(workerAddress)
	if err != nil {
		t.Error(err)
		return
	}
	tsc, err := collector.AddWorker(owner, worker)
	if err != nil {
		t.Error(err)
		return
	}

	//feeLimit, err := cli.GetEnergyCost(tsc)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Logf("feeLimit: %d", feeLimit)

	energyCost := int64(6260900)
	feeLimit := energyCost * 100
	feeLimit = feeLimit * 12
	feeLimit = feeLimit / 10

	tx, err := cli.BuildContractTransaction(tsc, feeLimit)
	if err != nil {
		return
	}

	privateKey, err := crypto.HexToECDSA(ownerPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = cli.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	txId := tx.Txid
	t.Logf("txId: %x", txId)
}

func TestRemoveWorker(t *testing.T) {
	collector := NewCollector(getContract())
	cli := newClient()

	owner, err := common.DecodeAddress(ownerAddress)
	if err != nil {
		t.Error(err)
		return
	}

	worker, err := common.DecodeAddress(workerAddress)
	if err != nil {
		t.Error(err)
		return
	}

	tsc, err := collector.RemoveWorker(owner, worker)
	if err != nil {
		t.Error(err)
		return
	}

	//feeLimit, err := cli.GetEnergyCost(tsc)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Logf("feeLimit: %d", feeLimit)

	energyCost := int64(23798)
	feeLimit := energyCost * 100
	feeLimit = feeLimit * 12
	feeLimit = feeLimit / 10

	tx, err := cli.BuildContractTransaction(tsc, feeLimit)
	if err != nil {
		return
	}

	privateKey, err := crypto.HexToECDSA(ownerPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = cli.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	txId := tx.Txid
	t.Logf("txId: %x", txId)
}

func TestApprove(t *testing.T) {
	cli := newClient()
	contractAddress, err := common.DecodeAddress(USDTContractAddressNile)
	if err != nil {
		t.Error(err)
		return
	}
	from, err := common.DecodeAddress(userAddress)
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := crypto.HexToECDSA(userPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}

	spender, err := common.DecodeAddress(collectContractAddress)
	if err != nil {
		t.Error(err)
		return
	}

	energyCost := int64(22506)
	feeLimit := energyCost * 100
	feeLimit = feeLimit * 12
	feeLimit = feeLimit / 10

	tx, err := cli.Trc20Approve(contractAddress, from, spender, common.UnlimitedApproveAmount, int32(0), feeLimit)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	err = cli.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Transaction ID: %x", tx.Txid)
}

func TestCollectSingle(t *testing.T) {
	collector := NewCollector(getContract())
	cli := newClient()

	worker, err := common.DecodeAddress(workerAddress)
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := crypto.HexToECDSA(workerPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}

	user, err := common.DecodeAddress(userAddress)
	if err != nil {
		t.Error(err)
	}

	tsc, err := collector.CollectSingle(worker, user)
	if err != nil {
		t.Error(err)
		return
	}

	//energyCost, err := cli.GetEnergyCost(tsc)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Logf("Energy Cost: %d", energyCost)
	energyCost := int64(28752)
	feeLimit := energyCost * 100
	feeLimit = feeLimit * 12
	feeLimit = feeLimit / 10

	tx, err := cli.BuildContractTransaction(tsc, feeLimit)
	if err != nil {
		t.Error(err)
		return
	}

	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = cli.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Transaction ID: %x", tx.Txid)
}

func TestCollectBatch(t *testing.T) {
	collector := NewCollector(getContract())
	cli := newClient()
	worker, err := common.DecodeAddress(workerAddress)
	if err != nil {
		t.Error(err)
	}

	privateKey, err := crypto.HexToECDSA(workerPrivateKey)
	if err != nil {
		t.Error(err)
		return
	}

	var addresses []common.Address
	for _, item := range userAddresses {
		address, err := common.DecodeAddress(item)
		if err != nil {
			t.Error(err)
			return
		}
		addresses = append(addresses, address)
	}

	tsc, err := collector.CollectBatch(worker, addresses)
	if err != nil {
		t.Error(err)
		return
	}

	energyCost, err := cli.GetEnergyCost(tsc)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Energy Cost: %d", energyCost)
	feeLimit := energyCost * 100
	feeLimit = feeLimit * 12
	feeLimit = feeLimit / 10
	tx, err := cli.BuildContractTransaction(tsc, feeLimit)
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.SignWithPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	err = cli.BroadcastTransaction(tx.Transaction)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Transaction ID: %x", tx.Txid)
}
