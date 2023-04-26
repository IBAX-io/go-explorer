package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ecoAmountObject struct {
	sync.Map
}

type EcosystemInfoMap struct {
	sync.Map
}

var (
	Tokens            *EcosystemInfoMap
	EcoNames          *EcosystemInfoMap
	countrys          *EcosystemInfoMap
	ecoTags           *EcosystemInfoMap
	ecoTypes          *EcosystemInfoMap
	ecoCascades       *EcosystemInfoMap
	registrationTypes *EcosystemInfoMap
	registrations     *EcosystemInfoMap
	EcosystemIdList   []int64

	allKeyAmount *ecoAmountObject //key:ecosystem value:decimal.Decimal all keys circulations and staking amount
	EcoTxCount   *EcosystemInfoMap
	EcoDigits    *EcosystemInfoMap
	EcoFuelRate  *EcosystemInfoMap
)

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

var ecoTagMap = map[int]string{
	1: "-", 2: "GameFi", 3: "DeFi", 4: "DEX", 5: "NFT", 6: "Exchange Market", 7: "SocialFi",
	8: "Swap DEX", 9: "DEX Protocol", 10: "Stablecoins", 11: "Sandbox", 12: "Gambling",
}

var ecoTypeMap = map[int]string{
	1: "-", 2: "Finance", 3: "Grants", 4: "Scocial",
	5: "Collector", 6: "Game", 7: "Ventures", 8: "Gambling",
	9: "Media", 10: "SocialMedia", 11: "Enterainment", 12: "Technology",
	13: "Government", 14: "Service", 15: "Business", 16: "Other",
}

var ecoCascadeMap = map[int]string{
	1: "Lending & Borrowing", 2: "Token Swaps", 3: "Trading & Prediction Markets", 4: "Investments", 5: "Payments", 6: "Crowdfunding", 7: "Insurance",
	8: "Portfolios", 9: "Art & Fashion", 10: "Digital Collectibles", 11: "Music", 12: "Action Games", 13: "Role-Playing Games", 14: "Strategy Games",
	15: "Puzzle Games", 16: "Simulation Games", 17: "Adventure Games", 18: "Shooting Games", 19: "Sports Games", 20: "Racing Games", 21: "Music Games",
	22: "Composite Genre", 23: "Metaverse", 24: "Games", 25: "Casino", 26: "Sports", 27: "Others", 28: "Music",
	29: "Video", 30: "Others", 31: "Utilities", 32: "Marketplaces", 33: "Developer Tools", 34: "Browsers",
}

var registrationTypeMap = map[int]string{
	1: "Sole-Proprietorship /Sole Trader",
	2: "Ordinary Business Partnership",
	3: "Limited Partnership (LP)",
	4: "Limited Liability Partnership (LLP)",
	5: "Private Company Limited By Shares",
}

var registrationMap = map[int]string{
	1: "no",
	2: "yes",
}

func InitEcosystemInfo() {
	Tokens = &EcosystemInfoMap{}
	EcoNames = &EcosystemInfoMap{}
	countrys = &EcosystemInfoMap{}
	ecoTags = &EcosystemInfoMap{}
	ecoTypes = &EcosystemInfoMap{}
	ecoCascades = &EcosystemInfoMap{}
	registrationTypes = &EcosystemInfoMap{}
	registrations = &EcosystemInfoMap{}
	allKeyAmount = &ecoAmountObject{}
	EcoTxCount = &EcosystemInfoMap{}
	EcoDigits = &EcosystemInfoMap{}
	EcoFuelRate = &EcosystemInfoMap{}

	for k, v := range countryMap {
		countrys.Store(k, v)
	}
	for k, v := range ecoTagMap {
		ecoTags.Store(k, v)
	}
	for k, v := range ecoTypeMap {
		ecoTypes.Store(k, v)
	}
	for k, v := range ecoCascadeMap {
		ecoCascades.Store(k, v)
	}
	for k, v := range registrationTypeMap {
		registrationTypes.Store(k, v)
	}
	for k, v := range registrationMap {
		registrations.Store(k, v)
	}
}

func (p *EcosystemInfoMap) Get(ecosystem int64) string {
	if p == nil {
		return ""
	}
	value, ok := p.Load(ecosystem)
	if ok {
		if cp, ok := value.(string); !ok {
			return ""
		} else {
			return cp
		}
	}
	return ""
}

func (p *EcosystemInfoMap) GetId(infoId int, defaultValue string) string {
	if p == nil {
		return ""
	}
	value, ok := p.Load(infoId)
	if ok {
		if cp, ok := value.(string); !ok {
			return ""
		} else {
			return cp
		}
	}
	return defaultValue
}

