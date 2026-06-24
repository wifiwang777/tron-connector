package common

import "github.com/shopspring/decimal"

const (
	addressLength  = 21
	prefix         = 0x41
	checksumLength = 4

	TronDecimals = 6
)

var (
	UnlimitedApproveAmount = decimal.NewFromInt(2).Pow(decimal.NewFromInt(256)).Sub(decimal.NewFromInt(1)) // ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
)
