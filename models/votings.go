/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/smart"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

const DaoVoteChart = "dao_vote_chart"

var VotingReady bool

type Voting struct {
	Id          int64  `gorm:"primary_key;not null"`
	Creator     string `gorm:"column:lock;type:jsonb"`
	DateEnded   int64  `gorm:"column:date_ended"`
	DateStarted int64  `gorm:"column:date_started"`
	Deleted     int64  `gorm:"column:deleted"`
	Ecosystem   int64  `gorm:"column:ecosystem"`
	Flags       string `gorm:"column:flags;type:jsonb"`
	Optional    string `gorm:"column:optional;type:jsonb"`
	Progress    string `gorm:"column:progress;type:jsonb"`
	Status      int64  `gorm:"column:status"`
	Voting      string `gorm:"column:voting;type:jsonb"`
}

type CreatorInfo struct {
	Account    string `json:"account"`
	MemberName string `json:"member_name"`
}

type FlagsInfo struct {
	Success  string `json:"success"`
	Decision string `json:"decision"`
	Notifics string `json:"notifics"`
	FullData string `json:"full_data"`
}

type ProgressInfo struct {
	NumberVoters       int    `json:"number_voters"`
	PercentVoters      int    `json:"percent_voters"`
	PercentSuccess     int    `json:"percent_success"`
	NumberParticipants string `json:"number_participants"`
}

type OptionalInfo struct {
	ContractAccept string `json:"contract_accept"`
	ContractReject string `json:"contract_reject"`
}

type VotingInfo struct {
	Name             string `json:"name"`
	Type             int    `json:"type"`
	Quorum           int    `json:"quorum"`
	Rating           int    `json:"rating"`
	Volume           int    `json:"volume"`
	RoleId           string `json:"role_id"`
	Description      string `json:"description"`
	TypeDecision     int    `json:"type_decision"`
	CountTypeVoters  int    `json:"count_type_voters"`
	TypeParticipants int    `json:"type_participants"`
}

func (p *Voting) TableName() string {
	return "1_votings"
}

func (p *Voting) GetLast() (bool, error) {
	return isFound(GetDB(nil).Where("deleted = 0").Last(p))
}

func GetLatestDaoVoting() (VotingResponse, error) {
	var (
		rets  VotingResponse
		total int64
		p     Voting
	)
	err := GetDB(nil).Table(p.TableName()).Where("deleted = 0 AND voting->>'name' like '%voting_for_control_mode_template%'").Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get Last Dao Voting Total Failed")
		return rets, err
	}
	rets.Total = total
	if total > 0 {
		err = GetDB(nil).Raw(`
SELECT vs.title,vs.id,vs.created,vs.voted_rate,sub.result_rate,sub.rejected_rate FROM(
	select id,coalesce(voting->>'name','') AS title,date_started AS created,
	CAST(coalesce(progress->>'percent_voters','0') as numeric)as voted_rate
	from "1_votings" WHERE deleted = 0 AND voting->>'name' like '%voting_for_control_mode_template%' ORDER BY id DESC LIMIT 1
)AS vs
LEFT JOIN (
	SELECT round(CAST(coalesce(results->>'percent_accepted','0') as numeric),2) AS result_rate,
		round(CAST(coalesce(results->>'percent_rejected','0') as numeric),2) AS rejected_rate,voting_id FROM "1_votings_subject"
)AS sub ON(sub.voting_id = vs.id)
`).Take(&rets).Error
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Get Last Dao Voting Failed")
			return rets, err
		}
	}
	return rets, nil
}

