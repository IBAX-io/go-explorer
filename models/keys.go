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
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/shopspring/decimal"

	//"github.com/IBAX-io/go-explorer/conf"
	"strconv"
	"strings"

	"github.com/IBAX-io/go-ibax/packages/converter"
)

// Key is model
type Key struct {
	Ecosystem int64
	ID        int64           `gorm:"primary_key;not null"`
	PublicKey []byte          `gorm:"column:pub;not null"`
	Amount    decimal.Decimal `gorm:"not null"`
	Maxpay    decimal.Decimal `gorm:"not null"`
	Deposit   decimal.Decimal `gorm:"not null"`
	Multi     int64           `gorm:"not null"`
	Deleted   int64           `gorm:"not null"`
	Blocked   int64           `gorm:"not null"`
	AccountID string          `gorm:"column:account;not null"`
	Lock      string          `gorm:"column:lock;type:jsonb"`
}

type KeyHex struct {
	Ecosystem     int64           `json:"ecosystem"`
	ID            string          `json:"id"`
	PublicKey     string          `json:"publickey"`
	Amount        string          `json:"amount"`
	Maxpay        decimal.Decimal `json:"maxpay"`
	Multi         int64           `json:"multi"`
	Deleted       int64           `json:"deleted"`
	Blocked       int64           `json:"blocked"`
	Ecosystemname string          `json:"ecosystemname"`
	TokenSymbol   string          `json:"token_symbol"`
}

type EcosyKeyHex struct {
	Ecosystem int64 `json:"ecosystem"`
	//Ecosyname string `json:"Ecosyname"`
	IsValued        int64  `json:"isvalued"`
	Ecosystemname   string `json:"ecosystemname"`
	TokenSymbol     string `json:"token_symbol"`
	Amount          string `json:"amount"`
	Info            string `json:"info"`
	Emission_amount string `json:"emission_amount"`
	Type_emission   int64
	Type_withdraw   int64
}

type EcosyKeyTotalHex struct {
	Chart           *AccountHistoryChart `json:"chart,omitempty"`
	Wallet          string               `json:"wallet"`
	Ecosystem       int64                `json:"ecosystem"`
	IsValued        int64                `json:"isvalued"`
	Ecosystemname   string               `json:"ecosystemname"`
	TokenSymbol     string               `json:"token_symbol"`
	Amount          string               `json:"amount"`
	Info            string               `json:"info"`
	Emission_amount string               `json:"emission_amount"`
	MemberName      string               `json:"member_name"`
	MemberHash      string               `json:"member_hash"`
	LogoHash        string               `json:"logo_hash"`
	Type_emission   int64
	Type_withdraw   int64
	Transaction     int64           `json:"transaction"`
	InTx            int64           `json:"in_tx"`
	OutTx           int64           `json:"out_tx"`
	StakeAmount     decimal.Decimal `json:"stake_amount"`
	FreezeAmount    decimal.Decimal `json:"freeze_amount"`
	Inamount        decimal.Decimal `json:"inamount"`
	Outamount       decimal.Decimal `json:"outamount"`
	Totalamount     decimal.Decimal `json:"total_amount"`
	JoinTime        int64           `json:"join_time"`
	RolesName       string          `json:"roles_name"`
}

