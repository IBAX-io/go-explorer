/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
)

func (b *Block) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(b); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (b *Block) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &b); err != nil {
		return err
	}
	return nil
}

func (b *Block) GetRedisByhash(hash []byte) (bool, error) {
	rd := RedisParams{
		Key:   "block-" + string(hash),
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return false, nil
		}
		return false, err
	}

	if err := b.Unmarshal([]byte(rd.Value)); err != nil {
		return true, err
	}
	return true, nil
}

func (b *Block) GetRedisByid(id int64) (bool, error) {
	rd := RedisParams{
		Key:   "block-" + strconv.FormatInt(id, 10),
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return false, nil
		}
		return false, err
	}
	if err := b.Unmarshal([]byte(rd.Value)); err != nil {
		return true, err
	}
	return true, nil
}

func (b *Block) InsertRedis() error {
	val, err := b.Marshal()
	if err != nil {
		return err
	}
	rd := RedisParams{
		Key:   "block-" + string(b.Hash),
		Value: string(val),
	}
	err = rd.Set()
	if err != nil {
		return err
	}

	rd = RedisParams{
		Key:   "block-" + strconv.FormatInt(b.ID, 10),
		Value: string(val),
	}
	err = rd.Set()
	if err != nil {
		return err
	}
	return nil
}

func GetBlockDetialRespones(page int, limit int, txs *[]TxDetailedInfoResponse) *BlockDetailedInfoHexRespone {
	var (
		ret     BlockDetailedInfoHexRespone
		st, end int
	)
	transactions := *txs
	ret.Limit = limit
	ret.Page = page
	ret.Total = int64(len(transactions))

	if page > 0 {
		st = (page - 1) * limit
		end = page * limit
	} else {
		st = 0
		end = limit
	}
	if end > int(ret.Total) {
		end = int(ret.Total)
	}
	if st < int(ret.Total) {
		ret.Transactions = transactions[st:end]
	}
	return &ret
}
