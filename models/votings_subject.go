/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

type VotingSubject struct {
	Id           int64  `gorm:"primary_key;not null"`
	Ecosystem    int64  `gorm:"column:ecosystem"`
	NumberAccept int64  `gorm:"column:number_accept"`
	Results      string `gorm:"column:results;type:jsonb"`
	Subject      string `gorm:"column:subject;type:jsonb"`
	VotingId     int64  `gorm:"column:voting_id"`
}

type SubjectInfo struct {
	ContractAccept       string `json:"contract_accept"`
	ContractReject       string `json:"contract_reject"`
	ContractAcceptParams struct {
		Value          string `json:"Value"`
		VotingId       int    `json:"VotingId"`
		Conditions     string `json:"Conditions"`
		TemplateId     int    `json:"TemplateId"`
		ApplicationId  int    `json:"ApplicationId"`
		TokenEcosystem int    `json:"TokenEcosystem"`
	} `json:"contract_accept_params"`
	ContractRejectParams struct {
	} `json:"contract_reject_params"`
}

type ResultsInfo struct {
	RatingAccepted  string `json:"rating_accepted"`
	RatingRejected  string `json:"rating_rejected"`
	PercentAccepted string `json:"percent_accepted"`
	PercentRejected string `json:"percent_rejected"`
}

func (p *VotingSubject) TableName() string {
	return "1_votings_subject"
}

func (p *VotingSubject) GetByVotingId(votingId int64, ecosystem int64) (bool, error) {
	return isFound(GetDB(nil).Where("voting_id = ? and ecosystem = ?", votingId, ecosystem).Last(p))
}
