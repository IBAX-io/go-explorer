/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	//history
	FifteenDaysGasFee          = "15_days_gas_fee_chart"
	FifteenDaysNewCirculations = "15_days_new_circulations_chart"
	FifteenDaysBlockSize       = "15_days_block_size_chart"
	FifteenDaysTransaction     = "15_days_transaction_chart"
	FifteenDaysBlockNumber     = "15_days_block_number_chart"
	NewNfTMinerChange          = "new_nft_miner_change_chart"
	NftMinerRewardChange       = "nft_miner_reward_change_chart"
	TopTenEcosystemTx          = "top_ten_ecosystem_tx_chart"
	NodeVoteChange             = "node_vote_change_chart"
	NodeStakingChange          = "node_staking_change_chart"

	//realTime
	NewKey           = "new_key_chart"
	AccountChange    = "account_change_chart"
	NftMinerInterval = "nft_miner_interval_chart"
	NodeRegion       = "node_region_chart"
)

func InsertRedis(name string, data string) {
	if data == "" {
		return
	}
	rd := RedisParams{
		Key:   name,
		Value: data,
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err, "name": name}).Error("Insert Redis Failed")
		return
	}
}

func GetRedis(name string) (any, error) {
	var (
		rets any
	)
	rd := RedisParams{
		Key: name,
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err, "name": name}).Warn("Get Redis Failed")
		return nil, err
	}
	err := json.Unmarshal([]byte(rd.Value), &rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "name": name}).Warn("Get Redis json Unmarshal Failed")
		return nil, err
	}
	return rets, nil
}

func GetDataChart(date any, err error) string {
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Get Data Chart Failed")
		return ""
	}
	value, err := json.Marshal(date)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Get Data Chart Data json marshal Failed")
		return ""
	}
	return string(value)
}

func DataChartHistoryServer() {
	InsertRedis(FifteenDaysGasFee, GetDataChart(Get15DayGasFeeChart()))
	InsertRedis(FifteenDaysNewCirculations, GetDataChart(Get15DayNewCirculationsChart()))
	InsertRedis(FifteenDaysBlockSize, GetDataChart(Get15DayBlockSizeChart()))
	InsertRedis(FifteenDaysTransaction, GetDataChart(GetEco15DayTransactionChart(1)))
	InsertRedis(FifteenDaysBlockNumber, GetDataChart(Get15DayBlockNumberChart()))
	InsertRedis(NewNfTMinerChange, GetDataChart(GetNewNftMinerChangeChart()))
	InsertRedis(NftMinerRewardChange, GetDataChart(GetNftMinerRewardChangeChart()))
	InsertRedis(TopTenEcosystemTx, GetDataChart(GetTopTenEcosystemTxChart()))

	//node chart
	InsertRedis(NodeVoteChange, GetDataChart(getNodeVoteChangeChangeChart()))
	InsertRedis(NodeStakingChange, GetDataChart(getNodeStakingChangeChangeChart()))

}

func DataChartRealtimeSever() {
	InsertRedis(NewKey, GetDataChart(GetNewKeysChart()))
	InsertRedis(AccountChange, GetDataChart(GetAccountChangeChart()))
	InsertRedis(NftMinerInterval, GetDataChart(GetNftMinerIntervalChart()))

	//node chart
	InsertRedis(DaoVoteChart, GetDataChart(getDaoVoteChart()))
	InsertRedis(NodeRegion, GetDataChart(getTopTenNodeRegionChart()))
}

func GetDayNumberFormat(days int64) (string, string) {
	timeDbFormat := "yyyy-MM-dd" //days
	layout := "2006-01-02"
	if days >= 720 { //years
		timeDbFormat = "yyyy"
		layout = "2006"
	} else if days >= 60 { //months
		timeDbFormat = "yyyy-MM"
		layout = "2006-01"
	}
	return timeDbFormat, layout
}

func Get15DayGasFeeChart() (GasFeeChangeResponse, error) {
	var rets GasFeeChangeResponse
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	rets.Time = make([]int64, getDays)
	rets.GasFee = make([]string, getDays)

	var list []DaysAmount
	err := GetDB(nil).Raw(`
SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS days,sum(amount) AS amount
FROM "1_history" WHERE (type = 1 OR type = 2) AND ecosystem = 1 AND created_at >= ? GROUP BY days ORDER BY days DESC LIMIT ?
`, t1.UnixMilli(), getDays).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get 15Day Gas Fee Chart Failed")
		return rets, err
	}

	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.GasFee[i] = GetDaysAmount(rets.Time[i], list)
	}

	return rets, nil
}

func GetHonorNodeChart(page, limit int) (HonorNodeChartResponse, error) {
	var rets HonorNodeChartResponse
	type nodeGasFee struct {
		NodePosition int64           `json:"node_position" gorm:"column:node_position"`
		GasFee       decimal.Decimal `json:"gas_fee" gorm:"column:gas_fee"`
	}
	var list []nodeGasFee
	err := GetDB(nil).Raw(`SELECT node_position,
(SELECT sum(amount) AS gas_fee FROM "1_history" WHERE (type = 1 OR type = 2) AND  ecosystem = 1 AND recipient_id = bk.key_id)
FROM block_chain AS bk GROUP BY node_position,key_id`).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get 15Day Gas Fee Chart Failed")
		return rets, err
	}
	rets.Total = int64(len(HonorNodes))
	rets.Page = page
	rets.Limit = limit

	sort.Sort(LeaderboardSlice(HonorNodes))
	offset := (page - 1) * limit
	if len(HonorNodes) >= offset {
		data := HonorNodes[offset:]
		if len(data) >= limit {
			data = data[:limit]
		}
		rets.List = make([]HonorNodeListResponse, len(data))
		for i := 0; i < len(data); i++ {
			rets.List[i].NodePosition = data[i].NodePosition
			rets.List[i].KeyID = data[i].KeyID
			rets.List[i].NodeName = data[i].NodeName
			rets.List[i].City = data[i].City
			rets.List[i].IconUrl = data[i].IconUrl
			rets.List[i].NodeBlocks = data[i].NodeBlock
			rets.List[i].PkgAccountedFor = data[i].PkgAccountedFor
		}
		//rets.Data = data
		//ret.Return(rets, CodeSuccess)
	}

	for i := 0; i < len(list); i++ {
		for key, value := range rets.List {
			if value.NodePosition == list[i].NodePosition+1 {
				rets.List[key].GasFee = list[i].GasFee.String()
			}
		}
	}
	if len(HonorNodes) >= 10 {
		rets.NodeBlock = make([]int64, 10)
		rets.Name = make([]string, 10)
	} else {
		rets.NodeBlock = make([]int64, len(HonorNodes))
		rets.Name = make([]string, len(HonorNodes))
	}
	for i, value := range HonorNodes {
		if i >= len(rets.Name) {
			break
		}
		rets.Name[i] = value.City
		rets.NodeBlock[i] = value.NodeBlock
	}

	return rets, nil
}

