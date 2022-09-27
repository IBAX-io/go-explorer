/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/shopspring/decimal"
	"strconv"
)

// SpentInfo is model
type SpentInfo struct {
	InputTxHash  []byte `gorm:"default:(-)"`
	InputIndex   int32
	OutputTxHash []byte `gorm:"not null"`
	OutputIndex  int32  `gorm:"not null"`
	OutputKeyId  int64  `gorm:"not null"`
	OutputValue  string `gorm:"not null"`
	Ecosystem    int64
	BlockId      int64
}

// TableName returns name of table
func (si *SpentInfo) TableName() string {
	return "spent_info"
}

func (si *SpentInfo) GetAmountByKeyId(keyId int64, ecosystem int64) (decimal.Decimal, error) {
	var utxoAmount SumAmount
	f, err := isFound(GetDB(nil).Table(si.TableName()).
		Select("coalesce(sum(output_value),'0') as sum").Where("input_tx_hash is NULL AND output_key_id = ? AND ecosystem = ?", keyId, ecosystem).
		Take(&utxoAmount))
	if err != nil {
		return decimal.Zero, err
	}
	if !f {
		return decimal.Zero, nil
	}
	return utxoAmount.Sum, nil
}

func (si *SpentInfo) GetExplorer(txHash []byte) (*UtxoExplorer, error) {
	var (
		txData      TransactionData
		rets        UtxoExplorer
		ecoGasExist bool
		outputList  []SpentInfo
	)

	f, err := txData.GetTxDataByHash(txHash)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errors.New("waiting for transactions to sync")
	}

	info, err := si.UnmarshalTransaction(txData.TxData)
	if err != nil {
		return nil, err
	}
	rets.UtxoType = info.UtxoType
	rets.Sender = info.Sender
	rets.Recipient = info.Recipient
	rets.Amount = info.Amount
	rets.Comment = info.Comment
	rets.Expedite = info.Expedite
	rets.TokenSymbol = info.TokenSymbol
	rets.Ecosystem = info.Ecosystem
	rets.Size = strconv.FormatInt(info.Size, 10) + " bit"

	rets.Inputs, err = si.GetInputs(txHash, converter.StringToAddress(info.Sender))
	if err != nil {
		return nil, err
	}
	rets.Outputs, outputList, err = si.GetOutputs(txHash)
	if err != nil {
		return nil, err
	}

	if rets.UtxoType == UtxoTx {
		var (
			index       int
			indexSet    bool
			changeList  []utxoDetail
			ecoGasFee   FeesInfo
			basisGasFee FeesInfo
			unit        = "/bit"
			ecoCount    int
		)

		for _, v := range outputList {
			if v.Ecosystem != 1 {
				ecoCount += 1
			}
		}
		if ecoCount >= 3 {
			ecoGasExist = true
		}

		for _, v := range outputList {
			amount, _ := decimal.NewFromString(v.OutputValue)
			recipient := converter.AddressToString(v.OutputKeyId)
			outputTxHash := hex.EncodeToString(v.OutputTxHash)
			if rets.Ecosystem == 1 {
				if v.Ecosystem == 1 {
					switch index {
					case 0:
						basisGasFee.Fees.Amount = amount.String()
						basisGasFee.Fees.Recipient = recipient
						basisGasFee.Fees.Sender = rets.Sender
						basisGasFee.Fees.TokenSymbol = rets.TokenSymbol

						basisGasFee.TokenSymbol = rets.TokenSymbol
						basisGasFee.Amount = basisGasFee.Amount.Add(amount)
						index += 1
					case 1:
						basisGasFee.Taxes.Amount = amount.String()
						basisGasFee.Taxes.Recipient = recipient
						basisGasFee.Taxes.Sender = rets.Sender
						basisGasFee.Taxes.TokenSymbol = rets.TokenSymbol

						basisGasFee.TokenSymbol = rets.TokenSymbol
						basisGasFee.Amount = basisGasFee.Amount.Add(amount)
						index += 1
					case 2:
						index += 1
					case 3:
						var change utxoDetail
						change.Address = recipient
						change.TokenSymbol = rets.TokenSymbol
						change.Amount = amount.String()
						change.Hash = outputTxHash

						changeList = append(changeList, change)
					}
				}
			} else {
				if v.Ecosystem == 1 {
					switch index {
					case 0:
						basisGasFee.Fees.Amount = amount.String()
						basisGasFee.Fees.Recipient = recipient
						basisGasFee.Fees.Sender = rets.Sender
						basisGasFee.Fees.TokenSymbol = Tokens.Get(v.Ecosystem)

						basisGasFee.TokenSymbol = basisGasFee.Fees.TokenSymbol
						basisGasFee.Amount = basisGasFee.Amount.Add(amount)
						index += 1
					case 1:
						basisGasFee.Taxes.Amount = amount.String()
						basisGasFee.Taxes.Recipient = recipient
						basisGasFee.Taxes.Sender = rets.Sender
						basisGasFee.Taxes.TokenSymbol = Tokens.Get(v.Ecosystem)

						basisGasFee.TokenSymbol = basisGasFee.Taxes.TokenSymbol
						basisGasFee.Amount = basisGasFee.Amount.Add(amount)
						index += 1
					case 2:
						var change utxoDetail
						change.Address = recipient
						change.TokenSymbol = Tokens.Get(v.Ecosystem)
						change.Amount = amount.String()
						change.Hash = outputTxHash

						changeList = append(changeList, change)
					}
				} else {
					if !indexSet {
						if ecoGasExist {
							index = 0
						} else {
							index = 2
						}
						indexSet = true
					}
					switch index {
					case 0:
						ecoGasFee.Fees.Amount = amount.String()
						ecoGasFee.Fees.Recipient = recipient
						ecoGasFee.Fees.Sender = rets.Sender
						ecoGasFee.Fees.TokenSymbol = rets.TokenSymbol

						ecoGasFee.TokenSymbol = rets.TokenSymbol
						ecoGasFee.Amount = ecoGasFee.Amount.Add(amount)
						index += 1
					case 1:
						ecoGasFee.Taxes.Amount = amount.String()
						ecoGasFee.Taxes.Recipient = recipient
						ecoGasFee.Taxes.Sender = rets.Sender
						ecoGasFee.Taxes.TokenSymbol = rets.TokenSymbol

						ecoGasFee.TokenSymbol = rets.TokenSymbol
						ecoGasFee.Amount = ecoGasFee.Amount.Add(amount)
						index += 1
					case 2:
						index += 1
					case 3:
						var change utxoDetail
						change.Address = recipient
						change.TokenSymbol = rets.TokenSymbol
						change.Amount = amount.String()
						change.Hash = outputTxHash

						changeList = append(changeList, change)
					}
				}
			}
		}
		rets.Change = changeList
		txSize := decimal.NewFromInt(info.Size)
		if ecoGasFee.Amount.GreaterThan(decimal.Zero) {
			ecoGasFee.FuelRate = FuelRateResponse{ecoGasFee.Amount.DivRound(txSize, 0).String(), ecoGasFee.TokenSymbol + unit}
		}
		if basisGasFee.Amount.GreaterThan(decimal.Zero) {
			basisGasFee.FuelRate = FuelRateResponse{basisGasFee.Amount.DivRound(txSize, 0).String(), basisGasFee.TokenSymbol + unit}
		}
		rets.EcoGasFee = ecoGasFee
		rets.BasisGasFee = basisGasFee
	}

	return &rets, nil
}

