/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"errors"
	"github.com/IBAX-io/go-explorer/storage"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"sort"
	"strconv"
	"time"
)

var SendWebsocketSignal chan bool

type BlockListChart struct {
	Time  []int64 `json:"time"`
	Block []int64 `json:"block"`
}

type TxListChart struct {
	Name string  `json:"name,omitempty"`
	Time []int64 `json:"time"`
	Tx   []int64 `json:"tx"`
}

type Circulations struct {
	TotalAmount        string `json:"total_amount"`
	CirculationsAmount string `json:"circulations_amount"`
	StakeAmounts       string `json:"stake_amounts"`
	FreezeAmount       string `json:"freeze_amount"`
}

type NftMinerInfoChart struct {
	Count           int64 `json:"count"`
	NowStakingCount int64 `json:"now_staking_count"`
	UnStakingCount  int64 `json:"un_staking_count"`
}

type EcoListChart struct {
	Time    []int64 `json:"time"`
	EcoLibs []int64 `json:"eco_libs"`
}

type KeyInfoChart struct {
	Name      string  `json:"name,omitempty"`
	Time      []int64 `json:"time"`
	NewKey    []int64 `json:"new_key,omitempty"`
	ActiveKey []int64 `json:"active_key,omitempty"`
}

type HonorNodeChart struct {
	Time          []string `json:"time"`
	HonorNode     []string `json:"honor_node"`
	CandidateNode []string `json:"candidate_node"`
}

type NewKeyHistoryChart struct {
	Time   []string `json:"time"`
	NewKey []int64  `json:"new_key"`
}

type DashboardChartData struct {
	BlockChart        BlockListChart    `json:"block_chart"`
	TxChart           TxListChart       `json:"tx_chart"`
	CirculationsChart Circulations      `json:"circulations_chart"`
	NftMinerChart     NftMinerInfoChart `json:"nft_miner_chart"`
	EcoLibsChart      EcoListChart      `json:"eco_libs_chart"`
	KeyChart          KeyInfoChart      `json:"key_chart"`
	HonerNodeChart    HonorNodeChart    `json:"honer_node_chart"`
}

type DaysNumber struct {
	Days string `gorm:"column:days"`
	Num  int64  `gorm:"column:num"`
}

func FindDaysNumber(sql string, values ...interface{}) ([]DaysNumber, error) {
	var (
		list []DaysNumber
		err  error
	)
	if values == nil {
		err = GetDB(nil).Raw(sql).Find(&list).Error
	} else {
		err = GetDB(nil).Raw(sql, values...).Find(&list).Error
	}
	return list, err
}

func FindDaysAmount(sql string) ([]DaysAmount, error) {
	var (
		list []DaysAmount
	)
	err := GetDB(nil).Raw(sql).Find(&list).Error
	return list, err
}

func (s *DashboardChartData) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(s); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (s *DashboardChartData) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &s); err != nil {
		return err
	}
	return nil
}

func GetDashboardChartDataToRedis() {
	ChartWG.Add(1)
	defer func() {
		ChartWG.Done()
	}()
	rets, err := GetDashboardChartData()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetDashboardChartData error")
		return
	}

	val, err := rets.Marshal()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetDashboardChartData marshal error")
		return
	}

	//date, err := json.Marshal(rets)
	//if err != nil {
	//	return err
	//}

	rd := RedisParams{
		Key:   "dashboard-chart",
		Value: string(val),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetDashboardChartData set redis error")
		return
	}
}

func GetDashboardChartDataFromRedis() (*DashboardChartData, error) {
	var err error
	rets := &DashboardChartData{}
	rd := RedisParams{
		Key:   "dashboard-chart",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("Get Dashboard Chart Data From Redis getDb err")
		return nil, err
	}
	err = rets.Unmarshal([]byte(rd.Value))
	if err != nil {
		return nil, err
	}
	//if err = json.Unmarshal([]byte(rd.Value), &dash); err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("GetDashboardChartDataFromRedis json err")
	//	return nil, err
	//}

	cir, err := GetCirculations(1)
	if err != nil {
		return rets, err
	}

	var miner NftMinerStaking
	stakingNum, nftStaking, err := miner.GetAllStakeAmount()
	if err != nil {
		return rets, err
	}

	var agi AssignGetInfo
	agm, err := agi.GetAllBalance(nil)
	if err != nil {
		return rets, err
	}

	if NftMinerReady {
		err := GetDB(nil).Table("1_nft_miner_items").Where("merge_status = ? ", 1).Count(&rets.NftMinerChart.Count).Error
		if err != nil {
			return rets, err
		}
	}
	var nodeStaking SumAmount
	if NodeReady {
		err = GetDB(nil).Table("1_candidate_node_decisions").Select("coalesce(sum(earnest),'0')as sum").Where("decision <> 3").Take(&nodeStaking.Sum).Error
		if err != nil {
			return rets, err
		}
	}

	rets.CirculationsChart.TotalAmount = TotalSupplyToken
	rets.CirculationsChart.CirculationsAmount = cir
	rets.CirculationsChart.StakeAmounts = nftStaking.Add(nodeStaking.Sum).String()
	rets.CirculationsChart.FreezeAmount = agm.String()

	rets.NftMinerChart.UnStakingCount = rets.NftMinerChart.Count - stakingNum
	rets.NftMinerChart.NowStakingCount = stakingNum

	return rets, nil
}

