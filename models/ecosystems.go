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
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	Tokens   *EcosystemInfoMap
	EcoNames *EcosystemInfoMap
)

type Ecosystem struct {
	ID             int64 `gorm:"primary_key;not null"`
	Name           string
	Info           string `gorm:"type:jsonb(PostgreSQL)"`
	FeeModeInfo    string
	IsValued       int64
	EmissionAmount string `gorm:"type:jsonb(PostgreSQL)"`
	TokenSymbol    string
	TokenName      string
	TypeEmission   int64
	TypeWithdraw   int64
	ControlMode    int64
}

type Combustion struct {
	Flag    int64 `json:"flag"`
	Percent int64 `json:"percent"`
}

type FeeModeInfo struct {
	FeeModeDetail map[string]sqldb.FeeModeFlag `json:"fee_mode_detail"`
	Combustion    Combustion                   `json:"combustion"`
	FollowFuel    float64                      `json:"follow_fuel"`
}

type KeyChangeChart struct {
	Time   []int64 `json:"time"`
	NewKey []int64 `json:"new_key"`
}

type TxChangeChart struct {
	Time []int64 `json:"time"`
	Tx   []int64 `json:"tx"`
}

type BasisEcosystemChartDataResponse struct {
	KeyInfo KeyChangeChart `json:"key_info"`
	TxInfo  TxChangeChart  `json:"tx_info"`
}

type BasisEcosystemResponse struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	LogoHash string         `json:"logo_hash"`
	KeyCount int64          `json:"key_count"`
	TotalTx  int64          `json:"total_tx"`
	KeyInfo  KeyChangeChart `json:"key_info"`
	TxInfo   TxChangeChart  `json:"tx_info"`
}

type EcosystemTotalResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	LogoHash    string `json:"logo_hash"`
	GovernModel int64  `json:"govern_model"`
	FeeModel    int    `json:"fee_model"`
	TokenSymbol string `json:"token_symbol"`
	Creator     string `json:"creator"`
	TotalAmount string `json:"total_amount"`
	Member      int64  `json:"member"`
	Contract    int64  `json:"contract"`
}

//EcosystemTotalResult example
type EcosystemTotalResult struct {
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	Limit    int                       `json:"limit"`
	Sysecosy *BasisEcosystemResponse   `json:"sysecosy,omitempty"`
	Rets     *[]EcosystemTotalResponse `json:"rets,omitempty"`
}

// EcosystemKeyTotalRet is model
type EcosystemKeyTotalRet struct {
	ID              int64  `json:"id" example:"1"`               //ID
	Name            string `json:"name" example:""`              //
	Info            string `json:"info" example:""`              //
	IsValued        int64  `json:"is_valued" example:""`         //
	Emission_amount string `json:"emission_amount" example:""`   //
	TokenSymbol     string `json:"token_symbol" example:""`      //
	Token_name      string `json:"token_name" example:""`        //
	Type_emission   int64  `json:"type_emission" example:""`     //
	Type_withdraw   int64  `json:"type_withdraw" example:""`     //
	Amount          string `json:"amount" example:""`            //
	NickName        string `json:"nick_name" example:""`         //
	MemberImageID   int64  `json:"member_image_id" example:""`   //ID
	MemberImageHash string `json:"member_image_hash" example:""` //url
	MemberInfo      string `json:"member_info" example:""`       //
	Account         string `json:"account" example:""`           //
}

type EcosystemTxCount struct {
	Ecosystem int64  `gorm:"column:ecosystem"`
	Tx        int64  `gorm:"column:tx"`
	Name      string `gorm:"column:name"`
	Total     int64  `gorm:"column:total"`
}

