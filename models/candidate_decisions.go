/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

type CandidateNodeDecisions struct {
	Id            int64  `gorm:"primary_key;not_null"`
	RequestId     int64  `gorm:"not_null"`
	Account       string `gorm:"not_null"`
	Decision      int64  `gorm:"not_null"`
	Earnest       string `gorm:"not_null"`
	DateCreated   int64  `gorm:"not_null"`
	DateUpdated   int64  `gorm:"not_null"`
	DecisionType  int64  `gorm:"not_null"`
	DateWithdrawd int64  `gorm:"not_null"`
}

func (p *CandidateNodeDecisions) TableName() string {
	return "1_candidate_node_decisions"
}
