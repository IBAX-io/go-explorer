/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Get15DayGasFeeChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.FifteenDaysGasFee)
	if err != nil {
		ret.ReturnFailureString("get 15 Day Gas Fee Chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeContributionListHandler(c *gin.Context) {
	ret := &Response{}
	var req GeneralRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString("request params marshal json failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetNodeContributionList(req.Page, req.Limit)
	if err != nil {
		ret.ReturnFailureString("Get Node Contribution List Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeContributionChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetNodeContributionChart()
	if err != nil {
		ret.ReturnFailureString("Get Honor Node Contribution Chart Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNewCirculationsChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.FifteenDaysNewCirculations)
	if err != nil {
		ret.ReturnFailureString("Get Circulations Chart Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetAccountChangeChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.AccountChange)
	if err != nil {
		ret.ReturnFailureString("Get Account Change Chart Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetAccountActiveChartHandler(c *gin.Context) {
	ret := &Response{}
	var req GeneralRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString("request params marshal json failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Search == nil && (req.StartTime == 0 || req.EndTime == 0) {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetDailyActiveAccountChangeChart(req.Search, req.StartTime, req.EndTime)
	if err != nil {
		ret.ReturnFailureString("Get Daily Active Account Change Chart Failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetAccountActiveListHandler(c *gin.Context) {
	ret := &Response{}
	var req GeneralRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString("request params marshal json failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.StartTime == 0 || req.EndTime == 0 || req.Page == 0 || req.Limit == 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetDailyActiveReport(req.Page, req.Limit, req.StartTime, req.EndTime)
	if err != nil {
		ret.ReturnFailureString("Get Daily Active Report Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetStakingAccountHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetTopTenStakingAccount()
	if err != nil {
		ret.ReturnFailureString("Get Top Ten Staking Account Failed")
		JsonResponse(c, ret)
		return
	}

	//server-timing
	//c.Header("Server-Timing", fmt.Sprintf("db;dur=%d, app;dur=47.2", time.Now().Sub(t1).Milliseconds()))

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetBlockSizeChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.FifteenDaysBlockSize)
	if err != nil {
		ret.ReturnFailureString("Get Block Size Chart Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetBlockSizeListHandler(c *gin.Context) {
	ret := &Response{}
	var req GeneralRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString("request params marshal json failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page == 0 || req.Limit == 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.Get15DayBlockSizeList(req.Page, req.Limit)
	if err != nil {
		ret.ReturnFailureString("Get Block Size List Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetTxChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.Get15DaysTxCountFromRedis(1)
	if err != nil {
		ret.ReturnFailureString("Get Tx Chart Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetTxListHandler(c *gin.Context) {
	ret := &Response{}
	var req GeneralRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString("request params marshal json failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page == 0 || req.Limit == 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.Get15DayTxList(req.Page, req.Limit)
	if err != nil {
		ret.ReturnFailureString("Get Tx List Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNewKeysHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.NewKey)
	if err != nil {
		ret.ReturnFailureString("Get New Keys Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetBlockNumberHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.FifteenDaysBlockNumber)
	if err != nil {
		ret.ReturnFailureString("Get Block Number Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NewNftMinerHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.NewNfTMinerChange)
	if err != nil {
		ret.ReturnFailureString("Get New Nft Miner History Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NftMinerRewardHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.NftMinerRewardChange)
	if err != nil {
		ret.ReturnFailureString("Get Nft Miner Reward History Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NftMinerIntervalHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetNftMinerIntervalChart()
	if err != nil {
		ret.ReturnFailureString("Get Nft Miner Energy Point Interval Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NftMinerIntervalListHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetNftMinerIntervalListChart()
	if err != nil {
		ret.ReturnFailureString("Get Nft Miner Energy Point Interval List Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NftMinerEnergyPowerChangeHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetNftEnergyPowerChangeChart()
	if err != nil {
		ret.ReturnFailureString("Get Nft Miner Energy Power Change Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NftMinerStakedChangeHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetNftMinerStakedChangeChart()
	if err != nil {
		ret.ReturnFailureString("Get Nft Miner Staked Change Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetHistoryNewEcosystemHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.NewEcosystemChange)
	if err != nil {
		ret.ReturnFailureString("Get History New Ecosystem Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetTokenEcosystemRatioHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetTokenEcosystemRatioChart()
	if err != nil {
		ret.ReturnFailureString("Get Token Ecosystem Ratio Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetTopTenEcosystemTxHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.TopTenEcosystemTx)
	if err != nil {
		ret.ReturnFailureString("Get Top Ten Ecosystem Tx Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNewEcosystemChartListHandler(c *gin.Context) {
	ret := &Response{}
	var req GeneralRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString("request params marshal json failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page == 0 || req.Limit == 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetNewEcosystemChartList(req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString("Get New Ecosystem Chart List Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetTopTenMaxKeysEcosystemHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetTopTenMaxKeysEcosystem()
	if err != nil {
		ret.ReturnFailureString("Get Top Ten Max Keys Ecosystem Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetMultiFeeEcosystemChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetMultiFeeEcosystemChart()
	if err != nil {
		ret.ReturnFailureString("Get Multi Fee Ecosystem Chart Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcosystemGovernModelChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetEcosystemGovernModelChart()
	if err != nil {
		ret.ReturnFailureString("Get Ecosystem Govern Model Chart Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}
