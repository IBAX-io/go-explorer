/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/smart"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type ScanOut struct {
	BlockSizes        int64
	BlockTransactions int64
	Hash              string
	RollbacksHash     string
	KeyID             string
	NodePosition      int64
	ConsensusMode     int32
	Time              int64
	CurrentVersion    string

	TotalCounts          int64 //total count power
	BlockTransactionSize int64
	HonorNode            int64 //honorNode number
	CastNodes            int64
	SubNodes             int64
	CLBNodes             int64

	MintAmounts string

	Blockid          int64
	MaxTps           int64
	MaxBlockSize     string
	StorageCapacitys string

	TotalTx         int64
	TwentyFourTx    int64
	WeekAverageTx   int64
	MaxActiveEcoLib string

	Circulations            string
	TodayCirculationsAmount float64
	TwentyFourAmount        string

	NftMinerCount   int64
	NftBlockReward  float64
	HalveNumber     int64
	NftStakeAmounts string

	EcosystemCount int64
}

type ScanOutRet struct {
	Blockid           int64  `json:"block_id"`    //Block Id
	BlockSizes        int64  `json:"block_sizes"` //Block Size
	BlockTranscations int64  `json:"block_transcations"`
	Hash              string `json:"hash"`
	RollbacksHash     string `json:"rollbacks_hash"`
	EcosystemID       int64  `json:"ecosystem_id"`
	KeyID             string `json:"key_id"`
	NodePosition      int64  `json:"node_position"`
	ConsensusMode     int32  `json:"consensus_mode"`
	Time              int64  `json:"time"`
	CurrentVersion    string `json:"current_version"`

	TotalCounts          int64  `json:"total_counts"` //total count
	BlockTranscationSize int64  `json:"block_transcation_size"`
	GuardianNodes        int64  `json:"guardian_nodes"`
	SubNodes             int64  `json:"sub_nodes"`
	CLBNodes             int64  `json:"clb_nodes"`
	MintAmounts          string `json:"mint_amount"`

	BlockInfo         BlockRet              `json:"block_info"`
	TxInfo            TxRet                 `json:"tx_info"`
	CirculationsInfo  CirculationsRet       `json:"circulations_info"`
	NftMinerInfo      NftMinerRet           `json:"nft_miner_info"`
	EcoLibsInfo       EcoLibsRet            `json:"eco_libs_info"`
	KeysInfo          KeysRet               `json:"keys_info"`
	CandidateNodeInfo CandidateHonorNodeRet `json:"candidate_node_info"`
	CastNodeInfo      CastNodeRet           `json:"cast_node_info"` //todo:Need Add
}

type BlockRet struct {
	Blockid          int64  `json:"block_id"` //
	MaxTps           int64  `json:"max_tps"`
	MaxBlockSize     string `json:"max_block_size"`
	StorageCapacitys string `json:"storage_capacitys"`
	MaxTpsId         int64  `json:"max_tps_id,omitempty"`
	MaxBlockSizeId   int64  `json:"max_block_size_id,omitempty"`
}

type TxRet struct {
	TotalTx        int64 `json:"total_tx"`
	TwentyFourTx   int64 `json:"twenty_four_tx"`
	SingleDayMaxTx int64 `json:"single_day_max_tx"`
	WeekAverageTx  int64 `json:"week_average_tx"`
}

type CirculationsRet struct {
	Circulations            string  `json:"circulations_amount"`
	TotalAmounts            string  `json:"total_amount"`
	TodayCirculationsAmount float64 `json:"today_circulations_amount"`
	TwentyFourAmount        string  `json:"twenty_four_amount"`
}

//Nft Miner
type NftMinerRet struct {
	Count        int64   `json:"count"`
	BlockReward  float64 `json:"block_reward"`
	HalveNumber  int64   `json:"halve_number"`
	StakeAmounts string  `json:"stake_amounts"`
}

type EcoLibsRet struct {
	Ecosystems    int64 `json:"ecosystems"`
	EcoTokenTotal int64 `json:"eco_token_total"`
	DaoGovern     int64 `json:"dao_govern"`
	Contract      int64 `json:"contract"`
}

type KeysRet struct {
	KeyCount       int64 `json:"key_count"`
	HasTokenKey    int64 `json:"has_token_key"`
	MonthActiveKey int64 `json:"month_active_key"`
	TwentyFourKey  int64 `json:"twenty_four_key"`
}

