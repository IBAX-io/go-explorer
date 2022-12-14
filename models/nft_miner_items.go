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
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/smart"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type NftMinerItems struct {
	ID          int64  `gorm:"primary_key;not null"`         //NFT ID
	EnergyPoint int    `gorm:"column:energy_point;not null"` //power
	Owner       string `gorm:"column:owner;not null"`        //owner account
	Creator     string `gorm:"column:creator;not null"`      //create account
	MergeCount  int64  `gorm:"column:merge_count;not null"`  //merage count
	MergeStatus int64  `gorm:"column:merge_status;not null"` //1:un merge(valid) 0:merge(invalid)
	Attributes  string `gorm:"column:attributes;not null"`   //SVG Params
	TempId      int64  `gorm:"column:temp_id;not null"`
	DateCreated int64  `gorm:"column:date_created;not null"` //create time
	Coordinates string `gorm:"column:coordinates;type:jsonb"`
	TokenHash   []byte `gorm:"column:token_hash"`
	TxHash      []byte `gorm:"column:tx_hash"`
}

type SvgParams struct {
	Owner       string `json:"owner"`
	Point       string `json:"point"`
	Number      int64  `json:"number"`
	Location    string `json:"location"`
	DateCreated int64  `json:"date_created"` //milliseconds
}

type StrReplaceStruct struct {
	CapitalLetter    int `json:"capital_letter"`
	LowercaseLetters int `json:"lowercase_letters"`
	Number           int `json:"number"`
	OtherString      int `json:"other_string"`
}

var NftMinerReady bool

func (p *NftMinerItems) TableName() string {
	return "1_nft_miner_items"
}

func (p *NftMinerItems) GetById(id int64) (bool, error) {
	return isFound(GetDB(nil).Where("id = ? ", id).First(p))
}

func (p *NftMinerItems) GetByTokenHash(tokenHash string) (bool, error) {
	hash, _ := hex.DecodeString(tokenHash)
	return isFound(GetDB(nil).Where("token_hash = ? ", hash).First(p))
}

func (p *NftMinerItems) GetAllPower() (int64, error) {
	type powerCount struct {
		Count int64 `gorm:"column:count"`
	}
	var count powerCount
	if HasTableOrView(p.TableName()) {
		if err := GetDB(nil).Table(p.TableName()).Select("sum(energy_point) count").Where("merge_status = 1").First(&count).Error; err != nil {
			return 0, err
		}
		return count.Count, nil
	} else {
		return 0, nil
	}
}

func (p *NftMinerItems) ParseSvgParams() (string, error) {
	var (
		ret SvgParams
		app AppParam
	)
	err := json.Unmarshal([]byte(p.Attributes), &ret)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Parse Svg Params json unmarshal failed")
		return "", err
	}
	f, err := app.GetById(nil, p.TempId)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "tempId": p.TempId}).Error("Parse Svg Params get app id failed")
		return "", err
	}
	if !f {
		return "", nil
	}
	star := "★"
	if p.EnergyPoint <= 20 {
	} else if p.EnergyPoint <= 40 {
		star = "★★"
	} else if p.EnergyPoint <= 60 {
		star = "★★★"
	} else if p.EnergyPoint <= 80 {
		star = "★★★★"
	} else {
		star = "★★★★★"
	}
	return fmt.Sprintf(app.Value, ret.Point, "#"+strconv.FormatInt(ret.Number, 10), formatLocation(ret.Location), smart.Date(SvgTimeFormat, MsToSeconds(ret.DateCreated)), ret.Owner, star), nil
}

// formatLocation example:NorthAmerica result:North America
func formatLocation(location string) string {
	strs, indexList := StrReplaceAllString(location)
	if strs.CapitalLetter >= 2 {
		for i := 2; i <= strs.CapitalLetter; i++ {
			location = location[:indexList[i-1].CapitalLetter] + " " + location[indexList[i-1].CapitalLetter:]
		}
	}

	return location
}