var countryMap = map[int]string{
	1: "Afghanistan", 2: "Albania", 3: "Algeria", 4: "American Samoa", 5: "Andorra", 6: "Angola", 7: "Anguilla",
	8: "Antigua and Barbuda", 9: "Argentina", 10: "Armenia", 11: "Aruba", 12: "Australia", 13: "Austria", 14: "Azerbaijan",
	15: "Bahamas", 16: "Bahrain", 17: "Bangladesh", 18: "Barbados", 19: "Belarus", 20: "Belgium", 21: "Belize",
	22: "Benin", 23: "Bermuda", 24: "Bhutan", 25: "Bolivia", 26: "Bosnia and Herzegovina", 27: "Botswana", 28: "Brazil",
	29: "British Virgin Islands", 30: "Brunei", 31: "Bulgaria", 32: "Burkina Faso", 33: "Burundi", 34: "Cambodia", 35: "Cameroon",
	36: "Canada", 37: "Cape Verde", 38: "Cayman Islands", 39: "Central African Republic", 40: "Chad", 41: "Chile", 42: "China",
	43: "Colombia", 44: "Comoros", 45: "Cook Islands", 46: "Costa Rica", 47: "Croatia", 48: "Cuba", 49: "Curacao",
	50: "Cyprus", 51: "Czech Republic", 52: "Denmark", 53: "Djibouti", 54: "Dominica", 55: "Dominican Republic", 56: "DR Congo",
	57: "Ecuador", 58: "Egypt", 59: "El Salvador", 60: "Equatorial Guinea", 61: "Eritrea", 62: "Estonia", 63: "Eswatini",
	64: "Ethiopia", 65: "Falkland Islands", 66: "Faroe Islands", 67: "Fiji", 68: "Finland", 69: "France", 70: "French Guiana",
	71: "French Polynesia", 72: "Gabon", 73: "Gambia", 74: "Georgia", 75: "Germany", 76: "Ghana", 77: "Gibraltar",
	78: "Greece", 79: "Greenland", 80: "Grenada", 81: "Guadeloupe", 82: "Guam", 83: "Guatemala", 84: "Guinea",
	85: "Guinea-Bissau", 86: "Guyana", 87: "Haiti", 88: "Honduras", 89: "Hong Kong", 90: "Hungary", 91: "Iceland",
	92: "India", 93: "Indonesia", 94: "Iran", 95: "Iraq", 96: "Ireland", 97: "Isle of Man", 98: "Israel",
	99: "Italy", 100: "Ivory Coast", 101: "Jamaica", 102: "Japan", 103: "Jordan", 104: "Kazakhstan", 105: "Kenya",
	106: "Kiribati", 107: "Kuwait", 108: "Kyrgyzstan", 109: "Laos", 110: "Latvia", 111: "Lebanon", 112: "Lesotho",
	113: "Liberia", 114: "Libya", 115: "Liechtenstein", 116: "Lithuania", 117: "Luxembourg", 118: "Macau", 119: "Madagascar",
	120: "Malawi", 121: "Malaysia", 122: "Maldives", 123: "Mali", 124: "Malta", 125: "Marshall Islands", 126: "Martinique",
	127: "Mauritania", 128: "Mauritius", 129: "Mayotte", 130: "Mexico", 131: "Micronesia", 132: "Moldova", 133: "Monaco",
	134: "Mongolia", 135: "Montenegro", 136: "Montserrat", 137: "Morocco", 138: "Mozambique", 139: "Myanmar", 140: "Namibia",
	141: "Nauru", 142: "Nepal", 143: "Netherlands", 144: "New Caledonia", 145: "New Zealand", 146: "Nicaragua", 147: "Niger",
	148: "Nigeria", 149: "Niue", 150: "North Korea", 151: "North Macedonia", 152: "Northern Mariana Islands", 153: "Norway", 154: "Oman",
	155: "Pakistan", 156: "Palau", 157: "Palestine", 158: "Panama", 159: "Papua New Guinea", 160: "Paraguay", 161: "Peru",
	162: "Philippines", 163: "Poland", 164: "Portugal", 165: "Puerto Rico", 166: "Qatar", 167: "Republic of the Congo", 168: "Reunion",
	169: "Romania", 170: "Russia", 171: "Rwanda", 172: "Saint Kitts and Nevis", 173: "Saint Lucia", 174: "Saint Martin", 175: "Saint Pierre and Miquelon",
	176: "Saint Vincent and the Grenadines", 177: "Samoa", 178: "San Marino", 179: "Sao Tome and Principe", 180: "Saudi Arabia", 181: "Senegal",
	182: "Serbia", 183: "Seychelles", 184: "Sierra Leone", 185: "Singapore", 186: "Sint Maarten", 187: "Slovakia", 188: "Slovenia", 189: "Solomon Islands",
	190: "Somalia", 191: "South Africa", 192: "South Korea", 193: "South Sudan", 194: "Spain", 195: "Sri Lanka", 196: "Sudan",
	197: "Suriname", 198: "Sweden", 199: "Switzerland", 200: "Syria", 201: "Taiwan", 202: "Tajikistan", 203: "Tanzania", 204: "Thailand",
	205: "Timor-Leste", 206: "Togo", 207: "Tokelau", 208: "Tonga", 209: "Trinidad and Tobago", 210: "Tunisia", 211: "Turkey", 212: "Turkmenistan",
	213: "Turks and Caicos Islands", 214: "Tuvalu", 215: "Uganda", 216: "Ukraine", 217: "United Arab Emirates", 218: "United Kingdom", 219: "United States",
	220: "United States Virgin Islands", 221: "Uruguay", 222: "Uzbekistan", 223: "Vanuatu", 224: "Vatican City", 225: "Venezuela", 226: "Vietnam",
	227: "Wallis and Futuna", 228: "Western Sahara", 229: "Yemen", 230: "Zambia", 231: "Zimbabwe",
}

type EcosystemInfoMap struct {
	sync.RWMutex
	Map map[int64]string
}

func (sys *Ecosystem) TableName() string {
	return ecosysTable
}

func (sys *Ecosystem) Get(id int64) (bool, error) {
	return isFound(conf.GetDbConn().Conn().First(sys, "id = ?", id))
}

func (sys *Ecosystem) GetTokenSymbol(id int64) (bool, error) {
	return isFound(GetDB(nil).Select("token_symbol,name").First(sys, "id = ?", id))
}

func GetActiveEcoLibs() (string, error) {
	//default ecoLibs: 1
	type countEcosystem struct {
		Ecosystem int64 `gorm:"column:ecosystem"`
		Count     int64 `gorm:"column:count"`
	}
	var rets countEcosystem
	f, err := isFound(GetDB(nil).Raw(`SELECT count(*),ecosystem FROM "1_history" WHERE type = 1 GROUP BY ecosystem ORDER BY count DESC`).Limit(1).Take(&rets))
	if err != nil {
		log.Info("get active ecoLibs err:", err.Error())
		return SysEcosystemName, err
	}
	if !f {
		return SysEcosystemName, nil
	}

	name := EcoNames.Get(rets.Ecosystem)

	return name, nil
}

func GetActiveEcoLibsToRedis() {
	HistoryWG.Add(1)
	defer func() {
		HistoryWG.Done()
	}()
	rets, err := GetActiveEcoLibs()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get active ecosystem error")
		return
	}

	res, err := msgpack.Marshal(rets)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get active ecosystem msgpack error")
		return
	}

	rd := RedisParams{
		Key:   "active_ecosystem",
		Value: string(res),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get active ecosystem set redis error")
		return
	}
}

