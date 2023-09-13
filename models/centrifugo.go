/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/IBAX-io/go-explorer/conf"
)

var centrifugoTimeout = time.Second * 5
var SendWebsocketData chan bool

const (
	ChannelDashboard = "dashboard"

	ChannelBlockTpsList         = "blockTpsList"
	ChannelBlockTransactionList = "blockTxList"
	ChannelBlockList            = "blockList"
	ChannelStatistical          = "statistical"
	ChannelNodeNewest           = "nodeNewest"
	ChannelNodePkgRate          = "nodePkgRate"
	ChannelNodeMap              = "nodeMap"
)

func WriteChannelByte(channel string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), centrifugoTimeout)
	defer cancel()
	return conf.GetCentrifugoConn().Conn().Publish(ctx, channel, data)
}

func SendAllWebsocketData() {
	var scanOut ScanOut
	ret1, err := scanOut.GetDashboardFromRedis()
	if err != nil {
		log.Info("Get Dashboard Redis Failed:", err.Error())
	} else {
		err := SendDashboardDataToWebsocket(ret1, ChannelStatistical)
		if err != nil {
			//log.Info("Send Websocket Failed:", err.Error(), "cmd:", ChannelStatistical)
		}
	}

	ret2, err := GetTxInfoFromRedis(30)
	if err != nil {
		log.Info("Get Tx info From Redis Failed:", err.Error())
	} else {
		err = SendTpsListToWebsocket(ret2)
		if err != nil {
			//log.Info("Send Websocket Failed:", err.Error(), "cmd:", ChannelBlockTpsList)
		}
	}

	ret3, _, err := GetTransactionBlockFromRedis()
	if err != nil {
		log.Info("Get Transaction Block From Redis Failed:", err.Error())
	} else {
		err = SendTransactionListToWebsocket(ret3)
		if err != nil {
			//log.Info("Send Websocket Failed:", err.Error(), "cmd:", ChannelBlockTransactionList)
		}
	}

	ret4, err := GetBlockListFromRedis()
	if err != nil {
		log.Info("Get Transaction Block From Redis Failed:", err.Error())
	} else {
		if err := SendBlockListToWebsocket(&ret4.List); err != nil {
			//log.WithFields(log.Fields{"INFO": err, " channel": ChannelBlockList}).Info("Send Websocket Failed")
		}
	}

	ret5, err := GetHonorNodeMapFromRedis()
	if err != nil {
		log.Info("Get Honor Node Map From Redis Failed:", err.Error())
	} else {
		err = SendDashboardDataToWebsocket(ret5, ChannelNodeMap)
		if err != nil {
			//log.WithFields(log.Fields{"INFO": err, " channel": ChannelNodeMap}).Info("Send Websocket Failed")
		}
	}

	cmd := "pkg_rate"
	ret6, err := GetHonorListFromRedis(cmd)
	if err != nil {
		log.Info("Get Honor List From Redis Failed:", err.Error(), "cmd", cmd)
	} else {
		err = SendDashboardDataToWebsocket(ret6.List, ChannelNodePkgRate)
		if err != nil {
			//log.WithFields(log.Fields{"INFO": err, " cmd": ChannelNodePkgRate}).Info("Send Websocket Failed")
		}
	}

	cmd = "newest"
	ret7, err := GetHonorListFromRedis(cmd)
	if err != nil {
		log.Info("Get Honor List From Redis Failed:", err.Error(), "cmd", cmd)
	} else {
		err = SendDashboardDataToWebsocket(ret7.List, ChannelNodeNewest)
		if err != nil {
			//log.WithFields(log.Fields{"INFO": err, " cmd": ChannelNodeNewest}).Info("Send Websocket Failed")
		}
	}
}
