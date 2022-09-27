/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/IBAX-io/go-ibax/packages/transaction"
	"github.com/IBAX-io/go-ibax/packages/types"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SpentInfoHistory struct {
	Id          int64  `gorm:"primary_key;not null"`
	Block       int64  `gorm:"column:block;not null"`
	Hash        []byte `gorm:"column:hash;not null"`
	SenderId    int64  `gorm:"column:sender_id;not null"`
	RecipientId int64  `gorm:"column:recipient_id;not null"`
	Amount      string `gorm:"column:amount;type:decimal(40);default:'0';not null"`
	CreatedAt   int64  `gorm:"column:created_at;not null"`
	Ecosystem   int64  `gorm:"not null"`
	Type        int    `gorm:"not null"` //1:UTXO_Transfer 2:UTXO_Tx
}

type spentInfoTxData struct {
	OutputTxHash []byte
	BlockId      int64
	Time         int64
	Data         []byte
}

type utxoTxInfo struct {
	UtxoType    string
	SenderId    int64
	RecipientId int64
	Amount      string
	Ecosystem   int64
}

var getUtxoTxData chan bool

const (
	FeesType    = "fees"
	TaxesType   = "taxes"
	StartUpType = "startUp"
)

func (p *SpentInfoHistory) TableName() string {
	return "spent_info_history"
}

func (p *SpentInfoHistory) CreateTable() (err error) {
	err = nil
	if !HasTableOrView(p.TableName()) {
		if err = GetDB(nil).Migrator().CreateTable(p); err != nil {
			return err
		}
	}
	return err
}

func (p *SpentInfoHistory) GetLast() (bool, error) {
	return isFound(GetDB(nil).Last(p))
}

func (p *SpentInfoHistory) RollbackTransaction() error {
	return GetDB(nil).Where("block > ?", p.Block).Delete(&SpentInfoHistory{}).Error
}

func InitSpentInfoHistory() error {
	var p SpentInfoHistory
	err := p.CreateTable()
	if err != nil {
		return err
	}
	go utxoDataSyncSignalReceive()

	return nil
}

func utxoDataSyncSignalReceive() {
	if getUtxoTxData == nil {
		getUtxoTxData = make(chan bool)
	}
	for {
		select {
		case <-getUtxoTxData:
			err := utxoTxSync()
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("utxo tx Sync Failed")
			}
		}
	}
}

func SendUtxoTxSyncSignal() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	select {
	case getUtxoTxData <- true:
	default:
		//fmt.Printf("Get utxo tx Data len:%d,cap:%d\n", len(getUtxoTxData), cap(getUtxoTxData))
	}
}

func utxoTxSync() error {
	var insertData []SpentInfoHistory
	var (
		si SpentInfo
		st SpentInfo
	)

	tr := &SpentInfoHistory{}
	_, err := tr.GetLast()
	if err != nil {
		return fmt.Errorf("[utxo sync]get spent info history last failed:%s", err.Error())
	}
	f, err := si.GetLast()
	if err != nil {
		return fmt.Errorf("[utxo sync]get spent info last failed:%s", err.Error())
	}
	if f {
		if tr.Block >= si.BlockId {
			utxoTxCheck(si.BlockId)
			return nil
		}
	}

	f, err = st.GetFirst(tr.Block)
	if err != nil {
		return fmt.Errorf("[utxo sync]get spent info block:%d first failed:%s", tr.Block, err.Error())
	}
	if !f {
		return nil
	}

	txList, err := getSpentInfoHashList(st.BlockId, st.BlockId+100)
	if err != nil {
		return fmt.Errorf("[utxo sync]get spent info hash list failed:%s", err.Error())
	}
	if txList == nil {
		return nil
	}

	for _, val := range *txList {
		var (
			data       SpentInfoHistory
			outputList []SpentInfo
		)
		info, err := val.UnmarshalTransaction()
		if err != nil {
			return fmt.Errorf("[utxo sync]unmarshal utxo transaction failed:%s", err.Error())
		}
		data.CreatedAt = val.Time
		data.Hash = val.OutputTxHash
		data.Block = val.BlockId

		_, outputList, err = si.GetOutputs(val.OutputTxHash)
		if err != nil {
			return fmt.Errorf("[utxo sync]get out puts failed:%s", err.Error())
		}

		if info.UtxoType == UtxoTx {
			var (
				index       int
				indexSet    bool
				ecoCount    int
				ecoGasExist bool
			)

			for _, v := range outputList {
				if v.Ecosystem != 1 {
					ecoCount += 1
				}
			}
			if ecoCount >= 3 {
				ecoGasExist = true
			}

			for _, v := range outputList {
				amount, _ := decimal.NewFromString(v.OutputValue)
				recipientId := v.OutputKeyId
				if info.Ecosystem == 1 {
					if v.Ecosystem == 1 {
						switch index {
						case 0:
							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = 1
							data.Type = getSpentInfoHistoryType(FeesType)
							insertData = append(insertData, data)

							index += 1
						case 1:
							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = 1
							data.Type = getSpentInfoHistoryType(TaxesType)
							insertData = append(insertData, data)

							index += 1
						case 2:
							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = 1
							data.Type = getSpentInfoHistoryType(UtxoTx)
							insertData = append(insertData, data)

							index += 1
						case 3:
						}
					}
				} else {
					if v.Ecosystem == 1 {
						switch index {
						case 0:
							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = 1
							data.Type = getSpentInfoHistoryType(FeesType)
							insertData = append(insertData, data)

							index += 1
						case 1:
							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = 1
							data.Type = getSpentInfoHistoryType(TaxesType)
							insertData = append(insertData, data)

							index += 1
						case 2:
						}
					} else {
						if !indexSet {
							if ecoGasExist {
								index = 0
							} else {
								index = 2
							}
							indexSet = true
						}
						switch index {
						case 0:
							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = v.Ecosystem
							data.Type = getSpentInfoHistoryType(FeesType)
							insertData = append(insertData, data)

							index += 1
						case 1:

							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = v.Ecosystem
							data.Type = getSpentInfoHistoryType(TaxesType)
							insertData = append(insertData, data)

							index += 1
						case 2:
							data.Amount = amount.String()
							data.SenderId = info.SenderId
							data.RecipientId = recipientId
							data.Ecosystem = v.Ecosystem
							data.Type = getSpentInfoHistoryType(UtxoTx)
							insertData = append(insertData, data)

							index += 1
						case 3:
						}
					}
				}
			}
		} else if info.UtxoType == StartUpType {
			var lt LogTransaction
			f, err := lt.GetByHash(val.OutputTxHash)
			if err != nil {
				return err
			}
			if !f {
				return fmt.Errorf("[utxo sync]get log hash doesn't exist hash:%s", hex.EncodeToString(val.OutputTxHash))
			}

			data.Type = getSpentInfoHistoryType(info.UtxoType)
			data.SenderId = 5555
			data.RecipientId = lt.Address
			data.Amount = decimal.New(consts.FounderAmount, int32(consts.MoneyDigits)).String()
			data.Ecosystem = lt.EcosystemID

			insertData = append(insertData, data)
		} else {
			data.Type = getSpentInfoHistoryType(info.UtxoType)
			data.SenderId = info.SenderId
			data.RecipientId = info.RecipientId
			data.Amount = info.Amount
			data.Ecosystem = info.Ecosystem

			insertData = append(insertData, data)
		}
		if len(insertData) > 5000 {
			err = createUtxoTxBatches(GetDB(nil), &insertData)
			if err != nil {
				return err
			}
			insertData = nil
		}
	}
	err = createUtxoTxBatches(GetDB(nil), &insertData)
	if err != nil {
		return err
	}

	return utxoTxSync()
}

