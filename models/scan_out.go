/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/smart"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

var MaxBlockId int64

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

	BlockId         int64
	MaxTps          int32
	MaxBlockSize    string
	StorageCapacity string

	TotalTx         int64
	TwentyFourTx    int64
	WeekAverageTx   int64
	MaxActiveEcoLib string
	SingleDayMaxTx  int64

	Circulations            string
	TodayCirculationsAmount float64
	TwentyFourAmount        string

	NftMinerCount   int64
	NftBlockReward  float64
	HalveNumber     int64
	NftStakeAmounts string

	EcoLibsInfo       EcoLibsRet
	KeysInfo          KeysRet
	CandidateNodeInfo CandidateHonorNodeRet
}

type ScanOutRet struct {
	BlockId           int64  `json:"block_id"`    //Block Id
	BlockSizes        int64  `json:"block_sizes"` //Block Size
	BlockTransactions int64  `json:"block_transactions"`
	Hash              string `json:"hash"`
	RollbacksHash     string `json:"rollbacks_hash"`
	EcosystemID       int64  `json:"ecosystem_id"`
	KeyID             string `json:"key_id"`
	NodePosition      int64  `json:"node_position"`
	ConsensusMode     int32  `json:"consensus_mode"`
	Time              int64  `json:"time"`
	CurrentVersion    string `json:"current_version"`

	TotalCounts          int64  `json:"total_counts"` //total count
	BlockTransactionSize int64  `json:"block_transaction_size"`
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
	BlockId          int64  `json:"block_id"` //
	MaxTps           int32  `json:"max_tps"`
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

// Nft Miner
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
	BlockTransactions int64 `json:"block_transactions" ` //
}

type maxBlock struct {
	ID     int64
	Length int64
}

type maxTx struct {
	ID int64
	Tx int32
}

var (
	ScanOutStPrefix = "scan-out-"
	Latest          = "latest"
	GetScanOut      chan bool
)

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

func (m *ScanOut) InsertRedis() error {
	MaxBlockId = m.BlockId
	errCs := m.Changes()
	if errCs != nil {
		return fmt.Errorf("changes err:%s\n", errCs.Error())
	}
	val, err := m.Marshal()
	if err != nil {
		return err
	}

	rd := RedisParams{
		Key:   ScanOutStPrefix + Latest,
		Value: string(val),
	}
	err = rd.SetExpire(time.Millisecond * time.Duration(getTomorrowGapMilliseconds()))
	if err != nil {
		return err
	}

	return err
}

func GetRedisByName(name string) (any, error) {
	var rets any
	rd := RedisParams{
		Key:   name,
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return rets, nil
		}
		return rets, err
	}
	if err := msgpack.Unmarshal([]byte(rd.Value), &rets); err != nil {
		if err = json.Unmarshal([]byte(rd.Value), &rets); err != nil {
			return rets, err
		} else {
			return rets, nil
		}
	}
	return rets, nil
}

func GetScanOutDataToRedis() error {
	err := processScanOutBlocks()
	return err
}

func InitGlobalSwitch() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	NodeReady = CandidateTableExist()
	NftMinerReady = NftMinerTableIsExist()
	VotingReady = VotingTableExist()
	AssignReady = AssignTableExist()
	AirdropReady = AirdropTableExist()
	INameReady = INameTableExist()
}

