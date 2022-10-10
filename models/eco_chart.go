/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
	"time"
)

type DaysAmount struct {
	Days   string          `gorm:"column:days"`
	Amount decimal.Decimal `gorm:"column:amount"`
}

func SyncChartDataMain() {

}

func GetDaysAmount(dayTime int64, list []DaysAmount) string {
	for i := 0; i < len(list); i++ {
		times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
		if dayTime == times.Unix() {
			return list[i].Amount.String()
		}
	}
	return "0"
}

func GetDaysAmountEqual(findTime int64, list []DaysAmount, layout string, areEqual bool) decimal.Decimal {
	for i := 0; i < len(list); i++ {
		times, _ := time.ParseInLocation(layout, list[i].Days, time.Local)
		if areEqual {
			if findTime == times.Unix() {
				return list[i].Amount
			}
		} else {
			if findTime >= times.Unix() {
				return list[i].Amount
			}
		}
	}
	return decimal.New(0, 0)
}

// GetAccountTokenChangeChart
// findTime: 0:Get All  order:Find Start TIME
func GetAccountTokenChangeChart(ecosystem, keyId int64, findTime int64) (AccountAmountChangeBarChart, error) {
	var (
		rets        AccountAmountChangeBarChart
		balanceList []DaysAmount
	)
	//account getTime balance + utxo(input time < getTime <= output time) balance
	err := GetDB(nil).Raw(`
SELECT v2.days,(
	SELECT COALESCE(sum(v1.amount),0)+COALESCE((
				SELECT CASE WHEN (sender_id = ?) THEN
				sender_balance
			ELSE
				recipient_balance
			END
				FROM "1_history"
				WHERE (recipient_id = ? OR sender_id = ?) AND ecosystem = ? 
				AND to_char(to_timestamp(COALESCE(created_at/1000,0)),'yyyy-MM-dd') <= days ORDER BY id DESC LIMIT 1
	),0) AS amount
	FROM(
		SELECT CASE WHEN (sender_id = ?) THEN
				sender_balance
			ELSE
				recipient_balance
			END AS amount
				FROM spent_info_history
				WHERE (recipient_id = ? OR sender_id = ?) AND ecosystem = ?
				AND to_char(to_timestamp(COALESCE(created_at,0)),'yyyy-MM-dd') <= days ORDER BY id DESC LIMIT 1
	)AS v1
				
)AS amount 
FROM(
	SELECT v1.days
	FROM(
		SELECT to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd')AS days FROM "1_history" 
		WHERE (recipient_id = ? OR sender_id = ?) AND ecosystem = ? AND created_at >= ? GROUP BY days
			UNION
		SELECT to_char(to_timestamp(created_at), 'yyyy-mm-dd')AS days FROM spent_info_history 
		WHERE (recipient_id = ? OR sender_id = ?) AND ecosystem = ? AND created_at >= ? GROUP BY days
	)AS v1
)AS v2
ORDER BY days ASC
`,
		keyId, keyId, keyId, ecosystem,
		keyId, keyId, keyId, ecosystem,
		keyId, keyId, ecosystem, findTime,
		keyId, keyId, ecosystem, time.UnixMilli(findTime).Unix()).Find(&balanceList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Account Token Change Chart Failed")
		return rets, nil
	}
	for i := 0; i < len(balanceList); i++ {
		t1, err := time.ParseInLocation("2006-01-02", balanceList[i].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "day": balanceList[i].Days}).Error("Get Account Token Change Chart ParseInLocation Failed")
			return rets, err
		}
		rets.Time = append(rets.Time, t1.Unix())
		rets.Balance = append(rets.Balance, GetDaysAmount(rets.Time[i], balanceList))
	}
	rets.TokenSymbol, rets.Name = Tokens.Get(ecosystem), EcoNames.Get(ecosystem)

	return rets, nil
}

