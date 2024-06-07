/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package get_data

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/models/buffer"
)

type History struct {
	Signal chan bool
}

func (p *History) SendSignal() {
	select {
	case p.Signal <- true:
	default:
		//If there is still unprocessed content in the channel, not continue to send
	}
}

func (p *History) ReceiveSignal() {
	if p.Signal == nil {
		p.Signal = make(chan bool)
	}
	for {
		select {
		case <-p.Signal:
			models.HistoryWG.Wait()
			historyDataServer()
		}
	}
}

func historyDataServer() {
	go models.GetScanOutKeyInfoToRedis()
	go models.GetSingleDayMaxTxToRedis()
	go models.GetWeekAverageValueTxToRedis()
	go models.GetActiveEcoLibsToRedis()
	go models.GetMaxBlockSizeToRedis()
	go models.GetMaxTxToRedis()

	//Delayed chart data
	go models.SyncNationalFlagIcon()
	go buffer.StartServer(buffer.GetBufferType(2))
	go models.UpdatePairBuffer()
}
