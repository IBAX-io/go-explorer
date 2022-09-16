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
	"gorm.io/gorm/clause"
)

type TransactionData struct {
	Hash      []byte `gorm:"primary_key;not null"`
	Block     int64  `gorm:"not null"`
	TxData    []byte `gorm:"not null"`
	Amount    string `gorm:"column:amount;type:decimal(40);default:'0';not null"` //the transaction occurs ecosystem generates transaction amount
	Ecosystem int64  `gorm:"not null"`
	TxTime    int64  `gorm:"not null"`
}

var getTransactionData chan bool

func (p *TransactionData) TableName() string {
	return "transaction_data"
}

func (p *TransactionData) CreateTable() (err error) {
	err = nil
	if !HasTableOrView(p.TableName()) {
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
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	select {
	case getTransactionData <- true:
	default:
		//fmt.Printf("Get Transaction Data len:%d,cap:%d\n", len(getTransactionData), cap(getTransactionData))
	}
}

func transactionDataSync() error {
	var insertData []TransactionData
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

	bkList, err := GetBlockchain(tr.Block, tr.Block+100, "asc")
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
		for _, data := range txList {
			data.Block = val.ID
			if data.TxTime == 0 {
				data.TxTime = val.Time
			}
			if data.Ecosystem == 0 {
				data.Ecosystem = 1
			}
			insertData = append(insertData, data)
		}
	}
	err = createTransactionDataBatches(GetDB(nil), &insertData)
	if err != nil {
		return err
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

func createTransactionDataBatches(dbTx *gorm.DB, data *[]TransactionData) error {
	if data == nil {
		return nil
	}
	return dbTx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(data, 1000).Error
}

func UnmarshallBlockTxData(blockBuffer *bytes.Buffer) (map[string]TransactionData, error) {
	var (
		block = &types.BlockData{}
	)
	if err := block.UnmarshallBlock(blockBuffer.Bytes()); err != nil {
		return nil, err
	}

	txList := make(map[string]TransactionData)
	for i := 0; i < len(block.TxFullData); i++ {
		var info TransactionData

		tx, err := transaction.UnmarshallTransaction(bytes.NewBuffer(block.TxFullData[i]))
		if err != nil {
			return nil, err
		}
		info.TxData = block.TxFullData[i]
		info.Hash = tx.Hash()

		if tx.IsSmartContract() {
			if tx.SmartContract().TxSmart.UTXO != nil {
				info.Amount = tx.SmartContract().TxSmart.UTXO.Value
			} else if tx.SmartContract().TxSmart.TransferSelf != nil {
				//info.Amount = tx.SmartContract().TxSmart.TransferSelf.Value
			} else {
				var his History
				info.Amount = his.GetHashSum(info.Hash, tx.SmartContract().TxSmart.EcosystemID, block.Header.BlockId)
			}
			info.TxTime = MsToSeconds(tx.Timestamp())
			info.Ecosystem = tx.SmartContract().TxSmart.EcosystemID
		}
		txList[hex.EncodeToString(tx.Hash())] = info
	}
	return txList, nil
}

func GetUtxoTxContractNameByHash(hash []byte) string {
	if hash == nil {
		return ""
	}
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
			return UtxoTx
		} else if tx.SmartContract().TxSmart.TransferSelf != nil {
			return UtxoTransfer
		}
	}
	return ""
}

func UnmarshallTransaction(blockBuffer *bytes.Buffer) (*transaction.Transaction, error) {
	return transaction.UnmarshallTransaction(blockBuffer)
}

func IsUtxoTransaction(txData []byte) (bool, error) {
	tx, err := UnmarshallTransaction(bytes.NewBuffer(txData))
	if err != nil {
		return false, err
	}
	if tx.IsSmartContract() {
		if tx.SmartContract().TxSmart.UTXO != nil || tx.SmartContract().TxSmart.TransferSelf != nil {
			return true, nil
		}
	}
	return false, nil
}