// todo: Need add freeze amount
func GetEcosystemCirculationsChart(ecosystem int64) (EcoCirculationsResponse, error) {
	var (
		cycleDay     int64
		timeDbFormat string
		bk           Block
		his          History
		ret          EcoCirculationsResponse
		err          error
		layout       string
	)
	tz := time.Now()

	if ecosystem == 1 {
		firstBk, err := bk.GetSystemTime()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart system time Failed")
			return ret, err
		}
		cycleDay = int64(time.Unix(tz.Unix(), 0).Sub(time.Unix(firstBk, 0)).Hours() / 24)
	} else {
		f, err := isFound(GetDB(nil).Select("created_at").Where("type = 6").First(&his))
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart first Create token time Failed")
			return ret, err
		}
		if !f {
			return ret, nil
		}
		cycleDay = int64(time.Unix(tz.Unix(), 0).Sub(time.Unix(his.Createdat/1000, 0)).Hours() / 24)
	}

	type circulations struct {
		Days             string
		Circulations     decimal.Decimal
		NftBalanceSupply decimal.Decimal
	}
	type nowChartDataResponse struct {
		Circulations          string
		StakeAmount           string
		FreezeAmount          string //TODO:NEED ADD
		NftMinerBalanceSupply string
		BurningTokens         string
		Combustion            string
		TokenSymbol           string
		Name                  string
		SupplyToken           string
		Emission              string
	}

	escapeAmount := func(findTime int64, list []DaysAmount, layout string, areEqual bool) decimal.Decimal {
		for i := 0; i < len(list); i++ {
			times, _ := time.ParseInLocation(layout, list[i].Days, time.Local)
			if areEqual {
				if findTime == times.Unix() {
					return list[i].Amount
				}
			} else {
				if findTime >= times.Unix() {
					return list[i].Amount
				}
			}
		}
		return decimal.Zero
	}

	escapeCirculations := func(findTime string, list []circulations) decimal.Decimal {
		for i := 0; i < len(list); i++ {
			if list[i].Days == findTime {
				return list[i].Circulations
			}
		}
		return decimal.Zero
	}
	getListTime := func(days, layout string) int64 {
		times, _ := time.ParseInLocation(layout, days, time.Local)
		return times.Unix()
	}

	if cycleDay >= 720 { //years
		timeDbFormat = "yyyy"
		layout = "2006"
	} else if cycleDay >= 60 { //months
		timeDbFormat = "yyyy-MM"
		layout = "2006-01"
	} else { //days
		timeDbFormat = "yyyy-MM-dd"
		layout = "2006-01-02"
	}

	var cir []circulations
	var delCir []DaysAmount
	var newStaked []DaysAmount
	var deleStaked []DaysAmount
	var burning []DaysAmount
	var combustion []DaysAmount
	var emission []DaysAmount
	var supplyToken []DaysAmount
	var nowChart nowChartDataResponse

	if ecosystem == 1 {
		if NodeReady && NftMinerReady {
			err = GetDB(nil).Raw(`
SELECT sum(amount)+COALESCE((SELECT sum(output_value) FROM "spent_info" 
WHERE input_tx_hash is null AND ecosystem = 1 AND output_key_id <> 0),0) AS circulations,
	(SELECT coalesce(sum(stake_amount),0)+
	(SELECT coalesce(sum(earnest),0) FROM "1_candidate_node_decisions" WHERE decision <> 3) AS stake_amount FROM "1_nft_miner_staking" WHERE staking_status = 1),
	(SELECT value AS nft_miner_balance_supply FROM "1_app_params" WHERE "name" = 'nft_miner_balance_supply' AND ecosystem = 1),
 	coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem)),'IBXC') as token_symbol,
	(SELECT name FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem))
FROM "1_keys" AS k1 WHERE ecosystem = 1 and deleted = 0 and blocked = 0 AND id <> 0`).Take(&nowChart).Error
		} else if !NodeReady && !NftMinerReady {
			err = GetDB(nil).Raw(`
SELECT sum(amount)+COALESCE((SELECT sum(output_value) FROM "spent_info" WHERE 
input_tx_hash is null AND ecosystem = 1 AND output_key_id <> 0),0) AS circulations,
	COALESCE((SELECT value FROM "1_app_params" WHERE "name" = 'nft_miner_balance_supply' AND ecosystem = 1),'0') AS nft_miner_balance_supply,
	coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem)),'IBXC') as token_symbol,
(SELECT name FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem))
		FROM "1_keys" AS k1 WHERE ecosystem = 1 and deleted = 0 and blocked = 0 AND id <> 0`).Take(&nowChart).Error
		} else if NodeReady {
			err = GetDB(nil).Raw(`
SELECT sum(amount)+COALESCE((SELECT sum(output_value) FROM "spent_info"
WHERE input_tx_hash is null AND ecosystem = 1 AND output_key_id <> 0),0) AS circulations,
		(SELECT coalesce(sum(earnest),0) FROM "1_candidate_node_decisions" WHERE decision <> 3) AS stake_amount,
		COALESCE((SELECT value FROM "1_app_params" WHERE "name" = 'nft_miner_balance_supply' AND ecosystem = 1),'0') AS nft_miner_balance_supply,
 coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem)),'IBXC') as token_symbol,
(SELECT name FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem))
		FROM "1_keys" AS k1 WHERE ecosystem = 1 and deleted = 0 and blocked = 0 AND id <> 0`).Take(&nowChart).Error
		} else {
			//Nft Miner Ready
			err = GetDB(nil).Raw(`
SELECT sum(amount)+COALESCE((SELECT sum(output_value) FROM "spent_info" 
WHERE input_tx_hash is null AND ecosystem = 1 AND output_key_id <> 0),0) AS circulations,
	(SELECT coalesce(sum(stake_amount),0) AS stake_amount FROM "1_nft_miner_staking" WHERE staking_status = 1),
	COALESCE((SELECT value FROM "1_app_params" WHERE "name" = 'nft_miner_balance_supply' AND ecosystem = 1),'0') AS nft_miner_balance_supply,
	coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem)),'IBXC') as token_symbol,
	(SELECT name FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem))
FROM "1_keys" AS k1 WHERE ecosystem = 1 and deleted = 0 and blocked = 0 AND id <> 0`).Take(&nowChart).Error
		}
		nowChart.SupplyToken = TotalSupplyToken
		nowChart.Emission = "0"
	} else {
		err = GetDB(nil).Raw(`SELECT v1.circulations,v1.burning_tokens,v1.combustion,v1.token_symbol,v1.name,
COALESCE(v2.supply_token,0)supply_token,COALESCE(v2.emission,0)emission FROM(
	SELECT sum(amount)+COALESCE((SELECT sum(output_value) FROM "spent_info" WHERE input_tx_hash is null AND ecosystem = ? AND output_key_id <> 0),0) AS circulations,max(ecosystem) eco_id,
		coalesce((SELECT sum(amount) FROM "1_history" WHERE type = 7 AND ecosystem = max(k1.ecosystem)),0) AS burning_tokens,
		coalesce((SELECT sum(amount) FROM "1_history" WHERE type = 16 AND ecosystem = max(k1.ecosystem)),0) +
		COALESCE((SELECT sum(amount) FROM spent_info_history WHERE type = 6 AND ecosystem = max(k1.ecosystem)),0) AS combustion,
		 (SELECT COALESCE(token_symbol,'') FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem)) as token_symbol,
	(SELECT name FROM "1_ecosystems" as ec WHERE ec.id = max(k1.ecosystem))
	FROM "1_keys" AS k1 WHERE ecosystem = ? and deleted = 0 and blocked = 0 AND id <> 0
)AS v1
LEFT JOIN(
	SELECT COALESCE(amount,0)AS supply_token,ecosystem,(SELECT COALESCE(sum(amount),0) AS emission FROM "1_history" AS h2  
		WHERE h2.id > h1.id AND h2.type = h1.type AND h2.ecosystem = h1.ecosystem) FROM "1_history" AS h1 WHERE type = 6 AND ecosystem = ? ORDER BY id ASC LIMIT 1
)AS v2 ON(v1.eco_id = v2.ecosystem)`, ecosystem, ecosystem, ecosystem).Take(&nowChart).Error
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart Failed")
		return ret, err
	}
	ret.Circulations = nowChart.Circulations
	ret.TokenSymbol = nowChart.TokenSymbol
	ret.StakeAmount = nowChart.StakeAmount
	ret.FreezeAmount = nowChart.FreezeAmount
	ret.NftBalanceSupply = nowChart.NftMinerBalanceSupply
	ret.Combustion = nowChart.Combustion
	ret.BurningTokens = nowChart.BurningTokens
	ret.Name = nowChart.Name
	ret.SupplyToken = nowChart.SupplyToken
	ret.Emission = nowChart.Emission
	//get In the day Circulations
	err = GetDB(nil).Raw(`
SELECT cir.days,cir.circulations,sy.nft_balance_supply
 FROM (
	WITH "1_history" AS (SELECT sum(amount) as amount,max(ecosystem) ecosystem,
	to_char(to_timestamp(created_at/1000),?) AS days
	FROM "1_history" WHERE type IN(4,6,12,14,21,22) AND ecosystem = ?
	GROUP BY days)
	SELECT s1.days,s1.amount,s1.ecosystem,
	CASE WHEN s1.ecosystem = 1 THEN
			5250000000000000000+(SELECT SUM(amount) FROM "1_history" s2 WHERE s2.days <= s1.days AND SUBSTRING(s1.days,0,5) = SUBSTRING(s2.days,0,5))
		ELSE
			(SELECT SUM(amount) FROM "1_history" s2 WHERE s2.days <= s1.days AND SUBSTRING(s1.days,0,5) = SUBSTRING(s2.days,0,5))
	END AS circulations
	FROM "1_history" AS s1 
)AS cir
LEFT JOIN(
	WITH "1_history" AS (SELECT sum(amount) as amount,max(ecosystem) ecosystem,
	to_char(to_timestamp(created_at/1000),?) AS days
	FROM "1_history" WHERE type IN(12) AND ecosystem = ?
	GROUP BY days
	ORDER BY days)
	SELECT s1.days,s1.amount,s1.ecosystem,
			CAST((SELECT value AS nft_miner_total_supply FROM "1_app_params" WHERE "name" = 'nft_miner_total_supply' AND ecosystem = s1.ecosystem) as numeric)-
				(SELECT SUM(amount) FROM "1_history" s2 WHERE s2.days <= s1.days AND SUBSTRING(s1.days,0,5) = SUBSTRING(s2.days,0,5)) AS nft_balance_supply
	FROM "1_history" AS s1 
)AS sy ON(sy.days = cir.days)
ORDER BY cir.days asc
`, timeDbFormat, ecosystem, timeDbFormat, ecosystem).Find(&cir).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart cir Failed")
		return ret, err
	}

	err = GetDB(nil).Raw(`
SELECT del.days,del.total_amount as amount
 FROM (WITH "1_history" AS (SELECT sum(amount) as amount,max(ecosystem) ecosystem,
to_char(to_timestamp(created_at/1000),?) AS days
FROM "1_history" WHERE type IN(7,13,16,17,18,19,20) AND ecosystem = ?
GROUP BY days
ORDER BY days desc)
SELECT s1.days,s1.amount,s1.ecosystem,
		(SELECT SUM(amount) FROM "1_history" s2 WHERE s2.days <= s1.days AND SUBSTRING(s1.days,0,5) = SUBSTRING(s2.days,0,5)) AS total_amount
FROM "1_history" AS s1 
)AS del
`, timeDbFormat, ecosystem).Find(&delCir).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart delCir Failed")
		return ret, err
	}

	if ecosystem != 1 {
		//get Burning Tokens by days
		err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT CASE WHEN s1.days > '' THEN
	 s1.days
	ELSE
	 s2.days
END days,COALESCE(s1.amount,0)+COALESCE(s2.amount) AS amount
FROM(
	SELECT to_char(to_timestamp(created_at/1000),'%s') AS days,sum(amount)AS amount FROM "1_history" WHERE type = 7 AND 
	ecosystem = ? GROUP BY days 
)AS s1
FULL JOIN(
	SELECT to_char(to_timestamp(created_at),'%s') AS days,sum(amount)AS amount FROM 
		spent_info_history WHERE type = 2 AND recipient_id = 0 AND ecosystem = ? GROUP BY days 
)AS s2 ON(s2.days = s1.days)
ORDER BY days desc
`, timeDbFormat, timeDbFormat), ecosystem, ecosystem).Find(&burning).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart burning Failed")
			return ret, err
		}

		err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT CASE WHEN s1.days > '' THEN
	 s1.days
	ELSE
	 s2.days
END days,COALESCE(s1.amount,0)+COALESCE(s2.amount) AS amount
FROM(
	SELECT to_char(to_timestamp(created_at/1000),'%s') AS days,sum(amount)AS amount FROM "1_history" WHERE type = 16 AND 
	ecosystem = ? GROUP BY days 
)AS s1
FULL JOIN(
	SELECT to_char(to_timestamp(created_at),'%s') AS days,sum(amount)AS amount FROM 
		spent_info_history WHERE type = 6 AND ecosystem = ? GROUP BY days 
)AS s2 ON(s2.days = s1.days)
`, timeDbFormat, timeDbFormat), ecosystem, ecosystem).Find(&combustion).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart combustion Failed")
			return ret, err
		}

		var supply History
		f, err := isFound(GetDB(nil).Select("created_at,amount,id").Where("type = 6 AND ecosystem = ?", ecosystem).
			Order("id asc").Limit(1).Take(&supply))
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart supply token Failed")
			return ret, err
		}
		if f {
			supplyToken = append(supplyToken, DaysAmount{Days: time.UnixMilli(supply.Createdat).Format(layout), Amount: supply.Amount})

			err = GetDB(nil).Table(his.TableName()).Select("to_char(to_timestamp(created_at/1000),?) AS days,sum(amount) as amount", timeDbFormat).
				Where("type = 6 AND ecosystem = ? AND id > ?", ecosystem, supply.ID).Group("days").Order("days desc").Find(&emission).Error
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart emission Failed")
				return ret, err
			}
		}
	} else {
		//get Create staked by days
		q := GetDB(nil).Table(his.TableName()).Select("to_char(to_timestamp(created_at/1000),?) AS days,sum(amount) as amount", timeDbFormat).
			Where("ecosystem = ?", ecosystem).Group("days").Order("days desc")
		err = q.Where("type IN(13,18,19,20)").Find(&newStaked).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart cirStaked Failed")
			return ret, err
		}

		//get Transfer out staked by days
		err = GetDB(nil).Table(his.TableName()).Select("to_char(to_timestamp(created_at/1000),?) AS days,sum(amount) as amount", timeDbFormat).
			Where("ecosystem = ? AND type IN(14,21,22)", ecosystem).Group("days").Order("days desc").Find(&deleStaked).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart cirStaked Failed")
			return ret, err
		}
	}

	lastCirAmount := decimal.Zero
	lastDelAmount := decimal.Zero
	stakingAmount := decimal.Zero
	lastStakingAmount := decimal.Zero
	burnAmount := decimal.Zero
	lastBurnAmount := decimal.Zero
	supplyAmount := decimal.Zero
	emissionAmount := decimal.Zero
	lastSupply := decimal.Zero
	lastEmission := decimal.Zero
	combusAmount := decimal.Zero
	lastCombusAmount := decimal.Zero
	lastNftBanlance := decimal.Zero
	var startTime time.Time
	end := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	ret.Change.Time = make([]string, 0)
	ret.Change.Circulations = make([]string, 0)
	ret.Change.StakeAmount = make([]string, 0)
	ret.Change.FreezeAmount = make([]string, 0)
	ret.Change.NftBalanceSupply = make([]string, 0)
	ret.Change.BurningTokens = make([]string, 0)
	ret.Change.Combustion = make([]string, 0)
	ret.Change.SupplyToken = make([]string, 0)
	ret.Change.Emission = make([]string, 0)

	getCir := func(cir []circulations, findTime string) {
		var isFindout bool

		for i := 0; i < len(cir); i++ {
			t2 := getListTime(cir[i].Days, layout)
			t1 := cir[i].Days
			if findTime == t1 {
				isFindout = true
				ret.Change.Time = append(ret.Change.Time, t1)

				var circulations decimal.Decimal
				var delCirculations decimal.Decimal
				circulations = escapeCirculations(t1, cir)

				delCirculations = escapeAmount(t2, delCir, layout, false)
				if ecosystem != 1 {
					burnAmount = escapeAmount(t2, burning, layout, true)
					supplyAmount = escapeAmount(t2, supplyToken, layout, true)
					emissionAmount = escapeAmount(t2, emission, layout, true)
					combusAmount = escapeAmount(t2, combustion, layout, true)

					if burnAmount.Equal(decimal.Zero) {
						ret.Change.BurningTokens = append(ret.Change.BurningTokens, lastBurnAmount.String())
					} else {
						ret.Change.BurningTokens = append(ret.Change.BurningTokens, burnAmount.Add(lastBurnAmount).String())
						lastBurnAmount = burnAmount.Add(lastBurnAmount)
					}

					if combusAmount.Equal(decimal.Zero) {
						ret.Change.Combustion = append(ret.Change.Combustion, lastCombusAmount.String())
					} else {
						ret.Change.Combustion = append(ret.Change.Combustion, combusAmount.Add(lastCombusAmount).String())
						lastCombusAmount = combusAmount.Add(lastCombusAmount)
					}

					if supplyAmount.Equal(decimal.Zero) {
						ret.Change.SupplyToken = append(ret.Change.SupplyToken, lastSupply.String())
					} else {
						ret.Change.SupplyToken = append(ret.Change.SupplyToken, supplyAmount.Add(lastSupply).String())
						lastSupply = supplyAmount.Add(lastSupply)
					}

					if emissionAmount.Equal(decimal.Zero) {
						ret.Change.Emission = append(ret.Change.Emission, lastEmission.String())
					} else {
						ret.Change.Emission = append(ret.Change.Emission, emissionAmount.Add(lastEmission).String())
						lastEmission = emissionAmount.Add(lastEmission)
					}
				} else {
					if !cir[i].NftBalanceSupply.Equal(decimal.Zero) {
						lastNftBanlance = cir[i].NftBalanceSupply
					}
					ret.Change.NftBalanceSupply = append(ret.Change.NftBalanceSupply, lastNftBanlance.String())

					stakingAmount = escapeAmount(t2, newStaked, layout, true).Sub(escapeAmount(t2, deleStaked, layout, true))
					if stakingAmount.Equal(decimal.Zero) {
						ret.Change.StakeAmount = append(ret.Change.StakeAmount, lastStakingAmount.String())
					} else {
						ret.Change.StakeAmount = append(ret.Change.StakeAmount, stakingAmount.Add(lastStakingAmount).String())
						lastStakingAmount = stakingAmount.Add(lastStakingAmount)
					}
					ret.Change.SupplyToken = append(ret.Change.SupplyToken, TotalSupplyToken)
					ret.Change.Emission = append(ret.Change.Emission, "0")
				}

				if !delCirculations.Equal(decimal.Zero) {
					lastDelAmount = delCirculations
				}
				if circulations.Equal(decimal.Zero) {
					if !delCirculations.Equal(decimal.Zero) {
						ret.Change.Circulations = append(ret.Change.Circulations, lastCirAmount.Sub(delCirculations).String())
					} else {
						ret.Change.Circulations = append(ret.Change.Circulations, lastCirAmount.String())
					}
				} else {
					lastCirAmount = circulations.Sub(lastDelAmount)
					ret.Change.Circulations = append(ret.Change.Circulations, circulations.Sub(lastDelAmount).String())
				}
				break
			}
		}
		if !isFindout {
			times, _ := time.ParseInLocation(layout, findTime, time.Local)
			t1 := times.Unix()
			var delCirculations decimal.Decimal
			delCirculations = escapeAmount(t1, delCir, layout, true)
			if !delCirculations.Equal(decimal.Zero) {
				ret.Change.Circulations = append(ret.Change.Circulations, lastCirAmount.Sub(delCirculations.Sub(lastDelAmount)).String())
			} else {
				ret.Change.Circulations = append(ret.Change.Circulations, ret.Change.Circulations[len(ret.Change.Circulations)-1])
			}
			combusAmount = escapeAmount(t1, combustion, layout, true)
			supplyAmount = escapeAmount(t1, supplyToken, layout, true)
			emissionAmount = escapeAmount(t1, emission, layout, true)
			burnAmount = escapeAmount(t1, burning, layout, true)

			ret.Change.Time = append(ret.Change.Time, findTime)
			if ecosystem != 1 {
				if burnAmount.Equal(decimal.Zero) {
					ret.Change.BurningTokens = append(ret.Change.BurningTokens, lastBurnAmount.String())
				} else {
					ret.Change.BurningTokens = append(ret.Change.BurningTokens, burnAmount.Add(lastBurnAmount).String())
					lastBurnAmount = burnAmount.Add(lastBurnAmount)
				}
				if combusAmount.Equal(decimal.Zero) {
					ret.Change.Combustion = append(ret.Change.Combustion, lastCombusAmount.String())
				} else {
					ret.Change.Combustion = append(ret.Change.Combustion, combusAmount.Add(lastCombusAmount).String())
					lastCombusAmount = combusAmount.Add(lastCombusAmount)
				}

				if supplyAmount.Equal(decimal.Zero) {
					ret.Change.SupplyToken = append(ret.Change.SupplyToken, lastSupply.String())
				} else {
					ret.Change.SupplyToken = append(ret.Change.SupplyToken, supplyAmount.Add(lastSupply).String())
					lastSupply = supplyAmount.Add(lastSupply)
				}

				if emissionAmount.Equal(decimal.Zero) {
					ret.Change.Emission = append(ret.Change.Emission, lastEmission.String())
				} else {
					ret.Change.Emission = append(ret.Change.Emission, emissionAmount.Add(lastEmission).String())
					lastEmission = emissionAmount.Add(lastEmission)
				}
			} else {
				stakingAmount = escapeAmount(t1, newStaked, layout, true).Sub(escapeAmount(t1, deleStaked, layout, true))
				if stakingAmount.Equal(decimal.Zero) {
					ret.Change.StakeAmount = append(ret.Change.StakeAmount, lastStakingAmount.String())
				} else {
					ret.Change.StakeAmount = append(ret.Change.StakeAmount, stakingAmount.Add(lastStakingAmount).String())
					lastStakingAmount = stakingAmount.Add(lastStakingAmount)
				}
				ret.Change.NftBalanceSupply = append(ret.Change.NftBalanceSupply, lastNftBanlance.String())
				//ret.Change.StakeAmount = append(ret.Change.StakeAmount, lastStakingAmount.String())
				ret.Change.SupplyToken = append(ret.Change.SupplyToken, TotalSupplyToken)
				ret.Change.Emission = append(ret.Change.Emission, "0")
			}

		}

	}

	if len(cir) > 0 {
		startTime, err = time.ParseInLocation(layout, cir[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Ecosystem Circulations Chart ParseInLocation Failed")
			return ret, err
		}
		switch layout {
		case "2006-01": //month
			end = time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
		case "2006": //year
			end = time.Date(tz.Year(), tz.Month(), 0, 0, 0, 0, 0, tz.Location())
		}
		for startTime.Unix() <= end.Unix() {
			getCir(cir, startTime.Format(layout))
			switch layout {
			case "2006-01-02":
				startTime = startTime.AddDate(0, 0, 1)
			case "2006-01":
				startTime = startTime.AddDate(0, 1, 0)
			default:
				startTime = startTime.AddDate(1, 0, 0)
			}
		}
	}

	return ret, nil
}

