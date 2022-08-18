/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package daemons

import (
	//"fmt"
	"context"
	"sync/atomic"
	"time"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"

	"github.com/IBAX-io/go-explorer/services"
	//"github.com/IBAX-io/go-explorer/models"
	log "github.com/sirupsen/logrus"
)

//recvice websocket request sendTopData
func EcosystemDealupdate(ctx context.Context) error {
	//bk := &models.TransactionStatus{}
	services.Deal_Redis_Dashboard()
	models.SendWebsocketData = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			log.Error("NodeTranStatusSumupdate done his work")
			return nil
		case <-models.SendWebsocketData:
			if err := services.Deal_Redis_Dashboard(); err != nil {
				log.Info("send topdata err:", err)
			}

		}
	}
}

func Sys_BlockWork(ctx context.Context) {
	go services.SyncDealBlock()
	models.GetScanOut = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			return
		case <-models.GetScanOut:
			if err := models.GetScanOutDataToRedis(); err != nil {
				log.Info("GetScanOutDataToRedis failed:", err)
			}
			err := services.WorkDealBlock()
			if err != nil {
				log.Info("WorkDealBlock:", err)
			}
		}
	}
}

func GetFirstBlockTimeService() {
	for {
		if err := models.InitFristBlockTime(); err != nil {
			log.Info("InitFirstBlockService Failed:", err)
			time.Sleep(5 * time.Second)
		} else {
			return
		}
	}
}

func Sys_Work_ChainValidBlock(ctx context.Context) {
	ChainValidBlockInit()

	for {
		select {
		case <-ctx.Done():
			log.Info("Sys_Work_ChainValidBlock done his work")
			return
		case <-time.After(time.Second * 4):
			//err := ChainValidBlock()
			//if err != nil {
			//	log.Info("ChainValidBlock")
			//}
		}
	}
}

var bgOnChaimRun uint32

//ChainValidBlock get chain valid block
//Scan do not need to add a consensus mechanism,Redis stores a large number of Scan-block_id
func ChainValidBlock() error {
	if atomic.CompareAndSwapUint32(&bgOnChaimRun, 0, 1) {
		defer atomic.StoreUint32(&bgOnChaimRun, 0)
	} else {
		return nil
	}
	var cf sqldb.Confirmation
	var bc models.BlockID
	f, err := cf.GetGoodBlockLast()
	if err != nil {
		return err
	}
	if f {
		fc, err := bc.GetbyName(models.ChainMax)
		if err != nil {
			if err.Error() == "redis: nil" || err.Error() == "EOF" {
			} else {
				return err
			}
		}
		if !fc {
			bc.Time = cf.Time
			bc.Name = models.ChainMax
			bc.ID = cf.BlockID
			return bc.InsertRedis()
		}
		if cf.BlockID > bc.ID {
			bc.ID = cf.BlockID
			bc.Time = cf.Time
			return bc.InsertRedis()
		}
	}

	return nil
}

func ChainValidBlockInit() {
	var bc models.BlockID
	var bk models.Block
	fc, er := bk.GetMaxBlock()
	if er != nil {
		log.Info("Chain Valid Block Init Last Block Failed:", er.Error())
		return
	}
	if !fc {
		log.Info("Chain Valid Block Init Block Doesn't Not Exist")
		return
	}
	fm, err := bc.GetbyName(models.MintMax)
	if err == nil && fm {
		if bc.ID > bk.ID {
			err = bc.DelbyName(models.MintMax)
			if err != nil {
				log.Info("Chain Valid Block Init Block Del mint_max Failed:", er.Error())
				return
			}
			//bc.DelbyName(models.ChainMax)
		}
	}
}

func Sys_CentrifugoWork(ctx context.Context) {
	models.SendScanOut = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			return
		case <-models.SendScanOut:
			var scanOut models.ScanOut
			rets, err := scanOut.GetRedisdashboard()
			if err != nil {
				log.Info("Get Redis Dashboard Failed:", err.Error())
			} else {
				if err := SendToWebsocket(rets, &scanOut); err != nil {
					log.Info("Send Dashboard To Websocket Failed:", err.Error())
				}
			}
		}
	}
}

func SendToWebsocket(rets *models.ScanOutRet, scanOut *models.ScanOut) error {
	err := models.SendDashboardDataToWebsocket(rets, models.ChannelStatistical)
	if err != nil {
		return err
	}

	return nil
}

func InitReport() error {
	var de models.DailyNodeReport
	err := de.CreateTable()
	if err != nil {
		return err
	}
	return models.InitDailyActiveReport()
}
