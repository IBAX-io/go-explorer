/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-ibax/packages/converter"
	log "github.com/sirupsen/logrus"
)

// Contract represents record of 1_contracts table
type Contract struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
	WalletID    int64  `json:"wallet_id,omitempty"`
	Active      bool   `json:"active,omitempty"`
	TokenID     int64  `json:"token_id,omitempty"`
	Conditions  string `json:"conditions,omitempty"`
	AppID       int64  `json:"app_id,omitempty"`
	EcosystemID int64  `gorm:"column:ecosystem" json:"ecosystem_id,omitempty"`
}

// TableName returns name of table
func (c *Contract) TableName() string {
	return `1_contracts`
}

func (c *Contract) GetById(id int64) (bool, error) {
	return isFound(GetDB(nil).Where("id = ?", id).First(c))
}

func GetContractCodeByName(contractName string) string {
	if contractName == UtxoTx || contractName == UtxoTransfer {
		return ""
	}
	ecosystem, name := converter.ParseName(contractName)
	var c Contract
	f, err := isFound(GetDB(nil).Select("value").Where("name = ? AND ecosystem = ?", name, ecosystem).First(&c))
	if err == nil && f {
		return c.Value
	}
	return ""
}

func (c *Contract) GetContractsByEcoLibs(ecosystem int64) int64 {
	var total int64
	if err := GetDB(nil).Table(c.TableName()).Where("ecosystem = ?", ecosystem).Count(&total).Error; err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("GetContractsByEcoLibs err")
		return 0
	}
	return total
}

func (c *Contract) GetByApp(appID int64, ecosystemID int64) ([]Contract, error) {
	var result []Contract
	err := GetDB(nil).Select("id,name,value,conditions").Where("app_id = ? and ecosystem = ?", appID, ecosystemID).Find(&result).Error
	return result, err
}