type EcosyKeyTotalDetail struct {
	Account      string          `json:"account"`
	Ecosystem    int64           `json:"ecosystem"`
	LogoHash     string          `json:"logo_hash"`
	Name         string          `json:"name"`
	JoinTime     int64           `json:"join_time"`
	TokenSymbol  string          `json:"token_symbol"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	FreezeAmount decimal.Decimal `json:"freeze_amount"`
	StakeAmount  decimal.Decimal `json:"stake_amount"`
	RolesName    string          `json:"roles_name"`
}

type AccountHistoryChart struct {
	Time      []int64  `json:"time"`
	Inamount  []string `json:"inamount"`
	Outamount []string `json:"outamount"`
}

type KeysResult struct {
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Rets  []EcosyKeyTotalHex `json:"rets"`
}

type EcosyKeyList struct {
	Ecosystem    int64           `json:"ecosystem"`
	Account      string          `json:"account"`
	AccountName  string          `json:"account_name"`
	Amount       string          `json:"amount"`
	AccountedFor decimal.Decimal `json:"accounted_for"`
	StakeAmount  string          `json:"stake_amount"`
	FreezeAmount string          `json:"freeze_amount"`
	TokenSymbol  string          `json:"token_symbol"`
}

type KeysListResult struct {
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
	KeysInfo KeysRet        `json:"keys_info"`
	KeyChart KeyInfoChart   `json:"key_chart"`
	Rets     []EcosyKeyList `json:"rets"`
}

type KeyLocks struct {
	NftMinerStake       string `json:"nft_miner_stake"`
	CandidateReferendum string `json:"candidate_referendum"`
	CandidateSubstitute string `json:"candidate_substitute"`
}

// SetTablePrefix is setting table prefix
func (m *Key) SetTablePrefix(prefix int64) *Key {
	m.Ecosystem = prefix
	return m
}

// TableName returns name of table
func (m Key) TableName() string {
	if m.Ecosystem == 0 {
		m.Ecosystem = 1
	}
	return `1_keys`
}

func (m *Key) Get(id int64, wallet string) (*EcosyKeyHex, error) {

	var (
		ecosystems []Ecosystem
	)
	da := EcosyKeyHex{}

	key := strconv.FormatInt(id, 10)
	wid, err := strconv.ParseInt(wallet, 10, 64)
	if err == nil {
		//
		err := conf.GetDbConn().Conn().Table("1_ecosystems").Where("id = ?", key).Find(&ecosystems).Error
		if err == nil {
			da.Ecosystem = id
			da.Ecosystemname = ecosystems[0].Name
			da.IsValued = ecosystems[0].IsValued
			if da.Ecosystem == 1 {
				da.TokenSymbol = SysTokenSymbol
			} else {
				da.TokenSymbol = ecosystems[0].TokenSymbol
			}
			//da.Token_title = ecosystems[0].Token_title
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Where("id = ?", wid).Find(m).Error
		if err == nil {
			da.Ecosystem = id
			da.Amount = m.Amount.String()
		}
	}

	//err := conf.GetDbConn().Conn().Where("id = ? and ecosystem = ?", wallet, m.ecosystem).First(m).Error
	return &da, err
}

func (ts *Key) GetKeys(id int64, page int, size int, order string) (*[]KeyHex, int64, error) {
	var (
		tss    []Key
		ret    []KeyHex
		num    int64
		ioffet int
		i      int
	)

	if order == "" {
		order = "id asc"
	}
	num = 0

	key := strconv.FormatInt(id, 10)
	//err := conf.GetDbConn().Conn().Table(key + "_keys").Order(order).Find(&tss).Error
	err := conf.GetDbConn().Conn().Table("1_keys").Where("ecosystem = ?", key).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
	if page < 1 || size < 1 {
		return &ret, num, err
	}
	ioffet = (page - 1) * size
	num = int64(len(tss))
	if num < int64(page*size) {
		size = int(num) % size
	}
	if num < int64(ioffet) || num < 1 {
		return &ret, num, err
	}

	tokenSymbol, name := GetEcosystemTokenSymbol(id)
	if tokenSymbol == "" && name == "" {
		return &ret, num, err
	}
	for i = 0; i < size; i++ {
		da := KeyHex{}
		da.Ecosystem = id
		if da.Ecosystem == 1 {
			da.TokenSymbol = SysTokenSymbol
		} else {
			da.TokenSymbol = tokenSymbol
		}
		//da.Token_title = es.Token_title
		da.Ecosystemname = name
		da.ID = strconv.FormatInt(tss[ioffet].ID, 10)
		da.PublicKey = hex.EncodeToString(tss[ioffet].PublicKey)
		da.Maxpay = tss[ioffet].Maxpay
		da.Amount = tss[ioffet].Amount.String()
		da.Deleted = tss[ioffet].Deleted
		da.Multi = tss[ioffet].Multi
		da.Blocked = tss[ioffet].Blocked
		//fmt.Println("ecosystem %d", id)
		ret = append(ret, da)
		ioffet++
	}

	return &ret, num, err
}

func (m *Key) GetTotal(page, limit int, order, wallet string) (int64, int, *[]EcosyKeyTotalHex, error) {

	var (
		tss   []Key
		total int64
	)
	var da []EcosyKeyTotalHex

	wid := converter.StringToAddress(wallet)
	err := errors.New("wallet err ")
	//wid, err := strconv.ParseInt(wallet, 10, 64)
	if wid != 0 {
		err = conf.GetDbConn().Conn().Table("1_keys").
			Where("id = ?", wid).
			Count(&total).Error
		if err != nil {
			return 0, 0, &da, err
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Where("id = ?", wid).Order(order).Offset((page - 1) * limit).Limit(limit).Find(&tss).Error
		if err == nil {
			dlen := len(tss)
			for i := 0; i < dlen; i++ {
				ds := tss[i]
				d := EcosyKeyTotalHex{}
				d.Ecosystem = ds.Ecosystem
				d.Amount = ds.Amount.String()

				//
				ems := Ecosystem{}
				f, err := ems.Get(ds.Ecosystem)
				if err != nil {
					return 0, 0, &da, err
				}
				if f {
					d.Ecosystemname = ems.Name
					d.IsValued = ems.IsValued
					if d.Ecosystem == 1 {
						d.TokenSymbol = SysTokenSymbol
					} else {
						d.TokenSymbol = ems.TokenSymbol
					}
				}
				//
				ts := &History{}
				dh, err := ts.GetAccountHistoryTotals(ds.Ecosystem, wid)
				if err != nil {
					return 0, 0, &da, err
				}
				d.Transaction = dh.Transaction
				d.Inamount = dh.Inamount
				d.Outamount = dh.Outamount
				d.InTx = dh.InTx
				d.OutTx = dh.OutTx

				da = append(da, d)
			}

			return total, limit, &da, nil

		}
	}

	return 0, 0, &da, err
}

func (m *Key) GetEcosyKey(keyid int64, wallet string) (*EcosyKeyTotalHex, error) {
	d := EcosyKeyTotalHex{}
	d.Ecosystem = m.Ecosystem
	d.Amount = m.Amount.String()
	d.Wallet = m.AccountID

	mb := Member{}
	fm, _ := mb.GetAccount(m.Ecosystem, wallet)
	if fm {
		if mb.ImageID != nil {
			if *mb.ImageID != int64(0) {
				hash, err := GetFileHash(*mb.ImageID)
				if err != nil {
					return &d, err
				}
				d.MemberHash = hash
			}
		}
		if mb.MemberName != "" {
			d.MemberName = mb.MemberName
		} else {
			d.MemberName = "iName"
		}
	}

	escape := func(value any) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}

	//
	ems := Ecosystem{}
	f, err := ems.Get(m.Ecosystem)
	if err != nil {
		return &d, err
	}
	if f {
		d.Ecosystemname = ems.Name
		d.IsValued = ems.IsValued

		if ems.Info != "" {
			minfo := make(map[string]any)
			err := json.Unmarshal([]byte(ems.Info), &minfo)
			if err != nil {
				return &d, err
			}
			usid, ok := minfo["logo"]
			if ok {
				urid := escape(usid)
				uid, err := strconv.ParseInt(urid, 10, 64)
				if err != nil {
					return &d, err
				}

				hash, err := GetFileHash(uid)
				if err != nil {
					return &d, err
				}
				d.LogoHash = hash

			}
		}
		if d.Ecosystem == 1 {
			d.TokenSymbol = SysTokenSymbol
			if d.Ecosystemname == "" {
				d.Ecosystemname = SysEcosystemName
			}
		} else {
			d.TokenSymbol = ems.TokenSymbol
		}
	}
	//
	ts := &History{}
	dh, err := ts.GetAccountHistoryTotals(m.Ecosystem, keyid)
	if err != nil {
		return &d, err
	}

	ag := &AssignGetInfo{}
	ba, fa, _, err := ag.GetBalance(nil, keyid)
	if err != nil {
		return &d, err
	}
	if m.Lock != "" {
		var stake KeyLocks
		if err := json.Unmarshal([]byte(m.Lock), &stake); err != nil {
			return &d, err
		}
		nftLock, _ := decimal.NewFromString(stake.NftMinerStake)
		referendumLock, _ := decimal.NewFromString(stake.CandidateReferendum)
		substituteLock, _ := decimal.NewFromString(stake.CandidateSubstitute)

		d.StakeAmount = nftLock.Add(referendumLock).Add(substituteLock)
	}

	if m.Ecosystem == 1 {
		chart, err := ts.GetWalletTimeLineHistoryTotals(m.Ecosystem, keyid)
		if err != nil {
			return &d, err
		}
		d.Chart = chart
	}
	d.Transaction = dh.Transaction
	d.InTx = dh.InTx
	d.OutTx = dh.OutTx
	d.Inamount = dh.Inamount
	d.Outamount = dh.Outamount
	if ba {
		d.FreezeAmount = fa
	}
	amount := decimal.New(0, 0)
	if len(d.Amount) > 0 {
		amount, _ = decimal.NewFromString(d.Amount)
	}
	d.Totalamount = d.Totalamount.Add(amount)
	d.Totalamount = d.Totalamount.Add(d.StakeAmount)
	//d.Totalamount = d.Totalamount.Add(d.FreezeAmount)
	d.JoinTime = getJoinTime(keyid, m.Ecosystem)
	d.RolesName = getRolesName(converter.AddressToString(keyid), m.Ecosystem)
	return &d, err
}

func getJoinTime(keyId int64, ecosystem int64) int64 {
	var tableId string
	ecosystemStr := strconv.FormatInt(ecosystem, 10)
	keyIdStr := strconv.FormatInt(keyId, 10)
	tableId = strings.Join([]string{keyIdStr, ecosystemStr}, ",")
	var rollback sqldb.RollbackTx
	f, err := isFound(GetDB(nil).Select("block_id").Table("rollback_tx").Where("table_name = '1_keys' AND table_id = ? AND data = ''", tableId).First(&rollback))
	if err != nil {
		log.Info("get join time err:", err.Error(), " table_id:", tableId)
		return 0
	}
	var bk Block
	if f {
		f, err := isFound(GetDB(nil).Select("time").Where("id = ?", rollback.BlockID).First(&bk))
		if err != nil {
			log.Info("get join time blockId err:", err.Error(), " blockId:", rollback.BlockID)
			return 0
		}

		if f {
			return bk.Time
		}
	} else {
		return FirstBlockTime
	}
	return 0
}

func getRolesName(account string, ecosystem int64) string {
	var roles sqldb.RolesParticipants
	var rolesName string
	var nameList []string
	roles.SetTablePrefix(ecosystem)
	list, err := roles.GetActiveMemberRoles(account)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "key": account}).Warn("GetActiveMemberRoles err")
		return ""
	}
	if len(list) == 0 {
		return ""
	}
	type roleData struct {
		Name string `json:"name"`
	}
	var role roleData
	for i := 0; i < len(list); i++ {
		if err := json.Unmarshal([]byte(list[i].Role), &role); err != nil {
			log.WithFields(log.Fields{"warn": err, "role": list[i].Role}).Warn("getRolesName json err")
			return ""
		}
		nameList = append(nameList, role.Name)
	}
	if len(nameList) > 0 {
		rolesName = strings.Join(nameList, " / ")
	}
	return rolesName
}

func (m *Key) GetWalletTotalBasisEcosystem(wallet string) (*EcosyKeyTotalHex, error) {
	var (
		ft  Key
		ret EcosyKeyTotalHex
	)
	wid := converter.StringToAddress(wallet)
	if wid != 0 || wallet == "0000-0000-0000-0000-0000" {
		f, err := isFound(GetDB(nil).Table("1_keys").Where("id = ? and ecosystem = ?", wid, 1).First(&ft))
		if err != nil {
			return &ret, err
		}
		if !f {
			if wallet != "0000-0000-0000-0000-0000" {
				return nil, err
			} else {
				ft.Ecosystem = 1
				ft.AccountID = "0000-0000-0000-0000-0000"
			}
		}
		df, err := ft.GetEcosyKey(wid, wallet)
		if err != nil {
			return &ret, err
		}
		ret = *df
	} else {
		return &ret, errors.New("wallet invalid")
	}
	return &ret, nil
}

//GetWalletTotalEcosystem response 10-20ms The modified 4-9ms
func (m *Key) GetWalletTotalEcosystem(page, limit int, order string, wallet string) (*GeneralResponse, error) {
	var (
		total int64
		ret   GeneralResponse
	)
	ret.Limit = limit
	ret.Page = page
	da := []EcosyKeyTotalDetail{}

	if page <= 0 || limit <= 0 {
		return nil, errors.New("request params invalid")
	}
	if order == "" {
		order = "total_amount desc"
	}

	wid := converter.StringToAddress(wallet)
	if wid != 0 || wallet == "0000-0000-0000-0000-0000" {
		err := conf.GetDbConn().Conn().Table("1_keys").
			Where("id = ?", wid).
			Count(&total).Error
		if err != nil {
			return &ret, err
		}
		if NftMinerReady {
			err = GetDB(nil).Table(`"1_keys" AS k1`).Select(`account,ecosystem,
	(SELECT hash AS logo_hash FROM "1_binaries" as bs WHERE bs.id = coalesce((SELECT cast(es.info->>'logo' as numeric) FROM "1_ecosystems" as es WHERE es.id = k1.ecosystem LIMIT 1),0)),
	(SELECT name AS name FROM "1_ecosystems" as es WHERE es.id = k1.ecosystem),
	(SELECT array_to_string(array(SELECT rs.role->>'name' FROM "1_roles_participants" as rs 
	WHERE rs.ecosystem=k1.ecosystem and rs.member->>'account' = k1.account AND rs.deleted = 0),' / '))as roles_name,
	(SELECT time from block_chain b WHERE b.id = coalesce((SELECT rt.block_id FROM rollback_tx rt WHERE table_name = '1_keys' AND table_id = cast(k1.id as VARCHAR) 
	|| ','||  cast(k1.ecosystem as VARCHAR) AND data = '' LIMIT 1),1)) AS join_time,
	case WHEN k1.ecosystem = 1 THEN
	 'IBXC'
	ELSE
	 (SELECT token_symbol FROM "1_ecosystems" as es WHERE es.id = k1.ecosystem)
	END as token_symbol,
	k1.amount +(SELECT COALESCE(sum(output_value),0) FROM "spent_info" WHERE input_tx_hash is null AND output_key_id = k1.id)+ 
		to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')as total_amount,
	to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999') AS stake_amount`).
				Where("id = ?", wid).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&da).Error
		} else {
			err = GetDB(nil).Table(`"1_keys" AS k1`).Select(`account,ecosystem,
(SELECT hash AS logo_hash FROM "1_binaries" as bs WHERE bs.id = coalesce((SELECT cast(es.info->>'logo' as numeric) FROM "1_ecosystems" as es WHERE es.id = k1.ecosystem LIMIT 1),0)),
(SELECT name AS name FROM "1_ecosystems" as es WHERE es.id = k1.ecosystem),
(SELECT array_to_string(array(SELECT rs.role->>'name' FROM "1_roles_participants" as rs 
WHERE rs.ecosystem=k1.ecosystem and rs.member->>'account' = k1.account AND rs.deleted = 0),' / '))as roles_name,
(SELECT time from block_chain b WHERE b.id = coalesce((SELECT rt.block_id FROM rollback_tx rt WHERE table_name = '1_keys' AND table_id = cast(k1.id as VARCHAR) 
|| ','||  cast(k1.ecosystem as VARCHAR) AND data = '' LIMIT 1),1)) AS join_time,
case WHEN k1.ecosystem = 1 THEN
 'IBXC'
ELSE
 (SELECT token_symbol FROM "1_ecosystems" as es WHERE es.id = k1.ecosystem)
END as token_symbol,
k1.amount+(SELECT COALESCE(sum(output_value),0) FROM "spent_info" WHERE input_tx_hash is null AND output_key_id = k1.id) as total_amount`).
				Where("id = ?", wid).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&da).Error
		}
		if err != nil {
			return &ret, err
		}

		ret.Total = int64(total)
		ret.List = da
		return &ret, nil
	} else {
		return &ret, errors.New("wallet invalid")
	}
}

