/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/vmihailenco/msgpack/v5"
)

type BlockID struct {
	ID   int64
	Time int64
	Name string
}

func (b *BlockID) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(b); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (b *BlockID) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &b); err != nil {
		return err
	}
	return nil
}

func (b *BlockID) GetByName(name string) (bool, error) {
	rp := &RedisParams{
		Key: name,
	}
	if err := rp.Get(); err != nil {
		return false, err
	}
	if err := b.Unmarshal([]byte(rp.Value)); err != nil {
		return false, err
	}
	return true, nil
}

func (b *BlockID) DelByName(name string) error {
	rp := &RedisParams{
		Key: name,
	}
	err := rp.Del()
	return err
}

func (b *BlockID) InsertRedis() error {
	val, err := b.Marshal()
	if err != nil {
		return err
	}
	rp := RedisParams{
		Key:   b.Name,
		Value: string(val),
	}
	for i := 0; i < 5; i++ {
		err = rp.Set()
		if err == nil {
			break
		}
	}
	return err
}
