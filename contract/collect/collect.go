package collect

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/wifiwang777/tron-connector/common"
	"github.com/wifiwang777/tron-connector/contract"
	"github.com/wifiwang777/tron-connector/protos/core"
)

type Collector struct {
	contract.Contract
}

func NewCollector(contract contract.Contract) *Collector {
	return &Collector{
		Contract: contract,
	}
}

func (c *Collector) AddWorker(owner, newWorker common.Address) (*core.TriggerSmartContract, error) {
	data, err := parsedABI.Pack("addWorker", ethcommon.Address(newWorker[1:]))
	if err != nil {
		return nil, err
	}
	return &core.TriggerSmartContract{
		OwnerAddress:    owner,
		ContractAddress: c.Address,
		Data:            data,
	}, nil
}

func (c *Collector) RemoveWorker(owner, workerToRemove common.Address) (*core.TriggerSmartContract, error) {
	data, err := parsedABI.Pack("removeWorker", ethcommon.Address(workerToRemove[1:]))
	if err != nil {
		return nil, err
	}
	return &core.TriggerSmartContract{
		OwnerAddress:    owner,
		ContractAddress: c.Address,
		Data:            data,
	}, nil
}

func (c *Collector) CollectSingle(worker, userAddress common.Address) (*core.TriggerSmartContract, error) {
	data, err := parsedABI.Pack("collectSingle", ethcommon.Address(userAddress[1:]))
	if err != nil {
		return nil, err
	}
	return &core.TriggerSmartContract{
		OwnerAddress:    worker,
		ContractAddress: c.Address,
		Data:            data,
	}, nil
}

func (c *Collector) CollectBatch(worker common.Address, userAddresses []common.Address) (*core.TriggerSmartContract, error) {
	var ethAddresses []ethcommon.Address
	for _, userAddress := range userAddresses {
		ethAddresses = append(ethAddresses, ethcommon.Address(userAddress[1:]))
	}
	data, err := parsedABI.Pack("collectBatch", ethAddresses)
	if err != nil {
		return nil, err
	}
	return &core.TriggerSmartContract{
		OwnerAddress:    worker,
		ContractAddress: c.Address,
		Data:            data,
	}, nil
}
