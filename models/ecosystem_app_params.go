/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"errors"
	"github.com/IBAX-io/go-ibax/packages/smart"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	"github.com/IBAX-io/go-ibax/packages/types"
)

type ContractsParams struct {
	Id         int64  `json:"id,omitempty"`
	Type       string `json:"Type"`
	Name       string `json:"Name"`
	Value      string `json:"Value"`
	Conditions string `json:"Conditions"`
}

type PageParams struct {
	Id         int64  `json:"id,omitempty"`
	Type       string `json:"Type"`
	Name       string `json:"Name"`
	Value      string `json:"Value"`
	Conditions string `json:"Conditions"`
	Menu       string `json:"Menu"`
}

type SnippetsParams struct {
	Id         int64  `json:"id,omitempty"`
	Type       string `json:"Type"`
	Name       string `json:"Name"`
	Value      string `json:"Value"`
	Conditions string `json:"Conditions"`
}

type TableParams struct {
	Id          int64  `json:"id,omitempty"`
	Type        string `json:"Type"`
	Name        string `json:"Name"`
	Columns     string `json:"Columns"`
	Permissions string `json:"Permissions"`
	Conditions  string `json:"conditions"`
}

type AppParams struct {
	Id         int64  `json:"id,omitempty"`
	Type       string `json:"Type"`
	Name       string `json:"Name"`
	Conditions string `json:"Conditions"`
	Value      string `json:"Value"`
}

type MenuParams struct {
	Type       string `json:"Type"`
	Name       string `json:"Name"`
	Conditions string `json:"Conditions"`
	Value      string `json:"Value"`
}

type LanguagesParams struct {
	Type  string `json:"Type"`
	Name  string `json:"Name"`
	Trans string `json:"trans"`
}

func getAppTableByApp(appID int64, ecosystemID int64) ([]sqldb.Table, error) {
	var result []sqldb.Table
	err := GetDB(nil).Select("id,name,permissions,conditions,ecosystem,columns").Where("app_id = ? and ecosystem = ?", appID, ecosystemID).Find(&result).Error
	return result, err
}

func getAppParamsByApp(appID int64, ecosystemID int64) ([]sqldb.AppParam, error) {
	var result []sqldb.AppParam
	err := GetDB(nil).Select("id,name,conditions,value").Where("app_id = ? and ecosystem = ?", appID, ecosystemID).Find(&result).Error
	return result, err
}

func getLanguagesParamsByApp(appID int64, ecosystemID int64) ([]sqldb.Language, error) {
	var result []sqldb.Language
	err := GetDB(nil).Select("name,res").Where("app_id = ? and ecosystem = ?", appID, ecosystemID).Find(&result).Error
	return result, err
}

func getPageByApp(appID int64, ecosystemID int64) ([]sqldb.Page, error) {
	var result []sqldb.Page
	err := GetDB(nil).Select("id,name,value,menu,conditions").Where("app_id = ? and ecosystem = ?", appID, ecosystemID).Find(&result).Error
	return result, err
}

func getSnippetsByApp(appID int64, ecosystemID int64) ([]sqldb.Snippet, error) {
	var result []sqldb.Snippet
	err := GetDB(nil).Select("id,name,value,conditions").Where("app_id = ? and ecosystem = ?", appID, ecosystemID).Find(&result).Error
	return result, err
}

func getMenuParamsByApp(ecosystemID int64) ([]sqldb.Menu, error) {
	var result []sqldb.Menu
	err := GetDB(nil).Select("name,value,conditions").Where("ecosystem = ?", ecosystemID).Find(&result).Error
	return result, err
}

func getAppContractsParams(appid, ecosystemId int64, isExport bool) ([]ContractsParams, error) {
	var ret Contract

	list, err := ret.GetByApp(appid, ecosystemId)
	if err != nil {
		return nil, err
	}
	rets := make([]ContractsParams, len(list))
	for i := 0; i < len(list); i++ {
		if !isExport {
			rets[i].Id = list[i].ID
		}
		rets[i].Type = "contracts"
		rets[i].Name = list[i].Name
		rets[i].Conditions = list[i].Conditions
		rets[i].Value = list[i].Value
	}
	return rets, nil
}

func getAppPageParams(appid, ecosystemId int64, isExport bool) ([]PageParams, error) {
	list, err := getPageByApp(appid, ecosystemId)
	if err != nil {
		return nil, err
	}
	rets := make([]PageParams, len(list))
	for i := 0; i < len(list); i++ {
		if !isExport {
			rets[i].Id = list[i].ID
		}
		rets[i].Type = "pages"
		rets[i].Name = list[i].Name
		rets[i].Conditions = list[i].Conditions
		rets[i].Value = list[i].Value
		rets[i].Menu = list[i].Menu
	}
	return rets, nil
}

