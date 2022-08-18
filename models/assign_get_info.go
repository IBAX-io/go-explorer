/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	//"time"
)

type AssignRules struct {
	StartBlockID    int64  `json:"start_blockid"`
	EndBlockID      int64  `json:"end_blockid"`
	IntervalBlockID int64  `json:"interval_blockid"`
	Count           int64  `json:"count"`
	TotalAmount     string `json:"total_amount"`
}

// AssignGetInfo is model
type AssignGetInfo struct {
	ID            int64           `gorm:"primary_key;not null"`
	Type          int64           `gorm:"not null"`
	Account       int64           `gorm:"not null"`
	TotalAmount   decimal.Decimal `gorm:"not null"`
	BalanceAmount decimal.Decimal `gorm:"not null"`
	Amount        decimal.Decimal `gorm:"not null"`
	Latestid      int64           `gorm:"not null"`
	Deleted       int64           `gorm:"not null"`
	DateUpdated   int64           `gorm:"not null" `
	DateCreated   int64           `gorm:"not null" `
}

// TableName returns name of table
func (m AssignGetInfo) TableName() string {
	return `1_assign_get_info`
}

// GetId is retrieving model from database
func (m *AssignGetInfo) GetBalance(db *DbTransaction, wallet int64) (bool, decimal.Decimal, decimal.Decimal, error) {

	var mps []AssignGetInfo
	var balance, totalBalance decimal.Decimal
	balance = decimal.NewFromFloat(0)
	totalBalance = decimal.NewFromFloat(0)
	if !HasTable(m) {
		return false, balance, totalBalance, nil
	}
	err := GetDB(nil).Table(m.TableName()).
		Where("keyid = ? and deleted =? ", wallet, 0).
		Find(&mps).Error
	if err != nil {
		return false, balance, totalBalance, err
	}
	if len(mps) == 0 {
		return false, balance, totalBalance, err
	}

	//newblockid
	block := &Block{}
	found, err := block.GetMaxBlock()
	if err != nil {
		return false, balance, totalBalance, err
	}
	if !found {
		return false, balance, totalBalance, errors.New("maxblockid not found")
	}

	//assign_rule
	var sp StateParameter
	sp.SetTablePrefix(`1`)
	found1, err1 := sp.Get(`assign_rule`)
	if err1 != nil {
		return false, balance, totalBalance, err1
	}

	if !found1 || len(sp.Value) == 0 {
		return false, balance, totalBalance, errors.New("assign_rule not found or not exist assign_rule")
	}

	rules := make(map[int64]AssignRules, 10)
	err = json.Unmarshal([]byte(sp.Value), &rules)
	if err != nil {
		return false, balance, totalBalance, err
	}

	maxblockid := block.ID
	for _, t := range mps {
		am := decimal.NewFromFloat(0)
		tm := t.BalanceAmount
		rule, ok := rules[t.Type]
		if ok {
			sid := rule.StartBlockID
			iid := rule.IntervalBlockID
			eid := rule.EndBlockID

			if maxblockid >= eid {
				am = tm
			} else {
				if t.Latestid == 0 {
					count := int64(0)
					if maxblockid > sid {
						count = (maxblockid - sid) / iid
						count += 1
					}
					if count > 0 {
						if t.Type == 4 {
							//first
							if t.Latestid == 0 {
								am = t.BalanceAmount.Mul(decimal.NewFromFloat(0.1))
								sm := t.Amount.Mul(decimal.NewFromFloat(float64(count - 1)))
								am = am.Add(sm)
							} else {
								am = t.Amount.Mul(decimal.NewFromFloat(float64(count)))
							}

						} else {
							am = t.Amount.Mul(decimal.NewFromFloat(float64(count)))
						}
					}

				} else {
					if maxblockid > t.Latestid {
						count := (maxblockid - t.Latestid) / iid
						am = t.Amount.Mul(decimal.NewFromFloat(float64(count)))
					}
				}
			}

			tm = tm.Sub(am)
			balance = balance.Add(am)
			totalBalance = totalBalance.Add(tm)
		}
	}
	return true, balance, totalBalance, err
}
func (m *AssignGetInfo) GetAllBalance(db *DbTransaction) (decimal.Decimal, error) {

	var mps []AssignGetInfo
	var balance decimal.Decimal
	balance = decimal.NewFromFloat(0)
	if HasTableOrView(nil, m.TableName()) {
		err := GetDB(db).Table(m.TableName()).Select("balance_amount").
			Where("deleted =?", 0).
			Find(&mps).Error
		if err != nil {
			return balance, err
		}
		if len(mps) == 0 {
			return balance, err
		}

		for _, ai := range mps {
			balance.Add(ai.BalanceAmount)
		}
	}
	return balance, nil
}
