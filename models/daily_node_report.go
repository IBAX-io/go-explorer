/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type DailyNodeReport struct {
	ID            int64  `gorm:"primary_key;not null" json:"id"`
	Time          int64  `gorm:"column:time;not null" json:"time"`
	HonorNode     string `gorm:"not null;type:jsonb" json:"honor_node"`
	CandidateNode string `gorm:"not null;type:jsonb" json:"candidate_node"`
}

type nodeDetail struct {
	Id int64 `json:"id"`
}

func (p *DailyNodeReport) TableName() string {
	return "daily_node_report"
}

func (p *DailyNodeReport) CreateTable() (err error) {
	err = nil
	if !HasTableOrView(nil, p.TableName()) {
		if err = GetDB(nil).Migrator().CreateTable(p); err != nil {
			return err
		}
	}
	return err
}

func (p *DailyNodeReport) GetLast() (bool, error) {
	return isFound(GetDB(nil).Last(p))
}

func (dt *DailyNodeReport) Insert() error {
	if dt.CandidateNode == "" && dt.HonorNode == "" {
		return nil
	}
	if err := GetDB(nil).Model(&DailyNodeReport{}).Create(&dt).Error; err != nil {
		return err
	}
	return nil
}

func InsertDailyNodeReport() {
	var (
		dt DailyNodeReport
	)
	if !NodeReady {
		return
	}
	now := time.Now()
	f, err := dt.GetLast()
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("Get Last Daily Node Report Failed")
		return
	}
	if f {
		lastTime := time.Unix(dt.Time, 0)
		diffDay := int(now.Sub(lastTime).Hours() / 24)
		if diffDay < 1 {
			return
		}
	}

	var rt DailyNodeReport
	err = rt.getDailyNodeReport()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get Daily Node Report Report Failed")
		return
	}

	err = rt.Insert()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Insert Daily Node Report Insert Failed")
		return
	}
}

func (p *DailyNodeReport) getDailyNodeReport() error {
	type nodeList struct {
		Id        int64
		Committee bool
	}
	var list []nodeList
	err := GetDB(nil).Raw(`
SELECT cs.id,CASE WHEN cs.earnest_total >= ? AND row_number() OVER (ORDER BY vote DESC) < 102 THEN 
	TRUE
ELSE
	FALSE
END committee
 FROM (
	SELECT id,earnest_total,(
		SELECT case WHEN coalesce(sum(earnest),0) > 0 THEN 
			round(coalesce(sum(earnest),0) / 1e12,0)
		ELSE
			0
		END as vote FROM "1_candidate_node_decisions" WHERE decision_type = 1 AND decision <> 3 AND request_id = cs.id
	) FROM "1_candidate_node_requests" AS cs WHERE deleted = 0
) AS cs
`, PledgeAmount).Find(&list).Error
	if err != nil {
		return err
	}
	p.Time = GetNowTimeUnix()
	var honorList []nodeDetail
	var candidateList []nodeDetail
	for i := 0; i < len(list); i++ {
		if list[i].Committee {
			honorList = append(honorList, nodeDetail{Id: list[i].Id})
		} else {
			candidateList = append(candidateList, nodeDetail{Id: list[i].Id})
		}
	}
	if len(honorList) == 0 && len(candidateList) == 0 {
		return nil
	}
	hor, err := json.Marshal(honorList)
	if err != nil {
		return err
	}
	cte, err := json.Marshal(candidateList)
	if err != nil {
		return err
	}
	p.HonorNode = string(hor)
	p.CandidateNode = string(cte)

	return nil
}

func GetNodeChangeChart() (NodeChangeResponse, error) {
	var (
		rets NodeChangeResponse
	)
	type nodeChange struct {
		Days         string `json:"days"`
		HonorLen     int64
		CandidateLen string
		Num          string
	}
	if !NodeReady {
		return rets, nil
	}
	var list []nodeChange
	err := GetDB(nil).Raw(`
SELECT c1.days,c1.honor_len,c1.candidate_len,coalesce(c2.num,0)num FROM (
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
		END candidate_len FROM daily_node_report order by time asc
)AS c1
LEFT JOIN(
	SELECT to_char(to_timestamp(rt.time),'yyyy-MM-dd') AS days,count(1) num
	FROM(
		SELECT id,coalesce((SELECT time FROM daily_node_report WHERE honor_node @> ('[{"id":'||c1.id||'}]')::jsonb ORDER BY time asc LIMIT 1),0)AS time 
		FROM "1_candidate_node_requests" AS c1 WHERE deleted = 0
	)AS rt WHERE rt.time > 0 GROUP BY days
)AS c2 ON(c2.days = c1.days) order by days asc
`).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Node Change Chart Failed")
		return rets, err
	}
	for i := 0; i < len(list); i++ {
		rets.Time = append(rets.Time, list[i].Days)
		rets.HonorNode = append(rets.HonorNode, strconv.FormatInt(list[i].HonorLen, 10))
		rets.CandidateNode = append(rets.CandidateNode, list[i].CandidateLen)
		rets.NewHonor = append(rets.NewHonor, list[i].Num)
	}

	return rets, nil

}
