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
	"github.com/IBAX-io/go-ibax/packages/types"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"reflect"
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
	Status       int32  `gorm:"not null"`
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

func (lt *LogTransaction) GetContract(hash []byte) (bool, error) {
	return isFound(GetDB(nil).Select("contract_name").Where("hash = ?", hash).First(lt))
}

func (lt *LogTransaction) GetBlockIdByHash(hash []byte) (bool, error) {
	return isFound(GetDB(nil).Select("block").Where("hash = ?", hash).First(lt))
}
func (lt *LogTransaction) GetStatus(hash []byte) (bool, error) {
	return isFound(GetDB(nil).Select("status").Where("hash = ?", hash).First(lt))
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

	err = GetDB(nil).Select("hash,block,status").Offset((page - 1) * limit).Limit(limit).Order("block desc").Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}

	if num < 1 {
		return &ret, num, err
	}
	TBlock := make(map[string]int64)
	Thash := make(map[string]bool)
	hashStatus := make(map[string]int32)
	for i = 0; i < len(tss); i++ {
		hash := hex.EncodeToString(tss[i].Hash)
		Thash[hash] = true
		hashStatus[hash] = tss[i].Status

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
				bh.Hash = rt.Transactions[j].Hash
				bh.KeyID = rt.Transactions[j].KeyID
				//bh.Params = rt.Transactions[j].Params
				if reqType == 1 {
					bh.Time = bk.Time
				} else {
					bh.Time = rt.Transactions[j].Time
				}
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
					bh.Status = hashStatus[rt.Transactions[j].Hash]
					if reqType == 0 {
						if bh.ContractName == UtxoTx {
							var params types.UTXO
							err := json.Unmarshal([]byte(rt.Transactions[j].Params), &params)
							if err != nil {
								return nil, 0, err
							}
							bh.Amount, _ = decimal.NewFromString(params.Value)
							//bh.GasFee todo: need add
						} else if bh.ContractName == UtxoTransfer {
							var params types.TransferSelf
							err := json.Unmarshal([]byte(rt.Transactions[j].Params), &params)
							if err != nil {
								return nil, 0, err
							}
							bh.Amount, _ = decimal.NewFromString(params.Value)
						} else {
							var his History
							gasFee, amount, err := his.GetTxListExplorer(converter.HexToBin(rt.Transactions[j].Hash))
							if err != nil {
								return nil, 0, err
							}
							bh.GasFee = gasFee
							bh.Amount = amount
						}
					}
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

func GetTransactionBlockToRedis() error {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
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
			txDetailedInfo.ContractName = UtxoTx
			dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.UTXO)
			txDetailedInfo.Params = string(dataBytes)
		} else if tx.SmartContract().TxSmart.TransferSelf != nil {
			txDetailedInfo.ContractName = UtxoTransfer
			dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.TransferSelf)
			txDetailedInfo.Params = string(dataBytes)
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
		txDetailedInfo.TokenSymbol, txDetailedInfo.Ecosystemname = Tokens.Get(txDetailedInfo.Ecosystem), EcoNames.Get(txDetailedInfo.Ecosystem)
	} else {
		if txDetailedInfo.Ecosystem == 0 {
			txDetailedInfo.Ecosystem = 1
			txDetailedInfo.TokenSymbol, txDetailedInfo.Ecosystemname = Tokens.Get(txDetailedInfo.Ecosystem), EcoNames.Get(txDetailedInfo.Ecosystem)
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
	var tx TransactionData
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
		rets.HashType = 1
	} else {
		f, err = tx.GetTxDataByHash(hashHex)
		if err != nil {
			return rets, err
		}
		if !f {
			return rets, errors.New("transaction data synchronization")
		}
		f, err = IsUtxoTransaction(tx.TxData)
		if err != nil {
			return rets, err
		}
		if f {
			rets.HashType = 2
		} else {
			rets.HashType = 3
		}
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

func (lt *LogTransaction) GetEcosystemAccountTransaction(ecosystem int64, page int, size int, wallet, order string, where map[string]any) (*GeneralResponse, error) {
	var (
		tss   []LogTransaction
		ret   []AccountTxListResponse
		count int64
		keyId int64
		err   error
		rets  GeneralResponse
		q     *gorm.DB
	)
	rets.Limit = size
	rets.Page = page
	if order == "" {
		order = "timestamp desc"
	}

	keyId = converter.StringToAddress(wallet)
	if wallet == "0000-0000-0000-0000-0000" {
	} else if keyId == 0 {
		return &rets, errors.New("wallet does not meet specifications")
	}
	if page < 1 || size < 1 {
		return &rets, err
	}
	if ecosystem != 0 {
		if where == nil {
			where = make(map[string]any)
		}
		where["ecosystem_id ="] = ecosystem
		dayTime := int64(60 * 60 * 24)
		if value, ok := where["timestamp >="]; ok {
			if reflect.TypeOf(value).String() == "json.Number" {
				val, err := value.(json.Number).Int64()
				if err != nil {
					return nil, err
				}
				where["timestamp >="] = val * 1000
			}
		}
		if value, ok := where["timestamp <="]; ok {
			if reflect.TypeOf(value).String() == "json.Number" {
				val, err := value.(json.Number).Int64()
				if err != nil {
					return nil, err
				}
				where["timestamp <="] = (val + dayTime) * 1000
			}
		}
	}
	if len(where) != 0 {
		cond, vals, err := WhereBuild(where)
		if err != nil {
			return &rets, err
		}
		fmt.Printf("cond:%s\n", cond)
		fmt.Println("vals:", vals)
		q = GetDB(nil).Table(lt.TableName()).Where(cond, vals...).Where("address = ?", keyId)
	} else {
		q = GetDB(nil).Table(lt.TableName()).Where("address = ?", keyId)
	}
	if err = q.Count(&count).Error; err != nil {
		return &rets, err
	}
	if count > 0 {
		err = q.Order(order).Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}

	if err != nil {
		return &rets, err
	}

	length := len(tss)
	for i := 0; i < length; i++ {
		da := AccountTxListResponse{}
		da.Hash = hex.EncodeToString(tss[i].Hash)
		da.BlockId = tss[i].Block
		da.Timestamp = MsToSeconds(tss[i].Timestamp)
		da.Address = converter.AddressToString(tss[i].Address)
		da.ContractName = tss[i].ContractName
		if da.ContractName == "" {
			da.ContractName = GetUtxoTxContractNameByHash(tss[i].Hash)
		}
		da.Status = tss[i].Status

		ret = append(ret, da)
	}
	if length < rets.Limit {
		rets.Limit = length
	}
	rets.Total = count
	rets.List = ret
	return &rets, nil
}

func (lt *LogTransaction) GetEcosystemAccountTransactionNew(ecosystem int64, page int, size int, wallet, order string, where map[string]any) (*GeneralResponse, error) {
	var (
		ret   []AccountTxListResponse
		count int64
		keyId int64
		err   error
		rets  GeneralResponse
		q1    *gorm.DB
		q2    *gorm.DB
	)
	rets.Limit = size
	rets.Page = page
	if order == "" {
		order = "timestamp desc"
	} else {
		if !CheckSql(order) {
			return nil, errors.New("request params invalid")
		}
	}

	keyId = converter.StringToAddress(wallet)
	if wallet == "0000-0000-0000-0000-0000" {
	} else if keyId == 0 {
		return &rets, errors.New("wallet does not meet specifications")
	}
	if page < 1 || size < 1 {
		return &rets, err
	}
	recipientLike := "%" + fmt.Sprintf("%d", keyId) + "%"
	if ecosystem != 0 {
		if where == nil {
			where = make(map[string]any)
		}
		where["ecosystem ="] = ecosystem
		dayTime := int64(60 * 60 * 24)
		if value, ok := where["timestamp >="]; ok {
			if reflect.TypeOf(value).String() == "json.Number" {
				val, err := value.(json.Number).Int64()
				if err != nil {
					return nil, err
				}
				where["timestamp >="] = val * 1000
			}
		}
		if value, ok := where["timestamp <="]; ok {
			if reflect.TypeOf(value).String() == "json.Number" {
				val, err := value.(json.Number).Int64()
				if err != nil {
					return nil, err
				}
				where["timestamp <="] = (val + dayTime) * 1000
			}
		}
		where["recipient_id like"] = recipientLike
	}
	type accountTxList struct {
		Hash         []byte
		Block        int64
		SenderId     int64
		Timestamp    int64
		ContractName string
		Ecosystem    int64
		Status       int32
	}
	var (
		list       []accountTxList
		sqlQuery   string
		countQuery string
	)

	if len(where) != 0 {
		cond, vals, err := WhereBuild(where)
		if err != nil {
			return &rets, err
		}
		countQuery = fmt.Sprintf(`
SELECT count(1) FROM(
	SELECT v1.*,v2.recipient_id,v2.ecosystem FROM(
		SELECT hash,block,address AS sender_id,timestamp,contract_name,status FROM log_transactions AS log 
	)AS v1
	LEFT JOIN(
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash,ecosystem FROM(
			SELECT recipient_id AS keyid,txhash AS hash,ecosystem FROM "1_history" GROUP BY txhash,ecosystem,recipient_id
		)AS v1 GROUP BY hash,ecosystem
			UNION
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash,ecosystem FROM(
			SELECT output_key_id AS keyid,output_tx_hash AS hash,ecosystem FROM spent_info GROUP BY output_tx_hash,ecosystem,output_key_id
		)AS v1 GROUP BY hash,ecosystem

	)AS v2 ON(v2.hash = v1.hash)
)AS v1 
WHERE %s OR (sender_id = %d AND ecosystem = %d) 
`, cond, keyId, ecosystem)
		sqlQuery = fmt.Sprintf(`
SELECT * FROM(
	SELECT v1.*,v2.recipient_id,v2.ecosystem FROM(
		SELECT hash,block,address AS sender_id,timestamp,contract_name,status FROM log_transactions AS log 
	)AS v1
	LEFT JOIN(
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash,ecosystem FROM(
			SELECT recipient_id AS keyid,txhash AS hash,ecosystem FROM "1_history" GROUP BY txhash,ecosystem,recipient_id
		)AS v1 GROUP BY hash,ecosystem
			UNION
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash,ecosystem FROM(
			SELECT output_key_id AS keyid,output_tx_hash AS hash,ecosystem FROM spent_info GROUP BY output_tx_hash,ecosystem,output_key_id
		)AS v1 GROUP BY hash,ecosystem

	)AS v2 ON(v2.hash = v1.hash)
	ORDER BY %s
)AS v1 
WHERE %s OR (sender_id = %d AND ecosystem = %d) 
`, order, cond, keyId, ecosystem)
		q1 = GetDB(nil).Raw(sqlQuery, vals...)
		q2 = GetDB(nil).Raw(countQuery, vals...)
	} else {
		countQuery = `
SELECT count(1) FROM(
	SELECT v1.*,v2.recipient_id FROM(
		SELECT hash,block,address AS sender_id,timestamp,contract_name,status,ecosystem_id as ecosystem FROM log_transactions AS log 
	)AS v1
	LEFT JOIN(
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash FROM(
			SELECT recipient_id AS keyid,txhash AS hash FROM "1_history" GROUP BY txhash,recipient_id
		)AS v1 GROUP BY hash
			UNION
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash FROM(
			SELECT output_key_id AS keyid,output_tx_hash AS hash FROM spent_info GROUP BY output_tx_hash,output_key_id
		)AS v1 GROUP BY hash

	)AS v2 ON(v2.hash = v1.hash)
)AS v1 
WHERE (recipient_id like ?) OR(sender_id = ?)
`

		sqlQuery = fmt.Sprintf(`
SELECT * FROM(
	SELECT v1.*,v2.recipient_id FROM(
		SELECT hash,block,address AS sender_id,timestamp,contract_name,status,ecosystem_id as ecosystem FROM log_transactions AS log 
	)AS v1
	LEFT JOIN(
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash FROM(
			SELECT recipient_id AS keyid,txhash AS hash FROM "1_history" GROUP BY txhash,recipient_id
		)AS v1 GROUP BY hash
			UNION
		SELECT array_to_string(array_agg(keyid),',') AS recipient_id,hash FROM(
			SELECT output_key_id AS keyid,output_tx_hash AS hash FROM spent_info GROUP BY output_tx_hash,output_key_id
		)AS v1 GROUP BY hash

	)AS v2 ON(v2.hash = v1.hash)
	ORDER BY %s
)AS v1 
WHERE (recipient_id like ?) OR(sender_id = ?)
OFFSET %d LIMIT %d
`, order, (page-1)*size, size)
		q1 = GetDB(nil).Raw(sqlQuery, recipientLike, keyId)
		q2 = GetDB(nil).Raw(countQuery, recipientLike, keyId)
	}
	if err = q2.Take(&count).Error; err != nil {
		return &rets, err
	}
	if count > 0 {
		err = q1.Find(&list).Error
	}

	if err != nil {
		return &rets, err
	}

	length := len(list)
	for i := 0; i < length; i++ {
		da := AccountTxListResponse{}
		da.Hash = hex.EncodeToString(list[i].Hash)
		da.BlockId = list[i].Block
		da.Timestamp = MsToSeconds(list[i].Timestamp)
		da.Address = converter.AddressToString(list[i].SenderId)
		da.ContractName = list[i].ContractName
		if da.ContractName == "" {
			da.ContractName = GetUtxoTxContractNameByHash(list[i].Hash)
		}
		da.Status = list[i].Status
		da.EcosystemName = EcoNames.Get(list[i].Ecosystem)
		da.Ecosystem = list[i].Ecosystem

		ret = append(ret, da)
	}
	if length < rets.Limit {
		rets.Limit = length
	}
	rets.Total = count
	rets.List = ret
	return &rets, nil
}
