package ink

import "math/big"

type InkAlg interface {
	CalcInk(textLength int) (*big.Int, error)
}
