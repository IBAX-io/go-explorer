/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package cmd

import (
	"context"
	"fmt"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/daemons"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/models/crontab"
	"github.com/IBAX-io/go-explorer/route"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-explorer",
	Short: "scan application",
}

func init() {
	rootCmd.AddCommand(
		initDatabaseCmd,
		startCmd,
		versionCmd,
		initRedis,
		initRedisAll,
	)
	models.InitBuildInfo()
	// This flags are visible for all child commands
	rFlag := rootCmd.PersistentFlags()
	rFlag.StringVar(&conf.GetEnvConf().ConfigPath, "config", defaultConfigPath(), "filepath to config.yml")
}

func defaultConfigPath() string {
	p, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal("getting cur wd")
	}
	return filepath.Join(p, "conf")
}

// Execute executes rootCmd command.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Executing root command")
	}
}

func loadStartRun() error {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{"panic": r, "type": consts.PanicRecoveredError}).Error("recovered panic")
			panic(r)
		}
	}()
	conf.Initer()

	exitErr := func() {
		route.SeverShutdown()
		sqldb.GormClose()
		models.GormClose()
		models.GeoIpClose()
		os.Exit(1)
	}

	daemons.StartDaemons(context.Background())
	go crontab.CreateCrontab()

	go func() {
		err := route.Run(conf.GetEnvConf().ServerInfo.Str())
		if err != nil {
			daemons.ExitCh <- fmt.Errorf("route run err:%s\n", err.Error())
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-daemons.ExitCh:
		log.WithFields(log.Fields{"err:": err}).Error("Start Daemons Failed")
		exitErr()
		return err
	case sig := <-sigChan:
		log.WithFields(log.Fields{"info:": sig}).Info("receive exit signal")
		exitErr()
		return nil
	}
}

// Load the configuration from file
func loadInitDatabase() error {
	return models.InitDatabase()
}

func loadConfigWKey(cmd *cobra.Command, args []string) {
	conf.LoadConfig(conf.GetEnvConf().ConfigPath)
}

func initRedisDatabase() error {
	return models.InitRedisDb()
}

func initRedisDatabaseAll() error {
	return models.InitRedisDbAll()
}