type CandidateHonorNodeRet struct {
	CandidateNode    int64  `json:"candidate_node"`
	NodeStakeAmounts string `json:"node_stake_amounts"`
	NodeVote         int64  `json:"node_vote"`
	TwentyFourNode   int64  `json:"twenty_four_node"`
}

type CastNodeRet struct {
	CastNodes          int64  `json:"cast_nodes"`
	DistributionRegion int64  `json:"distribution_region"`
	CastCapacitys      string `json:"cast_capacitys"`
	BandwidthTraffic   string `json:"bandwidth_traffic"`
}

type ScanOutBlockTransactionRet struct {
	BlockId           int64 `json:"block_id"`            //
	BlockSizes        int64 `json:"block_size" `         //
	BlockTranscations int64 `json:"block_transcations" ` //
}

var ScanPrefix = "scan-"
var ScanOutStPrefix = "scan-out-"
var ScanOutLastest = "lastest"
var GetScanOut chan bool
var SendScanOut chan bool

func (s *ScanOutBlockTransactionRet) Marshal(q []ScanOutBlockTransactionRet) (string, error) {
	if res, err := msgpack.Marshal(q); err != nil {
		return "", err
	} else {
		return string(res), err
	}
}

func (s *ScanOutBlockTransactionRet) Unmarshal(bt string) (q []ScanOutBlockTransactionRet, err error) {
	if err := msgpack.Unmarshal([]byte(bt), &q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *ScanOut) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(s); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (s *ScanOut) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &s); err != nil {
		return err
	}
	return nil
}

func (s *ScanOut) Get(id int64) (bool, error) {
	rp := &RedisParams{
		Key: ScanPrefix + strconv.FormatInt(id, 10),
	}
	for i := 0; i < 10; i++ {
		err := rp.Get()
		if err == nil {
			err = s.Unmarshal([]byte(rp.Value))
			return true, err
		}
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			break
		} else {
			time.Sleep(200 * time.Millisecond)
		}

	}

	return false, nil
}

func (m *ScanOut) Del(id int64) error {
	rp := &RedisParams{
		Key: ScanPrefix + strconv.FormatInt(id, 10),
	}
	var err error

	for i := 0; i < 5; i++ {
		err = rp.Del()
		if err == nil {
			break
		}
	}

	return err
}

func (m *ScanOut) DelRange(id, count int64) error {

	for i := int64(0); i < count; i++ {
		rp := &RedisParams{
			Key: ScanPrefix + strconv.FormatInt(id+i, 10),
		}

		err := rp.Del()
		if err != nil {
			return err
		}
	}

	return nil
}
func (m *ScanOut) Del_Redis(id int64) error {
	rd := RedisParams{
		Key:   ScanOutStPrefix + strconv.FormatInt(id, 10),
		Value: "",
	}
	if err := rd.Del(); err != nil {
		//log.WithFields(log.Fields{"err": err}).Warn("Del_Redis failed")
		return err
	}
	return nil
}

func (m *ScanOut) Insert_Redis() error {
	val, err := m.Marshal()
	if err != nil {
		return err
	}

	rd := RedisParams{
		Key:   ScanOutStPrefix + ScanOutLastest,
		Value: string(val),
	}
	err = rd.Set()
	if err != nil {
		return err
	}

	//rd = RedisParams{
	//	Key:   ScanOutStPrefix + strconv.FormatInt(m.Blockid, 10),
	//	Value: string(val),
	//}
	//err = rd.Set()

	return err
}

func (m *ScanOut) Get_Redis(id int64) (bool, error) {
	rd := RedisParams{
		Key:   ScanOutStPrefix + strconv.FormatInt(id, 10),
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return false, nil
		}
		return false, err
	}
	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return false, err
	}
	return true, err
}

func GetScanOutDataToRedis() error {
	err := processScanOutBlocks()
	return err
}

func initGlobalSwitch() {
	NodeReady = CandidateTableExist()
	NftMinerReady = NftMinerTableIsExist()
	VotingReady = VotingTableExist()
}

