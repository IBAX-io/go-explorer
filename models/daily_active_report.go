/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type DailyActiveReport struct {
	ID            int64  `gorm:"primary_key;not null" json:"id"`
	Time          int64  `gorm:"column:time;not null" json:"time"`
	ActiveAccount int64  `gorm:"column:active;not null" json:"active_account"`
	Ratio         string `gorm:"column:ratio;type:varchar(30);not null" json:"ratio"`
	TotalKey      int64  `gorm:"column:total_key;not null" json:"total_key"`
	RelativeRatio string `gorm:"column:relative_ratio;type:varchar(30);not null" json:"relative_ratio"`
	TxNumber      int64  `gorm:"column:tx_number;not null" json:"tx_number"`
	TxAmount      string `gorm:"column:tx_amount;type:decimal(30);not null" json:"tx_amount"`
}

type DaysActiveReport struct {
	Days     string          `gorm:"column:days"`
	Active   int64           `gorm:"column:active"`
	TotalKey int64           `gorm:"column:total_key"`
	Ratio    decimal.Decimal `gorm:"column:ratio"`
	TxNumber int64           `gorm:"column:tx_number"`
	TxAmount decimal.Decimal `gorm:"column:tx_amount"`
}

var (
	getDailyActiveReportListEnd bool = true
)

func (dt *DailyActiveReport) TableName() string {
	return "daily_active_report"
}

func (p *DailyActiveReport) CreateTable() (err error) {
	err = nil
	if !HasTableOrView(p.TableName()) {
		if err = GetDB(nil).Migrator().CreateTable(p); err != nil {
			return err
		}
	}
	return err
}

func InitDailyActiveReport() error {
	var (
		dt DailyActiveReport
	)
	err := dt.CreateTable()
	if err != nil {
		return err
	}
	return nil
}

func (dt *DailyActiveReport) GetLast() (bool, error) {
	return isFound(GetDB(nil).Last(dt))
}

func (dt *DailyActiveReport) GetFirst() (bool, error) {
	return isFound(GetDB(nil).First(dt))
}

func (dt *DailyActiveReport) GetList() ([]DailyActiveReport, error) {
	var rets []DailyActiveReport
	err := GetDB(nil).Order("id asc").Find(&rets).Error
	return rets, err
}

func (dt *DailyActiveReport) GetTimeLine(stTime, edTime int64) ([]DailyActiveReport, error) {
	var rets []DailyActiveReport
	err := GetDB(nil).Order("id asc").Where("time >= ? AND time < ?", stTime, edTime).Find(&rets).Error
	return rets, err
}

