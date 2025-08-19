package common

import (
	"crypto/sha256"
	"fmt"
	"github.com/mr-tron/base58"
	"slices"
)

const (
	addressLength = 25
	prefix        = 0x41
)

type Address []byte

func EncodeAddress(b Address) string {
	s0 := sha256.New()
	s0.Write(b)
	b0 := s0.Sum(nil)

	s1 := sha256.New()
	s1.Write(b0)
	b1 := s1.Sum(nil)

	check := b
	check = append(check, b1[:4]...)

	return base58.Encode(check)
}

func DecodeAddress(address string) (Address, error) {
	decode, err := base58.Decode(address)
	if err != nil {
		return nil, err
	}

	if len(decode) != addressLength {
		return nil, fmt.Errorf("invalid address")
	}

	if decode[0] != prefix {
		return nil, fmt.Errorf("invalid address prefix")
	}

	check := decode[len(decode)-4:]
	decodeData := decode[:len(decode)-4]

	s0 := sha256.New()
	s0.Write(decodeData)
	b0 := s0.Sum(nil)

	s1 := sha256.New()
	s1.Write(b0)
	b1 := s1.Sum(nil)
	b1 = b1[:4]

	if slices.Equal(check, b1) {
		return decodeData, nil
	} else {
		return nil, fmt.Errorf("address check error")
	}
}
