/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"time"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/converter"

	"github.com/shopspring/decimal"
)

// History represent record of history table
type History struct {
	ID               int64           `gorm:"primary_key;not null"`
	Senderid         int64           `gorm:"column:sender_id;not null"`
	Recipientid      int64           `gorm:"column:recipient_id;not null"`
	SenderBalance    decimal.Decimal `gorm:"column:sender_balance;not null"`
	RecipientBalance decimal.Decimal `gorm:"column:recipient_balance;not null"`
	Amount           decimal.Decimal `gorm:"column:amount;not null"`
	ValueDetail      string          `gorm:"column:value_detail;not null"`
	Comment          string          `gorm:"column:comment;not null"`
	Blockid          int64           `gorm:"column:block_id;not null"`
	Txhash           []byte          `gorm:"column:txhash;not null"`
	Createdat        int64           `gorm:"column:created_at;not null"`
	Ecosystem        int64           `gorm:"not null"`
	Type             int64           `gorm:"not null"`
	Status           int32           `gorm:"not null"`
	//Createdat   time.Time       `gorm:"column:created_at;not null"`
}

type fuelDetail_old struct {
	FuelRate  string `json:"fuel_rate"`
	TaxesSize int    `json:"taxes_size"`
	VmCostFee struct {
		Value      string `json:"value"`
		Decimal    string `json:"decimal"`
		Percentage int    `json:"percentage"`
	} `json:"vmCost_fee"`
	ElementFee struct {
		Value      string `json:"value"`
		Decimal    string `json:"decimal"`
		Percentage int    `json:"percentage"`
	} `json:"element_fee"`
	StorageFee struct {
		Value      string `json:"value"`
		Decimal    string `json:"decimal"`
		Percentage int    `json:"percentage"`
	} `json:"storage_fee"`
	ExpediteFee struct {
		Value      string `json:"value"`
		Decimal    string `json:"decimal"`
		Percentage int    `json:"percentage"`
	} `json:"expedite_fee"`
	PaymentType string `json:"payment_type"`
}

type FeeDetail struct {
	Flag           int     `json:"flag"`
	Value          string  `json:"value"`
	Convert        bool    `json:"convert"`
	Decimal        string  `json:"decimal"`
	ConversionRate float64 `json:"conversion_rate"`
}

type fuelDetail struct {
	FuelRate   string `json:"fuel_rate"`
	Combustion struct {
		Flag    int     `json:"flag"`
		Percent float64 `json:"percent"`
	} `json:"combustion"`
	TaxesSize   int       `json:"taxes_size"`
	VmCostFee   FeeDetail `json:"vmCost_fee"`
	ElementFee  FeeDetail `json:"element_fee"`
	StorageFee  FeeDetail `json:"storage_fee"`
	ExpediteFee FeeDetail `json:"expedite_fee"`
	PaymentType string    `json:"payment_type"`
	TokenSymbol string    `json:"token_symbol"`
}

type CombustionDetail struct {
	FuelRate   string `json:"fuel_rate"`
	Combustion struct {
		Flag    int    `json:"flag"`
		After   string `json:"after"`
		Value   string `json:"value"`
		Before  string `json:"before"`
		Percent int    `json:"percent"`
	} `json:"combustion"`
	VmCostFee   string `json:"vmCost_fee"`
	ElementFee  string `json:"element_fee"`
	StorageFee  string `json:"storage_fee"`
	ExpediteFee string `json:"expedite_fee"`
	TokenSymbol string `json:"token_symbol"`
}

type HistoryHex struct {
	ID               int64           `json:"id,omitempty"`
	Senderid         string          `json:"sender_id"`
	Recipientid      string          `json:"recipient_id"`
	SenderBalance    decimal.Decimal `json:"sender_balance"`
	RecipientBalance decimal.Decimal `json:"recipient_balance"`
	Amount           decimal.Decimal `json:"amount"`
	Comment          string          `json:"comment"`
	BlockId          int64           `json:"block_id"`
	Txhash           string          `json:"txhash"`
	//Createdat        time.Time       `json:"created_at"`
	Createdat     int64  `json:"created_at"`
	Ecosystem     int64  `json:"ecosystem"`
	Type          int64  `json:"type"`
	Ecosystemname string `json:"ecosystemname"`
	TokenSymbol   string `json:"token_symbol"`
	ContractName  string `json:"contract_name"`
}

type HistoryMergeHex struct {
	Ecosystem     int64           `json:"ecosystem"`
	ID            int64           `json:"id"`
	Senderid      string          `json:"sender_id"`
	Recipientid1  string          `json:"recipientid1"`
	Recipientid2  string          `json:"recipientid2,omitempty"`
	Recipientid3  string          `json:"recipientid3,omitempty"`
	Recipientid4  string          `json:"recipientid4,omitempty"`
	Amount1       decimal.Decimal `json:"amount1,omitempty"`
	Amount2       decimal.Decimal `json:"amount2,omitempty"`
	Amount3       decimal.Decimal `json:"amount3,omitempty"`
	Amount4       decimal.Decimal `json:"amount4,omitempty"`
	Comment       string          `json:"comment"`
	Blockid       int64           `json:"blockid"`
	Txhash        string          `json:"txhash"`
	Createdat     time.Time       `json:"created_at"`
	Ecosystemname string          `json:"ecosystemname"`
	TokenSymbol   string          `json:"token_symbol"`
	//Ecosystem     int64    `json:"ecosystem"`
}

type HistoryItem struct {
	Senderid    string          `json:"sender_id"`
	Recipientid string          `json:"recipient_id"`
	Amount      decimal.Decimal `json:"amount"`
	Events      int64           `json:"events,omitempty"`
	Scale       float64         `json:"scale,omitempty"`
	Flag        int             `json:"flag,omitempty"`
	TokenSymbol string          `json:"token_symbol,omitempty"`
	Combustion  string          `json:"combustion,omitempty"`
	FuelRate    int64           `json:"fuel_rate,omitempty"`
}

type transDetail struct {
	TokenSymbol string      `json:"token_symbol"`
	VmCostFee   HistoryItem `json:"vmCost_fee"`
	ElementFee  HistoryItem `json:"element_fee"`
	StorageFee  HistoryItem `json:"storage_fee"`
	ExpediteFee HistoryItem `json:"expedite_fee"`
}

type ecoExplorer struct {
	GasFee struct {
		Amount      decimal.Decimal `json:"amount"`
		TokenSymbol string          `json:"token_symbol,omitempty"`
	} `json:"gas_fee"`
	Fees         HistoryItem       `json:"fees"`
	Taxes        HistoryItem       `json:"taxes"`
	Combustion   HistoryItem       `json:"combustion"`
	Detail       transDetail       `json:"detail"`
	Exchange     transDetail       `json:"exchange"`
	EcosystemPay *ecosystemPayInfo `json:"ecosystem_pay,omitempty"`
}

type ecosystemPayInfo struct {
	Paid    HistoryItem `json:"paid"`
	Payment HistoryItem `json:"payment"`
}