func (p *EcosystemInfoMap) GetInt64(ecosystem int64, defaultValue int64) int64 {
	if p == nil {
		return 0
	}
	value, ok := p.Load(ecosystem)
	if ok {
		if cp, ok := value.(int64); !ok {
			return 0
		} else {
			return cp
		}
	}
	return defaultValue
}

func (p *EcosystemInfoMap) GetFloat64(ecosystem int64, defaultValue float64) float64 {
	if p == nil {
		return 0
	}
	value, ok := p.Load(ecosystem)
	if ok {
		if cp, ok := value.(float64); !ok {
			return 0
		} else {
			return cp
		}
	}
	return defaultValue
}

func (p *EcosystemInfoMap) Len() int {
	var count int
	p.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

func (p *ecoAmountObject) Get(ecosystem int64) (decimal.Decimal, error) {
	if p == nil {
		return decimal.Zero, errors.New("eco amount object null")
	}
	value, ok := p.Load(ecosystem)
	if ok {
		return value.(decimal.Decimal), nil
	}
	return decimal.Zero, fmt.Errorf("eco[%d] amount not exist", ecosystem)
}

func GetAllEcosystemInfo() {
	var (
		list  []Ecosystem
		total int64
	)
	err := GetDB(nil).Model(Ecosystem{}).Count(&total).Error
	if err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("get all ecosystem id total failed")
		return
	}
	err = GetDB(nil).Select("id,token_symbol,name,digits,fee_mode_info").Find(&list).Error
	if err != nil {
		log.WithFields(log.Fields{"INFO": err}).Info("get all ecosystem id list failed")
		return
	}
	EcosystemIdList = nil
	for _, v := range list {
		EcosystemIdList = append(EcosystemIdList, v.ID)
		Tokens.Store(v.ID, v.TokenSymbol)
		EcoNames.Store(v.ID, v.Name)
		EcoDigits.Store(v.ID, v.Digits)
		if v.FeeModeInfo != "" {
			var feeInfo FeeModeInfo
			err := json.Unmarshal([]byte(v.FeeModeInfo), &feeInfo)
			if err == nil {
				followFuel, _ := decimal.NewFromFloat(feeInfo.FollowFuel).Mul(decimal.NewFromInt(100)).Float64()
				EcoFuelRate.Store(v.ID, followFuel)
			} else {
				log.Info("get all ecosystem fee mode failed:", err.Error())
				return
			}
		}
	}
}

func SyncEcosystemInfo() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	GetAllEcosystemInfo()
	getEcosystemTxCount()
}

func GetAllKeysTotalAmount(ecosystem int64) error {
	var totalAmount decimal.Decimal

	err := GetDB(nil).Raw(`
SELECT sum(amount) +
	COALESCE((SELECT sum(output_value) FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ?),0) AS total_amount
FROM "1_keys" WHERE ecosystem = ? AND id <> 0 AND id <> 5555
`, ecosystem, ecosystem).Take(&totalAmount).Error
	if err != nil {
		return err
	}

	if ecosystem == 1 {
		//all staking
		if NodeReady || NftMinerReady {
			var staking decimal.Decimal
			err = GetDB(nil).Raw(`
					SELECT sum(to_number(coalesce(NULLIF(lock->>'nft_miner_stake',''),'0'),'999999999999999999999999999999') +
						to_number(coalesce(NULLIF(lock->>'candidate_referendum',''),'0'),'999999999999999999999999999999') +
						to_number(coalesce(NULLIF(lock->>'candidate_substitute',''),'0'),'999999999999999999999999999999'))
					FROM "1_keys" WHERE ecosystem = 1
			`).Take(&staking).Error
			if err != nil {
				return err
			}
			totalAmount = totalAmount.Add(staking)
		}
		if AirdropReady {
			var staking decimal.Decimal
			err = GetDB(nil).Model(AirdropInfo{}).Select("sum(stake_amount)").Take(&staking).Error
			if err != nil {
				return err
			}
			totalAmount = totalAmount.Add(staking)
		}

	}
	if !totalAmount.IsZero() {
		allKeyAmount.Store(ecosystem, totalAmount)
	}

	return nil
}

func getEcosystemTxCount() {
	for _, id := range EcosystemIdList {
		var total int64
		err := GetDB(nil).Model(LogTransaction{}).Where("ecosystem_id = ?", id).Count(&total).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err, "ecosystem": id}).Error("get ecosystem tx count failed")
			return
		}
		EcoTxCount.Store(id, total)
	}
}