func GetDaoVoteList(search any, page, limit int) (GeneralResponse, error) {
	var (
		can     CandidateNodeRequests
		list    []NodeVoteResponse
		rets    GeneralResponse
		total   CountInt64
		err     error
		nodeId  int64
		account string
		vote    Voting
	)

	rets.Page = page
	rets.Limit = limit

	if search != nil {
		switch reflect.TypeOf(search).String() {
		case "json.Number":
			nodeId, err = search.(json.Number).Int64()
			if err != nil {
				log.WithFields(log.Fields{"err": err}).Warn("Get Node Dao Vote List Json Number Failed")
				return rets, err
			}
		default:
			log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Node Dao Vote List Search Failed")
			return rets, errors.New("request params invalid")
		}

		f, err := can.GetPubKeyById(nodeId)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Get Node Dao Voting Pub Key Failed")
			return rets, err
		}
		if !f {
			return rets, errors.New("get Node Pub Key Failed")
		}

		account = converter.IDToAddress(smart.PubToID(can.NodePubKey))
		if account == "invalid" {
			log.WithFields(log.Fields{"pub_key": can.NodePubKey}).Error("Node Pub Key Invalid")
			return rets, errors.New("candidate requests pub_key invalid")
		}

		err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT count(1) FROM "1_votings_participants" WHERE 
	voting_id IN(SELECT id FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%%voting_for_control_mode_template%%' AND ecosystem = 1) 
	AND member->>'account'='%s'
`, account)).Take(&total.Count).Error
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Get Node Dao Vote List Total Failed")
			return rets, err
		}
		rets.Total = total.Count

		err = GetDB(nil).Raw(`
SELECT vs.title,vs.created,vs.voted_rate,sub.result_rate,sub.rejected_rate,vs.id,vs.result,coalesce((
	SELECT CASE WHEN decision = 1 THEN
		1
	WHEN decision = -1 THEN
		2
	ELSE
		3
	END 
	 FROM "1_votings_participants" WHERE member->>'account'=? AND voting_id = vs.id),0)AS status FROM(
	select id,coalesce(voting->>'name','') AS title,date_started AS created,
	CAST(coalesce(progress->>'percent_voters','0') as numeric)as voted_rate,
	CASE WHEN cast(flags->>'success' AS numeric) = 1 THEN
		case WHEN cast(flags->>'decision' AS numeric) = 1 THEN
		 1
		ELSE
		 2
		END
	ELSE
		3
	END result
	FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%voting_for_control_mode_template%' AND id IN(
		SELECT voting_id FROM "1_votings_participants" WHERE voting_id IN(
			SELECT id FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%voting_for_control_mode_template%' AND ecosystem = 1
		) AND member->>'account'=?
	)
)AS vs
LEFT JOIN (
	SELECT round(CAST(coalesce(results->>'percent_accepted','0') as numeric),2) AS result_rate,
			round(CAST(coalesce(results->>'percent_rejected','0') as numeric),2) AS rejected_rate,voting_id
	FROM "1_votings_subject"
)AS sub ON(sub.voting_id = vs.id)
ORDER BY id desc OFFSET ? LIMIT ?
`, account, account, (page-1)*limit, limit).Find(&list).Error
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Get Node Dao Vote List Failed")
			return rets, err
		}
	} else {
		err = GetDB(nil).Table(vote.TableName()).Where("voting->>'name' like '%voting_for_control_mode_template%'").Count(&rets.Total).Error
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Get Dao Vote List Total Failed")
			return rets, err
		}

		err = GetDB(nil).Raw(`
SELECT vs.title,vs.created,vs.voted_rate,sub.result_rate,sub.rejected_rate,vs.id,vs.result FROM(
					SELECT id,coalesce(voting->>'name','') AS title,date_started AS created,
					CAST(coalesce(progress->>'percent_voters','0') as numeric)as voted_rate,
					CASE WHEN cast(flags->>'success' AS numeric) = 1 THEN
									case WHEN cast(flags->>'decision' AS numeric) = 1 THEN
									 1
									ELSE
									 2
									END
					ELSE
									3
					END result
					FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%voting_for_control_mode_template%'
				)AS vs
LEFT JOIN (
        SELECT round(CAST(coalesce(results->>'percent_accepted','0') as numeric),2) AS result_rate,
               round(CAST(coalesce(results->>'percent_rejected','0') as numeric),2) AS rejected_rate,voting_id
        FROM "1_votings_subject"
)AS sub ON(sub.voting_id = vs.id)
ORDER BY id desc OFFSET ? LIMIT ?
`, (page-1)*limit, limit).Find(&list).Error
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Get Dao Vote List Failed")
			return rets, err
		}
	}
	rets.List = list

	return rets, nil
}

func getDaoVoteChart() (DaoVoteChartResponse, error) {
	var (
		rets DaoVoteChartResponse
		err  error
	)
	if !VotingReady {
		return rets, nil
	}

	err = GetDB(nil).Raw(`
SELECT v1.date_ended AS time,v2.agree,v2.rejected,v2.did_not_vote FROM(
	SELECT id,(SELECT max(decision_date) FROM "1_votings_participants" WHERE voting_id = vs.id)AS date_ended
		FROM "1_votings" vs WHERE deleted = 0 AND voting->>'name' like '%%voting_for_control_mode_template%%' AND cast(progress->>'percent_success' as numeric) = 100
)AS v1 
LEFT JOIN(
	SELECT sum(CASE WHEN decision = 1 THEN 1 ELSE 0 END)agree,
		 sum(CASE WHEN decision = -1 THEN 1 ELSE 0 END)rejected,
		 sum(CASE WHEN decision = 0 THEN 1 ELSE 0 END)did_not_vote,voting_id FROM "1_votings_participants" GROUP BY voting_id
)as v2 ON(v2.voting_id = v1.id)
ORDER BY date_ended ASC
`).Find(&rets.List).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get Dao Chart list To Redis Failed")
		return rets, err
	}

	type decision struct {
		FlagsDecision int
		Num           int64
	}
	var decList []decision
	err = GetDB(nil).Raw(`
SELECT CAST(flags->>'decision' AS numeric)flags_decision,count(1) num
	FROM "1_votings" vs WHERE deleted = 0 AND voting->>'name' like '%%voting_for_control_mode_template%%' GROUP BY flags_decision
`).Find(&decList).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Get Dao Chart decision list To Redis Failed")
		return rets, err
	}
	for key, value := range rets.List {
		rets.List[key].Total = value.Agree + value.DidNotVote + value.Rejected
	}
	for _, value := range decList {
		number := value.Num
		switch value.FlagsDecision {
		case -2:
			rets.NotEnoughVotes = number
		case -1:
			rets.Rejected = number
		case 1:
			rets.Accept = number
		default: //0
			rets.No = number
		}
	}

	return rets, nil
}

func GetDaoVoteDetail(search any, language string) (DaoVoteDetailResponse, error) {
	var (
		rets   DaoVoteDetailResponse
		voteId int64
		err    error
	)

	switch reflect.TypeOf(search).String() {
	case "json.Number":
		voteId, err = search.(json.Number).Int64()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("Get Dao Vote Detail Json Number Failed")
			return rets, err
		}
	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Dao Vote Detail Search Failed")
		return rets, errors.New("request params invalid")
	}

	type voteDetail struct {
		Id                   int64
		Title                string
		FlagsDecision        int
		AppId                int64
		Creator              string
		Member               string
		Description          string
		FullData             int
		TypeDecision         int
		DateStarted          int64
		DateEnded            int64
		Quorum               int
		Volume               int
		CountTypeVoters      int
		Status               int
		Type                 int
		VotedRate            int
		ResultRate           float64
		RejectedRate         float64
		Progress             int
		ContractAccept       string
		ContractAcceptParams string
		ContractReject       string
		ContractRejectParams string
		NumberParticipants   int64
		Agree                int64
		Rejected             int64
		DidNotVote           int64
	}
	var info voteDetail

	err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT * FROM(
	SELECT id,voting->>'name' title,
		CAST(voting->>'type_decision' as numeric) type_decision,
		(SELECT id FROM "1_applications" WHERE ecosystem = 1 AND name = 'Basic')app_id,
		creator->>'account' creator,
		CAST(coalesce(progress->>'percent_voters','0') as numeric)as voted_rate,
		array_to_string(array(SELECT member->>'account' FROM "1_votings_participants" ps WHERE vs.id = ps.voting_id),',') member,
		voting->>'description' description,CAST(flags->>'full_data' AS numeric)full_data,CAST(flags->>'decision' AS numeric)flags_decision,
		date_started,date_ended,CAST(voting->>'quorum' AS numeric)quorum,CAST(voting->>'volume' AS numeric)volume,
		CAST(voting->>'count_type_voters' AS numeric)count_type_voters,status,CAST(voting->>'type' AS numeric)AS type,
		CAST(progress->>'percent_success' AS numeric)progress,CAST(progress->>'number_participants' AS numeric)number_participants
	FROM "1_votings" vs WHERE id = %d AND deleted = 0 AND voting->>'name' like '%%voting_for_control_mode_template%%'
)AS v1 
LEFT JOIN(
	SELECT subject->>'contract_accept' contract_accept,
		round(CAST(coalesce(results->>'percent_accepted','0') as numeric),2) AS result_rate,
		round(CAST(coalesce(results->>'percent_rejected','0') as numeric),2) AS rejected_rate,
		subject->>'contract_accept_params' contract_accept_params,
		subject->>'contract_reject' contract_reject,voting_id FROM "1_votings_subject"
)as v2 ON(v2.voting_id = v1.id)
LEFT JOIN(
	SELECT sum(CASE WHEN decision = 1 THEN 1 ELSE 0 END)agree,
		 sum(CASE WHEN decision = -1 THEN 1 ELSE 0 END)rejected,sum(CASE WHEN decision = 0 THEN 1 ELSE 0 END)did_not_vote,voting_id FROM "1_votings_participants" GROUP BY voting_id
)as v3 ON(v3.voting_id = v1.id)
`, voteId)).Take(&info).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("Get Dao Vote Detail Info Failed")
		return rets, err
	}

	decisions, err := GetAppValue(info.AppId, "type_voting_decisions", 1)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "name": "type_voting_decisions"}).Warn("Get Dao Vote Detail App Params Failed")
		return rets, err
	}
	decisionList := strings.Split(decisions, ",")
	if len(decisionList) >= info.TypeDecision {
		id, name := converter.ParseName(decisionList[info.TypeDecision-1])
		rets.TypeDecision = getLanguageValue(language, name, id)
		if err != nil {
			log.WithFields(log.Fields{"err": err, "name": name, "language": language, "ecosystem": id}).Warn("Get Dao Vote Detail Language Value Failed")
		}
	}

	rets.Id = info.Id
	rets.Title = info.Title
	//rets.DescriptionText = getLanguageValue(language, "description", 1)
	rets.Description = info.Description
	noText := getLanguageValue(language, "no", 1)
	yesText := getLanguageValue(language, "yes", 1)

	if info.ContractAccept != "" {
		rets.VotingSubject.ContractAccept = info.ContractAccept
		if info.ContractAcceptParams == "" || info.ContractAcceptParams == "{}" {
			rets.VotingSubject.Arguments = noText
		} else {
			rets.VotingSubject.Arguments = info.ContractAcceptParams
		}
	} else {
		rets.VotingSubject.ContractAccept = noText
	}
	if info.ContractReject != "" {
		rets.VotingSubject.ContractOfReject = info.ContractReject
		if info.ContractRejectParams == "" || info.ContractRejectParams == "{}" {
			rets.VotingSubject.Arguments = noText
		} else {
			rets.VotingSubject.Arguments = info.ContractRejectParams
		}
	} else {
		rets.VotingSubject.ContractOfReject = noText
	}

	voteType, err := GetAppValue(info.AppId, "type_voting", 1)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "name": "type_voting"}).Warn("Get Dao Vote Detail App Params Failed")
		return rets, err
	}
	typeList := strings.Split(voteType, ",")
	if len(typeList) >= info.Type {
		id, name := converter.ParseName(typeList[info.Type-1])
		rets.Voting.Type = getLanguageValue(language, name, id)
	}

	notFilled := getLanguageValue(language, "not_filled", 1)
	if info.Status == 1 && (info.FullData == 0 || info.NumberParticipants == 0) {
		rets.Voting.Status = notFilled
	} else {
		status, err := GetAppValue(info.AppId, "voting_statuses", 1)
		if err != nil {
			log.WithFields(log.Fields{"err": err, "name": "voting_statuses_classes"}).Warn("Get Dao Vote Detail App Params Failed")
			return rets, err
		}
		statusList := strings.Split(status, ",")
		if len(statusList) >= info.Status {
			id, name := converter.ParseName(statusList[info.Status-1])
			rets.Voting.Status = getLanguageValue(language, name, id)
		}
	}

	if info.CountTypeVoters == 1 {
		rets.Voting.VoteCountType = getLanguageValue(language, "number_votes", 1)
	} else {
		rets.Voting.VoteCountType = getLanguageValue(language, "percent_votes", 1)
	}
	rets.Voting.CountTypeVoters = info.CountTypeVoters

	if info.FullData == 1 {
		rets.Voting.Filled = yesText
	} else {
		rets.Voting.Filled = noText
	}

	switch info.FlagsDecision {
	case -2:
		rets.Voting.Decision = getLanguageValue(language, "not_enough_votes", 1)
	case -1:
		rets.Voting.Decision = getLanguageValue(language, "rejected", 1)
	case 1:
		rets.Voting.Decision = getLanguageValue(language, "accept", 1)
	default: //0
		rets.Voting.Decision = noText
	}
	rets.Voting.DateStart = info.DateStarted
	rets.Voting.DateEnd = info.DateEnded
	rets.Voting.Quorum = info.Quorum
	if info.CountTypeVoters != 1 && info.TypeDecision != 1 && info.TypeDecision != 2 {
		rets.Voting.Volume = info.Volume
	}

	if info.NumberParticipants > 0 {
		rets.Voting.Participants = info.NumberParticipants
	}
	rets.Voting.Creator = info.Creator
	rets.Voting.Member = strings.Split(info.Member, ",")
	rets.Agree = info.Agree
	rets.Rejected = info.Rejected
	rets.DidNotVote = info.DidNotVote
	rets.VotedRate = info.VotedRate
	rets.RejectedRate = info.RejectedRate
	rets.ResultRate = info.ResultRate
	rets.Created = info.DateStarted
	rets.Progress = info.Progress

	return rets, nil
}

