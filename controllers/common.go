/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"net/http"

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
	name := c.Param("name")
	if name == "" {
		ret.ReturnFailureString("request params ivalid")
		JsonResponse(c, ret)
		return
	}

	rets, err := models.GetRedisByName(name)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}
