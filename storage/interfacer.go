/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package storage

type DbConner interface {
	GormInit() error
	Close() error
}

// TransactionStatus is model
type TransactionStatus struct {
	Hash      []byte `gorm:"primary_key;not null"  json:"hash"`
	Time      int64  `gorm:"not null" json:"time"`
	Type      int64  `gorm:"not null"  json:"type"`
	Ecosystem int64  `gorm:"not null"  json:"ecosystem"`
	WalletID  int64  `gorm:"not null"  json:"wallet_id"`
	BlockID   int64  `gorm:"not null;index:tsblockid_idx"  json:"block_id"`
	Error     string `gorm:"not null"  json:"error"`
	Penalty   int64  `gorm:"not null"  json:"penalty"`
}

type BlockTxDetailedInfoHex struct {
	BlockID      int64  `gorm:"not null;index:blockid_idx" json:"block_id"`
	Hash         string `gorm:"primary_key;not null" json:"hash"`
	ContractName string `gorm:"not null" json:"contract_name"`
	//Params       map[string]any `json:"params"`
	Params string `gorm:"not null" json:"params"`
	KeyID  string `gorm:"not null" json:"key_id"`
	Time   int64  `gorm:"not null" json:"time"`
	Type   int64  `gorm:"not null" json:"type"`
	Size   int64  `gorm:"not null" json:"size"`

	Ecosystemname string `gorm:"null" json:"ecosystemname"`
	Token_title   string `gorm:"null" json:"token_title"`
	Ecosystem     int64  `gorm:"null" json:"ecosystem"`
}
