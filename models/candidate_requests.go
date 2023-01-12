/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"errors"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/smart"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"reflect"
)

var NodeReady bool
var PledgeAmount int64 = 1e18

type CandidateNodeRequests struct {
	Id           int64  `gorm:"primary_key;not_null"`
	TcpAddress   string `gorm:"not_null"`
	ApiAddress   string `gorm:"not_null"`
	NodePubKey   string `gorm:"not_null"`
	DateCreated  int64  `gorm:"not_null"`
	Deleted      int    `gorm:"not_null"`
	DateDeleted  int64  `gorm:"not_null"`
	Website      string `gorm:"not_null"`
	ReplyCount   int64  `gorm:"not_null"`
	DateReply    int64  `gorm:"not_null"`
	EarnestTotal string `gorm:"not_null"`
	NodeName     string `gorm:"not_null"`
}

type nodeDetailInfo struct {
	Id           int64
	NodeName     string
	Website      string
	ApiAddress   string
	Address      string
	Vote         decimal.Decimal
	VoteTrend    int
	EarnestTotal decimal.Decimal
	Committee    bool
	NodePubKey   string
	Ranking      int64
	Packed       int64
	PackedRate   string
}

func (p *CandidateNodeRequests) TableName() string {
	return "1_candidate_node_requests"
}

func InitPledgeAmount() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	if NodeReady {
		pledgeAmount, err := sqldb.GetPledgeAmount()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("init Pledge Amount Failed")
			return
		}
		PledgeAmount = pledgeAmount
	}
}

func (p *CandidateNodeRequests) GetById(id int64) (bool, error) {
	return isFound(GetDB(nil).Where("id = ?", id).First(&p))
}

func (p *CandidateNodeRequests) GetPubKeyById(id int64) (bool, error) {
	return isFound(GetDB(nil).Select("node_pub_key").Where("id = ?", id).First(&p))
}

func (p *CandidateNodeRequests) GetAll() ([]CandidateNodeRequests, error) {
	var list []CandidateNodeRequests
	err := GetDB(nil).Where("deleted = 0").Find(&list).Error
	return list, err
}

