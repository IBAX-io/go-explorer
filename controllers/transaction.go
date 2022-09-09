/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"unicode/utf8"
)

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce      json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func Get_transaction(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}
	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	ret, num, err := services.GetGroupTransactionStatus(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order)
	if err == nil && ret != nil {
		rb.Data = ret
		rb.Total = num
		//rb.Page_size = req.Params.Page_size
		//rb.Current_page = req.Params.Current_page
		GenResponse(c, req.Head, rb)
	} else {
		if err != nil {
			rb.Retinfo = err.Error()
		}
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}
}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func Get_transaction_history(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	ret, num, err := services.GetGroupTransactionHistory(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order)
	if err == nil {
		rb.Data = ret
		rb.Total = num

		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func Get_transaction_block(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	ret, num, err := models.Get_Group_TransactionBlock(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order, 0)
	if err == nil && ret != nil {
		rb.Data = ret
		rb.Total = num
		GenResponse(c, req.Head, rb)
	} else {
		if err != nil {
			rb.Retinfo = err.Error()
		}
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func Get_transaction_details(c *gin.Context) {

	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		//rb.Body.Retinfo = err.Error()
		GenResponse(c, req.Head, rb)
	}

	ret, err := services.GetTransactionDetailedInfoHash(req.Params.Hash)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = ret
		rb.PageSize = req.Params.PageSize
		rb.CurrentPage = req.Params.CurrentPage
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func GetTransactionDetails(c *gin.Context) {

	ret := &Response{}
	hash := c.Param("hash")
	if hash == "" || utf8.RuneCountInString(hash) > 100 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := services.GetTransactionDetailedInfoHash(hash)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

func SearchHash(c *gin.Context) {
	ret := &Response{}
	hash := c.Param("hash")
	if hash == "" || utf8.RuneCountInString(hash) > 100 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := models.SearchHash(hash)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

func GetTransactionHead(c *gin.Context) {

	ret := &Response{}
	hash := c.Param("hash")
	if hash == "" || utf8.RuneCountInString(hash) > 100 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := services.GetTransactionHeadInfoHash(hash)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func Get_Find_history(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	rb.PageSize = req.Params.PageSize
	rb.CurrentPage = req.Params.CurrentPage
	rb.RetDataType = req.Params.SearchType

	//ts := &History{}
	//ret, num, err := ts.GetWallets(req.Params.Current_page, req.Params.Page_size, req.Params.Wallet, req.Params.SearchType)
	ret, num, total, err := services.GetGroupTransactionWallet(req.Params.CurrentPage, req.Params.PageSize, req.Params.Wallet, req.Params.SearchType)
	if err == nil {
		rb.Data = ret
		rb.Total = num
		rb.Sum = total
		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func Get_Find_Ecosytemhistory(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	rb.Ecosystem = req.Params.Ecosystem
	rb.PageSize = req.Params.PageSize
	rb.CurrentPage = req.Params.CurrentPage
	rb.RetDataType = req.Params.SearchType

	ret, num, total, err := services.GetGroupTransactionEcosystemWallet(req.Params.Ecosystem, req.Params.CurrentPage, req.Params.PageSize, req.Params.Wallet, req.Params.SearchType)
	if err == nil {
		rb.Data = ret
		rb.Total = num
		rb.Sum = total
		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func GetAccountTransactionHistory(c *gin.Context) {
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 || req.Wallet == "" {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	lt := &models.LogTransaction{}
	rets, err := lt.GetEcosystemAccountTransaction(req.Ecosystem, req.Page, req.Limit, req.Wallet, req.Order, req.Where)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

// @tags
// @Description
// @Summary
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /auth/admin/{id} [get]
func Get_Find_Wallethistory(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	//rb.Page_size = req.Params.Page_size
	rb.Ecosystem = req.Params.Ecosystem
	rb.Wallet = req.Params.Wallet

	ret, err := services.GetGroupWalletHistory(req.Params.Ecosystem, req.Params.Wallet)
	if err == nil {
		rb.Data = ret
		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary  Find a list of all currencies under the account
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /api/get_wallettotal [post]
func Get_Wallet_Total(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	//rb.Page_size = req.Params.Page_size
	//rb.Ecosystem = req.Params.Ecosystem
	rb.Wallet = req.Params.Wallet
	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	total, page, ret, err := services.GetGroupWalletTotal(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order, req.Params.Wallet)
	if err == nil {
		rb.Data = ret
		rb.Total = total
		rb.PageSize = page
		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary  Find a list of all currencies under the account
// @Accept   json
// @Produce  json
// @Success  200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router   /api/get_wallettotal [post]
func GetAccountDetailEcosystem(c *gin.Context) {
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ts := &models.Key{}
	rets, err := ts.GetWalletTotalEcosystem(req.Page, req.Limit, req.Order, req.Wallet)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetAccountDetailBasisEcosystem(c *gin.Context) {
	ret := &Response{}
	//wallet
	ts := &models.Key{}
	wallet := c.Param("wallet")
	rets, err := ts.GetWalletTotalBasisEcosystem(wallet)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetAccountDetailBasisTokenChange(c *gin.Context) {
	ret := &Response{}
	wallet := c.Param("wallet")
	rets, err := models.GetWalletTokenChangeBasis(wallet)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

// @tags         common_transaction_search
// @Description  common_transaction_search
// @Summary      common_transaction_search
// @Accept       json
// @Produce  json
// @Success      200  {string}  json  "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router       /auth/admin/{id} [get]
func CommonTransactionSearch(c *gin.Context) {

	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.ReqType != 0 && req.ReqType != 1 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	ts := &models.BlockTxDetailedInfoHex{}
	rets, err := ts.GetCommonTransactionSearch(req.Page, req.Limit, req.Search, req.Order, req.ReqType)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

/*
// @tags  transaction history
// @Description transaction history
// @Summary transaction history
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_Transaction_queue(c *gin.Context) {
	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	rb.PageSize = req.Params.PageSize
	rb.CurrentPage = req.Params.CurrentPage
	rb.RetDataType = req.Params.SearchType

	ret, num, err := models.GetTransactionpages(req.Params.CurrentPage, req.Params.PageSize)
	if err == nil {
		rb.Data = ret
		rb.Total = num

		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}
*/

func GetNodeTransactionListHandler(c *gin.Context) {
	ret := &Response{}
	req := &DataBaseFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	if req.Page < 1 || req.Limit < 1 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, count, err := models.GetNodeTransactionList(req.Limit, req.Page)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	var bk models.Block
	f, err := bk.GetMaxBlock()
	if err != nil || !f {
		ret.ReturnFailureString("get max block id failed")
		JsonResponse(c, ret)
		return
	}
	type listResponse struct {
		Count     int64                     `json:"count"`
		LastBlock int64                     `json:"last_block"`
		Rets      *[]models.TransactionList `json:"rets"`
		Page      int                       `json:"page"`
		Limit     int                       `json:"limit"`
	}
	var list listResponse
	list.Count = count
	list.Rets = rets
	list.LastBlock = bk.ID
	list.Page = req.Page
	list.Limit = req.Limit

	ret.Return(&list, CodeSuccess)
	JsonResponse(c, ret)

}

func GetUtxoTransactionDetails(c *gin.Context) {
	ret := &Response{}
	hash := c.Param("hash")
	if hash == "" || utf8.RuneCountInString(hash) > 100 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets, err := services.GetUtxoTransactionDetailedInfo(hash)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}