func GetTopTenHasTokenAccountToRedis() {
	HistoryWG.Add(1)
	defer func() {
		HistoryWG.Done()
	}()

	var eco []Ecosystem
	err := GetDB(nil).Select("id").Find(&eco).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get top ten has token account all ecosystem failed")
	}

	for _, v := range eco {
		rets, err := getEcoTopTenHasTokenAccount(v.ID)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("get top ten has token account failed")
			return
		}

		res, err := msgpack.Marshal(rets)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("get top ten has token account msgpack failed")
			return
		}

		rd := RedisParams{
			Key:   "top-ten-has-token-" + strconv.FormatInt(v.ID, 10),
			Value: string(res),
		}
		if err := rd.Set(); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("get top ten has token account set redis error")
			return
		}
	}
}

func GetTopTenHasTokenAccountFromRedis(ecosystem int64) (*EcoTopTenHasTokenResponse, error) {
	var rets EcoTopTenHasTokenResponse
	var err error
	rd := RedisParams{
		Key:   "top-ten-has-token-" + strconv.FormatInt(ecosystem, 10),
		Value: "",
	}
	if err := rd.Get(); err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return nil, nil
		}
		log.WithFields(log.Fields{"warn": err}).Warn("get top ten has token account From Redis getDb err")
		return nil, err
	}
	err = msgpack.Unmarshal([]byte(rd.Value), &rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get top ten has token account From Redis msgpack err")
		return nil, err
	}

	return &rets, nil
}