func (dt *DailyActiveReport) Insert() error {
	p := DailyActiveReport{
		Time:          dt.Time,
		ActiveAccount: dt.ActiveAccount,
		Ratio:         dt.Ratio,
		TotalKey:      dt.TotalKey,
		RelativeRatio: dt.RelativeRatio,
		TxNumber:      dt.TxNumber,
		TxAmount:      dt.TxAmount,
	}
	if err := GetDB(nil).Model(&DailyActiveReport{}).Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func InsertDailyActiveReport() {
	ChartWG.Add(1)
	defer func() {
		ChartWG.Done()
	}()
	var (
		dt DailyActiveReport
	)
	now := time.Now()
	f, err := dt.GetLast()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Last Daily Active Report Failed")
		return
	}
	if f {
		lastTime := time.Unix(dt.Time, 0)
		diffDay := int(now.Sub(lastTime).Hours() / 24)
		var lastTotalKey = dt.TotalKey
		//fmt.Printf("dt time:%s,now time:%s\n", lastTime.String(), now.String())
		//fmt.Printf("diffDay:%d\n", diffDay)

		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if diffDay > 1 {
			for i := 0; i < diffDay; i++ {
				t1 := lastTime.AddDate(0, 0, i+1)
				if t1.Unix() >= today.Unix() {
					continue
				}

				endTime := t1.AddDate(0, 0, 1)
				report, _, err := GetTimeLineDaysActiveReport(t1.UnixMilli(), endTime.UnixMilli())
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Insert Daily Active Report Failed")
					continue
				}
				var rt DailyActiveReport
				//dayTime, err := time.ParseInLocation("2006-01-02", report.Days, time.Local)
				//if err != nil {
				//	log.WithFields(log.Fields{"error": err, "days": report.Days}).Error("Insert Daily Active Report ParseInLocation Failed")
				//	continue
				//}
				if report.TotalKey != 0 {
					lastTotalKey = report.TotalKey
				} else {
					report.TotalKey = lastTotalKey
				}
				rt.Time = t1.Unix()
				rt.TxAmount = report.TxAmount.String()
				rt.TxNumber = report.TxNumber
				rt.TotalKey = report.TotalKey
				rt.Ratio = "0"
				active := decimal.NewFromInt(report.Active)
				if active.GreaterThan(decimal.Zero) {
					totalKey := decimal.NewFromInt(report.TotalKey)
					if totalKey.GreaterThan(decimal.Zero) {
						rt.Ratio = active.Mul(decimal.NewFromInt(100)).DivRound(totalKey, 2).String()
					} else {
						rt.Ratio = "100"
					}
				}
				rt.ActiveAccount = report.Active
				//fmt.Printf("rt time:%d,now time:%s\n", rt.Time, time.Now().String())
				preSt := t1.AddDate(0, 0, -1)
				preEndTime := preSt.AddDate(0, 0, 1).UnixMilli()

				preReport, f1, err := GetTimeLineDaysActiveReport(preSt.UnixMilli(), preEndTime)
				if err != nil {
					log.WithFields(log.Fields{"error": err, "preSt": preSt.String()}).Error("Insert Daily Active Report pre year Failed")
					continue
				}
				if f1 {
					diff := rt.ActiveAccount - preReport.Active
					if diff != 0 {
						if preReport.Active > 0 {
							diffDec := decimal.NewFromInt(diff)
							preDec := decimal.NewFromInt(preReport.Active)

							rt.RelativeRatio = diffDec.DivRound(preDec, 2).String()
						} else {
							rt.RelativeRatio = "100"
						}
					} else {
						rt.RelativeRatio = "0"
					}
				} else {
					rt.RelativeRatio = "0"
				}

				if err := rt.Insert(); err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Insert Daily Active Report Insert Failed")
					continue
				}
			}
		}

		getActiveReportToRedis("one_month")
		getActiveReportToRedis("three_month")
		getActiveReportToRedis("one_year")
		getActiveReportToRedis("all")
	} else {
		if getDailyActiveReportListEnd {
			list, err := GetDailyActiveReportList()
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Get Daily Active Report List Failed")
				return
			}
			if len(list) > 0 {
				err = dt.InsertList(list)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Daily Active Report Insert List Failed")
					return
				}
			}
		}
	}
}

func GetTimeLineDaysActiveReport(st, ed int64) (DaysActiveReport, bool, error) {
	//fmt.Printf("st:%d,ed:%d\n", st, ed)
	var (
		ratio DaysActiveReport
		his   History
		total int64
	)
	f, err := isFound(GetDB(nil).Table(his.TableName()).Where("created_at < ?", ed).Count(&total))
	if err != nil {
		return ratio, f, err
	}
	if !f {
		return ratio, f, nil
	}

	f, err = isFound(GetDB(nil).Raw(fmt.Sprintf(`
SELECT h2.ds as days,coalesce(h3.num, 0) AS active,coalesce(h2.num+5,0) AS total_key,h2.tx_number,h2.tx_amount

	FROM (
			(
				SELECT CASE WHEN v1.ds <> '' THEN
					v1.ds
				ELSE
					v2.ds
				END,COALESCE(v1.tx_amount,0)+COALESCE(v2.tx_amount,0)AS tx_amount
				FROM(
					SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS ds,sum(amount)tx_amount FROM "1_history" WHERE ecosystem = 1 AND type <> 24 AND created_at >= %d AND created_at < %d GROUP BY ds
				)AS v1
				FULL JOIN(
					SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS ds,sum(amount)tx_amount FROM spent_info_history WHERE ecosystem = 1 AND type <> 1 AND created_at >= %d AND created_at < %d GROUP BY ds 
				)AS v2 ON(v2.ds = v1.ds)
				ORDER BY ds DESC
			)AS h1
		
			LEFT JOIN(
				WITH rollback_tx AS(
					SELECT to_char(to_timestamp(log.time), 'yyyy-mm-dd') AS days,count(1) num
					FROM (SELECT tx_hash,table_id FROM rollback_tx WHERE table_name = '1_keys' AND table_id like '%%,1' AND data = '') AS rb LEFT JOIN(
						SELECT timestamp/1000 as time,hash FROM log_transactions 
					)AS log ON (log.hash = rb.tx_hash) GROUP BY days
				)
				SELECT rk1.days,
					(SELECT SUM(num) FROM rollback_tx s2 WHERE s2.days <= rk1.days) as num
				FROM rollback_tx AS rk1 
			) AS rk ON (h1.ds = rk.days)

			LEFT JOIN(
				SELECT count(1)AS tx_number,to_char(to_timestamp(timestamp/1000), 'yyyy-mm-dd') AS days FROM log_transactions WHERE ecosystem_id = 1 GROUP BY days
			)AS v1 ON(v1.days = h1.ds)
	
		) AS h2 
		
	LEFT JOIN(
			SELECT h1.days,count(1)as num FROM(
					SELECT sender_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd') AS days FROM "1_history" 
					WHERE sender_id <> 0 AND ecosystem = 1 AND created_at >= %d AND created_at < %d GROUP BY sender_id,days
					 UNION 
					SELECT recipient_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd') AS days FROM "1_history" 
					WHERE recipient_id <> 0 AND ecosystem = 1 AND created_at >= %d AND created_at < %d GROUP BY recipient_id,days
					 UNION
					SELECT v1.output_key_id AS keyid,to_char(to_timestamp(v2."timestamp"/1000), 'yyyy-mm-dd') AS days FROM spent_info AS v1 
						LEFT JOIN log_transactions AS v2 ON(v1.input_tx_hash = v2.hash) WHERE 
						v1.ecosystem = 1 AND v2."timestamp" >= %d AND v2."timestamp" < %d AND input_tx_hash is not NULL GROUP BY v1.output_key_id,days
					 UNION
					SELECT v1.output_key_id AS keyid,to_char(to_timestamp(v2."timestamp"/1000), 'yyyy-mm-dd') AS days FROM spent_info AS v1 
						LEFT JOIN log_transactions AS v2 ON(v1.output_tx_hash = v2.hash) 
						WHERE v1.ecosystem = 1 AND v2."timestamp" >= %d AND v2."timestamp" < %d GROUP BY v1.output_key_id,days
			) AS h1 GROUP BY h1.days
			
	) AS h3 ON(h3.days = h2.ds)
	
	ORDER BY h2.ds DESC`, st, ed, st, ed, st, ed, st, ed, st, ed, st, ed)).First(&ratio))
	return ratio, f, err
}