func Get15DayNewCirculationsChart() (CirculationsChartResponse, error) {
	var (
		rets CirculationsChartResponse
		list []DaysAmount
	)

	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	rets.Change.Time = make([]int64, getDays)
	rets.Change.FreezeAmount = make([]string, getDays)
	rets.Change.Circulations = make([]string, getDays)
	totalCir := decimal.New(0, 0)
	totalFreeze := decimal.New(0, 0)

	err := GetDB(nil).Raw(`
SELECT to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd') AS days,sum(amount) AS amount 
FROM "1_history" WHERE ecosystem = 1 AND type = 12 GROUP BY days ORDER BY days DESC limit ?
`, getDays).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get 15 Day New Circulations Chart Failed")
		return rets, err
	}

	for i := 0; i < len(rets.Change.Time); i++ {
		rets.Change.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		cir, _ := decimal.NewFromString(GetDaysAmount(rets.Change.Time[i], list))
		freeze, _ := decimal.NewFromString(rets.Change.FreezeAmount[i])
		rets.Change.Circulations[i] = cir.Add(freeze).String()
		totalCir = totalCir.Add(cir)
		totalFreeze = totalFreeze.Add(freeze)
	}
	rets.Circulations = totalCir.String()
	rets.FreezeAmount = totalFreeze.String()
	rets.TotalCirculations = totalCir.Add(totalFreeze).String()

	return rets, nil

}

func GetAccountChangeChart() (AccountChangeChartResponse, error) {
	var (
		rets     AccountChangeChartResponse
		htList   []DaysNumber
		acList   []DaysNumber
		minTime  int64
		maxTime  int64
		nowTotal int64
	)
	err := GetDB(nil).Raw(`
WITH rollback_tx AS(
	SELECT to_char(to_timestamp(log.time), 'yyyy-mm-dd') AS days,count(1) num
	FROM (SELECT tx_hash,table_id FROM rollback_tx WHERE table_name = '1_keys' AND table_id like '%,1' AND data = '') AS rb LEFT JOIN (
		SELECT timestamp/1000 as time,hash FROM log_transactions 
	)AS log ON (log.hash = rb.tx_hash) GROUP BY days ORDER BY days ASC
)
SELECT rk1.days,
(SELECT SUM(num)+3 FROM rollback_tx s2 WHERE s2.days <= rk1.days AND SUBSTRING(rk1.days,0,5) = SUBSTRING(s2.days,0,5)) as num
FROM rollback_tx AS rk1
`).Find(&acList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Account Change Chart acList Failed")
		return rets, nil
	}

	err = GetDB(nil).Raw(`
with "1_history" AS(
	SELECT h3.days,count(1) num FROM(
		SELECT min(h2.days) days,h2.keyid FROM(
				SELECT h1.days,h1.keyid FROM (
					SELECT sender_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd') AS days FROM "1_history" WHERE sender_id <> 0 AND sender_balance > 0 AND ecosystem = 1 GROUP BY sender_id,days
					 UNION 
					SELECT recipient_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd') AS days FROM "1_history" WHERE recipient_id <> 0 AND recipient_balance > 0 AND ecosystem = 1 GROUP BY recipient_id,days
				) AS h1 ORDER BY h1.days DESC
		)as h2
		GROUP BY h2.keyid ORDER BY days ASC
	)AS h3 GROUP BY days ORDER BY days ASC
)
SELECT h4.days,
(SELECT sum(num) FROM "1_history" h5 WHERE h5.days <= h4.days AND SUBSTRING(h4.days,0,5) = SUBSTRING(h5.days,0,5)) as num
FROM "1_history" AS h4
`).Find(&htList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Account Change Chart htList Failed")
		return rets, nil
	}
	maxAccountedFor := decimal.New(0, 0)
	minAccountedFor := decimal.New(100, 0)

	var startTime time.Time
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	getDaysCount := func(findTime int64, list []DaysNumber) int64 {
		return GetDaysNumber(findTime, list)
	}
	if len(acList) > 0 {
		startTime, err = time.ParseInLocation("2006-01-02", acList[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get NewKeys Chart Block Id List ParseInLocation Failed")
			return rets, err
		}
		nowTotal = acList[len(acList)-1].Num
		var lastCount int64
		var lastHasToken int64

		for startTime.Unix() <= today.Unix() {
			rets.Time = append(rets.Time, startTime.Unix())
			count := getDaysCount(startTime.Unix(), acList)
			if count > 0 {
				lastCount = count
				rets.Total = append(rets.Total, count)
			} else {
				rets.Total = append(rets.Total, lastCount)
				count = lastCount
			}
			hasToken := getDaysCount(startTime.Unix(), htList)
			if hasToken > 0 {
				lastHasToken = hasToken
				rets.HasToken = append(rets.HasToken, hasToken)
			} else {
				rets.HasToken = append(rets.HasToken, lastHasToken)
				hasToken = lastHasToken
			}
			accountedFor := decimal.New(0, 0)

			if hasToken > 0 {
				if count > 0 {
					accountedFor = decimal.NewFromInt(hasToken*100).DivRound(decimal.NewFromInt(count), 2)
				}
			}
			rets.AccountedFor = append(rets.AccountedFor, accountedFor)
			if accountedFor.LessThan(minAccountedFor) {
				minAccountedFor = accountedFor
				minTime = startTime.Unix()
			}
			if accountedFor.GreaterThan(maxAccountedFor) {
				maxAccountedFor = accountedFor
				maxTime = startTime.Unix()
			}

			startTime = startTime.AddDate(0, 0, 1)
		}
	}

	rets.NowTotal = nowTotal
	rets.MinAccountedFor = minAccountedFor
	rets.MinTime = minTime
	rets.MaxAccountedFor = maxAccountedFor
	rets.MaxTime = maxTime

	return rets, nil

}

func GetTopTenStakingAccount() ([]StakingAccountResponse, error) {
	var (
		rets         []StakingAccountResponse
		totalStaking SumAmount
		key          Key
	)

	err := GetDB(nil).Raw(`
SELECT sum(to_number(coalesce(NULLIF(lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(lock->>'candidate_substitute',''),'0'),'999999999999999999999999')) FROM "1_keys" WHERE ecosystem = 1
`).Take(&totalStaking.Sum).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Top Ten Staking Account Total Staking Amount Failed")
		return nil, err
	}

	err = GetDB(nil).Table(key.TableName()).
		Select(`account,to_number(coalesce(NULLIF(lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(lock->>'candidate_substitute',''),'0'),'999999999999999999999999') AS stake_amount`).
		Where("ecosystem = 1").Order("stake_amount desc").Limit(10).Find(&rets).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Top Ten Staking Account list Failed")
		return nil, err
	}
	zero := decimal.New(0, 0)
	if totalStaking.Sum.GreaterThan(zero) {
		for i := 0; i < len(rets); i++ {
			if rets[i].StakeAmount.GreaterThan(zero) {
				rets[i].AccountedFor = rets[i].StakeAmount.Mul(decimal.NewFromInt(100)).DivRound(totalStaking.Sum, 2)
			}
		}
	}

	return rets, nil
}