func GetActiveEcoLibsFromRedis() (string, error) {
	var rets string
	var err error
	rd := RedisParams{
		Key:   "active_ecosystem",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get active ecosystem From Redis getDb err")
		return rets, err
	}
	err = msgpack.Unmarshal([]byte(rd.Value), &rets)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("get active ecosystem From Redis msgpack err")
		return rets, err
	}

	return rets, nil
}

func (p *Ecosystem) GetBasisEcosystem() (*BasisEcosystemResponse, error) {
	var (
		rets BasisEcosystemResponse
		key  Key
		lg   LogTransaction
	)
	f, err := p.Get(1)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errors.New("get basis ecosystem doesn't exist")
	}
	rets.ID = p.ID
	rets.Name = p.Name

	escape := func(value any) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}
	if p.Info != "" {
		minfo := make(map[string]any)
		err := json.Unmarshal([]byte(p.Info), &minfo)
		if err != nil {
			return nil, err
		}
		usid, ok := minfo["logo"]
		if ok {
			urid := escape(usid)
			uid, err := strconv.ParseInt(urid, 10, 64)
			if err != nil {
				return nil, err
			}

			hash, err := GetFileHash(uid)
			if err != nil {
				return nil, err
			}
			rets.LogoHash = hash

		}
	}

	if err := GetDB(nil).Table(key.TableName()).Where("ecosystem = 1").Count(&rets.KeyCount).Error; err != nil {
		return nil, err
	}

	if err := GetDB(nil).Table(lg.TableName()).Where("ecosystem_id = 1").Count(&rets.TotalTx).Error; err != nil {
		return nil, err
	}

	return &rets, nil
}

func (p *Ecosystem) GetBasisEcosystemChart() (*BasisEcosystemChartDataResponse, error) {
	var (
		basis BasisEcosystemChartDataResponse
	)

	const getDay = 15

	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	t1 := yesterday.AddDate(0, 0, -1*getDay)

	var txList []DaysNumber
	err := GetDB(nil).Raw(`SELECT to_char(to_timestamp(timestamp/1000),'yyyy-MM-dd') days,count(*) num 
FROM "log_transactions" WHERE ecosystem_id = 1 AND timestamp >= ? GROUP BY days`, t1.UnixMilli()).Find(&txList).Error
	if err != nil {
		return nil, err
	}

	basis.TxInfo.Time = make([]int64, getDay)
	basis.TxInfo.Tx = make([]int64, getDay)
	for i := 0; i < len(basis.TxInfo.Time); i++ {
		basis.TxInfo.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		basis.TxInfo.Tx[i] = GetDaysNumber(basis.TxInfo.Time[i], txList)
	}

	var keyList []DaysNumber
	err = GetDB(nil).Raw(`SELECT to_char(to_timestamp(created_at/1000),'yyyy-MM-dd') days,count(*) num FROM "1_history" WHERE ecosystem = 1 
AND comment = 'New User' AND type = 4 AND created_at >= ? GROUP BY days`, t1.UnixMilli()).Find(&keyList).Error
	if err != nil {
		return nil, err
	}

	basis.KeyInfo.Time = make([]int64, getDay)
	basis.KeyInfo.NewKey = make([]int64, getDay)
	for i := 0; i < len(basis.KeyInfo.Time); i++ {
		basis.KeyInfo.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		basis.KeyInfo.NewKey[i] = GetDaysNumber(basis.KeyInfo.Time[i], keyList)
	}

	return &basis, nil
}

