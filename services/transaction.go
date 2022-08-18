/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"encoding/hex"
	"errors"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/shopspring/decimal"
)

func GetGroupTransactionStatus(ids int, icount int, order string) (*[]models.TransactionStatusHex, int64, error) {
	ts := &models.TransactionStatus{}
	ret, num, err := ts.GetTransactions(ids, icount, order)
	return ret, num, err
}

func GetGroupTransactionHistory(ids int, icount int, order string) (*[]models.HistoryHex, int64, error) {
	ts := &models.History{}
	ret, num, err := ts.GetHistory(ids, icount, order)
	return ret, num, err
}

func GetGroupTransactionWallet(ids int, icount int, wallet string, searchType string) (*[]models.HistoryHex, int64, decimal.Decimal, error) {
	ts := &models.History{}

	ret, num, total, err := ts.GetWallets(ids, icount, wallet, searchType)
	return ret, num, total, err
}

func GetGroupTransactionEcosystemWallet(id int64, ids int, icount int, wallet string, searchType string) (*[]models.HistoryHex, int64, decimal.Decimal, error) {
	ts := &models.History{}
	return ts.GetEcosytemWallets(id, ids, icount, wallet, searchType)
}

func GetGroupWalletHistory(id int64, wallet string) (*models.WalletHistoryHex, error) {
	ts := &models.History{}
	key := &models.Key{}
	ret, err := ts.GetWalletTotals(wallet)
	if err != nil {
		return ret, err
	}

	dat, err := key.Get(id, wallet)
	if err != nil {
		return ret, err
	}
	amount := decimal.New(0, 0)
	if len(dat.Amount) > 0 {
		amount, _ = decimal.NewFromString(dat.Amount)
	}
	ret.Amount = amount

	return ret, err
}

func GetGroupWalletTotal(ids int, icount int, order string, wallet string) (int64, int, *[]models.EcosyKeyTotalHex, error) {
	key := &models.Key{}
	return key.GetTotal(ids, icount, order, wallet)
}

func GetTransactionDetailedInfoHash(hash string) (*models.HistoryExplorer, error) {
	hashByte := converter.HexToBin(hash)
	ts := &models.History{}
	ret2, err2 := ts.GetExplorer(hashByte)
	if err2 != nil {
		return nil, err2
	}
	return ret2, nil
}

func GetTransactionHeadInfoHash(hash string) (*models.TxDetailedInfoHeadResponse, error) {
	var rets models.TxDetailedInfoHeadResponse

	var lt models.LogTransaction
	hashHex, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	f, err := lt.GetByHash(hashHex)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errors.New("unknown tx hash")
	}
	var txData models.TransactionData
	f, err = txData.GetTxDataByHash(hashHex)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errors.New("waiting for transactions to sync")
	}

	info, err := lt.UnmarshalTxTransaction(txData.TxData)
	if err != nil {
		return nil, err
	}
	rets.ContractCode = models.GetContractCodeByName(info.ContractName)
	rets.LogoHash, rets.TokenSymbol = models.GetEcosystemLogoHash(info.Ecosystem)
	rets.Hash = info.Hash
	rets.Ecosystem = info.Ecosystem
	rets.Time = info.Time
	rets.ContractName = info.ContractName
	rets.EcosystemName = info.Ecosystemname
	rets.Params = info.Params
	rets.BlockID = lt.Block
	rets.Address = converter.AddressToString(lt.Address)
	rets.Size = models.TocapacityString(info.Size)

	return &rets, nil
}
