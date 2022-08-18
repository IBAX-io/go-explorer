/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-explorer/conf"
	"gorm.io/gorm"
)

// GetDB is returning gorm.DB
func GetDB(db *DbTransaction) *gorm.DB {
	if db != nil && db.conn != nil {
		return db.conn
	}
	return conf.GetDbConn().Conn()
}