func createUtxoTxBatches(dbTx *gorm.DB, data *[]SpentInfoHistory) error {
	if data == nil {
		return nil
	}
	return dbTx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(data, 1000).Error
}

func (si *spentInfoTxData) UnmarshalTransaction() (*utxoTxInfo, error) {
	var (
		block = &types.BlockData{}
	)
	blockBuffer := bytes.NewBuffer(si.Data)
	if err := block.UnmarshallBlock(blockBuffer.Bytes()); err != nil {
		return nil, err
	}

	var result utxoTxInfo

	for i := 0; i < len(block.TxFullData); i++ {
		tx, err := transaction.UnmarshallTransaction(bytes.NewBuffer(block.TxFullData[i]))
		if err != nil {
			return nil, err
		}

		if hex.EncodeToString(tx.Hash()) == hex.EncodeToString(si.OutputTxHash) {

			if tx.IsSmartContract() {
				result.Ecosystem = tx.SmartContract().TxSmart.Header.EcosystemID
				if tx.SmartContract().TxSmart.UTXO != nil {
					result.SenderId = tx.KeyID()
					result.UtxoType = UtxoTx
				} else if tx.SmartContract().TxSmart.TransferSelf != nil {
					result.UtxoType = UtxoTransfer
					result.SenderId = tx.KeyID()
					result.RecipientId = tx.KeyID()
					result.Amount = tx.SmartContract().TxSmart.TransferSelf.Value
				} else {
					return &result, errors.New("doesn't not UTXO transaction")
				}
			} else {
				if si.BlockId == 1 {
					result.UtxoType = StartUpType
					return &result, nil
				}
				return nil, errors.New("doesn't not Smart Contract")
			}
		}
	}
	return &result, nil
}

func utxoTxCheck(lastBlockId int64) {
	tx := &SpentInfoHistory{}
	f, err := tx.GetLast()
	if err == nil && f {
		logTran := &LogTransaction{}
		f, err = logTran.GetByHash(tx.Hash)
		if err == nil {
			if !f {
				if tx.Block > lastBlockId {
					tx.Block = lastBlockId
				}
				if tx.Block > 0 {
					log.WithFields(log.Fields{"log hash doesn't exist": hex.EncodeToString(tx.Hash), "block": tx.Block}).Info("[utxo tx check] rollback data")
					tx.Block -= 1
					err = tx.RollbackTransaction()
					if err == nil {
						utxoTxCheck(tx.Block)
					} else {
						log.WithFields(log.Fields{"error": err, "block": tx.Block}).Error("[utxo tx check] rollback Failed")
					}
				}
			}
		} else {
			log.WithFields(log.Fields{"error": err, "hash": hex.EncodeToString(tx.Hash)}).Error("[utxo tx check] get log transaction failed")
		}
	}
}

func getUtxoTxBasisGasFee(hash []byte) decimal.Decimal {
	var hi SpentInfoHistory
	gasFee := decimal.Zero
	_, err := isFound(GetDB(nil).Table(hi.TableName()).Select("COALESCE(sum(amount),0)").
		Where("hash = ? AND ecosystem = 1 AND (type = ? OR type = ?)", hash, getSpentInfoHistoryType(FeesType), getSpentInfoHistoryType(TaxesType)).Take(&gasFee))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "hash": hex.EncodeToString(hash)}).Error("get utxo transaction gas fee failed")
	}
	return gasFee
}

func getSpentInfoHistoryType(utxoType string) int {
	if utxoType == UtxoTransfer {
		return 1
	} else if utxoType == UtxoTx {
		return 2
	} else if utxoType == FeesType {
		return 3
	} else if utxoType == TaxesType {
		return 4
	} else {
		return 5
	}
}