func (dt *DailyActiveReport) InsertList(list []DailyActiveReport) error {
	if err := GetDB(nil).Create(&list).Error; err != nil {
		return err
	}
	return nil
}

func GetDailyActiveReportList() ([]DailyActiveReport, error) {
	var (
		dtList []DailyActiveReport
	)
	getDailyActiveReportListEnd = false
	defer func() {
		getDailyActiveReportListEnd = true
	}()

	var list []DaysActiveReport
	//now := time.Now()
	err := GetDB(nil).Raw(
		`SELECT h2.ds as days,coalesce(h3.num, 0) AS active,coalesce(h2.num+5,0) AS total_key,h2.tx_number,h2.tx_amount
	FROM (
			(
				SELECT CASE WHEN v1.ds <> '' THEN
					v1.ds
				ELSE
					v2.ds
				END,COALESCE(v1.tx_amount,0)+COALESCE(v2.tx_amount,0)AS tx_amount
				FROM(
					SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS ds,sum(amount)tx_amount FROM "1_history" WHERE ecosystem = 1 AND type <> 24 GROUP BY ds
				)AS v1
				FULL JOIN(
					SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS ds,sum(amount)tx_amount FROM spent_info_history WHERE ecosystem = 1 AND type <> 1 GROUP BY ds 
				)AS v2 ON(v2.ds = v1.ds)
				ORDER BY ds asc
			)AS h1

			LEFT JOIN(
				WITH rollback_tx AS(
					SELECT to_char(to_timestamp(log.time), 'yyyy-mm-dd') AS days,count(1) num
					FROM (SELECT tx_hash,table_id FROM rollback_tx WHERE table_name = '1_keys' AND table_id like '%,1' AND data = '') AS rb LEFT JOIN(
						SELECT timestamp/1000 as time,hash FROM log_transactions 
					)AS log ON (log.hash = rb.tx_hash) GROUP BY days
				)
				SELECT rk1.days,
					(SELECT SUM(num) FROM rollback_tx s2 WHERE s2.days <= rk1.days) as num
				FROM rollback_tx AS rk1 
			) AS rk ON (h1.ds = rk.days)
			
			LEFT JOIN(
				SELECT count(1)AS tx_number,to_char(to_timestamp(timestamp/1000), 'yyyy-mm-dd') AS days FROM log_transactions WHERE ecosystem_id = 1 GROUP BY days
			)AS v1 ON(v1.days = h1.ds)
	
		) AS h2

	LEFT JOIN(
			SELECT h1.days,count(1)as num FROM(
					SELECT sender_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd') AS days FROM "1_history" WHERE sender_id <> 0 AND ecosystem = 1 GROUP BY sender_id,days
					 UNION 
					SELECT recipient_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm-dd') AS days FROM "1_history" WHERE recipient_id <> 0 AND ecosystem = 1 GROUP BY recipient_id,days
					 UNION
					SELECT output_key_id AS keyid,to_char(to_timestamp(timestamp/1000),'yyyy-MM-dd') days FROM spent_info AS s1 LEFT JOIN 
					 log_transactions AS l1 ON(l1.hash = s1.output_tx_hash)	WHERE ecosystem = 1 GROUP BY days,output_key_id
					 UNION
					SELECT output_key_id AS keyid,to_char(to_timestamp(timestamp/1000),'yyyy-MM-dd') days FROM spent_info AS s1 LEFT JOIN 
					 log_transactions AS l1 ON(l1.hash = s1.input_tx_hash)	WHERE ecosystem = 1 GROUP BY days,output_key_id
			) AS h1 GROUP BY h1.days
			
	) AS h3 ON(h3.days = h2.ds)
	
	ORDER BY h2.ds ASC`).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Daily Active Report List Failed")
		return nil, err
	}

	var lastTotalKey int64
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	getDaysReport := func(findTime time.Time, list []DaysActiveReport) (DailyActiveReport, error) {
		var dt DailyActiveReport
		dt.Time = findTime.Unix()
		for i := 0; i < len(list); i++ {
			dayTime, err := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
			if err != nil {
				log.WithFields(log.Fields{"error": err, "days": list[i].Days}).Error("Get Daily Active Report List ParseInLocation Failed")
				continue
			}
			if dayTime.Unix() == findTime.Unix() {
				if list[i].TotalKey != 0 {
					lastTotalKey = list[i].TotalKey
				}
				dt.TotalKey = lastTotalKey
				dt.ActiveAccount = list[i].Active
				dt.Ratio = "0"
				active := decimal.NewFromInt(list[i].Active)
				if active.GreaterThan(decimal.Zero) {
					totalKey := decimal.NewFromInt(lastTotalKey)
					if totalKey.GreaterThan(decimal.Zero) {
						dt.Ratio = active.Mul(decimal.NewFromInt(100)).DivRound(totalKey, 2).String()
					} else {
						dt.Ratio = "100"
					}
				}
				dt.TxAmount = list[i].TxAmount.String()
				dt.TxNumber = list[i].TxNumber

				preSt := dayTime.AddDate(0, 0, -1)
				endTime := preSt.AddDate(0, 0, 1).UnixMilli()

				preReport, f1, err := GetTimeLineDaysActiveReport(preSt.UnixMilli(), endTime)
				if err != nil {
					log.WithFields(log.Fields{"error": err, "preYear": preSt.UnixMilli()}).Error("Get Time Line Days Active Report List Failed")
					continue
				}

				//fmt.Printf("nowRatio:%s,preRatio:%s\n", list[i].Ratio.String(), relativeRatio.Ratio.String())

				if f1 {
					diff := dt.ActiveAccount - preReport.Active
					if diff != 0 {
						if preReport.Active > 0 {
							diffDec := decimal.NewFromInt(diff)
							preDec := decimal.NewFromInt(preReport.Active)

							dt.RelativeRatio = diffDec.DivRound(preDec, 2).String()
						} else {
							dt.RelativeRatio = "100"
						}
					} else {
						dt.RelativeRatio = "0"
					}
				} else {
					dt.RelativeRatio = "0"
				}
				return dt, nil
			}
		}

		preSt := findTime.AddDate(0, 0, -1)
		endTime := preSt.AddDate(0, 0, 1).UnixMilli()
		preReport, f1, err := GetTimeLineDaysActiveReport(preSt.UnixMilli(), endTime)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "preYear": preSt.UnixMilli()}).Error("Get Time Line Days Active Report List Failed")
			return dt, err
		}
		if f1 {
			diff := dt.ActiveAccount - preReport.Active
			if diff != 0 {
				if preReport.Active > 0 {
					diffDec := decimal.NewFromInt(diff)
					preDec := decimal.NewFromInt(preReport.Active)

					dt.RelativeRatio = diffDec.DivRound(preDec, 2).String()
				} else {
					dt.RelativeRatio = "100"
				}
			} else {
				dt.RelativeRatio = "0"
			}
		} else {
			dt.RelativeRatio = "0"
		}
		dt.Ratio = "0"
		dt.TxAmount = "0"
		dt.TotalKey = lastTotalKey

		return dt, nil
	}

	var startTime time.Time
	if len(list) >= 1 {
		startTime, err = time.ParseInLocation("2006-01-02", list[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get NewKeys Chart Block Id List ParseInLocation Failed")
			return nil, err
		}
		for startTime.Unix() < today.Unix() {
			dt, err := getDaysReport(startTime, list)
			if err != nil {
				return nil, err
			}
			dtList = append(dtList, dt)
			startTime = startTime.AddDate(0, 0, 1)
		}
	}

	return dtList, nil

}

