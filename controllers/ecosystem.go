/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"unicode/utf8"

	//"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetPlatformEcosystemParam(c *gin.Context) {
	var (
		//params []models.AppParam
		rets models.SystemParameterResult
	)
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	//eid := strconv.FormatInt(req.Ecosystem, 10)
	//err := models.GetALL(eid+"_app_params", "", &rets.Rets)
	//var ap models.AppParam
	//num, rs, err := ap.FindAppParameters(req.Ecosystem, req.Page, req.Limit, req.Search, req.Order)
	//if err != nil {
	//	ret.ReturnFailureString(err.Error())
	//	JsonResponse(c, ret)
	//	return
	//}

	//len([]rune(req.Search))
	if req.Page <= 0 || req.Limit <= 0 || utf8.RuneCountInString(req.Search) > 100 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	var ap models.PlatformParameter
	num, rs, err := ap.FindAppParameters(req.Page, req.Limit, req.Search, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	rets.Total = int64(num)
	rets.Page = req.Page
	rets.Limit = req.Limit
	rets.Rets = rs
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcosystemParam(c *gin.Context) {
	var (
		rets models.EcosystemParameterResult
	)
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	if req.Page <= 0 || req.Limit <= 0 || utf8.RuneCountInString(req.Search) > 100 || req.Ecosystem == 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	var ap models.StateParameter
	num, rs, err := ap.FindStateParameters(req.Page, req.Limit, req.Search, req.Order, req.Ecosystem)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	rets.Total = num
	rets.Page = req.Page
	rets.Limit = req.Limit
	rets.Rets = rs
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcosystemList(c *gin.Context) {
	var req EcosytemTranscationHistoryFind
	ret := &Response{}
	var rets models.EcosystemTotalResult

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	rets.Limit = req.Limit
	rets.Page = req.Page
	var eco models.Ecosystem
	total, list, err := eco.GetEcoSystemList(req.Limit, req.Page, req.Order, req.Where)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	} else {
		rets.Total = total
		rets.Rets = list
		ret.Return(&rets, CodeSuccess)
		JsonResponse(c, ret)
	}
}

func GetEcosystemBasis(c *gin.Context) {
	ret := &Response{}
	var rets models.EcosystemTotalResult

	basisEcoLib, err := models.GetEcoLibsChartData()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	} else {
		rets.Total = 1
		rets.Sysecosy = basisEcoLib
		ret.Return(&rets, CodeSuccess)
		JsonResponse(c, ret)
	}
}