func (p *Ecosystem) GetEcoSystemList(limit, page int, order string, where map[string]any) (int64, *[]EcosystemTotalResponse, error) {
	var (
		total int64
		list  []EcosystemTotalResponse
	)
	type ecoListResponse struct {
		Ecosystem
		Creator  string `json:"creator"`
		Member   int64  `json:"member"`
		Contract int64  `json:"contract"`
	}
	var ecoList []ecoListResponse

	if order == "" {
		order = "id desc"
	}

	if strings.Contains(order, "member") || strings.Contains(order, "contract") {
		str := strings.Split(order, " ")
		if len(str) != 2 || (str[1] != "desc" && str[1] != "DESC" && str[1] != "ASC" && str[1] != "asc") {
			return 0, nil, errors.New("order by request params invalid")
		}

	} else {
		if strings.Contains(order, "fee_model") {
			str := strings.Split(order, " ")
			if len(str) != 2 || (str[1] != "desc" && str[1] != "DESC" && str[1] != "ASC" && str[1] != "asc") {
				return 0, nil, errors.New("order by request params invalid")
			}
			order = `coalesce(fee_mode_info->'fee_mode_detail'->'vmCost_fee'->>'flag','0')||
					coalesce(fee_mode_info->'fee_mode_detail'->'element_fee'->>'flag','0')||
					coalesce(fee_mode_info->'fee_mode_detail'->'storage_fee'->>'flag','0')||
					coalesce(fee_mode_info->'fee_mode_detail'->'expedite_fee'->>'flag','0') ` + str[1]
		} else if strings.Contains(order, "govern_model") {
			str := strings.Split(order, " ")
			if len(str) != 2 || (str[1] != "desc" && str[1] != "DESC" && str[1] != "ASC" && str[1] != "asc") {
				return 0, nil, errors.New("order by request params invalid")
			}
			order = "control_mode " + str[1]
		}
	}

	if len(where) == 0 {
		if err := GetDB(nil).Table(p.TableName()).Count(&total).Error; err != nil {
			return 0, nil, err
		}
		if err := GetDB(nil).Table(`"1_ecosystems" AS e`).Select(`*,
(SELECT count(*) from "1_keys" AS k WHERE k.ecosystem = e.id)as member,
(SELECT count(*) from "1_contracts" AS c WHERE c.ecosystem = e.id)as contract,
(SELECT value from "1_parameters" AS p WHERE p.name = 'founder_account' AND e.id = p.ecosystem LIMIT 1)as creator`).
			Order(order).Offset((page - 1) * limit).Limit(limit).Find(&ecoList).Error; err != nil {
			return 0, nil, err
		}
	} else {
		cond, vals, err := WhereBuild(where)
		if err != nil {
			return 0, nil, err
		}
		if err := GetDB(nil).Table(p.TableName()).Where(cond, vals...).Count(&total).Error; err != nil {
			return 0, nil, err
		}

		if err := GetDB(nil).Table(`"1_ecosystems" AS e`).Select(`*,
(SELECT count(*) from "1_keys" AS k WHERE k.ecosystem = e.id)as member,
(SELECT count(*) from "1_contracts" AS c WHERE c.ecosystem = e.id)as contract,
(SELECT value from "1_parameters" AS p WHERE p.name = 'founder_account' AND e.id = p.ecosystem LIMIT 1)as creator`).
			Where(cond, vals...).Order(order).Offset((page - 1) * limit).Limit(limit).Find(&ecoList).Error; err != nil {
			return 0, nil, err
		}
	}
	list = make([]EcosystemTotalResponse, len(ecoList))
	type emsAmount struct {
		Val  decimal.Decimal `json:"val"`
		Time string          `json:"time"`
		Type string          `json:"type"`
	}
	var emissionAmount []emsAmount
	escape := func(value any) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}
	for i := 0; i < len(ecoList); i++ {
		list[i].ID = ecoList[i].ID
		list[i].TokenSymbol = ecoList[i].TokenSymbol
		list[i].Name = ecoList[i].Name
		if ecoList[i].EmissionAmount != "" {
			info := ecoList[i].EmissionAmount
			if err := json.Unmarshal([]byte(info), &emissionAmount); err != nil {
				return 0, nil, err
			}
			amount := decimal.New(0, 0)
			for k, v := range emissionAmount {
				if v.Type == "emission" && k == 0 {
					amount = amount.Add(v.Val)
				}
			}
			list[i].TotalAmount = amount.String()
		}
		if ecoList[i].ID == 1 {
			list[i].TokenSymbol = SysTokenSymbol
			list[i].TotalAmount = TotalSupplyToken
		}

		if ecoList[i].Info != "" {
			minfo := make(map[string]any)
			err := json.Unmarshal([]byte(ecoList[i].Info), &minfo)
			if err != nil {
				log.Info("json failed:", err.Error())
				return 0, nil, err
			}
			usid, ok := minfo["logo"]
			if ok {
				urid := escape(usid)
				uid, err := strconv.ParseInt(urid, 10, 64)
				if err != nil {
					return 0, nil, err
				}

				hash, err := GetFileHash(uid)
				if err != nil {
					return 0, nil, err
				}
				list[i].LogoHash = hash

			}
		}
		feeMode, err := getEcosystemFeeMode(ecoList[i].FeeModeInfo)
		if err != nil {
			return 0, nil, err
		}
		list[i].FeeModel = feeMode
		list[i].GovernModel = ecoList[i].ControlMode
		var sp StateParameter
		sp.ecosystem = list[i].ID
		found1, err1 := sp.Get(`founder_account`)
		if err1 != nil {
			return 0, nil, err1
		}

		if !found1 || len(sp.Value) == 0 {
			return 0, nil, errors.New("get ecosystem creator invalid")
		}
		keyId, err := strconv.ParseInt(sp.Value, 10, 64)
		if err != nil {
			return 0, nil, errors.New("get ecosystem creator keyId invalid")
		}
		list[i].Creator = converter.AddressToString(keyId)
		var key Key
		if err := GetDB(nil).Table(key.TableName()).Where("ecosystem = ?", list[i].ID).Count(&list[i].Member).Error; err != nil {
			return 0, nil, err
		}
		var crt Contract
		list[i].Contract = crt.GetContractsByEcoLibs(list[i].ID)
	}
	//if isSort {
	//	offset := (page - 1) * limit
	//	sort.Sort(EcoList(list))
	//	if len(list) >= offset {
	//		list = list[offset:]
	//		if len(list) >= limit {
	//			list = list[:limit]
	//		}
	//	} else {
	//		return 0, nil, nil
	//	}
	//}

	return total, &list, nil
}

func getEcosystemFeeMode(info string) (int, error) {
	var (
		feeMode int = 1
	)
	if info != "" {
		var feeInfo FeeModeInfo
		err := json.Unmarshal([]byte(info), &feeInfo)
		if err != nil {
			log.Info("get ecosystem fee mode json failed:", err.Error())
			return 0, err
		}
		for key, value := range feeInfo.FeeModeDetail {
			switch key {
			case "vmCost_fee", "element_fee", "storage_fee", "expedite_fee":
				if value.FlagToInt() > 1 {
					feeMode = 2
				}
			}
		}
	}
	return feeMode, nil
}

