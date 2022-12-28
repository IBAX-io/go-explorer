/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package crontab

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	honorNode     GetDataProvider
	realTime      GetDataProvider
	history       GetDataProvider
	loadContracts GetDataProvider
	chart         GetDataProvider
)

func CreateCrontab() {
	InitCrontabTask()
	CrontabInfo := conf.GetEnvConf().Crontab
	if CrontabInfo != nil {
		go CreateCrontabFromHonorNode(CrontabInfo.HonorNode)

		//go CreateCrontabFromNodeTransaction(CrontabInfo.NodeTransaction)
		go CreateCrontabFromChartData(CrontabInfo.ChartData)
		go CreateCrontabFromRealtime(CrontabInfo.Realtime)
		go CreateCrontabFromLoadContracts(CrontabInfo.LoadContracts)
		go CreateCrontabFromHistoryData(CrontabInfo.HistoryData)
	}

}

func CreateCrontabFromHonorNode(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		honorNode.SendSignal()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Create Crontab From Honor Node Add Function failed")
	}
	c.Start()
}

func NewWithSecond() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func CreateCrontabFromWebsocket(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		select {
		case models.SendWebsocketData <- true:
		default:
		}
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("CreateCrontabFromWebsocket addFunction failed")
	}
	c.Start()
}

func CreateCrontabFromTransaction(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		if err := models.GetTransactionBlockToRedis(); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("getTransactionBlockToRedis failed")
		}
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("CreateCrontabFromTransaction addFunction failed")
	}
	c.Start()
}

func CreateCrontabFromNodeTransaction(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		if err := models.GetALLNodeTransactionList(); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("GetALLNodeTransactionList failed")
		}
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("CreateCrontabFromNodeTransaction addFunction failed")
	}
	c.Start()
}

func CreateCrontabFromRealtime(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		realTime.SendSignal()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("Create Crontab From Realtime Add Function failed")
	}
	c.Start()
}

// CreateCrontabFromChartData It can't be real-time data
func CreateCrontabFromChartData(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		chart.SendSignal()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("Create Crontab From Chart Data Add Function failed")
	}
	c.Start()
}

func CreateCrontabFromLoadContracts(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		loadContracts.SendSignal()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("Create Crontab From Load Contracts Add Function failed")
	}
	c.Start()
}

func CreateCrontabFromHistoryData(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		history.SendSignal()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("Create Crontab From history data Add Function failed")
	}
	c.Start()
}

func InitCrontabTask() {
	honorNode = NewGetDataProvider(0)
	history = NewGetDataProvider(1)
	realTime = NewGetDataProvider(2)
	loadContracts = NewGetDataProvider(3)
	chart = NewGetDataProvider(4)

	go realTime.ReceiveSignal()
	go honorNode.ReceiveSignal()
	go loadContracts.ReceiveSignal()
	go chart.ReceiveSignal()
	go history.ReceiveSignal()

	//wait receive channel start up finish
	time.Sleep(3 * time.Second)

	realTime.SendSignal()
	loadContracts.SendSignal()
	honorNode.SendSignal()
	chart.SendSignal()
	history.SendSignal()

	models.InitHonorNodeByRedis("newest")
	models.InitHonorNodeByRedis("pkg_rate")
	models.InitHonorNodeByRedis("map")
}
