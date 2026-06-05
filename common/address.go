package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

type Address []byte

func (a Address) String() string {
	if len(a) == 0 {
		return ""
	}

	// 1. 连续两次 SHA256 计算校验和
	h1 := sha256.Sum256(a)
	h2 := sha256.Sum256(h1[:])

	// 2. 优雅拼接：预先分配好 25 字节空间，绝对不污染原切片 a，且只有一次内存分配
	check := make([]byte, len(a)+checksumLength)
	copy(check, a)
	copy(check[len(a):], h2[:checksumLength])

	// 3. Base58 编码
	return base58.Encode(check)
}

func DecodeAddress(address string) (Address, error) {
	decode, err := base58.Decode(address)
	if err != nil {
		return nil, fmt.Errorf("base58 decode error: %w", err)
	}

	// 长度校验
	if len(decode) != addressLength {
		return nil, fmt.Errorf("invalid address length")
	}

	// 前缀校验
	if decode[0] != prefix {
		return nil, fmt.Errorf("invalid address prefix")
	}

	// 拆分数据与校验和
	splitIdx := len(decode) - checksumLength
	decodeData := decode[:splitIdx]
	actualCheck := decode[splitIdx:]

	// 重新计算校验和
	h1 := sha256.Sum256(decodeData)
	h2 := sha256.Sum256(h1[:])
	expectedCheck := h2[:checksumLength]

	// 校验和比对 (Go 自带的 slice 比对，比循环更高效)
	// 这里不需要用 slices.Equal，直接利用 string 或者原生循环，为减少依赖直接利用简易比对
	for i := 0; i < checksumLength; i++ {
		if actualCheck[i] != expectedCheck[i] {
			return nil, fmt.Errorf("address checksum mismatch")
		}
	}

	return decodeData, nil
}

func PublicKeyToAddress(publicKey ecdsa.PublicKey) Address {
	address := crypto.PubkeyToAddress(publicKey)

	addressTron := make([]byte, 1+len(address.Bytes()))
	addressTron[0] = prefix
	copy(addressTron[1:], address.Bytes())

	return addressTron
}