func (p *CandidateNodeRequests) GetNodeMap() (*NodeMapResponse, error) {
	var (
		rets         NodeMapResponse
		eligibleNode int64
	)
	err := GetDB(nil).Table(p.TableName()).Where("deleted = 0 AND earnest_total >= ?", PledgeAmount).Count(&eligibleNode).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Node Map Eligible Node Failed")
		return nil, err
	}

	err = GetDB(nil).Table(p.TableName()).Where("deleted = 0 AND earnest_total < ?", PledgeAmount).Count(&rets.CandidateTotal).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Node Map Candidate Total Failed")
		return nil, err
	}

	if eligibleNode > 101 {
		rets.HonorTotal = 101
		rets.CandidateTotal += eligibleNode - 101
	} else {
		rets.HonorTotal = eligibleNode
	}

	type nodeMap struct {
		Positioning
		EarnestTotal decimal.Decimal `json:"earnest_total"`
	}
	var list1 []nodeMap
	var list2 []nodeMap

	err = GetDB(nil).Raw(`
SELECT coalesce(ns.latitude,0)lat,coalesce(ns.longitude,0)lng,cs.earnest_total,cs.node_name as val FROM (
	SELECT id,node_name,earnest_total,
		CASE WHEN coalesce(referendum_total,0)>0 THEN
			round(coalesce(referendum_total,0) / 1e12,12)
		ELSE
			0
		END
		as vote,date_updated_referendum
	FROM "1_candidate_node_requests" WHERE deleted = 0 AND earnest_total >= ?
)AS cs
LEFT JOIN (
	SELECT latitude,longitude,value FROM honor_node_info
)AS ns ON(cast(ns.value->>'id' as numeric) = cs.id AND cast(ns.value->>'consensus_mode' as numeric) = 2)
ORDER BY vote desc,date_updated_referendum asc
`, PledgeAmount).Find(&list1).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Node Map List1 Failed")
		return nil, err
	}

	err = GetDB(nil).Raw(`
SELECT coalesce(ns.latitude,0)lat,coalesce(ns.longitude,0)lng,cs.earnest_total,cs.node_name as val FROM (
	SELECT id,earnest_total,node_name,
		CASE WHEN coalesce(referendum_total,0)>0 THEN
			round(coalesce(referendum_total,0) / 1e12,12)
		ELSE
			0
		END
		as vote,date_updated_referendum
	FROM "1_candidate_node_requests" WHERE deleted = 0 AND earnest_total < ?
)AS cs
LEFT JOIN (
	SELECT latitude,longitude,value FROM honor_node_info
)AS ns ON(cast(ns.value->>'id' as numeric) = cs.id AND cast(ns.value->>'consensus_mode' as numeric) = 2)
ORDER BY vote desc,date_updated_referendum asc
`, PledgeAmount).Find(&list2).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Node Map List1 Failed")
		return nil, err
	}

	eligibleDecl := decimal.NewFromInt(PledgeAmount)
	for i := 0; i < len(list1); i++ {
		if list1[i].EarnestTotal.GreaterThanOrEqual(eligibleDecl) {
			if i < 101 {
				list1[i].Positioning.Val = list1[i].Positioning.Val + "(Honor Node)"
				rets.NodeList = append(rets.NodeList, list1[i].Positioning)
			} else {
				list1[i].Positioning.Val = list1[i].Positioning.Val + "(Candidate Node)"
				rets.NodeList = append(rets.NodeList, list1[i].Positioning)
			}
		} else {
			list1[i].Positioning.Val = list1[i].Positioning.Val + "(Candidate Node)"
			rets.NodeList = append(rets.NodeList, list1[i].Positioning)
		}
	}
	for _, value := range list2 {
		value.Positioning.Val = value.Positioning.Val + "(Candidate Node)"
		rets.NodeList = append(rets.NodeList, value.Positioning)
	}

	return &rets, nil

}

func (p *nodeDetailInfo) getNodeDetailInfo(voteTotal decimal.Decimal, account string) (NodeListResponse, error) {
	var rets NodeListResponse
	zeroDec := decimal.New(0, 0)

	type committee struct {
		FrontCommittee bool
	}
	var rlt committee
	f, _ := isFound(GetDB(nil).Raw(`
	SELECT CASE WHEN (SELECT control_mode FROM "1_ecosystems" WHERE id = 1) = 2 THEN
		CASE WHEN (SELECT count(1) FROM "1_votings_participants" WHERE
				voting_id = (SELECT id FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%%voting_for_control_mode_template%%' AND ecosystem = 1 ORDER BY id DESC LIMIT 1)
				AND member->>'account'=?) > 0 THEN
			TRUE
		ELSE
			FALSE
		END
	ELSE
		FALSE
	END AS front_committee
	`, account).Take(&rlt))
	if f && rlt.FrontCommittee {
		rets.FrontCommittee = true
	}
	rets.Committee = p.Committee
	rets.ApiAddress = p.ApiAddress
	rets.Name = p.NodeName
	rets.Vote = p.Vote.String()
	voteDec := p.Vote.Mul(decimal.NewFromInt(100))
	if voteDec.GreaterThan(zeroDec) && voteTotal.GreaterThan(zeroDec) {
		rets.VoteRate = voteDec.DivRound(voteTotal, 2).String()
	} else {
		rets.VoteRate = zeroDec.String()
	}
	rets.VoteTrend = p.VoteTrend
	rets.IconUrl = getIconNationalFlag(getCountry(p.Address))
	rets.Ranking = p.Ranking
	rets.Website = p.Website
	rets.Address = p.Address
	rets.Staking = p.EarnestTotal.String()
	rets.Id = p.Id
	rets.Packed = p.Packed
	rets.PackedRate = p.PackedRate
	return rets, nil
}

