/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/shopspring/decimal"
)

type ResObject struct{}

type WebRequest struct {
	Head   *RequestHead   ` json:"head"`
	Params *RequestParams ` json:"params"`
}

type WebResponse struct {
	Head *RequestHead  ` json:"head"`
	Body *ResponseBoby ` json:"body"`
}

type RequestHead struct {
	Version   string ` json:"version"`
	Msgtype   string `json:"msgtype"`
	Interface string `json:"interface"`
	Remark    string `json:"remark"`
}

type RequestParams struct {
	Cmd          string `json:"cmd,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	CurrentPage  int    `json:"current_page,omitempty"`
	DatabaseId   string `json:"database_id,omitempty"`
	Ecosystem    int64  `json:"ecosystem,omitempty"`
	Wallet       string `json:"wallet,omitempty"`
	SearchType   string `json:"searchType,omitempty"`
	Block_id     int64  `json:"block_id,omitempty"`
	Table_name   string `json:"table_name,omitempty"`
	NodePosition int64  `json:"nodeposition,omitempty"`
	Where        string `json:"where,omitempty"`
	Hash         string `json:"hash,omitempty"`
	Order        string `json:"order,omitempty"`
}

type ResponseBoby struct {
	Cmd          string          `json:"cmd,omitempty"`
	PageSize     int             `json:"page_size,omitempty"`
	CurrentPage  int             `json:"current_page,omitempty"`
	RetDataType  string          `json:"ret_data_type,omitempty"`
	NodePosition int64           `json:"nodeposition,omitempty"`
	Wallet       string          `json:"wallet,omitempty"`
	Ecosystem    int64           `json:"ecosystem,omitempty"`
	Block_id     int64           `json:"database_id,omitempty"`
	TableName    string          `json:"table_name,omitempty"`
	Order        string          `json:"order,omitempty"`
	Hash         string          `json:"hash,omitempty"`
	Total        int64           `json:"total,omitempty"`
	Sum          decimal.Decimal `json:"sum,omitempty"`
	Data         any             `json:"data,omitempty"`
	Ret          string          `json:"ret,omitempty"`
	Retcode      int             `json:"retcode,omitempty"`
	Retinfo      string          `json:"retinfo,omitempty"`
}

type DBWebInfo struct {
	Id         string `json:"id,omitempty"`
	Nodename   string `yaml:"nodename" json:"nodename"`
	IconUrl    string `yaml:"icon_url" json:"icon_url"`
	Name       string `json:"name,omitempty"`
	Engine     string `json:"engine,omitempty"`
	Version    string `json:"backend_version,omitempty"`
	APIAddress string `json:"api_address,omitempty"`
}

type DashboardTopInfo struct {
	Title  string `json:"title"`
	Number int64  `json:"number"`
	Icon   string `json:"icon"`
	Color  string `json:"color"`
}

type DashboardNumInfo struct {
	MaxBlockID       int64 `json:"max_block_id"`
	MaxEcosystemID   int64 `json:"max_ecosystem_id"`
	MaxTransactionID int64 `json:"max_Transaction_id"`
	MaxNodeID        int64 `json:"max_node_id"`
}

type BlockccdataInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
type BlockPriceInfo struct {
	Name           string  `json:"name"`
	Symbol         string  `json:"symbol"`
	Price          float32 `json:"price"`
	High           float32 `json:"high"`
	Low            float32 `json:"low"`
	Hist_high      float32 `json:"hist_high"`
	Hist_low       float32 `json:"hist_low"`
	Timestamps     int64   `json:"timestamps"`
	Volume         float32 `json:"volume"`
	Display_volume float32 `json:"display_volume"`
	Usd_volume     float32 `json:"usd_volume"`
	Change_hourly  float32 `json:"change_hourly"`
	Change_daily   float32 `json:"change_daily"`
	Change_weekly  float32 `json:"change_weekly"`
	Change_monthly float32 `json:"change_monthly"`
}

type FindForm struct {
	Where map[string]any `json:"where"`                             //
	Order string         `json:"order" example:"date_created desc"` //
	Page  int            `json:"page"`                              //
	Limit int            `json:"limit"`                             //
}

type EcosytemTranscationHistoryFind struct {
	Ecosystem int64          `json:"ecosystem"`
	Wallet    string         `json:"wallet"`
	Search    string         `json:"search"`
	Where     map[string]any `json:"where"`
	Order     string         `json:"order"`
	Page      int            `json:"page"`
	Limit     int            `json:"limit"`
	ReqType   int            `json:"type"`
	AppId     int64          `json:"app_id"`
	Hash      string         `json:"hash"`
}

type DataBaseFind struct {
	Cmd          string `json:"cmd,omitempty"`
	Page_size    int    `json:"page_size,omitempty"`
	Current_page int    `json:"current_page,omitempty"`
	Database_id  string `json:"database_id,omitempty"`
	Ecosystem    int64  `json:"ecosystem,omitempty"`
	Wallet       string `json:"wallet,omitempty"`
	SearchType   string `json:"searchType,omitempty"`
	Block_id     int64  `json:"block_id,omitempty"`
	Table_name   string `json:"table_name,omitempty"`
	NodePosition int64  `json:"nodeposition,omitempty"`
	Hash         string `json:"hash,omitempty"`
	Where        string `json:"where"`                             //
	Order        string `json:"order" example:"date_created desc"` //
	Page         int    `json:"page"`                              //
	Limit        int    `json:"limit"`                             //
}

type DataBaseRespone struct {
	Cmd           string          `json:"cmd"`
	Page_size     int             `json:"page_size"`
	Current_page  int             `json:"current_page"`
	Ret_data_type string          `json:"ret_data_type,omitempty"`
	NodePosition  int64           `json:"nodeposition"`
	Wallet        string          `json:"wallet,omitempty"`
	Ecosystem     int64           `json:"ecosystem,omitempty"`
	Block_id      int64           `json:"database_id,omitempty"`
	Table_name    string          `json:"table_name,omitempty"`
	Order         string          `json:"order,omitempty"`
	Hash          string          `json:"hash,omitempty"`
	Total         int             `json:"total"`
	Sum           decimal.Decimal `json:"sum,omitempty"`
	Data          any             `json:"data,omitempty"`
	Ret           string          `json:"ret,omitempty"`
	Retcode       int             `json:"retcode,omitempty"`
	Retinfo       string          `json:"retinfo,omitempty"`
}

type AccountAmountChangeRequest struct {
	Ecosystem int64  `json:"ecosystem"`
	Wallet    string `json:"wallet"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	IsDefault int    `json:"is_default"`
}

type GeneralRequest struct {
	Search    any            `json:"search"`
	Page      int            `json:"page"`
	Limit     int            `json:"limit"`
	Order     string         `json:"order"`
	Where     map[string]any `json:"where"`
	StartTime int64          `json:"startTime"`
	EndTime   int64          `json:"endTime"`
	Language  string         `json:"language"`
}

func init() {
	binding.EnableDecoderUseNumber = true
}