func getEcoTopTenHasTokenAccount(ecosystem int64) (*EcoTopTenHasTokenResponse, error) {
	var (
		err  error
		rets EcoTopTenHasTokenResponse
	)
	type accountHold struct {
		KeyId       int64
		Amount      decimal.Decimal
		StakeAmount decimal.Decimal
	}
	var list []accountHold

	if NftMinerReady || NodeReady {
		err = GetDB(nil).Raw(`
SELECT v1.key_id,COALESCE(v2.amount,0)+
	COALESCE(v2.stake_amount,0)+
	COALESCE((SELECT sum(output_value) FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = v1.ecosystem AND output_key_id = v1.key_id),0)
	AS amount,COALESCE(v2.stake_amount,0)AS stake_amount 
FROM(							
	SELECT id AS key_id,ecosystem FROM "1_keys" WHERE ecosystem = ? AND blocked = 0 AND deleted = 0 AND amount +
			to_number(coalesce(NULLIF(lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+
			to_number(coalesce(NULLIF(lock->>'candidate_referendum',''),'0'),'999999999999999999999999') +
			to_number(coalesce(NULLIF(lock->>'candidate_substitute',''),'0'),'999999999999999999999999')>0
		UNION
	SELECT output_key_id AS key_id,ecosystem FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ? GROUP BY output_key_id,ecosystem
)AS v1
LEFT JOIN(
		SELECT id,amount,ecosystem,
			to_number(coalesce(NULLIF(lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+
			to_number(coalesce(NULLIF(lock->>'candidate_referendum',''),'0'),'999999999999999999999999') +
			to_number(coalesce(NULLIF(lock->>'candidate_substitute',''),'0'),'999999999999999999999999') AS stake_amount
		FROM "1_keys"
)AS v2 ON(v2.id = v1.key_id AND v2.ecosystem = v1.ecosystem)
ORDER BY amount DESC
`, ecosystem, ecosystem).Find(&list).Error
	} else {
		err = GetDB(nil).Raw(`
SELECT v1.key_id,COALESCE(v2.amount,0)+
	COALESCE((SELECT sum(output_value) FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = v1.ecosystem AND output_key_id = v1.key_id),0)
	AS amount
FROM(							
	SELECT id AS key_id,ecosystem FROM "1_keys" WHERE ecosystem = ? AND blocked = 0 AND deleted = 0 AND amount > 0
		UNION
	SELECT output_key_id AS key_id,ecosystem FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ? GROUP BY output_key_id,ecosystem
)AS v1
LEFT JOIN(
		SELECT id,amount,ecosystem FROM "1_keys"
)AS v2 ON(v2.id = v1.key_id AND v2.ecosystem = v1.ecosystem)
ORDER BY amount DESC
`, ecosystem, ecosystem).Find(&list).Error
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Eco TopTen HasToken Account Chart Failed")
		return nil, err
	}

	type totalAmountRet struct {
		TotalAmount decimal.Decimal `json:"total_amount"`
	}
	var totalAmount totalAmountRet
	err = GetDB(nil).Raw(`
SELECT sum(k2.amount)+sum(to_number(coalesce(NULLIF(k2.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999'))+
	sum(to_number(coalesce(NULLIF(k2.lock->>'candidate_referendum',''),'0'),'999999999999999999999999')) + 
	sum(to_number(coalesce(NULLIF(k2.lock->>'candidate_substitute',''),'0'),'999999999999999999999999'))+
	COALESCE((SELECT sum(output_value) FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ?),0) 
as total_amount FROM "1_keys" AS k2 WHERE ecosystem = ? AND blocked = 0 AND deleted = 0
`, ecosystem, ecosystem).Take(&totalAmount).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("Get Eco TopTen HasToken Account Chart Total Amount Failed")
		return nil, err
	}

	otherAmount := decimal.Zero
	otherStaking := decimal.Zero
	for key, val := range list {
		if key >= 10 {
			otherAmount = otherAmount.Add(val.Amount)
			otherStaking = otherStaking.Add(val.StakeAmount)
		} else {
			var rt accountRatio
			rt.Account = converter.AddressToString(val.KeyId)
			rt.StakeAmount = val.StakeAmount.String()
			rt.Amount = val.Amount.String()
			rt.AccountedFor = val.Amount.Mul(decimal.NewFromInt(100)).DivRound(totalAmount.TotalAmount, 2)
			rets.List = append(rets.List, rt)
		}
	}
	if !otherAmount.IsZero() {
		var rt accountRatio
		rt.Account = "Other"
		rt.StakeAmount = otherStaking.String()
		rt.Amount = otherAmount.String()
		rt.AccountedFor = otherAmount.Mul(decimal.NewFromInt(100)).DivRound(totalAmount.TotalAmount, 2)
		rets.List = append(rets.List, rt)
	}
	rets.TokenSymbol, rets.Name = Tokens.Get(ecosystem), EcoNames.Get(ecosystem)

	return &rets, nil
}

