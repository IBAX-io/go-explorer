/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package storage

import (
	"github.com/shopspring/decimal"
	"time"
)

type HonorNodeModel struct {
	NodeName     string `json:"node_name"`
	TCPAddress   string `json:"tcp_address,omitempty"`
	APIAddress   string `json:"api_address"`
	City         string `json:"city"`
	Icon         string `json:"icon"`
	IconUrl      string `json:"icon_url"`
	NodePosition int64  `json:"node_position"`
	KeyID        string `json:"key_id"`
	Display      bool   `json:"display"`
	//PublicKey       string          `json:"public_key"`
	Latitude        string          `json:"latitude,omitempty"`
	Longitude       string          `json:"longitude,omitempty"`
	NodeBlock       int64           `json:"node_block"`
	PkgAccountedFor decimal.Decimal `json:"pkg_accounted_for"`
	ReplyRate       string          `json:"reply_rate"`
	ConsensusMode   int32           `json:"consensus_mode"`

	NodeStatusTime time.Time `json:"node_status_time,omitempty"`
}

type Crontab struct {
	HonorNode     string `yaml:"honor_node"`
	LoadContracts string `yaml:"load_contracts"`
	ChartData     string `yaml:"chart_data"`
	Dashboard     string `yaml:"dashboard"`
}

type CryptoSettings struct {
	Cryptoer string `yaml:"cryptoer"`
	Hasher   string `yaml:"hasher"`
}
