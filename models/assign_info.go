/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"errors"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
	//"time"
)

var (
	AssignReady bool

	AssignTotalBalance     decimal.Decimal
	AssignTotalSupplyToken decimal.Decimal
)

type AssignRules struct {
	StartBlockID    int64  `json:"start_blockid"`
	EndBlockID      int64  `json:"end_blockid"`
	IntervalBlockID int64  `json:"interval_blockid"`
	Count           int64  `json:"count"`
	TotalAmount     string `json:"total_amount"`
}

// AssignInfo is model
type AssignInfo struct {
	ID            int64           `gorm:"primary_key;not null"`
	Type          int64           `gorm:"not null"`
	Account       string          `gorm:"not null"`
	TotalAmount   decimal.Decimal `gorm:"not null"`
	BalanceAmount decimal.Decimal `gorm:"not null"`
	Detail        string          `gorm:"not null;type:jsonb"`
	Deleted       int64           `gorm:"not null"`
	DateDeleted   int64           `gorm:"not null"`
	DateUpdated   int64           `gorm:"not null"`
	DateCreated   int64           `gorm:"not null"`
}

// TableName returns name of table
func (m AssignInfo) TableName() string {
	return `1_assign_info`
}

// GetTotalBalance is retrieving model from database
func (m *AssignInfo) GetTotalBalance(dbTx *DbTransaction, account string) (decimal.Decimal, decimal.Decimal, error) {

	var mps []AssignInfo
	var amount, balance decimal.Decimal
	amount = decimal.NewFromFloat(0)
	balance = decimal.NewFromFloat(0)
	if !HasTable(m) {
		return amount, balance, nil
	}
	err := GetDB(nil).Table(m.TableName()).
		Where("account =? AND deleted = 0 AND balance_amount > 0", account).
		Find(&mps).Error
	if err != nil {
		return amount, balance, err
	}
	if len(mps) == 0 {
		return amount, balance, nil
	}

	now := time.Now()
	for _, t := range mps {
		list, err := getAssignDetail(t.Detail, t.Type)
		if err != nil {
			return amount, balance, err
		}

		for _, v := range list {
			st, _ := strconv.ParseInt(v.StartAt, 10, 64)
			if st >= FirstBlockTime && st <= now.Unix() && v.Status == 1 {
				am, _ := decimal.NewFromString(v.Amount)
				amount = amount.Add(am)
			}
		}
		balance = balance.Add(t.BalanceAmount)
	}
	return amount, balance, err
}

func (m *AssignInfo) GetBalance(db *DbTransaction, account string) (decimal.Decimal, error) {

	var totalBalance decimal.Decimal
	if AssignReady {
		query := GetDB(db).Table(m.TableName()).Select("coalesce(sum(balance_amount),0)").
			Where("deleted =?", 0)
		if account != "" {
			query = query.Where("account = ?", account)
		}
		err := query.Take(&totalBalance).Error
		if err != nil {
			return totalBalance, err
		}
	}
	return totalBalance, nil
}

type assignDetail struct {
	Amount  string `json:"amount"`
	Status  int    `json:"status"`
	StartAt string `json:"startAt"`
	ClaimAt string `json:"claimAt"`
}

func getAssignDetail(detail string, assignType int64) ([]assignDetail, error) {
	var list []assignDetail
	switch assignType {
	case 1, 2, 3, 4, 5, 6:
		err := json.Unmarshal([]byte(detail), &list)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("not support assign type")
	}
	return list, nil

}

func AssignTableExist() bool {
	var p AssignInfo
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}

func GetAssignTotalBalanceAmount() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	if AssignTotalSupplyToken.IsZero() {
		getAssignTotalSupplyToken()
	}
	if AssignReady {
		getAssignNowTotalBalance()
	} else {
		AssignTotalBalance = AssignTotalSupplyToken
	}
}

func getAssignNowTotalBalance() {
	err := GetDB(nil).Raw(`SELECT CAST((SELECT value FROM "1_app_params" WHERE "name" = 'balance_supply_foundation' AND ecosystem = 1) as numeric)+
			CAST((SELECT value FROM "1_app_params" WHERE "name" = 'balance_supply_partners' AND ecosystem = 1) as numeric)+
			CAST((SELECT value FROM "1_app_params" WHERE "name" = 'balance_supply_private_round1' AND ecosystem = 1) as numeric)+
			CAST((SELECT value FROM "1_app_params" WHERE "name" = 'balance_supply_private_round2' AND ecosystem = 1) as numeric)+
			CAST((SELECT value FROM "1_app_params" WHERE "name" = 'balance_supply_public_round' AND ecosystem = 1) as numeric)+
			CAST((SELECT value FROM "1_app_params" WHERE "name" = 'balance_supply_dev_team' AND ecosystem = 1) as numeric)+
			(SELECT COALESCE(sum(balance_amount),0) FROM "1_assign_info")
		AS lock_amount
	`).Take(&AssignTotalBalance).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get now assign Lock All Failed")
	}
}

func getAssignTotalSupplyToken() {
	AssignTotalSupplyToken = decimal.New(AssignTotalSupply, int32(consts.MoneyDigits))
}
