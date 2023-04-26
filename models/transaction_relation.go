/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
)

var getTxRelation chan bool
var txRelationStart bool = true

type TransactionRelation struct {
	Hash         []byte `gorm:"column:hash;not null;index"`
	SenderIds    string `gorm:"column:sender_ids;not null"`
	RecipientIds string `gorm:"column:recipient_ids"`
	Ecosystem    int64  `gorm:"column:ecosystem;not null"`
	Block        int64  `gorm:"column:block;not null;index"`
	CreatedAt    int64  `gorm:"column:created_at;not null"`
}

type txRelationInfo struct {
	Hash         []byte
	Block        int64
	ContractName string
	Address      int64
	EcosystemID  int64
	Timestamp    int64
	TxData       []byte
}

func (p *TransactionRelation) TableName() string {
	return "transaction_relation"
}

func (p *TransactionRelation) CreateIndex() error {
	extensionName := "pg_trgm"
	indexName := "idx_tx_relation_sender_recipient_ids"

	var (
		extname string
		relname string
	)
	f, err := isFound(GetDB(nil).Table("pg_extension").Select("extname").Where("extname = ?", extensionName).Take(&extname))
	if err != nil {
		return err
	}
	if !f {
		err = GetDB(nil).Exec(fmt.Sprintf("CREATE EXTENSION %s", extensionName)).Error
		if err != nil {
			return err
		}
	}
	f, err = isFound(GetDB(nil).Table("pg_stat_user_indexes").Select("relname").Where("relname = ?", indexName).Take(&relname))
	if err != nil {
		return err
	}
	if !f {
		err = GetDB(nil).Exec(fmt.Sprintf(`
CREATE INDEX %s
             ON %s using gin ((sender_ids || recipient_ids) gin_trgm_ops)
`, indexName, p.TableName())).Error
		if err != nil {
			return err
		}
	}
	return nil

}

func (p *TransactionRelation) CreateTable() (err error) {
	err = nil
	if !HasTableOrView(p.TableName()) {
		if err = GetDB(nil).Migrator().CreateTable(p); err != nil {
			return err
		}
	}
	//return p.CreateIndex()
	return nil
}

func (p *TransactionRelation) RollbackTransaction() error {
	return GetDB(nil).Where("block >= ?", p.Block).Delete(&TransactionRelation{}).Error
}

//func (p *TransactionRelation) CreateIndex() error {
//	GetDB(nil).
//}

func (p *TransactionRelation) GetLast() (bool, error) {
	return isFound(GetDB(nil).Order("block desc").Take(p))
}

func (p *TransactionRelation) RollbackOne() error {
	if p.Block > 0 {
		err := p.RollbackTransaction()
		if err != nil {
			log.WithFields(log.Fields{"error": err, "block": p.Block}).Error("[tx relation] rollback one Failed")
			return err
		}
	}
	return nil
}

func createTxRelationBatches(dbTx *gorm.DB, data *[]TransactionRelation) error {
	if data == nil {
		return nil
	}
	return dbTx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(data, 1000).Error
}

func InitTransactionRelation() error {
	var p TransactionRelation
	err := p.CreateTable()
	if err != nil {
		return err
	}
	go txRelationSignalReceive()

	return nil
}

func txRelationSignalReceive() {
	if getTxRelation == nil {
		getTxRelation = make(chan bool)
	}
	for {
		select {
		case <-getTxRelation:
			err := txRelationSync()
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Transaction Relation Sync Failed")
			}
		}
	}
}

func SendTxRelationSignal() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	select {
	case getTxRelation <- true:
	default:
	}
}

