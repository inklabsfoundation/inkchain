/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package impl

var InkFeeK float32
var InkFeeX0 float32
var InkFeeB float32

type SimpleInkAlg struct {
}

func NewSimpleInkAlg() *SimpleInkAlg {
	return &SimpleInkAlg{}
}
func (s *SimpleInkAlg) CalcInk(textLength int) (int64, error) {
	ink := int64((float32(textLength)-InkFeeX0)*InkFeeK + InkFeeB)
	if ink < 0 {
		ink = 0
	}
	return ink, nil
}
