/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"errors"
	"github.com/IBAX-io/go-explorer/models"
	log "github.com/sirupsen/logrus"
)

type NodeBlockData struct {
	Data *[]models.Block
}

var (
	Sqlite_MaxBlockid int64
	bgOnWorkRun       uint32
)

func WorkDealBlock() error {
	//if atomic.CompareAndSwapUint32(&bgOnWorkRun, 0, 1) {
	//	defer atomic.StoreUint32(&bgOnWorkRun, 0)
	//} else {
	//	return nil
	//}

	var bm models.BlockID
	//var bc models.BlockID
	var bk models.Block
	fm, errm := bm.GetbyName(models.MintMax)
	if errm != nil {
		if (errm.Error() == "redis: nil" || errm.Error() == "EOF") && !fm {
			bm.ID = 0
			bm.Time = 1
			bm.Name = models.MintMax
			err := bm.InsertRedis()
			if err != nil {
				return err
			}
		} else {
			return errm
		}
	}
	fc, err := bk.GetMaxBlock()
	if err != nil {
		return err
	}
	//fc, errc := bc.GetbyName(models.ChainMax)
	//if errc != nil {
	//	return errc
	//}
	//if !fc || !fm {
	//	return errors.New("mint or chain block id not found")
	//}
	if !fm || !fc {
		return errors.New("mint or max block id not found")
	}

	//count := bc.ID - bm.ID
	count := bk.ID - bm.ID
	sc := bm.ID + 1
	elen := sc + count
	//fmt.Printf("sc:%d,elen:%d,count:%d\n", sc, elen, count)

	for i := sc; i <= elen; i++ {
		bid := i

		var mc models.ScanOut
		f, err := mc.Get(bid)
		if err != nil {
			log.Info("get scan bid:", bid, " error:", err.Error())
			continue
		}

		if f {
			err = mc.Insert_Redis()
			if err != nil {
				log.Info(err.Error())
			}
			if bid > bm.ID {
				bm.ID = bid
				bm.Time = mc.Time
				bm.InsertRedis()
			}
		}
		err = mc.Del(i - 1)
		if err != nil {
			log.Info("redis delete scan-:", bid, " failed:", err.Error())
		}
	}

	return nil
}

func SyncDealBlock() {
	var bm models.BlockID
	fm, errm := bm.GetbyName(models.MintMax)
	if errm != nil {
		if (errm.Error() == "redis: nil" || errm.Error() == "EOF") && !fm {
			return
		} else {
			log.Info("SyncDealBlock redis err:", errm.Error())
			return
		}
	}
	if fm {
		for i := bm.ID; i > 0; i-- {
			bid := i

			var mc models.ScanOut
			f, _ := mc.Get(bid)
			if f {
				err := mc.Del(i)
				if err != nil {
					log.Info("Sync Deal Block Redis Delete Failed:", errm.Error())
				}
			}
			//time.Sleep(time.Microsecond * 10)
		}
	}
	return
}