func Get15DayBlockSizeChart() (StorageCapacitysChart, error) {
	var (
		rets StorageCapacitysChart
		list []DaysNumber
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	err := GetDB(nil).Raw(`SELECT to_char(to_timestamp(time),'yyyy-MM-dd') days,
	sum(length("data")) num FROM block_chain WHERE time >= ? GROUP BY days`, t1.Unix()).Find(&list).Error
	if err != nil {
		return rets, err
	}
	rets.Time = make([]int64, getDays)
	rets.StorageCapacitys = make([]string, getDays)
	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.StorageCapacitys[i] = ToCapcityMb(GetDaysNumber(rets.Time[i], list))
	}
	return rets, nil
}

func Get15DayBlockSizeList(page, limit int) (GeneralResponse, error) {
	var (
		list []BlockSizeListResponse
		rets GeneralResponse
		bk   Block
	)
	rets.Page = page
	rets.Limit = limit
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays+1)

	err := GetDB(nil).Table(bk.TableName()).
		Where("time >= ? AND time < ?", t1.Unix(), yesterday.AddDate(0, 0, 1).Unix()).Count(&rets.Total).Error
	if err != nil {
		return rets, err
	}

	err = GetDB(nil).Table(bk.TableName()).Select("id,time,length(data) size,tx").
		Where("time >= ? AND time < ?", t1.Unix(), yesterday.AddDate(0, 0, 1).Unix()).
		Offset((page - 1) * limit).Limit(limit).Order("id desc").Find(&list).Error
	if err != nil {
		return rets, err
	}
	type blockSizeList struct {
		Id   int64  `json:"id"`
		Time int64  `json:"time"`
		Size string `json:"size"`
		Tx   int64  `json:"tx"`
	}
	var cut []blockSizeList
	for i := 0; i < len(list); i++ {
		var rts blockSizeList
		rts.Id = list[i].Id
		rts.Tx = list[i].Tx
		rts.Size = TocapacityString(list[i].Size)
		rts.Time = list[i].Time
		cut = append(cut, rts)
	}
	rets.List = cut

	return rets, nil
}

func Get15DayTxList(page, limit int) (GeneralResponse, error) {
	var (
		list []TransactionListResponse
		rets GeneralResponse
		lt   LogTransaction
	)
	rets.Page = page
	rets.Limit = limit
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays+1)

	err := GetDB(nil).Table(lt.TableName()).
		Where("timestamp >= ? AND timestamp < ?", t1.UnixMilli(), yesterday.AddDate(0, 0, 1).UnixMilli()).Count(&rets.Total).Error
	if err != nil {
		return rets, err
	}

	err = GetDB(nil).Raw(`
SELECT log.hash,log.block,log.timestamp as time,log.address,(SELECT name AS name FROM "1_ecosystems" as es WHERE es.id = log.ecosystem_id) 
FROM (
	SELECT encode(hash,'hex') hash,block,timestamp,address,ecosystem_id FROM 
	log_transactions WHERE timestamp >= ? AND timestamp < ?
)as log ORDER BY log.block DESC OFFSET ? LIMIT ?
`, t1.UnixMilli(), yesterday.AddDate(0, 0, 1).UnixMilli(), (page-1)*limit, limit).Find(&list).Error
	if err != nil {
		return rets, err
	}
	type transactionList struct {
		Hash    string `json:"hash"`
		Block   int64  `json:"block"`
		Time    int64  `json:"time"`
		Address string `json:"address"`
		Name    string `json:"name"`
	}
	var cut []transactionList
	for i := 0; i < len(list); i++ {
		var rts transactionList
		rts.Hash = list[i].Hash
		rts.Block = list[i].Block
		rts.Time = MsToSeconds(list[i].Time)
		rts.Address = converter.AddressToString(list[i].Address)
		rts.Name = list[i].Name
		cut = append(cut, rts)
	}
	rets.List = cut

	return rets, nil
}

func GetNewKeysChart() (NewKeyHistoryChart, error) {
	var (
		keyChart    NewKeyHistoryChart
		idList      []int64
		bk          Block
		blockIdList []DaysNumber
		total       CountInt64
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	err := GetDB(nil).
		Raw("SELECT count(1) from (SELECT to_char(to_timestamp(time),'yyyy-MM-dd') AS days FROM block_chain GROUP BY days)as bk").Take(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get NewKeys Chart Total Failed")
		return keyChart, err
	}
	timeDbFormat, layout := GetDayNumberFormat(total.Count)

	err = GetDB(nil).Table(bk.TableName()).Select(fmt.Sprintf("to_char(to_timestamp(time),'%s') AS days,min(id) num", timeDbFormat)).
		Group("days").Order("days asc").Find(&blockIdList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get NewKeys Chart Block Id List Failed")
		return keyChart, err
	}
	var startTime time.Time
	if len(blockIdList) >= 1 {
		startTime, err = time.ParseInLocation(layout, blockIdList[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get NewKeys Chart Block Id List ParseInLocation Failed")
			return keyChart, err
		}
		for startTime.Unix() <= today.Unix() {
			keyChart.Time = append(keyChart.Time, startTime.Format(layout))
			idList = append(idList, GetDaysNumberLike(startTime.Unix(), blockIdList, false, "asc"))
			startTime = addTimeFromLayout(layout, startTime)
		}
	}

	_, err = bk.GetByTimeBlockId(today.Unix())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("get Ecosystem New Key Chart Today Block Failed")
		return keyChart, err
	}

	for i := 0; i < len(idList); i++ {
		if i == len(idList)-1 {
			keyChart.NewKey = append(keyChart.NewKey, getNewKeyNumber(idList[i], bk.ID, 1))
		} else {
			keyChart.NewKey = append(keyChart.NewKey, getNewKeyNumber(idList[i], idList[i+1], 1))
		}
	}

	return keyChart, nil
}

//Get15DayBlockNumberChart response time:200-1500ms TODO:NEED TO Redis
func Get15DayBlockNumberChart() (DaysNumberResponse, error) {
	var rets DaysNumberResponse
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	rets.Time = make([]int64, getDays)
	rets.Number = make([]int64, getDays)

	sql := `
SELECT to_char(to_timestamp(time),'yyyy-MM-dd') AS days,count(1) num FROM block_chain
WHERE time >= ?
GROUP BY days ORDER BY days DESC LIMIT ?
`
	list, err := FindDaysNumber(sql, t1.Unix(), getDays)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get 15 Day Block Number Chart Failed")
		return rets, err
	}

	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.Number[i] = GetDaysNumber(rets.Time[i], list)
	}

	return rets, nil
}