func GetEcoTopTenTxAccountChart(ecosystem int64) (*EcoTopTenTxAmountResponse, error) {
	var (
		err   error
		rets  EcoTopTenTxAmountResponse
		total decimal.Decimal
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := today.AddDate(0, 0, -1*getDays)
	type findStruct struct {
		Keyid  int64           `json:"keyid"`
		Amount decimal.Decimal `json:"amount"`
	}
	var ret []findStruct
	err = GetDB(nil).Raw(`
SELECT COALESCE(sum(amount),0)+
	COALESCE((SELECT sum(amount) FROM spent_info_history WHERE ecosystem = ? AND type <> 1 AND created_at >= ? AND sender_id <> 0 AND sender_id <> 5555),0)+
	COALESCE((SELECT sum(amount) FROM spent_info_history WHERE ecosystem = ? AND type <> 1 AND created_at >= ? AND recipient_id <> 0 AND recipient_id <> 5555),0)+
	COALESCE((SELECT sum(amount) FROM "1_history" WHERE ecosystem = ? AND created_at >= ? AND sender_id <> 0 AND sender_id <> 5555),0)
AS total
FROM "1_history" WHERE ecosystem = ? AND created_at >= ? AND recipient_id <> 0 AND recipient_id <> 5555
`, ecosystem, t1.Unix(),
		ecosystem, t1.Unix(),
		ecosystem, t1.UnixMilli(),
		ecosystem, t1.UnixMilli()).Take(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Eco TopTen Tx Account Chart Total Failed")
		return nil, err
	}

	err = GetDB(nil).Raw(`
SELECT keyid,amount
FROM(
	SELECT keyid,ecosystem,
		(SELECT COALESCE(sum(amount),0)+
			(SELECT COALESCE(sum(amount),0) FROM "1_history" WHERE recipient_id = t1.keyid AND ecosystem = t1.ecosystem AND created_at >= ?)+
			(SELECT COALESCE(sum(amount),0) FROM spent_info_history WHERE recipient_id = t1.keyid AND ecosystem = t1.ecosystem AND type <> 1 AND created_at >= ?)+
			(SELECT COALESCE(sum(amount),0) FROM spent_info_history WHERE sender_id = t1.keyid AND ecosystem = t1.ecosystem AND type <> 1 AND created_at >= ?)
		FROM "1_history" WHERE sender_id = t1.keyid AND ecosystem = t1.ecosystem AND created_at >= ?)AS amount
	FROM(
		SELECT sender_id as keyid,max(ecosystem) ecosystem FROM "1_history" WHERE ecosystem = ? AND created_at >= ? GROUP BY sender_id
		 UNION
		SELECT recipient_id as keyid,max(ecosystem) ecosystem FROM "1_history" WHERE ecosystem = ? AND created_at >= ? GROUP BY recipient_id
		 UNION
		SELECT output_key_id AS keyid,max(ecosystem) ecosystem FROM spent_info AS s1 LEFT JOIN log_transactions AS l1 ON(l1.hash = s1.output_tx_hash) 
			WHERE ecosystem = ? AND timestamp >= ? GROUP BY output_key_id
		 UNION
		SELECT output_key_id AS keyid,max(ecosystem) ecosystem FROM spent_info AS s1 LEFT JOIN log_transactions AS l1 ON(l1.hash = s1.input_tx_hash) 
			WHERE ecosystem = ? AND timestamp >= ? GROUP BY output_key_id
	)AS t1
)AS tt
WHERE keyid <> 0 AND keyid <> 5555
order by amount desc limit 10
`, t1.UnixMilli(),
		t1.Unix(),
		t1.Unix(),
		t1.UnixMilli(),
		ecosystem, t1.UnixMilli(),
		ecosystem, t1.UnixMilli(),
		ecosystem, t1.UnixMilli(),
		ecosystem, t1.UnixMilli()).Find(&ret).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Eco TopTen Tx Account Chart Failed")
		return nil, err
	}
	for _, value := range ret {
		var qt accountRatio
		qt.Amount = value.Amount.String()
		qt.Account = converter.AddressToString(value.Keyid)
		qt.AccountedFor = value.Amount.Mul(decimal.NewFromInt(100)).DivRound(total, 2)
		rets.List = append(rets.List, qt)
	}
	rets.TokenSymbol, rets.Name = Tokens.Get(ecosystem), EcoNames.Get(ecosystem)

	return &rets, nil
}

