/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package storage

import (
	"fmt"
	"time"

	"github.com/IBAX-io/go-ibax/packages/smart"

	"github.com/IBAX-io/go-ibax/packages/conf/syspar"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

type DatabaseModel struct {
	Enable  bool   `yaml:"enable"`
	DBType  string `yaml:"type"`
	Connect string `yaml:"connect"`
	Name    string `yaml:"name"`
	Ver     string `yaml:"ver"`
	MaxIdle int    `yaml:"max_idle"`
	MaxOpen int    `yaml:"max_open"`
}

func (d *DatabaseModel) GormInit() (err error) {
	dsn := fmt.Sprintf("%s TimeZone=UTC", d.Connect)
	DBConn, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		//AllowGlobalUpdate: true,                                  //allow global update
		Logger: logger.Default.LogMode(logger.Silent), // start Logger，show detail log
	})
	if err != nil {
		return err
	}
	sqlDB, err := DBConn.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqldb.DBConn = DBConn
	if err = syspar.SysUpdate(nil); err != nil {
		return err
	}
	smart.InitVM()

	// Stats returns database statistics.
	//sqlDB.Stats()
	return nil

}

func (d *DatabaseModel) Conn() *gorm.DB {
	return DBConn
}

func (d *DatabaseModel) Close() error {
	if DBConn != nil {
		sqlDB, err := DBConn.DB()
		if err != nil {
			return err
		}
		if err = sqlDB.Close(); err != nil {
			return err
		}
		DBConn = nil
	}
	return nil
}

func GormDBInit(engine, connect string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connect,
	}), &gorm.Config{
		AllowGlobalUpdate: true,                                  //allow global update
		Logger:            logger.Default.LogMode(logger.Silent), // start Logger，show detail log
	})
	if err != nil {
		return nil, err
	}
	db, err := conn.DB()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	return conn, nil
}
