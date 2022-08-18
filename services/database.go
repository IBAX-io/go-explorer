/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"encoding/json"
	"strconv"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/sirupsen/logrus"
)

func Set_ALLTables(dat []map[string]string) error {
	lg1, err := json.Marshal(dat)
	if err == nil {
		rp := models.RedisParams{
			Key:   "ALLTables",
			Value: string(lg1),
		}
		err = rp.Set()
		if err != nil {
			logrus.Info("redis Setdb3 err key: %s  value:%s", rp.Key, rp.Value)
		}
	}
	return err

}

func Get_ALLTables(id int64) (*[]map[string]string, error) {
	var fs []map[string]string
	sid := strconv.FormatInt(id, 10)
	rp := &models.RedisParams{
		Key: sid + "ALLTables",
	}
	err := rp.Get()
	if err == nil {
		err = json.Unmarshal([]byte(rp.Value), &fs)
		return &fs, err
	}
	return &fs, err
}

func Set_ColumnTypeTables(name string, dat []map[string]string) error {
	lg1, err := json.Marshal(dat)
	if err == nil {
		rp := models.RedisParams{
			Key:   name + "_ColumnTypeTables",
			Value: string(lg1),
		}
		err = rp.Set()
		if err != nil {
			logrus.Info("redis Setdb3 err key: %s  value:%s", rp.Key, rp.Value)
		}
	}
	return err

}
func Get_ColumnTypeTables(name string) (*[]map[string]string, error) {
	var fs []map[string]string
	rp := &models.RedisParams{
		Key: name + "_ColumnTypeTables",
	}
	err := rp.Get()
	if err == nil {
		err = json.Unmarshal([]byte(rp.Value), &fs)
		return &fs, err
	}
	return &fs, err
}