func GetEcosystemDetailInfo(search any) (*EcosystemDetailInfoResponse, error) {
	var (
		rets EcosystemDetailInfoResponse
	)
	type ecoinfo struct {
		Ecosystem
		Creator string `json:"creator"`
	}
	var eco ecoinfo

	switch reflect.TypeOf(search).String() {
	case "string":
		name := search.(string)
		if utf8.RuneCountInString(name) > 300 {
			return nil, errors.New("request params invalid")
		}
		f, err := isFound(GetDB(nil).Table(`"1_ecosystems" as e`).Select(`*,
(SELECT value from "1_parameters" AS p WHERE p.name = 'founder_account' AND e.id = p.ecosystem LIMIT 1)as creator`).Where("name = ?", name).First(&eco))
		if err != nil {
			return nil, err
		}
		if !f {
			err = GetDB(nil).Table(`"1_ecosystems" as e`).Select(`*,
(SELECT value from "1_parameters" AS p WHERE p.name = 'founder_account' AND e.id = p.ecosystem LIMIT 1)as creator`).Where("token_symbol = ?", name).First(&eco).Error
			if err != nil {
				return nil, err
			}
		}
	case "json.Number":
		ecosystemId, err := search.(json.Number).Int64()
		if err != nil {
			return nil, err
		}
		if err := GetDB(nil).Table(`"1_ecosystems" as e`).Select(`*,
(SELECT value from "1_parameters" AS p WHERE p.name = 'founder_account' AND e.id = p.ecosystem LIMIT 1)as creator`).Where("id = ?", ecosystemId).First(&eco).Error; err != nil {
			return nil, err
		}
	default:
		log.WithFields(log.Fields{"search type": reflect.TypeOf(search).String()}).Warn("Get Ecosystem Detail Info Failed")
		return nil, errors.New("request params invalid")
	}

	type emsAmount struct {
		Val  decimal.Decimal `json:"val"`
		Time string          `json:"time"`
		Type string          `json:"type"`
	}
	var emissionAmount []emsAmount
	escape := func(value any) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}
	rets.EcosystemId = eco.ID
	rets.TokenSymbol = eco.TokenSymbol
	rets.Ecosystem = eco.Name
	if eco.TypeWithdraw == 2 {
		rets.IsWithdraw = true
	}
	if eco.TypeEmission == 2 {
		rets.IsEmission = true
	}
	if eco.EmissionAmount != "" {
		info := eco.EmissionAmount
		if err := json.Unmarshal([]byte(info), &emissionAmount); err != nil {
			return nil, err
		}
		total := decimal.New(0, 0)
		withdraw := decimal.New(0, 0)
		emission := decimal.New(0, 0)
		for _, v := range emissionAmount {
			switch v.Type {
			case "emission":
				if total.GreaterThan(decimal.Zero) {
					emission = emission.Add(v.Val)
				} else {
					total = total.Add(v.Val)
				}
			case "burn":
				withdraw = withdraw.Add(v.Val)
			}
		}
		rets.TotalAmount = total.String()
		rets.Emission = emission.String()
		rets.Withdraw = withdraw.String()
	}
	if eco.Info != "" {
		minfo := make(map[string]any)
		err := json.Unmarshal([]byte(eco.Info), &minfo)
		if err != nil {
			log.Info("Get Ecosystem Detail json failed:", err.Error())
			return nil, err
		}
		usid, ok := minfo["logo"]
		if ok {
			urid := escape(usid)
			uid, err := strconv.ParseInt(urid, 10, 64)
			if err != nil {
				return nil, err
			}

			hash, err := GetFileHash(uid)
			if err != nil {
				return nil, err
			}
			rets.LogoHash = hash
		}

		for k, v := range minfo {
			switch k {
			case "description":
				rets.EcoIntroduction = fmt.Sprint(v)
			case "type", "tag", "cascade", "registered", "country", "registration_type":
				value, _ := strconv.Atoi(fmt.Sprint(v))
				if value > 0 {
					switch k {
					case "type":
						rets.EcoType = value
					case "tag":
						rets.EcoTag = value
					case "cascade":
						rets.EcoCascade = value
					case "registered":
						rets.Registered = value
					case "country":
						rets.Country = rets.getCountry(value)
					case "registration_type":
						rets.RegistrationType = value
					}
				}
			case "registration_no":
				rets.RegistrationNo = fmt.Sprint(v)
			case "discord", "twitter", "youtube", "facebook", "telegram", "github", "medium", "web_page":
				value := fmt.Sprint(v)
				if rets.Social == nil {
					rets.Social = make(map[string]string)
				}
				if k == "web_page" {
					rets.WebPage = value
				} else {
					rets.Social[k] = value
				}
			}

		}
	}
	//defalut
	rets.FeeModel = 1
	rets.FollowFuel = 100
	rets.FeeModelVmcost.ConversionRate = "100"
	rets.FeeModelVmcost.Flag = "1"
	rets.FeeModeElement.ConversionRate = "100"
	rets.FeeModeElement.Flag = "1"
	rets.FeeModeStorage.ConversionRate = "100"
	rets.FeeModeStorage.Flag = "1"
	rets.FeeModeExpedite.ConversionRate = "100"
	rets.FeeModeExpedite.Flag = "1"
	if eco.FeeModeInfo != "" {
		var feeInfo FeeModeInfo
		err := json.Unmarshal([]byte(eco.FeeModeInfo), &feeInfo)
		if err != nil {
			log.Info("Get Ecosystem Detail feeInfo json failed:", err.Error())
			return nil, err
		}
		rets.FollowFuel = feeInfo.FollowFuel * 100
		for key, value := range feeInfo.FeeModeDetail {
			switch key {
			case "vmCost_fee":
				if value.FlagToInt() > 1 {
					rets.MultiFee = true
					rets.FeeModel = 2
					rets.FeeModeAccount = getFeeModeAccount(eco.ID)
				}
				rets.FeeModelVmcost.ConversionRate = value.ConversionRate
				rets.FeeModelVmcost.Flag = value.Flag
			case "element_fee":
				if value.FlagToInt() > 1 {
					rets.MultiFee = true
					rets.FeeModel = 2
					rets.FeeModeAccount = getFeeModeAccount(eco.ID)
				}
				rets.FeeModeElement.ConversionRate = value.ConversionRate
				rets.FeeModeElement.Flag = value.Flag
			case "storage_fee":
				if value.FlagToInt() > 1 {
					rets.MultiFee = true
					rets.FeeModel = 2
					rets.FeeModeAccount = getFeeModeAccount(eco.ID)
				}
				rets.FeeModeStorage.ConversionRate = value.ConversionRate
				rets.FeeModeStorage.Flag = value.Flag
			case "expedite_fee":
				if value.FlagToInt() > 1 {
					rets.MultiFee = true
					rets.FeeModel = 2
					rets.FeeModeAccount = getFeeModeAccount(eco.ID)
				}
				rets.FeeModeExpedite.ConversionRate = value.ConversionRate
				rets.FeeModeExpedite.Flag = value.Flag
			}
		}
		if rets.FeeModel == 2 {
			if feeInfo.FeeModeDetail["vmCost_fee"].FlagToInt() > 1 &&
				feeInfo.FeeModeDetail["element_fee"].FlagToInt() > 1 &&
				feeInfo.FeeModeDetail["storage_fee"].FlagToInt() > 1 &&
				feeInfo.FeeModeDetail["expedite_fee"].FlagToInt() > 1 {
				rets.WithholdingMode = 2
			} else {
				rets.WithholdingMode = 1
			}
		}
		if feeInfo.Combustion.Flag > 1 {
			rets.IsCombustion = true
			rets.CombustionPercent = feeInfo.Combustion.Percent
		}
	}
	rets.Combustion = getEcosystemCombustion(eco.ID)
	cir, err := GetCirculations(eco.ID)
	if err != nil {
		return nil, err
	}
	rets.Circulations = cir

	rets.GovernModel = eco.ControlMode
	creatorId, err := strconv.ParseInt(eco.Creator, 10, 64)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "creator:": eco.Creator}).Error("Get Ecosystem Creator Failed")
		return nil, err
	}
	rets.Creator = converter.AddressToString(creatorId)

	if eco.ID != 1 {
		var rk sqldb.RollbackTx
		var his History
		findStr := fmt.Sprintf(`{"type":"NewEcosystem","id":%d}`, eco.ID)
		if err := GetDB(nil).Select("block_id").Where("table_name = '@system' AND data = ?", findStr).First(&rk).Error; err != nil {
			log.Info("Get Ecosystem Detail New Ecosystem rollback Failed:", err)
			return nil, errors.New("get ecosystem detail Create blockId Failed")
		}
		rets.BlockId = rk.BlockID
		newEcosystemComment := "taxes for execution of @1NewEcosystem contract"
		if err := GetDB(nil).Select("created_at,txhash").Where("type = 1 AND comment = ? AND block_id = ?", newEcosystemComment, rets.BlockId).First(&his).Error; err != nil {
			log.Info("Get Ecosystem Detail New Ecosystem history Failed:", err)
			return nil, errors.New("get ecosystem detail Create Time Failed")
		}
		rets.Hash = hex.EncodeToString(his.Txhash)
		rets.Time = MsToSeconds(his.Createdat)
	} else {
		rets.BlockId = 1
		rets.Time = FirstBlockTime
		rets.TotalAmount = TotalSupplyToken
		rets.TokenSymbol = SysTokenSymbol
		rets.EcoType = 1
		rets.Country = rets.getCountry(185)
		if rets.Social == nil {
			rets.Social = make(map[string]string)
		}
		rets.Social["discord"] = ""
		rets.Social["twitter"] = ""
		rets.Social["youtube"] = ""
		rets.Social["facebook"] = ""
		rets.Social["telegram"] = ""
		rets.Social["github"] = ""
		rets.Social["medium"] = ""
		rets.Social["web_page"] = ""

		var ts LogTransaction
		if err := GetDB(nil).Select("hash").Where("block = ?", rets.BlockId).First(&ts).Error; err != nil {
			log.Info("Get Ecosystem Detail detail hash Failed:", err)
			return nil, errors.New("get ecosystem detail hash failed")
		}
		rets.Hash = hex.EncodeToString(ts.Hash)
	}

	//todo:need add rets.GovernCommittee

	return &rets, nil
}

