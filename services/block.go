/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (

	//"strings"
	//	"fmt"

	"github.com/IBAX-io/go-explorer/models"
	//"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
)

func GetGroupBlockLists() (*models.BlockListHeaderResponse, error) {
	ret, err := models.GetBlockListFromRedis()
	return ret, err

}

func GetGroupBlockTpsLists() (*[]models.ScanOutBlockTransactionRet, error) {
	ret, err := models.GetBlockTpsListsFromRedis()
	return ret, err
}
