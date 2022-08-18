/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"gorm.io/gorm"
	"strconv"
	"time"
	"unsafe"
)

// LogTransaction is model
type LogTransaction struct {
	Hash         []byte `gorm:"primary_key;not null"`
	Block        int64  `gorm:"not null"`
	Timestamp    int64  `gorm:"not null"`
	ContractName string `gorm:"not null"`
	Address      int64  `gorm:"not null"`
	EcosystemID  int64  `gorm:"not null"`
	Status       int64  `gorm:"not null"`
}

var (
	GLogTranHash map[string]int64
)

// TableName returns name of table
func (m LogTransaction) TableName() string {
	return `log_transactions`
}

// GetByHash returns LogTransactions existence by hash
func (lt *LogTransaction) GetByHash(hash []byte) (bool, error) {
	return isFound(GetDB(nil).Where("hash = ?", hash).First(lt))
}

func (lt *LogTransaction) GetBlockIdByHash(hash []byte) (bool, error) {
	return isFound(GetDB(nil).Select("block").Where("hash = ?", hash).First(lt))
}

func (lt *LogTransaction) GetBlockTransactions(page int, limit int, order string, reqType int) (*[]BlockTxDetailedInfoHex, int64, error) {
	var (
		tss []LogTransaction
		ret []BlockTxDetailedInfoHex
		num int64
		i   int
		j   int32
		err error
	)
	if page < 1 || limit < 1 {
		return &ret, num, err
	}
	err = GetDB(nil).Table(lt.TableName()).Count(&num).Error
	if err != nil {
		return &ret, num, err
	}

	err = GetDB(nil).Select("hash,block").Offset((page - 1) * limit).Limit(limit).Order("block desc").Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}

	if num < 1 {
		return &ret, num, err
	}
	TBlock := make(map[string]int64)
	Thash := make(map[string]bool)
	for i = 0; i < len(tss); i++ {
		hash := hex.EncodeToString(tss[i].Hash)
		Thash[hash] = true

		key := strconv.FormatInt(tss[i].Block, 10)
		TBlock[key] = tss[i].Block
	}

	var Blocks []int64
	for _, k := range TBlock {
		Blocks = append(Blocks, k)
	}

	quickSort(Blocks, 0, int64(len(Blocks)-1))

	for i = len(Blocks); i > 0; i-- {
		bk := &Block{}
		found, err := bk.GetId(Blocks[i-1])
		if err == nil && found {
			rt, err := GetBlocksDetailedInfoHex(bk)
			if err != nil {
				return nil, 0, err
			}
			for j = 0; j < rt.Tx; j++ {
				bh := BlockTxDetailedInfoHex{}
				bh.BlockID = rt.Header.BlockId
				bh.ContractName = rt.Transactions[j].ContractName
				if bh.ContractName == "" {
					bh.ContractName = GetTxContractNameByHash(converter.HexToBin(rt.Transactions[j].Hash))
				}
				bh.Hash = rt.Transactions[j].Hash
				bh.KeyID = rt.Transactions[j].KeyID
				//bh.Params = rt.Transactions[j].Params
				bh.Time = rt.Transactions[j].Time
				bh.Type = rt.Transactions[j].Type
				bh.Ecosystem = rt.Transactions[j].Ecosystem
				bh.Ecosystemname = rt.Transactions[j].Ecosystemname
				if bh.Ecosystem == 1 {
					bh.Token_symbol = SysTokenSymbol
					if bh.Ecosystemname == "" {
						bh.Ecosystemname = SysEcosystemName
					}
				} else {
					bh.Token_symbol = rt.Transactions[j].TokenSymbol
				}
				Ten := unsafe.Sizeof(rt.Transactions[j])
				bh.Size = int64(Ten)
				if Thash[rt.Transactions[j].Hash] {
					var his History
					det, err := his.GetTxListExplorer(converter.HexToBin(rt.Transactions[j].Hash))
					if err != nil {
						return nil, 0, err
					}
					bh.GasFee = det.GasFee
					bh.Amount = det.Amount
					bh.Status = det.Status
					ret = append(ret, bh)
				}
			}
		} else {
			if err != nil {
				return nil, 0, err
			}
		}
	}
	return &ret, num, err
}