func GetDashboardChartData() (*DashboardChartData, error) {
	var (
		rets     DashboardChartData
		bkList   []DaysNumber
		txList   []DaysNumber
		bkChart  BlockListChart
		txChart  TxListChart
		ecoChart EcoListChart
		keyChart KeyInfoChart
	)
	const getWeek = 7
	const getFifteenDays = 15

	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	t1 := yesterday.AddDate(0, 0, -1*getWeek)
	t2 := yesterday.AddDate(0, 0, -1*getFifteenDays)

	err := GetDB(nil).Raw(`SELECT to_char(to_timestamp(time),'yyyy-MM-dd') days,sum(tx) num FROM 
"block_chain" WHERE time >= ? GROUP BY days`, t2.Unix()).Find(&txList).Error
	if err != nil {
		return &rets, err
	}
	txChart.Time = make([]int64, getFifteenDays)
	txChart.Tx = make([]int64, getFifteenDays)
	for i := 0; i < len(txChart.Time); i++ {
		txChart.Time[i] = t2.AddDate(0, 0, i+1).Unix()
		txChart.Tx[i] = GetDaysNumber(txChart.Time[i], txList)
	}

	err = GetDB(nil).Raw(`SELECT to_char(to_timestamp(time),'yyyy-MM-dd') days,count(*) num FROM 
"block_chain" WHERE time >= ? GROUP BY days`, t1.Unix()).Find(&bkList).Error
	if err != nil {
		return &rets, err
	}
	bkChart.Time = make([]int64, getWeek)
	bkChart.Block = make([]int64, getWeek)
	for i := 0; i < len(bkChart.Time); i++ {
		bkChart.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		bkChart.Block[i] = GetDaysNumber(bkChart.Time[i], bkList)
	}
	rets.BlockChart = bkChart
	rets.TxChart = txChart

	var his []History
	if err := GetDB(nil).Select("created_at,id").Where("created_at >= ? AND comment = 'taxes for execution of @1NewEcosystem contract' AND type = 1", t1.UnixMilli()).Find(&his).Error; err != nil {
		return &rets, err
	}
	ecoChart.Time = make([]int64, getWeek)
	ecoChart.EcoLibs = make([]int64, getWeek)
	for i := 0; i < len(ecoChart.Time); i++ {
		ecoChart.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		stTime := t1.AddDate(0, 0, i+1).UnixMilli()
		endTime := t1.AddDate(0, 0, i+2).UnixMilli()
		ecoChart.EcoLibs[i] = getTimeLineEcoListInfo(stTime, endTime, his)
	}
	rets.EcoLibsChart = ecoChart

	keyChart.Time = make([]int64, getFifteenDays)
	keyChart.ActiveKey = make([]int64, getFifteenDays)
	keyChart.NewKey = make([]int64, getFifteenDays)
	idList := make([]int64, getFifteenDays)

	var activeList []DaysNumber
	err = GetDB(nil).Raw(`SELECT days,count(keyid) as num  FROM (

SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') days ,sender_id as keyid FROM "1_history" WHERE sender_id <> 0 AND created_at >= ? AND ecosystem = 1 GROUP BY days, sender_id
 UNION 
SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') days , recipient_id as keyid  FROM "1_history" WHERE recipient_id <> 0 AND created_at >= ? AND ecosystem = 1 GROUP BY days,  recipient_id 

) as tt GROUP BY days ORDER BY days desc `, t2.Unix(), t2.Unix()).Find(&activeList).Error
	if err != nil {
		return &rets, err
	}

	for i := 0; i < len(keyChart.Time); i++ {
		keyChart.Time[i] = t2.AddDate(0, 0, i+1).Unix()
		var bks Block
		f, err := bks.GetByTimeBlockId(nil, keyChart.Time[i])
		if err != nil {
			return &rets, err
		}
		if f {
			idList[i] = bks.ID
		}
		keyChart.ActiveKey[i] = GetDaysNumber(keyChart.Time[i], activeList)

	}
	for i := 0; i < len(idList); i++ {
		if i == len(idList)-1 {
			keyChart.NewKey[i] = getNewKeyNumber(idList[i], 0, 1)
		} else {
			keyChart.NewKey[i] = getNewKeyNumber(idList[i], idList[i+1], 1)
		}
	}
	rets.KeyChart = keyChart

	if NodeReady {
		type honorNodeChange struct {
			Days         string
			HonorLen     string
			CandidateLen string
		}
		var nelist []honorNodeChange
		getNodeNumber := func(getTime string, list []honorNodeChange) (string, string) {
			for i := 0; i < len(list); i++ {
				if getTime == list[i].Days {
					return list[i].HonorLen, list[i].CandidateLen
				}
			}
			return "0", "0"
		}
		err = GetDB(nil).Raw(`
SELECT to_char(to_timestamp(time),'yyyy-MM-dd') AS days,
		CASE WHEN honor_node::text != 'null' THEN 
			JSONB_ARRAY_LENGTH("honor_node"::jsonb)
		ELSE
			0
		END	honor_len,
		CASE WHEN candidate_node::text != 'null' THEN
			JSONB_ARRAY_LENGTH("candidate_node"::jsonb)
		ELSE
			0
		END candidate_len FROM daily_node_report WHERE time >= ?
ORDER BY days DESC LIMIT 15
`, t2.Unix()).Find(&nelist).Error
		if err != nil {
			return &rets, err
		}
		var lastHonor = "0"
		var lastCandidate = "0"
		rets.HonerNodeChart.Time = make([]string, getFifteenDays)
		rets.HonerNodeChart.HonorNode = make([]string, getFifteenDays)
		rets.HonerNodeChart.CandidateNode = make([]string, getFifteenDays)
		for i := 0; i < len(rets.HonerNodeChart.Time); i++ {
			rets.HonerNodeChart.Time[i] = t2.AddDate(0, 0, i+1).Format("2006-01-02")
			honorLen, candidateLen := getNodeNumber(rets.HonerNodeChart.Time[i], nelist)
			if honorLen == "0" && candidateLen == "0" {
				rets.HonerNodeChart.HonorNode[i] = lastHonor
				rets.HonerNodeChart.CandidateNode[i] = lastCandidate
			} else {
				rets.HonerNodeChart.HonorNode[i] = honorLen
				rets.HonerNodeChart.CandidateNode[i] = candidateLen
				lastHonor = honorLen
				lastCandidate = candidateLen
			}
		}
	} else {
		rets.HonerNodeChart.Time = make([]string, 0)
		rets.HonerNodeChart.HonorNode = make([]string, 0)
		rets.HonerNodeChart.CandidateNode = make([]string, 0)
	}

	return &rets, nil
}

