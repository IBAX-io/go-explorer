package controllers

import (
	"fmt"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type tokenPriceInfo struct {
	Ecosystem  int64  `json:"ecosystem"`
	PriceInUsd string `json:"price_in_usd"`
}

func GetTokenPriceHandler(c *gin.Context) {
	ret := &Response{}
	str := c.Query("ecosystems")
	if str == "" {
		ret.ReturnFailureString("ecosystems cannot be empty")
		JsonResponse(c, ret)
		return
	}
	ecosystemList := strings.Split(str, ",")
	var ecosystems = make([]int64, len(ecosystemList))
	for i, v := range ecosystemList {
		var err error
		ecosystems[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			ret.ReturnFailureString(fmt.Sprintf("ecosytem %s invalid", v))
			JsonResponse(c, ret)
			return
		}
		if ecosystems[i] <= 0 {
			ret.ReturnFailureString(fmt.Sprintf("ecosytem %d invalid", ecosystems[i]))
			JsonResponse(c, ret)
			return
		}
		if models.EcoNames.Get(ecosystems[i]) == "" || models.EcoDigits.GetInt(ecosystems[i], 999) == 999 {
			ret.ReturnFailureString(fmt.Sprintf("the ecosystem %d does not exist or there is no token information", ecosystems[i]))
			JsonResponse(c, ret)
			return
		}
	}

	prices, err := models.GetTokenPrices(ecosystems)
	if err != nil {
		ret.ReturnFailureString(fmt.Sprintf("get prices %s", err))
		JsonResponse(c, ret)
		return
	}
	list := make([]tokenPriceInfo, len(prices))
	for i, ecosystem := range ecosystems {
		list[i].Ecosystem = ecosystem
		list[i].PriceInUsd = prices[ecosystem]
	}
	ret.Return(list, CodeSuccess)
	JsonResponse(c, ret)
	return
}

func GetTokenLogoHandler(c *gin.Context) {
	ret := &Response{}
	chainIdStr := c.Query("chainId")
	tokenAddress := c.Query("tokenAddress")
	ecosystem := c.Query("ecosystem")
	if (chainIdStr == "" || tokenAddress == "") && ecosystem == "" {
		ret.ReturnFailureString("param cannot be empty")
		JsonResponse(c, ret)
		return
	}
	var logoUri string
	ecosystemId, err := strconv.ParseUint(ecosystem, 10, 64)
	if err != nil {
		ret.ReturnFailureString("ecosystem " + ecosystem + " invalid")
		JsonResponse(c, ret)
		return
	}
	if ecosystemId > 0 {
		info := models.Info.Get(int64(ecosystemId))
		if info.ChainName != "" {
			logoUri = info.LogoURI
		} else {
			if info.LogoHash != "" {
				logoUri = conf.GetEnvConf().Url.Base + models.ApiPath + "get_eco_attachment_export/" + info.LogoHash
			}
		}
	} else {
		chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
		if err != nil {
			ret.ReturnFailureString("chainId " + chainIdStr + " invalid")
			JsonResponse(c, ret)
			return
		}
		if chainId == 0 {
			ret.ReturnFailureString("param invalid")
			JsonResponse(c, ret)
			return
		}
		token := &models.TokensInfo{}
		logoUri, _, err = token.GetLogoURI(chainIdStr, tokenAddress)
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		}
	}

	ret.Return(logoUri, CodeSuccess)
	JsonResponse(c, ret)
	return
}