func (ret *ScanOut) Changes() error {
	//var mh MineIncomehistory
	//f, err := mh.GetID(ret.Blockid)
	//if err != nil {
	//	return err
	//}
	//if f {
	//	ret.TotalCounts = mh.Nonce
	//}
	initGlobalSwitch()

	var ne NftMinerItems
	powerCount, err := ne.GetAllPower()
	if err != nil {
		return errors.New("GetAllPower failed:" + err.Error())
	}
	ret.TotalCounts = powerCount

	tm, err := GetTotalAmount(1)
	if err != nil {
		return errors.New("GetTotalAmount failed:" + err.Error())
	}
	var nft NftMinerStaking
	_, nftStaking, err := nft.GetAllStakeAmount()
	if err != nil {
		return errors.New("GetAllStakeAmount failed:" + err.Error())
	}

	//stakeamount, err := key.GetStakeAmount()
	//if err != nil {
	//	return err
	//}

	ret.Circulations = tm.String()
	ret.NftStakeAmounts = nftStaking.String()

	systemCount, err := GetAllSystemCount()
	if err != nil {
		return errors.New("GetAllSystemCount failed:" + err.Error())
	}
	ret.EcosystemCount = systemCount

	var mst MinePledgeStatus
	honor, casts, nftCount, err := mst.GetCastNodeandGuardianNode()
	if err != nil {
		return errors.New("GetCastNodeandGuardianNode failed:" + err.Error())
	}
	ret.HonorNode = honor
	ret.CastNodes = casts
	ret.NftMinerCount = nftCount

	capacity, err := getDatabaseSize()
	if err != nil {
		return errors.New("getDatabaseSize failed:" + err.Error())
	}
	ret.StorageCapacitys = capacity

	var sp StateParameter
	sp.ecosystem = 1
	mb, err := sp.GetMintAmount()
	if err != nil {
		return errors.New("GetMintAmount failed:" + err.Error())
	}
	ret.MintAmounts = mb
	ret.HalveNumber, ret.NftBlockReward = getHalveNumber()

	var his History
	ret.TodayCirculationsAmount = his.GetTodayCirculationsAmount(ret.NftBlockReward)
	ret.TwentyFourAmount = his.Get24HourTxAmount()

	ret.MaxTps, _ = getMaxTps()
	maxSize, _ := getMaxBlockSize()
	ret.MaxBlockSize = TocapacityString(maxSize)
	ret.TotalTx = getTotalTx()
	ret.TwentyFourTx = getTwentyFourTx()
	ret.WeekAverageTx = getWeekAverageValueTx()
	ret.MaxActiveEcoLib = GetActiveEcoLibs()

	return nil
}

func processScanOutBlocks() error {
	//var bk Block
	var ret ScanOut
	cbk := InfoBlock{}
	if _, err := cbk.Get(); err != nil {
		return err
	}
	if cbk.BlockID == 2 {
		err := processScanOutFirstBlocks()
		if err != nil {
			return err
		}
	}
	ret.Blockid = cbk.BlockID
	ret.KeyID = converter.AddressToString(cbk.KeyID)
	ret.Time = cbk.Time
	ret.CurrentVersion = cbk.CurrentVersion
	ret.Hash = hex.EncodeToString(cbk.Hash)
	ret.NodePosition, _ = strconv.ParseInt(cbk.NodePosition, 10, 64)
	ret.ConsensusMode = cbk.ConsensusMode

	type blockTxInfo struct {
		BlockSize     int64
		TxSize        int64
		RollbacksHash string
		Tx            int64
	}
	var bk blockTxInfo
	err := GetDB(nil).Table("block_chain").Raw(`
SELECT length(data)block_size,tx,encode(rollbacks_hash, 'hex')rollbacks_hash,
        (SELECT sum(length(tx_data))tx_size FROM transaction_data WHERE block = bk.id)
FROM block_chain AS bk WHERE id = ?
`, cbk.BlockID).Take(&bk).Error
	if err != nil {
		return fmt.Errorf("process Scan Out Blocks Info Filed:%s\n", err.Error())
	}
	ret.BlockTransactions = bk.Tx
	ret.BlockTransactionSize = bk.TxSize
	ret.BlockSizes = bk.BlockSize
	ret.RollbacksHash = bk.RollbacksHash

	err = ret.InsertRedisDb()
	if err != nil {
		return fmt.Errorf("Insert Redisdb scanout:%s\n", err.Error())
	}

	return nil
}