func GetDaysNumber(getTime int64, list []DaysNumber) int64 {
	for i := len(list); i > 0; i-- {
		times, _ := time.ParseInLocation("2006-01-02", list[i-1].Days, time.Local)
		if getTime == times.Unix() {
			return list[i-1].Num
		}
	}
	return 0
}

func GetDaysNumberLike(getTime int64, list []DaysNumber, getTimeFront bool, order string) int64 {
	if getTimeFront {
		switch order {
		case "desc":
			for i := 0; i < len(list); i++ {
				times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
				if times.Unix() <= getTime {
					return list[i].Num
				}
			}
		case "asc":
			for i := len(list) - 1; i > 0; i-- {
				times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
				if times.Unix() <= getTime {
					return list[i].Num
				}
			}
		default:
			return 0
		}
	} else {
		switch order {
		case "desc":
			for i := len(list) - 1; i > 0; i-- {
				times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
				if times.Unix() >= getTime {
					return list[i].Num
				}
			}
		case "asc":
			for i := 0; i < len(list); i++ {
				times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
				if times.Unix() >= getTime {
					return list[i].Num
				}
			}
		default:
			return 0
		}
	}
	return 0
}

func GetDaysAmountResponse(list []DaysAmount) DaysAmountResponse {
	var (
		rets DaysAmountResponse
		err  error
	)
	layout := "2006-01-02"
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	var startTime time.Time
	if len(list) > 0 {
		startTime, err = time.ParseInLocation(layout, list[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"warn": err, "Days": list[0].Days}).Warn("Get Days Amount Response ParseInLocation Failed")
			return rets
		}
		for startTime.Unix() <= today.Unix() {
			rets.Time = append(rets.Time, startTime.Unix())
			rets.Amount = append(rets.Amount, GetDaysAmount(startTime.Unix(), list))
			startTime = startTime.AddDate(0, 0, 1)
		}
	}
	return rets
}

