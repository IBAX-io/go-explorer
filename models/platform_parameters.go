/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"strconv"
	"sync"

	"github.com/IBAX-io/go-explorer/conf"
)

type FuelRateInfo struct {
	Map map[int64]decimal.Decimal
	sync.RWMutex
}

// PlatformParameter is model
type PlatformParameter struct {
	ID         int64  `gorm:"primary_key;not null;" json:"id"`
	Name       string `gorm:"not null;size:255" json:"name"`
	Value      string `gorm:"not null" json:"value"`
	Conditions string `gorm:"not null" json:"conditions"`
}

type SystemParameterResult struct {
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
	Rets  []PlatformParameter `json:"rets"`
}

// TableName returns name of table
func (sp PlatformParameter) TableName() string {
	return "1_platform_parameters"
}

// Get is retrieving model from database
func (sp *PlatformParameter) Get(name string) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("name = ?", name).First(sp))
}

// GetJSONField returns fields as json
func (sp *PlatformParameter) GetJSONField(jsonField string, name string) (string, error) {
	var result string
	err := GetDB(nil).Table(sp.TableName()).Where("name = ?", name).Select(jsonField).Row().Scan(&result)
	return result, err
}

// GetValueParameterByName returns value parameter by name
func (sp *PlatformParameter) GetValueParameterByName(name, value string) (*string, error) {
	var result *string
	//FROM "1_platform_parameters"
	err := GetDB(nil).Table(sp.TableName()).Raw(`SELECT value->'`+value+`' WHERE name = ?`, name).Row().Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ToMap is converting PlatformParameter to map
func (sp *PlatformParameter) ToMap() map[string]string {
	result := make(map[string]string, 0)
	result["name"] = sp.Name
	result["value"] = sp.Value
	result["conditions"] = sp.Conditions
	return result
}

func (sp *PlatformParameter) FindAppParameters(page int, size int, name, order string) (num int64, rets []PlatformParameter, err error) {
	ns := "%" + name + "%"
	if err := GetDB(nil).Table(sp.TableName()).Where("name like ?", ns).Count(&num).Error; err != nil {
		return num, rets, err
	}

	err = GetDB(nil).Table(sp.TableName()).Where("name like ?", ns).
		Order(order).Offset((page - 1) * size).Limit(size).Find(&rets).Error

	return num, rets, err
}

func GetFuelRate() (rlt map[int64]decimal.Decimal) {
	var pla PlatformParameter
	f, err := pla.Get("fuel_rate")
	if err == nil && f {
		rlt = make(map[int64]decimal.Decimal)
		var values [][]string
		err = json.Unmarshal([]byte(pla.Value), &values)
		if err == nil {
			for _, v1 := range values {
				if len(v1) == 2 {
					ecoId, _ := strconv.ParseInt(v1[0], 10, 64)
					if ecoId > 0 {
						fuelRate := v1[1]
						rlt[ecoId], _ = decimal.NewFromString(fuelRate)
					}
				}
			}
		}
	}
	return
}