// GetKeysCount returns common count of keys
func GetTotalAmount(ecosystem int64) (decimal.Decimal, error) {
	var err error
	type result struct {
		Amount decimal.Decimal
	}
	var (
		res  result
		utxo SumAmount
	)
	err = GetDB(nil).Table("1_keys").
		Select("coalesce(sum(amount),0) as amount").Where("ecosystem = ?", ecosystem).Scan(&res).Error
	if err != nil {
		return decimal.Zero, err
	}
	err = GetDB(nil).Table("spent_info").Select("sum(output_value)").Where("input_tx_hash is NULL AND ecosystem = ?", ecosystem).Take(&utxo).Error
	if err != nil {
		return decimal.Zero, err
	}
	return res.Amount.Add(utxo.Sum), err
}

func (m *Key) GetStakeAmount() (string, error) {
	type result struct {
		Amount decimal.Decimal
	}
	var agi AssignGetInfo
	agm, err := agi.GetAllBalance(nil)
	if err != nil {
		return "0", err
	}

	if HasTableOrView(nil, "1_mine_stake") {
		var mine, pool result
		err = conf.GetDbConn().Conn().Table("1_keys").Select("SUM(mine_lock) as amount").Where("ecosystem = 1").Scan(&mine).Error
		if err != nil {
			b := strings.ContainsAny(err.Error(), "column mine_lock does not exist")
			if b {
				return agm.String(), nil
			}
			return "0", err
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Select("SUM(pool_lock) as amount").Where("ecosystem = 1").Scan(&pool).Error
		if err != nil {
			b := strings.ContainsAny(err.Error(), "column pool_lock does not exist")
			if b {
				return agm.String(), nil
			}
			return "0", err
		}

		rt := mine.Amount.Add(pool.Amount)
		rt = rt.Add(agm)
		return rt.String(), nil
	}
	return agm.String(), nil
}

//GetAccountList response 40ms-90ms The modified:40ms-60ms todo:need optimize newkeys chart and FreezeAmount
func (m *Key) GetAccountList(page, limit int, reqOrder string, ecosystem int64) (*KeysListResult, error) {

	var (
		total int64
		ret   KeysListResult
		order string
	)
	if reqOrder == "" {
		order = "amount desc"
	} else {
		order = reqOrder
	}
	ret.Limit = limit
	ret.Page = page

	err := GetDB(nil).Table("1_keys").
		Where("ecosystem = ?", ecosystem).
		Count(&total).Error
	if err != nil {
		return &ret, err
	}
	ret.KeysInfo = getScanOutKeyInfo(ecosystem)
	ret.KeyChart = GetEcosystemNewKeyChart(ecosystem, 15)

	if NftMinerReady || NodeReady {
		err = GetDB(nil).Table(`"1_keys" as k1`).Select(`
	account,ecosystem,
	coalesce((SELECT ms.member_name FROM "1_members" as ms WHERE ms.ecosystem = k1.ecosystem AND ms.account = k1.account LIMIT 1),'iName') as account_name,
	k1.amount +  to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999') as amount,
		
	to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999') AS stake_amount,
	case WHEN k1.ecosystem = 1 THEN
	 coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem),'IBXC')
	ELSE
	 (SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem) 
	END as token_symbol,
	
	CASE WHEN (k1.amount +  to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999') = 0) OR 
		((SELECT sum(k2.amount)+sum(to_number(coalesce(NULLIF(k2.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem) = 0) THEN
			0
	ELSE
		round(
		(k1.amount +  to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')) * 100 / 
			(SELECT sum(k2.amount)+sum(to_number(coalesce(NULLIF(k2.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem), 2) 
	END as accounted_for
`).Where("ecosystem = ?", ecosystem).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&ret.Rets).Error
	} else {
		err = GetDB(nil).Table(`"1_keys" as k1`).Select(`account,ecosystem,
coalesce((SELECT ms.member_name FROM "1_members" as ms WHERE ms.ecosystem = k1.ecosystem AND ms.account = k1.account LIMIT 1),'iName') as account_name,
k1.amount as amount,
case WHEN k1.ecosystem = 1 THEN
 coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem),'IBXC')
ELSE
 (SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem) 
END as token_symbol,

case WHEN (k1.amount = 0) OR 
((SELECT sum(k2.amount) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem) = 0) THEN
0
ELSE
round(
k1.amount * 100 / 
  (SELECT sum(k2.amount) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem) , 2) 
	END as accounted_for`).Where("ecosystem = ?", ecosystem).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&ret.Rets).Error
	}
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ret.Rets); i++ {
		if ecosystem == 1 {
			var ag AssignGetInfo
			ba, fa, _, err := ag.GetBalance(nil, converter.StringToAddress(ret.Rets[i].Account))
			if err != nil {
				return &ret, err
			}
			if ba {
				ret.Rets[i].FreezeAmount = fa.String()
			}
		}
	}
	ret.Total = total

	return &ret, nil
}

