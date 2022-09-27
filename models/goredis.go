/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/IBAX-io/go-explorer/conf"
)

var ctx = context.Background()

type RedisParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HashParams struct {
	Hash   string
	HArray []string       //[]string{"key1", "value1", "key2", "value2",...}
	HMap   map[string]any //map[string]interface{}{"key1": "value1", "key2": value2,...}

	Value    string
	ValueMap map[string]string
}

func initRedisServer() error {
	redisInfo := conf.GetEnvConf().RedisInfo
	return redisInfo.Init()
}

func InitRedisDb() error {
	if err := initRedisServer(); err != nil {
		return err
	}
	res, err := GetRdDb().FlushDB(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("res:", res)
	return nil
}

func InitRedisDbAll() error {
	if err := initRedisServer(); err != nil {
		return err
	}
	res, err := GetRdDb().FlushAll(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("res all:", res)
	return nil
}

func GetRdDb() *redis.Client {
	return conf.GetRedisDbConn().Conn()
}

func (rp *RedisParams) Set() error {
	return GetRdDb().Set(ctx, rp.Key, rp.Value, 0).Err()
}

func (rp *RedisParams) SetExpire(expire time.Duration) error {
	return GetRdDb().Set(ctx, rp.Key, rp.Value, expire).Err()
}

func (rp *RedisParams) TTL() (time.Duration, error) {
	return GetRdDb().TTL(ctx, rp.Key).Result()
}

func (rp *RedisParams) Get() error {
	val, err := conf.GetRedisDbConn().Conn().Get(ctx, rp.Key).Result()
	//if err != nil && err != redis.Nil {
	//	return err
	//}
	if err != nil {
		return err
	}
	rp.Value = val
	return nil
}

func (rp *RedisParams) Exist() (bool, error) {
	return RdExist(rp.Key)
}

func RdExist(key string) (bool, error) {
	n, err := GetRdDb().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

//RdExists all key exists return true else return false
func RdExists(keys ...string) (bool, error) {
	keyLen := int64(len(keys))
	if keyLen <= 0 {
		return false, errors.New("input keys invalid")
	}
	n, err := GetRdDb().Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	if n > 0 {
		if keyLen == n {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

//RdRange intercept string[start:end]
//if key not exist return ""
func RdRange(key string, start, end int64) (string, error) {
	return GetRdDb().GetRange(ctx, key, start, end).Result()
}

func (rp *RedisParams) Del() error {
	return GetRdDb().Del(ctx, rp.Key).Err()
}

func (rp *RedisParams) Size() (int64, error) {
	return GetRdDb().DBSize(ctx).Result()
}

func (rp *HashParams) HMapSet() error {
	return GetRdDb().HSet(ctx, rp.Hash, rp.HMap).Err()
}

func (rp *HashParams) HArraySet() error {
	return GetRdDb().HSet(ctx, rp.Hash, rp.HArray).Err()
}

// HSet example:
// "key1", "value1", "key2", value2 ...
func (rp *HashParams) HSet(values ...any) error {
	return GetRdDb().HSet(ctx, rp.Hash, values...).Err()
}

func (rp *HashParams) HGet(key string) error {
	val, err := GetRdDb().HGet(ctx, rp.Hash, key).Result()
	if err != nil {
		return err
	}
	rp.Value = val

	return nil
}

func (rp *HashParams) HGetAll() error {
	val, err := GetRdDb().HGetAll(ctx, rp.Hash).Result()
	if err != nil {
		return err
	}
	rp.ValueMap = val

	return nil
}

// HDel delete hash table keys
func (rp *HashParams) HDel(keys ...string) error {
	return GetRdDb().HDel(ctx, rp.Hash, keys...).Err()
}

// HExists hash table key exists
func (rp *HashParams) HExists(key string) (bool, error) {
	return GetRdDb().HExists(ctx, rp.Hash, key).Result()
}

// HLen return hash table keys len
func (rp *HashParams) HLen() (int64, error) {
	return GetRdDb().HLen(ctx, rp.Hash).Result()
}
