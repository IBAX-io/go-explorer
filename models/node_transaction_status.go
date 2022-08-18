/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"github.com/IBAX-io/go-explorer/conf"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

var (
	GNodeStatusTranHash map[string]TransactionStatus
)

func (ts *TransactionStatus) GetNodecount(db *gorm.DB) (int64, error) {
	var (
		count int64
	)
	err := db.Table("transactions_status").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

func (ts *TransactionStatus) DbconngetSqlite(transactionHash []byte) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("hash = ?", transactionHash).First(ts))
}

func (ts *TransactionStatus) DBconnGetcount_Sqlite() (int64, error) {
	var (
		count int64
	)
	err := conf.GetDbConn().Conn().Table("transactions_status").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

func DbconnbatchupdateSqlite(objarr *[]TransactionStatus) error {
	for _, val := range *objarr {
		err := conf.GetDbConn().Conn().Model(&TransactionStatus{}).Where("hash = ?", val.Hash).Updates(val).Error
		if err != nil {
			log.Info("TransactionStatus update false: "+err.Error()+" data:", val)
			return err
		}
	}
	return nil
}

func (ts *TransactionStatus) DbconnbatchinsertSqlites(objArr *[]TransactionStatus) error {
	ret, ret1 := DbconndealReduplictionTransactionstatus(objArr)
	count := len(*ret)
	if len(*ret) != 0 {
		dat := *ret
		for i := 0; i < count; {
			if i+100 < count {
				s := dat[i : i+100]
				err := DbconnbatchupdateSqlite(&s)
				if err != nil {
					log.Info("node TransactionStatus update count err: " + err.Error())
				}
				i += 100
			} else {
				s := dat[i:]
				err := DbconnbatchupdateSqlite(&s)
				if err != nil {
					log.Info("node TransactionStatus update count err: " + err.Error())
				}
				i = count
			}
		}
	}

	if len(*ret1) != 0 {
		DbconnbatchupdateSqlite(ret1)
	}
	return nil
}

func DbconndealReduplictionTransactionstatus(objArr *[]TransactionStatus) (*[]TransactionStatus, *[]TransactionStatus) {
	var (
		ret  []TransactionStatus
		ret1 []TransactionStatus
	)
	if GNodeStatusTranHash == nil {
		GNodeStatusTranHash = make(map[string]TransactionStatus)
	}
	for _, val := range *objArr {
		key := hex.EncodeToString(val.Hash)
		dat, ok := GNodeStatusTranHash[key]
		if ok {
			if val.Error != dat.Error || val.BlockID != dat.BlockID {
				ret1 = append(ret1, val)
				GNodeStatusTranHash[key] = val
			}
		} else {
			GNodeStatusTranHash[key] = val
			ret = append(ret, val)
		}
	}
	return &ret, &ret1
}