func processScanOutFirstBlocks() error {
	var so ScanOut
	var bk Block

	fb, err := bk.GetId(1)
	if err != nil {
		return err
	}
	if fb {
		so.Blockid = 1
		so.KeyID = converter.AddressToString(bk.KeyID)
		so.Time = bk.Time
		so.Hash = hex.EncodeToString(bk.Hash)
		so.NodePosition = bk.NodePosition

		so.BlockTransactions = int64(bk.Tx)
		so.RollbacksHash = hex.EncodeToString(bk.RollbacksHash)
		so.BlockSizes = int64(len(bk.Data))
		ts, bdt, err := GetTransactionTxDetial(&bk)
		if err != nil {
			return err
		}
		so.CurrentVersion = strconv.FormatInt(int64(bdt.Header.Version), 10)
		so.BlockTransactionSize = ts
	}

	err = so.InsertRedisDb()
	if err != nil {
		return err
	}
	return nil
}

func (s *ScanOut) InsertRedisDb() error {
	errCs := s.Changes()
	if errCs != nil {
		return fmt.Errorf("changes err:%s\n", errCs.Error())
	}
	val, err := s.Marshal()
	if err != nil {
		return err
	}
	rp := RedisParams{
		Key:   ScanPrefix + strconv.FormatInt(s.Blockid, 10),
		Value: string(val),
	}

	for i := 0; i < 10; i++ {
		err = rp.SetExpire(time.Second * 10)
		if err == nil {
			break
		} else {
			time.Sleep(5 * time.Millisecond)
		}

	}

	return err
}

func (m *ScanOut) GetRedisdashboard() (*ScanOutRet, error) {
	var rets ScanOutRet
	f, err := m.GetRedisLastest()
	if err != nil {
		return &rets, err
	}
	if !f {
		return &rets, nil
	}

	//var bc BlockID
	//fc, _ := bc.GetbyName(consts.TransactionsMax)
	//if fc {
	//	rets.QueueTranscations = bc.ID
	//}

	rets.NodePosition = m.NodePosition
	rets.ConsensusMode = m.ConsensusMode

	rets.Blockid = m.Blockid
	rets.Hash = m.Hash
	rets.RollbacksHash = m.RollbacksHash
	//rets.Tx    = m.Blockid
	rets.KeyID = m.KeyID
	//rets.NodePosition = m.NodePosition
	rets.Time = m.Time
	rets.CurrentVersion = m.CurrentVersion

	rets.TotalCounts = m.TotalCounts
	//rets.TotalCapacitys = m.b
	rets.BlockSizes = m.BlockSizes
	rets.BlockTranscations = m.BlockTransactions
	rets.BlockTranscationSize = m.BlockTransactionSize
	rets.GuardianNodes = m.HonorNode
	rets.CastNodeInfo.CastNodes = m.CastNodes
	rets.SubNodes = m.SubNodes
	rets.CLBNodes = m.CLBNodes
	rets.MintAmounts = m.MintAmounts

	//dashboard date
	rets.BlockInfo.Blockid = m.Blockid
	rets.BlockInfo.MaxTps = m.MaxTps
	rets.BlockInfo.MaxBlockSize = m.MaxBlockSize
	rets.BlockInfo.StorageCapacitys = m.StorageCapacitys

	rets.TxInfo.TotalTx = m.TotalTx
	rets.TxInfo.TwentyFourTx = m.TwentyFourTx
	rets.TxInfo.SingleDayMaxTx = getSingleDayMaxTx()
	rets.TxInfo.WeekAverageTx = m.WeekAverageTx

	rets.CirculationsInfo.Circulations = m.Circulations
	rets.CirculationsInfo.TotalAmounts = TotalSupplyToken
	rets.CirculationsInfo.TodayCirculationsAmount = m.TodayCirculationsAmount
	rets.CirculationsInfo.TwentyFourAmount = m.TwentyFourAmount

	rets.NftMinerInfo.Count = m.NftMinerCount
	rets.NftMinerInfo.BlockReward = m.NftBlockReward
	rets.NftMinerInfo.HalveNumber = m.HalveNumber
	rets.NftMinerInfo.StakeAmounts = m.NftStakeAmounts

	rets.EcoLibsInfo.Ecosystems = m.EcosystemCount
	rets.EcoLibsInfo.EcoTokenTotal, rets.EcoLibsInfo.Contract, rets.EcoLibsInfo.DaoGovern = getEcoLibsInfo()

	rets.KeysInfo = getScanOutKeyInfo(1)
	rets.CandidateNodeInfo.CandidateNode, rets.CandidateNodeInfo.NodeVote, rets.CandidateNodeInfo.TwentyFourNode, rets.CandidateNodeInfo.NodeStakeAmounts = getScanOutNodeInfo()

	return &rets, err
}