func EcosystemSearch(search string, account string) (*[]EcosystemSearchResponse, error) {
	var list []Ecosystem
	var rets []EcosystemSearchResponse
	like := "%" + search + "%"
	wid := converter.StringToAddress(account)
	if account != "" {
		if err := GetDB(nil).Select("name,id").Where(`id in(SELECT ecosystem FROM "1_keys" WHERE 
id = ?) and name like ?`, wid, like).Limit(10).Find(&list).Error; err != nil {
			log.Info("EcosystemSearch failed:", err, " like:", like, " account:", account)
			return nil, errors.New("search account ecosystem failed")
		}
	} else {
		if err := GetDB(nil).Select("id,name").Where("name like ?", like).Limit(10).Find(&list).Error; err != nil {
			log.Info("EcosystemSearch failed:", err, " like:", like)
			return nil, errors.New("search ecosystem failed")
		}
	}
	rets = make([]EcosystemSearchResponse, len(list))
	for i, value := range list {
		rets[i].Id = value.ID
		rets[i].Name = value.Name
	}
	return &rets, nil
}

//GetEcosystemDatabase
//reqType params: 1:tableName 2:tableColumns 3:tableRows
func GetEcosystemDatabase(page, limit, reqType int, ecosystemId int64, search, order string) (*GeneralResponse, error) {
	var str sqldb.Table
	var total int64
	var ret GeneralResponse
	ret.Page = page
	ret.Limit = limit
	if reqType == 1 {
		if page <= 0 || limit <= 0 {
			return nil, errors.New("request params invalid")
		}
		var list []sqldb.Table
		var rets []string
		var like string
		like = "%" + search + "%"
		if err := GetDB(nil).Table(str.TableName()).Where("ecosystem = ? AND name like ?", ecosystemId, like).Count(&total).Error; err != nil {
			log.Info("Ecosystem dataBase failed:", err, " like:", like, " search:", search)
			return nil, err
		}
		ret.Total = total

		if err := GetDB(nil).Select("name").Where("ecosystem = ? AND name like ?", ecosystemId, like).Offset((page - 1) * limit).Limit(limit).Find(&list).Error; err != nil {
			return nil, err
		}
		rets = make([]string, len(list))
		for key, value := range list {
			rets[key] = value.Name
		}
		ret.List = rets

	} else if reqType == 2 {
		if search == "" {
			return nil, errors.New("request params invalid")
		}
		search = strconv.FormatInt(ecosystemId, 10) + "_" + search
		order = "ordinal_position ASC"
		sqlQuery := fmt.Sprintf("SELECT column_name,data_type,column_default FROM information_schema.columns WHERE table_name='%s' ORDER BY %s", search, order)
		rows, err := sqldb.GetDB(nil).Raw(sqlQuery).Rows()
		if err != nil {
			return nil, err
		}
		list, err := sqldb.GetResult(rows)
		if err != nil {
			return nil, err
		}
		ret.Total = int64(len(list))
		ret.List = list
	} else {
		if search == "" || page <= 0 || limit <= 0 || order == "" {
			return nil, errors.New("request params invalid")
		}
		search = strconv.FormatInt(ecosystemId, 10) + "_" + search
		num, err := sqldb.GetNodeRows(search)
		if err != nil {
			return nil, err
		}

		ret.Total = num
		var sqlQuest string
		sqlQuest = fmt.Sprintf(`select * from "%s" order by %s offset %d limit %d`, search, order, (page-1)*limit, limit)
		rows, err := sqldb.GetDB(nil).Raw(sqlQuest).Rows()
		if err != nil {
			return nil, fmt.Errorf("getEcoDatabase rows err:%s in query %s", err, sqlQuest)
		}

		ret.List, err = sqldb.GetRowsInfo(rows, sqlQuest)
		if err != nil {
			return nil, err
		}

	}
	return &ret, nil
}

