/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"strconv"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

type BlockTxDetailedInfo struct {
	BlockID      int64  `gorm:"not null;index:btdblockid_idx" json:"block_id"`
	Hash         string `gorm:"primary_key;not null" json:"hash"`
	ContractName string `gorm:"not null" json:"contract_name"`
	Params       string `gorm:"not null" json:"params"`
	KeyID        string `gorm:"not null" json:"key_id"`
	Time         int64  `gorm:"not null" json:"time"`
	Type         int64  `gorm:"not null" json:"type"`
	Size         int64  `gorm:"not null" json:"size"`

	Ecosystemname string `gorm:"null" json:"ecosystemname"`
	TokenSymbol   string `gorm:"null" json:"token_symbol"`
	Ecosystem     int64  `gorm:"null" json:"ecosystem"`
}

type BlockTxDetailedInfoHex struct {
	BlockID      int64  `gorm:"not null;index:blockid_idx" json:"block_id"`
	Hash         string `gorm:"primary_key;not null" json:"hash"`
	ContractName string `gorm:"not null" json:"contract_name"`
	ContractCode string `json:"contract_code"`
	//Params       map[string]any `json:"params"`
	Params string `gorm:"not null" json:"params"`
	KeyID  string `gorm:"not null" json:"key_id"`
	Time   int64  `gorm:"not null" json:"time"`
	Type   int64  `gorm:"not null" json:"type"`
	Size   int64  `gorm:"not null" json:"size"`

	Ecosystemname string          `gorm:"null" json:"ecosystemname"`
	Token_symbol  string          `gorm:"null" json:"token_symbol"`
	Ecosystem     int64           `gorm:"null" json:"ecosystem"`
	GasFee        decimal.Decimal `json:"gas_fee"`
	Amount        decimal.Decimal `json:"amount"`
	Status        int32           `json:"status"`
	Digits        int             `json:"digits"`
}

type TxListRet struct {
	TotalTx       int64  `json:"total_tx"`
	TwentyFourTx  int64  `json:"twenty_four_tx"`
	ActiveEcoLib  string `json:"active_eco_lib"`
	WeekAverageTx int64  `json:"week_average_tx"`
}

type HashTransactionResult struct {
	Total  int64                    `json:"total" `
	Page   int                      `json:"page" `
	Limit  int                      `json:"limit"`
	TxInfo *TxListRet               `json:"tx_info,omitempty"`
	Rets   []BlockTxDetailedInfoHex `json:"rets"`
}

func (bt *BlockTxDetailedInfoHex) GetByHashDb(hash string) (bool, error) {
	return bt.GetDbTxDetailedHash(hash)
}

func (bt *BlockTxDetailedInfoHex) GetByBlockIdBlockTransactionsLastDB(id int64, page int, limit int, order string) (int64, *[]BlockTxDetailedInfoHex, error) {
	var (
		ret   []BlockTxDetailedInfoHex
		total int64
		err   error
	)
	ret, err = bt.GetDbTxDetailedId(id)
	if err != nil {
		return total, &ret, err
	}
	//if len(ret) > limit {
	//	ret = ret[:limit]
	//}
	total = int64(len(ret))
	return total, &ret, err
}

func (bt *BlockTxDetailedInfoHex) GetByKeyIdBlockTransactionsLastDb(id string, page int, size int, order string) (int64, *[]BlockTxDetailedInfoHex, error) {
	var (
		ret   []BlockTxDetailedInfoHex
		total int64
		err   error
	)
	ret, total, err = bt.GetDb_txdetailedKey(id, order, size, page)
	if err != nil {
		return total, &ret, err
	}
	//total = int64(len(ret))
	return total, &ret, err

}