func GetGasCombustionPieChart(ecosystem int64) (EcoGasFeeResponse, error) {
	var rets EcoGasFeeResponse

	err := GetDB(nil).Raw(`
SELECT COALESCE(sum(amount),0) + 
		COALESCE((SELECT sum(amount) FROM spent_info_history WHERE ecosystem = ? AND "type" IN(3,4)),0) AS gas_fee,
	coalesce((SELECT sum(amount) FROM "1_history" WHERE type = 16 AND ecosystem = ?),'0')+
	COALESCE((SELECT sum(amount) FROM spent_info_history WHERE type = 6 AND ecosystem = ?),'0') AS combustion 
FROM "1_history" WHERE type IN(1,2) AND ecosystem = ?
`, ecosystem, ecosystem, ecosystem, ecosystem).Take(&rets).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("Get Gas Combustion Pie Chart Failed")
		return rets, err
	}
	rets.Name = EcoNames.Get(ecosystem)
	rets.TokenSymbol = Tokens.Get(ecosystem)

	return rets, nil
}

func GetGasCombustionLineChart(ecosystem int64) (EcoGasFeeChangeResponse, error) {
	var (
		his  History
		rets EcoGasFeeChangeResponse
	)

	tz := time.Unix(GetNowTimeUnix(), 0)
	end := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	type gasFeeChange struct {
		EcoGasFeeResponse
		Time string `json:"time"`
	}

	getDaysGasFee := func(list []gasFeeChange, getDay string) gasFeeChange {
		for i := 0; i < len(list); i++ {
			t1 := list[i].Time
			if getDay == t1 {
				return list[i]
			}
		}
		return gasFeeChange{}
	}

	var list []gasFeeChange

	err := GetDB(nil).Select("created_at").Where("ecosystem = ?", ecosystem).Order("id asc").Limit(1).Take(&his).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Gas Combustion Line Chart Total Failed")
		return rets, err
	}

	created := time.UnixMilli(his.Createdat)
	if his.Createdat <= 0 {
		created = time.Unix(FirstBlockTime, 0)
	}
	diffDay := int64(tz.Sub(created).Hours() / 24)
	startTime := time.Date(created.Year(), created.Month(), created.Day(), 0, 0, 0, 0, created.Location())

	dbFormat, layout := GetDayNumberFormat(diffDay)
	err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT time,sum(gas_fee)gas_fee,sum(combustion)combustion FROM(
	SELECT CASE WHEN h1.days > '' THEN
	 h1.days
	ELSE
	 h2.days
	END AS time,coalesce(h1.gas_fee) AS gas_fee,coalesce(h2.combustion,0) as combustion

	FROM (
		SELECT to_char(to_timestamp(created_at/1000),'%s') AS days,sum(amount)AS gas_fee FROM "1_history" WHERE (type = 1 or type = 2) AND 
	ecosystem = ? GROUP BY days
	)AS h1

	FULL JOIN(
		SELECT to_char(to_timestamp(created_at/1000),'%s') AS days,sum(amount)AS combustion FROM "1_history" WHERE type = 16 AND 
	ecosystem = ? GROUP BY days 
	) as h2 ON(h1.days = h2.days)

		UNION
	SELECT CASE WHEN s1.days > '' THEN
	 s1.days
	ELSE
	 s2.days
	END AS time,COALESCE(s1.gas_fee,0) AS gas_fee,COALESCE(s2.combustion,0) AS combustion
	FROM(
		SELECT to_char(to_timestamp("created_at"),'%s') days,COALESCE(sum(amount),0)AS gas_fee FROM spent_info_history 
		WHERE ecosystem = ? AND "type" IN(3,4) GROUP BY days
	)AS s1
	FULL JOIN(
		SELECT to_char(to_timestamp(created_at),'%s') AS days,sum(amount)AS combustion FROM 
		spent_info_history WHERE type = 6 AND ecosystem = ? GROUP BY days 
	)AS s2 ON(s2.days = s1.days)
)AS v1
GROUP BY time
ORDER BY time ASC
`, dbFormat, dbFormat, dbFormat, dbFormat), ecosystem, ecosystem, ecosystem, ecosystem).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Gas Combustion Line Chart Failed")
		return rets, err
	}

	switch layout {
	case "2006-01": //month
		end = time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	case "2006": //year
		end = time.Date(tz.Year(), tz.Month(), 0, 0, 0, 0, 0, tz.Location())
	}
	for startTime.Unix() <= end.Unix() {
		findTime := startTime.Format(layout)
		info := getDaysGasFee(list, findTime)

		rets.Time = append(rets.Time, findTime)
		if info.GasFee == "" {
			rets.GasFee = append(rets.GasFee, decimal.Zero.String())
		} else {
			rets.GasFee = append(rets.GasFee, info.GasFee)
		}

		if info.Combustion == "" {
			rets.Combustion = append(rets.Combustion, decimal.Zero.String())
		} else {
			rets.Combustion = append(rets.Combustion, info.Combustion)
		}

		switch layout {
		case "2006-01-02":
			startTime = startTime.AddDate(0, 0, 1)
		case "2006-01":
			startTime = startTime.AddDate(0, 1, 0)
		default:
			startTime = startTime.AddDate(1, 0, 0)
		}
	}

	rets.TokenSymbol, rets.Name = Tokens.Get(ecosystem), EcoNames.Get(ecosystem)

	return rets, nil
}

func GetEco15DayTxAmountChart(ecosystem int64) (EcoTxAmountDiffResponse, error) {
	var rets EcoTxAmountDiffResponse
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	rets.Time = make([]int64, getDays)
	rets.Amount = make([]string, getDays)

	var list []DaysAmount
	err := GetDB(nil).Raw(`
	SELECT CASE WHEN v1.days <> '' THEN
		v1.days
	ELSE
		v2.days
	END,COALESCE(v1.tx_amount,0)+COALESCE(v2.tx_amount,0)AS amount
	FROM(
		SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS days,sum(amount)tx_amount FROM "1_history" WHERE ecosystem = ? GROUP BY days
	)AS v1
	FULL JOIN(
		SELECT to_char(to_timestamp(created_at),'yyyy-MM-dd') AS days,sum(amount)tx_amount FROM spent_info_history WHERE ecosystem = ? AND type <> 1 GROUP BY days 
	)AS v2 ON(v2.days = v1.days)
	ORDER BY days DESC LIMIT ?