func StrReplaceAllString(s2 string) (strReplace StrReplaceStruct, indexList []StrReplaceStruct) {
	indexList = make([]StrReplaceStruct, len(s2))
	for i := 0; i < len(s2); i++ {
		switch {
		case 64 < s2[i] && s2[i] < 91:
			strReplace.CapitalLetter += 1
			indexList[strReplace.CapitalLetter-1].CapitalLetter = i
		case 96 < s2[i] && s2[i] < 123:
			strReplace.LowercaseLetters += 1
			indexList[strReplace.LowercaseLetters-1].LowercaseLetters = i
		case 47 < s2[i] && s2[i] < 58:
			strReplace.Number += 1
			indexList[strReplace.Number-1].Number = i
		default:
			strReplace.OtherString += 1
			indexList[strReplace.OtherString-1].OtherString = i
		}
	}
	return strReplace, indexList
}

func (p *NftMinerItems) GetAccountDetailNftMinerInfo(keyid, order string, page, limit int) (*AccountNftMinerListResult, error) {
	if order == "" {
		order = "id desc"
	}
	q := GetDB(nil).Table(p.TableName()).Where("owner = ? AND merge_status = 1", keyid)
	var total int64
	var res AccountNftMinerListResult
	res.Page = page
	res.Limit = limit

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []NftMinerItems
	err := q.Order(order).Limit(limit).Offset((page - 1) * limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	rets := make([]AccountNftMinerList, len(list))
	nowTime := time.Now().Unix()
	for i := 0; i < len(list); i++ {
		var stak NftMinerStaking
		var stakeAmount int64
		var cycle int64
		f, err := stak.GetByTokenId(list[i].ID)
		if err != nil {
			log.Info("get wallet nft miner key by id err:", err.Error(), " id:", list[i].ID)
		}
		if f {
			if nowTime >= stak.StartDated && nowTime <= stak.EndDated {
				stakeAmount = stak.StakeAmount
				cycle = int64(time.Unix(stak.EndDated, 0).Sub(time.Unix(stak.StartDated, 0)).Hours() / 24)
			}
		}
		var nftIns SumAmount
		var his History
		kid := converter.StringToAddress(keyid)
		_, err = isFound(GetDB(nil).Table(his.TableName()).Select("sum(amount)").
			Where("recipient_id = ? AND type = 12 AND comment = ?", kid, fmt.Sprintf("NFT Miner #%d", list[i].ID)).Take(&nftIns))
		if err != nil {
			log.Info("get wallet nft miner key nftIns err:", err.Error())
		}

		rets[i].ID = list[i].ID
		rets[i].Hash = hex.EncodeToString(list[i].TokenHash)
		rets[i].EnergyPoint = list[i].EnergyPoint
		rets[i].Time = list[i].DateCreated
		rets[i].StakeAmount = stakeAmount
		rets[i].Cycle = cycle
		rets[i].Ins = nftIns.Sum.String()
	}

	res.Rets = rets
	res.Total = total
	return &res, nil
}

func (p *NftMinerItems) GetNftMinerBySearch(search string) (NftMinerInfoResponse, error) {
	var (
		rets NftMinerInfoResponse
		f    bool
		err  error
	)

	nftId := converter.StrToInt64(search)
	if nftId > 0 {
		f, err = p.GetById(nftId)
	} else {
		f, err = p.GetByTokenHash(search)
	}
	if err != nil {
		return rets, err
	}
	if !f {
		return rets, errors.New("NFT doesn't not exist")
	}

	rets.ID = p.ID
	rets.Hash = hex.EncodeToString(p.TokenHash)
	rets.EnergyPoint = p.EnergyPoint
	rets.DateCreated = p.DateCreated
	rets.Creator = p.Creator
	rets.Owner = p.Owner

	var stak NftMinerStaking
	err = GetDB(nil).Table(stak.TableName()).Where("token_id = ?", p.ID).Count(&rets.StakeCount).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft info stakeCount err:", err.Error(), " nftId:", p.ID)
		}
		return rets, err
	}

	nowTime := time.Now().Unix()
	f, err = isFound(GetDB(nil).Where("staking_status = 1 AND token_id = ?", p.ID).Last(&stak))
	if err != nil {
		log.Info("get nft info stakeAmount err:", err.Error(), " nftId:", p.ID)
		return rets, err
	}
	if f {
		rets.StakeAmount = stak.StakeAmount
		rets.Cycle = int64(time.Unix(stak.EndDated, 0).Sub(time.Unix(stak.StartDated, 0)).Hours() / 24)
		if stak.StartDated <= nowTime && stak.EndDated >= nowTime {
			rets.EnergyPower = stak.EnergyPower
		}
	}

	var reward SumAmount
	var his History
	//kid := converter.StringToAddress(p.Owner)
	q := GetDB(nil).Table(his.TableName()).Where("type = 12 AND comment = ?", fmt.Sprintf("NFT Miner #%d", p.ID))
	err = q.Select("sum(amount)").Take(&reward).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft miner reward info err:", err.Error(), " nftId:", p.ID)
			return rets, err
		}
	}
	err = q.Count(&rets.RewardCount).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft miner info reward Count err:", err.Error(), " nftId:", p.ID)
			return rets, err
		}
	}
	rets.Ins = reward.Sum.String()

	return rets, nil
}