func (ret *ScanOut) Changes() error {

	var ne NftMinerItems
	powerCount, err := ne.GetAllPower()
	if err != nil {
		return errors.New("Get All Power failed:" + err.Error())
	}
	ret.TotalCounts = powerCount

	cir, err := GetCirculations(1)
	if err != nil {
		return errors.New("Get Circulations failed:" + err.Error())
	}

	var nft NftMinerStaking
	_, nftStaking, err := nft.GetAllStakeAmount()
	if err != nil {
		return errors.New("Get All Stake Amount failed:" + err.Error())
	}

	//stakeamount, err := key.GetStakeAmount()
	//if err != nil {
	//	return err
	//}

	ret.Circulations = cir
	ret.NftStakeAmounts = nftStaking.String()

	ret.EcoLibsInfo, err = getEcoLibsInfo()
	if err != nil {
		return err
	}

	ret.KeysInfo, err = getScanOutKeyInfoFromRedis()
	if err != nil {
		return err
	}

	ret.CandidateNodeInfo, err = getScanOutNodeInfo()
	if err != nil {
		return err
	}

	var mst MinePledgeStatus
	honor, casts, nftCount, err := mst.GetCastNodeandGuardianNode()
	if err != nil {
		return fmt.Errorf("get Cast Nodeand Guardian Node failed:%s", err.Error())
	}
	ret.HonorNode = honor
	ret.CastNodes = casts
	ret.NftMinerCount = nftCount

	capacity, err := getDatabaseSize()
	if err != nil {
		return fmt.Errorf("get Database Size failed:%s", err.Error())
	}
	ret.StorageCapacity = capacity

	var sp StateParameter
	sp.ecosystem = 1
	mb, err := sp.GetMintAmount()
	if err != nil {
		return fmt.Errorf("get Mint Amount failed:%s", err.Error())
	}
	ret.MintAmounts = mb
	ret.HalveNumber, ret.NftBlockReward, err = getHalveNumber()
	if err != nil {
		return err
	}

	var his History
	ret.TodayCirculationsAmount, err = his.GetTodayCirculationsAmount(ret.NftBlockReward)
	if err != nil {
		return err
	}

	ret.TwentyFourAmount, err = his.Get24HourTxAmount()
	if err != nil {
		return err
	}

	maxTx, err := getMaxTxFromRedis()
	if err != nil {
		return err
	}
	ret.MaxTps = maxTx.Tx

	maxBlock, err := getMaxBlockSizeFromRedis()
	if err != nil {
		return err
	}

	ret.MaxBlockSize = ToCapacityString(maxBlock.Length)

	ret.TotalTx, err = getTotalTx()
	if err != nil {
		return err
	}

	ret.TwentyFourTx, err = getTwentyFourTx()
	if err != nil {
		return err
	}

	ret.WeekAverageTx, err = getWeekAverageValueTxFromRedis()
	if err != nil {
		return err
	}

	ret.MaxActiveEcoLib, err = GetActiveEcoLibsFromRedis()
	if err != nil {
		return err
	}

	ret.SingleDayMaxTx, err = getSingleDayMaxTxFromRedis()
	if err != nil {
		return err
	}

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
	ret.BlockId = cbk.BlockID
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

	err = ret.InsertRedis()
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
		so.BlockId = 1
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

	err = so.InsertRedis()
	if err != nil {
		return err
	}
	return nil
}

func (m *ScanOut) GetDashboardFromRedis() (*ScanOutRet, error) {
	var rets ScanOutRet
	f, err := m.GetRedisLatest()
	if err != nil {
		return &rets, err
	}
	if !f {
		return &rets, nil
	}

	rets.NodePosition = m.NodePosition
	rets.ConsensusMode = m.ConsensusMode

	rets.BlockId = m.BlockId
	rets.Hash = m.Hash
	rets.RollbacksHash = m.RollbacksHash
	rets.KeyID = m.KeyID
	rets.Time = m.Time
	rets.CurrentVersion = m.CurrentVersion

	rets.TotalCounts = m.TotalCounts
	rets.BlockSizes = m.BlockSizes
	rets.BlockTransactions = m.BlockTransactions
	rets.BlockTransactionSize = m.BlockTransactionSize
	rets.GuardianNodes = m.HonorNode
	rets.CastNodeInfo.CastNodes = m.CastNodes
	rets.SubNodes = m.SubNodes
	rets.CLBNodes = m.CLBNodes
	rets.MintAmounts = m.MintAmounts

	//dashboard date
	rets.BlockInfo.BlockId = m.BlockId
	rets.BlockInfo.MaxTps = m.MaxTps
	rets.BlockInfo.MaxBlockSize = m.MaxBlockSize
	rets.BlockInfo.StorageCapacitys = m.StorageCapacity

	rets.TxInfo.TotalTx = m.TotalTx
	rets.TxInfo.TwentyFourTx = m.TwentyFourTx
	rets.TxInfo.SingleDayMaxTx = m.SingleDayMaxTx
	rets.TxInfo.WeekAverageTx = m.WeekAverageTx

	rets.CirculationsInfo.Circulations = m.Circulations
	rets.CirculationsInfo.TotalAmounts = TotalSupplyToken
	rets.CirculationsInfo.TodayCirculationsAmount = m.TodayCirculationsAmount
	rets.CirculationsInfo.TwentyFourAmount = m.TwentyFourAmount

	rets.NftMinerInfo.Count = m.NftMinerCount
	rets.NftMinerInfo.BlockReward = m.NftBlockReward
	rets.NftMinerInfo.HalveNumber = m.HalveNumber
	rets.NftMinerInfo.StakeAmounts = m.NftStakeAmounts
	rets.EcoLibsInfo = m.EcoLibsInfo
	rets.KeysInfo = m.KeysInfo
	rets.CandidateNodeInfo = m.CandidateNodeInfo

	return &rets, err
}

