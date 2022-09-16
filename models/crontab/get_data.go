/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package crontab

import (
	"fmt"
	"github.com/IBAX-io/go-explorer/models/crontab/get_data"
)

type DataType int32

const (
	HonorNodeData     DataType = 0
	HistoryData       DataType = 1
	RealTimeData      DataType = 2
	LoadContractsData DataType = 3
	ChartData         DataType = 4
)

type GetDataProvider interface {
	SendSignal()
	ReceiveSignal()
}

func NewGetDataProvider(p DataType) GetDataProvider {
	switch p {
	case HonorNodeData:
		return &get_data.HonorNode{}
	case HistoryData:
		return &get_data.History{}
	case RealTimeData:
		return &get_data.RealTime{}
	case LoadContractsData:
		return &get_data.LoadContracts{}
	case ChartData:
		return &get_data.Chart{}
	}
	panic(fmt.Errorf("get data [%v] is not supported yet", p))
}
