package client

import (
	"context"
	"crypto/x509"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/wifiwang777/tron-connector/common"
	"github.com/wifiwang777/tron-connector/protos/api"
	"github.com/wifiwang777/tron-connector/protos/core"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	Conn api.WalletClient
}

type options struct {
	useSSL         bool
	certPool       *x509.CertPool
	apiHeader      string
	apiKey         string
	customDialOpts []grpc.DialOption
}

type Option func(*options)

func WithSSL(certPool *x509.CertPool) Option {
	return func(o *options) {
		o.useSSL = true
		o.certPool = certPool
	}
}

func WithAPIKey(header, key string) Option {
	return func(o *options) {
		o.apiHeader = header
		o.apiKey = key
	}
}

func WithCustomDialOptions(opts ...grpc.DialOption) Option {
	return func(o *options) {
		o.customDialOpts = append(o.customDialOpts, opts...)
	}
}

func NewClient(endpoint string, setter ...Option) (*Client, error) {
	opts := &options{
		useSSL: false,
	}

	for _, set := range setter {
		set(opts)
	}

	var grpcOpts []grpc.DialOption

	if opts.useSSL {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(opts.certPool, "")))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if opts.apiKey != "" && opts.apiHeader != "" {
		grpcOpts = append(grpcOpts, grpc.WithChainUnaryInterceptor(UnaryAPIKeyInterceptor(opts.apiHeader, opts.apiKey)))
	}

	if len(opts.customDialOpts) > 0 {
		grpcOpts = append(grpcOpts, opts.customDialOpts...)
	}

	conn, err := grpc.NewClient(endpoint, grpcOpts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn: api.NewWalletClient(conn),
	}, nil
}

func UnaryAPIKeyInterceptor(apiHeader, apiKey string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx, apiHeader, apiKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
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

func (c *Client) Transfer(from, to common.Address, amount decimal.Decimal) (*common.Transaction, error) {
	var err error
	if amount.Sign() < 1 {
		err = fmt.Errorf("amount must be positive")
		return nil, err
	}

	actualDecimals := -amount.Exponent()
	if actualDecimals > common.TronDecimals {
		err = fmt.Errorf("amount exceeds maximum precision (%d decimals)", common.TronDecimals)
	}

	amount = amount.Shift(common.TronDecimals)

	in := &core.TransferContract{
		OwnerAddress: from,
		ToAddress:    to,
		Amount:       amount.IntPart(),
	}
	tx, err := c.Conn.CreateTransaction2(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return common.NewTransaction(tx), nil
}

func (c *Client) DelegateResource(from, to common.Address, resource core.ResourceCode, balance int64) (*common.Transaction, error) {
	dr := &core.DelegateResourceContract{
		OwnerAddress:    from,
		ReceiverAddress: to,
		Resource:        resource,
		Balance:         balance,
	}

	tx, err := c.Conn.DelegateResource(context.Background(), dr)
	if err != nil {
		return nil, err
	}
	return common.NewTransaction(tx), nil
}

func (c *Client) BroadcastTransaction(tx *core.Transaction) error {
	_, err := c.Conn.BroadcastTransaction(context.Background(), tx)
	return err
}

func (c *Client) GetCurrentEnergyPrice() (int64, error) {
	result, err := c.Conn.GetEnergyPrices(context.Background(), &api.EmptyMessage{})
	if err != nil {
		return 0, err
	}

	if result.Prices == "" {
		err = fmt.Errorf("empty energy price result, check network connection")
		return 0, err
	}

	pairs := strings.Split(result.Prices, ",")
	if len(pairs) == 0 {
		err = fmt.Errorf("invalid energy price format: %s", result.Prices)
		return 0, err
	}

	lastPair := pairs[len(pairs)-1]
	parts := strings.Split(lastPair, ":")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid energy price result: %s", result.Prices)
		return 0, err
	}

	price, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse energy price: %w", err)
	}

	return price, nil
}

func (c *Client) GetCurrentBandwidthPrice() (int64, error) {
	result, err := c.Conn.GetBandwidthPrices(context.Background(), &api.EmptyMessage{})
	if err != nil {
		return 0, err
	}
	if result.Prices == "" {
		err = fmt.Errorf("empty bandwidth price result, check network connection")
		return 0, err
	}

	pairs := strings.Split(result.Prices, ",")
	if len(pairs) == 0 {
		err = fmt.Errorf("invalid bandwidth price format: %s", result.Prices)
		return 0, err
	}

	lastPair := pairs[len(pairs)-1]
	parts := strings.Split(lastPair, ":")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid bandwidth price result: %s", result.Prices)
		return 0, err
	}

	price, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse bandwidth price: %w", err)
	}

	return price, nil
}

// Contract

func (c *Client) GetEnergyCost(tsc *core.TriggerSmartContract) (int64, error) {
	result, err := c.Conn.TriggerConstantContract(context.Background(), tsc)
	if err != nil {
		return 0, err
	}

	if result.Result.Message != nil {
		err = fmt.Errorf("invalid tx: %s", result.Result.Message)
		return 0, err
	}

	return result.EnergyUsed, nil
}

func (c *Client) BuildContractTransaction(tsc *core.TriggerSmartContract, feeLimit int64) (*common.Transaction, error) {
	tx, err := c.Conn.TriggerContract(context.Background(), tsc)
	if err != nil {
		return nil, err
	}

	res := common.NewTransaction(tx)

	if feeLimit > 0 {
		res.Transaction.RawData.FeeLimit = feeLimit
		res.Txid, err = res.CalculateDigestHash()
		if err != nil {
			return nil, err
		}
	}

	return common.NewTransaction(tx), nil
}