func getDailyActiveReportChart(startTime, endTime int64) (DailyActiveChartResponse, error) {
	var (
		dt        DailyActiveReport
		rets      DailyActiveChartResponse
		maxActive int64
		maxTime   int64
		minActive int64
		minTime   int64
		list      []DailyActiveReport
		err       error
	)
	if startTime != 0 && endTime != 0 {
		list, err = dt.GetTimeLine(startTime, endTime)
		if err != nil {
			return rets, err
		}
	} else {
		list, err = dt.GetList()
		if err != nil {
			return rets, err
		}
	}
	rets.Time = make([]int64, len(list))
	rets.Info = make([]DailyActiveReport, len(list))
	rets.Info = list
	for i := 0; i < len(list); i++ {
		if i == 0 || list[i].ActiveAccount < minActive {
			minActive = list[i].ActiveAccount
			minTime = list[i].Time
		}
		if i == 0 || list[i].ActiveAccount > maxActive {
			maxActive = list[i].ActiveAccount
			maxTime = list[i].Time
		}
		rets.Time[i] = list[i].Time
	}
	rets.MinActive = minActive
	rets.MinTime = minTime
	rets.MaxActive = maxActive
	rets.MaxTime = maxTime

	return rets, nil
}

func getActiveReportToRedis(search string) {
	var (
		rets  DailyActiveChartResponse
		err   error
		rdKey string
	)
	switch search {
	case "one_month":
		rdKey = "one_month_active_report"
		rets, err = getOneMonthYearActiveReport(search)
	case "three_month":
		rdKey = "three_month_active_report"
		rets, err = getTimeLineThreeMonthActiveReport()
	case "one_year":
		rdKey = "one_year_active_report"
		rets, err = getOneMonthYearActiveReport(search)
	default:
		rdKey = "daily_active_report"
		rets, err = getDailyActiveReportChart(0, 0)
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err, "search": search}).Error("Get Active Report To Redis Error")
		return
	}

	val, err := MsgpackMarshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Active Report To Redis Marshal Error")
		return
	}

	rd := RedisParams{
		Key:   rdKey,
		Value: string(val),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get Active Report To Redis set redis error")
		return
	}
}

