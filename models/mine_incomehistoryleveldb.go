/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"strconv"

	"github.com/vmihailenco/msgpack/v5"
)

func (b *MineIncomehistory) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(b); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (b *MineIncomehistory) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &b); err != nil {
		return err
	}
	return nil
}

func (b *MineIncomehistory) GetRedisByhash(hash []byte) (bool, error) {
	rd := RedisParams{
		Key:   "mih-" + string(hash),
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

func (b *MineIncomehistory) GetRedisbyid(id int64) (bool, error) {
	rd := RedisParams{
		Key:   "mih-" + strconv.FormatInt(id, 10),
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

func (b *MineIncomehistory) Insert_redis() error {
	val, err := b.Marshal()
	if err != nil {
		return err
	}
	rd := RedisParams{
		Key:   "mih-" + string(b.Mine_incomehistory_hash),
		Value: string(val),
	}
	err = rd.Set()
	if err != nil {
		return err
	}

	rd = RedisParams{
		Key:   "mih-" + strconv.FormatInt(b.ID, 10),
		Value: string(val),
	}
	err = rd.Set()
	if err != nil {
		return err
	}

	return nil
}