func getNewKeyNumber(startId, endId int64, ecosystem int64) int64 {
	var rk sqldb.RollbackTx
	var number int64
	var like string
	if ecosystem != 0 {
		like = "%," + strconv.FormatInt(ecosystem, 10)
	}

	req := GetDB(nil).Table(rk.TableName())
	if endId == 0 {
		if startId == 0 {
			return number
		}
		if err := req.Where("table_name = '1_keys' AND data = '' AND block_id >= ? AND table_id like ?", startId, like).Count(&number).Error; err != nil {
			log.WithFields(log.Fields{"warn": err}).Warn("get New Key Number err")
			return 0
		}
	} else {
		if startId == endId {
			if err := req.Where("table_name = '1_keys' AND data = '' AND block_id >= ? AND block_id <= ? AND table_id like ?", startId, endId, like).Count(&number).Error; err != nil {
				log.WithFields(log.Fields{"warn": err}).Warn("get New Key Number err")
				return 0
			}
		} else {
			if err := req.Where("table_name = '1_keys' AND data = '' AND block_id >= ? AND block_id < ? AND table_id like ?", startId, endId, like).Count(&number).Error; err != nil {
				log.WithFields(log.Fields{"warn": err}).Warn("get New Key Number err")
				return 0
			}
		}
	}
	return number
}

func getTimeLineEcoListInfo(stTime int64, endTime int64, his []History) int64 {
	var ecoNum int64
	for i := 0; i < len(his); i++ {
		if his[i].Createdat >= stTime && his[i].Createdat < endTime {
			ecoNum += 1
		}
	}
	return ecoNum

}

func GetHonorListToRedis(cmd string) {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	var (
		bk   Block
		list []any
	)
	rets := &GeneralResponse{}
	switch cmd {
	case "newest":
		bkList, err := bk.GetBlocksFrom(1, 10, "desc")
		if err != nil {
			log.WithFields(log.Fields{"INFO": err}).Info("Get Honor List To Redis Get Blocks Failed")
			return
		}

		for i := 0; i < len(bkList); i++ {
			replyRate, err := GetNodeBlockReplyRate(&bkList[i])
			if err != nil {
				log.WithFields(log.Fields{"INFO": err, "block id": bkList[i].ID}).Info("Get Honor List To Redis Get Reply Rate Failed")
				return
			}
			for _, value := range HonorNodes {
				if value.NodePosition == bkList[i].NodePosition && value.ConsensusMode == bkList[i].ConsensusMode {
					value.ReplyRate = replyRate
					list = append(list, value)
				}
			}
		}
		rets.Total = int64(len(list))
		rets.List = list
		rets.Page = 1
		rets.Limit = 10
	case "pkg_rate":
		sort.Sort(LeaderboardSlice(HonorNodes))
		offset := 0
		if len(HonorNodes) >= offset {
			data := HonorNodes[offset:]
			var list []storage.HonorNodeModel
			for _, val := range data {
				if !val.PkgAccountedFor.IsZero() {
					list = append(list, val)
				}
			}
			if len(list) >= 5 {
				list = list[:5]
			}
			rets.Total = int64(len(data))
			rets.List = list
		} else {
			rets.Total = 0
			rets.List = nil
		}
		rets.Page = 1
		rets.Limit = 5
	default:
		log.WithFields(log.Fields{"INFO": cmd}).Info("get honor node list to redis unknown cmd")
		return
	}

	value, err := json.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("Get Honor List To Redis Json Marshal Failed")
		return
	}

	rd := RedisParams{
		Key:   "dashboard-node-" + cmd,
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("Get Honor List To Redis Failed")
		return
	}

	return
}