// GetCommonTransactionSearch is retrieving model from database
func (bt *BlockTxDetailedInfoHex) GetCommonTransactionSearch(page, limit int, search, order string, reqType int) (*HashTransactionResult, error) {
	var (
		ret HashTransactionResult
		err error
	)
	ret.Page = page
	ret.Limit = limit

	bid, err := strconv.ParseInt(search, 10, 64)
	if err == nil && bid > 0 {
		//blockid
		total, rets, err := bt.GetByBlockIdBlockTransactionsLastDB(bid, page, limit, order)
		if err != nil {
			return &ret, err
		}
		ret.Total = total
		ret.Rets = *rets
		return &ret, err
	} else {
		keyid := converter.StringToAddress(search)
		if keyid != 0 {
			//wallet

			total, rets, err := bt.GetByKeyIdBlockTransactionsLastDb(search, page, limit, order)
			if err != nil {
				return &ret, err
			}
			ret.Total = total
			ret.Rets = *rets
			return &ret, err
		} else {
			//hash
			if search == "" {
				if page == 1 && limit == 10 && reqType == 1 {
					rets, total, err := GetTransactionBlockFromRedis()
					if err != nil {
						return &ret, err
					}
					ret.Total = total
					ret.Rets = *rets
					return &ret, err
				}
				rets, total, err := Get_Group_TransactionBlock(page, limit, order, reqType)
				if err != nil {
					return &ret, err
				}
				if reqType == 0 {
					var m ScanOut
					f, err := m.GetRedisLatest()
					if err != nil {
						return &ret, err
					}
					if f {
						txRet := &TxListRet{}
						txRet.TotalTx = m.TotalTx
						txRet.TwentyFourTx = m.TwentyFourTx
						txRet.WeekAverageTx = m.WeekAverageTx
						txRet.ActiveEcoLib = m.MaxActiveEcoLib
						ret.TxInfo = txRet
					}
				}

				ret.Total = total
				ret.Rets = *rets
				return &ret, err
			} else {
				f, err := bt.GetByHashDb(search)
				if err != nil {
					return &ret, err
				}
				if f {
					ret.Total = 1
					ret.Rets = append(ret.Rets, *bt)
					return &ret, err
				} else {
					return &ret, errors.New("not found hash")
				}
			}
		}
	}
}
func Get_Group_TransactionBlock(ids int, icount int, order string, reqType int) (*[]BlockTxDetailedInfoHex, int64, error) {
	ts := &LogTransaction{}
	//bt := &models.BlockTxDetailedInfoHex{}

	//if models.GsqliteIsactive {
	//	ret, num, err := bt.Get_BlockTransactions_Sqlite(ids, icount, order)
	//	if err == nil && *ret != nil && num > 0 {
	//		//fmt.Println("Get_BlockTransactions_Sqlite  ok ids:%d icount:%d", ids, icount)
	//		return ret, num, err
	//	}
	//}
	ret, num, err := ts.GetBlockTransactions(ids, icount, order, reqType)
	//fmt.Println("Get_BlockTransactions pg  ok ids:%d icount:%d", ids, icount)
	return ret, int64(num), err
	//return nil, 0, nil
}
func (bt *BlockTxDetailedInfoHex) Get_BlockTransactions_Sqlite(page int, size int, order string) (*[]BlockTxDetailedInfoHex, int, error) {
	var (
		ret []BlockTxDetailedInfoHex
		tss []BlockTxDetailedInfo
	)

	err := conf.GetDbConn().Conn().Limit(size).Offset((page - 1) * size).Order(order).Find(&tss).Error
	if err == nil {

		for i := 0; i < len(tss); i++ {

			bh := BlockTxDetailedInfoHex{}
			//params := map[string]any
			bh.BlockID = tss[i].BlockID
			bh.ContractName = tss[i].ContractName
			bh.Hash = tss[i].Hash
			bh.KeyID = tss[i].KeyID
			//bh.Params = rt.Transactions[j].Params
			bh.Time = tss[i].Time
			bh.Type = tss[i].Type
			bh.Size = tss[i].Size

			bh.Ecosystem = tss[i].Ecosystem
			bh.Ecosystemname = tss[i].Ecosystemname
			if bh.Ecosystem == 1 {
				bh.Token_symbol = SysTokenSymbol
			} else {
				bh.Token_symbol = tss[i].TokenSymbol
			}
			//bh.Token_title = tss[i].Token_title
			//es := Ecosystem{}
			//f, err := es.Get(tss[i]..EcosystemID)
			//if f && err == nil {
			//	bh.Ecosystem = tss[i].TxHeader.EcosystemID
			//	bh.Ecosystemname = es.Name
			//	bh.Token_title = es.Token_title
			//}

			if err := json.Unmarshal([]byte(tss[i].Params), &bh.Params); err == nil {
				//bh.Params
			}
			ret = append(ret, bh)
		}

	}

	return &ret, len(GLogTranHash), err

}

