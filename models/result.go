/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"time"
)

type BlockHeaderInfo struct {
	BlockID      int64  `json:"block_id"`
	Time         int64  `json:"time"`
	EcosystemID  int64  `json:"ecosystem_id"`
	KeyID        int64  `json:"key_id"`
	NodePosition int64  `json:"node_position"`
	Sign         []byte `json:"sign"`
	Hash         []byte `json:"hash"`
	Version      int    `json:"version"`
}

type TxDetailedInfo struct {
	Hash         []byte         `json:"hash"`
	ContractName string         `json:"contract_name"`
	Params       map[string]any `json:"params"`
	KeyID        int64          `json:"key_id"`
	Time         int64          `json:"time"`
	Type         int64          `json:"type"`
}

type BlockDetailedInfo struct {
	Header        BlockHeaderInfo  `json:"header"`
	Hash          []byte           `json:"hash"`
	EcosystemID   int64            `json:"ecosystem_id"`
	NodePosition  int64            `json:"node_position"`
	KeyID         int64            `json:"key_id"`
	Time          int64            `json:"time"`
	Tx            int32            `json:"tx_count"`
	RollbacksHash []byte           `json:"rollbacks_hash"`
	MrklRoot      []byte           `json:"mrkl_root"`
	BinData       []byte           `json:"bin_data"`
	SysUpdate     bool             `json:"sys_update"`
	GenBlock      bool             `json:"gen_block"`
	StopCount     int              `json:"stop_count"`
	Transactions  []TxDetailedInfo `json:"transactions"`
}

type paramValue struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Conditions string `json:"conditions"`
}

type ecosystemParamsResult struct {
	List []paramValue `json:"list"`
}

type listResult struct {
	Count string              `json:"count"`
	List  []map[string]string `json:"list"`
}

type BlocksResult struct {
	BlockID      int64  `json:"id"`
	Time         int64  `json:"time"`
	EcosystemID  int64  `json:"ecosystem_id"`
	KeyID        string `json:"key_id"`
	NodePosition int64  `json:"node_position"`
	PreHash      string `json:"pre_hash"`
	Hash         string `json:"hash"`
	Tx           int32  `json:"tx"`
}

type BlockHeaderInfoHex struct {
	BlockId       int64  `json:"block_id"`
	Time          int64  `json:"time"`
	EcosystemId   int64  `json:"ecosystem_id"`
	KeyId         string `json:"key_id"`
	NodePosition  int64  `json:"node_position"`
	PreHash       string `json:"pre_hash"`
	Sign          string `json:"sign"`
	BlockHash     string `json:"block_hash"`
	Version       int32  `json:"version"`
	Address       string `json:"address"`
	IconUrl       string `json:"icon_url"`
	ConsensusMode int32  `json:"consensus_mode"`
}

type TxDetailedInfoHex struct {
	Hash         string `json:"hash"`
	ContractName string `json:"contract_name"`
	//Params       map[string]any `json:"params"`
	Params string `json:"params"`
	KeyID  string `json:"key_id"`
	Time   int64  `json:"time"`
	Type   int64  `json:"type"`
	Size   int64  `json:"size"`

	Ecosystemname string `json:"ecosystemname"`
	TokenSymbol   string `json:"token_symbol"`
	Ecosystem     int64  `json:"ecosystem"`
}

type TxDetailedInfoResponse struct {
	Hash         string `json:"hash"`
	ContractName string `json:"contract_name"`
	//Params       map[string]any `json:"params"`
	Params string `json:"params"`
	KeyID  string `json:"key_id"`
	Time   int64  `json:"time"`
	Type   int64  `json:"type"`
	Size   string `json:"size"`

	Ecosystemname string `json:"ecosystemname"`
	TokenSymbol   string `json:"token_symbol"`
	Ecosystem     int64  `json:"ecosystem"`
}

type BlockDetailedInfoHex struct {
	Header        BlockHeaderInfoHex  `json:"header"`
	Hash          string              `json:"hash"`
	EcosystemID   int64               `json:"ecosystem_id"`
	NodePosition  int64               `json:"node_position"`
	KeyID         string              `json:"key_id"`
	Time          int64               `json:"time"`
	Tx            int32               `json:"tx_count"`
	RollbacksHash string              `json:"rollbacks_hash"`
	MerkleRoot    string              `json:"merkle_root"`
	BinData       string              `json:"bin_data"`
	SysUpdate     bool                `json:"sys_update"`
	GenBlock      bool                `json:"gen_block"`
	StopCount     int                 `json:"stop_count"`
	BlockSize     string              `json:"block_size"`
	TxTotalSize   string              `json:"tx_total_size"`
	Transactions  []TxDetailedInfoHex `json:"transactions"`
}