func NodeListSearch(page, limit int, order string) (*GeneralResponse, error) {
	var (
		list []NodeListResponse
		rets GeneralResponse
	)
	if order == "" {
		order = "vote desc,date_updated_referendum asc"
	}

	var info []nodeDetailInfo

	rets.Page = page
	rets.Limit = limit

	err := GetDB(nil).Raw(`
SELECT cs.node_name,cs.id,cs.website,cs.api_address,hr.address,cs.vote,RANK() OVER (ORDER BY vote DESC,date_updated_referendum ASC) AS ranking,cs.packed,
	CASE WHEN cs.packed > 0 THEN
		round(cs.packed*100 / cast( (SELECT max(id) FROM block_chain)  as numeric),2) 
	ELSE
		0
	END packed_rate,
	CASE WHEN cs.vote > cast(coalesce(hr.value->>'vote','0') as numeric) THEN
		1
	WHEN cs.vote < cast(coalesce(hr.value->>'vote','0') as numeric) THEN
		2
	ELSE
		3
	END vote_trend,cs.earnest_total,
	CASE WHEN cs.earnest_total >= ? AND row_number() OVER (ORDER BY vote DESC) < 102 THEN 
		TRUE
	ELSE
		FALSE
	END committee,cs.node_pub_key
 FROM (
	SELECT id,api_address,(SELECT count(1) FROM block_chain WHERE node_position = cs.id AND consensus_mode = 2)packed,website,node_name,earnest_total,node_pub_key,
		CASE WHEN coalesce(referendum_total,0)>0 THEN
			round(coalesce(referendum_total,0) / 1e12,12)
		ELSE
			0
		END
		as vote,date_updated_referendum
	FROM "1_candidate_node_requests" AS cs WHERE deleted = 0
) AS cs
LEFT JOIN(
	SELECT value,address FROM honor_node_info AS he
)AS hr ON (cs.id = CAST(hr.value->>'id' AS numeric) AND CAST(hr.value->>'consensus_mode' AS numeric) = 2)
ORDER BY ? OFFSET ? LIMIT ?
`, PledgeAmount, order, (page-1)*limit, limit).Find(&info).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Node List Search Failed")
		return &rets, errors.New("candidate requests pubkey invalid")
	}

	type totalInfo struct {
		Total     int64
		VoteTotal decimal.Decimal
	}
	var ti totalInfo
	err = GetDB(nil).Raw(`
SELECT (SELECT count(1) FROM "1_candidate_node_requests" WHERE deleted = 0) total,case WHEN coalesce(sum(earnest),0) > 0 THEN 
	round(coalesce(sum(earnest),0) / 1e12,12)
ELSE
	0
END as vote_total FROM "1_candidate_node_decisions" WHERE decision_type = 1 AND decision <> 3
AND request_id IN (SELECT id FROM "1_candidate_node_requests" WHERE deleted = 0)
`).Take(&ti).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Node List Search Total Info Failed")
		return &rets, err
	}
	rets.Total = ti.Total
	for i := 0; i < len(info); i++ {
		account := converter.IDToAddress(smart.PubToID(info[i].NodePubKey))
		if account == "invalid" {
			log.WithFields(log.Fields{"pub_key": info[i].NodePubKey}).Error("Node List Search Pub Key Failed")
			return &rets, errors.New("candidate requests pub_key invalid")
		}
		rts, err := info[i].getNodeDetailInfo(ti.VoteTotal, account)
		if err != nil {
			log.WithFields(log.Fields{"err": err, "info": info[i]}).Error("Get Node Detail Info Failed")
			return &rets, err
		}

		list = append(list, rts)

	}
	rets.List = list
	return &rets, nil
}