func (m *ScanOut) GetRedisLatest() (bool, error) {
	rd := RedisParams{
		Key:   ScanOutStPrefix + Latest,
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

func ToCapacityString(count int64) string {
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
func ToCapacityMb(count int64) string {
	rs := float64(count) / float64(1024)
	rs = rs / float64(1024)
	return strconv.FormatFloat(rs, 'f', 6, 64)
}

func ToCapacityKb(count int64) string {
	rs := float64(count) / float64(1024)
	return strconv.FormatFloat(rs, 'f', 5, 64)
}

func DealRedisBlockTpsList() error {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	return GetDBDealTraninfo(30)
}

func GetStatisticsSignal() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	for len(GetScanOut) > 0 {
		<-GetScanOut
	}
	select {
	case GetScanOut <- true:
	default:
	}

	for len(SendWebsocketSignal) > 0 {
		<-SendWebsocketSignal
	}
	select {
	case SendWebsocketSignal <- true:
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

func getMaxTx() (maxTx, error) {
	var (
		bk   Block
		rets maxTx
	)
	if err := GetDB(nil).Select("id,tx").Order("tx DESC").Last(&bk).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get Max tx err")
		return rets, err
	}
	rets.ID = bk.ID
	rets.Tx = bk.Tx
	return rets, nil
}

func GetMaxTxToRedis() {
	HistoryWG.Add(1)
	defer func() {
		HistoryWG.Done()
	}()
	rets, err := getMaxTx()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get max tx error")
		return
	}

	res, err := msgpack.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get max tx msgpack error")
		return
	}

	rd := RedisParams{
		Key:   "max_tx",
		Value: string(res),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get max tx set redis error")
		return
	}
}

func getMaxTxFromRedis() (maxTx, error) {
	var rets maxTx
	var err error
	rd := RedisParams{
		Key:   "max_tx",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get max tx From Redis getDb err")
		return rets, err
	}
	err = msgpack.Unmarshal([]byte(rd.Value), &rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get max tx From Redis msgpack err")
		return rets, err
	}

	return rets, nil
}

func getMaxBlock() (maxBlock, error) {
	var bk Block
	var info maxBlock
	if err := GetDB(nil).Table(bk.TableName()).Select("id,length(data)").Order("length(data) DESC").Last(&info).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get max block err")
		return info, err
	}

	return info, nil
}

func GetMaxBlockSizeToRedis() {
	HistoryWG.Add(1)
	defer func() {
		HistoryWG.Done()
	}()
	rets, err := getMaxBlock()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get max block error")
		return
	}

	res, err := msgpack.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get max block msgpack error")
		return
	}

	rd := RedisParams{
		Key:   "max_block_size",
		Value: string(res),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get max block set redis error")
		return
	}
}

func getMaxBlockSizeFromRedis() (maxBlock, error) {
	var rets maxBlock
	var err error
	rd := RedisParams{
		Key:   "max_block_size",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get max block From Redis getDb err")
		return rets, err
	}
	err = msgpack.Unmarshal([]byte(rd.Value), &rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get max block From Redis msgpack err")
		return rets, err
	}

	return rets, nil
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

func getTotalTx() (int64, error) {
	var bk Block
	var ret SumInt64
	if err := GetDB(nil).Table(bk.TableName()).Select("sum(tx)").Take(&ret).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getTotalTx failed")
		return 0, err
	}
	return ret.Sum, nil
}

func getSingleDayMaxTx() (int64, error) {
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
		return 0, err
	}
	return oneTx.TxNum, nil
}

