/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initDatabaseCmd represents the initDatabase command
var initDatabaseCmd = &cobra.Command{
	Use:    "initDatabase",
	Short:  "Initializing database",
	PreRun: loadConfigWKey,
	Run: func(cmd *cobra.Command, args []string) {
		if err := loadInitDatabase(); err != nil {
			log.WithError(err).Fatal("init db")
		}
		log.Info("initDatabase completed")
	},
}

var initRedis = &cobra.Command{
	Use:    "initRedis",
	Short:  "Initializing redis database",
	PreRun: loadConfigWKey,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initRedisDatabase(); err != nil {
			log.WithError(err).Fatal("init redis db")
		}
		log.Info("init redis database completed")
	},
}

var initRedisAll = &cobra.Command{
	Use:    "initRedisAll",
	Short:  "Initializing redis database all",
	PreRun: loadConfigWKey,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initRedisDatabaseAll(); err != nil {
			log.WithError(err).Fatal("init redis db all")
		}
		log.Info("init redis database all completed")
	},
}
