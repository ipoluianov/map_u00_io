package utils

import "math/big"

func ParseDatabigInt(data []byte) *big.Int {
	v := big.NewInt(0)
	if len(data) == 32 {
		v = v.SetBytes(data)
	}
	return v
}
