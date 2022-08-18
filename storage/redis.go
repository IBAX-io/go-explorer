/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rc *redis.Client
var ctx = context.Background()

type RedisModel struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

func (r *RedisModel) Str() string {
	return fmt.Sprintf("%s:%d", r.Address, r.Port)
}

func (r *RedisModel) Init() error {
	rc = redis.NewClient(&redis.Options{
		Addr:     r.Str(),
		Password: r.Password,
		DB:       r.Db,
	})
	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisModel) Conn() *redis.Client {
	return rc
}
func (l *RedisModel) Close() error {
	return rc.Close()
}