func (si *SpentInfo) UnmarshalTransaction(txData []byte) (*UtxoExplorerInfo, error) {
	if len(txData) == 0 {
		return nil, errors.New("tx data length is empty")
	}
	var result UtxoExplorerInfo
	tx, err := UnmarshallTransaction(bytes.NewBuffer(txData))
	if err != nil {
		return nil, err
	}

	if tx.IsSmartContract() {
		result.Ecosystem = tx.SmartContract().TxSmart.Header.EcosystemID
		result.TokenSymbol = Tokens.Get(result.Ecosystem)
		result.Expedite = tx.SmartContract().TxSmart.Expedite
		result.Size = int64(len(tx.Payload()))

		if tx.SmartContract().TxSmart.UTXO != nil {
			result.Comment = tx.SmartContract().TxSmart.UTXO.Comment
			result.Sender = converter.AddressToString(tx.KeyID())
			result.Recipient = converter.AddressToString(tx.SmartContract().TxSmart.UTXO.ToID)
			result.Amount = tx.SmartContract().TxSmart.UTXO.Value
			result.UtxoType = UtxoTx
		} else if tx.SmartContract().TxSmart.TransferSelf != nil {
			result.Sender = converter.AddressToString(tx.KeyID())
			result.Recipient = converter.AddressToString(tx.KeyID())
			if tx.SmartContract().TxSmart.TransferSelf.Source == "Account" && tx.SmartContract().TxSmart.TransferSelf.Target == "UTXO" {
				result.Comment = "Account-UTXO"
			} else {
				result.Comment = "UTXO-Account"
			}
			result.Amount = tx.SmartContract().TxSmart.TransferSelf.Value
			result.UtxoType = UtxoTransfer
		} else {
			return &result, errors.New("doesn't not UTXO transaction")
		}
	} else {
		return &result, errors.New("doesn't not Smart Contract transaction")
	}
	return &result, nil
}

