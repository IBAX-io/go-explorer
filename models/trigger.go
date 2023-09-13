package models

import (
	"strings"
)

func InitDBTrigger() (err error) {
	err = GetDB(nil).Exec(getTriggerScript()).Error
	return
}

func getTriggerScript() string {
	var scripts []string
	if !AccountDetailTableExist() {
		scripts = append(scripts, CreateAccountDetailTable)
	}
	scripts = append(scripts, AccountDetailSQL)
	return strings.Join(scripts, "\r\n")
}