func SendDashboardDataToWebsocket(data any, cmd string) error {
	dat := ResponseDashboardTitle{}
	dat.Cmd = cmd
	dat.List = data

	return sendChannelDashboardData(dat)
}

func sendChannelDashboardData(data ResponseDashboardTitle) error {
	ds, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = WriteChannelByte(ChannelDashboard, ds)
	if err != nil {
		return err
	}
	return nil
}

func quickSort(arr []int64, start, end int64) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}

		if start < j {
			quickSort(arr, start, j)
		}
		if end > i {
			quickSort(arr, i, end)
		}
	}
}

func InitDashboardTx() error {
	rd := RedisParams{
		Key:   "dashboard-tx",
		Value: "",
	}
	if err := rd.Del(); err != nil {
		return err
	}
	return nil
}

func GetTransactionBlockFromRedis() (*[]BlockTxDetailedInfoHex, int64, error) {
	var rets HashTransactionResult
	var ret []BlockTxDetailedInfoHex
	rd := RedisParams{
		Key:   "dashboard-tx",
		Value: "",
	}
	var num int64
	//if err := GetDB(nil).Model(LogTransaction{}).Order("block desc").Count(&num).Error; err != nil {
	//	return nil, 0, err
	//}

	if err := rd.Get(); err != nil {
		return nil, 0, err
	}
	if err := json.Unmarshal([]byte(rd.Value), &rets); err != nil {
		return nil, 0, err
	}
	ret = rets.Rets
	num = rets.Total
	return &ret, num, nil
}

func getTransactionBlockToRedis() error {
	var ret HashTransactionResult
	rets, total, err := Get_Group_TransactionBlock(1, 10, "", 1)
	if err != nil {
		return err
	}
	ret.Total = total
	ret.Rets = *rets

	value, err := json.Marshal(ret)
	if err != nil {
		return err
	}

	rd := RedisParams{
		Key:   "dashboard-tx",
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		return err
	}
	if err := SendTransactionListToWebsocket(&ret.Rets); err != nil {
		return fmt.Errorf("send transaction list err:%s", err.Error())
	}

	return nil
}

func SendTransactionListToWebsocket(ret *[]BlockTxDetailedInfoHex) error {
	err := SendTransactionList(ret)
	if err != nil {
		return err
	}
	return nil
}

func SendTransactionList(txList *[]BlockTxDetailedInfoHex) error {
	err := SendDashboardDataToWebsocket(txList, ChannelBlockTransactionList)
	if err != nil {
		return err
	}
	return nil
}

func (lt *LogTransaction) getTransactionIdFromHash(hash []byte) (bool, error) {
	f, err := isFound(GetDB(nil).Select("block").Where("hash = ?", hash).First(&lt))
	if f && err == nil {
		return f, err
	}
	return f, err
}

func (lt *LogTransaction) getHashListByBlockId(blockId int64) ([]LogTransaction, error) {
	var rets []LogTransaction
	if err := GetDB(nil).Select("hash").Where("block = ?", blockId).Find(&rets).Error; err != nil {
		return nil, err
	}
	return rets, nil
}