func GetSingleDayMaxTxToRedis() {
	HistoryWG.Add(1)
	defer func() {
		HistoryWG.Done()
	}()
	rets, err := getSingleDayMaxTx()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get single day max tx error")
		return
	}

	res, err := msgpack.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get single day max tx msgpack error")
		return
	}

	rd := RedisParams{
		Key:   "day_max_tx",
		Value: string(res),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get single day max tx set redis error")
		return
	}
}

func getSingleDayMaxTxFromRedis() (int64, error) {
	var rets int64
	var err error
	rd := RedisParams{
		Key:   "day_max_tx",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get Single Day Max Tx From Redis getDb err")
		return rets, err
	}
	err = msgpack.Unmarshal([]byte(rd.Value), &rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get Single Day Max Tx From Redis msgpack err")
		return rets, err
	}

	return rets, nil
}

func getTwentyFourTx() (int64, error) {
	tz := time.Unix(GetNowTimeUnix(), 0)
	t1 := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	var res SumInt64
	var bk Block
	if err := GetDB(nil).Table(bk.TableName()).Select("sum(tx)").Where("time >= ?", t1.Unix()).Take(&res).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getSingleDayMaxTx failed")
		return 0, err
	}
	return res.Sum, nil
}

func getWeekAverageValueTx() (int64, error) {
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	t1 := yesterday.AddDate(0, 0, -1*7)

	var res SumInt64
	var bk Block
	if err := GetDB(nil).Table(bk.TableName()).Select("sum(tx)").Where("time >= ? and time < ?", t1.AddDate(0, 0, 1).Unix(), yesterday.AddDate(0, 0, 1).Unix()).Take(&res).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getSingleDayMaxTx failed")
		return 0, err
	}
	if res.Sum != 0 {
		return res.Sum / 7, nil
	} else {
		return 0, nil
	}
}

func GetWeekAverageValueTxToRedis() {
	HistoryWG.Add(1)
	defer func() {
		HistoryWG.Done()
	}()
	rets, err := getWeekAverageValueTx()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get week avg tx error")
		return
	}

	rd := RedisParams{
		Key:   "week_tx_avg",
		Value: strconv.FormatInt(rets, 10),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get week avg tx set redis error")
		return
	}
}

func getWeekAverageValueTxFromRedis() (int64, error) {
	rd := RedisParams{
		Key:   "week_tx_avg",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get week avg tx From Redis getDb err")
		return 0, err
	}
	avgTx, err := strconv.ParseInt(rd.Value, 10, 64)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get week avg tx From Redis prase int err")
		return 0, err
	}

	return avgTx, nil
}

func getHalveNumber() (int64, float64, error) {
	if !NftMinerReady {
		return 0, 0, nil
	}
	var (
		halvingNumber int64
		interval      int64
		reward        float64
	)
	var stak NftMinerStaking
	var stakNumber int64

	var app AppParam
	f, err := app.GetByName(1, "nft_miner_per_reward")
	if err != nil {
		return 0, 0, err
	}
	if !f {
		return 0, 0, nil
	}
	reward, err = strconv.ParseFloat(app.Value, 64)
	if err != nil {
		return 0, 0, err
	}

	intervalStr, err := GetAppValue(app.AppID, "nft_miner_halving_interval", 1)
	if err != nil {
		return 0, 0, err
	}
	interval, err = strconv.ParseInt(intervalStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	if err := GetDB(nil).Table(stak.TableName()).Select("token_id").Group("token_id").Count(&stakNumber).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.WithFields(log.Fields{"warn": err}).Warn("get Havle Number staking failed")
		}
		return 0, 0, err
	}
	if stakNumber >= interval {
		st1, _ := smart.Log(int64(4))
		st2, _ := smart.Log(smart.Float(stakNumber) / smart.Float(interval))
		num, _ := smart.Int(st2 / st1)
		halvingNumber = num + 1
	}
	if (halvingNumber) < 0 {
		halvingNumber = 0
	}
	return halvingNumber, reward / 1e12, nil

}