func VotingTableExist() bool {
	var p Voting
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}

func GetNodeVotingHistory(search any, page, limit int, order string) (*GeneralResponse, error) {
	var (
		nodeId int64
		err    error
		rets   []NodeVoteHistoryResponse
		total  int64
		result GeneralResponse
	)
	if order == "" {
		order = "id DESC"
	} else {
		if !CheckSql(order) {
			return nil, errors.New("params invalid")
		}
	}

	type voteHistory struct {
		Id      int64  `json:"id"`
		Vote    int64  `json:"vote"`
		TxHash  string `json:"tx_hash"`
		Address int64  `json:"address"`
		Time    int64  `json:"time"`
		Events  int64  `json:"events"`
		Amount  string `json:"amount"`
	}
	var list []voteHistory
	switch reflect.TypeOf(search).String() {
	case "json.Number":
		nodeId, err = search.(json.Number).Int64()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("Get Node Voting History Json Number Failed")
			return nil, err
		}
	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Node Voting History Search Failed")
		return nil, errors.New("request params invalid")
	}
	err = GetDB(nil).Table("1_history").
		Where(fmt.Sprintf(
			`(type = 20 AND comment = 'Candidate Node Referendum #%d') OR (type = 22 AND comment = 'Candidate Node Withdraw Referendum #%d')`, nodeId, nodeId)).
		Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err, "node id": nodeId}).Warn("Get Node Voting History Total Failed")
		return nil, err
	}

	err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT encode(txhash, 'hex') AS tx_hash,id,created_at as time,
