/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

type IName struct {
	Id          int64 `gorm:"primary_key;not_null"`
	Account     string
	DateCreated int64
	Ecosystem   int64
	Name        string
	TxHash      []byte
}

var INameReady bool

func (p IName) TableName() string {
	return `1_iname`
}

func INameTableExist() bool {
	var p IName
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}

func (p *IName) Get(account string) (bool, error) {
	return isFound(GetDB(nil).Where("account = ? AND ecosystem = 1", account).Take(p))
}