func getEcoLibsInfo() (EcoLibsRet, error) {
	var rets EcoLibsRet

	var ecosystems Ecosystem
	if !HasTableOrView(ecosystems.TableName()) {
		return rets, nil
	}

	contractNum := getAllContracts()
	var symbolTotal int64
	var daoGovern int64
	var total int64

	if err := GetDB(nil).Table(ecosystems.TableName()).Count(&total).Error; err != nil {
		return rets, err
	}

	err := GetDB(nil).Table(ecosystems.TableName()).Where("id = 1 or token_symbol != ''").Count(&symbolTotal).Error
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getAllTokenCount symbolTotal failed")
		return rets, err
	}

	err = GetDB(nil).Table(ecosystems.TableName()).Where("control_mode = 2").Count(&daoGovern).Error
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getAllTokenCount daoGovern failed")
		return rets, err
	}
	rets.DaoGovern = daoGovern
	rets.Contract = contractNum
	rets.EcoTokenTotal = symbolTotal
	rets.Ecosystems = total

	return rets, nil
}

func getScanOutKeyInfo(ecosystem int64) (KeysRet, error) {
	var (
		key Key
		ret KeysRet
		err error
	)

	tz := time.Unix(GetNowTimeUnix(), 0)
	nowDay := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	t1 := nowDay.AddDate(0, 0, -1*30)
	if NftMinerReady || NodeReady {
		err = GetDB(nil).Table(key.TableName()).Select(`count(1) AS key_count,
(SELECT count(1) AS has_token_key FROM(
	SELECT id FROM "1_keys" WHERE (amount > 0 OR 
		to_number(coalesce(NULLIF(lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') > 0 OR
		to_number(coalesce(NULLIF(lock->>'candidate_referendum',''),'0'),'999999999999999999999999') > 0 OR 
		to_number(coalesce(NULLIF(lock->>'candidate_substitute',''),'0'),'999999999999999999999999') > 0) AND ecosystem = ?
	UNION
	SELECT output_key_id FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ?
)AS v1),
(SELECT count(1) AS month_active_key FROM(
	SELECT sender_id as keyid FROM "1_history" WHERE created_at >= ? and ecosystem = ? GROUP BY sender_id
 		UNION 
	SELECT recipient_id as keyid FROM "1_history" WHERE created_at >= ? AND ecosystem = ? GROUP BY recipient_id
		UNION
	SELECT output_key_id AS keyid FROM spent_info AS s1 LEFT JOIN 
	 log_transactions AS l1 ON(l1.hash = s1.output_tx_hash)	WHERE ecosystem = ? AND timestamp >= ? GROUP BY output_key_id
		UNION
	SELECT output_key_id AS keyid FROM spent_info AS s1 LEFT JOIN 
	 log_transactions AS l1 ON(l1.hash = s1.input_tx_hash)	WHERE ecosystem = ? AND timestamp >= ? GROUP BY output_key_id
) AS tt)`, ecosystem, ecosystem, t1.Unix(), ecosystem, t1.Unix(), ecosystem, ecosystem, t1.UnixMilli(), ecosystem, t1.UnixMilli()).
			Where("ecosystem = ?", ecosystem).Take(&ret).Error
	} else {
		err = GetDB(nil).Table(key.TableName()).Select(`count(1) AS key_count,
(SELECT count(1) AS has_token_key FROM(
	SELECT id FROM "1_keys" WHERE amount > 0 AND ecosystem = ?
	UNION
	SELECT output_key_id FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ?
)AS v1),
(SELECT count(1) AS month_active_key FROM(
	SELECT sender_id as keyid FROM "1_history" WHERE created_at >= ? and ecosystem = ? GROUP BY sender_id
	 UNION 
	SELECT recipient_id as keyid FROM "1_history" WHERE created_at >= ? AND ecosystem = ? GROUP BY recipient_id
	 UNION
	SELECT output_key_id AS keyid FROM spent_info AS s1 LEFT JOIN 
	 log_transactions AS l1 ON(l1.hash = s1.output_tx_hash)	WHERE ecosystem = ? AND timestamp >= ? GROUP BY output_key_id
	 UNION
	SELECT output_key_id AS keyid FROM spent_info AS s1 LEFT JOIN 
	 log_transactions AS l1 ON(l1.hash = s1.input_tx_hash)	WHERE ecosystem = ? AND timestamp >= ? GROUP BY output_key_id
) AS tt)`, ecosystem, ecosystem, t1.Unix(), ecosystem, t1.Unix(), ecosystem, ecosystem, t1.UnixMilli(), ecosystem, t1.UnixMilli()).
			Where("ecosystem = ?", ecosystem).Take(&ret).Error
	}
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("getScanOutKeyInfo ecosystem keysRet failed")
		return ret, err
	}

	bk := &Block{}
	f, err := bk.GetByTimeBlockId(nil, nowDay.Unix())
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("getScanOutKeyInfo nowDay block failed")
		return ret, err
	}

	var rk sqldb.RollbackTx
	req := GetDB(nil).Table(rk.TableName())
	like := "%," + strconv.FormatInt(ecosystem, 10)
	if f {
		if err := req.Where("table_name = '1_keys' AND data = '' AND block_id >= ? AND table_id like ?", bk.ID, like).Count(&ret.TwentyFourKey).Error; err != nil {
			log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("getScanOutKeyInfo rollback oneDay failed")
			return ret, err
		}
	}
	return ret, nil
}

