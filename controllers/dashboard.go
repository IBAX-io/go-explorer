/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	"github.com/gin-gonic/gin"
)

func DashboardGetToken(c *gin.Context) {
	ret := &Response{}
	rets, err := services.GetJWTCentToken(1, 60*60)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
	} else {
		ret.Return(rets, CodeSuccess)
		JsonResponse(c, ret)
	}
}

// GetDashboard godoc
// @Summary      get dashboard
// @Description  get dashboard statistical data
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200 {object} Response{data=models.ScanOutRet} code:0
// @Failure      200 {object} Response{data=models.ScanOutRet} code:1
// @Router       /api/v2/dashboard [get]
func GetDashboard(c *gin.Context) {
	ret := &Response{}
	var scanout models.ScanOut
	rets, err := scanout.GetDashboardFromRedis()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}

func GetBlockTpsLists(c *gin.Context) {
	ret := &Response{}
	rets, err := services.GetGroupBlockTpsLists()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}

func GetDashboardChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetDashboardChartDataFromRedis()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}
