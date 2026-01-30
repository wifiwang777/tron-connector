package common

import "math/big"

func MultiplyBy10Power(num *big.Int, n int) *big.Int {
	result := new(big.Int).Set(num) // 复制num

	// 计算 10^n
	power := new(big.Int).Exp(
		big.NewInt(10),       // 基数 10
		big.NewInt(int64(n)), // 指数 n
		nil,                  // 模运算，nil表示不需要
	)

	// num * 10^n
	result.Mul(result, power)

	return result
}