func Deal_LogTransactionBlockTxDetial(objArr *[]LogTransaction) (*[]BlockTxDetailedInfo, error) {
	var (
		ret    []BlockTxDetailedInfo
		Blocks []int64
	)

	ret1 := Deal_Redupliction_LogTransaction(objArr)
	count := len(*ret1)
	if len(*ret1) == 0 {
		logrus.Info("logtran redup:")
		return &ret, nil
	}
	dat := *ret1

	TBlock := make(map[string]int64)
	Thash := make(map[string]bool)
	for i := 0; i < count; i++ {
		hash := hex.EncodeToString(dat[i].Hash)
		Thash[hash] = true
		key := strconv.FormatInt(dat[i].Block, 10)
		TBlock[key] = dat[i].Block
	}

	for _, k := range TBlock {
		Blocks = append(Blocks, k)
	}

	for i := int64(len(Blocks)); i > 0; i-- {
		bk := &Block{}
		found, err := bk.GetId(Blocks[i-1])
		if err == nil && found {
			rt, err := GetBlocksDetailedInfoHex(bk)
			if err == nil {
				for j := 0; j < int(rt.Tx); j++ {
					bh := BlockTxDetailedInfo{}
					bh.BlockID = rt.Header.BlockId
					bh.ContractName = rt.Transactions[j].ContractName
					bh.Hash = rt.Transactions[j].Hash
					bh.KeyID = rt.Transactions[j].KeyID
					//bh.Params = rt.Transactions[j].Params
					bh.Time = rt.Transactions[j].Time
					bh.Type = rt.Transactions[j].Type

					bh.Ecosystemname = rt.Transactions[j].Ecosystemname
					bh.Ecosystem = rt.Transactions[j].Ecosystem
					if bh.Ecosystem == 1 {
						bh.TokenSymbol = SysTokenSymbol
					} else {
						bh.TokenSymbol = rt.Transactions[j].TokenSymbol
					}
					bh.Size = rt.Transactions[j].Size

					lg1, err := json.Marshal(rt.Transactions[j].Params)
					if err == nil {
						bh.Params = string(lg1)
					}

					if Thash[rt.Transactions[j].Hash] {
						ret = append(ret, bh)
					}
				}
			} else {
				//logrus.Info("logtran GetBlocksDetailedInfoHex: %s", err.Error())
			}

		} else if err != nil {
			//logrus.Info("logtran GetBlocks  DetailedInfoHex: %s", err.Error())
		}
	}

	return &ret, nil
}

func Deal_Redupliction_LogTransaction(objArr *[]LogTransaction) *[]LogTransaction {
	var (
		ret []LogTransaction
	)
	if GLogTranHash == nil {
		GLogTranHash = make(map[string]int64)
	}
	for _, val := range *objArr {
		key := hex.EncodeToString(val.Hash)
		dat, ok := GLogTranHash[key]
		if ok {
			logrus.Info("GLogTranHash exist block:%d block1:%d key: "+key, dat, val.Block)
		} else {
			GLogTranHash[key] = val.Block
			ret = append(ret, val)
		}
	}
	return &ret
}

