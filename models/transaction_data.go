/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"encoding/hex"
	"github.com/IBAX-io/go-ibax/packages/transaction"
	"github.com/IBAX-io/go-ibax/packages/types"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionData struct {
	Hash   []byte `gorm:"primary_key;not null"`
	Block  int64  `gorm:"not null"`
	TxData []byte `gorm:"not null"`
}

var getTransactionData chan bool

func (p *TransactionData) TableName() string {
	return "transaction_data"
}

func (p *TransactionData) CreateTable() (err error) {
	err = nil
	if !HasTableOrView(nil, p.TableName()) {
		if err = GetDB(nil).Migrator().CreateTable(p); err != nil {
			return err
		}
	}
	return err
}

func InitTransactionData() error {
	var p TransactionData
	err := p.CreateTable()
	if err != nil {
		return err
	}
	TxDataSyncSignalReceive()

	return nil
}

func (p *TransactionData) GetByHash(hash []byte) (bool, error) {
	return isFound(GetDB(nil).Where("hash = ?", hash).First(p))
}

func (p *TransactionData) GetTxDataByHash(hash []byte) (bool, error) {
	return isFound(GetDB(nil).Select("tx_data").Where("hash = ?", hash).First(p))
}

func (p *TransactionData) GetLast() (bool, error) {
	return isFound(GetDB(nil).Order("block desc").Limit(1).Take(p))
}

func (p *TransactionData) RollbackTransaction() error {
	return GetDB(nil).Where("block > ?", p.Block).Delete(&TransactionData{}).Error
}

func TxDataSyncSignalReceive() {
	if getTransactionData == nil {
		getTransactionData = make(chan bool)
	}
	for {
		select {
		case <-getTransactionData:
			if err := transactionDataSync(); err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Transaction Data Sync Failed")
			}
		}
	}
}

func SendTxDataSyncSignal() {
	select {
	case getTransactionData <- true:
	default:
		//fmt.Printf("Get Transaction Data len:%d,cap:%d\n", len(getTransactionData), cap(getTransactionData))
	}
}

func transactionDataSync() error {
	var insertData []*TransactionData
	var b1 Block

	tr := &TransactionData{}
	_, err := tr.GetLast()
	if err != nil {
		return err
	}
	f, err := b1.GetMaxBlock()
	if err != nil {
		return err
	}
	if f {
		if tr.Block >= b1.ID {
			transactionDataCheck(b1.ID)
			return nil
		}
	}

	bkList, err := GetBlockchain(tr.Block, tr.Block+5000, "asc")
	if err != nil {
		return err
	}
	if bkList == nil {
		return nil
	}
	for _, val := range *bkList {
		txList, err := UnmarshallBlockTxData(bytes.NewBuffer(val.Data))
		if err != nil {
			return err
		}
		for hash, data := range txList {
			var tran TransactionData
			tran.Hash, _ = hex.DecodeString(hash)
			tran.Block = val.ID
			tran.TxData = data
			insertData = append(insertData, &tran)
		}
		if len(insertData) >= 5000 {
			err = createTransactionDataBatches(GetDB(nil), insertData)
			if err != nil {
				return err
			}
			insertData = nil
		}
	}
	if insertData != nil {
		err = createTransactionDataBatches(GetDB(nil), insertData)
		if err != nil {
			return err
		}
	}

	return transactionDataSync()
}

func transactionDataCheck(lastBlockId int64) {
	tran := &TransactionData{}
	f, err := tran.GetLast()
	if err == nil && f {
		logTran := &LogTransaction{}
		f, err = logTran.GetByHash(tran.Hash)
		if err == nil {
			if !f {
				if tran.Block > lastBlockId {
					tran.Block = lastBlockId
				}
				if tran.Block > 0 {
					log.WithFields(log.Fields{"log hash doesn't exist": hex.EncodeToString(tran.Hash), "block": tran.Block}).Info("rollback transaction data")
					tran.Block -= 1
					err = tran.RollbackTransaction()
					if err == nil {
						transactionDataCheck(tran.Block)
					} else {
						log.WithFields(log.Fields{"error": err, "block": tran.Block}).Error("transaction Data rollback Failed")
					}
				}
			}
		} else {
			log.WithFields(log.Fields{"error": err, "hash": hex.EncodeToString(tran.Hash)}).Error("get log transaction failed")
		}
	}
}

func createTransactionDataBatches(dbTx *gorm.DB, data []*TransactionData) error {
	if len(data) == 0 {
		return nil
	}
	return dbTx.Model(&TransactionData{}).Create(&data).Error
}

func UnmarshallBlockTxData(blockBuffer *bytes.Buffer) (map[string][]byte, error) {
	var (
		block = &types.BlockData{}
	)
	if err := block.UnmarshallBlock(blockBuffer.Bytes()); err != nil {
		return nil, err
	}

	txList := make(map[string][]byte)
	for i := 0; i < len(block.TxFullData); i++ {
		tx, err := transaction.UnmarshallTransaction(bytes.NewBuffer(block.TxFullData[i]))
		if err != nil {
			return nil, err
		}

		txList[hex.EncodeToString(tx.Hash())] = block.TxFullData[i]
	}
	return txList, nil
}

func GetTxContractNameByHash(hash []byte) string {
	tr := &TransactionData{}
	f, err := tr.GetByHash(hash)
	if err != nil || !f {
		return ""
	}

	if len(tr.TxData) == 0 {
		return ""
	}
	tx, err := UnmarshallTransaction(bytes.NewBuffer(tr.TxData))
	if err != nil {
		return ""
	}
	if tx.IsSmartContract() {
		if tx.SmartContract().TxSmart.UTXO != nil {
			return "UTXO_Tx"
		} else if tx.SmartContract().TxSmart.TransferSelf != nil {
			return "UTXO_Transfer"
		}
	}
	return ""
}

func UnmarshallTransaction(blockBuffer *bytes.Buffer) (*transaction.Transaction, error) {
	return transaction.UnmarshallTransaction(blockBuffer)
}
