/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-explorer/models/limiter"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
	"time"
)

// redis key
const (
	EcosystemCirculations      = "ecosystem-circulations-"
	TopTenHoldings             = "top-ten-holdings-"
	TopTenTxAccount            = "top-ten-tx-account"
	FifteenDaysActiveKeys      = "15-days-active-keys-"
	FifteenDaysStorageCapacity = "15-days-storage-capacity-"
	FifteenDaysGasFeeChart     = "15-days-gas-fee-"
	GasCombustionPieChart      = "gas-combustion-pie-"
	GasCombustionLineChart     = "gas-combustion-line-"
	FifteenDaysTxAmountChart   = "15-days-tx-amount-"
	FifteenDaysNewKeyChart     = "15-days-new-key-"
	FifteenTxCountChart        = "15-days-tx-count-"
)

type RefreshObject struct {
	Key       string
	Ecosystem int64
	Cmd       int //0:start 1:over
}

var (
	RefreshRequest chan RefreshObject
)

func SendRefreshRequest(key string, ecosystem int64) {
	if EcoNames.Get(ecosystem) == "" {
		return
	}
	ref := RefreshObject{
		Key:       key,
		Ecosystem: ecosystem,
	}
	select {
	case RefreshRequest <- ref:
	default:
	}
}

func RefreshChartDaemons(key string, ecosystem int64) {
	switch key {
	case TopTenTxAccount:
		getTopTenTxAccountToRedis(ecosystem, true)
	case FifteenDaysActiveKeys:
		getFifteenDaysActiveKeysToRedis(ecosystem, true)
	case FifteenDaysStorageCapacity:
		getStorageCapacityToRedis(ecosystem, true)
	case FifteenDaysGasFeeChart:
		getGasFeeToRedis(ecosystem, true)
	case GasCombustionPieChart:
		getGasCombustionPieToRedis(ecosystem, true)
	case GasCombustionLineChart:
		getGasCombustionLineToRedis(ecosystem, true)
	case FifteenDaysTxAmountChart:
		get15DaysTxAmountToRedis(ecosystem, true)
	case FifteenDaysNewKeyChart:
		get15DaysNewKeyToRedis(ecosystem, true)
	case FifteenTxCountChart:
		get15DaysTxCountToRedis(ecosystem, true)
	}
	ref := RefreshObject{
		Key:       key,
		Cmd:       1,
		Ecosystem: ecosystem,
	}
	select {
	case RefreshRequest <- ref:
	default:
	}
}

func GetHistoryEcosystemChartInfo() {
	limit := limiter.NewRequestLimitService(10)
	defer limit.Delete()
	addLimit := func() {
		limiterAvailable := limit.IsAvailable()
		for !limiterAvailable {
			time.Sleep(100 * time.Millisecond)
			limiterAvailable = limit.IsAvailable()
		}
		limit.Increase()
	}
	startRoutine := func(id int64, f func(int64)) {
		go func(ecosystem int64) {
			addLimit()
			f(ecosystem)
			limit.Reduce()
		}(id)
	}
	startRefRoutine := func(id int64, f func(int64, bool)) {
		go func(ecosystem int64) {
			addLimit()
			f(ecosystem, false)
			limit.Reduce()
		}(id)
	}
	for _, v := range EcosystemIdList {

		startRoutine(v, getEcosystemCirculationsToRedis)
		startRoutine(v, getTopTenHasTokenAccountToRedis)

		go func(ecosystem int64) {
			addLimit()
			defer limit.Reduce()
			err := GetAllKeysTotalAmount(ecosystem)
			if err != nil {
				log.WithFields(log.Fields{"INFO": err, "ecosystem": ecosystem}).Info("[Get All Keys Total Amount] failed")
			}
		}(v)

		startRefRoutine(v, getTopTenTxAccountToRedis)
		startRefRoutine(v, getFifteenDaysActiveKeysToRedis)
		startRefRoutine(v, getStorageCapacityToRedis)
		startRefRoutine(v, getGasFeeToRedis)
		startRefRoutine(v, getGasCombustionPieToRedis)
		startRefRoutine(v, getGasCombustionLineToRedis)
		startRefRoutine(v, get15DaysTxAmountToRedis)
		startRefRoutine(v, get15DaysNewKeyToRedis)
		startRefRoutine(v, get15DaysTxCountToRedis)
	}
	time.Sleep(1 * time.Second)
	count := limit.GetCount()
	for count > 0 {
		time.Sleep(100 * time.Millisecond)
		count = limit.GetCount()
	}
}

