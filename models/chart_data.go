/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

func MsgpackMarshal(v any) ([]byte, error) {
	if res, err := msgpack.Marshal(v); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func MsgpackUnmarshal(bt []byte, v any) error {
	if err := msgpack.Unmarshal(bt, &v); err != nil {
		return err
	}
	return nil
}

func GetEcoLibsChartDataToRedis() {
	var eco Ecosystem
	rets, err := eco.GetBasisEcosystemChart()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetBasisEcoLibs error")
		return
	}

	val, err := MsgpackMarshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetBasisEcoLibs marshal error")
		return
	}

	rd := RedisParams{
		Key:   "ecoLibs-chart",
		Value: string(val),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetBasisEcoLibs set redis error")
		return
	}
}

func GetEcoLibsChartData() (*BasisEcosystemResponse, error) {
	var err error
	var eco Ecosystem
	rets := &BasisEcosystemResponse{}
	chart := &BasisEcosystemChartDataResponse{}
	rd := RedisParams{
		Key:   "ecoLibs-chart",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetEcoLibsChartDataFromRedis getDb err")
		return nil, err
	}
	err = MsgpackUnmarshal([]byte(rd.Value), chart)
	if err != nil {
		return nil, err
	}

	rets, err = eco.GetBasisEcosystem()
	if err != nil {
		return nil, err
	}
	rets.TxInfo = chart.TxInfo
	rets.KeyInfo = chart.KeyInfo
	return rets, nil
}

func GetEcoLibsTxChartDataToRedis() {
	rets, err := GetEcoLibsTransaction()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetEcoLibsTransaction error")
		return
	}
	val, err := MsgpackMarshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetEcoLibsTransaction marshal error")
		return
	}

	rd := RedisParams{
		Key:   "ecoLibs-tx-chart",
		Value: string(val),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetEcoLibsTransaction set redis error")
		return
	}
}

func GetEcoLibsTxChartDataFromRedis() ([]EcosystemTxRatioChart, error) {
	var err error
	var ret []EcosystemTxRatioChart

	rd := RedisParams{
		Key:   "ecoLibs-tx-chart",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetEcoLibsTxChartDataFromRedis getDb err")
		return nil, err
	}
	err = MsgpackUnmarshal([]byte(rd.Value), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func getBlockDiffChartDataFromDays(days int) (*BlockListChart, error) {
	var (
		rets BlockListChart
	)

	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	t1 := yesterday.AddDate(0, 0, -1*days)

	var list []DaysNumber
	err := GetDB(nil).Raw(fmt.Sprintf(`SELECT to_char(to_timestamp(time),'yyyy-MM-dd') days,count(*) num FROM 
"block_chain" WHERE time >= %d GROUP BY days`, t1.Unix())).Find(&list).Error
	if err != nil {
		return &rets, err
	}

	rets.Block = make([]int64, days)
	rets.Time = make([]int64, days)
	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.Block[i] = GetDaysNumber(rets.Time[i], list)
	}
	return &rets, nil
}

func Get15DayBlockDiffChartDataToRedis() {
	rets, err := getBlockDiffChartDataFromDays(15)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get15DayBlockDiffChartDataToRedis error")
		return
	}

	val, err := MsgpackMarshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get15DayBlockDiffChartDataToRedis marshal error")
		return
	}

	rd := RedisParams{
		Key:   "block-diff-chart",
		Value: string(val),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get15DayBlockDiffChartDataToRedis set redis error")
		return
	}
}

func Get15DayBlockDiffChartDataFromRedis() (BlockListChart, error) {
	var err error
	var rets BlockListChart
	rd := RedisParams{
		Key:   "block-diff-chart",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("Get15DayBlockDiffChartDataFromRedis getDb err")
		return rets, err
	}
	err = MsgpackUnmarshal([]byte(rd.Value), &rets)
	if err != nil {
		return rets, err
	}
	return rets, nil
}
