/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/IBAX-io/go-explorer/conf"
	log "github.com/sirupsen/logrus"
	//"github.com/IBAX-io/go-explorer/conf"
)

// TransactionStatus is model
type TransactionStatus struct {
	Hash      []byte `gorm:"primary_key;not null"  json:"hash"`
	Time      int64  `gorm:"not null" json:"time"`
	Type      int64  `gorm:"not null"  json:"type"`
	Ecosystem int64  `gorm:"not null"  json:"ecosystem"`
	WalletID  int64  `gorm:"not null"  json:"wallet_id"`
	BlockID   int64  `gorm:"not null;index:tsblockid_idx"  json:"block_id"`
	Error     string `gorm:"not null"  json:"error"`
	Penalty   int64  `gorm:"not null"  json:"penalty"`
}

var (
	GStatusTranHash map[string]TransactionStatus
)

// TableName returns name of table
func (ts *TransactionStatus) TableName() string {
	return "transactions_status"
}

func (ts *TransactionStatus) Create() error {
	return conf.GetDbConn().Conn().Create(ts).Error
}

func (ts *TransactionStatus) Get(transactionHash []byte) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("hash = ?", transactionHash).First(ts))
}

func (ts *TransactionStatus) Getcount() (int64, error) {
	var (
		count int64
	)
	err := conf.GetDbConn().Conn().Table("transactions_status").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

func (ts *TransactionStatus) GetTransactions(page int, size int, order string) (*[]TransactionStatusHex, int64, error) {
	var (
		tss []TransactionStatus
		ret []TransactionStatusHex
		num int64
	)

	err := conf.GetDbConn().Conn().Limit(size).Offset((page - 1) * size).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
	err = conf.GetDbConn().Conn().Table("transactions_status").Count(&num).Error
	if err != nil {
		return &ret, num, err
	}

	for i := 0; i < len(tss); i++ {

		var tx = TransactionStatusHex{
			Hash:      hex.EncodeToString(tss[i].Hash),
			Time:      tss[i].Time,
			Type:      tss[i].Type,
			Ecosystem: tss[i].Ecosystem,
			WalletID:  strconv.FormatInt(tss[i].WalletID, 10),
			BlockID:   tss[i].BlockID,
			Error:     tss[i].Error,
			Penalty:   tss[i].Penalty,
		}
		tx.TokenSymbol, tx.Ecosystemname = Tokens.Get(tss[i].Ecosystem), EcoNames.Get(tss[i].Ecosystem)
		ret = append(ret, tx)
	}
	return &ret, num, err
}

func (ts *TransactionStatus) GetTimelimit(time time.Time) (*[]TransactionStatus, error) {
	var (
		tss []TransactionStatus
	)

	err := conf.GetDbConn().Conn().Where("time >= ?", time.Unix()).Order("time desc").Find(&tss).Error
	if err != nil {
		return nil, err
	}

	return &tss, err
}

func (ts *TransactionStatus) Create_Sqlite() error {
	return conf.GetDbConn().Conn().Create(ts).Error
}

func (ts *TransactionStatus) Get_Sqlite(transactionHash []byte) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("hash = ?", transactionHash).First(ts))
}