func getEcosystemChartInfoToRedis(key string, marshal func() ([]byte, error), expire time.Duration) error {
	v, err := marshal()
	if err != nil {
		return err
	}
	rd := RedisParams{
		Key:   key,
		Value: string(v),
	}
	return rd.SetExpire(expire)
}

func getEcosystemChartInfoFromRedis(key string, unmarshal func([]byte) error) error {
	rd := RedisParams{
		Key: key,
	}
	err := rd.Get()
	if err != nil {
		return err
	}
	return unmarshal([]byte(rd.Value))
}

func getEcosystemCirculationsToRedis(ecosystem int64) {
	rets, err := GetEcosystemCirculationsChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("get circulations chart failed")
		return
	}
	key := EcosystemCirculations + strconv.FormatInt(ecosystem, 10)
	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 0)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("get circulations chart to redis failed")
		return
	}
}

func GetEcosystemCirculationsFromRedis(ecosystem int64) (*EcoCirculationsResponse, error) {
	var rets EcoCirculationsResponse
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}

	key := EcosystemCirculations + strconv.FormatInt(ecosystem, 10)
	err := getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get top ten has token account From Redis msgpack err")
		return &rets, err
	}

	return &rets, nil
}

func getTopTenHasTokenAccountToRedis(ecosystem int64) {
	rets, err := getEcoTopTenHasTokenAccount(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get top ten holdings account failed")
		return
	}

	key := TopTenHoldings + strconv.FormatInt(ecosystem, 10)
	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 0)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("get top ten holdings msgpack failed")
		return
	}
}

func GetTopTenHasTokenAccountFromRedis(ecosystem int64) (*EcoTopTenHasTokenResponse, error) {
	var rets EcoTopTenHasTokenResponse
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}

	key := TopTenHoldings + strconv.FormatInt(ecosystem, 10)
	err := getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get top ten has token account From Redis msgpack err")
		return &rets, err
	}

	return &rets, nil
}

func getTopTenTxAccountToRedis(ecosystem int64, isRefresh bool) {
	key := TopTenTxAccount + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}

	rets, err := getEcoTopTenTxAccountChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get top ten tx account failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("get top ten tx account msgpack failed")
		return
	}
}

func GetTopTenTxAccountFromRedis(ecosystem int64) (*EcoTopTenTxAmountResponse, error) {
	var rets EcoTopTenTxAmountResponse
	key := TopTenTxAccount + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(TopTenTxAccount, ecosystem)
	}

	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get top ten tx account From Redis msgpack err")
		return &rets, err
	}

	return &rets, nil
}

func getFifteenDaysActiveKeysToRedis(ecosystem int64, isRefresh bool) {
	key := FifteenDaysActiveKeys + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}

	rets, err := getEco15DayActiveKeysChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("[15-days-active-keys-chart] get failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-active-keys-chart] to redis failed")
		return
	}
}

func GetFifteenDaysActiveKeysFromRedis(ecosystem int64) (*KeyInfoChart, error) {
	var rets KeyInfoChart
	key := FifteenDaysActiveKeys + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(FifteenDaysActiveKeys, ecosystem)
	}

	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return &rets, nil
		}
		log.WithFields(log.Fields{"warn": err}).Warn("[15-days-active-keys-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}

func getStorageCapacityToRedis(ecosystem int64, isRefresh bool) {
	key := FifteenDaysStorageCapacity + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}
	rets, err := getEco15DayStorageCapacityChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("[15-days-storage-capacity-chart] get failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-storage-capacity-chart] to redis failed")
		return
	}
}

func Get15DaysStorageCapacityFromRedis(ecosystem int64) (*StorageCapacitysChart, error) {
	var rets StorageCapacitysChart
	key := FifteenDaysStorageCapacity + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(FifteenDaysStorageCapacity, ecosystem)
	}
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("[15-days-storage-capacity-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}

func getGasFeeToRedis(ecosystem int64, isRefresh bool) {
	key := FifteenDaysGasFeeChart + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}
	rets, err := GetEco15DayGasFeeChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("[15-days-gas-fee-chart] get failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-gas-fee-chart] to redis failed")
		return
	}
}