func txRelationSync() error {
	var (
		insertData []TransactionRelation
	)
begin:
	tr := &TransactionRelation{}
	_, err := tr.GetLast()
	if err != nil {
		return fmt.Errorf("[tx relation sync]get last failed:%s", err.Error())
	}
	lg := &LogTransaction{}
	f, err := lg.GetLast()
	if err != nil {
		return fmt.Errorf("[tx relation sync]get logtransaction last failed:%s", err.Error())
	}
	insertTxData := func() error {
		err = createTxRelationBatches(GetDB(nil), &insertData)
		if err != nil {
			return fmt.Errorf("insert tx relation batches failed:%s", err.Error())
		}
		insertData = nil
		return nil
	}
	if f {
		if txRelationStart {
			err = tr.RollbackOne()
			if err != nil {
				return err
			} else {
				txRelationStart = false
				goto begin
			}
		}
		if tr.Block >= lg.Block {
			TxRelationCheck(lg.Block)
			return nil
		}
	}

	st := &LogTransaction{}
	f, err = st.GetBlockFirst(tr.Block)
	if err != nil {
		return fmt.Errorf("[tx relation sync]get logtransaction block:%d first failed:%s", tr.Block, err.Error())
	}
	if !f {
		return nil
	}

	txList, err := getTxListByBlockNew(st.Block, st.Block+100)
	if err != nil {
		return fmt.Errorf("[tx relation sync]get tx list failed:%s", err.Error())
	}
	if txList == nil {
		return nil
	}

	for _, tx := range *txList {
		if tx.ContractName != "" {
			relations, err := tx.getTransactionRelation(false)
			if err != nil {
				return fmt.Errorf("[tx relation sync]get transaction relation failed:%s", err.Error())
			}
			insertData = append(insertData, *relations...)
		} else {
			f, err = IsUtxoTransaction(tx.TxData, tx.Block)
			if err != nil {
				return fmt.Errorf("[tx relation sync]unmarshal transaction failed:%s", err.Error())
			}
			relations, err := tx.getTransactionRelation(f)
			if err != nil {
				return fmt.Errorf("[tx relation sync]get transaction relation failed:%s", err.Error())
			}

			insertData = append(insertData, *relations...)
		}

	}
	err = insertTxData()
	if err != nil {
		return err
	}

	return txRelationSync()
}

