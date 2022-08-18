/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX  All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"strconv"

	"github.com/IBAX-io/go-explorer/conf"
)

// StateParameter is model
type StateParameter struct {
	ecosystem  int64
	ID         int64  `gorm:"primary_key;not null" json:"id"`
	Name       string `gorm:"not null;size:100" json:"name"`
	Value      string `gorm:"not null" json:"value"`
	Conditions string `gorm:"not null" json:"conditions"`
}

type EcosystemParameterResult struct {
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Rets  []StateParameter `json:"rets"`
}

// TableName returns name of table
func (sp *StateParameter) TableName() string {
	if sp.ecosystem == 0 {
		sp.ecosystem = 1
	}
	return `1_parameters`
}

// SetTablePrefix is setting table prefix
func (sp *StateParameter) SetTablePrefix(prefix string) {
	pre, _ := strconv.ParseInt(prefix, 10, 64)
	sp.ecosystem = pre
}

// SetTablePrefix is setting table prefix
func (sp *StateParameter) SetTableFix(fix int64) {
	sp.ecosystem = fix
}

// Get is retrieving model from database
func (sp *StateParameter) Get(name string) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("ecosystem = ? and name = ?", sp.ecosystem, name).First(sp))
}

// Get is retrieving model from database
func (sp *StateParameter) GetMintAmount() (string, error) {
	var sp1, sp2 StateParameter
	f, err := isFound(GetDB(nil).Where("ecosystem = ? and name = ?", sp.ecosystem, "mint_balance").First(sp))
	if err != nil {
		return "0", err
	}

	f1, err := isFound(GetDB(nil).Where("ecosystem = ? and name = ?", sp.ecosystem, "foundation_balance").First(&sp1))
	if err != nil {
		return "0", err
	}

	f2, err := isFound(GetDB(nil).Where("ecosystem = ? and name = ?", sp.ecosystem, "assign_rule").First(&sp2))
	if err != nil {
		return "0", err
	}
	if !f || !f1 || !f2 {
		return "0", nil
	}

	ret := make(map[int64]AssignRules, 1)
	err = json.Unmarshal([]byte(sp2.Value), &ret)
	if err != nil {
		return "0", err
	}

	as3, ok3 := ret[3]
	as6, ok6 := ret[6]
	if !ok3 || !ok6 {
		return "0", nil
	}

	tfa, err := decimal.NewFromString(as3.TotalAmount)
	if err != nil {
		return "0", err
	}
	tma, err := decimal.NewFromString(as6.TotalAmount)
	if err != nil {
		return "0", err
	}

	ma, err := decimal.NewFromString(sp.Value)
	if err != nil {
		return "0", err
	}
	fa, err := decimal.NewFromString(sp1.Value)
	if err != nil {
		return "0", err
	}
	if fa.LessThanOrEqual(tfa) && ma.LessThanOrEqual(tma) {
		mb := tma.Sub(ma)
		fb := tfa.Sub(fa)
		tt := mb.Add(fb)
		return tt.String(), nil
	} else {
		return "0", errors.New("assign rules err")
	}

	return "0", nil
}

func (sp *StateParameter) FindStateParameters(page int, size int, name, order string, ecosystem int64) (num int64, rets []StateParameter, err error) {
	ns := "%" + name + "%"
	if err := GetDB(nil).Table(sp.TableName()).Where("name like ? AND ecosystem = ?", ns, ecosystem).Count(&num).Error; err != nil {
		return num, rets, err
	}

	err = GetDB(nil).Table(sp.TableName()).Where("name like ? AND ecosystem = ?", ns, ecosystem).
		Order(order).Offset((page - 1) * size).Limit(size).Find(&rets).Error

	return num, rets, err
}
