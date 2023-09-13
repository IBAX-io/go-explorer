package models

type BridgeSettings struct {
	Id              int64  `gorm:"primary_key;not null"`
	BridgeAddress   string `gorm:"not null"`
	BridgeName      string `gorm:"not null"`
	ChainId         int64  `gorm:"not null"`
	ChainName       string `gorm:"not null"`
	DepositConfirm  int64  `gorm:"not null"`
	Owners          string `gorm:"not null"`
	Required        int64  `gorm:"not null"`
	Status          int    `gorm:"not null"` //1:enable 0:disable
	WithdrawConfirm int64  `gorm:"not null"`
}

var BridgeReady bool

func (b *BridgeSettings) TableName() string {
	return `1_bridge_settings`
}

func (b *BridgeSettings) Get(id int64) (bool, error) {
	return isFound(GetDB(nil).Where("id = ? AND status = 1", id).First(b))
}

func (b *BridgeSettings) GetByChainName(name string) (bool, error) {
	return isFound(GetDB(nil).Where("chain_name = ? AND status = 1", name).First(b))
}

func BridgeTableExist() bool {
	var p BridgeSettings
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}