func (s *BlockTxDetailedInfoHex) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(s); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (s *BlockTxDetailedInfoHex) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &s); err != nil {
		return err
	}
	return nil
}

func (bt *BlockTxDetailedInfoHex) GetDbTxDetailedHash(hash string) (bool, error) {
	var bk Block
	hs, _ := hex.DecodeString(hash)

	dt := &LogTransaction{}
	fb, err := dt.getTransactionIdFromHash(hs)
	if !fb || err != nil {
		return false, err
	}
	fb, err = bk.GetId(dt.Block)
	if err != nil {
		return false, err
	}
	if fb {
		_, rt, bkts, err := DealTransactionBlockTxDetial(&bk)
		if err != nil {
			return false, err
		}

		for _, obj := range *bkts {
			if obj.Hash == hash {
				for _, rts := range rt.Transactions {
					if rts.Hash == hash {
						lg1, err := json.Marshal(rts.Params)
						if err == nil {
							obj.Params = string(lg1)
						}
					}
				}

				val, err := obj.Marshal()
				if err != nil {
					return false, err
				}
				err1 := bt.Unmarshal(val)
				if err1 != nil {
					return false, err1
				}
			}
		}
		if bt == nil {
			return false, nil
		}
	} else {
		return false, nil
	}
	return true, nil
}

func (bt *BlockTxDetailedInfoHex) GetDbTxDetailedId(id int64) ([]BlockTxDetailedInfoHex, error) {
	var bk Block
	var ret []BlockTxDetailedInfoHex
	fb, err := bk.GetId(id)
	if err != nil {
		return nil, err
	}
	if fb {
		_, _, bkts, err := DealTransactionBlockTxDetial(&bk)
		if err != nil {
			return nil, err
		}
		ret = make([]BlockTxDetailedInfoHex, len(*bkts))
		for i, obj := range *bkts {
			ret[i] = obj
		}
	} else {
		return nil, nil
	}

	return ret, nil
}

func (bt *BlockTxDetailedInfoHex) GetDb_txdetailedKey(key string, order string, limit, page int) ([]BlockTxDetailedInfoHex, int64, error) {
	var bk Block
	var ret []BlockTxDetailedInfoHex
	var needBlock []Block
	var total int64
	fb, err := bk.GetBlocksKey(converter.StringToAddress(key), order)
	if err != nil {
		return nil, total, err
	}
	ioffet := (page - 1) * limit
	var txData int32
	isAddOneBlock := false
	findFirstBlock := false
	var startTxId int
	var endTxiD int
	for i := 0; i < len(fb); i++ {
		txData += fb[i].Tx
		if int(txData) > ioffet {
			if int(txData) < ioffet+limit {
				needBlock = append(needBlock, fb[i])
				if !findFirstBlock {
					findFirstBlock = true
					startTxId = int(txData-(txData-fb[i].Tx)) - (int(txData) - ioffet)
				}
			} else {
				if !isAddOneBlock {
					isAddOneBlock = true
					needBlock = append(needBlock, fb[i])
					endTxiD = (int(txData) - (ioffet + limit))
					if !findFirstBlock {
						findFirstBlock = true
						startTxId = int(txData-(txData-fb[i].Tx)) - (int(txData) - ioffet)
					}
				}
			}
		}
	}
	total = int64(txData)

	for i := 0; i < len(needBlock); i++ {
		_, _, bkts, err := DealTransactionBlockTxDetial(&needBlock[i])
		if err != nil {
			return nil, total, err
		}
		for _, obj := range *bkts {
			ret = append(ret, obj)
		}
	}
	if len(ret) >= endTxiD && len(ret) >= startTxId && (len(ret)-endTxiD) >= startTxId {
		ret = ret[startTxId : len(ret)-endTxiD]
	}
	return ret, total, nil
}
