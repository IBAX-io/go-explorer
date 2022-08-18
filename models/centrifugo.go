/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"context"
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