func Get15DaysGasFeeFromRedis(ecosystem int64) (*EcoTxGasFeeDiffResponse, error) {
	var rets EcoTxGasFeeDiffResponse
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	key := FifteenDaysGasFeeChart + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(FifteenDaysGasFeeChart, ecosystem)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("[15-days-gas-fee-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}

func getGasCombustionPieToRedis(ecosystem int64, isRefresh bool) {
	key := GasCombustionPieChart + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}
	rets, err := getGasCombustionPieChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[gas-combustion-pie-chart] get failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[gas-combustion-pie-chart] to redis failed")
		return
	}
}

func GetGasCombustionPieFromRedis(ecosystem int64) (*EcoGasFeeResponse, error) {
	var rets EcoGasFeeResponse
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	key := GasCombustionPieChart + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(GasCombustionPieChart, ecosystem)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("[gas-combustion-pie-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}

func getGasCombustionLineToRedis(ecosystem int64, isRefresh bool) {
	key := GasCombustionLineChart + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}
	rets, err := getGasCombustionLineChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[gas-combustion-line-chart] get failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[gas-combustion-line-chart] to redis failed")
		return
	}
}

func GetGasCombustionLineFromRedis(ecosystem int64) (*EcoGasFeeChangeResponse, error) {
	var rets EcoGasFeeChangeResponse
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	key := GasCombustionLineChart + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(GasCombustionLineChart, ecosystem)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("[gas-combustion-line-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}

func get15DaysTxAmountToRedis(ecosystem int64, isRefresh bool) {
	key := FifteenDaysTxAmountChart + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}
	rets, err := getEco15DayTxAmountChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-tx-amount-chart] get failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-tx-amount-chart] to redis failed")
		return
	}
}

func Get15DaysTxAmountFromRedis(ecosystem int64) (*EcoTxAmountDiffResponse, error) {
	var rets EcoTxAmountDiffResponse
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	key := FifteenDaysTxAmountChart + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(FifteenDaysTxAmountChart, ecosystem)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("[15-days-tx-amount-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}

func get15DaysNewKeyToRedis(ecosystem int64, isRefresh bool) {
	key := FifteenDaysNewKeyChart + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}
	rets := getEcosystemNewKeyChart(ecosystem, 15)
	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err := getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-new-key-chart] to redis failed")
		return
	}
}

func Get15DaysNewKeyFromRedis(ecosystem int64) (*KeyInfoChart, error) {
	var rets KeyInfoChart
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	key := FifteenDaysNewKeyChart + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(FifteenDaysNewKeyChart, ecosystem)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("[15-days-new-key-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}

func get15DaysTxCountToRedis(ecosystem int64, isRefresh bool) {
	key := FifteenTxCountChart + strconv.FormatInt(ecosystem, 10)
	if !isRefresh {
		rd := RedisParams{
			Key: key,
		}
		exist, err := rd.Exist()
		if err != nil {
			return
		}
		if exist {
			return
		}
	}
	rets, err := getEco15DayTransactionChart(ecosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-tx-count-chart] get failed")
		return
	}

	f := func() ([]byte, error) {
		return msgpack.Marshal(rets)
	}

	err = getEcosystemChartInfoToRedis(key, f, 10*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem": ecosystem}).Error("[15-days-tx-count-chart] to redis failed")
		return
	}
}

func Get15DaysTxCountFromRedis(ecosystem int64) (*TxListChart, error) {
	var rets TxListChart
	f := func(data []byte) error {
		return msgpack.Unmarshal(data, &rets)
	}
	key := FifteenTxCountChart + strconv.FormatInt(ecosystem, 10)
	rd := RedisParams{
		Key: key,
	}
	exist, err := rd.Exist()
	if err != nil {
		return &rets, err
	}
	if !exist {
		SendRefreshRequest(FifteenTxCountChart, ecosystem)
	}
	err = getEcosystemChartInfoFromRedis(key, f)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "ecosystem": ecosystem}).Warn("[15-days-tx-count-chart] From Redis err")
		return &rets, err
	}

	return &rets, nil
}