func GetScanOutKeyInfoToRedis() {
	HistoryWG.Add(1)
	defer func() {
		HistoryWG.Done()
	}()
	rets, err := getScanOutKeyInfo(1)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get Scan Out Key Info error")
		return
	}

	res, err := msgpack.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get Scan Out Key Info msgpack error")
		return
	}

	rd := RedisParams{
		Key:   "key_info",
		Value: string(res),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get Scan Out Key Info set redis error")
		return
	}
}

func getScanOutKeyInfoFromRedis() (KeysRet, error) {
	var rets KeysRet
	var err error
	rd := RedisParams{
		Key:   "key_info",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return rets, nil
		}
		log.WithFields(log.Fields{"warn": err}).Warn("get Scan Out Key Info From Redis getDb err")
		return rets, err
	}
	err = msgpack.Unmarshal([]byte(rd.Value), &rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get Scan Out Key Info From Redis msgpack err")
		return rets, err
	}

	return rets, nil
}

func idIsRepeat(id int64, list []int64) bool {
	for j := 0; j < len(list); j++ {
		if list[j] == id || id == 0 {
			return true
		}
	}
	return false
}

func getScanOutNodeInfo() (CandidateHonorNodeRet, error) {
	var rets CandidateHonorNodeRet
	if !NodeReady {
		return rets, nil
	}
	var nodeStakeAmounts SumAmount
	var vote SumAmount

	var cn CandidateNodeRequests
	var ds CandidateNodeDecisions
	if err := GetDB(nil).Table(cn.TableName()).Where("deleted = 0").Count(&rets.CandidateNode).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get candidate failed")
		return rets, err
	}

	tz := time.Unix(GetNowTimeUnix(), 0)
	nowDay := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	if err := GetDB(nil).Table(cn.TableName()).Where("deleted = 0 AND date_created > ?", nowDay.Unix()).Count(&rets.TwentyFourNode).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get twenty Four failed")
		return rets, err
	}

	if err := GetDB(nil).Table(cn.TableName()).Select("sum(earnest_total)").Where("deleted = 0").Take(&nodeStakeAmounts).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get node Stake Amounts failed")
		return rets, err
	}
	rets.NodeStakeAmounts = nodeStakeAmounts.Sum.String()

	if err := GetDB(nil).Table(ds.TableName()).Select("sum(earnest)").Where(`decision_type = 1 AND decision <> 3 
		AND request_id IN (SELECT id FROM "1_candidate_node_requests" WHERE deleted = 0)`).Take(&vote).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get node Vote failed")
		return rets, err
	}
	moneyDec := decimal.NewFromInt(1e12)
	rets.NodeVote = vote.Sum.Div(moneyDec).IntPart()

	return rets, nil
}

func getDatabaseSize() (size string, err error) {
	if err = GetDB(nil).Raw("SELECT pg_size_pretty(pg_database_size(?))", conf.GetDbConn().Name).Scan(&size).Error; err != nil {
		return "", err
	}
	return
}

func getTomorrowGapMilliseconds() int64 {
	d1 := time.Now()
	tomorrow := time.Date(d1.Year(), d1.Month(), d1.Day(), 0, 0, 0, 0, d1.Location()).AddDate(0, 0, 1)
	return tomorrow.Sub(d1).Milliseconds()
}