func getActiveReportFromRedis(search string) (*DailyActiveChartResponse, error) {
	var (
		err error
	)
	chart := &DailyActiveChartResponse{}
	rd := RedisParams{
		Key:   "daily_active_report",
		Value: "",
	}
	switch search {
	case "one_month":
		rd.Key = "one_month_active_report"
	case "three_month":
		rd.Key = "three_month_active_report"
	case "one_year":
		rd.Key = "one_year_active_report"
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get Active Report From Redis getDb err")
		return nil, err
	}
	err = MsgpackUnmarshal([]byte(rd.Value), chart)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

func getOneMonthYearActiveReport(search string) (DailyActiveChartResponse, error) {
	var (
		list      []DaysActiveReport
		rets      DailyActiveChartResponse
		maxActive int64
		maxTime   int64
		minActive int64
		minTime   int64
		dtList    []DailyActiveReport
		sql       string
		layout    string
	)

	if search == "one_year" {
		layout = "2006"
		sql = `
SELECT dt.days,coalesce(h3.num,0) as active,dt.total_key,CASE WHEN (h3.num > 0) THEN
		case when (dt.total_key = 0) THEN
			100
		ELSE
			round(
				h3.num * 100 / dt.total_key , 2) 
		END
	ELSE
		0
	END AS ratio,dt.tx_number,dt.tx_amount FROM
(SELECT to_char(to_timestamp(time), 'yyyy') AS days,sum(tx_amount) tx_amount,sum(tx_number) tx_number,max(total_key) total_key FROM "daily_active_report" GROUP BY days) AS dt
LEFT JOIN(
			SELECT h1.days,count(1)as num FROM(
					SELECT sender_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy') AS days FROM "1_history" WHERE sender_id <> 0 AND ecosystem = 1 GROUP BY sender_id,days
					 UNION 
					SELECT recipient_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy') AS days FROM "1_history" WHERE recipient_id <> 0 AND ecosystem = 1  GROUP BY recipient_id,days
					 UNION
					SELECT v1.output_key_id AS keyid,to_char(to_timestamp(v2."timestamp"/1000), 'yyyy-mm-dd') AS days FROM spent_info AS v1 
						LEFT JOIN log_transactions AS v2 ON(v1.input_tx_hash = v2.hash) WHERE v1.ecosystem = 1 AND input_tx_hash is not NULL GROUP BY v1.output_key_id,days
					 UNION
					SELECT v1.output_key_id AS keyid,to_char(to_timestamp(v2."timestamp"/1000), 'yyyy-mm-dd') AS days FROM spent_info AS v1 
						LEFT JOIN log_transactions AS v2 ON(v1.output_tx_hash = v2.hash) WHERE v1.ecosystem = 1 GROUP BY v1.output_key_id,days
			) AS h1 GROUP BY h1.days
			
	) AS h3 ON(h3.days = dt.days)
ORDER BY dt.days ASC
`
	} else {
		layout = "2006-01"
		sql = `
SELECT dt.days,coalesce(h3.num,0) as active,dt.total_key,CASE WHEN (h3.num > 0) THEN
		case when (dt.total_key = 0) THEN
			100
		ELSE
			round(
				h3.num * 100 / dt.total_key , 2) 
		END
	ELSE
		0
	END AS ratio,dt.tx_number,dt.tx_amount FROM
(SELECT to_char(to_timestamp(time), 'yyyy-mm') AS days,sum(tx_amount) tx_amount,sum(tx_number) tx_number,max(total_key) total_key FROM "daily_active_report" GROUP BY days) AS dt
LEFT JOIN(
			SELECT h1.days,count(1)as num FROM(
					SELECT sender_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm') AS days FROM "1_history" WHERE sender_id <> 0 AND ecosystem = 1 GROUP BY sender_id,days
					 UNION 
					SELECT recipient_id as keyid,to_char(to_timestamp(created_at/1000), 'yyyy-mm') AS days FROM "1_history" WHERE recipient_id <> 0 AND ecosystem = 1  GROUP BY recipient_id,days
					 UNION
					SELECT v1.output_key_id AS keyid,to_char(to_timestamp(v2."timestamp"/1000), 'yyyy-mm-dd') AS days FROM spent_info AS v1 
						LEFT JOIN log_transactions AS v2 ON(v1.input_tx_hash = v2.hash) WHERE v1.ecosystem = 1 AND input_tx_hash is not NULL GROUP BY v1.output_key_id,days
					 UNION
					SELECT v1.output_key_id AS keyid,to_char(to_timestamp(v2."timestamp"/1000), 'yyyy-mm-dd') AS days FROM spent_info AS v1 
						LEFT JOIN log_transactions AS v2 ON(v1.output_tx_hash = v2.hash) WHERE v1.ecosystem = 1 GROUP BY v1.output_key_id,days
			) AS h1 GROUP BY h1.days
			
	) AS h3 ON(h3.days = dt.days)
ORDER BY dt.days ASC
`
	}
	err := GetDB(nil).Raw(sql).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get One Month Active Report Failed")
		return rets, err
	}

	now := time.Now()
	for i := 0; i < len(list); i++ {
		var dt DailyActiveReport
		dayTime, err := time.ParseInLocation(layout, list[i].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "days": list[i].Days}).Error("Get One Month Active Report ParseInLocation Failed")
			continue
		}
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if dayTime.Unix() >= today.Unix() {
			continue
		}
		dt.Time = dayTime.Unix()
		dt.ID = int64(i + 1)
		dt.TxAmount = list[i].TxAmount.String()
		dt.TxNumber = list[i].TxNumber
		dt.Ratio = list[i].Ratio.String()
		dt.TotalKey = list[i].TotalKey
		dt.ActiveAccount = list[i].Active

		preSt := dayTime.AddDate(-1, 0, 0)
		endTime := preSt.AddDate(1, 0, 0).UnixMilli()
		if layout == "2006-01" {
			preSt = dayTime.AddDate(0, -1, 0)
			endTime = preSt.AddDate(0, 1, 0).UnixMilli()
		}

		preReport, f1, err := GetTimeLineDaysActiveReport(preSt.UnixMilli(), endTime)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "preYear": preSt.String()}).Error("Get Time Line Days Active Report Failed")
			continue
		}

		//fmt.Printf("nowRatio:%s,preRatio:%s\n", list[i].Ratio.String(), relativeRatio.Ratio.String())
		if f1 {
			diff := dt.ActiveAccount - preReport.Active
			if diff != 0 {
				if preReport.Active > 0 {
					diffDec := decimal.NewFromInt(diff)
					preDec := decimal.NewFromInt(preReport.Active)

					dt.RelativeRatio = diffDec.DivRound(preDec, 2).String()
				} else {
					dt.RelativeRatio = "100"
				}
			} else {
				dt.RelativeRatio = "0"
			}
		} else {
			dt.RelativeRatio = "0"
		}

		dtList = append(dtList, dt)
	}

	rets.Time = make([]int64, len(dtList))
	rets.Info = make([]DailyActiveReport, len(dtList))
	rets.Info = dtList
	for i := 0; i < len(dtList); i++ {
		if i == 0 || dtList[i].ActiveAccount < minActive {
			minActive = dtList[i].ActiveAccount
			minTime = dtList[i].Time
		}
		if i == 0 || dtList[i].ActiveAccount > maxActive {
			maxActive = dtList[i].ActiveAccount
			maxTime = dtList[i].Time
		}
		rets.Time[i] = dtList[i].Time
	}
	rets.MinActive = minActive
	rets.MinTime = minTime
	rets.MaxActive = maxActive
	rets.MaxTime = maxTime

	return rets, nil
}