func (m *Key) GetEcosystemTokenSymbolList(page, limit int, order string, ecosystem int64) (*GeneralResponse, error) {

	var (
		total int64
		rets  GeneralResponse
		err   error
	)
	if order == "" {
		order = "amount desc"
	}
	rets.Limit = limit
	rets.Page = page

	if NftMinerReady {
		err = GetDB(nil).Table(m.TableName()).
			Where("ecosystem = ? AND (amount > 0 OR lock->>'nft_miner_stake' != '')", ecosystem).
			Count(&total).Error
	} else {
		err = GetDB(nil).Table(m.TableName()).
			Where("ecosystem = ? AND amount > 0", ecosystem).
			Count(&total).Error
	}
	if err != nil {
		return nil, err
	}
	rets.Total = total

	var ret []EcosystemTokenSymbolList
	if NftMinerReady || NodeReady {
		err = GetDB(nil).Table(`"1_keys" AS k1`).Select(fmt.Sprintf(
			`k1,id,k1.account, 
 coalesce((SELECT ms.member_name FROM "1_members" as ms WHERE ms.ecosystem = k1.ecosystem AND ms.account = k1.account),'iName') as account_name,
 k1.amount +  to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999') as amount,
 
case WHEN (k1.amount +  to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999') = 0) OR 
((SELECT sum(k2.amount)+sum(to_number(coalesce(NULLIF(k2.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem) = 0) THEN
0
ELSE
round(
(k1.amount +  to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')) * 100 / 
  (SELECT sum(k2.amount)+sum(to_number(coalesce(NULLIF(k2.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k2.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem) , 2) 
	END as accounted_for,

case WHEN k1.ecosystem = 1 THEN
 coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem),'%s')
ELSE
 (SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem) 
END as token_symbol,
CASE WHEN (SELECT control_mode FROM "1_ecosystems" WHERE id = k1.ecosystem AND id > 1) = 2 THEN
		CASE WHEN (SELECT count(1) FROM "1_votings_participants" WHERE 
				voting_id = (SELECT id FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%%voting_for_control_mode_template%%' AND ecosystem = k1.ecosystem ORDER BY id DESC LIMIT 1) 
				AND member->>'account'=k1.account) > 0 THEN
			TRUE
		ELSE
			FALSE
		END
ELSE
	FALSE
END AS front_committee
`, SysTokenSymbol)).
			Where("ecosystem = ? AND (amount > 0 OR lock->>'nft_miner_stake' != '' OR lock->>'candidate_substitute' != '' OR lock->>'candidate_referendum' != '')", ecosystem).
			Order(order).Offset((page - 1) * limit).Limit(limit).
			Find(&ret).Error
	} else {
		err = GetDB(nil).Table(`"1_keys" AS k1`).Select(fmt.Sprintf(
			`k1,id,k1.account, 
 coalesce((SELECT ms.member_name FROM "1_members" as ms WHERE ms.ecosystem = k1.ecosystem AND ms.account = k1.account),'iName') as account_name,
 k1.amount as amount,
CASE WHEN (SELECT control_mode FROM "1_ecosystems" WHERE id = k1.ecosystem AND id > 1) = 2 THEN
		CASE WHEN (SELECT count(1) FROM "1_votings_participants" WHERE 
				voting_id = (SELECT id FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%%voting_for_control_mode_template%%' AND ecosystem = k1.ecosystem ORDER BY id DESC LIMIT 1) 
				AND member->>'account'=k1.account) > 0 THEN
			TRUE
		ELSE
			FALSE
		END
ELSE
	FALSE
END AS front_committee,
case WHEN (k1.amount = 0) OR 
((SELECT sum(k2.amount) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem) = 0) THEN
0
ELSE
round(
k1.amount  * 100 / 
  (SELECT sum(k2.amount) FROM "1_keys" AS k2 WHERE k1.ecosystem = k2.ecosystem) , 2) 
	END as accounted_for,

case WHEN k1.ecosystem = 1 THEN
 coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem),'%s')
ELSE
 (SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem) 
END as token_symbol`, SysTokenSymbol)).
			Where("ecosystem = ? AND amount > 0", ecosystem).
			Order(order).Offset((page - 1) * limit).Limit(limit).
			Find(&ret).Error
	}
	if err != nil {
		return nil, err
	}
	type ids struct {
		Id int64
	}
	var committeeList []ids
	if ecosystem > 1 {
		err = GetDB(nil).Raw(`
SELECT id FROM "1_keys" AS k1 WHERE (SELECT control_mode FROM "1_ecosystems" WHERE id = k1.ecosystem AND id > 1) = 2 AND
ecosystem = ? AND deleted = 0 AND blocked = 0 AND amount > 0 ORDER BY row_number() OVER (ORDER BY amount DESC) <= 50
`, ecosystem).Find(&committeeList).Error
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < len(ret); i++ {
		for _, v := range committeeList {
			if v.Id == ret[i].Id {
				ret[i].Committee = true
			}
		}
	}
	rets.List = ret

	return &rets, err
}

func (m *Key) GetEcosystemDetailMemberList(page, limit int, order string, ecosystem int64) (*GeneralResponse, error) {
	var (
		total int64
		rets  GeneralResponse
	)
	if order == "" {
		order = "join_time desc"
	}
	rets.Limit = limit
	rets.Page = page

	err := GetDB(nil).Table(m.TableName()).
		Where("ecosystem = ?", ecosystem).
		Count(&total).Error
	if err != nil {
		return nil, err
	}
	rets.Total = total

	var ret []EcosystemMemberList
	err = GetDB(nil).Table(`"1_keys" as k1`).Select(
		`account,
(SELECT array_to_string(array(SELECT rs.role->>'name' FROM "1_roles_participants" as rs 
WHERE rs.ecosystem=k1.ecosystem and rs.member->>'account' = k1.account AND rs.deleted = 0),' / '))
as roles_name,
coalesce((SELECT ms.member_name FROM "1_members" as ms WHERE ms.ecosystem = k1.ecosystem AND ms.account = k1.account LIMIT 1),'iName') as account_name ,

(SELECT time from block_chain b WHERE b.id = coalesce((SELECT rt.block_id FROM rollback_tx rt WHERE table_name = '1_keys' AND table_id = cast(k1.id as VARCHAR) 
|| ','||  cast(k1.ecosystem as VARCHAR) AND data = '' LIMIT 1),1)) AS join_time`).
		Where("ecosystem = ?", ecosystem).
		Order(order).Offset((page - 1) * limit).Limit(limit).
		Find(&ret).Error
	rets.List = ret

	return &rets, err
}

func GetEcosystemNewKeyChart(ecosystem int64, getDay int) KeyInfoChart {
	var rets KeyInfoChart
	var bk Block
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	t1 := yesterday.AddDate(0, 0, -1*getDay)

	rets.Time = make([]int64, getDay)
	rets.NewKey = make([]int64, getDay)
	idList := make([]int64, getDay)

	var daysId []DaysNumber
	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		if i == 0 {
			err := GetDB(nil).Raw(fmt.Sprintf(`SELECT to_char(to_timestamp(time),'yyyy-MM-dd') AS days,min(id) num 
FROM block_chain WHERE time >= %d and time < %d GROUP BY days ORDER BY days asc
`, rets.Time[i], today.Unix())).Find(&daysId).Error
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Warn("Get ecosystem New Key Chart Days Id List Failed")
				return rets
			}
		}
		id := GetDaysNumberLike(rets.Time[i], daysId, false, "asc")
		if id == 0 {
			fmt.Println("get days number like failed!")
		} else {
			idList[i] = id
		}
		_, err := bk.GetByTimeBlockId(today.Unix())
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("get Ecosystem New Key Chart Today Block Failed")
			return rets
		}
		//if !f {
		//	return rets
		//}
	}

	for i := 0; i < len(idList); i++ {
		if i == len(idList)-1 {
			rets.NewKey[i] = getNewKeyNumber(idList[i], bk.ID, ecosystem)
		} else {
			rets.NewKey[i] = getNewKeyNumber(idList[i], idList[i+1], ecosystem)
		}
	}
	_, rets.Name = GetEcosystemTokenSymbol(ecosystem)

	return rets
}

