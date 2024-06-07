/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package conf

import (
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"os"
	"path"
	"time"

	"github.com/IBAX-io/go-explorer/storage"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

var (
	configInfo EnvConf // all server config information
)

type EnvConf struct {
	ConfigPath     string
	ServerInfo     *serverModel              `yaml:"server"`
	DatabaseInfo   *storage.DatabaseModel    `yaml:"database"`
	RedisInfo      *storage.RedisModel       `yaml:"redis"`
	Url            *UrlModel                 `yaml:"url"`
	Centrifugo     *storage.CentrifugoConfig `yaml:"centrifugo"`
	Crontab        *storage.Crontab          `yaml:"crontab"`
	CryptoSettings storage.CryptoSettings    `yaml:"crypto_settings"`
	Defi           defiInfo                  `yaml:"defi"`
}

type defiInfo struct {
	Enable    bool  `yaml:"enable"`
	Ecosystem int64 `yaml:"ecosystem"`
}

func GetEnvConf() *EnvConf {
	return &configInfo
}

func GetDbConn() *storage.DatabaseModel {
	return GetEnvConf().DatabaseInfo
}

func GetRedisDbConn() *storage.RedisModel {
	return GetEnvConf().RedisInfo
}
func GetCentrifugoConn() *storage.CentrifugoConfig {
	return GetEnvConf().Centrifugo
}

func LoadConfig(configPath string) {
	filePath := path.Join(configPath, "config.yml")
	configData, err := os.ReadFile(filePath)
	if err != nil {
		logrus.WithError(err).Fatal("config file read failed")
	}
	// expand environment variables
	configData = []byte(os.ExpandEnv(string(configData)))
	err = yaml.Unmarshal(configData, &configInfo)
	data, _ := json.Marshal(&configInfo)
	fmt.Printf("config: %v\n", string(data))
	if err != nil {
		logrus.WithError(err).Fatal("config parse failed")
	}
	registerCrypto(GetEnvConf().CryptoSettings)
}

func Initer() {
	databaseInfo := GetEnvConf().DatabaseInfo
	redisInfo := GetEnvConf().RedisInfo
	centrifugo := GetEnvConf().Centrifugo

	if err := databaseInfo.GormInit(); err != nil {
		logrus.WithError(err).Fatal("postgres database connect failed: %v", databaseInfo.Connect)
	}
	if err := redisInfo.Init(); err != nil {
		logrus.WithError(err).Fatal("redis database config information: %v", redisInfo)
	}
	if err := centrifugo.Init(); err != nil {
		logrus.WithError(err).Fatal("centrifugo config information: %v", centrifugo)
	}
	if err := initLogs(); err != nil {
		logrus.WithError(err).Fatal("init log file")
	}
}

func init() {
	time.Local = time.UTC
}

func initLogs() error {
	fileName := path.Join(GetEnvConf().ConfigPath, "logrus.log")
	openMode := os.O_APPEND
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		openMode = os.O_CREATE
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|openMode, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't open log file: ", fileName)
		return err
	}
	logrus.SetOutput(f)
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func registerCrypto(c storage.CryptoSettings) {
	crypto.InitAsymAlgo(c.Cryptoer)
	crypto.InitHashAlgo(c.Hasher)
}