func (ts *TransactionStatus) Getcount_Sqlite() (int64, error) {
	var (
		count int64
	)
	err := conf.GetDbConn().Conn().Table("transactions_status").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

func (ts *TransactionStatus) GetTimelimit_Sqlite(time time.Time) (*[]TransactionStatus, error) {
	var (
		tss []TransactionStatus
	)
	err := conf.GetDbConn().Conn().Where("time >= ?", time.Unix()).Find(&tss).Order("time").Error
	if err != nil {
		return nil, err
	}

	return &tss, err
	//return isFound(conf.GetDbConn().Conn().Where("time >= ?", time).Find(ts))
}

func StUpdate_Sqlite(objarr *[]TransactionStatus) error {

	for _, val := range *objarr {
		err := conf.GetDbConn().Conn().Model(&TransactionStatus{}).Updates(val).Error
		if err != nil {
			log.Info("TransactionStatus update false: "+err.Error(), val)
		}
	}
	return nil
}

func (ts *TransactionStatus) BatchUpdate_Sqlite(reportForms *[]TransactionStatus) error {
	for _, val := range *reportForms {
		err := conf.GetDbConn().Conn().Create(&val).Error
		if err != nil {
			log.Info("conf.GetDbConn().Conn().NewRecord false: " + err.Error())
			//
			conf.GetDbConn().Conn().Save(&val)
		}
	}
	return nil
}

func (ts *TransactionStatus) BatchinsertSqlite(objArr *[]TransactionStatus) error {
	if len(*objArr) == 0 {
		return nil
	}
	//dat := *objArr
	//mainObj := dat[0]
	return conf.GetDbConn().Conn().Create(objArr).Error
	//
	//mainScope := conf.GetDbConn().Conn().NewScope(mainObj)
	//mainFields := mainScope.Fields()
	//quoted := make([]string, 0, len(mainFields))
	//vquoted := make([]string, 0, len(mainFields))
	//upstr := " ON DUPLICATE KEY UPDATE  "
	//for i := range mainFields {
	//	// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
	//	// If field is ignore field, skip it.
	//	if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
	//		continue
	//	}
	//	quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	//	vquoted = append(vquoted, mainScope.Quote(mainFields[i].DBName)+"=VALUES("+mainScope.Quote(mainFields[i].DBName)+")")
	//	//upstr += mainFields[i].DBName +"=VALUES("+mainFields[i].DBName+"),"
	//}
	//
	//placeholdersArr := make([]string, 0, len(*objArr))
	//
	//for _, obj := range *objArr {
	//	scope := conf.GetDbConn().Conn().NewScope(obj)
	//	fields := scope.Fields()
	//	placeholders := make([]string, 0, len(fields))
	//	for i := range fields {
	//		if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
	//			continue
	//		}
	//		placeholders = append(placeholders, scope.AddToVars(fields[i].Field.Interface()))
	//	}
	//	placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
	//	placeholdersArr = append(placeholdersArr, placeholdersStr)
	//	// add real variables for the replacement of placeholders' '?' letter later.
	//	mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
	//}
	//
	//mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ",
	//	mainScope.QuotedTableName(),
	//	strings.Join(quoted, ", "),
	//	strings.Join(placeholdersArr, ", "),
	//))
	//
	//upstr += strings.Join(vquoted, ", ")
	//if _, err := mainScope.SQLDB().Exec(mainScope.SQL, mainScope.SQLVars...); err != nil {
	//	return err
	//}
	//return nil
}

func (ts *TransactionStatus) BatchInsert_Sqlites(objArr *[]TransactionStatus) error {
	ret, ret1 := Deal_Redupliction_Transactionstatus(objArr)
	count := len(*ret)
	if count > 0 {
		dat := *ret
		for i := 0; i < count; {
			if i+100 < count {
				s := dat[i : i+100]
				err := ts.BatchinsertSqlite(&s)
				if err != nil {
					StUpdate_Sqlite(&s)
					log.Info("BatchInsert_Sqlite update count err: " + err.Error())
				}
				i += 100
			} else {
				s := dat[i:]
				err := ts.BatchinsertSqlite(&s)
				if err != nil {
					StUpdate_Sqlite(&s)
					log.Info("BatchInsert_Sqlite update count err: " + err.Error())
				}
				i = count
			}

		}
	}
	if len(*ret1) > 0 {
		StUpdate_Sqlite(ret1)
	}

	return nil
}

func Deal_Redupliction_Transactionstatus(objArr *[]TransactionStatus) (*[]TransactionStatus, *[]TransactionStatus) {
	var (
		ret  []TransactionStatus
		ret1 []TransactionStatus
	)
	if GStatusTranHash == nil {
		GStatusTranHash = make(map[string]TransactionStatus)
	}
	for _, val := range *objArr {
		key := hex.EncodeToString(val.Hash)
		dat, ok := GStatusTranHash[key]
		if ok {
			if dat.BlockID > 0 {

			} else if val.Error != dat.Error || val.BlockID != dat.BlockID {
				//
				GStatusTranHash[key] = val
				ret1 = append(ret1, val)
			}

		} else {
			GStatusTranHash[key] = val
			ret = append(ret, val)
		}
	}
	return &ret, &ret1
}
