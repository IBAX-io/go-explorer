/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/shopspring/decimal"
)

// SpentInfo is model
type SpentInfo struct {
	InputTxHash  []byte `gorm:"default:(-)"`
	InputIndex   int32
	OutputTxHash []byte `gorm:"not null"`
	OutputIndex  int32  `gorm:"not null"`
	OutputKeyId  int64  `gorm:"not null"`
	OutputValue  string `gorm:"not null"`
	Scene        string
	Ecosystem    int64
	Contract     string
	BlockId      int64
	Asset        string
	Action       string `gorm:"-"` // UTXO operation control : change
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
		txData TransactionData
		rets   UtxoExplorer
	)

	f, err := txData.GetTxDataByHash(txHash)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errors.New("waiting for transactions to sync")
	}

	info, err := si.UnmarshalTxTransaction(txData.TxData)
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
	fuels := GetFuelRate()
	if _, ok := fuels[rets.Ecosystem]; ok {
		if rets.Ecosystem == 1 {
			rets.BasisFuelRate = FuelRate{1, fuels[1].DivRound(decimal.NewFromInt(10), 0).String(), SysTokenSymbol}
		} else {
			rets.EcoFuelRate = FuelRate{rets.Ecosystem, fuels[rets.Ecosystem].DivRound(decimal.NewFromInt(10), 0).String(),
				Tokens.Get(rets.Ecosystem)}
		}
	}
	if rets.Ecosystem != 1 {
		if _, ok := fuels[1]; ok {
			rets.BasisFuelRate = FuelRate{1, fuels[1].DivRound(decimal.NewFromInt(10), 0).String(), SysTokenSymbol}
		}
	}
	rets.Inputs, err = si.GetInputs(txHash, converter.StringToAddress(info.Sender))
	if err != nil {
		return nil, err
	}
	rets.Outputs, err = si.GetOutputs(txHash)
	if err != nil {
		return nil, err
	}

	for _, val := range rets.Outputs {
		if rets.TokenSymbol == val.TokenSymbol && rets.Sender == val.Address && val.TokenSymbol != SysTokenSymbol {
			rets.Change = append(rets.Change, val)
		}
		if rets.TokenSymbol == val.TokenSymbol && rets.Sender == val.Address && val.TokenSymbol == SysTokenSymbol {
			rets.Change = append(rets.Change, val)
		}
	}

	return &rets, nil
}

func (si *SpentInfo) UnmarshalTxTransaction(txData []byte) (*UtxoExplorer, error) {
	if len(txData) == 0 {
		return nil, errors.New("tx data length is empty")
	}
	var result UtxoExplorer
	tx, err := UnmarshallTransaction(bytes.NewBuffer(txData))
	if err != nil {
		return nil, err
	}

	if tx.IsSmartContract() {
		result.Ecosystem = tx.SmartContract().TxSmart.Header.EcosystemID
		result.TokenSymbol = Tokens.Get(result.Ecosystem)
		result.Expedite = tx.SmartContract().TxSmart.Expedite
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

func (si *SpentInfo) GetOutputs(txHash []byte) (rlts []utxoDetail, err error) {
	var (
		list []SpentInfo
	)
	err = GetDB(nil).Table(si.TableName()).Select("output_key_id,output_value,output_tx_hash,ecosystem").Where("output_tx_hash = ?", txHash).Find(&list).Error
	for _, val := range list {
		rlts = append(rlts, utxoDetail{Address: converter.AddressToString(val.OutputKeyId), Amount: val.OutputValue, Hash: hex.EncodeToString(val.OutputTxHash), TokenSymbol: Tokens.Get(val.Ecosystem)})
	}
	return
}
