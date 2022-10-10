/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"fmt"
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
	TotalSupplyToken = "2100000000000000000000"
	UtxoTx           = "UTXO_Tx"
	UtxoTransfer     = "UTXO_Transfer"

	AccountUTXO = "Account-UTXO"
	UTXOAccount = "UTXO-Account"
)

var (
	buildBranch = ""
	buildDate   = ""
	commitHash  = ""
)

var BuildInfo string

func Version() string {
	status := `scan server status is running`
	fmt.Printf("BuildInfo:%s\n", BuildInfo)
	return strings.TrimSpace(strings.Join([]string{status, VERSION, BuildInfo}, " "))
}

// go build -ldflags "-X 'github.com/IBAX-io/go-explorer/cmd.buildBranch=main' -X 'github.com/IBAX-io/go-explorer/cmd.buildDate=2022-06-17' -X 'github.com/IBAX-io/go-explorer/cmd.commitHash=2141saf'"
func InitBuildInfo() {
	BuildInfo = func() string {
		if buildBranch == "" {
			return fmt.Sprintf("branch.%s commit.%s time.%s", "unknown", "unknown", "unknown")
		}
		return fmt.Sprintf("branch.%s commit.%s time.%s", buildBranch, commitHash, buildDate)
	}()
}
