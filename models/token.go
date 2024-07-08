package models

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"sync"
)

const (
	NameAllowRankEcosystem = "allow_rank_ecosystem"
)

var AllowRankEcosystem int64

type Param struct {
	Id    int64 `gorm:"primary_key;not null"`
	Name  string
	Value string
}

func (p *Param) SetTablePrefix(id int64) (prefix string) {
	return
}

func (p *Param) TableName() string {
	return strconv.FormatInt(conf.GetEnvConf().Defi.Ecosystem, 10) + "_" + `param`
}

func (p *Param) Get(name string) (bool, error) {
	return isFound(GetDB(nil).Where("name = ?", name).First(p))
}

type Pair struct {
	Id          int64 `gorm:"primary_key;not null"`
	Ecosystem1  int64
	Ecosystem2  int64
	Reserve1    decimal.Decimal
	Reserve2    decimal.Decimal
	TxHash      string
	PairAccount string
	Status      string //1:close 2:stop
	Liquidity   decimal.Decimal
	Topic       string
}

func (p *Pair) TableName() string {
	return strconv.FormatInt(conf.GetEnvConf().Defi.Ecosystem, 10) + "_" + `pair`
}

func updateAllowRankEcosystem() {
	if !conf.GetEnvConf().Defi.Enable {
		return
	}
	p := &Param{}
	exist, err := p.Get(NameAllowRankEcosystem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("get allow rank ecosystem")
		return
	}
	if exist {
		AllowRankEcosystem, err = strconv.ParseInt(p.Value, 10, 64)
		if err != nil {
			log.WithFields(log.Fields{"info": err, "name": NameAllowRankEcosystem, "value": p.Value}).Info("parse int64")
			return
		}
	}
}

type PairEcosystem struct {
	Id         int64
	Ecosystem1 int64
	Ecosystem2 int64
	Reserve1   decimal.Decimal
	Reserve2   decimal.Decimal
	Liquidity  decimal.Decimal
}

type pairList struct {
	pairs []PairEcosystem
	sync.RWMutex
}

var allPair pairList

func UpdatePairBuffer() {
	HistoryWG.Add(1)
	defer HistoryWG.Done()
	if !conf.GetEnvConf().Defi.Enable {
		return
	}
	var pair Pair
	var count int64
	err := GetDB(nil).Table(pair.TableName()).Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Update Pair Buffer")
		return
	}
	var pairBuffer []PairEcosystem
	limit := 1000
	for page := 1; int64((page-1)*limit) < count; page++ {
		var listBuffer []PairEcosystem
		err = GetDB(nil).Table(pair.TableName()).Select("ecosystem1,ecosystem2,id,reserve1,reserve2,liquidity").Where("status = 'run'").
			Offset((page - 1) * limit).Limit(limit).Order("id asc").Find(&listBuffer).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Update Pair Buffer DbFind")
			return
		}
		pairBuffer = append(pairBuffer, listBuffer...)
	}
	updatePair(pairBuffer)
}

func updatePair(list []PairEcosystem) {
	allPair.Lock()
	defer allPair.Unlock()
	allPair.pairs = list
}

func GetTokenPrices(ecosystems []int64) (prices map[int64]string, err error) {
	prices = make(map[int64]string)
	if !conf.GetEnvConf().Defi.Enable {
		for _, id := range ecosystems {
			prices[id] = "0"
		}
		return
	}
	allPair.RLock()
	defer allPair.RUnlock()
	digits := int(EcoDigits.GetInt(AllowRankEcosystem, 0))
	num1 := decimal.NewFromFloat(math.Pow10(digits))
	for _, ecosystemId := range ecosystems {
		var price = "0"
		if ecosystemId == AllowRankEcosystem {
			prices[ecosystemId] = "1"
			continue
		}
		for _, pair := range allPair.pairs {
			if (pair.Ecosystem1 == AllowRankEcosystem && pair.Ecosystem2 == ecosystemId ||
				pair.Ecosystem2 == AllowRankEcosystem && pair.Ecosystem1 == ecosystemId) &&
				!pair.Reserve2.IsZero() && !pair.Reserve1.IsZero() {
				num2 := decimal.NewFromFloat(math.Pow10(int(EcoDigits.GetInt(ecosystemId, 0)))) //Value of each token
				if pair.Ecosystem1 == ecosystemId {
					price = pair.Reserve2.DivRound(pair.Reserve1, 30).Mul(num2).DivRound(num1, 30).String()
				} else {
					price = pair.Reserve1.DivRound(pair.Reserve2, 30).Mul(num2).DivRound(num1, 30).String()
				}
			}
		}
		prices[ecosystemId] = price
	}
	return
}

type TokensInfo struct {
	Name      string `gorm:"column:name"`
	Address   string `gorm:"column:address"`
	Symbol    string `gorm:"column:symbol"`
	Decimals  int32  `gorm:"column:decimals"`
	ChainId   string `gorm:"column:chainId"`
	LogoURI   string `gorm:"column:logoURI"`
	ChainName string `gorm:"column:chainName"`
}

// TableName returns name of table
func (ti *TokensInfo) TableName() string {
	return "tokens"
}

var tokensReady bool

func UpdateTokensTableStatus() {
	if tokensReady {
		return
	}
	var p TokensInfo
	if HasTableOrView(p.TableName()) {
		tokensReady = true
		return
	}
	return
}

func (ti *TokensInfo) GetLogoURI(chainId string, address string) (logoURI string, exist bool, err error) {
	exist, err = isFound(GetDB(nil).Model(TokensInfo{}).Select("logoURI").
		Where(`"chainId" = ? AND address = ?`, chainId, address).Take(&logoURI))
	return
}

func getLogoURI(info EcosystemInfo) (logoURI string) {
	if !tokensReady {
		return
	}
	if info.TokenAddress == "" {
		return
	}
	ti := &TokensInfo{}
	var (
		err error
	)
	logoURI, _, err = ti.GetLogoURI(strconv.FormatInt(info.chainId, 10), info.TokenAddress)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "info": info}).Info("Get Logo URL Failed")
		return
	}
	return
}
