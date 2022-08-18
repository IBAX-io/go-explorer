/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

type CountResponse struct {
	Count int64 `json:"count"`
}

type TransactionsStatusResponse struct {
	CountResponse
	List []TransactionStatus `json:"list"`
}

type TransactionsResponse struct {
	CountResponse
	List []*Transaction `json:"list"`
}

//GetQueueTransactions get queue transactions from node
func GetQueueTransactionsCount(reqUrl string) (int64, error) {
	reqNew := make(url.Values)
	reqNew["table_name"] = []string{"transactions"}
	reqNew["page"] = []string{"1"}
	reqNew["limit"] = []string{"1"}

	data, err := sendPostFormRequest(reqUrl, reqNew)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetQueueTransactionsCount err:")
		return 0, err
	}

	var rets CountResponse
	if err := json.Unmarshal(data, &rets); err != nil {
		log.WithFields(log.Fields{"error": err, "data": string(data)}).Error("GetQueueTransactionsCount json err:")
		return 0, err
	}

	return rets.Count, nil
}

//GetTransactionsStatus get transactions status from node
func GetTransactionsStatus(reqUrl string, t1 int64, page, limit int) (*[]TransactionStatus, error) {
	reqNew := make(url.Values)
	reqNew["table_name"] = []string{"transactions_status"}
	if t1 > 0 {
		reqNew["where"] = []string{fmt.Sprintf("time >=%d", t1)}
	}
	reqNew["order"] = []string{"time desc"}
	reqNew["page"] = []string{fmt.Sprintf("%d", page)}
	reqNew["limit"] = []string{fmt.Sprintf("%d", limit)}

	data, err := sendPostFormRequest(reqUrl, reqNew)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetTransactionsStatus err:")
		return nil, err
	}

	var rets TransactionsStatusResponse
	if err := json.Unmarshal(data, &rets); err != nil {
		log.WithFields(log.Fields{"error": err, "data": string(data)}).Error("GetTransactionsStatus json err:")
		return nil, err
	}
	if rets.Count == 0 {
		return nil, nil
	}

	//if int64(page*limit) >= rets.Count {
	//	return &rets.List, nil
	//} else {
	//	return GetTransactionsStatus(reqUrl, t1, page+1, limit)
	//}

	if int64(page*limit) >= rets.Count {
		return &rets.List, nil
	} else {
		ret, err := GetTransactionsStatus(reqUrl, t1, page+1, limit)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			return nil, nil
		}
		retLen := *ret
		for i := 0; i < len(retLen); i++ {
			rets.List = append(rets.List, retLen[i])
		}
	}
	return &rets.List, nil
}

//GetQueueTransactions get queue transactions from node
func GetQueueTransactions(reqUrl string, page, limit int) (*TransactionsResponse, error) {
	reqNew := make(url.Values)
	reqNew["table_name"] = []string{"transactions"}
	reqNew["page"] = []string{strconv.Itoa(page)}
	reqNew["limit"] = []string{strconv.Itoa(limit)}

	data, err := sendPostFormRequest(reqUrl, reqNew)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetQueueTransactions err:")
		return nil, err
	}

	var rets TransactionsResponse
	if err := json.Unmarshal(data, &rets); err != nil {
		log.WithFields(log.Fields{"error": err, "data": string(data)}).Error("GetQueueTransactions json err:")
		return nil, err
	}
	//fmt.Printf("node transactions len:%d,count:%d\n", len(rets.List),rets.Count)
	if rets.Count == 0 {
		return nil, nil
	}

	if int64(page*limit) >= rets.Count {
		return &rets, nil
	} else {
		ret, err := GetQueueTransactions(reqUrl, page+1, limit)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(ret.List); i++ {
			rets.List = append(rets.List, ret.List[i])
		}
		rets.Count = ret.Count
	}
	return &rets, nil
}