func GetNftMinerIntervalChart() (NftMinerIntervalResponse, error) {
	var (
		rets NftMinerIntervalResponse
	)
	if !NftMinerReady {
		return rets, nil
	}
	sql := `
SELECT 
	sum(
		case when energy_point between 1 and 11 then 1 else 0 end
  ) as one_to_ten,
	sum(
		case when energy_point between 11 and 21 then 1 else 0 end
  ) as ten_to_twenty,
	sum(
		case when energy_point between 21 and 31 then 1 else 0 end
  ) as twenty_to_thirty,
	sum(
		case when energy_point between 31 and 41 then 1 else 0 end
  ) as thirty_to_forty,
	sum(
		case when energy_point between 41 and 51 then 1 else 0 end
  ) as forty_to_fifty,
	sum(
		case when energy_point between 51 and 61 then 1 else 0 end
  ) as fifty_to_sixty,
	sum(
		case when energy_point between 61 and 71 then 1 else 0 end
  ) as sixty_to_seventy,
	sum(
		case when energy_point between 71 and 81 then 1 else 0 end
  ) as seventy_to_eighty,
	sum(
		case when energy_point between 81 and 91 then 1 else 0 end
  ) as eighty_to_ninety,
	sum(
		case when energy_point between 91 and 101 then 1 else 0 end
  ) as ninety_to_hundred

FROM "1_nft_miner_items" where merge_status = 1
`
	err := GetDB(nil).Raw(sql).Take(&rets).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Energy Point Interval Chart Failed")
		return rets, err
	}
	return rets, nil
}

func GetNftMinerIntervalListChart() ([]NftMinerIntervalResponse, error) {
	var (
		rets  []NftMinerIntervalResponse
		total int64
	)
	if !NftMinerReady {
		return rets, nil
	}

	err := GetDB(nil).Raw(`
SELECT count(1) FROM(
	SELECT to_char(to_timestamp(date_created),'yyyy-MM-dd') days
	FROM "1_nft_miner_items" where merge_status = 1 GROUP BY days
)AS ns`).Take(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Energy Point Interval List Chart Total Failed")
		return rets, err
	}
	timeDbFormat, _ := GetDayNumberFormat(total)

	sql := fmt.Sprintf(`
 WITH "1_nft_miner_items" AS(
	SELECT to_char(to_timestamp(date_created),'%s') AS time,
	sum(
		case when energy_point between 1 and 11 then 1 else 0 end
	) as one_to_ten,
	sum(
		case when energy_point between 11 and 21 then 1 else 0 end
	) as ten_to_twenty,
	sum(
		case when energy_point between 21 and 31 then 1 else 0 end
	) as twenty_to_thirty,
	sum(
		case when energy_point between 31 and 41 then 1 else 0 end
	) as thirty_to_forty,
	sum(
		case when energy_point between 41 and 51 then 1 else 0 end
	) as forty_to_fifty,
	sum(
		case when energy_point between 51 and 61 then 1 else 0 end
	) as fifty_to_sixty,
	sum(
		case when energy_point between 61 and 71 then 1 else 0 end
	) as sixty_to_seventy,
	sum(
		case when energy_point between 71 and 81 then 1 else 0 end
	) as seventy_to_eighty,
	sum(
		case when energy_point between 81 and 91 then 1 else 0 end
	) as eighty_to_ninety,
	sum(
		case when energy_point between 91 and 101 then 1 else 0 end
	) as ninety_to_hundred
	FROM "1_nft_miner_items" WHERE merge_status = 1 GROUP BY time ORDER BY time ASC
)
	SELECT nf.time,
	(SELECT sum(one_to_ten) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS one_to_ten,
	(SELECT sum(ten_to_twenty) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS ten_to_twenty,
	(SELECT sum(twenty_to_thirty) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS twenty_to_thirty,
	(SELECT sum(thirty_to_forty) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS thirty_to_forty,
	(SELECT sum(forty_to_fifty) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS forty_to_fifty,
	(SELECT sum(fifty_to_sixty) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS fifty_to_sixty,
	(SELECT sum(sixty_to_seventy) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS sixty_to_seventy,
	(SELECT sum(seventy_to_eighty) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS seventy_to_eighty,
	(SELECT sum(eighty_to_ninety) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS eighty_to_ninety,
	(SELECT sum(ninety_to_hundred) FROM "1_nft_miner_items" ns WHERE ns.time <= nf.time AND SUBSTRING(nf.time,0,5) = SUBSTRING(ns.time,0,5))AS ninety_to_hundred
FROM "1_nft_miner_items" AS nf
`, timeDbFormat)
	err = GetDB(nil).Raw(sql).Find(&rets).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Energy Point Interval List Chart Failed")
		return rets, err
	}

	return rets, nil
}