func (p *NftMinerItems) GetNftMinerTxInfo(search any, page, limit int, order string) (*GeneralResponse, error) {
	var (
		rets  []NftMinerTxInfoResponse
		total int64
		ret   GeneralResponse
		f     bool
		err   error
	)
	if order == "" {
		order = "id desc"
	}
	switch reflect.TypeOf(search).String() {
	case "string":
		f, err = p.GetByTokenHash(search.(string))
	case "json.Number":
		tokenId, err := search.(json.Number).Int64()
		if err != nil {
			return nil, err
		}
		f, err = p.GetById(tokenId)
	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Nft Miner Tx Info Search Failed")
		return nil, errors.New("request params invalid")
	}

	if err != nil {
		log.Info("get nft miner txInfo err:", err.Error(), " nftId:", p.ID)
		return nil, err
	}
	if !f {
		return nil, errors.New("NFT Miner Doesn't Not Exist")
	}
	var his []History
	err = GetDB(nil).Table("1_history").Where("type = 12 AND comment = ?", fmt.Sprintf("NFT Miner #%d", p.ID)).Count(&total).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft tx info total err:", err.Error(), " nftId:", p.ID)
			return nil, err
		}
	}

	err = GetDB(nil).Select("id,block_id,created_at,amount,recipient_id").Where("type = 12 AND comment = ?", fmt.Sprintf("NFT Miner #%d", p.ID)).Limit(limit).Offset((page - 1) * limit).Order(order).Find(&his).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft miner tx info nftIns err:", err.Error(), " nft miner id:", p.ID)
			return nil, err
		}
	}
	rets = make([]NftMinerTxInfoResponse, len(his))
	for i := 0; i < len(his); i++ {
		rets[i].ID = his[i].ID
		rets[i].NftId = p.ID
		rets[i].BlockId = his[i].Blockid
		rets[i].Time = MsToSeconds(his[i].Createdat)
		rets[i].Ins = his[i].Amount.String()
		rets[i].Account = converter.AddressToString(his[i].Recipientid)
	}
	ret.Page = page
	ret.Limit = limit
	ret.Total = total
	ret.List = rets

	return &ret, nil
}

func (p *NftMinerItems) NftMinerMetaverseList(page, limit int, order string) (GeneralResponse, error) {
	if order == "" {
		order = "date_created desc"
	}
	var rets []NftMinerListResponse
	var ret GeneralResponse
	ret.Page = page
	ret.Limit = limit
	var list []NftMinerItems
	err := GetDB(nil).Table(p.TableName()).Where("merge_status = 1").Count(&ret.Total).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft miner metaverse list count err:", err.Error())
		}
		return ret, err
	}
	err = GetDB(nil).Where("merge_status = 1").Offset((page - 1) * limit).Limit(limit).Order(order).Find(&list).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Info("get nft miner metaverse list err:", err.Error())
		}
		return ret, err
	}
	rets = make([]NftMinerListResponse, len(list))
	nowTime := time.Now().Unix()

	for i := 0; i < len(list); i++ {
		rets[i].ID = list[i].ID
		rets[i].Hash = hex.EncodeToString(list[i].TokenHash)
		rets[i].Time = list[i].DateCreated
		rets[i].EnergyPoint = list[i].EnergyPoint
		rets[i].Owner = list[i].Owner
		var stak NftMinerStaking
		f, err := stak.GetNftStakeByTokenId(list[i].ID)
		if err != nil {
			log.Info("get nft miner metaverse list err:", err.Error(), " nft miner id:", list[i].ID)
		} else {
			if f {
				if nowTime >= stak.StartDated && nowTime < stak.EndDated {
					rets[i].EnergyPower = stak.EnergyPower
				}
				rets[i].StakeAmount = stak.StakeAmount
			}
		}

		var nftIns SumAmount
		var his History
		kid := converter.StringToAddress(list[i].Owner)
		q := GetDB(nil).Table(his.TableName()).Where("recipient_id = ? AND type = 12 AND comment = ?", kid, fmt.Sprintf("NFT Miner #%d", list[i].ID))
		err = q.Select("sum(amount)").Take(&nftIns).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Info("get nft miner metaverse list nftIns err:", err.Error(), " nft miner id:", p.ID)
			}
		}
		rets[i].Ins = nftIns.Sum.String()
	}
	ret.List = rets

	return ret, nil
}

