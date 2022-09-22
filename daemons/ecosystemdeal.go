/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package daemons

import (
	//"fmt"
	"context"
	"time"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	//"github.com/IBAX-io/go-explorer/models"
	log "github.com/sirupsen/logrus"
)

//recvice websocket request sendTopData
func EcosystemDealUpdate(ctx context.Context) error {
	services.DealRedisDashboard()
	models.SendWebsocketData = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			log.Error("Ecosystem Deal Update done his work")
			return nil
		case <-models.SendWebsocketData:
			if err := services.DealRedisDashboard(); err != nil {
				log.Info("Deal Redis Dashboard err:", err)
			}

		}
	}
}

func SyncBlockWork(ctx context.Context) {
	models.GetScanOut = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			return
		case <-models.GetScanOut:
			if err := models.GetScanOutDataToRedis(); err != nil {
				log.Info("GetScanOutDataToRedis failed:", err)
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

func SyncCentrifugoWork(ctx context.Context) {
	models.SendWebsocketSignal = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			return
		case <-models.SendWebsocketSignal:
			models.SendAllWebsocketData()
		}
	}
}

func InitReport() error {
	var de models.DailyNodeReport
	err := de.CreateTable()
	if err != nil {
		return err
	}
	return models.InitDailyActiveReport()
}