func (m *ScanOut) GetRedisLastest() (bool, error) {
	rd := RedisParams{
		Key:   ScanOutStPrefix + ScanOutLastest,
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return false, nil
		}
		return false, err
	}
	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return false, err
	}

	return true, err
}

func TocapacityString(count int64) string {
	rs := float64(count) / float64(1024)
	if rs >= float64(1024) && rs < float64(1048576) {
		rs = rs / float64(1024)
		return strconv.FormatFloat(rs, 'f', 2, 64) + " MB"
	} else if rs >= float64(1048576) {
		rs = rs / float64(1024*1024)
		return strconv.FormatFloat(rs, 'f', 2, 64) + " GB"
	} else {
		return strconv.FormatFloat(rs, 'f', 2, 64) + " KB"
	}
}
func ToCapcityMb(count int64) string {
	rs := float64(count) / float64(1024)
	rs = rs / float64(1024)
	return strconv.FormatFloat(rs, 'f', 6, 64)
}

func ToCapcityKb(count int64) string {
	rs := float64(count) / float64(1024)
	return strconv.FormatFloat(rs, 'f', 5, 64)
}

func DealRedisBlockTpsList() error {
	return GetDBDealTraninfo(30)
}

func SendStatisticsSignal() {
	for len(GetScanOut) > 0 {
		<-GetScanOut
	}
	select {
	case GetScanOut <- true:
	default:
	}
	for len(SendScanOut) > 0 {
		<-SendScanOut
	}
	select {
	case SendScanOut <- true:
	default:
	}
}

func getAllContracts() int64 {
	tableName := "1_contracts"
	var total int64
	if err := GetDB(nil).Table(tableName).Count(&total).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getAllContracts err")
		return 0
	}
	return total
}

func getBlockContracts(blockId int64) int64 {
	var tss LogTransaction
	txList, err := tss.getHashListByBlockId(blockId)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getBlockContracts err")
		return 0
	}
	Thash := make(map[string]bool)
	for k := 0; k < len(txList); k++ {
		hash := hex.EncodeToString(txList[k].Hash)
		Thash[hash] = true
	}

	bk := &Block{}
	found, err := bk.GetId(blockId)
	var name []string
	if err == nil && found {
		list, err := GetBlocksContractNameList(bk)
		if err != nil {
			log.WithFields(log.Fields{"warn": err}).Warn("getBlockContracts name err")
			return 0
		}
		for key, value := range list {
			if Thash[key] {
				for i := 0; i < len(name); i++ {
					if name[i] == value {
						break
					}
				}
				name = append(name, value)
			}
		}

	} else {
		log.WithFields(log.Fields{"warn": err}).Warn("getBlockContracts name err")
		return 0
	}
	return int64(len(name))
}

func getMaxTps() (int64, int64) {
	var bk Block
	if err := GetDB(nil).Select("id,tx").Order("tx DESC").Last(&bk).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getMaxTps err")
		return 0, 0
	}
	return int64(bk.Tx), bk.ID
}

func getMaxBlockSize() (int64, int64) {
	var bk Block
	type block struct {
		ID     int64
		Length int64
	}
	var info block
	if err := GetDB(nil).Table(bk.TableName()).Select("id,length(data)").Order("length(data) DESC").Last(&info).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getMaxTps err")
		return 0, 0
	}

	return info.Length, info.ID
}

func getNftBlockReward() string {
	type applications struct {
		Id int64 `gorm:"column:id"`
	}
	var app applications
	f, err := isFound(GetDB(nil).Table("1_applications").Where("name = 'NFT' AND ecosystem = 1").First(&app))
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getNftBlockReward applications err")
		return "0"
	}
	if !f {
		return "0"
	}

	var params sqldb.AppParam
	f, err = isFound(GetDB(nil).Table(params.TableName()).Where("app_id = ? AND name = 'nft_per_reward' AND ecosystem = 1", app.Id).First(&params))
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getNftBlockReward app params err")
		return "0"
	}
	if !f {
		return "0"
	}

	money, err := smart.FormatMoney(nil, params.Value, 12)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getNftBlockReward format failed")
		return "0"
	}
	return money
}

