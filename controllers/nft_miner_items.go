/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"encoding/json"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"unicode/utf8"
)

func GetAccountDetailNftMinerHandler(c *gin.Context) {
	var req models.MineHistoryRequest
	ret := &Response{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.Return(nil, CodeJsonformaterr.Errorf(err))
		JsonResponse(c, ret)
		return
	}

	if req.KeyId == "" || req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	items := &models.NftMinerItems{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	res, err := items.GetAccountDetailNftMinerInfo(req.KeyId, req.Order, req.Page, req.Limit)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	JsonResponse(c, ret)
}

func NftMinerInfoHandler(c *gin.Context) {
	ret := &Response{}
	search := c.Param("search")
	if search == "" || utf8.RuneCountInString(search) > 100 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	items := &models.NftMinerItems{}
	res, err := items.GetNftMinerBySearch(search)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	PureJsonResponse(c, ret)
}

func NftMinerHistoryInfoHandler(c *gin.Context) {
	var req models.NftMinerInfoRequest
	ret := &Response{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.Return(nil, CodeJsonformaterr.Errorf(err))
		JsonResponse(c, ret)
		return
	}

	if req.Search == nil || req.Limit <= 0 || req.Page <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	es := &models.NftMinerEvents{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	res, err := es.GetNftHistoryInfo(req.Search, req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNftMinerStakeInfoHandler(c *gin.Context) {
	var req models.NftMinerInfoRequest
	ret := &Response{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.Return(nil, CodeJsonformaterr.Errorf(err))
		JsonResponse(c, ret)
		return
	}

	if req.Search == nil || req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	stak := &models.NftMinerStaking{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	res, err := stak.GetNftMinerStakeInfo(req.Search, req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNftMinerTxInfoHandler(c *gin.Context) {
	var req models.NftMinerInfoRequest
	ret := &Response{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.Return(nil, CodeJsonformaterr.Errorf(err))
		JsonResponse(c, ret)
		return
	}

	if req.Search == nil || req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	items := &models.NftMinerItems{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	res, err := items.GetNftMinerTxInfo(req.Search, req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNftMinerMetaverse(c *gin.Context) {
	ret := &Response{}
	items := &models.NftMinerItems{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	rets, err := items.GetNftMetaverse()
	if err != nil {
		ret.Return(nil, CodeDBfinderr.Errorf(err))
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNftMinerMapHandler(c *gin.Context) {
	ret := &Response{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	rets, err := models.GetNftMinerMap()
	if err != nil {
		ret.Return(nil, CodeDBfinderr.Errorf(err))
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNftMinerMetaverseList(c *gin.Context) {
	var req models.NftMinerInfoRequest
	ret := &Response{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.Return(nil, CodeJsonformaterr.Errorf(err))
		JsonResponse(c, ret)
		return
	}

	if req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	items := &models.NftMinerItems{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	res, err := items.NftMinerMetaverseList(req.Page, req.Limit, req.Order)
	if err != nil {
		ret.Return(nil, CodeDBfinderr.Errorf(err))
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNftMinerRegionHandler(c *gin.Context) {
	ret := &Response{}
	params := c.Query("limit")
	var err error
	var limit int
	if params != "all" {
		limit, err = strconv.Atoi(params)
		if err != nil {
			ret.ReturnFailureString("request params invalid")
			JsonResponse(c, ret)
			return
		}
	}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}

	res, err := models.GetTopTwentyNftRegionChart(limit)
	if err != nil {
		ret.Return(nil, CodeDBfinderr.Errorf(err))
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	JsonResponse(c, ret)
}

func NftMinerRegionListHandler(c *gin.Context) {
	ret := &Response{}
	if !models.NftMinerReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}

	req := &GeneralRequest{}
	params := c.Query("params")
	if params == "" {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	err := json.Unmarshal([]byte(params), req)
	if err != nil {
		ret.ReturnFailureString("request params invalid:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	res, err := models.GetNftMinerRegionList(req.Page, req.Limit)
	if err != nil {
		ret.Return(nil, CodeDBfinderr.Errorf(err))
		JsonResponse(c, ret)
		return
	}
	ret.Return(res, CodeSuccess)
	JsonResponse(c, ret)
}