type BlockDetailedInfoHexRespone struct {
	Total        int64                    `json:"total"`
	Page         int                      `json:"page"`
	Limit        int                      `json:"limit"`
	Transactions []TxDetailedInfoResponse `json:"transactions"`
}

type BlockHeaderInfoDetailed struct {
	Header BlockHeaderInfoHex `json:"header"`
	//Info          BlockInfo           `json:"info"`
	Transactions []TxDetailedInfoHex `json:"transactions"`
}

type TransactionStatusHex struct {
	Hash          string `gorm:"primary_key;not null" json:"hash"`
	Time          int64  `gorm:"not null" json:"time"`
	Type          int64  `gorm:"not null" json:"type"`
	WalletID      string `gorm:"not null" json:"wallet_id"`
	BlockID       int64  `gorm:"not null" json:"block_id"`
	Error         string `gorm:"not null" json:"err"`
	Penalty       int64  `gorm:"not null"  json:"penalty"`
	Ecosystemname string `json:"ecosystemname"`
	TokenSymbol   string `json:"token_symbol"`
	Ecosystem     int64  `json:"ecosystem"`
}

type EcosystemList struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Info           string `json:"info"`
	IsValued       int64  `json:"isValued"`
	EmissionAuount string `json:"emission_auount"`
	TokenSymbol    string `json:"token_symbol"`
	TypeEmission   int64  `json:"type_emission"`
	TypeWithdraw   int64  `json:"type_withdraw"`
	Member         int64  `json:"member"`
	AppParams      any    `json:"app_params,omitempty"`
	Status         int    `json:"status"` //0:Not joined 1:join 2:unknown
}

type TransactionHistoryHex struct {
	Name      string  `json:"name,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

//because of PublicKey is byte
type FullNodeCityJSON struct {
	TCPAddress string `json:"tcp_address"`
	APIAddress string `json:"api_address"`
	City       string `json:"city"`
	Icon       string `json:"icon"`
	//KeyID      int64  `json:"key_id"`
	KeyID     json.Number `json:"key_id"`
	PublicKey string      `json:"public_key"`
	//UnbanTime  json.Number `json:"unban_time,er"`
	UnbanTime time.Time `json:"unbantime"`
}

//because of PublicKey is byte
type FullNodeCityJSONHex struct {
	TCPAddress string `json:"tcp_address"`
	APIAddress string `json:"api_address"`
	City       string `json:"city"`
	Icon       string `json:"icon"`
	//KeyID      int64  `json:"key_id"`
	KeyID     string `json:"key_id"`
	PublicKey string `json:"public_key"`
	//UnbanTime  json.Number `json:"unban_time,er"`
	UnbanTime time.Time `json:"unbantime"`
}

type MineGpsInfo struct {
	ID          int64  `gorm:"primary_key;not null"`
	Devid       int64  `gorm:"not null"`
	Ip          string `gorm:"not null"`
	Location    string `gorm:"not null"`
	Longitude   string `gorm:"not null"`
	Latitude    string `gorm:"not null"`
	DateUpdated int64  `gorm:"not null" example:"2019-07-19 17:45:31"`
	DateCreated int64  `gorm:"not null" example:"2019-07-19 17:45:31"`
}

type DBTransactionsInfo struct {
	Name        string `json:"name"`
	Transaction int32  `json:"transaction"`
}

type DashboardChainInfo struct {
	High string `json:"high"`
	Low  string `json:"low"`
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
	Last string `json:"last"`
	Vol  string `json:"vol"`
}

type RatesInfo struct {
	Base      string         `json:"base"`
	Timestamp int64          `json:"timestamp"`
	Rates     map[string]any `json:"rates"`
}

type BlockCCRatesInfo struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    RatesInfo `json:"data"`
}

type ResponseTopDataBoby struct {
	TopData           any `json:"topdata,omitempty"`
	TopBlocks         any `json:"topblocks,omitempty"`
	TopTransactions   any `json:"toptransactions,omitempty"`
	TopTransactiontps any `json:"toptransactiontps,omitempty"`
}

type ResponseDashboardTitle struct {
	Cmd  string `json:"cmd,omitempty"`
	List any    `json:"list,omitempty"`
}

type BlockCCPriceInfo struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    []map[string]any `json:"data"`
}

type StatisticsData struct {
	AllTransactionsNum int64 `json:"transactions"`
	ChainContractsNum  int64 `json:"contracts"`
	GuardNode          int64 `json:"node"`
	StorageCapacity    int64 `json:"storage"`
	EcosystemsNum      int64 `json:"ecosystems"`
}
