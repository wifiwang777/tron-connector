package common

import "encoding/base64"

func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
