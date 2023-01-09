/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/converter"

	//"strconv"
	"strings"
	"time"

	"github.com/IBAX-io/go-ibax/packages/block"
)

const (
	ChainMax = "chain_max"
	MintMax  = "mint_max"
)

var FirstBlockTime int64

// Block is model
type Block struct {
	ID             int64  `gorm:"primary_key;not_null"`
	Hash           []byte `gorm:"not null"`
	RollbacksHash  []byte `gorm:"not null"`
	Data           []byte `gorm:"not null"`
	EcosystemID    int64  `gorm:"not null"`
	KeyID          int64  `gorm:"not null"`
	NodePosition   int64  `gorm:"not null"`
	Time           int64  `gorm:"not null"`
	Tx             int32  `gorm:"not null"`
	ConsensusMode  int32  `gorm:"not null"`
	CandidateNodes []byte `gorm:"not null"`
}

// TableName returns name of table
func (Block) TableName() string {
	return "block_chain"
}

func (b *Block) GetId(blockID int64) (bool, error) {
	return isFound(GetDB(nil).Where("id = ?", blockID).First(b))
}

func (b *Block) GetByTimeBlockId(dbTx *DbTransaction, time int64) (bool, error) {
	return isFound(GetDB(dbTx).Select("id").Where("time >= ?", time).First(b))
}

// GetNodeBlocksAtTime returns records of blocks for time interval and position of node
func (b *Block) GetBlocksHash(hash []byte) (bool, error) {
	//f, err := b.GetRedisByhash(hash)
	//if f && err == nil {
	//	fmt.Println("return redis !!!\n")
	//	return f, err
	//}

	f, err := isFound(conf.GetDbConn().Conn().Where("hash = ?", hash).First(&b))
	if f && err == nil {
		//b.InsertRedis()
		return f, err
	}
	return f, err
}

func (b *Block) GetBlocksKey(key int64, order string) ([]Block, error) {
	var err error
	var blockchain []Block

	err = conf.GetDbConn().Conn().Order(order).Where("key_id = ?", key).Find(&blockchain).Error
	return blockchain, err
}

// GetMaxBlock returns last block existence
func (b *Block) GetMaxBlock() (bool, error) {
	return isFound(conf.GetDbConn().Conn().Last(b))
}

// GetBlockData is retrieving chain of blocks from database
func GetBlockData(startId int64, endId int64, order string) (*[]Block, error) {
	var err error
	blockchain := new([]Block)

	orderStr := "id " + string(order)
	query := conf.GetDbConn().Conn().Model(&Block{}).Order(orderStr)
	if endId > 0 {
		query = query.Select("id,time,data").Where("id > ? AND id <= ?", startId, endId).Find(&blockchain)
	} else {
		query = query.Select("id,time,data").Where("id > ?", startId).Find(&blockchain)
	}

	if query.Error != nil {
		return nil, err
	}
	return blockchain, nil
}

func (b *Block) GetBlockTotal() (int64, error) {
	var err error
	var total int64
	if err = GetDB(nil).Table(b.TableName()).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, err
}

// GetBlocksFrom is retrieving ordered chain of blocks from database
func (b *Block) GetBlocksFrom(page, limit int, ordering string) ([]Block, error) {
	var err error
	blockchain := new([]Block)

	if limit == 0 {
		err = GetDB(nil).Order("id " + ordering).Find(&blockchain).Error
	} else {
		err = GetDB(nil).Order("id " + ordering).Offset((page - 1) * limit).Limit(limit).Find(&blockchain).Error
	}
	return *blockchain, err
}

// GetReverseBlockchain returns records of blocks in reverse ordering
func (b *Block) GetReverseBlockchain(endBlockID int64, limit int) ([]Block, error) {
	var err error
	blockchain := new([]Block)

	err = conf.GetDbConn().Conn().Model(&Block{}).Order("id DESC").Where("id <= ?", endBlockID).Limit(limit).Find(&blockchain).Error
	return *blockchain, err
}

