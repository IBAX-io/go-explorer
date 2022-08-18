/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

/*
// Transaction is model
type Transaction struct {
	Hash     []byte          `gorm:"private_key;not null"`
	Data     []byte          `gorm:"not null"`
	Used     int8            `gorm:"not null"`
	HighRate transactionRate `gorm:"not null"`
	Type     int8            `gorm:"not null"`
	KeyID    int64           `gorm:"not null"`
	Counter  int8            `gorm:"not null"`
	Sent     int8            `gorm:"not null"`
	Attempt  int8            `gorm:"not null"`
	Verified int8            `gorm:"not null;default:1"`
}

type TransactionpageHex struct {
	Hash     string          `json:"hash"`
	Data     string          `json:"data"`
	Used     int8            `json:"used"`
	HighRate transactionRate `json:"highrate"`
	Type     int8            `json:"type"`
	KeyID    string          `json:"key_id"`
	Counter  int8            `json:"counter"`
	Sent     int8            `json:"sent"`
	Attempt  int8            `json:"attempt"`
	Verified int8            `json:"verified"`
}

func GetTransactionpages(page int, size int) (*[]TransactionpageHex, int64, error) {
	var (
		tss    []Transaction
		ret    []TransactionpageHex
		num    int64
		ioffet int
	)

	err := conf.GetDbConn().Conn().Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
	if page < 1 || size < 1 {
		return &ret, num, err
	}
	num = int64(len(tss))

	ioffet = (page - 1) * size

	if int(num) < page*size {
		size = int(num) % size
	}

	if int(num) < ioffet || num < 1 {
		return &ret, num, err
	}
	for i := 0; i < size; i++ {

		var tx = TransactionpageHex{
			Hash:     hex.EncodeToString(tss[ioffet].Hash),
			Type:     tss[ioffet].Type,
			Data:     hex.EncodeToString(tss[ioffet].Data),
			Used:     tss[ioffet].Used,
			HighRate: tss[ioffet].HighRate,
			KeyID:    strconv.FormatInt(tss[ioffet].KeyID, 10),
			Counter:  tss[ioffet].Counter,
			Sent:     tss[ioffet].Sent,
			Attempt:  tss[ioffet].Attempt,
			Verified: tss[ioffet].Verified,
		}
		ret = append(ret, tx)
		ioffet++
	}

	return &ret, num, err
}

*/
