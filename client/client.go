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
	"github.com/wifiwang777/tron-protocol/protos/api"
	"github.com/wifiwang777/tron-protocol/protos/core"
	"github.com/wifiwang777/tron-protocol/protos/core/contract"
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

// NewClient 开源库对外暴露的唯一构造函数
func NewClient(endpoint string, setter ...Option) (*Client, error) {
	// 1. 初始化默认参数（不走 SSL，不带 API Key）
	opts := &options{
		useSSL: false,
	}

	// 2. 执行用户传入的选项，修改默认参数
	for _, set := range setter {
		set(opts)
	}

	// 3. 构建真正的 gRPC DialOptions 数组
	var grpcOpts []grpc.DialOption

	// 4. 处理安全凭证逻辑
	if opts.useSSL {
		// 如果用户传了具体的 certPool 就用，没传（nil）就默认用系统根证书
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(opts.certPool, "")))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 5. 处理内置的 API Key 拦截器
	if opts.apiKey != "" && opts.apiHeader != "" {
		grpcOpts = append(grpcOpts, grpc.WithChainUnaryInterceptor(UnaryAPIKeyInterceptor(opts.apiHeader, opts.apiKey)))
	}

	if len(opts.customDialOpts) > 0 {
		grpcOpts = append(grpcOpts, opts.customDialOpts...)
	}

	// 7. 建立连接（2026年推荐使用 grpc.NewClient 代替已废弃的 grpc.Dial）
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
		// 1. 将 API Key 注入到外发的 metadata 中
		// gRPC 内部会自动将 "TRON-PRO-API-KEY" 转换为标准的 header 格式
		ctx = metadata.AppendToOutgoingContext(ctx, apiHeader, apiKey)

		// 2. 将带有 API Key 的 context 传给下一个处理器或实际的调用发起者
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

	in := &contract.TransferContract{
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

func (c *Client) DelegateResource(from, to common.Address, resource contract.ResourceCode, balance int64) (*common.Transaction, error) {
	dr := &contract.DelegateResourceContract{
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

	// 0:100,1572597600000:10,1606282800000:40,1612768800000:140,1612769400000:140,1612778400000:140,1628674200000:420,1635143400000:280,1669603800000:420,1726283400000:210,1754644200000:100
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

	// 3. 将价格转换为 int64
	price, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse energy price: %w", err)
	}

	return price, nil
}

// Contract

func (c *Client) GetEnergyCost(tsc *contract.TriggerSmartContract) (int64, error) {
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

func (c *Client) BuildContractTransaction(tsc *contract.TriggerSmartContract, feeLimit int64) (*common.Transaction, error) {
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

	tsc := &contract.TriggerSmartContract{
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

	ct := &contract.TriggerSmartContract{
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

func (c *Client) GenerateTrc20TransferTrigger(contractAddress, from, to common.Address, amount decimal.Decimal, decimals int32) (*contract.TriggerSmartContract, error) {
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
	tsc := &contract.TriggerSmartContract{
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

func (c *Client) GenerateTrc20ApproveTrigger(contractAddress, from, spender common.Address, amount decimal.Decimal, decimals int32) (*contract.TriggerSmartContract, error) {
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
	tsc := &contract.TriggerSmartContract{
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

func (c *Client) GenerateTrc20TransferFromTrigger(contractAddress, spender, sender, receiver common.Address, amount decimal.Decimal, decimals int32) (*contract.TriggerSmartContract, error) {
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
	tsc := &contract.TriggerSmartContract{
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
