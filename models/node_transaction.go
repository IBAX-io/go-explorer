/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/transaction"
	"github.com/shopspring/decimal"
)

var nodeTransaction []TransactionList

type transactionRate int8

type Transaction struct {
	Hash     []byte          `gorm:"private_key;not null"`
	Data     []byte          `gorm:"not null"`
	Used     int8            `gorm:"not null"`
	HighRate transactionRate `gorm:"not null"`
	Expedite decimal.Decimal `gorm:"not null"`
	Type     int8            `gorm:"not null"`
	KeyID    int64           `gorm:"not null"`
	Sent     int8            `gorm:"not null"`
	Verified int8            `gorm:"not null"`
	Time     int64           `gorm:"not null"`
}

// TableName returns name of table
func (p *Transaction) TableName() string {
	return "transactions"
}

func GetALLNodeTransactionList() error {
	var list []TransactionList
	rt := make(map[string]TransactionList)

	for i := 0; i < len(HonorNodes); i++ {
		rets, err := GetQueueTransactions(HonorNodes[i].APIAddress+"/api/v2/open/rowsInfo", 1, 100)
		if err != nil {
			return err
		}
		if rets == nil {
			continue
		}
		for _, value := range rets.List {
			key := hex.EncodeToString(value.Hash)
			if _, ok := rt[key]; !ok {
				li := TransactionList{}
				li.Time = value.Time
				li.Hash = hex.EncodeToString(value.Hash)

				//data, err := hex.DecodeString(value.Data)
				//if err != nil {
				//	return err
				//}
				t, err := transaction.UnmarshallTransaction(bytes.NewBuffer(value.Data))
				if err != nil {
					//if t != nil {
					//	transaction.MarkTransactionBad(t.DbTransaction, t.TxHash(), err.Error())
					//}
					return fmt.Errorf("parse transaction error(%s)", err)
				}
				if t.SmartContract() == nil {
					var cr Contract
					f, err := cr.GetById(int64(t.SmartContract().TxSmart.ID))
					if err != nil {
						return err
					}
					if f {
						li.ContractName = cr.Name
					}
				} else {
					li.ContractName = t.SmartContract().TxContract.Name
				}

				rt[key] = li

			}
		}
	}
	for _, value := range rt {
		list = append(list, value)
	}
	nodeTransaction = list
	//fmt.Printf("list len:%d\n", len(nodeTransaction))
	return nil
}

func GetNodeTransactionList(limit, page int) (*[]TransactionList, int64, error) {
	if limit < 1 || page < 1 {
		return nil, 0, errors.New("request params invalid")
	}
	offset := (page - 1) * limit
	ret := nodeTransaction
	var count int64
	count = int64(len(nodeTransaction))
	if len(ret) >= offset {
		ret = ret[offset:]
		if len(ret) >= limit {
			ret = ret[:limit]
		}
	} else {
		return nil, count, nil
	}
	return &ret, count, nil
}