CASE WHEN sender_id = 0 THEN
	recipient_id
ELSE 
	sender_id
END AS address,type AS events,case WHEN coalesce(amount,0) > 0 THEN 
			round(coalesce(amount,0) / 1e12,0)
		ELSE
			0
		END as vote,amount FROM "1_history" WHERE (type = 20 AND comment = 'Candidate Node Referendum #%d') OR 
		(type = 22 AND comment = 'Candidate Node Withdraw Referendum #%d') 
ORDER BY %s OFFSET ? LIMIT ?
`, nodeId, nodeId, order), (page-1)*limit, limit).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err, "node id": nodeId}).Warn("Get Node Voting History Failed")
		return nil, err
	}

	for _, value := range list {
		var rt NodeVoteHistoryResponse
		rt.Id = value.Id
		rt.Vote = value.Vote
		rt.Address = converter.AddressToString(value.Address)
		rt.Amount = value.Amount
		rt.Time = MsToSeconds(value.Time)
		switch value.Events {
		case 20:
			rt.Events = 1
		case 22:
			rt.Events = 2
		}
		rt.TxHash = value.TxHash
		rets = append(rets, rt)
	}
	result.Page = page
	result.Limit = limit
	result.List = rets
	result.Total = total

	return &result, nil
}

func GetNodeStakingHistory(search any, page, limit int, order string) (*GeneralResponse, error) {
	var (
		nodeId int64
		err    error
		rets   []NodeStakingHistoryResponse
		total  int64
		result GeneralResponse
	)
	if order == "" {
		order = "id DESC"
	} else {
		if !CheckSql(order) {
			return nil, fmt.Errorf("params invalid")
		}
	}

	switch reflect.TypeOf(search).String() {
	case "json.Number":
		nodeId, err = search.(json.Number).Int64()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("Get Node Staking History Json Number Failed")
			return nil, err
		}
	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Node Staking History Search Failed")
		return nil, errors.New("request params invalid")
	}

	err = GetDB(nil).Table("1_history").
		Where(fmt.Sprintf(`(type = 18 AND comment = 'Candidate Node Earnest #%d') OR 
		(type = 19 AND comment = 'Candidate Node Substitute #%d') OR 
		(type = 21 AND comment = 'Candidate Node Withdraw Substitute #%d')`, nodeId, nodeId, nodeId)).
		Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err, "node id": nodeId}).Warn("Get Node Staking History Total Failed")
		return nil, err
	}

	type stakingHistory struct {
		Id      int64  `json:"id"`
		TxHash  string `json:"tx_hash"`
		Address int64  `json:"address"`
		Time    int64  `json:"time"`
		Events  int64  `json:"events"`
		Amount  string `json:"amount"`
	}
	var list []stakingHistory
	err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT encode(txhash, 'hex') AS tx_hash,id,created_at as time,
CASE WHEN sender_id = 0 THEN
	recipient_id
ELSE 
	sender_id
END AS address,type AS events,amount FROM "1_history" WHERE 
		(type = 18 AND comment = 'Candidate Node Earnest #%d') OR 
		(type = 19 AND comment = 'Candidate Node Substitute #%d') OR 
		(type = 21 AND comment = 'Candidate Node Withdraw Substitute #%d')
ORDER BY %s OFFSET ? LIMIT ?
`, nodeId, nodeId, nodeId, order), (page-1)*limit, limit).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err, "node id": nodeId}).Warn("Get Node Staking History Failed")
		return nil, err
	}

	for _, value := range list {
		var rt NodeStakingHistoryResponse
		rt.Id = value.Id
		rt.Address = converter.AddressToString(value.Address)
		rt.Amount = value.Amount
		rt.Time = MsToSeconds(value.Time)
		switch value.Events {
		case 18:
			rt.Events = 1
		case 19:
			rt.Events = 2
		case 21:
			rt.Events = 3
		}
		rt.TxHash = value.TxHash
		rets = append(rets, rt)
	}
	result.Page = page
	result.Limit = limit
	result.Total = total
	result.List = rets

	return &result, nil
}
