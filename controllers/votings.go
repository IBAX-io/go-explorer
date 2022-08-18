/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetLatestDaoVotingHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetLatestDaoVoting()
	if err != nil {
		ret.ReturnFailureString("Get Latest Dao Voting Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeDaoVoteListHandler(c *gin.Context) {
	ret := &Response{}

	req := &GeneralRequest{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 || req.Search == nil {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetDaoVoteList(req.Search, req.Page, req.Limit)
	if err != nil {
		ret.ReturnFailureString("Get Node Dao Vote List Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetDaoVoteListHandler(c *gin.Context) {
	ret := &Response{}

	req := &GeneralRequest{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetDaoVoteList(nil, req.Page, req.Limit)
	if err != nil {
		ret.ReturnFailureString("Get Dao Vote List Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetDaoVoteChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.DaoVoteChart)
	if err != nil {
		ret.ReturnFailureString("Get Dao Vote Chart From Redis Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func DaoVoteDetailSearchHandler(c *gin.Context) {
	ret := &Response{}
	req := &GeneralRequest{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Language == "" || req.Search == nil {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetDaoVoteDetail(req.Search, req.Language)
	if err != nil {
		ret.ReturnFailureString("Get Dao Vote Detail Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeVoteHistoryHandler(c *gin.Context) {
	ret := &Response{}
	req := &GeneralRequest{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 || req.Search == nil {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetNodeVotingHistory(req.Search, req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString("Get Node Voting History Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeStakingHistoryHandler(c *gin.Context) {
	ret := &Response{}
	req := &GeneralRequest{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 || req.Search == nil {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetNodeStakingHistory(req.Search, req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString("Get Node Staking History Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}
