/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package buffer

import (
	"fmt"
	"github.com/IBAX-io/go-explorer/models"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"sync/atomic"
)

type BufferType struct {
	atomic uint32
}

var (
	ecoRealTime  = &BufferType{atomic: 0}
	ecoHistory   = &BufferType{atomic: 0}
	dataHistory  = &BufferType{atomic: 0}
	dataRealTime = &BufferType{atomic: 0}

	RefreshList sync.Map
)

func StartServer(reqType *BufferType) {
	if atomic.CompareAndSwapUint32(&reqType.atomic, 0, 1) {
		defer atomic.StoreUint32(&reqType.atomic, 0)
	} else {
		return
	}
	if reqType == nil {
		return
	}

	switch reqType {
	case dataHistory, dataRealTime:
		err := dataChartServer(reqType)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("[buffer server]Get data chart failed")
		}
	case ecoRealTime, ecoHistory:
		err := ecosystemChartServer(reqType)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("[buffer server]Get ecosystem chart failed")
		}
	}

}

func RefreshServer() {
	models.RefreshRequest = make(chan models.RefreshObject, 100)
	for {
		select {
		case req := <-models.RefreshRequest:
			if req.Key == "" {
				continue
			}
			key := req.Key + strconv.FormatInt(req.Ecosystem, 10)
			if req.Cmd == 1 {
				RefreshList.Delete(key)
				continue
			}
			_, exist := RefreshList.Load(key)
			if !exist {
				RefreshList.Store(key, 0)
				go models.RefreshChartDaemons(req.Key, req.Ecosystem)
			}
		}
	}
}

func GetBufferType(reqType int) *BufferType {
	switch reqType {
	case 1:
		return ecoRealTime
	case 2:
		return ecoHistory
	case 3:
		return dataRealTime
	case 4:
		return dataHistory
	}
	panic(fmt.Errorf("not support buffer server type[%d]", reqType))
}
