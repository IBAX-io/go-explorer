/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package storage

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var lpgdb *gorm.DB

type LDatabaseModel struct {
	Enable  bool   `yaml:"enable"`
	DBType  string `yaml:"type"`
	Connect string `yaml:"connect"`
	Name    string `yaml:"name"`
	Ver     string `yaml:"ver"`
	MaxIdle int    `yaml:"max_idle"`
	MaxOpen int    `yaml:"max_open"`
}

func (d *LDatabaseModel) GormInit() (err error) {
	dsn := fmt.Sprintf("%s TimeZone=UTC", d.Connect)
	lpgdb, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		AllowGlobalUpdate: true,                                  //allow global update
		Logger:            logger.Default.LogMode(logger.Silent), // start Logger，show detail log
	})
	if err != nil {
		return err
	}
	sqlDB, err := lpgdb.DB()
	if err != nil {
		return err
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)

	lpgdb.Migrator().DropTable(&TransactionStatus{})
	lpgdb.Migrator().DropTable(&BlockTxDetailedInfoHex{})
	lpgdb.Migrator().AutoMigrate(&TransactionStatus{})
	lpgdb.Migrator().AutoMigrate(&BlockTxDetailedInfoHex{})

	return nil
}

func (d *LDatabaseModel) Conn() *gorm.DB {
	return lpgdb
}

func (d *LDatabaseModel) Close() error {
	if lpgdb != nil {
		sqlDB, err := lpgdb.DB()
		if err != nil {
			return err
		}
		if err = sqlDB.Close(); err != nil {
			return err
		}
		lpgdb = nil
	}
	return nil
}