func (p *txRelationInfo) getTransactionRelation(isUtxo bool) (*[]TransactionRelation, error) {

	var (
		list      []TransactionRelation
		ecosystem int64
	)
	sMap := make(map[int64]string)
	rMap := make(map[int64]string)
	addList := func(ecosystem int64) {
		var (
			info  TransactionRelation
			sList []string
			rList []string
		)
		for eid, v := range sMap {
			if eid == ecosystem {
				sList = append(sList, v)
			}
		}
		for eid, v := range rMap {
			if eid == ecosystem {
				rList = append(rList, v)
			}
		}
		info.SenderIds = strings.Join(sList, ",")
		info.RecipientIds = strings.Join(rList, ",")
		info.Ecosystem = ecosystem
		info.Hash = p.Hash
		info.Block = p.Block
		info.CreatedAt = p.Timestamp

		list = append(list, info)
	}

	if isUtxo {
		var h1 []SpentInfoHistory
		err := GetDB(nil).Select("sender_id,ecosystem,recipient_id").
			Where("hash = ?", p.Hash).Group("sender_id,ecosystem,recipient_id").Order("ecosystem").Find(&h1).Error
		if err != nil {
			return nil, err
		}
		length := len(h1)
		if length == 0 {
			return nil, fmt.Errorf("waiting utxo tx sync,hash:%s", hex.EncodeToString(p.Hash))
		}
		lastIndex := length - 1
		for k, val := range h1 {
			if k == 0 {
				ecosystem = val.Ecosystem
			}
			sender := strconv.FormatInt(val.SenderId, 10)
			recipient := strconv.FormatInt(val.RecipientId, 10)
			if keyStr, ok := sMap[val.Ecosystem]; ok {
				sMap[val.Ecosystem] = keyStr + "," + sender
			} else {
				sMap[val.Ecosystem] = sender
			}
			if keyStr, ok := rMap[val.Ecosystem]; ok {
				rMap[val.Ecosystem] = keyStr + "," + recipient
			} else {
				rMap[val.Ecosystem] = recipient
			}

			if ecosystem != val.Ecosystem {
				addList(ecosystem)
				sMap[ecosystem] = ""
				rMap[ecosystem] = ""
				ecosystem = val.Ecosystem
			}
			if lastIndex == k {
				ecosystem = val.Ecosystem
			}
		}
		if sMap[ecosystem] != "" {
			addList(ecosystem)
		}
		return &list, nil
	}

	var h2 []History
	err := GetDB(nil).Select("sender_id,ecosystem,recipient_id").
		Where("txhash = ? AND block_id = ?", p.Hash, p.Block).Group("sender_id,ecosystem,recipient_id").Order("ecosystem").Find(&h2).Error
	if err != nil {
		return nil, err
	}
	length := len(h2)
	if length == 0 {
		var info TransactionRelation
		info.SenderIds = strconv.FormatInt(p.Address, 10)
		info.Ecosystem = p.EcosystemID
		info.Hash = p.Hash
		info.Block = p.Block
		info.CreatedAt = p.Timestamp

		list = append(list, info)
		return &list, nil
	}
	lastIndex := length - 1

	for k, val := range h2 {
		if k == 0 {
			ecosystem = val.Ecosystem
		}
		sender := strconv.FormatInt(val.Senderid, 10)
		recipient := strconv.FormatInt(val.Recipientid, 10)
		if keyStr, ok := sMap[val.Ecosystem]; ok {
			sMap[val.Ecosystem] = keyStr + "," + sender
		} else {
			sMap[val.Ecosystem] = sender
		}
		if keyStr, ok := rMap[val.Ecosystem]; ok {
			rMap[val.Ecosystem] = keyStr + "," + recipient
		} else {
			rMap[val.Ecosystem] = recipient
		}

		if ecosystem != val.Ecosystem {
			addList(ecosystem)
			sMap[ecosystem] = ""
			rMap[ecosystem] = ""
			ecosystem = val.Ecosystem
		}
		if lastIndex == k {
			ecosystem = val.Ecosystem
		}
	}
	if sMap[ecosystem] != "" {
		addList(ecosystem)
	}

	return &list, nil
}

func TxRelationCheck(lastBlockId int64) {
	tr := &TransactionRelation{}
	f, err := tr.GetLast()
	if err == nil && f {
		logTran := &LogTransaction{}
		f, err = logTran.GetByHash(tr.Hash)
		if err == nil {
			if !f {
				if tr.Block > lastBlockId {
					tr.Block = lastBlockId
				}
				if tr.Block > 0 {
					log.WithFields(log.Fields{"log hash doesn't exist": hex.EncodeToString(tr.Hash), "block": tr.Block}).Info("[tx relation check] rollback data")
					err = tr.RollbackTransaction()
					if err == nil {
						TxRelationCheck(tr.Block)
					} else {
						log.WithFields(log.Fields{"error": err, "block": tr.Block}).Warn("[tx relation check] rollback Failed")
					}
				}
			}
		} else {
			log.WithFields(log.Fields{"error": err, "hash": hex.EncodeToString(tr.Hash)}).Warn("[tx relation check] get log transaction failed")
		}
	}
}

func getTxListByBlock(startId, endId int64) (*[]LogTransaction, error) {
	var err error
	var list []LogTransaction
	err = GetDB(nil).Select("block,hash,address,ecosystem_id,contract_name").
		Where("block >= ? AND block < ?", startId, endId).Order("block asc,timestamp asc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func getTxListByBlockNew(startId, endId int64) (*[]txRelationInfo, error) {
	var err error
	var list []txRelationInfo

	err = GetDB(nil).Raw(`
SELECT lg.block,lg.hash,lg.address,lg.ecosystem_id,lg.contract_name,lg.timestamp,td.tx_data FROM log_transactions AS lg LEFT JOIN 
transaction_data AS td ON(td.hash = lg.hash AND lg.contract_name = '') WHERE lg.block >= ? AND lg.block < ? ORDER BY lg.block asc,lg.timestamp asc
`, startId, endId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}