func GetNftEnergyPowerChangeChart() (NftMinerEnergyPowerChangeResponse, error) {
	var (
		rets NftMinerEnergyPowerChangeResponse
	)
	rets.Time = make([]string, 0)
	rets.EnergyPower = make([]string, 0)
	if !NftMinerReady {
		return rets, nil
	}
	type nftStakEnergyPower struct {
		Days        string
		DelPower    decimal.Decimal
		EnergyPower decimal.Decimal
	}
	layout := "2006-01-02"
	tz := time.Now()
	endTime := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	sql := fmt.Sprintf(`
SELECT case WHEN s1.days > '' THEN
	s1.days
ELSE
	s2.days
END,coalesce(s1.energy_power,0) energy_power,coalesce(s2.energy_power,0) del_power
FROM (
		SELECT to_char(to_timestamp(start_dated), 'yyyy-mm-dd') AS days,sum(energy_power) energy_power
		FROM "1_nft_miner_staking" GROUP BY days ORDER BY days asc
)AS s1

FULL JOIN(
		SELECT to_char(to_timestamp(end_dated),'yyyy-MM-dd') AS days,
		sum(energy_power) energy_power FROM "1_nft_miner_staking" WHERE end_dated <= %d GROUP BY days ORDER BY days asc
)AS s2 ON(s2.days=s1.days)
`, endTime.AddDate(0, 0, 1).Unix())
	var list []nftStakEnergyPower
	err := GetDB(nil).Raw(sql).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Energy Power Change Chart Failed")
		return rets, err
	}
	if len(list) > 0 {
		var lastEnergyPower decimal.Decimal

		startTime, err := time.ParseInLocation(layout, list[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Energy Power Change Chart ParseInLocation Failed")
			return rets, err
		}
		for startTime.Unix() <= endTime.Unix() {
			findTime := startTime.Format(layout)
			rets.Time = append(rets.Time, findTime)
			for _, value := range list {
				if findTime == value.Days {
					if value.EnergyPower.String() != "0" {
						lastEnergyPower = lastEnergyPower.Add(value.EnergyPower)
					}
					if value.DelPower.String() != "0" {
						lastEnergyPower = lastEnergyPower.Sub(value.DelPower)
					}
				}
			}
			rets.EnergyPower = append(rets.EnergyPower, lastEnergyPower.String())
			startTime = startTime.AddDate(0, 0, 1)
		}
	}

	return rets, nil
}

func GetNftMinerStakedChangeChart() (NftMinerStakingChangeResponse, error) {
	var (
		rets  NftMinerStakingChangeResponse
		total int64
	)
	rets.Time = make([]string, 0)
	rets.StakeAmount = make([]int64, 0)
	rets.Number = make([]int64, 0)
	if !NftMinerReady {
		return rets, nil
	}
	type stakedChange struct {
		Days        string `json:"days"`
		StakeAmount int64  `json:"stakeAmount"`
		Number      int64  `json:"number"`
	}

	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	err := GetDB(nil).Raw(`
SELECT count(1) FROM(
	SELECT to_char(to_timestamp(start_dated),'yyyy-MM-dd') AS days FROM "1_nft_miner_staking" GROUP BY days
)AS t1`).Take(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Staked Change Days Total Failed")
		return rets, err
	}
	timeDbFormat, layout := GetDayNumberFormat(total)

	sql := fmt.Sprintf(`
SELECT to_char(to_timestamp(start_dated),'%s') AS days,
sum(stake_amount) stake_amount,count(1) number FROM "1_nft_miner_staking" GROUP BY days ORDER BY days asc
`, timeDbFormat)
	var list []stakedChange
	err = GetDB(nil).Raw(sql).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Staked Change Chart Failed")
		return rets, err
	}
	var startTime time.Time

	getNumber := func(getTime int64, list []stakedChange) (int64, int64) {
		for i := 0; i < len(list); i++ {
			times, _ := time.ParseInLocation(layout, list[i].Days, time.Local)
			if getTime == times.Unix() {
				return list[i].StakeAmount, list[i].Number
			}
		}
		return 0, 0
	}
	if len(list) > 0 {
		startTime, err = time.ParseInLocation(layout, list[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Nft Miner Staked Change Chart ParseInLocation Failed")
			return rets, err
		}
		for startTime.Unix() < today.Unix() {
			rets.Time = append(rets.Time, startTime.Format(layout))
			stakeAmount, number := getNumber(startTime.Unix(), list)
			rets.StakeAmount = append(rets.StakeAmount, stakeAmount)
			rets.Number = append(rets.Number, number)

			startTime = addTimeFromLayout(layout, startTime)
		}
	}

	return rets, nil
}

func addTimeFromLayout(layout string, startTime time.Time) time.Time {
	switch layout {
	case "2006-01-02":
		startTime = startTime.AddDate(0, 0, 1)
	case "2006-01":
		startTime = startTime.AddDate(0, 1, 0)
	case "2006":
		startTime = startTime.AddDate(1, 0, 0)
	}
	return startTime
}

func GetHistoryNewEcosystemChangeChart() (DaysNumberResponse, error) {
	var (
		rets DaysNumberResponse
		bk   Block
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	rets.Time = make([]int64, 0)
	rets.Number = make([]int64, 0)

	sql := `
SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS days,count(1)num FROM "1_history" 
WHERE comment = 'taxes for execution of @1NewEcosystem contract' AND type = 1 GROUP BY days ORDER BY days ASC
`
	firstSys, err := bk.GetSystemTime()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get History New Ecosystem Change Chart Sys Time Failed")
		return rets, err
	}

	list, err := FindDaysNumber(sql)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get History New Ecosystem Change Chart Failed")
		return rets, err
	}

	getNumber := func(getTime int64, list []DaysNumber) int64 {
		for i := 0; i < len(list); i++ {
			times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
			if getTime == times.Unix() {
				return list[i].Num
			}
		}
		return 0
	}
	if len(list) > 0 {
		ft := time.Unix(firstSys, 0)
		startTime := time.Date(ft.Year(), ft.Month(), ft.Day(), 0, 0, 0, 0, time.Local)
		var isFirst = true
		for startTime.Unix() <= today.Unix() {
			rets.Time = append(rets.Time, startTime.Unix())
			if isFirst {
				rets.Number = append(rets.Number, getNumber(startTime.Unix(), list)+1)
				isFirst = false
			} else {
				rets.Number = append(rets.Number, getNumber(startTime.Unix(), list))
			}
			startTime = startTime.AddDate(0, 0, 1)
		}
	}

	return rets, nil
}

func GetTokenEcosystemRatioChart() (TokenEcosystemResponse, error) {
	var (
		rets  TokenEcosystemResponse
		eco   Ecosystem
		total int64
	)
	q := GetDB(nil).Table(eco.TableName())
	err := q.Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Token Ecosystem Ratio Chart Total Failed")
		return rets, err
	}
	err = q.Where(`emission_amount @> '[{"type":"emission"}]'::jsonb or id = 1`).Count(&rets.Emission).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Token Ecosystem Ratio Chart Emission Failed")
		return rets, err
	}
	if rets.Emission > 0 && total > 0 {
		rets.UnEmission = total - rets.Emission
		totalDl := decimal.NewFromInt(total)
		rets.EmissionRatio = decimal.NewFromInt(rets.Emission*100).DivRound(totalDl, 2)
		rets.UnEmissionRatio = decimal.NewFromInt(rets.UnEmission*100).DivRound(totalDl, 2)
	}

	return rets, nil
}

