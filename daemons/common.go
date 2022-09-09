/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package daemons

import (
	"context"
	"fmt"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
)

var ExitCh chan error

func StartDaemons(ctx context.Context) {
	ExitCh = make(chan error)
	go func() {
		err := services.DealNodetransactionstatussqlite(ctx)
		if err != nil {
			ExitCh <- fmt.Errorf("Deal Node transaction status sqlite err:%s\n", err.Error())
		}
	}()

	go func() {
		err := services.DealNodeblocktransactionchsqlite(ctx)
		if err != nil {
			ExitCh <- fmt.Errorf("Deal Nodeblock transaction ch sqlite err:%s\n", err.Error())
		}
	}()
	go Sys_BlockWork(ctx)

	go Sys_Work_ChainValidBlock(ctx)

	go func() {
		err := EcosystemDealupdate(ctx)
		if err != nil {
			ExitCh <- fmt.Errorf("Ecosystem Dealupdate err:%s\n", err.Error())
		}
	}()

	go func() {
		GetFirstBlockTimeService()
	}()

	//go func() {
	//	err := NodeTranStatusSumupdate(ctx)
	//	if err != nil {
	//		ExitCh <- err
	//	}
	//}()

	go Sys_CentrifugoWork(ctx)

	go func() {
		err := InitReport()
		if err != nil {
			ExitCh <- fmt.Errorf("Init Report err:%s\n", err.Error())
		}
		err = models.InitTransactionData()
		if err != nil {
			ExitCh <- fmt.Errorf("Init transaction data err:%s\n", err.Error())
		}
	}()
	err := models.InitCountryLocator()
	if err != nil {
		ExitCh <- fmt.Errorf("Init Country Locator err:%s\n", err.Error())
	}

	err = models.GeoIpDatabaseInit()
	if err != nil {
		ExitCh <- fmt.Errorf("GeoIp Database Init err:%s\n", err.Error())
	}

	models.InitEcosystemInfo()

}