type HistoryExplorer struct {
	//Ecosystem int64       `json:"ecosystem"`
	//ID        int64       `json:"id"`
	//EcoSystemName string           `json:"eco_system_name"`
	//TokenSymbol   string           `json:"token_symbol"`

	//Senderid string      `json:"sender_id"`
	//Address string      `json:"address"`
	Comment string      `json:"comment"`
	Fees    HistoryItem `json:"fees"`
	Taxes   HistoryItem `json:"taxes"`
	GasFee  struct {
		Amount      decimal.Decimal `json:"amount"`
		TokenSymbol string          `json:"token_symbol,omitempty"`
	} `json:"gas_fee"`
	TxFee HistoryItem `json:"tx_fee"`
	//CreateSetup   int64        `json:"created_setup"`
	Detail    transDetail  `json:"detail"`
	EcoDetail *ecoExplorer `json:"eco_detail,omitempty"`
	Status    int32        `json:"status"`
}

type WalletHistoryHex struct {
	Transaction int64           `json:"transaction"`
	InTx        int64           `json:"in_tx"`
	OutTx       int64           `json:"out_tx"`
	Inamount    decimal.Decimal `json:"inamount"`
	Outamount   decimal.Decimal `json:"outamount"`
	Amount      decimal.Decimal `json:"amount,omitempty"`
}

type HistoryTransaction struct {
	ID            int64  `json:"id"`
	Keyid         string `json:"key_id"`
	Blockid       int64  `json:"block_id"`
	Txhash        string `json:"txhash"`
	Createdat     int64  `json:"created_at"`
	Ecosystem     int64  `json:"ecosystem"`
	Ecosystemname string `json:"ecosystemname"`
	ContractName  string `json:"contract_name"`
	//Ecosystem     int64    `json:"ecosystem"`
}

type Historys []History

