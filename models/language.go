/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Language struct {
	ecosystem  int64
	ID         int64  `gorm:"primary_key;not null"`
	Name       string `gorm:"not null;size:100"`
	Res        string `gorm:"type:jsonb"`
	Conditions string `gorm:"not null"`
}

func (l *Language) TableName() string {
	if l.ecosystem == 0 {
		l.ecosystem = 1
	}
	return `1_languages`
}

func getLanguageValue(language, name string, ecosystem int64) string {
	var lang Language
	err := GetDB(nil).Select("res").Where("name = ? AND ecosystem = ?", name, ecosystem).First(&lang).Error
	if err != nil {
		log.WithFields(log.Fields{"err": err, "name": name, "ecosystem": ecosystem}).Warn("get Language Value Failed")
		return ""
	}
	var resMap map[string]string
	err = json.Unmarshal([]byte(lang.Res), &resMap)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "name": name, "ecosystem": ecosystem, "res": lang.Res}).Warn("get Language Value json unmarshal Failed")
		return ""
	}
	val, ok := resMap[language]
	if !ok {
		return resMap["en"]
	}
	return val
}
