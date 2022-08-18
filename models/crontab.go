/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func CreateCrontab() {
	InitCrontabTask()
	CrontabInfo := conf.GetEnvConf().Crontab
	if CrontabInfo != nil {
		go CreateCrontabFromHonorNode(CrontabInfo.HonorNode)

		//go CreateCrontabFromNodeTransaction(CrontabInfo.NodeTransaction)
		go CreateCrontabFromChartData(CrontabInfo.ChartData)
		go CreateCrontabFromDashboard(CrontabInfo.Dashboard)
		go CreateCrontabFromLoadContracts(CrontabInfo.LoadContracts)
	}

}

func CreateCrontabFromHonorNode(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		getHonorNodeInfo()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("CreateCrontabFromHonorNode addFunction failed")
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
		case SendWebsocketData <- true:
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
		if err := getTransactionBlockToRedis(); err != nil {
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
		if err := GetALLNodeTransactionList(); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("GetALLNodeTransactionList failed")
		}
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("CreateCrontabFromNodeTransaction addFunction failed")
	}
	c.Start()
}

func CreateCrontabFromDashboard(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		realTimeDataServer()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("CreateCrontabFromHotEcosystemInfo addFunction failed")
	}
	c.Start()
}

//CreateCrontabFromChartData It can't be real-time data
func CreateCrontabFromChartData(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		historyDataServer()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("CreateCrontabFromDashboardChartData addFunction failed")
	}
	c.Start()
}

func CreateCrontabFromLoadContracts(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		SendLoadContractsSignal()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeSet": timeSet}).Error("Create Crontab From Load Contracts addFunction failed")
	}
	c.Start()
}

func InitCrontabTask() {
	initGlobalSwitch()
	historyDataServer()
	realTimeDataServer()
	InitHonorNodeRedis("newest")
	InitHonorNodeRedis("pkg_rate")
	InitHonorNodeRedis("map")
	getHonorNodeInfo()
}

func historyDataServer() {
	go GetDashboardChartDataToRedis()
	go GetEcoLibsChartDataToRedis()
	go GetEcoLibsTxChartDataToRedis()
	go Get15DayBlockDiffChartDataToRedis()
	go InsertDailyActiveReport()
	go InsertDailyNodeReport()
	go DataChartHistoryServer()
}

func realTimeDataServer() {
	go SyncBlockListToRedis()
	go DealRedisBlockTpsList()
	go getTransactionBlockToRedis()
	SendStatisticsSignal()
	go DataChartRealtimeSever()
	go GetHonorListToRedis("newest")
	go GetHonorListToRedis("pkg_rate")
	go GetHonorNodeMapToRedis()
	go updateHonorNodeInfo()
	go initPledgeAmount()
	go SendTxDataSyncSignal()
}
