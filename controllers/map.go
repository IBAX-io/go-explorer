/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"os"
	"path"
)

func GetMapInfo(c *gin.Context) {
	ret := &Response{}
	configPath := path.Join(conf.GetEnvConf().ConfigPath, "map")
	type reqRequest struct {
		Search string `json:"search"`
	}
	req := &reqRequest{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	search := req.Search
	var filePath string
	if search == "world" {
		filePath = path.Join(configPath, "world.json")
	} else {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		ret.ReturnFailureString("readFile failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(string(data), CodeSuccess)
	JsonResponse(c, ret)
}