//GetTopTenEcosystemTxChart return Fifteen days transactions TODO:200-800ms NEED TO Redis
func GetTopTenEcosystemTxChart() ([]EcosystemTxRatioResponse, error) {
	var (
		list []EcosystemTxCount
		rets []EcosystemTxRatioResponse
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)
	err := GetDB(nil).Raw(`
SELECT log.total,log.tx,log.ecosystem,log.name FROM(
		SELECT lg.tx,(
			SELECT count(1)AS total FROM "log_transactions" WHERE "timestamp" >= ?
		),lg.ecosystem,es.name,es.id FROM (
			SELECT count(1) tx,ecosystem_id AS ecosystem FROM "log_transactions" WHERE "timestamp" >= ? GROUP BY ecosystem_id
		) AS lg 
		
		INNER JOIN (
			SELECT name,id FROM "1_ecosystems"
		) AS es  ON (lg.ecosystem = es.id) 
)AS log 
ORDER BY log.tx DESC
LIMIT 10`, t1.UnixMilli(), t1.UnixMilli()).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Top Ten Ecosystem Tx Chart Failed")
		return nil, err
	}
	if len(list) >= 1 {
		type ecoinfo struct {
			SumInt64
			MinInt64
		}
		var other ecoinfo
		err = GetDB(nil).Raw(`
SELECT sum(ts.tx),min(ts.total) FROM(
	SELECT log.total,log.tx FROM(
			SELECT lg.tx,(
				SELECT count(1)AS total FROM "log_transactions" WHERE "timestamp" >= ?
			),lg.ecosystem,es.name,es.id FROM (
				SELECT count(1) tx,ecosystem_id AS ecosystem FROM "log_transactions" WHERE "timestamp" >= ? GROUP BY ecosystem_id
			) AS lg 
			
			INNER JOIN (
				SELECT name,id FROM "1_ecosystems"
			) AS es  ON (lg.ecosystem = es.id) 
			ORDER BY lg.tx desc
	)AS log 
	ORDER BY log.tx DESC OFFSET 10 
)AS ts
`, t1.UnixMilli(), t1.UnixMilli()).Take(&other).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Top Ten Ecosystem Tx Chart Other Failed")
			return nil, err
		}
		list = append(list, EcosystemTxCount{Name: "Other EcoLibs", Tx: other.Sum, Total: other.Min})
	} else {
		list = append(list, EcosystemTxCount{Name: "Other EcoLibs", Tx: 0, Total: 0})
	}
	for i := 0; i < len(list); i++ {
		var er EcosystemTxRatioResponse
		er.Name = list[i].Name
		er.Tx = list[i].Tx
		if list[i].Tx > 0 && list[i].Total > 0 {
			er.Ratio = decimal.NewFromInt(list[i].Tx*100).DivRound(decimal.NewFromInt(list[i].Total), 2)
		}
		rets = append(rets, er)
	}

	return rets, nil
}

func GetNewEcosystemChartList(page, limit int, order string) (GeneralResponse, error) {
	var (
		ret   []EcosystemListResponse
		total int64
		eco   Ecosystem
		rets  GeneralResponse
	)
	if order == "" {
		order = "id desc"
	}
	rets.Page = page
	rets.Limit = limit

	type ecosystemInfo struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Info     string `json:"info"`
		Contract int64  `json:"contract"`
		Block    int64  `json:"block"`
		Hash     []byte `json:"hash"`
	}
	var list []ecosystemInfo
	err := GetDB(nil).Table(eco.TableName()).Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get New Ecosystem Chart List Total Failed")
		return rets, err
	}
	err = GetDB(nil).Raw(`
SELECT e.id,e.name,e.info,
(SELECT count(*) from "1_contracts" AS c WHERE c.ecosystem = e.id)as contract,
CASE WHEN rk.block > 0 THEN
rk.block
ELSE
log.block
END,

CASE when length(rk.hash) > 0 THEN
rk.hash
ELSE
log.hash
END
 FROM "1_ecosystems" AS e 

LEFT JOIN(
			SELECT coalesce(block_id,0) as block,coalesce(tx_hash,'') as hash,coalesce(to_number(coalesce(NULLIF(table_id,''),'0'),'999999999999999999999999'),0)as ecoid FROM rollback_tx WHERE table_name = '1_ecosystems' AND data = ''
)AS rk ON(rk.ecoid = e.id)

LEFT JOIN(
		SELECT hash,block,block AS ecoid FROM log_transactions WHERE block = 1
)AS log ON(log.ecoid = e.id)
ORDER BY ? OFFSET ? LIMIT ? 
`, order, (page-1)*limit, limit).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get New Ecosystem Chart List Failed")
		return rets, err
	}
	escape := func(value any) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}
	for i := 0; i < len(list); i++ {
		var er EcosystemListResponse
		if list[i].Info != "" {
			minfo := make(map[string]any)
			err := json.Unmarshal([]byte(list[i].Info), &minfo)
			if err != nil {
				return rets, err
			}
			usid, ok := minfo["logo"]
			if ok {
				urid := escape(usid)
				uid, err := strconv.ParseInt(urid, 10, 64)
				if err != nil {
					return rets, err
				}

				hash, err := GetFileHash(uid)
				if err != nil {
					return rets, err
				}
				er.LogoHash = hash
			}
		}
		er.Id = list[i].Id
		er.Contract = list[i].Contract
		er.Name = list[i].Name
		er.Block = list[i].Block
		er.Hash = hex.EncodeToString(list[i].Hash)
		ret = append(ret, er)
	}

	rets.Total = total
	rets.List = ret

	return rets, nil
}

func GetTopTenMaxKeysEcosystem() ([]EcosystemKeysRatioResponse, error) {
	var (
		total int64
		key   Key
		rets  []EcosystemKeysRatioResponse
	)

	err := GetDB(nil).Table(key.TableName()).Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Top Ten Max Keys Ecosystem Total Failed")
		return rets, err
	}
	err = GetDB(nil).Raw(`
SELECT ecosystem as id,count(1) number,(SELECT name AS name FROM "1_ecosystems" as es WHERE es.id = k1.ecosystem) 
FROM "1_keys" as k1 GROUP BY ecosystem ORDER BY number desc,id asc limit 10
`).Find(&rets).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Top Ten Max Keys Ecosystem Failed")
		return rets, err
	}
	totalDec := decimal.NewFromInt(total)
	for i := 0; i < len(rets); i++ {
		if total > 0 && rets[i].Number > 0 {
			rets[i].Ratio = decimal.NewFromInt(rets[i].Number*100).DivRound(totalDec, 2)
		}
	}

	return rets, nil
}