func NodeDetail(id int64) (NodeDetailResponse, error) {
	var (
		nodeInfo  nodeDetailInfo
		rets      NodeDetailResponse
		voteTotal decimal.Decimal
	)
	err := GetDB(nil).Raw(`
SELECT case WHEN coalesce(sum(earnest),0) > 0 THEN 
	round(coalesce(sum(earnest),0) / 1e12,12)
ELSE
	0
END as vote_total FROM "1_candidate_node_decisions" WHERE decision_type = 1 AND decision <> 3
AND request_id IN (SELECT id FROM "1_candidate_node_requests" WHERE deleted = 0)
`).Take(&voteTotal).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err, "node id": id}).Error("Get Node Detail Vote Total Failed")
		return rets, err
	}

	err = GetDB(nil).Raw(`
SELECT cs.node_name,cs.id,cs.website,cs.api_address,hr.address,cs.vote,cs.packed,
	CASE WHEN cs.packed > 0 THEN
		round(cs.packed*100 / cast( (SELECT max(id) FROM block_chain)  as numeric),2) 
	ELSE
		0
	END packed_rate,
(SELECT CASE WHEN (SELECT count(1) FROM "1_candidate_node_decisions" WHERE decision_type = 1 AND decision <> 3 AND request_id = cs.id)  > 0 THEN
	(SELECT ct.ranking
	 FROM(
		SELECT coalesce(RANK() OVER (ORDER BY coalesce(sum(earnest),0) DESC),1)AS ranking,request_id FROM "1_candidate_node_decisions" WHERE decision_type = 1 AND decision <> 3  GROUP BY request_id
	 )AS ct WHERE ct.request_id = cs.id)
ELSE
	(SELECT coalesce(count(1),0) + 1 FROM (SELECT request_id FROM "1_candidate_node_decisions" WHERE decision_type = 1 AND decision <> 3  GROUP BY request_id)AS te)
END) ranking,
CASE WHEN cs.vote > cast(coalesce(hr.value->>'vote','0') as numeric) THEN
	1
WHEN cs.vote < cast(coalesce(hr.value->>'vote','0') as numeric) THEN
	2
ELSE
	3
END vote_trend,cs.earnest_total,CASE WHEN cs.earnest_total >= ? AND row_number() OVER (ORDER BY vote DESC,date_updated_referendum asc) < 102 THEN 
	TRUE
ELSE
	FALSE
END committee,cs.node_pub_key
 FROM (
	SELECT id,api_address,website,node_name,earnest_total,(SELECT count(1) FROM block_chain WHERE node_position = cs.id AND consensus_mode = 2)packed,node_pub_key,
		CASE WHEN coalesce(referendum_total,0)>0 THEN
			round(coalesce(referendum_total,0) / 1e12,12)
		ELSE
			0
		END
		as vote,date_updated_referendum
	FROM "1_candidate_node_requests" AS cs WHERE deleted = 0 AND id = ?
) AS cs
LEFT JOIN(
	SELECT value,address FROM honor_node_info AS he
)AS hr ON (cs.id = CAST(hr.value->>'id' AS numeric) AND CAST(hr.value->>'consensus_mode' AS numeric) = 2)
`, PledgeAmount, id).Take(&nodeInfo).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err, "node id": id}).Error("Get Node Detail Failed")
		return rets, err
	}
	account := converter.IDToAddress(smart.PubToID(nodeInfo.NodePubKey))
	if account == "invalid" {
		log.WithFields(log.Fields{"pub_key": nodeInfo.NodePubKey}).Error("Get Node Detail Pub Key Failed")
		return rets, errors.New("candidate requests pub_key invalid")
	}

	info, err := nodeInfo.getNodeDetailInfo(voteTotal, account)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "node id": id}).Error("Get Node Detail Info Failed")
		return rets, err
	}
	rets.NodeListResponse = info
	zeroDec := decimal.New(0, 0)
	eligibleDecl := decimal.NewFromInt(PledgeAmount)
	stakeDec, _ := decimal.NewFromString(info.Staking)
	if stakeDec.GreaterThan(zeroDec) {
		rets.StakeRate = stakeDec.Mul(decimal.NewFromInt(100)).DivRound(eligibleDecl, 2).String()
	} else {
		rets.StakeRate = "0"
	}
	rets.Account = account

	return rets, nil
}