func (p *NftMinerItems) GetNftMetaverse() (*NftMinerMetaverseInfoResponse, error) {
	var rets NftMinerMetaverseInfoResponse
	var m ScanOut
	f, err := m.GetRedisLatest()
	if err != nil {
		return &rets, err
	}
	if f {
		rets.BlockReward = m.NftBlockReward
		rets.StakeAmounts = m.NftStakeAmounts
		rets.HalveNumber = m.HalveNumber
	}
	err = GetDB(nil).Table("1_nft_miner_items").Count(&rets.Count).Error
	if err != nil {
		return &rets, err
	}

	var his History
	var nftIns SumAmount
	err = GetDB(nil).Table(his.TableName()).Select("sum(amount)").Where("type = 12").Take(&nftIns).Error
	if err != nil {
		log.Info("get nft metaverse sum failed:", err.Error())
		return nil, err
	}
	rets.RewardAmount = nftIns.Sum.String()

	err = GetDB(nil).Table(p.TableName()).Group("nation").Count(&rets.Region).Error
	if err != nil {
		log.Info("get nft nation failed:", err.Error())
		return nil, err
	}

	var stak NftMinerStaking
	if err := GetDB(nil).Table(stak.TableName()).Select("token_id").Group("token_id").Count(&rets.StakedCount).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &rets, err
		}
	}

	if err := GetDB(nil).Table(stak.TableName()).Where("staking_status = 1").Count(&rets.StakingCount).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &rets, err
		}
	}
	nowTime := time.Now().Unix()
	var energyPower SumAmount
	if err := GetDB(nil).Table(stak.TableName()).Select("sum(energy_power)").Where("? >= start_dated AND ? < end_dated", nowTime, nowTime).Take(&energyPower).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &rets, err
		}
	}
	rets.EnergyPower = energyPower.Sum

	return &rets, nil
}

func GetNftMinerMap() (*[]Positioning, error) {
	var rets []Positioning
	var list []NftMinerItems
	if err := GetDB(nil).Select("coordinates,id,energy_point").Where("merge_status = 1").Find(&list).Error; err != nil {
		return nil, err
	}
	rets = make([]Positioning, len(list))
	for i := 0; i < len(list); i++ {
		var ne NftMinerCoordinate
		if err := json.Unmarshal([]byte(list[i].Coordinates), &ne); err != nil {
			return nil, err
		}
		rets[i].Lat = ne.Latitude
		rets[i].Lng = ne.Longitude
		rets[i].Val = fmt.Sprintf("Miner:#%d EP:%d", list[i].ID, list[i].EnergyPoint)
	}
	return &rets, nil

}

func NftMinerTableIsExist() bool {
	var p NftMinerItems
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}

