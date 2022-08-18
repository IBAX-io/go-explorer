/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package daemons

import (
	"context"
	"time"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	//"encoding/hex"
)

const TransactionsMax = "transactionsMax_max"

func NodeTranStatusSumupdate(ctx context.Context) error {
	var maxlen int64
	for i := 0; i < len(models.HonorNodes); i++ {
		mlen, _ := services.DealGetnodetransactionstatus(models.HonorNodes[i])
		if mlen > maxlen || mlen == 0 {
			maxlen = mlen
		}
		var bc models.BlockID
		bc.Time = time.Now().Unix()
		bc.Name = TransactionsMax
		bc.ID = maxlen
		err := bc.InsertRedis()
		if err != nil {
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(1 * time.Second):
			dlen := len(models.HonorNodes)
			maxlen = 0
			for i := 0; i < dlen; i++ {
				mlen, _ := services.DealGetnodetransactionstatus(models.HonorNodes[i])
				if mlen > maxlen || mlen == 0 {
					maxlen = mlen
				}
			}
			//set
			var bc models.BlockID
			bc.Time = time.Now().Unix()
			bc.Name = TransactionsMax
			bc.ID = maxlen
			bc.InsertRedis()

		}
	}
}
