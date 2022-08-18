/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"context"
	"time"

	"github.com/IBAX-io/go-explorer/storage"

	"github.com/IBAX-io/go-explorer/models"
	log "github.com/sirupsen/logrus"
)

type NodeTransactionStatus struct {
	Nodename     string `yaml:"nodename" json:"nodename"`
	NodePosition int64  `yaml:"nodeposition" json:"nodeposition"`
	Data         *[]models.TransactionStatus
}

var (
	NodeTranStatusDaemonCh = make(chan *NodeTransactionStatus, 100)
)

func DealGetnodetransactionstatus(node storage.HonorNodeModel) (int64, error) {
	var count int64
	var err error
	count, err = models.GetQueueTransactionsCount(node.APIAddress + "/api/v2/open/rowsInfo")
	if err != nil {
		log.Info("models GetQueueTransactions transactions false:" + node.NodeName + err.Error())
	}

	if node.NodeStatusTime.IsZero() {
		ret, err := models.GetTransactionsStatus(node.APIAddress+"/api/v2/open/rowsInfo", 0, 1, 100)
		if err != nil {
			log.Info("models GetTransactionsStatus false: " + node.NodeName + err.Error())
			return 0, err
		}

		if ret != nil && len(*ret) > 0 {
			//log.Info("models GetTransactionsStatus len: ", len(*ret))

			dat := NodeTransactionStatus{}
			dat.Nodename = node.NodeName
			dat.NodePosition = node.NodePosition
			dat.Data = ret
			NodeTranStatusDaemonCh <- &dat

			node.NodeStatusTime = time.Now()
			node.NodeStatusTime = node.NodeStatusTime.AddDate(0, 0, -1)
		}
	} else {
		ret, err := models.GetTransactionsStatus(node.APIAddress+"/api/v2/open/rowsInfo", node.NodeStatusTime.Unix(), 1, 100)
		if err != nil {
			log.Info("models GetTransactionsStatus false:" + node.NodeName + err.Error())
		} else if ret != nil && len(*ret) > 0 {
			dat := NodeTransactionStatus{}
			dat.Nodename = node.NodeName
			dat.NodePosition = node.NodePosition
			dat.Data = ret
			NodeTranStatusDaemonCh <- &dat
			node.NodeStatusTime = time.Now()
			node.NodeStatusTime = node.NodeStatusTime.AddDate(0, 0, -1)
		}
	}

	return count, nil
}

func DealNodetransactionstatussqlite(ctx context.Context) error {
	bk := &models.TransactionStatus{}
	for {
		select {
		case <-ctx.Done():
			return nil
		case dat := <-NodeTranStatusDaemonCh:
			err := bk.DbconnbatchinsertSqlites(dat.Data)
			if err != nil {
				return err
			}
		}
	}
}