func getTimeLineActiveNumber(stTime, edTime int64) (int64, error) {
	var number CountInt64
	f, err := isFound(GetDB(nil).Raw(fmt.Sprintf(`
SELECT count(1) FROM(
	SELECT sender_id as keyid FROM "1_history" WHERE sender_id <> 0 AND ecosystem = 1 AND created_at >= %d AND created_at < %d GROUP BY sender_id
		 UNION 
	SELECT recipient_id as keyid FROM "1_history" WHERE recipient_id <> 0 AND ecosystem = 1 AND created_at >= %d AND created_at < %d GROUP BY recipient_id
		UNION
	SELECT v1.output_key_id AS keyid FROM spent_info AS v1 
			LEFT JOIN log_transactions AS v2 ON(v1.input_tx_hash = v2.hash) WHERE 
			v1.ecosystem = 1 AND input_tx_hash is not NULL AND v2."timestamp" >= %d AND v2."timestamp" < %d GROUP BY v1.output_key_id
		UNION
	SELECT v1.output_key_id AS keyid FROM spent_info AS v1 
		LEFT JOIN log_transactions AS v2 ON(v1.output_tx_hash = v2.hash) WHERE 
		v1.ecosystem = 1 AND v2."timestamp" >= %d AND v2."timestamp" < %d GROUP BY v1.output_key_id
) AS h1
`, stTime, edTime, stTime, edTime, stTime, edTime, stTime, edTime)).Take(&number))
	if err != nil {
		return 0, err
	}
	if !f {
		return 0, nil
	}
	return number.Count, nil
}