func getAppSnippetsParams(appid, ecosystemId int64, isExport bool) ([]SnippetsParams, error) {
	list, err := getSnippetsByApp(appid, ecosystemId)
	if err != nil {
		return nil, err
	}
	rets := make([]SnippetsParams, len(list))
	for i := 0; i < len(list); i++ {
		if !isExport {
			rets[i].Id = list[i].ID
		}
		rets[i].Type = "snippets"
		rets[i].Name = list[i].Name
		rets[i].Conditions = list[i].Conditions
		rets[i].Value = list[i].Value
	}
	return rets, nil
}
func getColumnsWithType(table_name, column string, ecosystem int64) string {
	var colsMap any
	var result []any
	var columns []any
	var sm = smart.SmartContract{
		TxSmart: &types.SmartTransaction{
			Header: new(types.Header),
		},
		DbTransaction: sqldb.NewDbTransaction(GetDB(nil)),
	}

	sm.TxSmart.EcosystemID = ecosystem

	//column:=`{"account": "ContractAccess(\"@1SocialBind\")", "social_info": "ContractAccess(\"@1SocialBind\")"}`
	colsMap, _ = smart.JSONDecode(column)
	columns = smart.GetMapKeys(colsMap.(*types.Map))
	var i int

	for i < len(columns) {
		//if len(columns[i]) > 0 {
		var col = make(map[string]any, 0)
		col["name"] = columns[i]
		if v, ok := colsMap.(*types.Map).Get(columns[i].(string)); ok {
			col["conditions"] = v
		}
		col["type"], _ = smart.GetColumnType(&sm, table_name, col["name"].(string))

		result = append(result, col)
		//}
		i = i + 1
	}
	str, _ := smart.JSONEncode(result)
	return str
}

func getAppTableParams(appid, ecosystemId int64, isExport bool) ([]TableParams, error) {
	list, err := getAppTableByApp(appid, ecosystemId)
	if err != nil {
		return nil, err
	}
	rets := make([]TableParams, len(list))
	for i := 0; i < len(list); i++ {
		if !isExport {
			rets[i].Id = list[i].ID
		}
		rets[i].Type = "tables"
		//tableName := strconv.FormatInt(list[i].Ecosystem, 10) + "_" + list[i].Name
		permissions, _ := json.Marshal(list[i].Permissions)
		rets[i].Name = list[i].Name
		rets[i].Permissions = string(permissions)
		rets[i].Columns = getColumnsWithType(list[i].Name, list[i].Columns, list[i].Ecosystem)
		rets[i].Conditions = list[i].Conditions
	}
	return rets, nil
}

func getAppParams(appid, ecosystemId int64, isExport bool) ([]AppParams, error) {
	list, err := getAppParamsByApp(appid, ecosystemId)
	if err != nil {
		return nil, err
	}
	rets := make([]AppParams, len(list))
	for i := 0; i < len(list); i++ {
		if !isExport {
			rets[i].Id = list[i].ID
		}
		rets[i].Type = "app_params"
		rets[i].Name = list[i].Name
		rets[i].Conditions = list[i].Conditions
		rets[i].Value = list[i].Value
	}
	return rets, nil
}

func getMenuParams(ecosystemId int64) ([]MenuParams, error) {
	list, err := getMenuParamsByApp(ecosystemId)
	if err != nil {
		return nil, err
	}
	rets := make([]MenuParams, len(list))
	for i := 0; i < len(list); i++ {
		rets[i].Type = "menu"
		rets[i].Name = list[i].Name
		rets[i].Conditions = list[i].Conditions
		rets[i].Value = list[i].Value
	}
	return rets, nil
}

func getLanguagesParams(appid, ecosystemId int64) ([]LanguagesParams, error) {
	list, err := getLanguagesParamsByApp(appid, ecosystemId)
	if err != nil {
		return nil, err
	}
	rets := make([]LanguagesParams, len(list))
	for i := 0; i < len(list); i++ {
		rets[i].Type = "tables"
		rets[i].Name = list[i].Name
		rets[i].Trans = list[i].Res
	}
	return rets, nil
}

func EcosystemAppExport(appid int64) (*ExportAppInfo, error) {
	var rets ExportAppInfo
	var app Applications
	f, err := app.GetById(appid)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errors.New("app doesn't not exist")
	}
	rets.Name = app.Name
	rets.Conditions = app.Conditions
	contract, err := getAppContractsParams(app.ID, app.Ecosystem, true)
	if err != nil {
		return nil, err
	}
	page, err := getAppPageParams(app.ID, app.Ecosystem, true)
	if err != nil {
		return nil, err
	}
	if len(page) > 0 {
		menu, err := getMenuParams(app.Ecosystem)
		if err != nil {
			return nil, err
		}
		for _, value := range menu {
			rets.Data = append(rets.Data, value)
		}
	}
	params, err := getAppParams(app.ID, app.Ecosystem, true)
	if err != nil {
		return nil, err
	}
	snippets, err := getAppSnippetsParams(app.ID, app.Ecosystem, true)
	if err != nil {
		return nil, err
	}
	table, err := getAppTableParams(app.ID, app.Ecosystem, true)
	if err != nil {
		return nil, err
	}
	//languages, err := getLanguagesParams(appid, app.Ecosystem)
	//if err != nil {
	//	return nil, err
	//}
	for _, value := range contract {
		rets.Data = append(rets.Data, value)
	}
	for _, value := range page {
		rets.Data = append(rets.Data, value)
	}
	for _, value := range params {
		rets.Data = append(rets.Data, value)
	}
	for _, value := range snippets {
		rets.Data = append(rets.Data, value)
	}
	for _, value := range table {
		rets.Data = append(rets.Data, value)
	}
	//rets.Data = append(rets.Data, languages)

	return &rets, nil
}
