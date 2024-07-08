/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package route

import (
	"context"
	"fmt"
	"github.com/IBAX-io/go-explorer/models"
	"golang.org/x/net/http2"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/IBAX-io/go-explorer/controllers"

	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/docs"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var server *http.Server

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			//c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//OPTIONS
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
func prefix(s string) string {
	return "/api/v2/" + s
}

func Run(host string) (err error) {
	r := gin.Default()
	//Ten requests per second
	limiter := tollbooth.NewLimiter(100, nil)
	r.Use(Cors(), tollbooth_gin.LimitHandler(limiter))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": models.Version(),
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	api := r.Group(models.ApiPath)

	// programatically set swagger info
	docs.SwaggerInfo.Title = "IBAX Explorer API"
	docs.SwaggerInfo.Description = "This is ibax explorer api server."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = conf.GetEnvConf().ServerInfo.DocsApi
	//docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// use ginSwagger middleware to serve the API docs
	//api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//dashboard
	api.GET(`/websocket_token`, controllers.DashboardGetToken)
	api.GET(`/dashboard`, controllers.GetDashboard)
	api.GET(`/get_dashboard_chart`, controllers.GetDashboardChartHandler)
	api.POST(`/common_transaction_search`, controllers.CommonTransactionSearch)
	api.POST("/block_list", controllers.GetBlockList)
	api.POST(`/get_map_info`, controllers.GetMapInfo)

	api.GET("/max_block_id", controllers.GetMaxBlockId)
	api.GET(`/honor_node_list`, controllers.GetHonorNodeLists)
	api.GET(`/honor_node_map`, controllers.GetHonorNodeMapHandler)

	api.GET(`/block_tps_list`, controllers.GetBlockTpsLists)

	//Global Search
	api.GET(`/search_hash/:hash`, controllers.SearchHash)
	api.GET(`/transaction_detail/:hash`, controllers.GetTransactionDetails)
	api.GET(`/contract_tx_detail_list/:hash`, controllers.GetContractTxDetailListHandler)
	api.GET(`/transaction_utxo_detail/:hash`, controllers.GetUtxoTransactionDetails)
	api.GET(`/utxo_inputs/:hash`, controllers.GetUtxoInputsHandler)
	api.GET(`/transaction_head/:hash`, controllers.GetTransactionHead)
	api.GET(`/block_detail/:block_id`, controllers.GetBlockDetails)
	api.POST(`/account_detail`, controllers.GetAccountDetailEcosystem)
	api.GET(`/account_detail_basis/:account`, controllers.GetAccountDetailBasisEcosystem)
	api.GET(`/account_detail_basis_chart/:account`, controllers.GetAccountDetailBasisTokenChange)
	api.GET(`/account_tx_count/:ecosystem/:account`, controllers.GetAccountTxCountHandler)
	//Nft Miner Global Search
	api.GET(`/nft_miner_info/:search`, controllers.NftMinerInfoHandler)
	api.POST(`/nft_miner_history_info`, controllers.NftMinerHistoryInfoHandler)
	api.POST(`/nft_miner_stake_info`, controllers.GetNftMinerStakeInfoHandler)
	api.POST(`/nft_miner_tx_info`, controllers.GetNftMinerTxInfoHandler)

	//Block Info
	api.POST(`/block_detail_tx_list`, controllers.GetBlockDetailTxList)
	api.GET(`/block_list_chart`, controllers.GetBlockListChart)
	api.GET(`/tx_list_chart`, controllers.GetTxListChart)
	//api.POST("/transaction_list", controllers.GetNodeTransactionListHandler) TODO:delete

	//Account Detail Info
	api.POST("/ecosystem_search", controllers.EcosystemSearchHandler)
	api.POST(`/account_list`, controllers.GetAccountList)
	api.GET(`/account_list_chart/:ecosystem`, controllers.GetAccountListChartHandler)
	api.POST(`/account_detail_tx`, controllers.GetAccountTransactionHistory)
	api.POST(`/account_detail_nft_miner`, controllers.GetAccountDetailNftMinerHandler)

	//Account Detail Info Chart
	api.POST("/account_total_amount_chart", controllers.GetAccountTotalAmountChart)
	api.POST("/account_amount_change_pie_chart", controllers.GetAmountChangePieChart)
	api.POST("/account_amount_change_bar_chart", controllers.GetAmountChangeBarChart)
	api.POST("/account_tx_chart", controllers.GetAccountTxChart)

	//Node Info
	api.GET("/node_map", controllers.GetNodeMap)
	api.POST("/node_list", controllers.NodeListSearchHandler)
	api.GET("/last_dao_voting", controllers.GetLatestDaoVotingHandler)
	api.GET("/node_detail/:id", controllers.NodeDetailHandler)
	api.POST("/node_dao_detail", controllers.GetNodeDaoVoteListHandler)
	api.POST("/node_block_list", controllers.GetNodeBlockListHandler)
	api.POST("/dao_vote_list", controllers.GetDaoVoteListHandler)
	api.GET("/dao_vote_chart", controllers.GetDaoVoteChartHandler)
	api.POST("/dao_vote_detail", controllers.DaoVoteDetailSearchHandler)
	api.POST("/node_vote_history", controllers.GetNodeVoteHistoryHandler)
	api.POST("/node_staking_history", controllers.GetNodeStakingHistoryHandler)

	//Nft Miner Metaverser
	api.GET(`/nft_miner_metaverse`, controllers.GetNftMinerMetaverse)
	api.GET(`/nft_miner_map`, controllers.GetNftMinerMapHandler)
	api.POST(`/nft_miner_metaverse_list`, controllers.GetNftMinerMetaverseList)
	api.GET("/nft_miner_file/:id", controllers.GetNftMinerFileHandler)
	api.GET("/nft_miner_region", controllers.GetNftMinerRegionHandler)

	//EcoLibs
	api.GET(`/get_basis_ecosystem`, controllers.GetEcosystemBasis)
	api.POST(`/ecosystem_list`, controllers.GetEcosystemList)
	api.POST(`/get_eco_detail_info`, controllers.GetEcosystemDetailInfoHandler)
	api.POST(`/get_eco_detail_tx`, controllers.GetEcosystemDetailTxHandler)
	api.POST(`/get_eco_detail_token_symbol`, controllers.GetEcosystemDetailTokenHandler)
	api.POST(`/get_eco_detail_member`, controllers.GetEcosystemDetailMemberHandler)
	api.POST(`/platform_ecosystem_param`, controllers.GetPlatformEcosystemParam)
	api.POST(`/ecosystem_param`, controllers.GetEcosystemParam)
	api.POST(`/get_eco_database`, controllers.GetEcosystemDatabaseHandler)
	api.POST(`/get_eco_app`, controllers.GetEcosystemAppHandler)
	api.GET(`/get_eco_app_export/:id`, controllers.GetEcosystemAppExportHandler)
	api.POST(`/get_eco_attachment`, controllers.GetEcosystemAttachmentHandler)
	api.GET(`/get_eco_attachment_export/:hash`, controllers.GetEcosystemAttachmentExportHandler)

	//EcoLibs Detail Chart
	ecoChartRoute := api.Group("/eco_chart")
	ecoChartRoute.POST(`/get_account_token_change`, controllers.GetAccountTokenChangeHandler)
	ecoChartRoute.GET(`/get_circulations/:ecosystem`, controllers.GetEcosystemCirculationsChartHandler)
	ecoChartRoute.GET(`/get_has_token/:ecosystem`, controllers.GetEcoTopTenHasTokenChartHandler)
	ecoChartRoute.GET(`/get_tx_account/:ecosystem`, controllers.GetEcoTopTenTxAccountChartHandler)
	ecoChartRoute.GET(`/get_gas_combustion_pie/:ecosystem`, controllers.GetGasCombustionPieChartHandler)
	ecoChartRoute.GET(`/get_gas_combustion_line/:ecosystem`, controllers.GetGasCombustionLineChartHandler)
	ecoChartRoute.GET(`/get_tx_amount/:ecosystem`, controllers.GetEcoTxAmountChartHandler)
	ecoChartRoute.GET(`/get_gas_fee/:ecosystem`, controllers.GetEcoGasFeeChartHandler)
	ecoChartRoute.GET(`/get_new_key/:ecosystem`, controllers.GetEcoNewKeyChartHandler)
	ecoChartRoute.GET(`/get_active_key/:ecosystem`, controllers.GetEcoActiveKeyChartHandler)
	ecoChartRoute.GET(`/get_transaction/:ecosystem`, controllers.GetEcoTransactionChartHandler)
	ecoChartRoute.GET(`/get_storage_capacity/:ecosystem`, controllers.GetEcoStorageCapacityChartHandler)

	//Data Chart(Basis EcoLibs)
	//Assets Chart
	api.GET("/get_gas_fee_chart", controllers.Get15DayGasFeeChartHandler)
	api.GET("/node_contribution_chart", controllers.GetNodeContributionChartHandler)
	api.POST("/node_contribution_list", controllers.GetNodeContributionListHandler)
	api.GET("/get_new_circulations_chart", controllers.GetNewCirculationsChartHandler)
	//Address Related
	api.GET("/get_staking_account", controllers.GetStakingAccountHandler)
	api.POST("/get_account_active_chart", controllers.GetAccountActiveChartHandler)
	api.POST("/get_account_active_list", controllers.GetAccountActiveListHandler)
	api.GET("/get_new_key", controllers.GetNewKeysHandler)
	api.GET("/get_account_change_chart", controllers.GetAccountChangeChartHandler)
	//Block Related
	api.GET("/get_block_number", controllers.GetBlockNumberHandler)
	api.GET("/get_block_size_chart", controllers.GetBlockSizeChartHandler)
	api.POST("/get_block_size_list", controllers.GetBlockSizeListHandler)
	api.GET("/get_tx_chart", controllers.GetTxChartHandler)
	api.POST("/get_tx_list", controllers.GetTxListHandler)
	//NFT Miner Related
	api.GET("/nft_miner_reward", controllers.NftMinerRewardHandler)
	api.GET("/new_nft_miner", controllers.NewNftMinerHandler)
	api.GET("/nft_miner_interval", controllers.NftMinerIntervalHandler)
	api.GET("/nft_miner_interval_list", controllers.NftMinerIntervalListHandler)
	api.GET("/nft_miner_energy_power", controllers.NftMinerEnergyPowerChangeHandler)
	api.GET("/nft_miner_staked_change", controllers.NftMinerStakedChangeHandler)
	api.GET("/nft_miner_region_list", controllers.NftMinerRegionListHandler)
	//EcoLibs Related
	api.POST("/get_new_ecosystem_list", controllers.GetNewEcosystemChartListHandler)
	api.GET("/get_new_ecosystem_change", controllers.GetHistoryNewEcosystemHandler)
	api.GET("/get_ecosystem_ratio", controllers.GetTokenEcosystemRatioHandler)
	api.GET("/get_multi_fee_ecosystem", controllers.GetMultiFeeEcosystemChartHandler)
	api.GET("/get_top_ten_ecosystem_tx", controllers.GetTopTenEcosystemTxHandler)
	api.GET("/get_max_key_ecosystem", controllers.GetTopTenMaxKeysEcosystemHandler)
	api.GET("/get_govern_model", controllers.GetEcosystemGovernModelChartHandler)
	//node relatedSettings
	api.GET("/get_node_vote_change", controllers.GetNodeVoteChartHandler)
	api.GET("/get_node_staking_change", controllers.GetNodeStakingChartHandler)
	api.GET("/get_node_region", controllers.GetNodeRegionChartHandler)
	api.GET("/get_node_statistical_change", controllers.GetNodeStatisticalChangeHandler)
	api.GET("/token_price", controllers.GetTokenPriceHandler)
	api.GET("/ecosystem_logo", controllers.GetEcosystemLogoHandler)

	//common
	common := api.Group("/common")
	common.POST("/history", controllers.GetHistoryHandler)
	common.GET("/token_logo", controllers.GetTokenLogoHandler)

	api.GET(`/get_redis/:name`, controllers.GetRedisKey) //get redis keys

	api.StaticFS("/flag", http.Dir("./flag"))

	server = &http.Server{
		Addr:    host,
		Handler: r,
	}
	err = http2.ConfigureServer(server, &http2.Server{})
	if err != nil {
		log.Errorf("http2 configure Server failed :%s", err.Error())
		return err
	}
	if conf.GetEnvConf().ServerInfo.EnableHttps {
		err = server.ListenAndServeTLS(conf.GetEnvConf().ServerInfo.CertFile, conf.GetEnvConf().ServerInfo.KeyFile)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil {
		log.Errorf("server http/https start failed :%s", err.Error())
		return err
	}

	return nil
}

func SeverShutdown() {
	if server != nil {
		if err := server.Shutdown(context.Background()); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("sever shutdown failed")
		}
	}
}