func getTimeLineThreeMonthActiveReport() (DailyActiveChartResponse, error) {
	var (
		rets      DailyActiveChartResponse
		list      []DaysActiveReport
		maxActive int64
		maxTime   int64
		minActive int64
		minTime   int64
		dtList    []DailyActiveReport
		dd        DailyActiveReport
		dr        DailyActiveReport
	)
	layout := "2006-01"
	f, err := dd.GetFirst()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Three Month Active Report First Failed")
		return rets, err
	}
	if !f {
		return rets, nil
	}
	f, err = dr.GetLast()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Three Month Active Report Last Failed")
		return rets, err
	}
	if !f {
		return rets, nil
	}

	str := time.Unix(dd.Time, 0).Format(layout)
	stTime, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "str": str}).Error("Get One Month Active Report ParseInLocation st Time Failed")
		return rets, err
	}
	//stTime := time.Unix(dd.Time, 0)
	edTime := stTime.AddDate(0, 3, 0)
	//fmt.Printf("stTime:%d,edTime:%d,dr time:%d\n", stTime.Unix(), edTime.Unix(), dr.Time)
	for edTime.Unix() < dr.Time {
		var info DaysActiveReport
		err := GetDB(nil).Table(dr.TableName()).Select("sum(tx_amount) tx_amount,sum(tx_number) tx_number,max(total_key) total_key").
			Where("time >=? AND time <?", stTime.Unix(), edTime.Unix()).Take(&info).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Three Month Active Report Info Failed")
			return rets, err
		}
		active, err := getTimeLineActiveNumber(stTime.UnixMilli(), edTime.UnixMilli())
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Three Month Active Report Active Number Failed")
			return rets, err
		}
		ave := decimal.NewFromInt(active * 100)
		total := decimal.NewFromInt(info.TotalKey)
		if active > 0 {
			if info.TotalKey == 0 {
				info.Ratio = decimal.New(100, 0)
			} else {
				info.Ratio = ave.DivRound(total, 2)
			}
		}
		//fmt.Printf("ave:%s,total:%s\n", ave.String(), total.String())
		info.Days = stTime.Format(layout)
		info.Active = active

		list = append(list, info)
		stTime = edTime
		edTime = edTime.AddDate(0, 3, 0)
	}

	//now := time.Now()
	for i := 0; i < len(list); i++ {
		var dt DailyActiveReport
		dayTime, err := time.ParseInLocation(layout, list[i].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "days": list[i].Days}).Error("Get One Month Active Report ParseInLocation Failed")
			continue
		}
		//today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		//if dayTime.Unix() >= today.Unix() {
		//	continue
		//}
		dt.Time = dayTime.Unix()
		dt.ID = int64(i + 1)
		dt.TxAmount = list[i].TxAmount.String()
		dt.TxNumber = list[i].TxNumber
		dt.Ratio = list[i].Ratio.String()
		dt.TotalKey = list[i].TotalKey
		dt.ActiveAccount = list[i].Active

		preResult := dayTime.AddDate(0, -3, 0)
		endTime := preResult.AddDate(0, 3, 0).UnixMilli()

		preReport, f1, err := GetTimeLineDaysActiveReport(preResult.UnixMilli(), endTime)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "preResult": preResult}).Error("Get Time Line Days Active Report Failed")
			continue
		}

		if f1 {
			diff := dt.ActiveAccount - preReport.Active
			if diff != 0 {
				if preReport.Active > 0 {
					diffDec := decimal.NewFromInt(diff)
					preDec := decimal.NewFromInt(preReport.Active)

					dt.RelativeRatio = diffDec.DivRound(preDec, 2).String()
				} else {
					dt.RelativeRatio = "100"
				}
			} else {
				dt.RelativeRatio = "0"
			}
		} else {
			dt.RelativeRatio = "0"
		}

		dtList = append(dtList, dt)
	}

	rets.Time = make([]int64, len(dtList))
	rets.Info = make([]DailyActiveReport, len(dtList))
	rets.Info = dtList
	for i := 0; i < len(dtList); i++ {
		if i == 0 || dtList[i].ActiveAccount < minActive {
			minActive = dtList[i].ActiveAccount
			minTime = dtList[i].Time
		}
		if i == 0 || dtList[i].ActiveAccount > maxActive {
			maxActive = dtList[i].ActiveAccount
			maxTime = dtList[i].Time
		}
		rets.Time[i] = dtList[i].Time
	}
	rets.MinActive = minActive
	rets.MinTime = minTime
	rets.MaxActive = maxActive
	rets.MaxTime = maxTime

	return rets, nil
}

