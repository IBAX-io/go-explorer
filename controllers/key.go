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

func GetAccountList(c *gin.Context) {
	req := &EcosytemTranscationHistoryFind{}
	ret := &Response{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 || req.Ecosystem <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	ts := &models.Key{}
	rets, err := ts.GetAccountList(req.Page, req.Limit, req.Order, req.Ecosystem)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return

}

func GetAccountTotalAmountChart(c *gin.Context) {
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

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

	rets, err := models.GetAccountTotalAmount(req.Ecosystem, req.Wallet)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

func GetAmountChangePieChart(c *gin.Context) {
	ret := &Response{}
	req := &AccountAmountChangeRequest{}

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

	rets, err := models.GetAmountChangePieChart(req.Ecosystem, req.Wallet, req.StartTime, req.EndTime)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

func GetAmountChangeBarChart(c *gin.Context) {
	ret := &Response{}
	req := &AccountAmountChangeRequest{}

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

	rets, err := models.GetAmountChangeBarChart(req.Ecosystem, req.Wallet, req.IsDefault)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

func GetAccountTxChart(c *gin.Context) {
	ret := &Response{}
	req := &AccountAmountChangeRequest{}

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

	rets, err := models.GetAccountTxChart(req.Ecosystem, req.Wallet)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}