func GetEcosystemApp(page, limit int, ecosystemId, appId int64, order, search string) (*GeneralResponse, error) {
	var app Applications
	var ret GeneralResponse
	ret.Page = page
	ret.Limit = limit
	if order == "" {
		order = "id desc"
	}

	if appId != 0 {
		if search == "" {
			return nil, errors.New("search request params invalid")
		}
		f, err := app.GetById(appId)
		if err != nil {
			return nil, err
		}
		if f {
			var rets EcosystemAppInfo
			rets.AppId = app.ID
			switch search {
			case "contracts":
				contract, err := getAppContractsParams(app.ID, app.Ecosystem, false)
				if err != nil {
					return nil, err
				}
				rets.Contracts = contract
			case "page":
				pageDate, err := getAppPageParams(app.ID, app.Ecosystem, false)
				if err != nil {
					return nil, err
				}
				rets.Page = pageDate
			case "snippets":
				snippets, err := getAppSnippetsParams(app.ID, app.Ecosystem, false)
				if err != nil {
					return nil, err
				}
				rets.Snippets = snippets
			case "table":
				table, err := getAppTableParams(app.ID, app.Ecosystem, false)
				if err != nil {
					return nil, err
				}
				rets.Table = table
			case "params":
				params, err := getAppParams(app.ID, app.Ecosystem, false)
				if err != nil {
					return nil, err
				}
				rets.Params = params
			default:
				return nil, errors.New("search request params unknown")
			}
			ret.List = rets
		}
	} else {
		var list []EcosystemAppList
		rets, total, err := app.FindApp(page, limit, order, fmt.Sprintf("ecosystem = %d AND deleted != %d", ecosystemId, 1))
		if err != nil {
			return nil, err
		}
		list = make([]EcosystemAppList, len(*rets))
		for key, value := range *rets {
			list[key].Name = value.Name
			list[key].AppId = value.ID
			list[key].Conditions = value.Conditions
		}
		ret.Total = total
		ret.List = list
	}
	return &ret, nil
}

