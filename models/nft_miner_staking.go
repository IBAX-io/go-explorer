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
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type NftMinerStaking struct {
	ID            int64           `gorm:"primary_key;not null"`         //ID
	TokenId       int64           `gorm:"column:token_id;not null"`     //NFT Miner ID
	StakeAmount   decimal.Decimal `gorm:"column:stake_amount;not null"` //starking
	EnergyPower   int64           `gorm:"column:energy_power;not null"`
	EnergyPoint   int64           `gorm:"column:energy_point;not null"`
	Source        int64           `gorm:"column:source;not null"`      //source
	StartDated    int64           `gorm:"column:start_dated;not null"` //start time
	EndDated      int64           `gorm:"column:end_dated;not null"`   //end time
	Staker        string          `gorm:"column:staker;not null"`      //owner account
	StakingStatus int64           `gorm:"column:staking_status;not null"`
	WithdrawDate  int64           `gorm:"column:withdraw_date;not null"` //withdraw time

}

func (p *NftMinerStaking) TableName() string {
	return "1_nft_miner_staking"
}

func (p *NftMinerStaking) GetAllStakeAmount() (int64, decimal.Decimal, error) {
	var nftStaking SumAmount
	var nftStakingNum int64
	zero := decimal.New(0, 0)
	if !NftMinerReady {
		return 0, zero, nil
	}
	if err := GetDB(nil).Table(p.TableName()).Where("staking_status = 1").Count(&nftStakingNum).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, zero, err
		}
	}
	if err := GetDB(nil).Table(p.TableName()).Select("coalesce(sum(stake_amount),'0')as sum").Where("staking_status = 1").Take(&nftStaking.Sum).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, zero, err
		}
	}
	return nftStakingNum, nftStaking.Sum, nil
}

func (p *NftMinerStaking) GetAllStakeAmountByStaker(keyid string) (int64, string, error) {
	var nftStaking string
	var nftStakingNum int64
	if !HasTableOrView(p.TableName()) {
		return 0, "", nil
	}
	if err := GetDB(nil).Table(p.TableName()).Where("staking_status = 1 AND staker = ?", keyid).Count(&nftStakingNum).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, "", err
		}
	}
	if err := GetDB(nil).Table(p.TableName()).Select("sum(stake_amount)").Where("staking_status = 1 AND staker = ?", keyid).Take(&nftStaking).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, "", err
		}
	}
	return nftStakingNum, nftStaking, nil
}

func (p *NftMinerStaking) GetById(id int64) (bool, error) {
	return isFound(GetDB(nil).Where("id = ?", id).First(p))
}

func (p *NftMinerStaking) GetByTokenId(tokenId int64) (bool, error) {
	return isFound(GetDB(nil).Where("token_id = ?", tokenId).First(p))
}

func (p *NftMinerStaking) GetNftStakeByTokenId(tokenId int64) (bool, error) {
	return isFound(GetDB(nil).Where("token_id = ? AND staking_status = 1", tokenId).First(p))
}

func (p *NftMinerStaking) GetNftMinerStakeInfo(search any, page, limit int, order string) (*GeneralResponse, error) {
	var (
		rets  []NftMinerStakeInfoResponse
		total int64
		ret   GeneralResponse
		nftId int64
	)
	if order == "" {
		order = "id desc"
	}

	switch reflect.TypeOf(search).String() {
	case "string":
		var item NftMinerItems
		f, err := item.GetByTokenHash(search.(string))
		if err != nil {
			log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Nft Miner Stake Info  Failed")
			return nil, err
		}
		if !f {
			return nil, errors.New("NFT Miner Doesn't Not Exist")
		}
		nftId = item.ID
	case "json.Number":
		tokenId, err := search.(json.Number).Int64()
		if err != nil {
			return nil, err
		}
		nftId = tokenId
	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Nft Miner Stake Info Search Failed")
		return nil, errors.New("request params invalid")
	}

	err := GetDB(nil).Table(p.TableName()).Where("token_id = ? ", nftId).Count(&total).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft Miner stake info total err:", err.Error(), " nftId:", nftId)
		}
		return nil, err
	}

	type nftStakedResponse struct {
		Id          int64
		NftId       int64
		StartDated  int64
		EndDated    int64
		EnergyPower decimal.Decimal
		Cycle       int64
		Hash        []byte
		Block       int64
		StakeAmount int64
	}
	var stak []nftStakedResponse
	err = GetDB(nil).Raw(`SELECT sk.id,sk.token_id AS nft_id,sk.start_dated,sk.end_dated,
sk.energy_power,sk.stake_amount,date_part('day',cast(to_char(to_timestamp(end_dated),'yyyy-MM-dd') as TIMESTAMP)-cast(to_char(to_timestamp(start_dated),'yyyy-MM-dd') as TIMESTAMP)) 
AS cycle,"log_transactions".hash,"log_transactions".block from "1_nft_miner_staking" 
AS sk left JOIN "log_transactions" ON (encode(hash, 'hex')= sk.tx_hash) WHERE token_id = ? order by ? offset ? limit ?`, nftId, order, (page-1)*limit, limit).Find(&stak).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft Miner stake info nftStaking err:", err.Error(), " nft Miner Id:", nftId)
		}
		return nil, err
	}

	nowTime := time.Now().Unix()
	rets = make([]NftMinerStakeInfoResponse, len(stak))
	for i := 0; i < len(stak); i++ {
		rets[i].ID = stak[i].Id
		rets[i].NftId = stak[i].NftId
		if nowTime >= stak[i].StartDated && nowTime <= stak[i].EndDated {
			rets[i].StakeStatus = true
			rets[i].EnergyPower = stak[i].EnergyPower.String()
		}
		rets[i].StakeAmount = stak[i].StakeAmount
		rets[i].Cycle = int64(time.Unix(stak[i].EndDated, 0).Sub(time.Unix(stak[i].StartDated, 0)).Hours() / 24)
		rets[i].Time = stak[i].StartDated

		if stak[i].Hash != nil {
			rets[i].TxHash = hex.EncodeToString(stak[i].Hash)
		}
		rets[i].BlockId = stak[i].Block
	}

	ret.Total = total
	ret.Page = page
	ret.Limit = limit
	ret.List = rets

	return &ret, nil
}