func GetMultiFeeEcosystemChart() (MultiFeeEcosystemRatioResponse, error) {
	var (
		rets  MultiFeeEcosystemRatioResponse
		total int64
		eco   Ecosystem
	)
	err := GetDB(nil).Table(eco.TableName()).Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Multi Fee Ecosystem Chart Total Failed")
		return rets, err
	}
	err = GetDB(nil).Raw(`
SELECT count(*)multi_fee FROM "1_ecosystems" WHERE 
	coalesce(fee_mode_info->'fee_mode_detail'->'vmCost_fee'->>'flag','0') > '1' OR
	coalesce(fee_mode_info->'fee_mode_detail'->'element_fee'->>'flag','0') > '1' OR
	coalesce(fee_mode_info->'fee_mode_detail'->'storage_fee'->>'flag','0') > '1' OR
	coalesce(fee_mode_info->'fee_mode_detail'->'expedite_fee'->>'flag','0') > '1'
`).Take(&rets).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Multi Fee Ecosystem Chart Failed")
		return rets, err
	}
	if total > 0 {
		rets.UnMultiFee = total - rets.MultiFee
		totalDec := decimal.NewFromInt(total)
		if rets.MultiFee > 0 {
			rets.MultiFeeRatio = decimal.NewFromInt(rets.MultiFee*100).DivRound(totalDec, 2)
		}
		if rets.UnMultiFee > 0 {
			rets.UnMultiFeeRatio = decimal.NewFromInt(rets.UnMultiFee*100).DivRound(totalDec, 2)
		}
	}

	return rets, nil
}

func GetNewNftMinerChangeChart() (NewNftChangeChartResponse, error) {
	var (
		rets NewNftChangeChartResponse
	)
	if !NftMinerReady {
		return rets, nil
	}

	type newNftChangeChart struct {
		Days   string `json:"days"`
		NewNft int64  `json:"new_nft"`
		Stake  int64  `json:"stake"`
	}
	var list []newNftChangeChart

	newSql := `
SELECT CASE WHEN coalesce(ns.days,'')<>'' THEN
ns.days
ELSE
sk.days
END
 FROM (
	SELECT to_char(to_timestamp(date_created),'yyyy-MM-dd') AS days FROM "1_nft_miner_items" GROUP BY days ORDER BY days ASC LIMIT 1
) as ns
full JOIN (
		SELECT to_char(to_timestamp(start_dated),'yyyy-MM-dd') AS days FROM "1_nft_miner_staking" GROUP BY days ORDER BY days ASC LIMIT 1
)AS sk ON(sk.days=ns.days)
`
	var first newNftChangeChart
	f, err := isFound(GetDB(nil).Raw(newSql).Take(&first))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get History New Nft Miner First Time Failed")
		return rets, err
	}
	if !f {
		log.WithFields(log.Fields{"error": err}).Error("Get History New Nft First Time Not Exist")
		return rets, nil
	}
	diff := GetDateDiffFromNow("2006-01-02", first.Days, 0)
	err = GetDB(nil).Raw(`
SELECT f2.days,coalesce(st.new_nft,0) new_nft,coalesce(st.stake,0) stake FROM(
		SELECT to_char(date_trunc('day', (?::TIMESTAMP + f1.offs::INTERVAL)), 'yyyy-MM-dd')as days FROM (
			SELECT generate_series(0, ?, 1) || ' d' as offs
		)as f1
	)as f2
LEFT JOIN(
	SELECT CASE WHEN coalesce(ns.days,'')<>'' THEN
		ns.days
	ELSE
		sk.days
	END,
		coalesce(ns.new_nft,0) as new_nft,coalesce(sk.stake,0)AS stake
	FROM (
		SELECT to_char(to_timestamp(date_created),'yyyy-MM-dd') AS days,count(1)new_nft FROM "1_nft_miner_items" GROUP BY days ORDER BY days ASC
	) as ns
	full JOIN (
		SELECT to_char(to_timestamp(start_dated),'yyyy-MM-dd') AS days,count(1)stake FROM "1_nft_miner_staking" GROUP BY days ORDER BY days ASC
	)AS sk ON(sk.days=ns.days)
)AS st ON(st.days = f2.days) ORDER BY f2.days asc
`, first.Days, diff).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get History New Nft Miner Change Chart Failed")
		return rets, err
	}

	zeroDec := decimal.New(0, 0)
	for i := 0; i < len(list); i++ {
		times, err := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get History New Nft Miner Change Chart ParseInLocation Failed")
			continue
		}
		rets.Time = append(rets.Time, times.Unix())
		rets.New = append(rets.New, list[i].NewNft)
		rets.Stake = append(rets.Stake, list[i].Stake)
		if list[i].Stake > 0 {
			if list[i].NewNft > 0 {
				rets.Ratio = append(rets.Ratio, decimal.NewFromInt(list[i].Stake*100).DivRound(decimal.NewFromInt(list[i].NewNft), 2))
			} else {
				rets.Ratio = append(rets.Ratio, decimal.NewFromInt(1))
			}
		} else {
			rets.Ratio = append(rets.Ratio, zeroDec)
		}
	}

	return rets, nil
}

func GetNftMinerRewardChangeChart() (DaysAmountResponse, error) {
	var (
		rets DaysAmountResponse
	)
	sql := `
SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS days,sum(amount)AS amount 
FROM "1_history" WHERE type = 12 GROUP BY days ORDER BY days ASC
`
	list, err := FindDaysAmount(sql)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get History Nft Reward Change Chart Failed")
		return rets, err
	}
	rets = GetDaysAmountResponse(list)
	return rets, nil
}

func GetEcosystemGovernModelChart() (GovernModelRatioResponse, error) {
	var (
		rets  GovernModelRatioResponse
		total int64
		eco   Ecosystem
	)
	err := GetDB(nil).Table(eco.TableName()).Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Govern Model Chart Total Failed")
		return rets, err
	}
	err = GetDB(nil).Raw(`
SELECT count(1)AS dao_governance FROM "1_ecosystems" WHERE control_mode = 2
`).Take(&rets).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Govern Model Chart Failed")
		return rets, err
	}
	if total > 0 {
		rets.CreatorModel = total - rets.DAOGovernance
		totalDec := decimal.NewFromInt(total)
		if rets.CreatorModel > 0 {
			rets.CreatorRatio = decimal.NewFromInt(rets.CreatorModel*100).DivRound(totalDec, 2)
		}
		if rets.DAOGovernance > 0 {
			rets.DAORatio = decimal.NewFromInt(rets.DAOGovernance*100).DivRound(totalDec, 2)
		}
	}

	return rets, nil
}