// TRC20

func (c *Client) GetTrc20Balance(account, contractAddress common.Address, decimals int32) (decimal.Decimal, error) {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte("balanceOf(address)"))
	methodId := hasher.Sum(nil)[:4]

	data := make([]byte, 0, 36)
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(account, 32)...)

	tsc := &core.TriggerSmartContract{
		OwnerAddress:    account,
		ContractAddress: contractAddress,
		Data:            data,
	}

	result, err := c.Conn.TriggerConstantContract(context.Background(), tsc)
	if err != nil {
		return decimal.Zero, err
	}
	if len(result.ConstantResult) == 0 || len(result.ConstantResult[0]) == 0 {
		return decimal.Zero, fmt.Errorf("empty constant result from contract, check contract address or network")
	}
	weiBalance := new(big.Int).SetBytes(result.ConstantResult[0])
	decWei := decimal.NewFromBigInt(weiBalance, 0)
	humanBalance := decWei.Shift(-decimals)

	return humanBalance, nil
}

func (c *Client) GetTrc20Allowance(contractAddress, owner, spender common.Address) (*big.Int, error) {
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

	ct := &core.TriggerSmartContract{
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

func (c *Client) GenerateTrc20TransferTrigger(contractAddress, from, to common.Address, amount decimal.Decimal, decimals int32) (*core.TriggerSmartContract, error) {
	var err error
	if amount.Sign() < 1 {
		err = fmt.Errorf("amount must be positive")
		return nil, err
	}

	actualDecimals := -amount.Exponent()
	if actualDecimals > decimals {
		err = fmt.Errorf("amount exceeds maximum precision (%d decimals)", decimals)
		return nil, err
	}

	amount = amount.Shift(decimals)

	methodSignature := []byte("transfer(address,uint256)")
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(methodSignature)
	methodId := keccak256.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(to, 32)...)
	data = append(data, common.LeftPadBytes(amount.BigInt().Bytes(), 32)...)
	tsc := &core.TriggerSmartContract{
		OwnerAddress:    from,
		ContractAddress: contractAddress,
		Data:            data,
	}
	return tsc, nil
}

func (c *Client) Trc20Transfer(contractAddress, from, to common.Address, amount decimal.Decimal, decimals int32, feeLimit int64) (*common.Transaction, error) {
	tsc, err := c.GenerateTrc20TransferTrigger(contractAddress, from, to, amount, decimals)
	if err != nil {
		return nil, err
	}
	return c.BuildContractTransaction(tsc, feeLimit)
}

func (c *Client) GenerateTrc20ApproveTrigger(contractAddress, from, spender common.Address, amount decimal.Decimal, decimals int32) (*core.TriggerSmartContract, error) {
	var err error
	if amount.Sign() < 1 {
		err = fmt.Errorf("amount must be positive")
		return nil, err
	}
	actualDecimals := -amount.Exponent()
	if actualDecimals > decimals {
		err = fmt.Errorf("amount exceeds maximum precision (%d decimals)", decimals)
		return nil, err
	}

	amount = amount.Shift(decimals)
	methodSignature := []byte("approve(address,uint256)")
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(methodSignature)
	methodId := keccak256.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(spender, 32)...)
	data = append(data, common.LeftPadBytes(amount.BigInt().Bytes(), 32)...)
	tsc := &core.TriggerSmartContract{
		OwnerAddress:    from,
		ContractAddress: contractAddress,
		Data:            data,
	}
	return tsc, nil
}

func (c *Client) Trc20Approve(contractAddress, from, to common.Address, amount decimal.Decimal, decimals int32, feeLimit int64) (*common.Transaction, error) {
	tsc, err := c.GenerateTrc20ApproveTrigger(contractAddress, from, to, amount, decimals)
	if err != nil {
		return nil, err
	}
	return c.BuildContractTransaction(tsc, feeLimit)
}

func (c *Client) GenerateTrc20TransferFromTrigger(contractAddress, spender, sender, receiver common.Address, amount decimal.Decimal, decimals int32) (*core.TriggerSmartContract, error) {
	var err error
	if amount.Sign() < 1 {
		err = fmt.Errorf("amount must be positive")
		return nil, err
	}
	actualDecimals := -amount.Exponent()
	if actualDecimals > decimals {
		err = fmt.Errorf("amount exceeds maximum precision (%d decimals)", decimals)
		return nil, err
	}

	amount = amount.Shift(decimals)
	methodSignature := []byte("transferFrom(address,address,uint256)")
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(methodSignature)
	methodId := keccak256.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(sender, 32)...)
	data = append(data, common.LeftPadBytes(receiver, 32)...)
	data = append(data, common.LeftPadBytes(amount.BigInt().Bytes(), 32)...)
	tsc := &core.TriggerSmartContract{
		OwnerAddress:    spender,
		ContractAddress: contractAddress,
		Data:            data,
	}
	return tsc, nil
}

func (c *Client) Trc20TransferFrom(contractAddress, spender, sender, receiver common.Address, amount decimal.Decimal, decimals int32, feeLimit int64) (*common.Transaction, error) {
	tsc, err := c.GenerateTrc20TransferFromTrigger(contractAddress, spender, sender, receiver, amount, decimals)
	if err != nil {
		return nil, err
	}
	return c.BuildContractTransaction(tsc, feeLimit)
}