func (lt *LogTransaction) UnmarshalTxTransaction(txData []byte) (*TxDetailedInfoHex, error) {
	if len(txData) == 0 {
		return nil, errors.New("tx data length is empty")
	}
	var result TxDetailedInfoHex
	tx, err := UnmarshallTransaction(bytes.NewBuffer(txData))
	if err != nil {
		return &result, err
	}

	txDetailedInfo := TxDetailedInfoHex{
		Hash: hex.EncodeToString(tx.Hash()),
	}

	if tx.IsSmartContract() {
		if tx.SmartContract().TxSmart.UTXO != nil {
			//TODO ADD
		} else if tx.SmartContract().TxSmart.TransferSelf != nil {
			//TODO ADD
		} else {
			txDetailedInfo.ContractName, txDetailedInfo.Params = GetMineParam(tx.SmartContract().TxSmart.EcosystemID, tx.SmartContract().TxContract.Name, tx.SmartContract().TxData, tx.Hash())
		}
		txDetailedInfo.KeyID = converter.AddressToString(tx.KeyID())
		txDetailedInfo.Time = MsToSeconds(tx.Timestamp())
		txDetailedInfo.Type = int64(tx.Type())
		txDetailedInfo.Size = int64(len(tx.FullData))
	}

	if txDetailedInfo.Time == 0 {
		txDetailedInfo.Time = MsToSeconds(lt.Timestamp)
	}
	if txDetailedInfo.KeyID == "" {
		txDetailedInfo.KeyID = converter.AddressToString(lt.Address)
	}

	if tx.IsSmartContract() {
		txDetailedInfo.Ecosystem = tx.SmartContract().TxSmart.EcosystemID
		if txDetailedInfo.Ecosystem == 0 {
			txDetailedInfo.Ecosystem = 1
		}
		txDetailedInfo.TokenSymbol, txDetailedInfo.Ecosystemname = GetEcosystemTokenSymbol(txDetailedInfo.Ecosystem)
	} else {
		if txDetailedInfo.Ecosystem == 0 {
			txDetailedInfo.Ecosystem = 1
		}
		if txDetailedInfo.Ecosystem == 1 || txDetailedInfo.Ecosystem == 0 {
			txDetailedInfo.TokenSymbol = SysTokenSymbol
			if txDetailedInfo.Ecosystemname == "" {
				txDetailedInfo.Ecosystemname = SysEcosystemName
			}
		}
	}
	result = txDetailedInfo
	return &result, nil
}

func SearchHash(hash string) (SearchHashResponse, error) {
	var (
		rets SearchHashResponse
	)
	var lt LogTransaction
	hashHex, err := hex.DecodeString(hash)
	if err != nil {
		return rets, err
	}
	f, err := lt.GetByHash(hashHex)
	if err != nil {
		return rets, err
	}
	if !f {
		var item NftMinerItems
		f, err = item.GetByTokenHash(hash)
		if err != nil {
			return rets, err
		}
		if !f {
			return rets, errors.New("doesn't not hash")
		}
	} else {
		rets.IsTxHash = true
	}

	return rets, nil
}

func (lt *LogTransaction) GetEcosystemTransactionFind(ecosystem int64, page, limit int, order, search string, where map[string]any) (*[]EcoTxListResponse, int64, error) {
	var (
		txList    []EcoTxListResponse
		total     int64
		q         *gorm.DB
		startTime time.Time
		endTime   time.Time
	)
	if order == "" {
		order = "block desc"
	}
	if search == "chart" {
		tz := time.Unix(GetNowTimeUnix(), 0)
		endTime = time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
		const getDays = 15
		t1 := endTime.AddDate(0, 0, -1*getDays)
		startTime = t1.AddDate(0, 0, 1)
		//fmt.Printf("start:%d,end:%d\n", startTime.Unix(), endTime.Unix())
	}
	if len(where) == 0 {
		if search == "chart" {
			q = GetDB(nil).Table(lt.TableName()+" as lg").Where("ecosystem_id = ? AND timestamp >= ? AND timestamp < ?", ecosystem, startTime.UnixMilli(), endTime.AddDate(0, 0, 1).UnixMilli())
		} else {
			q = GetDB(nil).Table(lt.TableName()+"  as lg").Where("ecosystem_id = ?", ecosystem)
		}
		err := q.Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = q.Select(`hash,block,timestamp,contract_name,address,status`).
			Order(order).Offset((page - 1) * limit).Limit(limit).Find(&txList).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		where["ecosystem_id ="] = ecosystem
		cond, vals, err := WhereBuild(where)
		if err != nil {
			return nil, 0, err
		}
		if search == "chart" {
			q = GetDB(nil).Table(lt.TableName()+"  as lg").Where(cond, vals...).Where("timestamp >= ? AND timestamp < ?", startTime.UnixMilli(), endTime.AddDate(0, 0, 1).UnixMilli())
		} else {
			q = GetDB(nil).Table(lt.TableName()+"  as lg").Where(cond, vals...)
		}
		err = q.Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = q.Select(`hash,block,timestamp,contract_name,address,status`).
			Order(order).Offset((page - 1) * limit).Limit(limit).Find(&txList).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return &txList, total, nil
}
