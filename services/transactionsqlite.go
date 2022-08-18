/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"context"

	"github.com/IBAX-io/go-explorer/models"
)

type NodeTranStatuslistData struct {
	Data *[]models.TransactionStatus
}

var (
	NodeTranStatusCh = make(chan *NodeTranStatuslistData, 10)
)

func DealNodeblocktransactionchsqlite(ctx context.Context) error {
	bk := &models.TransactionStatus{}
	for {
		select {
		case <-ctx.Done():
			return nil
		case dat := <-NodeTranStatusCh:
			bk.BatchInsert_Sqlites(dat.Data)
		}
	}
}
