package models

type BridgeToken struct {
	Id           int64  `gorm:"primary_key;not null"`
	Status       int64  `gorm:"not null"`
	SettingId    int64  `gorm:"not null"`
	ChainName    string `gorm:"not null"`
	Ecosystem    int64  `gorm:"not null"`
	TokenAddress string `gorm:"not null"`
	TokenDigits  int64  `gorm:"not null"`
	TokenName    string `gorm:"not null"`
	TokenSymbol  string `gorm:"not null"`
}

func (b *BridgeToken) TableName() string {
	return `1_bridge_token`
}

type bridgeInfo struct {
	From         string `json:"from"`
	ChainId      int64  `json:"chain_id"`
	TokenAddress string `json:"token_address"`
	TokenName    string `json:"token_name"`
	TokenSymbol  string `json:"token_symbol"`
	TokenDigits  int64  `json:"token_digits"`
}

func (b *BridgeToken) Get(ecosystem int64) (bool, error) {
	return isFound(GetDB(nil).Where("ecosystem = ? AND status = 1", ecosystem).First(b))
}

func (b *BridgeToken) GetByTokenSymbol(settingId int64, tokenAddress string) (bool, error) {
	return isFound(GetDB(nil).Where("setting_id = ? AND token_address = ? AND status = 1", settingId, tokenAddress).First(b))
}

func IsBridgeEcosystem(ecosystem int64) bool {
	if BridgeReady {
		b := &BridgeToken{}
		f, err := b.Get(ecosystem)
		if err == nil && f {
			return true
		}
	}
	return false
}

func getBridgeInfo(ecosystem int64) *bridgeInfo {
	var info bridgeInfo
	b := &BridgeToken{}
	f, err := b.Get(ecosystem)
	if err == nil && f {
		info.From = b.ChainName
		info.TokenAddress = b.TokenAddress
		info.TokenName = b.TokenName
		info.TokenSymbol = b.TokenSymbol
		info.TokenDigits = b.TokenDigits
		s := &BridgeSettings{}
		f, err = s.Get(b.SettingId)
		if err == nil && f {
			info.ChainId = s.ChainId
		}
		return &info
	}
	return nil
}
