/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Flagdir = "./flag/"
var UploadDir = "./upload/"

func init() {
	if err := utils.MakeDirectory(Flagdir); err != nil {
		log.WithFields(log.Fields{"error": err, "type": consts.IOError, "dir": Flagdir}).Error("can't create temporary directory")
	}
	if err := utils.MakeDirectory(UploadDir); err != nil {
		log.WithFields(log.Fields{"error": err, "type": consts.IOError, "dir": UploadDir}).Error("can't create temporary directory")
	}
}

func GetNftMinerFileHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("id")
	id := converter.StrToInt64(idStr)
	if id <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	var items models.NftMinerItems
	f, err := items.GetById(id)
	if err != nil {
		ret.ReturnFailureString("request error:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if !f {
		ret.ReturnFailureString("unknown nft Miner id:" + idStr)
		JsonResponse(c, ret)
		return
	}
	data, err := items.ParseSvgParams()
	if err != nil {
		ret.ReturnFailureString("Get Nft Miner File Failed")
		JsonResponse(c, ret)
		return
	}
	c.Header("Content-Type", "image/svg+xml;utf8")
	c.Header("Access-Control-Allow-Origin", "*")
	_, err = c.Writer.Write([]byte(data))
	if err != nil {
		ret.ReturnFailureString("Get Nft Miner File Handler Write Error:" + err.Error())
		JsonResponse(c, ret)
		return
	}

}