func getTotalTx() int64 {
	var bk Block
	var ret SumInt64
	if err := GetDB(nil).Table(bk.TableName()).Select("sum(tx)").Take(&ret).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getTotalTx failed")
		return 0
	}
	return ret.Sum
}

func getSingleDayMaxTx() int64 {
	type reqParams struct {
		Date  string `gorm:"column:date"`
		TxNum int64  `gorm:"column:tx_num"`
	}
	var oneTx reqParams
	//if err := GetDB(nil).Table(bk.TableName()).Select("to_char(to_timestamp(time),'yyyy-MM-dd') as date,sum(tx) as tx_num").Group("date").Order("tx_num desc").First(&oneTx).Error; err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("getSingleDayMaxTx failed")
	//	return 0
	//}
	if err := GetDB(nil).Raw(`select to_char(to_timestamp("time"),'yyyy-MM-dd') as date,sum(tx) as tx_num 
FROM block_chain GROUP BY date ORDER BY tx_num DESC LIMIT 1`).Take(&oneTx).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getSingleDayMaxTx failed")
		return 0
	}
	return oneTx.TxNum
}

func getTwentyFourTx() int64 {
	tz := time.Unix(GetNowTimeUnix(), 0)
	t1 := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	var res SumInt64
	var bk Block
	if err := GetDB(nil).Table(bk.TableName()).Select("sum(tx)").Where("time >= ?", t1.Unix()).Take(&res).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getSingleDayMaxTx failed")
		return 0
	}
	return res.Sum
}

func getWeekAverageValueTx() int64 {
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	t1 := yesterday.AddDate(0, 0, -1*7)

	var res SumInt64
	var bk Block
	if err := GetDB(nil).Table(bk.TableName()).Select("sum(tx)").Where("time >= ? and time < ?", t1.AddDate(0, 0, 1).Unix(), yesterday.AddDate(0, 0, 1).Unix()).Take(&res).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getSingleDayMaxTx failed")
		return 0
	}
	if res.Sum != 0 {
		return res.Sum / 7
	} else {
		return 0
	}
}

func getHalveNumber() (int64, float64) {
	if !NftMinerReady {
		return 0, 0
	}
	var halvingNumber int64
	const halvingInterval int64 = 80000
	var stak NftMinerStaking
	var stakNumber int64
	subsidy := float64(25) * 100000000000

	if err := GetDB(nil).Table(stak.TableName()).Select("token_id").Group("token_id").Count(&stakNumber).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.WithFields(log.Fields{"warn": err}).Warn("getHavleNumber staking failed")
		}
		return 0, 0
	}
	if stakNumber >= halvingInterval {
		st1, _ := smart.Log(int64(4))
		st2, _ := smart.Log(smart.Float(stakNumber) / smart.Float(halvingInterval))
		num, _ := smart.Int(st2 / st1)
		halvingNumber = num + 1
	}
	if (halvingNumber) < 0 {
		halvingNumber = 0
	}
	subsidy = subsidy / math.Pow(float64(2), float64(halvingNumber))
	//str := smart.Int(smart.Log(smart.Float(stakNumber)/smart.Float(halvingInterval))/smart.Log(4)) + 1
	return halvingNumber, subsidy / 100000000000

}

func getEcoLibsInfo() (int64, int64, int64) {
	contractNum := getAllContracts()

	var ecosystems Ecosystem
	var symbolTotal int64
	var daoGovern int64
	if !HasTableOrView(nil, ecosystems.TableName()) {
		return 0, 0, 0
	}

	err := GetDB(nil).Table(ecosystems.TableName()).Where("id = 1 or token_symbol != ''").Count(&symbolTotal).Error
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getAllTokenCount symbolTotal failed")
		return 0, 0, 0
	}

	err = GetDB(nil).Table(ecosystems.TableName()).Where("control_mode = 2").Count(&daoGovern).Error
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getAllTokenCount daoGovern failed")
		return 0, 0, 0
	}

	return symbolTotal, contractNum, daoGovern
}