func GetEcoLibsTransaction() ([]EcosystemTxRatioChart, error) {
	var rets []EcosystemTxRatioChart
	var list []EcosystemTxCount
	if err := GetDB(nil).Raw(`SELECT * FROM (SELECT count(*) as tx,ecosystem_id as ecosystem FROM "log_transactions" GROUP BY ecosystem_id) AS log 
INNER JOIN (
	SELECT name,id FROM "1_ecosystems"
) AS es  ON (log.ecosystem = es.id) ORDER BY log.tx DESC`).Find(&list).Error; err != nil {
		return nil, err
	}

	var total int64
	var logTx LogTransaction

	err := GetDB(nil).Table(logTx.TableName()).Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get EcoLibs Tx Chart Data From Redis Total Failed")
		return nil, err
	}
	totalDec := decimal.NewFromInt(total)
	zeroDec := decimal.New(0, 0)
	if totalDec.LessThan(zeroDec) {
		return nil, errors.New("tx Chart Total Is Zero")
	}

	var orderList []EcosystemTxCount
	var orderTxTotal int64
	for key, value := range list {
		var ratio float64
		txDec := decimal.NewFromInt(value.Tx)
		if txDec.GreaterThan(zeroDec) {
			ratio, _ = txDec.Mul(decimal.NewFromInt(100)).DivRound(totalDec, 2).Float64()
		}
		if key >= 10 {
			orderList = append(orderList, value)
		} else {
			rets = append(rets, EcosystemTxRatioChart{Name: value.Name, Value: ratio})
		}
	}
	for _, value := range orderList {
		orderTxTotal += value.Tx
	}
	if orderTxTotal > 0 {
		var ratio float64
		txDec := decimal.NewFromInt(orderTxTotal)
		if txDec.GreaterThan(zeroDec) {
			ratio, _ = txDec.Mul(decimal.NewFromInt(100)).DivRound(totalDec, 2).Float64()
		}
		rets = append(rets, EcosystemTxRatioChart{Name: "Other EcoLibs", Value: ratio})
	}

	return rets, nil
}

func GetEcosystemLogoHash(ecosystem int64) (string, string) {
	var tokenSymbol string
	es := Ecosystem{}
	escape := func(value any) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}
	if ecosystem == 1 {
		tokenSymbol = SysTokenSymbol
	}
	f, err := es.Get(ecosystem)
	if f && err == nil {
		if ecosystem != 1 {
			tokenSymbol = es.TokenSymbol
		}
		if es.Info != "" {
			minfo := make(map[string]any)
			err := json.Unmarshal([]byte(es.Info), &minfo)
			if err != nil {
				log.Info("GetEcosystemLogoHash json failed:", err.Error())
				return "", tokenSymbol
			}

			usid, ok := minfo["logo"]
			if ok {
				urid := escape(usid)
				uid, err := strconv.ParseInt(urid, 10, 64)
				if err != nil {
					log.Info("GetEcosystemLogoHash parse int failed:", err.Error())
					return "", tokenSymbol
				}

				hash, err := GetFileHash(uid)
				if err != nil {
					log.Info("GetEcosystemLogoHash GetFileHash failed:", err.Error())
					return "", tokenSymbol
				}
				return hash, tokenSymbol
			}
		}
	}
	return "", tokenSymbol
}

func getFeeModeAccount(ecosystem int64) string {
	var (
		param StateParameter
	)
	param.SetTableFix(ecosystem)
	f, err := param.Get("ecosystem_wallet")
	if err != nil {
		log.WithFields(log.Fields{"error": err, "ecosystem:": ecosystem}).Error("get Fee Mode Account Failed")
		return ""
	}
	if f {
		return param.Value
	}
	return ""
}

func GetEcosystemTokenSymbol(ecosystem int64) (tokenSymbol, name string) {
	if ecosystem == 1 {
		name = SysEcosystemName
		tokenSymbol = SysTokenSymbol
		return
	}
	var es Ecosystem
	f, err := es.GetTokenSymbol(ecosystem)
	if f && err == nil {
		name = es.Name
		tokenSymbol = es.TokenSymbol
	}
	return
}

func getEcosystemCombustion(ecosystem int64) string {
	var his History
	var sum SumAmount
	err := GetDB(nil).Table(his.TableName()).Select("sum(amount)").Where("type = 16 AND ecosystem = ?", ecosystem).Take(&sum).Error
	if err != nil {
		return "0"
	}
	return sum.Sum.String()
}

func (p *countryInfo) getCountry(country int) string {
	if val, ok := countryMap[country]; ok {
		return val
	}

	return ""
}

func GetAllTokenSymbol() ([]Ecosystem, error) {
	var (
		list []Ecosystem
	)
	err := GetDB(nil).Select("token_symbol,id").Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("get all token symbol failed")
		return nil, err
	}
	return list, nil
}

func GetAllEcosystemName() ([]Ecosystem, error) {
	var (
		list []Ecosystem
	)
	err := GetDB(nil).Select("name,id").Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("get all ecosystem name failed")
		return nil, err
	}
	return list, nil
}

func (p *EcosystemInfoMap) Get(ecosystem int64) string {
	p.RLock()
	defer p.RUnlock()
	value, ok := p.Map[ecosystem]
	if ok {
		return value
	}
	return ""
}

func (p *EcosystemInfoMap) Set(ecosystem int64, value string) {
	p.Lock()
	defer p.Unlock()
	p.Map[ecosystem] = value
}

func InitEcosystemInfo() {
	Tokens = &EcosystemInfoMap{
		Map: make(map[int64]string),
	}
	EcoNames = &EcosystemInfoMap{
		Map: make(map[int64]string),
	}
}

func SyncEcosystemInfo() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	list, err := GetAllTokenSymbol()
	if err == nil {
		for _, val := range list {
			Tokens.Set(val.ID, val.TokenSymbol)
		}
	}
	list, err = GetAllEcosystemName()
	if err == nil {
		for _, val := range list {
			EcoNames.Set(val.ID, val.Name)
		}
	}
}
