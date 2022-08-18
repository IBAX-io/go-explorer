/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import "github.com/shopspring/decimal"

type SumAmount struct {
	Sum decimal.Decimal `gorm:"column:sum"`
}

type MaxInt64 struct {
	Max int64 `gorm:"column:max"`
}

type SumInt64 struct {
	Sum int64 `gorm:"column:sum"`
}

type CountInt64 struct {
	Count int64 `gorm:"column:count"`
}

type MinInt64 struct {
	Min int64 `gorm:"column:min"`
}
