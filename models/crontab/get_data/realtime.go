/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package get_data

import (
	"github.com/IBAX-io/go-explorer/models"
)

type RealTime struct {
	Signal chan bool
}

func (p *RealTime) SendSignal() {
	select {
	case p.Signal <- true:
	default:
		//If there is still unprocessed content in the channel, not continue to send
	}
}

func (p *RealTime) ReceiveSignal() {
	if p.Signal == nil {
		p.Signal = make(chan bool)
	}
	for {
		select {
		case <-p.Signal:
			models.RealtimeWG.Wait()
			realTimeDataServer()
		}
	}
}

func realTimeDataServer() {
	go models.SyncBlockListToRedis()
	go models.DealRedisBlockTpsList()
	go models.GetTransactionBlockToRedis()
	go models.GetStatisticsSignal()
	go models.DataChartRealtimeSever()
	go models.GetHonorListToRedis("newest")
	go models.GetHonorListToRedis("pkg_rate")
	go models.GetHonorNodeMapToRedis()
	go models.UpdateHonorNodeInfo()
	go models.InitPledgeAmount()
	go models.SendTxDataSyncSignal()
	go models.InitGlobalSwitch()
	go models.SyncEcosystemInfo()
	go models.SendUtxoTxSyncSignal()
}