func GetTopTwentyNftRegionChart(limit int) (*RegionChangeResponse, error) {
	type NftNation struct {
		Nation string `json:"nation"`
		Total  string `json:"total"`
	}
	type daysNation struct {
		Days string `json:"days"`
		NftNation
	}
	var (
		rets  RegionChangeResponse
		where string
	)
	sql := `
SELECT v5.days,array_to_string(array_agg(v5.nation),',')nation,array_to_string(array_agg(v5.total),',')total FROM(
	SELECT v4.days,v4.nation,v4.total,ROW_NUMBER () OVER (PARTITION BY v4.days ORDER BY v4.total desc) AS rowd FROM(
		SELECT v3.days,v3.nation,count(1) total FROM (
			SELECT v1.days,v2.nation FROM(
				SELECT to_char(to_timestamp(date_created),'yyyy-MM-dd') AS days FROM "1_nft_miner_items" WHERE merge_status = 1 GROUP BY days
			)AS v1
			LEFT JOIN(
				SELECT date_created,nation FROM "1_nft_miner_items" WHERE merge_status = 1
			)AS v2 ON(to_char(to_timestamp(v2.date_created),'yyyy-MM-dd') <= v1.days)
		)AS v3 GROUP BY days,nation ORDER BY days ASC,total DESC
	)AS v4 
)AS v5 GROUP BY days ORDER BY days asc
`
	if limit > 0 {
		where = "rowd <= " + strconv.Itoa(limit)
		sql = fmt.Sprintf(`
SELECT v5.days,array_to_string(array_agg(v5.nation),',')nation,array_to_string(array_agg(v5.total),',')total FROM(
	SELECT v4.days,v4.nation,v4.total,ROW_NUMBER () OVER (PARTITION BY v4.days ORDER BY v4.total desc) AS rowd FROM(
		SELECT v3.days,v3.nation,count(1) total FROM (
			SELECT v1.days,v2.nation FROM(
				SELECT to_char(to_timestamp(date_created),'yyyy-MM-dd') AS days FROM "1_nft_miner_items" WHERE merge_status = 1 GROUP BY days
			)AS v1
			LEFT JOIN(
				SELECT date_created,nation FROM "1_nft_miner_items" WHERE merge_status = 1
			)AS v2 ON(to_char(to_timestamp(v2.date_created),'yyyy-MM-dd') <= v1.days)
		)AS v3 GROUP BY days,nation ORDER BY days ASC,total DESC
	)AS v4 
)AS v5 WHERE %s GROUP BY days ORDER BY days asc
`, where)
	}
	var list []daysNation
	err := GetDB(nil).Raw(sql).Where(where).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("get Top Twenty Nft Region Chart Failed")
		return nil, err
	}
	for _, val := range list {
		rets.Time = append(rets.Time, val.Days)

		var nt NftNation
		nt.Nation = val.Nation
		nt.Total = val.Total

		tList := strings.Split(val.Total, ",")
		vList := strings.Split(val.Nation, ",")
		var rts []any
		for key, value := range vList {
			var nt []any
			if len(tList)-1 >= key {
				nt = append(nt, tList[key])
			}
			nt = append(nt, value)
			nt = append(nt, val.Days)

			rts = append(rts, nt)
		}
		rets.List = append(rets.List, rts)
	}

	return &rets, nil
}

func GetNftMinerRegionList(page, limit int, order string) (*GeneralResponse, error) {
	var (
		list []NftMinerRegionListResponse
		rets GeneralResponse
	)

	if order == "" {
		order = "staking_amount DESC,total DESC"
	}
	err := GetDB(nil).Raw(`
SELECT v3.nation as region,max(v3.total)total,count(v3.token_id)staking_number,sum(coalesce(stake_amount,0))staking_amount FROM(
	SELECT v1.nation,v1.total,v2.token_id,stake_amount FROM(
		SELECT nation,array_agg(id)id_list,count(1)total FROM "1_nft_miner_items" WHERE merge_status = 1 GROUP BY nation
	)AS v1
	LEFT JOIN(
		SELECT token_id,stake_amount FROM "1_nft_miner_staking" WHERE staking_status = 1
	)AS v2 ON(v2.token_id = any(v1.id_list))
)AS v3 GROUP BY nation ORDER BY ? OFFSET ? LIMIT ?
`, order, (page-1)*limit, limit).Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("Get Nft Miner Region List Failed")
		return nil, err
	}
	err = GetDB(nil).Table("1_nft_miner_items").Group("nation").Count(&rets.Total).Error
	if err != nil {
		log.WithFields(log.Fields{"info": err}).Error("Get Nft Miner Region List Total Failed")
		return nil, err
	}
	rets.Page = page
	rets.Limit = limit
	rets.List = list

	return &rets, nil
}
