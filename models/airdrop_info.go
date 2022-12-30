package models

import (
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type AirdropInfo struct {
	Id            int64 `gorm:"primary_key;not_null"`
	Account       string
	BalanceAmount decimal.Decimal
	DateCreated   int64
	Detail        string `gorm:"type:jsonb"`
	DirectAmount  decimal.Decimal
	LatestAt      int64
	LockAmount    decimal.Decimal
	PeriodCount   int64
	Priority      int64
	StakeAmount   decimal.Decimal
	Surplus       int64
	TotalAmount   decimal.Decimal
}

var (
	AirdropReady bool
	//total amount: is not now amount
	AirdropLockAll decimal.Decimal

	nowAirdropLockAll    decimal.Decimal
	nowAirdropStakingAll decimal.Decimal
)

func (p AirdropInfo) TableName() string {
	return `1_airdrop_info`
}

func AirdropTableExist() bool {
	var p AirdropInfo
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}

func (p *AirdropInfo) Get(account string) (bool, error) {
	return isFound(GetDB(nil).Where("account = ?", account).First(p))
}

func (p *AirdropInfo) GetStaking(account string) decimal.Decimal {
	var staking decimal.Decimal
	f, err := isFound(GetDB(nil).Select("stake_amount").Where("account = ?", account).Take(&staking))
	if err == nil && f {
		return staking
	}
	return decimal.Zero
}

func GetAirdropLockAllTotal() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	if AirdropReady {
		err := GetDB(nil).Model(AirdropInfo{}).Select("COALESCE(sum(lock_amount),0)").Take(&AirdropLockAll).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Airdrop Lock All Total Failed")
		}
		getNowAirdropLockAll()
		getNowAirdropStakingAll()
	}
}

func getNowAirdropLockAll() {
	err := GetDB(nil).Model(AirdropInfo{}).Select("COALESCE(sum(balance_amount),0)").Take(&nowAirdropLockAll).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Airdrop Now Lock All Total Failed")
	}
}

func getNowAirdropStakingAll() {
	err := GetDB(nil).Model(AirdropInfo{}).Select("coalesce(sum(stake_amount),0)").Take(&nowAirdropStakingAll).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Airdrop Now Staking All Total Failed")
	}
}