func (si *SpentInfo) GetInputs(txHash []byte, kid int64) (rlts []utxoDetail, err error) {
	var (
		list []SpentInfo
	)
	err = GetDB(nil).Table(si.TableName()).Select("output_key_id,output_value,output_tx_hash,ecosystem").Where("input_tx_hash = ? AND output_key_id = ?", txHash, kid).Find(&list).Error
	for _, val := range list {
		rlts = append(rlts, utxoDetail{Address: converter.AddressToString(val.OutputKeyId), Amount: val.OutputValue, Hash: hex.EncodeToString(val.OutputTxHash),
			TokenSymbol: Tokens.Get(val.Ecosystem)})
	}
	return
}

func (si *SpentInfo) GetOutputs(txHash []byte) (rlts []utxoDetail, list []SpentInfo, err error) {
	err = GetDB(nil).Table(si.TableName()).
		Where("output_tx_hash = ?", txHash).Order("output_index ASC").Find(&list).Error
	for _, val := range list {
		rlts = append(rlts, utxoDetail{Address: converter.AddressToString(val.OutputKeyId), Amount: val.OutputValue, Hash: hex.EncodeToString(val.OutputTxHash), TokenSymbol: Tokens.Get(val.Ecosystem)})
	}
	return
}

func (si *SpentInfo) GetLast() (bool, error) {
	return isFound(GetDB(nil).Order("block_id desc").Take(si))
}

func getSpentInfoHashList(startId int64, limit int, order string) (*[]spentInfoTxData, error) {
	var (
		err error
	)
	var rlt []spentInfoTxData

	orderStr := "block_id " + string(order)

	err = GetDB(nil).Raw(fmt.Sprintf(`
SELECT v1.output_tx_hash,v2.tx_data,v1.block_id,v2.tx_time FROM(
	SELECT output_tx_hash,block_id FROM spent_info WHERE block_id > ? 
	GROUP BY output_tx_hash,block_id ORDER BY %s Limit ?
)AS v1
LEFT JOIN(
	SELECT tx_data,tx_time,hash,block FROM transaction_data
)AS v2 ON(v2.hash = v1.output_tx_hash)
WHERE v2.block > 0
`, orderStr), startId, limit).Find(&rlt).Error
	if err != nil {
		return nil, err
	}

	return &rlt, nil
}
