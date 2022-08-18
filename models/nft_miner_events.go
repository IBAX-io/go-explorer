/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type NftMinerEvents struct {
	ID           int64  `gorm:"primary_key;not null"`
	TokenId      int64  `gorm:"column:token_id;not null"`
	TokenHash    []byte `gorm:"column:token_hash;not null"`
	Event        string `gorm:"column:event;not null"`
	ContractName string `gorm:"column:contract_name;not null"`
	DateCreated  int64  `gorm:"column:date_created;not null"`
	TxHash       []byte `gorm:"column:tx_hash;not null"`
	Source       string `gorm:"column:source"`
}

func (p *NftMinerEvents) TableName() string {
	return "1_nft_miner_events"
}

func (p *NftMinerEvents) GetByTokenId(tokenId int64) (bool, error) {
	return isFound(GetDB(nil).Where("token_id = ?", tokenId).First(&p))
}

func (p *NftMinerEvents) GetByTokenHash(tokenHash string) (bool, error) {
	hash, _ := hex.DecodeString(tokenHash)
	return isFound(GetDB(nil).Where("token_hash = ?", hash).First(&p))
}

func (p *NftMinerEvents) GetNftHistoryInfo(search any, page, limit int, order string) (*GeneralResponse, error) {
	var (
		ret    []NftMinerHistoryInfoResponse
		rets   GeneralResponse
		total  int64
		events []NftMinerEvents
		f      bool
		err    error
	)
	if order == "" {
		order = "id desc"
	}
	switch reflect.TypeOf(search).String() {
	case "string":
		hash, _ := hex.DecodeString(search.(string))
		err = GetDB(nil).Table(p.TableName()).Where("token_hash = ?", hash).Count(&total).Error
		if err != nil {
			return nil, err
		}

		f, err = isFound(GetDB(nil).Table(p.TableName()).Where("token_hash = ?", hash).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&events))
		if err != nil {
			return nil, err
		}
	case "json.Number":
		tokenId, _ := search.(json.Number).Int64()
		if tokenId != 0 {
			err = GetDB(nil).Table(p.TableName()).Where("token_id = ?", tokenId).Count(&total).Error
			if err != nil {
				return nil, err
			}

			f, err = isFound(GetDB(nil).Table(p.TableName()).Where("token_id = ?", tokenId).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&events))
			if err != nil {
				return nil, err
			}
		}

	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Nft Miner History Search Failed")
		return nil, errors.New("request params invalid")
	}

	if !f {
		return nil, errors.New("NFT Miner Doesn't Not Exist")
	}
	for _, value := range events {
		var es NftMinerHistoryInfoResponse
		es.Events = value.Event
		es.TxHash = hex.EncodeToString(value.TxHash)
		es.ID = value.ID
		es.NftId = value.TokenId
		es.Time = value.DateCreated
		es.Source = value.Source
		es.Contract = value.ContractName
		es.NftHash = hex.EncodeToString(value.TokenHash)
		ret = append(ret, es)
	}

	rets.Page = page
	rets.Limit = limit
	rets.Total = total
	rets.List = ret
	return &rets, nil
}