func GetAccountTotalAmount(ecosystem int64, account string) (AccountTotalAmountChart, error) {
	var (
		rets AccountTotalAmountChart
		err  error
		f    bool
	)
	keyId := converter.StringToAddress(account)
	if keyId == 0 {
		return rets, errors.New("account invalid:" + account + " in ecosystem:" + strconv.FormatInt(ecosystem, 10))
	}
	if NftMinerReady {
		f, err = isFound(GetDB(nil).Table(`"1_keys" as k1`).Select(`k1.amount AS amount,
to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999')+ 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999') + 
		to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999')as stake_amount,
case WHEN k1.ecosystem = 1 THEN
 coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem),'IBXC')
ELSE
 (SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem) 
END as token_symbol`).Where("ecosystem = ? and id = ?", ecosystem, keyId).Take(&rets))
	} else {
		f, err = isFound(GetDB(nil).Table(`"1_keys" as k1`).Select(`k1.amount AS amount,
case WHEN k1.ecosystem = 1 THEN
 coalesce((SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem),'IBXC')
ELSE
 (SELECT token_symbol FROM "1_ecosystems" as ec WHERE ec.id = k1.ecosystem) 
END as token_symbol`).Where("ecosystem = ? and id = ?", ecosystem, keyId).Take(&rets))
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Account Total Amount Chart Failed")
		return rets, nil
	}
	if !f {
		return rets, errors.New("unknown account:" + account + " in ecosystem:" + strconv.FormatInt(ecosystem, 10))
	}
	rets.TotalAmount = rets.Amount.Add(rets.FreezeAmount).Add(rets.StakeAmount)
	zeroDec := decimal.New(0, 0)
	if rets.TotalAmount.GreaterThan(zeroDec) {
		if rets.Amount.GreaterThan(zeroDec) {
			rets.AmountRatio, _ = rets.Amount.Mul(decimal.NewFromInt(100)).DivRound(rets.TotalAmount, 2).Float64()
		}
		if rets.StakeAmount.GreaterThan(zeroDec) {
			rets.StakeAmountRatio, _ = rets.StakeAmount.Mul(decimal.NewFromInt(100)).DivRound(rets.TotalAmount, 2).Float64()
		}
		if rets.FreezeAmount.GreaterThan(zeroDec) {
			rets.FreezeAmountRatio, _ = rets.FreezeAmount.Mul(decimal.NewFromInt(100)).DivRound(rets.TotalAmount, 2).Float64()
		}
	}
	return rets, nil
}
