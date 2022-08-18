/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"net/http"

	"github.com/IBAX-io/go-ibax/packages/converter"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
)

func JsonResponse(c *gin.Context, body *Response) {
	c.JSON(http.StatusOK, body)
}

func PureJsonResponse(c *gin.Context, body *Response) {
	c.PureJSON(http.StatusOK, body)
}

//IndentedJsonResponse Json Format
func IndentedJsonResponse(c *gin.Context, body any) {
	c.IndentedJSON(http.StatusOK, body)
}

//GenResponse genrate reponse ,json format
func GenResponse(c *gin.Context, head *RequestHead, body *ResponseBoby) {
	c.JSON(http.StatusOK, gin.H{
		"body": body,
		"head": head,
	})
}

// @tags         ecosystem
// @Description  ecosystem
// @Summary      ecosystem
// @Accept       json
// @Produce      json
// @Success      200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router       /auth/admin/{id} [get]
func GetRedisKey(c *gin.Context) {
	ret := &Response{}
	id := c.Param("id")
	count := converter.StrToInt64(id)
	var scanout models.ScanOut
	f, err := scanout.Get_Redis(count)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if f {
		ret.Return(scanout, CodeSuccess)
		JsonResponse(c, ret)
		return
	} else {
		ret.ReturnFailureString("not found key in redis:" + models.ScanOutStPrefix + id)
		JsonResponse(c, ret)
		return
	}

}
