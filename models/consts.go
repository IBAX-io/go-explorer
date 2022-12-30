/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/shopspring/decimal"
	"strings"
)

const (
	// VERSION is current version
	VERSION = "1.0.0"

	// ApiPath is the beginning of the api url
	ApiPath = `/api/v2/`

	ecosysTable      = "1_ecosystems"
	SysTokenSymbol   = "IBXC"
	SysEcosystemName = "platform ecosystem"
	UtxoTx           = "UTXO_Tx"
	UtxoTransferSelf = "UTXO_Transfer_Self"
	UtxoBurning      = "UTXO_Burning"

	BlackHoleAddr = "0000-0000-0000-0000-0000"

	AccountUTXO = "Account-UTXO"
	UTXOAccount = "UTXO-Account"
)

const (
	publicRoundSupply = 42000000
	devTempSupply     = 315000000
	foundationSupply  = 315000000
	partnersSupply    = 210000000
	privateRound1     = 42000000
	privateRound2     = 147000000
)

const AssignTotalSupply = publicRoundSupply + devTempSupply + foundationSupply + partnersSupply + privateRound1 + privateRound2
const NftMinerTotalSupply = 393750000
const MintNodeTotalSupply = 630000000
const StartUpSupply = consts.FounderAmount

const TotalSupply = AssignTotalSupply + NftMinerTotalSupply + MintNodeTotalSupply + StartUpSupply

var (
	buildBranch = ""
	buildDate   = ""
	commitHash  = ""
)
var TotalSupplyToken decimal.Decimal

var BuildInfo string

func Version() string {
	status := `scan server status is running`
	fmt.Printf("BuildInfo:%s\n", BuildInfo)
	return strings.TrimSpace(strings.Join([]string{status, VERSION, BuildInfo}, " "))
}

// go build -ldflags "-X 'github.com/IBAX-io/go-explorer/models.buildBranch=main' -X 'github.com/IBAX-io/go-explorer/models.buildDate=2022-06-17' -X 'github.com/IBAX-io/go-explorer/models.commitHash=2141saf'"
func InitBuildInfo() {
	BuildInfo = func() string {
		if buildBranch == "" {
			return fmt.Sprintf("branch.%s commit.%s time.%s", "unknown", "unknown", "unknown")
		}
		return fmt.Sprintf("branch.%s commit.%s time.%s", buildBranch, commitHash, buildDate)
	}()
}

func init() {
	TotalSupplyToken = decimal.New(TotalSupply, int32(consts.MoneyDigits))
}
