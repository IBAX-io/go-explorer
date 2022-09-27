/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetAccountTokenChangeHandler(c *gin.Context) {
	req := &AccountAmountChangeRequest{}
	ret := &Response{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Ecosystem <= 0 || req.Wallet == "" {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	wid := converter.StringToAddress(req.Wallet)
	if wid == 0 && req.Wallet != "0000-0000-0000-0000-0000" {
		ret.ReturnFailureString("invalid account")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetAccountTokenChangeChart(req.Ecosystem, wid, 0)
	if err != nil {
		ret.ReturnFailureString("Get Account Token Change Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

func GetEcosystemCirculationsChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetEcosystemCirculationsChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("get ecosystem circulations chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoTopTenHasTokenChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetTopTenHasTokenAccountFromRedis(ecosystem)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoTopTenTxAccountChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetEcoTopTenTxAccountChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("get ecosystem top ten tx chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetGasCombustionPieChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetGasCombustionPieChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("Get Gas Combustion Pie Chart Handler failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetGasCombustionLineChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetGasCombustionLineChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("Get Gas Combustion Line Chart Handler Failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoTxAmountChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetEco15DayTxAmountChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("get ecosystem tx amount chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoGasFeeChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetEco15DayGasFeeChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("get ecosystem gas fee chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoNewKeyChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets := models.GetEcosystemNewKeyChart(ecosystem, 15)

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoActiveKeyChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetEco15DayActiveKeysChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("get ecosystem active key chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoTransactionChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetEco15DayTransactionChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("get ecosystem transaction chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcoStorageCapacityChartHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("ecosystem")
	ecosystem := converter.StrToInt64(idStr)
	if ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetEco15DayStorageCapacitysChart(ecosystem)
	if err != nil {
		ret.ReturnFailureString("get ecosystem storage capacity chart failed")
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}
