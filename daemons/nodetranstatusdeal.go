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

const TransactionsMax = "transactions_max"

func NodeTranStatusSumUpdate(ctx context.Context) error {
	var maxLen int64
	for i := 0; i < len(models.HonorNodes); i++ {
		mlen, _ := services.DealGetnodetransactionstatus(models.HonorNodes[i])
		if mlen > maxLen || mlen == 0 {
			maxLen = mlen
		}
		var bc models.BlockID
		bc.Time = time.Now().Unix()
		bc.Name = TransactionsMax
		bc.ID = maxLen
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
			dLen := len(models.HonorNodes)
			maxLen = 0
			for i := 0; i < dLen; i++ {
				mlen, _ := services.DealGetnodetransactionstatus(models.HonorNodes[i])
				if mlen > maxLen || mlen == 0 {
					maxLen = mlen
				}
			}
			//set
			var bc models.BlockID
			bc.Time = time.Now().Unix()
			bc.Name = TransactionsMax
			bc.ID = maxLen
			bc.InsertRedis()

		}
	}
}
