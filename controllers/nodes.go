/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"encoding/json"
	"github.com/IBAX-io/go-ibax/packages/converter"
	log "github.com/sirupsen/logrus"
	"sort"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//GetHonorNodeLists
//params: body ?body={"page":1,"limit":10,"order":"newest"}
//GET
func GetHonorNodeLists(c *gin.Context) {
	req := &FindForm{}
	ret := &Response{}
	rets := &models.GeneralResponse{}

	body := c.Query("body")
	err := json.Unmarshal([]byte(body), req)
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
	switch req.Order {
	case "newest":
		if req.Limit == 10 && req.Page == 1 {
			rets, err := models.GetHonorListFromRedis(req.Order)
			if err != nil {
				ret.ReturnFailureString("Get Honor List Newest Failed:" + err.Error())
				JsonResponse(c, ret)
				return
			}
			ret.Return(rets, CodeSuccess)
			JsonResponse(c, ret)
			return
		}
		var bk models.Block
		var list []any

		bkList, err := bk.GetBlocksFrom(req.Page, req.Limit, "desc")
		if err != nil {
			ret.ReturnFailureString("Get Blocks Nodes Failed")
			JsonResponse(c, ret)
			return
		}

		for i := 0; i < len(bkList); i++ {
			replyRate, err := models.GetNodeBlockReplyRate(&bkList[i])
			if err != nil {
				log.WithFields(log.Fields{"INFO": err, "block id": bkList[i].ID}).Info("Get Honor List To Redis Get Reply Rate Failed")
				return
			}
			for _, value := range models.HonorNodes {
				if value.NodePosition == bkList[i].NodePosition && value.ConsensusMode == bkList[i].ConsensusMode {
					value.ReplyRate = replyRate
					list = append(list, value)
				}
			}
		}

		rets.Page = req.Page
		rets.Limit = req.Limit
		rets.Total = int64(len(list))
		rets.List = list
		ret.Return(rets, CodeSuccess)
		JsonResponse(c, ret)
		return
	case "pkg_rate":
		if req.Page == 1 && req.Limit == 5 {
			rets, err := models.GetHonorListFromRedis(req.Order)
			if err != nil {
				ret.ReturnFailureString("Get Honor List Pkg Rate Failed:" + err.Error())
				JsonResponse(c, ret)
				return
			}
			ret.Return(rets, CodeSuccess)
			JsonResponse(c, ret)
			return
		}
		sort.Sort(models.LeaderboardSlice(models.HonorNodes))
	default:
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	offset := (req.Page - 1) * req.Limit
	rets.Page = req.Page
	rets.Limit = req.Limit
	if len(models.HonorNodes) >= offset {
		data := models.HonorNodes[offset:]
		if len(data) >= req.Limit {
			data = data[:req.Limit]
		}
		rets.Total = int64(len(data))
		rets.List = data
		ret.Return(rets, CodeSuccess)
	} else {
		ret.Return(nil, CodeSuccess)
	}
	JsonResponse(c, ret)

}

func GetNodeMap(c *gin.Context) {
	ret := &Response{}
	var cs models.CandidateNodeRequests
	if !models.NodeReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	rets, err := cs.GetNodeMap()
	if err != nil {
		ret.ReturnFailureString("Get Node Map Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NodeListSearchHandler(c *gin.Context) {
	req := &GeneralRequest{}
	ret := &Response{}

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
	if !models.NodeReady {
		ret.Return(nil, CodeSuccess)
		JsonResponse(c, ret)
		return
	}

	rets, err := models.NodeListSearch(req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString("Get Node List Search Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func NodeDetailHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("id")
	id := converter.StrToInt64(idStr)
	if id <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.NodeDetail(id)
	if err != nil {
		ret.ReturnFailureString("Get Node Detail Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

func GetNodeBlockListHandler(c *gin.Context) {
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
	rets, err := models.GetNodeBlockList(req.Search, req.Page, req.Limit, req.Order)
	if err != nil {
		ret.ReturnFailureString("Get Node Block List Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeVoteChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.NodeVoteChange)
	if err != nil {
		ret.ReturnFailureString("Get Node Vote Chart From Redis Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeStakingChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.NodeStakingChange)
	if err != nil {
		ret.ReturnFailureString("Get Node staking Chart From Redis Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeRegionChartHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetRedis(models.NodeRegion)
	if err != nil {
		ret.ReturnFailureString("Get Node Region Chart From Redis Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetNodeStatisticalChangeHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetNodeChangeChart()
	if err != nil {
		ret.ReturnFailureString("get Node Statistical Change Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetHonorNodeMapHandler(c *gin.Context) {
	ret := &Response{}
	rets, err := models.GetHonorNodeMapFromRedis()
	if err != nil {
		ret.ReturnFailureString("Get Honor Node Map Failed")
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}