type MineHistoryRequest struct {
	Order       string `json:"order"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	EcosystemID int64  `json:"ecosystem"`
	Opt         string `json:"opt"`
	KeyId       string `json:"keyid"`
}

// TableName returns name of table
func (th *History) TableName() string {
	return "1_history"
}

func (th *History) Get(txHash []byte) (*HistoryMergeHex, error) {
	var (
		ts  []History
		tss HistoryMergeHex
	)

	err := conf.GetDbConn().Conn().Where("txhash = ?", txHash).Order("id ASC").Find(&ts).Error
	count := len(ts)
	if err == nil && count > 0 {
		if ts[0].Blockid > 0 {
			sort.Sort(Historys(ts))

			//fmt.Println(ts)
			tss.Ecosystem = ts[0].Ecosystem
			tss.TokenSymbol, tss.Ecosystemname = Tokens.Get(tss.Ecosystem), EcoNames.Get(tss.Ecosystem)

			tss.ID = ts[0].ID
			tss.Senderid = converter.AddressToString(ts[0].Senderid) //strconv.FormatInt(ts[0].Senderid, 10)
			tss.Comment = ts[0].Comment
			tss.Blockid = ts[0].Blockid
			tss.Txhash = hex.EncodeToString(ts[0].Txhash)
			fmt.Println(tss.Txhash)
			fmt.Println(string(ts[0].Txhash))
			tss.Createdat = time.Unix(ts[0].Createdat, 0)
			if count == 4 {
				tss.Recipientid1 = converter.AddressToString(ts[3].Recipientid) //strconv.FormatInt(ts[3].Recipientid, 10)
				tss.Recipientid2 = converter.AddressToString(ts[2].Recipientid) //strconv.FormatInt(ts[2].Recipientid, 10)
				tss.Recipientid3 = converter.AddressToString(ts[1].Recipientid) //strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid4 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[3].Amount
				tss.Amount2 = ts[2].Amount
				tss.Amount3 = ts[1].Amount
				tss.Amount4 = ts[0].Amount
				//				fmt.Println(ts[2].Amount)
				//				fmt.Println(ts[1].Amount)
				//				fmt.Println(ts[0].Amount)
				//				fmt.Println(tss)
			} else if count == 3 {
				tss.Recipientid1 = converter.AddressToString(ts[2].Recipientid) //strconv.FormatInt(ts[2].Recipientid, 10)
				tss.Recipientid2 = converter.AddressToString(ts[1].Recipientid) //strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid3 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[2].Amount
				tss.Amount2 = ts[1].Amount
				tss.Amount3 = ts[0].Amount
				//				fmt.Println(ts[2].Amount)
				//				fmt.Println(ts[1].Amount)
				//				fmt.Println(ts[0].Amount)
				//				fmt.Println(tss)
			} else if count == 2 {
				tss.Recipientid2 = converter.AddressToString(ts[1].Recipientid) //strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid3 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount2 = ts[1].Amount
				tss.Amount3 = ts[0].Amount
				//tss.Amount1 =decimal.NewFromFloat(0)
			} else if count == 1 {
				tss.Recipientid1 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[0].Amount
			}
		}
	}
	return &tss, err
}

func (th *History) GetHashSum(txHash []byte, ecosystem int64) string {
	var sum SumAmount
	err := GetDB(nil).Table(th.TableName()).Select("sum(amount)").Where("txhash = ? AND ecosystem = ?", txHash, ecosystem).Take(&sum).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err, "hash": hex.EncodeToString(txHash), "ecosystem": ecosystem}).Error("Get Hash Sum Failed")
		return "0"
	}
	return sum.Sum.String()
}

func (th *History) GetByHash(txHash []byte) (bool, error) {
	return isFound(GetDB(nil).Where("txhash = ?", txHash).First(th))
}

func (th *History) GetByHashExist(txHash []byte) (bool, error) {
	return isFound(GetDB(nil).Where("txhash = ?", txHash))
}

//GetExplorer Not all transactions will exist in the history table
func (th *History) GetExplorer(txHash []byte) (*HistoryExplorer, error) {
	var (
		ts              []History
		tss             HistoryExplorer
		ecoInfo         ecoExplorer
		isEcosystemPaid bool
		ecoPaid         ecosystemPayInfo
	)
	type txInfo struct {
		Ecosystem   int64
		TokenSymbol string
	}
	var ecoIdList []txInfo
	err := GetDB(nil).Raw(`
SELECT h1.ecosystem,CASE WHEN h1.ecosystem = 1 THEN
		coalesce(es.token_symbol,'IBXC')
	ELSE
		es.token_symbol
	END token_symbol
FROM (
	SELECT ecosystem FROM "1_history" WHERE txhash = ? GROUP BY ecosystem
)AS h1
LEFT JOIN(
	SELECT id,name,token_symbol FROM "1_ecosystems"
)AS es on(es.id = h1.ecosystem)
`, txHash).Find(&ecoIdList).Error
	if err != nil {
		return nil, err
	}

	//f, err := isFound(GetDB(nil).Select("address").Where("hash = ?", txHash).First(&tss.Address))
	//if err != nil {
	//	return nil, err
	//}
	//if !f {
	//	return nil, errors.New("hash doesn't not exist")
	//}

	ecoTokenSymbol := make(map[int64]string)
	for _, value := range ecoIdList {
		ecoTokenSymbol[value.Ecosystem] = value.TokenSymbol
	}
	getFeeRate := func(fe FeeDetail) (float64, int) {
		if fe.Flag == 0 {
			return 0, 0
		}
		if fe.Flag > 1 {
			return fe.ConversionRate, fe.Flag
		}
		return 100, 1
	}
	getFeeAmount := func(fe FeeDetail) decimal.Decimal {
		ret, _ := decimal.NewFromString(fe.Value)
		return ret
	}
	getEcoExchangeDetail := func(fuel fuelDetail, senderid string, isExchange bool) transDetail {
		var det transDetail
		if fuel.VmCostFee.Flag == 2 || fuel.VmCostFee.Flag == 1 {
			det.VmCostFee.Amount = getFeeAmount(fuel.VmCostFee)
			det.VmCostFee.Senderid = senderid
			if isExchange {
				det.VmCostFee.Scale = fuel.VmCostFee.ConversionRate
			}
		}
		if fuel.ElementFee.Flag == 2 || fuel.ElementFee.Flag == 1 {
			det.ElementFee.Amount = getFeeAmount(fuel.ElementFee)
			det.ElementFee.Senderid = senderid
			if isExchange {
				det.ElementFee.Scale = fuel.ElementFee.ConversionRate
			}
		}
		if fuel.StorageFee.Flag == 2 || fuel.StorageFee.Flag == 1 {
			det.StorageFee.Amount = getFeeAmount(fuel.StorageFee)
			det.StorageFee.Senderid = senderid
			if isExchange {
				det.StorageFee.Scale = fuel.StorageFee.ConversionRate
			}
		}
		if fuel.ExpediteFee.Flag == 2 || fuel.ExpediteFee.Flag == 1 {
			det.ExpediteFee.Amount = getFeeAmount(fuel.ExpediteFee)
			det.ExpediteFee.Senderid = senderid
			if isExchange {
				det.ExpediteFee.Scale = fuel.ExpediteFee.ConversionRate
			}
		}
		if fuel.FuelRate != "" {
			combustion, _ := decimal.NewFromString(fuel.FuelRate)
			rate := decimal.NewFromInt(1e4)
			if !combustion.IsZero() {
				det.VmCostFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
				det.ElementFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
				det.StorageFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
				det.ExpediteFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
			}
		}
		return det
	}

	err = GetDB(nil).Where("txhash = ?", txHash).Order("id ASC").Find(&ts).Error
	if err != nil {
		return nil, err
	}
	count := len(ts)
	tss.Status = 0
	if count > 0 {
		//tss.CreateSetup = MsToSeconds(ts[0].Createdat)
		tss.Comment = ts[0].Comment

		//var ts1 TransactionStatus
		//found, err := ts1.DbconngetSqlite(txHash)
		//if err == nil && found {
		//	tss.CreateSetup = MsToSeconds(ts1.Time)
		//}

		//for
		for _, val := range ts {
			var item HistoryItem
			tss.Status = val.Status
			if val.Type == 1 {
				item.TokenSymbol = ecoTokenSymbol[val.Ecosystem]
				item.Amount = val.Amount
				item.Senderid = converter.AddressToString(val.Senderid)
				item.Recipientid = converter.AddressToString(val.Recipientid)
				if tss.Status == 1 || tss.Status == 2 {
					tss.TxFee.Senderid = item.Senderid
				}
				if val.Ecosystem == 1 {
					tss.Fees.Amount = tss.Fees.Amount.Add(item.Amount)
					tss.Fees.Recipientid = item.Recipientid
					tss.Fees.TokenSymbol = item.TokenSymbol
				} else {
					ecoInfo.Fees.Amount = item.Amount
					ecoInfo.Fees.Recipientid = item.Recipientid
					ecoInfo.Fees.TokenSymbol = item.TokenSymbol
				}

				if val.ValueDetail != "" {
					var fuel fuelDetail
					var det transDetail
					if err := json.Unmarshal([]byte(val.ValueDetail), &fuel); err != nil {
						return nil, errors.New("get tx value detail failed:" + err.Error())
					}
					switch fuel.PaymentType {
					case "ContractCaller":

					case "EcosystemAddress", "ContractBinder":
						ecoPaid.Payment = item
						isEcosystemPaid = true
					}
					det.TokenSymbol = item.TokenSymbol
					det.VmCostFee.Senderid = item.Senderid
					det.ElementFee.Senderid = item.Senderid
					det.StorageFee.Senderid = item.Senderid
					det.ExpediteFee.Senderid = item.Senderid
					det.VmCostFee.Amount = getFeeAmount(fuel.VmCostFee)
					det.ElementFee.Amount = getFeeAmount(fuel.ElementFee)
					det.StorageFee.Amount = getFeeAmount(fuel.StorageFee)
					det.ExpediteFee.Amount = getFeeAmount(fuel.ExpediteFee)
					det.VmCostFee.Scale, det.VmCostFee.Flag = getFeeRate(fuel.VmCostFee)
					det.ElementFee.Scale, det.ElementFee.Flag = getFeeRate(fuel.ElementFee)
					det.StorageFee.Scale, det.StorageFee.Flag = getFeeRate(fuel.StorageFee)
					det.ExpediteFee.Scale, det.ExpediteFee.Flag = getFeeRate(fuel.ExpediteFee)
					if val.Ecosystem != 1 {
						dtl := getEcoExchangeDetail(fuel, item.Senderid, false)
						ecoInfo.Detail.ElementFee.Amount = dtl.ElementFee.Amount
						ecoInfo.Detail.ElementFee.Senderid = dtl.ElementFee.Senderid
						ecoInfo.Detail.ElementFee.FuelRate = dtl.ElementFee.FuelRate
						ecoInfo.Detail.VmCostFee.Amount = dtl.VmCostFee.Amount
						ecoInfo.Detail.VmCostFee.Senderid = dtl.VmCostFee.Senderid
						ecoInfo.Detail.VmCostFee.FuelRate = dtl.VmCostFee.FuelRate
						ecoInfo.Detail.StorageFee.Amount = dtl.StorageFee.Amount
						ecoInfo.Detail.StorageFee.Senderid = dtl.StorageFee.Senderid
						ecoInfo.Detail.StorageFee.FuelRate = dtl.StorageFee.FuelRate
						ecoInfo.Detail.ExpediteFee.Amount = dtl.ExpediteFee.Amount
						ecoInfo.Detail.ExpediteFee.Senderid = dtl.ExpediteFee.Senderid
						ecoInfo.Detail.ExpediteFee.FuelRate = dtl.ExpediteFee.FuelRate
						ecoInfo.Detail.TokenSymbol = item.TokenSymbol
					} else {
						tss.Detail.TokenSymbol = item.TokenSymbol
						tss.Detail.VmCostFee.Amount = tss.Detail.VmCostFee.Amount.Add(det.VmCostFee.Amount)
						tss.Detail.ElementFee.Amount = tss.Detail.ElementFee.Amount.Add(det.ElementFee.Amount)
						tss.Detail.StorageFee.Amount = tss.Detail.StorageFee.Amount.Add(det.StorageFee.Amount)
						tss.Detail.ExpediteFee.Amount = tss.Detail.ExpediteFee.Amount.Add(det.ExpediteFee.Amount)
						if det.VmCostFee.Flag != 0 {
							tss.Detail.VmCostFee.Flag = det.VmCostFee.Flag
							if (tss.Detail.VmCostFee.Flag == 1 && !isEcosystemPaid) || (tss.Detail.VmCostFee.Flag == 2 && isEcosystemPaid) {
								tss.Detail.VmCostFee.Senderid = item.Senderid
							}
						}
						if det.ElementFee.Flag != 0 {
							tss.Detail.ElementFee.Flag = det.ElementFee.Flag
							if (tss.Detail.ElementFee.Flag == 1 && !isEcosystemPaid) || (tss.Detail.ElementFee.Flag == 2 && isEcosystemPaid) {
								tss.Detail.ElementFee.Senderid = item.Senderid
							}
						}
						if det.StorageFee.Flag != 0 {
							tss.Detail.StorageFee.Flag = det.StorageFee.Flag
							if (tss.Detail.StorageFee.Flag == 1 && !isEcosystemPaid) || (tss.Detail.StorageFee.Flag == 2 && isEcosystemPaid) {
								tss.Detail.StorageFee.Senderid = item.Senderid
							}
						}
						if det.ExpediteFee.Flag != 0 {
							tss.Detail.ExpediteFee.Flag = det.ExpediteFee.Flag
							if (tss.Detail.ExpediteFee.Flag == 1 && !isEcosystemPaid) || (tss.Detail.ExpediteFee.Flag == 2 && isEcosystemPaid) {
								tss.Detail.ExpediteFee.Senderid = item.Senderid
							}
						}
					}
				}
			} else if val.Type == 2 {
				//item.Senderid = converter.AddressToString(val.Senderid)
				item.Recipientid = converter.AddressToString(val.Recipientid)
				item.Amount = val.Amount
				if val.Ecosystem == 1 {
					tss.Taxes.TokenSymbol = tss.Fees.TokenSymbol
					tss.Taxes.Amount = tss.Taxes.Amount.Add(item.Amount)
					tss.Taxes.Recipientid = item.Recipientid
					if isEcosystemPaid == true {
						ecoPaid.Payment.Amount = ecoPaid.Payment.Amount.Add(item.Amount)
					}
				} else {
					item.TokenSymbol = ecoInfo.Fees.TokenSymbol
					ecoInfo.Taxes = item
				}
			} else if val.Type == 15 {
				item.Senderid = converter.AddressToString(val.Senderid)
				item.Recipientid = converter.AddressToString(val.Recipientid)
				item.Amount = val.Amount
				item.TokenSymbol = ecoTokenSymbol[val.Ecosystem]
				ecoPaid.Paid = item

				if val.ValueDetail != "" {
					var fuel fuelDetail
					var det transDetail
					if err := json.Unmarshal([]byte(val.ValueDetail), &fuel); err != nil {
						return nil, errors.New("get ecosystem paid value detail failed:" + err.Error())
					}
					det = getEcoExchangeDetail(fuel, item.Senderid, true)
					det.TokenSymbol = item.TokenSymbol
					ecoInfo.Exchange = det
				}
			} else if val.Type == 16 {
				item.Recipientid = converter.AddressToString(val.Recipientid)
				item.Amount = val.Amount
				ecoInfo.Combustion = item

				if val.ValueDetail != "" {
					var detail CombustionDetail
					if err := json.Unmarshal([]byte(val.ValueDetail), &detail); err != nil {
						return nil, errors.New("get ecosystem combustion value detail failed:" + err.Error())
					}
					if detail.FuelRate != "" {
						combustion, _ := decimal.NewFromString(detail.FuelRate)
						rate := decimal.NewFromInt(1e4)
						if !combustion.IsZero() {
							ecoInfo.Detail.VmCostFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
							ecoInfo.Detail.ElementFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
							ecoInfo.Detail.StorageFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
							ecoInfo.Detail.ExpediteFee.FuelRate = combustion.DivRound(rate, 0).IntPart()
						}
					}

					ecoInfo.Detail.VmCostFee.Combustion = detail.VmCostFee
					ecoInfo.Detail.ElementFee.Combustion = detail.ElementFee
					ecoInfo.Detail.StorageFee.Combustion = detail.StorageFee
					ecoInfo.Detail.ExpediteFee.Combustion = detail.ExpediteFee
					ecoInfo.Combustion.Scale = float64(detail.Combustion.Percent)
				}
			} else {
				item.Senderid = converter.AddressToString(val.Senderid)
				item.Recipientid = converter.AddressToString(val.Recipientid)
				item.Amount = val.Amount
				item.TokenSymbol = ecoTokenSymbol[val.Ecosystem]
				item.Events = val.Type
				tss.TxFee = item
			}
		}
	}
	if ecoInfo.Exchange.TokenSymbol != "" {
		ecoInfo.GasFee.Amount = ecoInfo.Taxes.Amount.Add(ecoInfo.Fees.Amount)
		ecoInfo.GasFee.TokenSymbol = ecoInfo.Taxes.TokenSymbol
		tss.EcoDetail = &ecoInfo
	}
	if tss.EcoDetail != nil {
		if isEcosystemPaid {
			tss.EcoDetail.EcosystemPay = &ecoPaid
		}
	}
	tss.GasFee.Amount = tss.Taxes.Amount.Add(tss.Fees.Amount)
	tss.GasFee.TokenSymbol = tss.Fees.TokenSymbol

	return &tss, err
}

func (th *History) GetGasFeeByTxHashList(txHash [][]byte) (*BlockListResponse, error) {
	var (
		tss    BlockListResponse
		gasfee SumAmount
	)

	err := GetDB(nil).Table(th.TableName()).Select("sum(amount)").Where("txhash in(?) AND type in(1,2)", txHash).Take(&gasfee).Error
	tss.GasFee = gasfee.Sum.String()

	return &tss, err
}

func (th *History) GetBlockRewardById(blockId int64) (*BlockListResponse, error) {
	var (
		ts  []History
		tss BlockListResponse
	)
	err := GetDB(nil).Select("type,recipient_id,amount").Where("block_id = ? AND type = 12", blockId).Order("id ASC").Find(&ts).Error
	count := len(ts)
	if err == nil && count > 0 {
		tss.ID = ts[0].Blockid
		//for
		for _, ret := range ts {
			if ret.Type == 12 {
				tss.Recipientid = converter.AddressToString(ret.Recipientid)
				tss.Reward = ret.Amount.String()
			}
		}
	}

	return &tss, err
}

func (th *History) GetTxListExplorer(txHash []byte) (*BlockTxDetailedInfoHex, error) {
	var (
		ts  []History
		tss BlockTxDetailedInfoHex
	)
	//get ecosystem 1 gasFee
	err := GetDB(nil).Select("type,amount,status,ecosystem").Where("txhash = ?", txHash).Order("id ASC").Find(&ts).Error
	count := len(ts)

	//TODO:status need add struct log table
	tss.Status = 0
	if err == nil && count > 0 {
		//for
		for _, ret := range ts {
			tss.Status = ret.Status
			if ret.Type == 1 || ret.Type == 2 {
				if ret.Ecosystem == 1 {
					tss.GasFee = tss.GasFee.Add(ret.Amount)
				}
			} else {
				if ret.Type != 15 && ret.Type != 16 {
					tss.Amount = tss.Amount.Add(ret.Amount)
				}
			}
		}
	}

	return &tss, err
}

func (th *History) GetHistoryTimeList(time time.Time) (*[]History, error) {
	var (
		tss []History
	)

	err := conf.GetDbConn().Conn().Model(&History{}).Where("created_at >?", time.UnixMilli()).Order("created_at desc").Find(&tss).Error
	return &tss, err
}

func (th *History) GetHistoryIdList(id int64) (*[]History, error) {
	var (
		tss []History
	)

	err := conf.GetDbConn().Conn().Model(&History{}).Where("id >?", id).Order("id desc").Find(&tss).Error
	return &tss, err
}

//GetHistory Get is retrieving model from database
func (th *History) GetHistory(page int, size int, order string) (*[]HistoryHex, int64, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
	)

	err := conf.GetDbConn().Conn().Limit(size).Offset((page - 1) * size).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").Count(&num).Error
	if err != nil {
		return &ret, num, err
	}
	for i := 0; i < len(tss); i++ {
		//fmt.Println("offset Error:%d ", offset)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		da.TokenSymbol, da.Ecosystemname = Tokens.Get(da.Ecosystem), EcoNames.Get(da.Ecosystem)
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.BlockId = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = MsToSeconds(tss[i].Createdat)
		ret = append(ret, da)
	}

	return &ret, num, err
}

//GetWallets Get is retrieving model from database
func (th *History) GetWallets(page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//offset int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)

	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}
	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? OR sender_id = ?", keyId, keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}

	//total = deal_history_total(&tss)

	count := int64(len(tss))
	//fmt.Println("tr_blocks Error: %d", num)
	//ioffet = (page - 1) * size
	//if num < page*size {
	//	size = num % size
	//}
	//if num < ioffet || num < 1 {
	//	return &ret, num, total, err
	//}
	for i = 0; i < count; i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		da.TokenSymbol, da.Ecosystemname = Tokens.Get(da.Ecosystem), EcoNames.Get(da.Ecosystem)
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.BlockId = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = MsToSeconds(tss[i].Createdat)
		ret = append(ret, da)
		//ioffet++
		//fmt.Println("offet Error:%d ", ioffet)
	}

	return &ret, num, total, err
}

//GetEcosytemWallets Get is retrieving model from database
func (th *History) GetEcosytemWallets(id int64, page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)

	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}
	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ?", keyId, keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ?", keyId, keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ?", keyId, keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}

	//total = deal_history_total(&tss)

	count := int64(len(tss))
	//fmt.Println("tr_blocks Error: %d", num)
	//ioffet = (page - 1) * size
	//if num < page*size {
	//	size = num % size
	//}
	//if num < ioffet || num < 1 {
	//	return &ret, num, total, err
	//}
	for i = 0; i < count; i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		da.TokenSymbol, da.Ecosystemname = Tokens.Get(da.Ecosystem), EcoNames.Get(da.Ecosystem)

		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.BlockId = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = MsToSeconds(tss[i].Createdat)
		ret = append(ret, da)
		//i++
		//fmt.Println("offet Error:%d ", ioffet)
	}

	return &ret, num, total, err
}

func (th *History) GetWalletTotals(wallet string) (*WalletHistoryHex, error) {
	var (
		tss1  []History
		tss2  []History
		ret   WalletHistoryHex
		keyId int64
		err   error
	)

	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("recipient_id = ?", keyId).
		Order("created_at desc").Find(&tss1).Error
	if err != nil {
		return &ret, err
	}
	err = conf.GetDbConn().Conn().Table("1_history").
		Where("sender_id = ?", keyId).
		Order("created_at desc").Find(&tss2).Error
	if err != nil {
		return &ret, err
	}
	ret.Transaction = int64(len(tss1)) + int64(len(tss2))
	ret.Inamount = deal_history_total(&tss1)
	ret.Outamount = deal_history_total(&tss2)

	return &ret, err
}

func (th *History) GetAccountHistoryTotals(id int64, keyId int64) (*WalletHistoryHex, error) {
	var (
		ret    WalletHistoryHex
		sCount int64
		rCount int64
		in     = "0"
		out    = "0"
		err    error
		his    History
	)
	type inOutInfo struct {
		InTx      int64
		OutTx     int64
		InAmount  string
		OutAmount string
	}
	var utxoInfo inOutInfo
	//in out tx
	err = GetDB(nil).Raw(`
SELECT count(v1.txhash)in_tx FROM(
	SELECT txhash FROM "1_history" WHERE recipient_id = ? and ecosystem = ? GROUP BY "txhash"
)AS v1
`, keyId, id).Row().Scan(&rCount)
	if err != nil {
		return &ret, err
	}

	err = GetDB(nil).Raw(`
SELECT count(v1.txhash)out_tx FROM(
	SELECT txhash FROM "1_history" WHERE sender_id = ? and ecosystem = ? GROUP BY "txhash"
)AS v1
`, keyId, id).Row().Scan(&sCount)
	if err != nil {
		return &ret, err
	}

	//in amount
	if rCount > 0 {
		err = GetDB(nil).Table(his.TableName()).Select("sum(amount)").Where("recipient_id = ? and ecosystem = ?", keyId, id).Row().Scan(&in)
		if err != nil {
			return &ret, err
		}
	}

	utxoQuery := `
SELECT COALESCE(sum(v4.out_tx),0)out_tx,COALESCE(sum(v4.in_tx),0)in_tx,
COALESCE(sum(v4.in_amount),0)in_amount,COALESCE(sum(v4.out_amount),0)out_amount
FROM(
	SELECT v3.*,CASE WHEN v3.diff >=0 THEN 
			v3.diff 
		ELSE
			0
		END AS in_amount,
		CASE WHEN v3.diff < 0 THEN 
			v3.diff*-1 
		ELSE
			0
		END AS out_amount,
		CASE WHEN v3.diff >=0 THEN 
			1
		ELSE
			0
		END AS in_tx,
		CASE WHEN v3.diff < 0 THEN 
			1
		ELSE
			0
		END AS out_tx
	FROM(
		SELECT v1.*,COALESCE(v2.pre_sum,0)pre_sum,CASE WHEN v1.num > 1 THEN
				CASE WHEN v2.pre_sum < v1.output_value THEN
					(v1.output_value - v2.pre_sum)*-1
				ELSE
					v2.pre_sum - v1.output_value
				END
			ELSE
				0
		END diff
		 FROM(
			SELECT sum(output_value)output_value,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
				FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  WHERE output_key_id = ? AND ecosystem = ? GROUP BY output_tx_hash 
			ORDER BY timestamp asc
		)AS v1
		LEFT JOIN(
			SELECT sum(output_value)pre_sum,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
				FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  WHERE output_key_id = ? AND ecosystem = ? GROUP BY output_tx_hash 
			ORDER BY timestamp asc
		)AS v2 ON(v2.num+1=v1.num)
	)AS v3
)as v4`

	err = GetDB(nil).Raw(utxoQuery, keyId, id, keyId, id).Take(&utxoInfo).Error
	if err != nil {
		return nil, err
	}

	//out amount
	if sCount > 0 {
		err = GetDB(nil).Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ?", keyId, id).Row().Scan(&out)
		if err != nil {
			return &ret, err
		}
	}

	in1Amount, _ := decimal.NewFromString(in)
	in2Amount, _ := decimal.NewFromString(utxoInfo.InAmount)
	out1Amount, _ := decimal.NewFromString(out)
	out2Amount, _ := decimal.NewFromString(utxoInfo.OutAmount)

	ret.InTx = rCount + utxoInfo.InTx
	ret.OutTx = sCount + utxoInfo.OutTx

	ret.Transaction = ret.InTx + ret.OutTx
	ret.Inamount = in1Amount.Add(in2Amount)
	ret.Outamount = out1Amount.Add(out2Amount)

	return &ret, err
}

func (th *History) GetWalletTimeLineHistoryTotals(ecosystem int64, keyId int64, getDay int) (*AccountHistoryChart, error) {
	var (
		ret AccountHistoryChart
		err error
		t1  time.Time
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	if getDay != 0 {
		t1 = yesterday.AddDate(0, 0, -1*getDay)
	}

	type daysAmount struct {
		Days   string          `gorm:"column:days"`
		Amount decimal.Decimal `gorm:"column:amount"`
	}

	type utxoDaysAmount struct {
		Days      string
		InAmount  decimal.Decimal
		OutAmount decimal.Decimal
	}
	var utxoInfo []utxoDaysAmount

	getDaysAmount := func(dayTime int64, list []daysAmount) decimal.Decimal {
		for i := 0; i < len(list); i++ {
			times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
			if dayTime == times.Unix() {
				return list[i].Amount
			}
		}
		return decimal.Zero
	}

	getUtxoDaysAmount := func(dayTime int64, list []utxoDaysAmount, isGetOut bool) decimal.Decimal {
		for i := 0; i < len(list); i++ {
			times, _ := time.ParseInLocation("2006-01-02", list[i].Days, time.Local)
			if dayTime == times.Unix() {
				if isGetOut {
					return list[i].OutAmount
				} else {
					return list[i].InAmount
				}
			}
		}
		return decimal.Zero
	}

	var (
		inAmount  []daysAmount
		outAmount []daysAmount
	)
	err = GetDB(nil).Table(th.TableName()).Raw(`SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') days ,sum(amount) amount  
FROM "1_history" WHERE recipient_id <> 0 AND recipient_id = ?
and ecosystem = ? AND created_at >= ? GROUP BY days ORDER BY days asc`, keyId, ecosystem, t1.UnixMilli()).Find(&inAmount).Error
	if err != nil {
		return &ret, err
	}

	err = GetDB(nil).Table(th.TableName()).Raw(`SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') days ,sum(amount) amount  
FROM "1_history" WHERE sender_id <> 0 AND sender_id = ?
and ecosystem = ? AND created_at >= ? GROUP BY days ORDER BY days asc`, keyId, ecosystem, t1.UnixMilli()).Find(&outAmount).Error
	if err != nil {
		return &ret, err
	}

	err = GetDB(nil).Raw(`
SELECT to_char(to_timestamp(v4.timestamp/1000), 'yyyy-mm-dd') AS days,
COALESCE(sum(v4.in_amount),0)in_amount,COALESCE(sum(v4.out_amount),0)out_amount
FROM(
	SELECT v3.*,CASE WHEN v3.diff >= 0 THEN 
			v3.diff 
		ELSE
			0
		END AS in_amount,
		CASE WHEN v3.diff < 0 THEN 
			v3.diff*-1 
		ELSE
			0
		END AS out_amount
	FROM(
		SELECT v1.*,COALESCE(v2.pre_sum,0)pre_sum,CASE WHEN v1.num > 1 THEN
				CASE WHEN v2.pre_sum < v1.output_value THEN
					(v1.output_value - v2.pre_sum)*-1
				ELSE
					v2.pre_sum - v1.output_value
				END
			ELSE
				0
		END diff
		 FROM(
			SELECT sum(output_value)output_value,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
				FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  WHERE output_key_id = ? AND ecosystem = ? AND timestamp >= ? GROUP BY output_tx_hash 
			ORDER BY timestamp asc
		)AS v1
		LEFT JOIN(
			SELECT sum(output_value)pre_sum,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
				FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  WHERE output_key_id = ? AND ecosystem = ? AND timestamp >= ? GROUP BY output_tx_hash 
			ORDER BY timestamp asc
		)AS v2 ON(v2.num+1 = v1.num)
	)AS v3
)as v4
GROUP BY days
`, keyId, ecosystem, t1.UnixMilli(), keyId, ecosystem, t1.UnixMilli()).Find(&utxoInfo).Error
	if err != nil {
		return nil, err
	}

	var startTime time.Time
	var stTime string

	if t1.UnixMilli() > 0 {
		stTime = t1.Format("2006-01-02")
	} else {
		stTime = time.Unix(FirstBlockTime, 0).Format("2006-01-02")
	}
	t2, err := time.ParseInLocation("2006-01-02", stTime, time.Local)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "day": stTime}).Error("Get Account Time Line History Totals startTime parseInLocation Failed")
		return &ret, err
	}
	startTime = t2
	today := yesterday.AddDate(0, 0, 1)
	for startTime.Unix() <= today.Unix() {
		ret.Time = append(ret.Time, startTime.Unix())
		ret.Inamount = append(ret.Inamount, getDaysAmount(startTime.Unix(), inAmount).Add(getUtxoDaysAmount(startTime.Unix(), utxoInfo, false)).String())
		ret.Outamount = append(ret.Outamount, getDaysAmount(startTime.Unix(), outAmount).Add(getUtxoDaysAmount(startTime.Unix(), utxoInfo, true)).String())
		startTime = startTime.AddDate(0, 0, 1)
	}

	return &ret, err
}

func (u Historys) Len() int {
	return len(u)
}

func (u Historys) Less(i, j int) bool {
	dat := u[i].Amount.Cmp(u[j].Amount)
	return dat < 0 //sort by id if id is the same sort by name...
}

func (u Historys) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (th *History) Get_Sqlite(txHash []byte) (*HistoryMergeHex, error) {
	var (
		ts  []History
		tss HistoryMergeHex
		//i   int
	)
	err := conf.GetDbConn().Conn().Where("txhash = ?", txHash).Find(&ts).Error
	count := len(ts)
	if err == nil && count > 0 {
		if ts[0].Blockid > 0 {
			//fmt.Println(ts)
			sort.Sort(Historys(ts))
			//fmt.Println(ts)
			tss.Ecosystem = ts[0].Ecosystem
			tss.TokenSymbol, tss.Ecosystemname = Tokens.Get(tss.Ecosystem), EcoNames.Get(tss.Ecosystem)
			tss.ID = ts[0].ID
			tss.Senderid = strconv.FormatInt(ts[0].Senderid, 10)
			tss.Comment = ts[0].Comment
			tss.Blockid = ts[0].Blockid
			tss.Txhash = hex.EncodeToString(ts[0].Txhash)
			tss.Createdat = time.Unix(ts[0].Createdat, 0)
			if count == 3 {
				tss.Recipientid1 = strconv.FormatInt(ts[2].Recipientid, 10)
				tss.Recipientid2 = strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid3 = strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[2].Amount
				tss.Amount2 = ts[1].Amount
				tss.Amount3 = ts[0].Amount
				//				fmt.Println(ts[2].Amount)
				//				fmt.Println(ts[1].Amount)
				//				fmt.Println(ts[0].Amount)
				//				fmt.Println(tss)
			} else if count == 2 {
				tss.Recipientid1 = strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid2 = strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[1].Amount
				tss.Amount2 = ts[0].Amount
			} else if count == 1 {
				tss.Recipientid1 = strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[0].Amount
			}
		}
	}
	return &tss, err
}

//Get is retrieving model from database
func (th *History) GetHistorys_Sqlite(page int, size int, order string) (*[]HistoryHex, int64, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
	)

	err := conf.GetDbConn().Conn().Limit(size).Offset((page - 1) * size).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").Count(&num).Error
	if err != nil {
		return &ret, num, err
	}
	for i := 0; i < len(tss); i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem

		da.TokenSymbol, da.Ecosystemname = Tokens.Get(da.Ecosystem), EcoNames.Get(da.Ecosystem)
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.BlockId = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = MsToSeconds(tss[i].Createdat)
		ret = append(ret, da)
	}

	return &ret, num, err
}

//Get is retrieving model from database
func (th *History) GetWallets_Sqlite(page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)
	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}

	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? OR sender_id = ?", keyId, keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}
	//total = deal_history_total(&tss)

	count := int64(len(tss))
	//num = int64(len(tss))
	//ioffet = (page - 1) * size
	//if num < page*size {
	//	size = num % size
	//}
	//if num < ioffet || num < 1 {
	//	return &ret, num, total, err
	//}

	for i = 0; i < count; i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		da.TokenSymbol, da.Ecosystemname = Tokens.Get(da.Ecosystem), EcoNames.Get(da.Ecosystem)
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.BlockId = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = MsToSeconds(tss[i].Createdat)
		ret = append(ret, da)
		//ioffet++
		//fmt.Println("offet Error:%d ", ioffet)
	}

	return &ret, num, total, err
}

//Get is retrieving model from database
func (th *History) GetWallets_EcosytemSqlite(id int64, page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)
	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}

	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? ", keyId, keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? ", keyId, keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? ", keyId, keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}
	//total = deal_history_total(&tss)

	count := int64(len(tss))

	for i = 0; i < count; i++ {
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		da.TokenSymbol, da.Ecosystemname = Tokens.Get(da.Ecosystem), EcoNames.Get(da.Ecosystem)
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.BlockId = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = MsToSeconds(tss[i].Createdat)
		ret = append(ret, da)
	}

	return &ret, num, total, err
}

func (th *History) GetWalletTotals_Sqlites(wallet string) (*WalletHistoryHex, error) {
	var (
		tss1  []History
		tss2  []History
		ret   WalletHistoryHex
		keyId int64
		err   error
	)

	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("recipient_id = ?", keyId).
		Order("created_at desc").Find(&tss1).Error
	if err != nil {
		return &ret, err
	}
	err = conf.GetDbConn().Conn().Table("1_history").
		Where("sender_id = ?", keyId).
		Order("created_at desc").Find(&tss2).Error
	if err != nil {
		return &ret, err
	}
	ret.Transaction = int64(len(tss1)) + int64(len(tss2))
	ret.Inamount = deal_history_total(&tss1)
	ret.Outamount = deal_history_total(&tss2)

	return &ret, err
}
func deal_history_total(objArr *[]History) decimal.Decimal {
	var (
		total decimal.Decimal
	)
	for _, val := range *objArr {
		total = total.Add(val.Amount)
	}
	return total
}

func (th *History) GetTodayCirculationsAmount(nftBlockReward float64) float64 {
	tz := time.Unix(GetNowTimeUnix(), 0)
	nowDay := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	var number int64
	if err := GetDB(nil).Table(th.TableName()).Where("type IN(12) AND created_at >= ? AND ecosystem = 1", nowDay.UnixMilli()).Count(&number).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetTodayCirculationsAmount err")
		return 0
	}
	//reward, err := strconv.ParseInt(nftBlockReward, 10, 64)
	//if err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("GetTodayCirculationsAmount parse int err")
	//	return 0
	//}
	return float64(number) * nftBlockReward
}

func (th *History) Get24HourTxAmount() string {
	tz := time.Unix(GetNowTimeUnix(), 0)
	t1 := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	type result struct {
		Amount decimal.Decimal
	}
	var res result

	err := GetDB(nil).Raw(`
SELECT COALESCE(sum(amount) + 
		(
			SELECT COALESCE(sum(v4.in_amount),0)+COALESCE(sum(v4.out_amount),0)amount
			FROM(
				SELECT v3.*,CASE WHEN v3.diff >= 0 THEN 
						v3.diff 
					ELSE
						0
					END AS in_amount,
					CASE WHEN v3.diff < 0 THEN 
						v3.diff*-1 
					ELSE
						0
					END AS out_amount
				FROM(
					SELECT v1.*,COALESCE(v2.pre_sum,0)pre_sum,CASE WHEN v1.num > 1 THEN
							CASE WHEN v2.pre_sum < v1.output_value THEN
								(v1.output_value - v2.pre_sum)*-1
							ELSE
								v2.pre_sum - v1.output_value
							END
						ELSE
							0
					END diff
					 FROM(
						SELECT sum(output_value)output_value,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
							FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  WHERE ecosystem = 1 AND timestamp >= ? GROUP BY output_tx_hash 
						ORDER BY timestamp asc
					)AS v1
					LEFT JOIN(
						SELECT sum(output_value)pre_sum,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
							FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  WHERE ecosystem = 1 AND timestamp >= ? GROUP BY output_tx_hash 
						ORDER BY timestamp asc
					)AS v2 ON(v2.num+1 = v1.num)
				)AS v3
			)as v4
		),0)as amount
FROM "1_history" AS h1 WHERE ecosystem = 1 AND created_at >= ?
`, t1.UnixMilli(), t1.UnixMilli(), t1.UnixMilli()).Scan(&res).Error
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("Get scan 24 Hour tx amount err")
		return "0"
	}

	return res.Amount.String()
}

func (th *History) GetEcosystem(ecosystem int64) (bool, error) {
	return isFound(GetDB(nil).Where("type = 1 AND ecosystem = ?", ecosystem).Last(th))
}

func GetAmountChangePieChart(ecosystem int64, account string, stTime, edTime int64) (AccountAmountChangePieChart, error) {
	var (
		rets      AccountAmountChangePieChart
		err       error
		f         bool
		endTime   time.Time
		startTime time.Time
		r1        AccountAmountChangePieChart
		r2        AccountAmountChangePieChart
	)
	keyId := converter.StringToAddress(account)
	if keyId == 0 {
		return rets, errors.New("account invalid")
	}
	if stTime == 0 && edTime == 0 {
		tz := time.Unix(GetNowTimeUnix(), 0)
		endTime = time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
		const getDays = 15
		t1 := endTime.AddDate(0, 0, -1*getDays)
		startTime = t1.AddDate(0, 0, 1)
	} else {
		startTime = time.Unix(stTime, 0)
		endTime = time.Unix(edTime, 0)
	}

	f, err = isFound(GetDB(nil).Raw(fmt.Sprintf(`SELECT sum(amount) AS outcome,
(SELECT sum(amount)AS income FROM "1_history" WHERE recipient_id = %d AND ecosystem = max(h1.ecosystem) AND created_at >= %d AND created_at < %d)
FROM "1_history" AS h1 WHERE sender_id = %d AND ecosystem = %d AND created_at >= %d AND created_at < %d`,
		keyId, startTime.UnixMilli(), endTime.UnixMilli(), keyId, ecosystem, startTime.UnixMilli(), endTime.UnixMilli())).Take(&r1))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Amount Change Pie Chart Failed")
		return rets, nil
	}
	if !f {
		return rets, errors.New("unknown account:" + converter.AddressToString(keyId) + " in ecosystem:" + strconv.FormatInt(ecosystem, 10))
	}

	f, err = isFound(GetDB(nil).Raw(`
SELECT COALESCE(sum(v4.in_amount),0)income,COALESCE(sum(v4.out_amount),0)outcome
FROM(
	SELECT v3.*,CASE WHEN v3.diff >= 0 THEN 
			v3.diff 
		ELSE
			0
		END AS in_amount,
		CASE WHEN v3.diff < 0 THEN 
			v3.diff*-1 
		ELSE
			0
		END AS out_amount
	FROM(
		SELECT v1.*,COALESCE(v2.pre_sum,0)pre_sum,CASE WHEN v1.num > 1 THEN
				CASE WHEN v2.pre_sum < v1.output_value THEN
					(v1.output_value - v2.pre_sum)*-1
				ELSE
					v2.pre_sum - v1.output_value
				END
			ELSE
				0
		END diff
		 FROM(
			SELECT sum(output_value)output_value,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
				FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  
				WHERE output_key_id = ? AND ecosystem = ? AND timestamp >= ? AND timestamp < ?  GROUP BY output_tx_hash 
			ORDER BY timestamp asc
		)AS v1
		LEFT JOIN(
			SELECT sum(output_value)pre_sum,output_tx_hash,max(output_key_id) output_key_id,max(lt."timestamp") AS timestamp,row_number() OVER (ORDER BY max(lt."timestamp") ASC) num
				FROM spent_info as si LEFT JOIN log_transactions AS lt ON(lt.hash = si.output_tx_hash)  
				WHERE output_key_id = ? AND ecosystem = ? AND timestamp >= ? AND timestamp < ? GROUP BY output_tx_hash 
			ORDER BY timestamp asc
		)AS v2 ON(v2.num+1 = v1.num)
	)AS v3
)as v4
`, keyId, ecosystem, startTime.UnixMilli(), endTime.UnixMilli(), keyId, ecosystem, startTime.UnixMilli(), endTime.UnixMilli()).Take(&r2))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Amount Change Pie Chart Failed")
		return rets, nil
	}
	if !f {
		return rets, errors.New("unknown account:" + converter.AddressToString(keyId) + " in ecosystem:" + strconv.FormatInt(ecosystem, 10))
	}

	rets.Outcome = r1.Outcome.Add(r2.Outcome)
	rets.Income = r1.Income.Add(r2.Income)
	rets.TokenSymbol = Tokens.Get(ecosystem)
	return rets, nil
}

func GetAmountChangeBarChart(ecosystem int64, account string, isDefault int) (AccountAmountChangeBarChart, error) {
	var (
		rets     AccountAmountChangeBarChart
		findTime int64
		his      History
		list     *AccountHistoryChart
		err      error
		balance  AccountAmountChangeBarChart
	)
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	keyId := converter.StringToAddress(account)
	if keyId == 0 && account != "0000-0000-0000-0000-0000" {
		return rets, errors.New("account invalid")
	}
	if isDefault == 1 {
		findTime = today.AddDate(0, 0, -1*14).UnixMilli()
		list, err = his.GetWalletTimeLineHistoryTotals(ecosystem, keyId, 15)
		if err != nil {
			return rets, err
		}
	} else {
		list, err = his.GetWalletTimeLineHistoryTotals(ecosystem, keyId, 0)
		if err != nil {
			return rets, err
		}
	}

	getBalance := func(dayTime int64, info AccountAmountChangeBarChart) string {
		for i := 0; i < len(info.Time); i++ {
			if dayTime == info.Time[i] {
				return info.Balance[i]
			}
		}
		return "0"
	}

	getAmount := func(dayTime int64, info *AccountHistoryChart, getInAmount bool) string {
		if info == nil {
			return "0"
		}
		for i := 0; i < len(info.Time); i++ {
			if dayTime == info.Time[i] {
				if getInAmount {
					return info.Inamount[i]
				} else {
					return info.Outamount[i]
				}
			}
		}
		return "0"
	}

	balance, err = GetAccountTokenChangeChart(ecosystem, keyId, findTime)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Amount Change Bar Chart balance list Failed")
		return rets, err
	}
	rets.TokenSymbol, rets.Name = Tokens.Get(ecosystem), EcoNames.Get(ecosystem)

	lastBalance := decimal.Zero
	var startTime time.Time
	var stTime string

	if findTime > 0 {
		stTime = time.UnixMilli(findTime).Format("2006-01-02")
	} else {
		stTime = time.Unix(FirstBlockTime, 0).Format("2006-01-02")
	}
	t1, err := time.ParseInLocation("2006-01-02", stTime, time.Local)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "day": stTime}).Error("Get Amount Change Bar Chart startTime parseInLocation Failed")
		return rets, err
	}
	startTime = t1
	for startTime.Unix() <= today.Unix() {
		rets.Time = append(rets.Time, startTime.Unix())
		balanceAmount := getBalance(startTime.Unix(), balance)

		//balance := GetDaysAmount(startTime.Unix(), balanceList)
		if balanceAmount != "0" {
			lastBalance, _ = decimal.NewFromString(balanceAmount)
			rets.Balance = append(rets.Balance, balanceAmount)
		} else {
			rets.Balance = append(rets.Balance, lastBalance.String())
		}
		rets.Outcome = append(rets.Outcome, getAmount(startTime.Unix(), list, false))
		rets.Income = append(rets.Income, getAmount(startTime.Unix(), list, true))
		startTime = startTime.AddDate(0, 0, 1)
	}

	return rets, nil
}

func GetAccountTxChart(ecosystem int64, account string) (AccountTxChart, error) {
	var (
		rets AccountTxChart
		list []DaysNumber
	)
	keyId := converter.StringToAddress(account)
	if keyId == 0 && account != "0000-0000-0000-0000-0000" {
		return rets, errors.New("account invalid")
	}
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	err := GetDB(nil).Raw(`
SELECT CASE WHEN n1.days = ''THEN
	n2.days
ELSE
	n1.days
END,COALESCE(n1.num,0)+COALESCE(n2.num,0)AS num
FROM(
	SELECT days,count(1)num FROM (
		SELECT txhash,to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') AS days FROM "1_history" WHERE (sender_id = ? or recipient_id = ?) AND ecosystem = ? GROUP BY txhash,created_at
		ORDER BY days asc
	)AS v1
	GROUP BY days
)AS n1
FULL JOIN(

	SELECT days,count(1)num FROM(
		SELECT output_tx_hash,to_char(to_timestamp(timestamp/1000),'yyyy-MM-dd') AS days FROM spent_info AS v1 
			LEFT JOIN log_transactions AS v2 ON(v1.output_tx_hash = v2.hash) WHERE output_key_id = ? AND ecosystem = ?
		GROUP BY output_tx_hash,"timestamp"
		ORDER BY days asc
	)AS v1
	GROUP BY days
)AS n2 ON(n2.days = n1.days)
`, keyId, keyId, ecosystem, keyId, ecosystem).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Amount Change Bar Chart balance list Failed")
		return rets, nil
	}
	rets.Time = make([]int64, 0)
	rets.Tx = make([]int64, 0)

	var startTime time.Time
	if len(list) >= 1 {
		t1, err := time.ParseInLocation("2006-01-02", list[0].Days, time.Local)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "day": list[0].Days}).Error("Get Account Tx Chart startTime parseInLocation Failed")
			return rets, err
		}
		startTime = t1
		for startTime.Unix() <= today.Unix() {
			rets.Time = append(rets.Time, startTime.Unix())
			rets.Tx = append(rets.Tx, GetDaysNumber(startTime.Unix(), list))
			startTime = startTime.AddDate(0, 0, 1)
		}
	}

	return rets, nil
}