`, ecosystem, ecosystem, getDays).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Eco 15 Day Tx Amount Chart Failed")
		return rets, err
	}

	rets.TokenSymbol, rets.Name = Tokens.Get(ecosystem), EcoNames.Get(ecosystem)

	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.Amount[i] = GetDaysAmount(rets.Time[i], list)
	}

	return rets, nil
}

func GetEco15DayGasFeeChart(ecosystem int64) (EcoTxGasFeeDiffResponse, error) {
	var (
		rets EcoTxGasFeeDiffResponse
		list []DaysAmount
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	err := GetDB(nil).Raw(`
SELECT days,sum(amount)amount FROM(
	SELECT to_char(to_timestamp("created_at"/1000),'yyyy-MM-dd') days,sum(amount) amount 
	FROM "1_history" WHERE ecosystem = ? AND created_at >= ? AND "type" IN(1,2) GROUP BY days
		UNION
	SELECT to_char(to_timestamp("created_at"),'yyyy-MM-dd') days,sum(amount) amount 
	FROM spent_info_history WHERE ecosystem = ? AND created_at >= ? AND "type" IN(3,4) GROUP BY days
)AS v1 GROUP BY days
ORDER BY days
`, ecosystem, t1.UnixMilli(), ecosystem, t1.Unix()).Find(&list).Error
	if err != nil {
		return rets, err
	}
	rets.TokenSymbol, rets.Name = Tokens.Get(ecosystem), EcoNames.Get(ecosystem)

	rets.Time = make([]int64, getDays)
	rets.EcoGasAmount = make([]string, getDays)
	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.EcoGasAmount[i] = GetDaysAmount(rets.Time[i], list)
	}
	return rets, nil
}

func GetEco15DayActiveKeysChart(ecosystem int64) (KeyInfoChart, error) {
	var keyChart KeyInfoChart
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	keyChart.Time = make([]int64, getDays)
	keyChart.ActiveKey = make([]int64, getDays)

	var activeList []DaysNumber
	err := GetDB(nil).Raw(`
