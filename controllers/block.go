/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"strconv"
	"unicode/utf8"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/converter"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type getMaxBlockIDResult struct {
	MaxBlockID int64 `json:"max_block_id"`
}

func GetMaxBlockId(c *gin.Context) {
	ret := &Response{}
	bid := getMaxBlockIDResult{}
	if models.MaxBlockId > 0 {
		bid.MaxBlockID = models.MaxBlockId
		ret.Return(bid, CodeSuccess)
		JsonResponse(c, ret)
	} else {
		bk := &models.Block{}
		f, err := bk.GetMaxBlock()
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
		}
		if f {
			bid.MaxBlockID = bk.ID
			ret.Return(bid, CodeSuccess)
			JsonResponse(c, ret)
		}
	}
}

// @tags         block detial
// @Description  block detial
// @Summary      block detial
// @Accept       json
// @Produce      json
// @Success      200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router       /auth/admin/{id} [get]
func GetBlockDetails(c *gin.Context) {
	ret := &Response{}
	blockId := c.Param("block_id")
	if blockId == "" || utf8.RuneCountInString(blockId) > 100 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	bid, err := strconv.ParseInt(blockId, 10, 64)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	var bk models.Block
	fb, err := bk.GetId(bid)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if fb {
		bdt, err := models.GetBlocksTransactionInfoByBlockInfo(&bk)
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		}
		var detail models.BlockDetailedInfoHex
		detail = *bdt
		ret.Return(detail, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	ret.ReturnFailureString("not found block id in db: " + blockId)
	JsonResponse(c, ret)
	return
}

func GetBlockDetailTxList(c *gin.Context) {
	ret := &Response{}

	req := &DataBaseFind{}

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
	if req.Block_id > 0 {
		var bk models.Block
		fb, err := bk.GetId(req.Block_id)
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		}
		if fb {
			bdt, err := models.GetBlocksTransactionListByBlockInfo(&bk)
			if err != nil {
				ret.ReturnFailureString(err.Error())
				JsonResponse(c, ret)
				return
			}
			ret.Return(models.GetBlockDetialRespones(req.Page, req.Limit, bdt), CodeSuccess)
			JsonResponse(c, ret)
			return
		} else {
			ret.ReturnFailureString("not found blockId in db: " + converter.Int64ToStr(req.Block_id))
			JsonResponse(c, ret)
			return
		}
	} else {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

}

func GetBlockList(c *gin.Context) {
	ret := &Response{}

	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page == 0 || req.Limit == 0 || (req.ReqType != 0 && req.ReqType != 1) {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	if req.Page == 1 && req.Limit == 10 && req.ReqType == 1 {
		rets, err := models.GetBlockListFromRedis()
		if err == nil {
			ret.Return(rets, CodeSuccess)
			JsonResponse(c, ret)
			return
		}
	}

	rets, err := models.GetBlockList(req.Page, req.Limit, req.ReqType, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}

func GetBlockListChart(c *gin.Context) {
	ret := &Response{}

	rets, err := models.Get15DayBlockDiffChartDataFromRedis()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}

func GetTxListChart(c *gin.Context) {
	ret := &Response{}

	rets, err := models.GetEcoLibsTxChartDataFromRedis()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}