func GetEcosystemDetailInfoHandler(c *gin.Context) {
	var req GeneralRequest
	ret := &Response{}
	rets := &models.EcosystemDetailInfoResponse{}

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Search == nil {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	detailInfo, err := models.GetEcosystemDetailInfo(req.Search)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	rets = detailInfo
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcosystemDetailTxHandler(c *gin.Context) {
	var req EcosytemTranscationHistoryFind
	ret := &Response{}
	var rets models.GeneralResponse
	var his models.LogTransaction
	var txList []models.EcosystemTxList

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Ecosystem <= 0 || req.Limit <= 0 || req.Page <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	list, total, err := his.GetEcosystemTransactionFind(req.Ecosystem, req.Page, req.Limit, req.Order, req.Search, req.Where)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ts := *list
	txList = make([]models.EcosystemTxList, len(ts))
	for i := 0; i < len(ts); i++ {
		txList[i].BlockId = ts[i].Block
		txList[i].Time = models.MsToSeconds(ts[i].Timestamp)
		txList[i].Contract = ts[i].ContractName
		if txList[i].Contract == "" {
			txList[i].Contract = models.GetUtxoTxContractNameByHash(ts[i].Hash)
		}
		txList[i].Address = converter.AddressToString(ts[i].Address)
		txList[i].Hash = hex.EncodeToString(ts[i].Hash)
		txList[i].Status = ts[i].Status
	}
	rets.Total = total
	rets.List = txList
	rets.Page = req.Page
	rets.Limit = req.Limit
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)

}

func EcosystemSearchHandler(c *gin.Context) {
	req := &EcosytemTranscationHistoryFind{}
	ret := &Response{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	rets, err := models.EcosystemSearch(req.Search, req.Wallet)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcosystemDetailTokenHandler(c *gin.Context) {
	var req EcosytemTranscationHistoryFind
	ret := &Response{}
	var key models.Key

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Ecosystem <= 0 || req.Limit <= 0 || req.Page <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	rets, err := key.GetEcosystemTokenSymbolList(req.Page, req.Limit, req.Order, req.Ecosystem)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(&rets, CodeSuccess)
	JsonResponse(c, ret)
}

func GetEcosystemDetailMemberHandler(c *gin.Context) {
	var req EcosytemTranscationHistoryFind
	ret := &Response{}
	var key models.Key

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if req.Ecosystem <= 0 || req.Limit <= 0 || req.Page <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	rets, err := key.GetEcosystemDetailMemberList(req.Page, req.Limit, req.Order, req.Ecosystem)
	if err != nil {
		ret.ReturnFailureString(err.Error())
	} else {
		ret.Return(&rets, CodeSuccess)
	}
	JsonResponse(c, ret)
}

func GetEcosystemDatabaseHandler(c *gin.Context) {
	var req EcosytemTranscationHistoryFind
	ret := &Response{}

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	if req.Ecosystem <= 0 || req.ReqType < 1 || req.ReqType > 3 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}
	if len([]rune(req.Search)) > 100 || len([]rune(req.Order)) > 100 {
		ret.ReturnFailureString("request params len failed")
		JsonResponse(c, ret)
		return
	}

	rets, err := models.GetEcosystemDatabase(req.Page, req.Limit, req.ReqType, req.Ecosystem, req.Search, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
	} else {
		ret.Return(&rets, CodeSuccess)
	}
	JsonResponse(c, ret)
}

func GetEcosystemAppHandler(c *gin.Context) {
	var req EcosytemTranscationHistoryFind
	ret := &Response{}

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	if req.Ecosystem <= 0 || req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	rets, err := models.GetEcosystemApp(req.Page, req.Limit, req.Ecosystem, req.AppId, req.Order, req.Search)
	if err != nil {
		ret.ReturnFailureString(err.Error())
	} else {
		ret.Return(&rets, CodeSuccess)
	}
	JsonResponse(c, ret)
}

func GetEcosystemAppExportHandler(c *gin.Context) {
	ret := &Response{}
	idStr := c.Param("id")
	id := converter.StrToInt64(idStr)
	if id <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	rets, err := models.EcosystemAppExport(id)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	data, err := json.Marshal(rets)
	if err != nil {
		ret.ReturnFailureString("EcosystemApp export json failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	c.Header("Content-Type", http.DetectContentType(data))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, rets.Name))
	c.Header("Access-Control-Allow-Origin", "*")
	_, err = c.Writer.Write(data)
	IndentedJsonResponse(c, rets)
	//if err != nil {
	//	ret.ReturnFailureString("Get Ecosystem App Export Handler Write Error:" + err.Error())
	//	JsonResponse(c, ret)
	//	return
	//}

}

func GetEcosystemAttachmentHandler(c *gin.Context) {
	var req EcosytemTranscationHistoryFind
	ret := &Response{}

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	if req.Page <= 0 || req.Limit <= 0 {
		ret.ReturnFailureString("request params invalid")
		JsonResponse(c, ret)
		return
	}

	var bin models.Binary
	rets, err := bin.FindAppNameByEcosystem(req.Page, req.Limit, req.Order, req.Ecosystem, req.Where)
	if err != nil {
		ret.ReturnFailureString(err.Error())
	} else {
		ret.Return(&rets, CodeSuccess)
	}
	JsonResponse(c, ret)
}

func GetEcosystemAttachmentExportHandler(c *gin.Context) {
	ret := &Response{}

	hash := c.Param("hash")

	if hash == "" || utf8.RuneCountInString(hash) > 100 {
		ret.ReturnFailureString("Request params invalid")
		JsonResponse(c, ret)
		return
	}

	fileName, id, err := models.GetFileNameByHash(hash)
	if err != nil {
		ret.ReturnFailureString("Get attachment failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}
	if fileName == "" {
		ret.ReturnFailureString("Get attachment failed:File doesn't not exist")
		JsonResponse(c, ret)
		return
	}
	//Save the file to the local. If the file does not exist, search for the file from the database
	if !models.IsExist(UploadDir + fileName) {
		fileName, err = models.LoadFile(id)
		if err != nil {
			ret.ReturnFailureString("loadFile failed:" + err.Error())
			JsonResponse(c, ret)
			return
		}
		if fileName == "" {
			ret.ReturnFailureString("hash doesn't not exist")
			JsonResponse(c, ret)
			return
		}

	}

	//var bin models.Binary
	//f, err := bin.GetByHash(hash)
	//if err != nil {
	//	ret.ReturnFailureString("Attachment Find Failed:" + err.Error())
	//	JsonResponse(c, ret)
	//	return
	//}
	//if !f {
	//	ret.ReturnFailureString("Attachment file doesn't not exist")
	//	JsonResponse(c, ret)
	//	return
	//}

	data, err := os.ReadFile(UploadDir + fileName)
	if err != nil {
		ret.ReturnFailureString("Get attachment readfile failed:" + err.Error())
		JsonResponse(c, ret)
		return
	}

	if !models.CompareHash(data, hash) {
		ret.ReturnFailureString("Hash is incorrect")
		JsonResponse(c, ret)
		return
	}

	//c.Header("Content-Type", bin.MimeType)
	c.Header("Content-Type", http.DetectContentType(data))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Header("Access-Control-Allow-Origin", "*")
	_, err = c.Writer.Write(data)
	if err != nil {
		ret.ReturnFailureString("Attachment File Write Error:" + err.Error())
		JsonResponse(c, ret)
		return
	}
}
