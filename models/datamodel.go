/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import "time"

//redisModel get redis information from config.yml
type TableModel struct {
	ID           int64     `gorm:"primary_key;not_null" json:"id" `
	NodePosition int64     `gorm:"not null" json:"nodeposition"`
	Enable       bool      `gorm:"not null" json:"enable"`
	Fmode        int       `gorm:"not null" json:"fmode"` //  <0      1  id    2    3
	Smode        bool      `gorm:"not null" json:"smode"` //  false no flash   true flash  memory
	Name         string    `gorm:"not null" json:"name"`
	Time         time.Time `gorm:"not null" json:"time"`
	Dataid       int64     `gorm:"not null" json:"id"`
	Count        int       `gorm:"not null" json:"count"`
	Updatetime   int       `gorm:"not null" json:"updatetime"`
	//Cmdsql string    `gorm:"not null" json:"cmdsql"`
}

type TableShowModel struct {
	ID     int    `gorm:"primary_key;not_null" json:"id" `
	Name   string `gorm:"not null" json:"name"`
	Cmdsql string `gorm:"not null" json:"_"`
}