SELECT days,count(keyid) as num  FROM (
	SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') days,sender_id as keyid FROM "1_history" WHERE sender_id <> 0 AND created_at >= ? AND ecosystem = ? GROUP BY days, sender_id
	 UNION 
	SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') days,recipient_id as keyid  FROM "1_history" WHERE recipient_id <> 0 AND created_at >= ? AND ecosystem = ? GROUP BY days, recipient_id 
	 UNION
	SELECT to_char(to_timestamp(timestamp/1000),'yyyy-MM-dd') days,output_key_id AS keyid FROM spent_info AS s1 LEFT JOIN log_transactions AS l1 ON(l1.hash = s1.output_tx_hash) 
	WHERE timestamp >= ? AND ecosystem = ? GROUP BY days,output_key_id
	 UNION
	SELECT to_char(to_timestamp(timestamp/1000),'yyyy-MM-dd') days,output_key_id AS keyid FROM spent_info AS s1 LEFT JOIN log_transactions AS l1 ON(l1.hash = s1.input_tx_hash) 
	WHERE timestamp >= ? AND ecosystem = ? GROUP BY days,output_key_id
) as tt GROUP BY days ORDER BY days DESC
`, t1.Unix(), ecosystem, t1.Unix(), ecosystem, t1.Unix(), ecosystem, t1.Unix(), ecosystem).Find(&activeList).Error
	if err != nil {
		return keyChart, err
	}

	for i := 0; i < len(keyChart.Time); i++ {
		keyChart.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		keyChart.ActiveKey[i] = GetDaysNumber(keyChart.Time[i], activeList)
	}
	keyChart.Name = EcoNames.Get(ecosystem)

	return keyChart, nil
}

func GetEco15DayTransactionChart(ecosystem int64) (TxListChart, error) {
	var (
		txChart TxListChart
		list    []DaysNumber
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	err := GetDB(nil).Raw(fmt.Sprintf(`SELECT to_char(to_timestamp("timestamp"/1000),'yyyy-MM-dd') days,
count(1) AS num FROM log_transactions WHERE ecosystem_id = %d AND "timestamp" >= %d GROUP BY days`, ecosystem, t1.UnixMilli())).Find(&list).Error
	if err != nil {
		return txChart, err
	}
	txChart.Time = make([]int64, getDays)
	txChart.Tx = make([]int64, getDays)
	for i := 0; i < len(txChart.Time); i++ {
		txChart.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		txChart.Tx[i] = GetDaysNumber(txChart.Time[i], list)
	}
	txChart.Name = EcoNames.Get(ecosystem)

	return txChart, nil
}

func GetEco15DayStorageCapacitysChart(ecosystem int64) (StorageCapacitysChart, error) {
	var (
		rets StorageCapacitysChart
		list []DaysNumber
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	const getDays = 15
	t1 := yesterday.AddDate(0, 0, -1*getDays)

	err := GetDB(nil).Raw(fmt.Sprintf(`SELECT to_char(to_timestamp("timestamp"/1000),'yyyy-MM-dd') days,
	sum(length("tx_data")) num FROM log_transactions as l1 LEFT JOIN transaction_data as t1 ON(t1.hash = l1.hash) 
	WHERE ecosystem_id = %d AND "timestamp" >= %d GROUP BY days`, ecosystem, t1.UnixMilli())).Find(&list).Error
	if err != nil {
		return rets, err
	}
	rets.Time = make([]int64, getDays)
	rets.StorageCapacitys = make([]string, getDays)
	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.StorageCapacitys[i] = ToCapacityMb(GetDaysNumber(rets.Time[i], list))
	}
	rets.Name = EcoNames.Get(ecosystem)
	return rets, nil
}
