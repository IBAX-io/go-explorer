/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/shopspring/decimal"
)

var (
	MintNodeReady bool

	MintNodeTotalBalance     decimal.Decimal
	MintNodeTotalSupplyToken decimal.Decimal
)

func GetMintNodeTotalBalance() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	if MintNodeTotalSupplyToken.IsZero() {
		MintNodeTotalSupplyToken = decimal.New(MintNodeTotalSupply, int32(consts.MoneyDigits))
	}
	if !MintNodeReady {
		MintNodeTotalBalance = MintNodeTotalSupplyToken
	}
}