func GetDailyActiveAccountChangeChart(search any, startTime, endTime int64) (*DailyActiveChartResponse, error) {
	if startTime != 0 && endTime != 0 {
		rets, err := getDailyActiveReportChart(startTime, endTime)
		if err != nil {
			return nil, err
		}
		return &rets, nil
	} else {
		if search != nil {
			searchMap := map[string]string{
				"1_month": "one_month",
				"3_month": "three_month",
				"1_year":  "one_year",
				"all":     "all",
			}
			switch reflect.TypeOf(search).String() {
			case "string":
				switch search.(string) {
				case "1_month", "3_month", "1_year", "all":
					return getActiveReportFromRedis(searchMap[search.(string)])
				default:
					return nil, errors.New("unknown request params")
				}
			default:
				log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Daily Active Account Change Chart Failed")
				return nil, errors.New("request params invalid")
			}
		} else {
			return nil, errors.New("request params invalid")
		}
	}
}

func GetDailyActiveReport(page, limit int, startTime, endTime int64) (GeneralResponse, error) {
	var (
		ret   []DailyActiveReport
		rets  GeneralResponse
		dt    DailyActiveReport
		total int64
	)
	endTime = time.Unix(endTime, 0).AddDate(0, 0, 1).Unix()
	err := GetDB(nil).Table(dt.TableName()).Where("time >= ? AND time < ?", startTime, endTime).Count(&total).Error
	if err != nil {
		return rets, err
	}
	rets.Total = total
	rets.Page = page
	rets.Limit = limit

	err = GetDB(nil).Order("id asc").Where("time >= ? AND time < ?", startTime, endTime).Offset((page - 1) * limit).Limit(limit).Find(&ret).Error
	if err != nil {
		return rets, err
	}
	rets.List = ret
	return rets, nil

}