func GetHonorListFromRedis(cmd string) (*GeneralResponse, error) {
	rets := &GeneralResponse{}

	switch cmd {
	case "newest", "pkg_rate":

	default:
		return nil, errors.New("get honor node list from redis unknown cmd")
	}
	rd := RedisParams{
		Key:   "dashboard-node-" + cmd,
		Value: "",
	}
	if err := rd.Get(); err != nil {
		return nil, err
	}
	err := json.Unmarshal([]byte(rd.Value), rets)
	if err != nil {
		return nil, err
	}
	return rets, err
}

func InitHonorNodeByRedis(cmd string) {
	rd := RedisParams{
		Key:   "dashboard-node-" + cmd,
		Value: "",
	}
	err := rd.Del()
	if err != nil {
		log.WithFields(log.Fields{"INFO": err, "cmd": cmd}).Info("Init Node List Redis Failed")
	}
}

func GetHonorNodeMapToRedis() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	rets, err := GetHonorNodeMap()
	value, err := json.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("Get Honor Node Map To Redis Json Marshal Failed")
		return
	}

	rd := RedisParams{
		Key:   "dashboard-node-" + "map",
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("Get Honor Node Map To Redis Failed")
		return
	}
}

func GetHonorNodeMapFromRedis() (*HonorNodeMapResponse, error) {
	rets := &HonorNodeMapResponse{}
	rd := RedisParams{
		Key:   "dashboard-node-" + "map",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		return nil, err
	}
	err := json.Unmarshal([]byte(rd.Value), rets)
	if err != nil {
		return nil, err
	}
	return rets, err
}

func SendAllWebsocketData() {
	var scanOut ScanOut
	ret1, err := scanOut.GetDashboardFromRedis()
	if err != nil {
		log.Info("Get Dashboard Redis Failed:", err.Error())
	} else {
		err := SendDashboardDataToWebsocket(ret1, ChannelStatistical)
		if err != nil {
			log.Info("Send Websocket Failed:", err.Error(), "cmd:", ChannelStatistical)
		}
	}

	ret2, err := GetTraninfoFromRedis(30)
	if err != nil {
		log.Info("Get Tran info From Redis Failed:", err.Error())
	} else {
		err = SendTpsListToWebsocket(ret2)
		if err != nil {
			log.Info("Send Websocket Failed:", err.Error(), "cmd:", ChannelBlockTpsList)
		}
	}

	ret3, _, err := GetTransactionBlockFromRedis()
	if err != nil {
		log.Info("Get Transaction Block From Redis Failed:", err.Error())
	} else {
		err = SendTransactionListToWebsocket(ret3)
		if err != nil {
			log.Info("Send Websocket Failed:", err.Error(), "cmd:", ChannelBlockTransactionList)
		}
	}

	ret4, err := GetBlockListFromRedis()
	if err != nil {
		log.Info("Get Transaction Block From Redis Failed:", err.Error())
	} else {
		if err := SendBlockListToWebsocket(&ret4.List); err != nil {
			log.WithFields(log.Fields{"warn": err}).Warn("sendBlockListToWebsocket err")
			log.Info("Send Websocket Failed:", err.Error(), "cmd:", ChannelBlockList)
		}
	}

	ret5, err := GetHonorNodeMapFromRedis()
	if err != nil {
		log.Info("Get Honor Node Map From Redis Failed:", err.Error())
	} else {
		err = SendDashboardDataToWebsocket(ret5, ChannelNodeMap)
		if err != nil {
			log.WithFields(log.Fields{"INFO": err, " channel": ChannelNodeMap}).Info("Send Websocket Failed")
		}
	}

	cmd := "pkg_rate"
	ret6, err := GetHonorListFromRedis(cmd)
	if err != nil {
		log.Info("Get Honor List From Redis Failed:", err.Error(), "cmd", cmd)
	} else {
		err = SendDashboardDataToWebsocket(ret6.List, ChannelNodePkgRate)
		if err != nil {
			log.WithFields(log.Fields{"INFO": err, " cmd": ChannelNodePkgRate}).Info("Send Websocket Failed")
			return
		}
	}

	cmd = "newest"
	ret7, err := GetHonorListFromRedis(cmd)
	if err != nil {
		log.Info("Get Honor List From Redis Failed:", err.Error(), "cmd", cmd)
	} else {
		err = SendDashboardDataToWebsocket(ret7.List, ChannelNodeNewest)
		if err != nil {
			log.WithFields(log.Fields{"INFO": err, " cmd": ChannelNodeNewest}).Info("Send Websocket Failed")
			return
		}
	}
}
