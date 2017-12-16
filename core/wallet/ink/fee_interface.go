/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ink

import "math/big"

type InkAlg interface {
	CalcInk(textLength int) (*big.Int, error)
}
