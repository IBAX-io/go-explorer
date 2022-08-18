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

func (rp *RedisParams) Set() error {
	return GetRdDb().Set(ctx, rp.Key, rp.Value, 0).Err()
}

func (rp *RedisParams) SetExpire(expire time.Duration) error {
	return GetRdDb().Set(ctx, rp.Key, rp.Value, expire).Err()
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