// GetNodeBlocksAtTime returns records of blocks for time interval and position of node
func (b *Block) GetNodeBlocksAtTime(from, to time.Time, node int64) ([]Block, error) {
	var err error
	blockchain := new([]Block)

	err = conf.GetDbConn().Conn().Model(&Block{}).Where("node_position = ? AND time BETWEEN ? AND ?", node, from.Unix(), to.Unix()).Find(&blockchain).Error
	return *blockchain, err
}

func (b *Block) GetBlocksByHash(hash []byte) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("hash = ?", hash).First(b))
}

func GetBlockList(page, limit, reqType int, order string) (*BlockListHeaderResponse, error) {
	var (
		ret  BlockListHeaderResponse
		err  error
		rets []BlockListResponse
		bk   Block
	)
	if order == "" {
		order = "id desc"
	}
	if page < 1 || limit < 1 {
		return &ret, err
	}
	ret.Page = page
	ret.Limit = limit
	if reqType == 0 {
		var m ScanOut
		f, err := m.GetRedisLatest()
		if err != nil {
			return &ret, err
		}
		if f {
			var bk Block
			f, err := bk.GetMaxBlock()
			if err != nil {
				return &ret, err
			}
			if !f {
				return &ret, errors.New("get block list max block failed")
			}
			maxBlock, err := getMaxBlockSizeFromRedis()
			if err != nil {
				return nil, err
			}
			maxTx, err := getMaxTxFromRedis()
			if err != nil {
				return nil, err
			}

			bkRet := &BlockRet{}
			bkRet.BlockId = bk.ID
			bkRet.MaxTps = maxTx.Tx
			bkRet.MaxBlockSize = m.MaxBlockSize
			bkRet.StorageCapacitys = m.StorageCapacity
			bkRet.MaxBlockSizeId = maxBlock.ID
			bkRet.MaxTpsId = maxTx.ID
			ret.BlockInfo = bkRet
		}
	}
	ret.Total, err = bk.GetBlockTotal()
	if err != nil {
		log.WithFields(log.Fields{"warn": err.Error()}).Warn("Get Block Total Failed")
		return &ret, err
	}

	type bkListResponse struct {
		Id            int64
		Time          int64
		Tx            int32
		NodePosition  int64
		GasFee        string
		EcoLib        int64
		RecipientId   string
		Reward        string
		Address       string
		ApiAddress    string
		ConsensusMode int32
		NodeName      string
		Hid           int64
	}
	if !CheckSql(order) {
		log.WithFields(log.Fields{"warn": err.Error(), "order": order}).Warn("Get Block list check sql Failed")
		return &ret, err
	}
	var list []bkListResponse
	err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT v1.id,v1.time,v1.tx,v1.node_position,v1.consensus_mode,coalesce((SELECT sum(amount) FROM "1_history" WHERE block_id = v1.id AND type IN(1,2) AND ecosystem = 1),0)gas_fee,
	coalesce((SELECT count(1) FROM (SELECT ecosystem_id FROM log_transactions WHERE block = v1.id GROUP BY ecosystem_id)AS lg),0)eco_lib,
	coalesce(CAST((SELECT min(recipient_id) FROM "1_history" WHERE block_id = v1.id AND type IN(12)) AS TEXT),'')recipient_id,
	coalesce((SELECT sum(amount) FROM "1_history" WHERE block_id = v1.id AND type IN(12)),0)reward,coalesce(t2.address,'')address,coalesce(t2.api_address,'')api_address,
	COALESCE(t2.node_id,0)hid,coalesce(node_name,'')node_name
FROM(
		SELECT id,time,tx,node_position,consensus_mode FROM block_chain ORDER BY %s OFFSET ? LIMIT ?
)AS v1
LEFT JOIN(
        SELECT CAST(value->>'id' AS numeric) node_id,CAST(value->>'consensus_mode' AS numeric) consensus_mode,address,id,
		value->>'api_address' api_address,value->>'node_name' node_name
		FROM honor_node_info
)AS t2 ON(t2.node_id = v1.node_position AND t2.consensus_mode = v1.consensus_mode)
ORDER BY %s
`, order, order), (page-1)*limit, limit).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"warn": err.Error()}).Warn("Get Block List Failed")
		return &ret, err
	}
	for _, value := range list {
		var rt BlockListResponse
		rt.ID = value.Id
		if value.RecipientId != "" {
			keyStr, _ := strconv.ParseInt(value.RecipientId, 10, 64)
			rt.Recipientid = converter.AddressToString(keyStr)
		}
		rt.Reward = value.Reward
		rt.Time = value.Time
		rt.Tx = value.Tx
		if value.NodeName != "" {
			rt.NodeName = value.NodeName
		} else {
			rt.NodeName = "HONOR_NODE" + strconv.FormatInt(value.Hid, 10)
		}
		rt.GasFee = value.GasFee
		rt.IconUrl = getIconNationalFlag(getCountry(value.Address))
		rt.APIAddress = value.ApiAddress
		rt.NodePosition = value.NodePosition
		rt.EcoLib = value.EcoLib
		rt.ConsensusMode = value.ConsensusMode

		rets = append(rets, rt)
	}
	ret.Digits = EcoDigits.GetInt64(1, 12)

	ret.List = rets

	return &ret, err
}

func GetBlockListFromRedis() (*BlockListHeaderResponse, error) {
	return GetBlocksListFromRedis()
}

func GetBlockTpsListsFromRedis() (*[]ScanOutBlockTransactionRet, error) {

	ret1, err := GetTxInfoFromRedis(30)
	if err == nil {
		return ret1, err
	} else {
		return nil, err
	}
}

func GetBlocksDetailedInfoHex(bk *Block) (*BlockDetailedInfoHex, error) {
	var (
		transize int64
	)
	result := BlockDetailedInfoHex{}

	blck, err := block.UnmarshallBlock(bytes.NewBuffer(bk.Data), false)
	if err != nil {
		return &result, err
	}

	txDetailedInfoCollection := make([]TxDetailedInfoHex, 0, len(blck.Transactions))
	for _, tx := range blck.Transactions {
		txDetailedInfo := TxDetailedInfoHex{
			Hash: hex.EncodeToString(tx.Hash()),
		}
		if tx.IsSmartContract() {
			if tx.SmartContract().TxSmart.UTXO != nil {
				txDetailedInfo.ContractName = UtxoTx
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.UTXO)
				txDetailedInfo.Params = string(dataBytes)
				if converter.AddressToString(tx.SmartContract().TxSmart.UTXO.ToID) == BlackHoleAddr {
					txDetailedInfo.ContractName = UtxoBurning
				}
			} else if tx.SmartContract().TxSmart.TransferSelf != nil {
				txDetailedInfo.ContractName = UtxoTransferSelf
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.TransferSelf)
				txDetailedInfo.Params = string(dataBytes)
			} else {
				txDetailedInfo.ContractName, txDetailedInfo.Params = GetMineParam(tx.SmartContract().TxSmart.EcosystemID, tx.SmartContract().TxContract.Name, tx.SmartContract().TxData, tx.Hash())
			}
			//txDetailedInfo.ContractName = tx.TxContract.Name
			//txDetailedInfo.Params = tx.TxData

		}
		//first block txcontract is nil,but time is not null
		txDetailedInfo.KeyID = converter.AddressToString(tx.KeyID())
		txDetailedInfo.Time = MsToSeconds(tx.Timestamp())
		txDetailedInfo.Type = int64(tx.Type())
		txDetailedInfo.Size = int64(len(tx.FullData))
		transize += txDetailedInfo.Size

		if tx.IsSmartContract() {
			txDetailedInfo.Ecosystem = tx.SmartContract().TxSmart.EcosystemID
			if txDetailedInfo.Ecosystem == 0 {
				txDetailedInfo.Ecosystem = 1
			}
			txDetailedInfo.TokenSymbol, txDetailedInfo.Ecosystemname = Tokens.Get(txDetailedInfo.Ecosystem), EcoNames.Get(txDetailedInfo.Ecosystem)
		}

		txDetailedInfoCollection = append(txDetailedInfoCollection, txDetailedInfo)

		//log.WithFields(log.Fields{"block_id": blockModel.ID, "tx hash": txDetailedInfo.Hash, "contract_name": txDetailedInfo.ContractName, "key_id": txDetailedInfo.KeyID, "time": txDetailedInfo.Time, "type": txDetailedInfo.Type, "params": txDetailedInfoCollection}).Debug("Block Transactions Information")
	}
	if blck.Header.EcosystemId == 0 {
		blck.Header.EcosystemId = 1
	}
	if bk.EcosystemID == 0 {
		bk.EcosystemID = 1
	}
	header := BlockHeaderInfoHex{
		BlockId:      blck.Header.BlockId,
		Time:         blck.Header.Timestamp,
		EcosystemId:  blck.Header.EcosystemId,
		KeyId:        converter.AddressToString(blck.Header.KeyId),
		NodePosition: blck.Header.NodePosition,
		Sign:         hex.EncodeToString(blck.Header.Sign),
		BlockHash:    hex.EncodeToString(blck.Header.BlockHash),
		Version:      blck.Header.Version,
	}

	bdi := BlockDetailedInfoHex{
		Header:        header,
		Hash:          hex.EncodeToString(bk.Hash),
		EcosystemID:   bk.EcosystemID,
		NodePosition:  bk.NodePosition,
		KeyID:         converter.AddressToString(bk.KeyID),
		Time:          bk.Time,
		Tx:            bk.Tx,
		RollbacksHash: hex.EncodeToString(bk.RollbacksHash),
		MerkleRoot:    hex.EncodeToString(blck.MerkleRoot),
		BinData:       hex.EncodeToString(blck.BinData),
		SysUpdate:     blck.SysUpdate,
		GenBlock:      blck.GenBlock,
		//StopCount:     blck.s,
		BlockSize:    ToCapacityString(int64(len(bk.Data))),
		TxTotalSize:  ToCapacityString(transize),
		Transactions: txDetailedInfoCollection,
	}

	return &bdi, nil
}
func GetMineParam(ecosystem int64, name string, param map[string]any, TxHash []byte) (string, string) {
	escape := func(value any) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}
	if name == "@1CallDelayedContract" && ecosystem == 1 {
		v, ok := param["Id"]
		if ok {
			idstr := escape(v)
			if idstr == "4" {
				return "@1Mint", GetMineIncomeParam(TxHash)
			}
		}

	}
	dataBytes, _ := json.Marshal(param)
	return name, string(dataBytes)
}

func GetMineIncomeParam(hash []byte) string {
	ret := make(map[string]any)
	ts := &MineIncomehistory{}
	f, err := ts.Get(hash)
	if err == nil && f {
		ret["miner"] = ts.Devid
		ret["minerowner"] = ts.Keyid
		ret["profiter"] = ts.Mineid
		ret["type"] = ts.Type
		ret["staked"] = ts.Nonce
		ret["earnings"] = ts.Amount
		//return ret
	}
	dataBytes, _ := json.Marshal(ret)
	return string(dataBytes)
}

func SyncBlockListToRedis() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	rets, err := GetBlockList(1, 10, 1, "")
	if err != nil {
		return
	}

	date, err := json.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("Sync Block List To Redis json marshal err")
		return
	}
	//value, err := msgpack.Marshal(blockInfo)
	//if err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis msgpack err")
	//	return
	//}
	//value, err := GzipEncode(blockInfo)
	//if err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis GzipEncode err")
	//	return
	//}
	rd := RedisParams{
		Key:   "dashboard-block",
		Value: string(date),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("Sync Block info To Redis set db err")
	}
}

func GetBlocksListFromRedis() (*BlockListHeaderResponse, error) {
	var err error
	var bk BlockListHeaderResponse
	rd := RedisParams{
		Key:   "dashboard-block",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis getdb err")
		return nil, err
	}
	//value, err1 := GzipDecode([]byte(rd.Value))
	//if err1 != nil {
	//	log.WithFields(log.Fields{"warn": err1}).Warn("GetBlocksListFromRedis GzipDecode err")
	//	return nil, err1
	//}

	//var blockInfo []byte
	//err = msgpack.Unmarshal([]byte(rd.Value), &blockInfo)
	//if err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis msgpack err")
	//	return nil, err
	//}

	if err = json.Unmarshal([]byte(rd.Value), &bk); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis json err")
		return nil, err
	}
	return &bk, nil
}

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		err = writer.Close()
		return out, err
	}
	if err = writer.Close(); err != nil {
		return out, err
	}
	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer func() {
		if err = reader.Close(); err != nil {
			println("Gzip unzip error", err.Error())
		}
	}()
	return ioutil.ReadAll(reader)
}

func SendBlockListToWebsocket(ret1 *[]BlockListResponse) error {
	var (
		err error
	)

	if err = SendBlockList(ret1); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("Send Block Transaction err")
		return err
	}
	return nil
}
func SendBlockList(data *[]BlockListResponse) error {
	err := SendDashboardDataToWebsocket(data, ChannelBlockList)
	if err != nil {
		return err
	}
	return nil
}

func InitFristBlockTime() error {
	var bk Block
	t1, err := bk.GetSystemTime()
	if err != nil {
		return err
	}
	if t1 == 0 {
		return errors.New("init first block time failed")
	}
	FirstBlockTime = t1
	return nil
}

// GetSystemTime is retrieving model from database
func (b *Block) GetSystemTime() (int64, error) {
	f, err := isFound(GetDB(nil).Select("time").Where("id = 1").First(b))
	if err == nil && f {
		return b.Time, nil
	}
	return 0, err
}

func getNodePkgInfo(nodePosition int64, consensusMode int32, ret []nodePkg) (decimal.Decimal, int64, string) {
	zero := decimal.New(0, 0)
	for i := 0; i < len(ret); i++ {
		if ret[i].NodePosition == nodePosition && ret[i].ConsensusMode == consensusMode {
			return ret[i].PkgFor, ret[i].Count, converter.AddressToString(ret[i].KeyId)
		}
	}
	return zero, 0, ""
}

func GetBlocksContractNameList(bk *Block) (map[string]string, error) {
	blck, err := block.UnmarshallBlock(bytes.NewBuffer(bk.Data), false)
	if err != nil {
		return nil, err
	}

	contractNameList := make(map[string]string)
	for _, tx := range blck.Transactions {
		txDetailedInfo := TxDetailedInfoHex{
			Hash: hex.EncodeToString(tx.Hash()),
		}
		if tx.IsSmartContract() {
			if tx.SmartContract().TxSmart.UTXO != nil {
				txDetailedInfo.ContractName = UtxoTx
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.UTXO)
				txDetailedInfo.Params = string(dataBytes)
				if converter.AddressToString(tx.SmartContract().TxSmart.UTXO.ToID) == BlackHoleAddr {
					txDetailedInfo.ContractName = UtxoBurning
				}
			} else if tx.SmartContract().TxSmart.TransferSelf != nil {
				txDetailedInfo.ContractName = UtxoTransferSelf
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.TransferSelf)
				txDetailedInfo.Params = string(dataBytes)
			} else {
				txDetailedInfo.ContractName, _ = GetMineParam(tx.SmartContract().TxSmart.EcosystemID, tx.SmartContract().TxContract.Name, tx.SmartContract().TxData, tx.Hash())
			}
		}

		contractNameList[txDetailedInfo.Hash] = txDetailedInfo.ContractName
	}
	return contractNameList, nil
}

func (bk *Block) GetLatestNodes(limit int) ([]Block, error) {
	var list []Block
	if err := GetDB(nil).Select("node_position,id,consensus_mode").Limit(limit).Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetNodeBlockReplyRate(bk *Block) (string, error) {
	if bk.ConsensusMode == 1 {
		return "100", nil
	}
	blck, err := block.UnmarshallBlock(bytes.NewBuffer(bk.Data), false)
	if err != nil {
		return "0", err
	}
	var list []CandidateNodeRequests
	if blck.Header.CandidateNodes != nil {
		err := json.Unmarshal(blck.Header.CandidateNodes, &list)
		if err != nil {
			return "0", err
		}
		nodeNum := decimal.NewFromInt(int64(len(list)))
		for _, val := range list {
			if val.Id == bk.NodePosition {
				replyNum := decimal.NewFromInt(val.ReplyCount)
				if replyNum.GreaterThanOrEqual(nodeNum) {
					return "100", nil
				} else {
					return replyNum.Mul(decimal.NewFromInt(100)).DivRound(nodeNum, 2).String(), nil
				}
			}
		}
	}

	return "0", nil
}

func CheckSql(query string) bool {
	if strings.Contains(query, "drop") || strings.Contains(query, ";") {
		return false
	}
	return true
}