func GetNodeBlockList(search any, page, limit int, order string) (GeneralResponse, error) {
	var (
		list   []NodeBlockListResponse
		rets   GeneralResponse
		err    error
		id     int64
		bk     Block
		bkList []Block
	)

	if order == "" {
		order = "id desc"
	}
	rets.Page = page
	rets.Limit = limit

	switch reflect.TypeOf(search).String() {
	case "json.Number":
		id, err = search.(json.Number).Int64()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("Get Node Block List Json Number Failed")
			return rets, err
		}
	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Node Block List Search Failed")
		return rets, errors.New("request params invalid")
	}
	if id <= 0 {
		return rets, errors.New("unknown node id 0")
	}
	err = GetDB(nil).Table(bk.TableName()).Where("node_position = ? AND consensus_mode = 2", id).Count(&rets.Total).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("Get Node Block List Total Failed")
		return rets, err
	}
	if rets.Total > 0 {
		err = GetDB(nil).Select("id,tx,time").Where("node_position = ? AND consensus_mode = 2", id).
			Offset((page - 1) * limit).Limit(limit).Order(order).Find(&bkList).Error
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("Get Node Block Block List Failed")
			return rets, err
		}

		for _, value := range bkList {
			var rts NodeBlockListResponse
			rts.BlockId = value.ID

			type txGasFee struct {
				Amount    string
				Ecosystem int64
			}
			var gasFee []txGasFee
			err = GetDB(nil).Raw(`
				SELECT v1.ecosystem,sum(v1.amount)as amount FROM(
					SELECT ecosystem,sum(amount)amount FROM "1_history" WHERE block_id = ? AND type IN(1,2) GROUP BY ecosystem
					UNION
					SELECT ecosystem,sum(amount)amount FROM "spent_info_history" WHERE block = ? AND type IN(3,4) GROUP BY ecosystem
				)AS v1
				GROUP BY ecosystem
			`, rts.BlockId, rts.BlockId).Find(&gasFee).Error
			if err != nil {
				log.WithFields(log.Fields{"err": err, "block_id": value.ID}).Warn("Get Node Block Tx List Failed")
				return rets, err
			}
			gasFeeCursor := 1
			rts.EcoNumber = len(gasFee)
			rts.Time = value.Time
			rts.Tx = value.Tx
			for _, vue := range gasFee {
				tokenSymbol := Tokens.Get(vue.Ecosystem)
				digits := EcoDigits.GetInt64(vue.Ecosystem, 0)
				if vue.Ecosystem == 1 {
					rts.GasFee1.Amount = vue.Amount
					rts.GasFee1.TokenSymbol = tokenSymbol
					rts.GasFee1.Digits = digits
				} else {
					gasFeeCursor += 1
					if gasFeeCursor > 5 {
						break
					}
					switch gasFeeCursor {
					case 2:
						rts.GasFee2.Amount = vue.Amount
						rts.GasFee2.TokenSymbol = tokenSymbol
						rts.GasFee2.Digits = digits
					case 3:
						rts.GasFee3.Amount = vue.Amount
						rts.GasFee3.TokenSymbol = tokenSymbol
						rts.GasFee3.Digits = digits
					case 4:
						rts.GasFee4.Amount = vue.Amount
						rts.GasFee4.TokenSymbol = tokenSymbol
						rts.GasFee4.Digits = digits
					case 5:
						rts.GasFee5.Amount = vue.Amount
						rts.GasFee5.TokenSymbol = tokenSymbol
						rts.GasFee5.Digits = digits
					}
				}
			}

			list = append(list, rts)
		}
		rets.List = list
	}

	return rets, nil
}

func CandidateTableExist() bool {
	var p CandidateNodeRequests
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}
