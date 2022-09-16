/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/converter"
)

//MineInfo example
type MineInfo struct {
	ID                int64  `gorm:"primary_key not null" example:"1"` //
	Devid             int64  `gorm:"not null unique" example:"1"`      //
	Number            string `gorm:"not null" example:"1"`             //
	Name              string `gorm:"not null" example:"1"`             //
	Mine_info_pub_key []byte `gorm:"not null" example:"1"`             //
	//Mine_info_pri_key        []byte `gorm:"not null" example:"1"`             //
	Mine_info_active_pub_key []byte `gorm:"not null" example:"1"` //
	//Mine_info_active_pri_key []byte `gorm:"not null" example:"1"`             //
	Active_key_id int64  `gorm:"not null" example:"1"` //
	Type          int64  `gorm:"not null" example:"1"` //
	Maxcapacity   int64  `gorm:"not null" example:"1"` //
	Capacity      int64  `gorm:"not null" example:"1"` //
	Mincapacity   int64  `gorm:"not null" example:"1"` //
	Status        int64  `gorm:"not null" example:"1"` //
	Ip            string `gorm:"not null" example:"1"` //
	Location      string `gorm:"not null" example:"1"` //
	Gps           string `gorm:"not null" example:"1"` //
	Version       string `gorm:"not null" example:"1"` //
	Ver           int64  `gorm:"not null" example:"1"` //
	Atime         int64  `gorm:"not null default 0"`   //
	ValidTime     int64  `gorm:"not null default 0"`   //
	Stime         int64  `gorm:"not null default 0"`   //
	Etime         int64  `gorm:"not null default 0"`   //
	Date_created  int64  `gorm:"not null default 0"`   //
}

// TableName returns name of table
func (m MineInfo) TableName() string {
	return `1_mine_info`
}

func (m *MineInfo) GetGuardianNodeCapacity() (int64, error) {
	var ret int64
	var rets string
	if HasTableOrView("1_mine_info") {
		err := conf.GetDbConn().Conn().Table("1_mine_info").Select("COALESCE(SUM(capacity),0)").Where("type = ?", 2).Row().Scan(&rets)
		if err != nil {
			return 0, err
		}
		if rets != "" {
			ret = converter.StrToInt64(rets)
		}
	}
	return ret, nil
}