func getScanOutKeyInfo(ecosystem int64) KeysRet {
	var (
		key Key
		ret KeysRet
		err error
	)

	tz := time.Unix(GetNowTimeUnix(), 0)
	nowDay := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	t1 := nowDay.AddDate(0, 0, -1*30)
	if NftMinerReady {
		err = GetDB(nil).Table(key.TableName()).Select(`count(1) AS key_count,
(SELECT count(1) AS has_token_key FROM "1_keys" as k2 WHERE (k2.amount > 0 OR to_number(coalesce(NULLIF(k2.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') > 0) AND ecosystem = ?),
(SELECT count(1) AS month_active_key FROM(
SELECT sender_id as keyid FROM "1_history" WHERE sender_id <> 0 AND created_at >= ? and ecosystem = ? GROUP BY sender_id
 UNION 
SELECT recipient_id as keyid FROM "1_history" WHERE recipient_id <> 0 AND created_at >= ? AND ecosystem = ? GROUP BY recipient_id 
) AS tt)`, ecosystem, t1.Unix(), ecosystem, t1.Unix(), ecosystem).Where("ecosystem = ? AND id <> 0", ecosystem).Take(&ret).Error
	} else {
		err = GetDB(nil).Table(key.TableName()).Select(`count(1) AS key_count,
(SELECT count(1) AS has_token_key FROM "1_keys" as k2 WHERE k2.amount > 0 AND ecosystem = ?),
(SELECT count(1) AS month_active_key FROM(
SELECT sender_id as keyid FROM "1_history" WHERE sender_id <> 0 AND created_at >= ? and ecosystem = ? GROUP BY sender_id
 UNION 
SELECT recipient_id as keyid FROM "1_history" WHERE recipient_id <> 0 AND created_at >= ? AND ecosystem = ? GROUP BY recipient_id 
) AS tt)`, ecosystem, t1.Unix(), ecosystem, t1.Unix(), ecosystem).Where("ecosystem = ? AND id <> 0", ecosystem).Take(&ret).Error
	}
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getScanOutKeyInfo ecosystem keysRet failed")
		return ret
	}

	var bk Block
	f, err := bk.GetByTimeBlockId(nowDay.Unix())
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("getScanOutKeyInfo nowDay block failed")
		return ret
	}

	var rk sqldb.RollbackTx
	req := GetDB(nil).Table(rk.TableName())
	like := "%," + strconv.FormatInt(ecosystem, 10)
	if f {
		if err := req.Where("table_name = '1_keys' AND data = '' AND block_id >= ? AND table_id like ?", bk.ID, like).Count(&ret.TwentyFourKey).Error; err != nil {
			log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("getScanOutKeyInfo rollback oneDay failed")
			return ret
		}
	}
	return ret
}

func idIsRepeat(id int64, list []int64) bool {
	for j := 0; j < len(list); j++ {
		if list[j] == id || id == 0 {
			return true
		}
	}
	return false
}

func getScanOutNodeInfo() (candidate, nodeVote, twentyFour int64, stakeAmounts string) {
	if !NodeReady {
		return
	}
	var nodeStakeAmounts SumAmount
	var vote SumAmount

	var cn CandidateNodeRequests
	var ds CandidateNodeDecisions
	if err := GetDB(nil).Table(cn.TableName()).Where("deleted = 0").Count(&candidate).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get candidate failed")
		return
	}

	tz := time.Unix(GetNowTimeUnix(), 0)
	nowDay := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	if err := GetDB(nil).Table(cn.TableName()).Where("deleted = 0 AND date_created > ?", nowDay.Unix()).Count(&twentyFour).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get twenty Four failed")
		return
	}

	if err := GetDB(nil).Table(cn.TableName()).Select("sum(earnest_total)").Where("deleted = 0").Take(&nodeStakeAmounts).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get node Stake Amounts failed")
		return
	}
	stakeAmounts = nodeStakeAmounts.Sum.String()

	if err := GetDB(nil).Table(ds.TableName()).Select("sum(earnest)").Where(`decision_type = 1 AND decision <> 3 
		AND request_id IN (SELECT id FROM "1_candidate_node_requests" WHERE deleted = 0)`).Take(&vote).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get node Vote failed")
		return
	}
	moneyDec := decimal.NewFromInt(1e12)
	nodeVote = vote.Sum.Div(moneyDec).IntPart()

	return
}

func getDatabaseSize() (size string, err error) {
	if err = GetDB(nil).Raw("SELECT pg_size_pretty(pg_database_size(?))", conf.GetDbConn().Name).Scan(&size).Error; err != nil {
		return "", err
	}
	return
}