func getNodeVoteChangeChangeChart() (NodeVoteChangeResponse, error) {
	var (
		rets NodeVoteChangeResponse
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	rets.Time = make([]int64, getDays)
	rets.Vote = make([]string, getDays)
	if !NodeReady {
		return rets, nil
	}

	sql := fmt.Sprintf(`
SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS days,case WHEN coalesce(sum(amount),0) > 0 THEN
	round(coalesce(sum(amount),0) / 1e12,0)
ELSE
	0
END num FROM "1_history" WHERE type = 20 AND created_at >= %d GROUP BY days ORDER BY days ASC
`, t1.UnixMilli())
	var list []DaysNumber
	err := GetDB(nil).Raw(sql).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get Node Vote Change Change Chart Failed")
		return rets, err
	}
	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.Vote[i] = strconv.FormatInt(GetDaysNumber(rets.Time[i], list), 10)
	}

	return rets, nil
}

func getNodeStakingChangeChangeChart() (*NodeStakingChangeResponse, error) {
	var (
		rets  NodeStakingChangeResponse
		total int64
	)
	rets.Time = make([]string, 0)
	rets.Staking = make([]string, 0)
	if !NodeReady {
		return nil, nil
	}

	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	err := GetDB(nil).Raw(`
SELECT count(1) FROM (
	SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS days FROM "1_history" WHERE type IN(18,19,21) GROUP BY days
)as h1
`).Take(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get Node Staking Change Change Chart Total Failed")
		return nil, err
	}
	timeDbFormat, layout := GetDayNumberFormat(total)

	stakingSql := fmt.Sprintf(`
WITH "1_history" AS (SELECT sum(amount) as amount,
	to_char(to_timestamp(created_at/1000),'%s') AS days
FROM "1_history" WHERE type IN(18,19) AND ecosystem = 1
GROUP BY days ORDER BY days)
SELECT s1.days,
	(SELECT SUM(amount) FROM "1_history" s2 WHERE s2.days <= s1.days AND SUBSTRING(s1.days,0,5) = SUBSTRING(s2.days,0,5)) AS amount
FROM "1_history" s1
`, timeDbFormat)

	withdrawSql := fmt.Sprintf(`
WITH "1_history" AS (SELECT sum(amount) as amount,
	to_char(to_timestamp(created_at/1000),'%s') AS days
FROM "1_history" WHERE type IN(21) AND ecosystem = 1
GROUP BY days ORDER BY days)
SELECT s1.days,
	(SELECT SUM(amount) FROM "1_history" s2 WHERE s2.days <= s1.days AND SUBSTRING(s1.days,0,5) = SUBSTRING(s2.days,0,5)) AS amount
FROM "1_history" s1
`, timeDbFormat)
	var staking []DaysAmount
	var withdraw []DaysAmount
	err = GetDB(nil).Raw(stakingSql).Find(&staking).Error
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("get Node Staking Change Change staking Chart Failed")
		return nil, err
	}

	err = GetDB(nil).Raw(withdrawSql).Find(&withdraw).Error
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("get Node Staking Change Change Withdraw Chart Failed")
		return nil, err
	}
	var startTime time.Time

	if len(staking) > 0 {
		startTime, err = time.ParseInLocation(layout, staking[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"info": err}).Error("get Node Staking Change Change Chart ParseInLocation Failed")
			return nil, err
		}
		lastStakingAmount := decimal.New(0, 0)
		lastDelAmount := decimal.New(0, 0)
		for startTime.Unix() < today.Unix() {
			rets.Time = append(rets.Time, startTime.Format(layout))
			stakingAmount := GetDaysAmount(startTime.Unix(), staking)
			withdrawAmount := GetDaysAmountEqual(startTime.Unix(), withdraw, layout, false)

			if withdrawAmount.String() != "0" {
				lastDelAmount = withdrawAmount
			}

			if stakingAmount == "0" {
				if withdrawAmount.String() != "0" {
					rets.Staking = append(rets.Staking, lastStakingAmount.Sub(withdrawAmount).String())
				} else {
					rets.Staking = append(rets.Staking, lastStakingAmount.String())
				}
			} else {
				stakingDec, err := decimal.NewFromString(stakingAmount)
				if err != nil {
					log.WithFields(log.Fields{"info": err, "staking amount": stakingAmount}).Error("get Node Staking staking amount decimal Failed")
					return nil, err
				}
				lastStakingAmount = stakingDec.Sub(lastDelAmount)
				rets.Staking = append(rets.Staking, lastStakingAmount.String())
			}

			startTime = addTimeFromLayout(layout, startTime)
		}
	}

	return &rets, nil
}

func getTopTenNodeRegionChart() (RegionChangeResponse, error) {
	var (
		rets  RegionChangeResponse
		total int64
		daily DailyNodeReport
	)
	rets.Time = make([]any, 0)
	rets.List = make([][]any, 0)
	if !NodeReady {
		return rets, nil
	}

	type regionInfo struct {
		Time   string `json:"time"`
		Region string `json:"region"`
		Total  string `json:"total"`
	}
	err := GetDB(nil).Table(daily.TableName()).Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("get Top Ten Node Region Chart Total Failed")
		return rets, err
	}
	var list []regionInfo
	timeDbFormat, _ := GetDayNumberFormat(total)
	err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT v4.days as time,array_to_string(array_agg(v4.address),',')region,array_to_string(array_agg(v4.total),',')total FROM(
	SELECT v3.days,v3.address,v3.total,ROW_NUMBER () OVER (PARTITION BY v3.days ORDER BY v3.total desc) AS rowd FROM(
		SELECT to_char(to_timestamp(v1.time),'%s') AS days,address,count(1)total FROM(
			SELECT time,honor_node AS node FROM daily_node_report
			UNION
			SELECT time,candidate_node as node FROM daily_node_report
		)AS v1
		LEFT JOIN(
			SELECT address,value FROM honor_node_info
		)AS v2 ON(v1.node @> ('[{"id":'||CAST(v2.value->>'id' AS numeric)||'}]')::jsonb)
		WHERE v2.address <> '' GROUP BY days,address ORDER BY days asc,total desc
	)AS v3 
)AS v4 WHERE rowd <= 10 GROUP BY time ORDER BY time asc
`, timeDbFormat)).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("get Top Ten Node Region Chart Failed")
		return rets, err
	}

	for i := 0; i < len(list); i++ {
		region := strings.Split(list[i].Region, ",")
		total := strings.Split(list[i].Total, ",")
		t1 := list[i].Time
		rets.Time = append(rets.Time, t1)
		var rts []any

		for key, value := range region {
			//sg.Region = value
			var sg []any
			if len(total)-1 >= key {
				//sg.Total = total[key]
				sg = append(sg, total[key])
			}
			sg = append(sg, value)

			sg = append(sg, t1)

			rts = append(rts, sg)
		}
		rets.List = append(rets.List, rts)
	}

	return rets, nil
}
