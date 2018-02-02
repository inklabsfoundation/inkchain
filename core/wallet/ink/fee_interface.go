/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ink

type InkAlg interface {
	CalcInk(textLength int) (int64, error)
}
