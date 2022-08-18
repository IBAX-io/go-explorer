/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	"github.com/shopspring/decimal"
)

type BlockListHeaderResponse struct {
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	Limit     int                 `json:"limit"`
	BlockInfo *BlockRet           `json:"block_info,omitempty"`
	List      []BlockListResponse `json:"list"`
}

type BlockListResponse struct {
	ID            int64  `json:"id"`
	Time          int64  `json:"time"`
	Tx            int32  `json:"tx"`
	NodePosition  int64  `json:"node_position"`
	NodeName      string `json:"node_name"`
	APIAddress    string `json:"api_address"`
	IconUrl       string `json:"icon_url"`
	Recipientid   string `json:"recipient_id"`
	GasFee        string `json:"gas_fee"`
	Reward        string `json:"reward"`
	EcoLib        int64  `json:"eco_lib"`
	ConsensusMode int32  `json:"consensus_mode"`
}

type TransactionList struct {
	Hash         string `json:"hash"`
	Time         int64  `json:"time"`
	ContractName string `json:"contract_name"`
}

type NftMinerSummaryResponse struct {
	NftMinerCount int64  `json:"nft_miner_count"`
	EnergyPower   int64  `json:"energy_power"`
	NftMinerIns   string `json:"nft_miner_ins"`
	StakeAmount   string `json:"stake_amount"`
}

type NftMinerOverviewResponse struct {
	Amount string `json:"amount"`
	Time   int64  `json:"time"`
}

type NftMinerFindRequest struct {
	Keyid string `json:"keyid"`
}

type TotalResult struct {
	Total int64 `json:"total"`
}

type CommonResult struct {
	TotalResult
	IsCreate bool                `json:"is_create"`
	Rets     []map[string]string `json:"rets"`
}

type AccountNftMinerListResult struct {
	TotalResult
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
	Rets  []AccountNftMinerList `json:"rets"`
}

type AccountNftMinerList struct {
	ID          int64  `json:"id"`
	EnergyPoint int    `json:"energy_point"`
	Hash        string `json:"hash"` //hash
	Time        int64  `json:"time"` //create time
	StakeAmount int64  `json:"stake_amount"`
	Cycle       int64  `json:"cycle"`
	Ins         string `json:"ins"`
}

type NftMinerInfoRequest struct {
	Search any    `json:"search"` //NFT Miner ID OR NFT Miner HASH
	Order  string `json:"order"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type NftMinerInfoResponse struct {
	ID          int64           `json:"id"`   //NFT Miner ID
	Hash        string          `json:"hash"` //NFT Miner hash
	EnergyPoint int             `json:"energy_point"`
	StakeCount  int64           `json:"stake_count"`
	StakeAmount int64           `json:"stake_amount"` //starking
	Owner       string          `json:"owner"`        //owner account
	Creator     string          `json:"creator"`      //owner account
	RewardCount int64           `json:"reward_count"`
	Cycle       int64           `json:"cycle"`
	DateCreated int64           `json:"date_created"` //create time
	Ins         string          `json:"ins"`
	EnergyPower decimal.Decimal `json:"energy_power"`
}

type NftMinerTxInfoResponse struct {
	ID      int64  `json:"id"`
	NftId   int64  `json:"nft_id"`
	BlockId int64  `json:"block_id"`
	Time    int64  `json:"time"`
	Ins     string `json:"ins"`
	Account string `json:"account"`
}

type NftMinerHistoryInfoResponse struct {
	ID       int64  `json:"id"`
	NftId    int64  `json:"nft_id"`
	TxHash   string `json:"tx_hash"`
	NftHash  string `json:"nft_hash"`
	Time     int64  `json:"time"`
	Events   string `json:"events"`
	Contract string `json:"contract"`
	Source   string `json:"source"`
}

type NftMinerListResponse struct {
	ID          int64           `json:"id"` //NFT Miner ID
	Hash        string          `json:"hash"`
	EnergyPoint int             `json:"energy_point"`
	Owner       string          `json:"owner"`
	StakeAmount int64           `json:"stake_amount"`
	EnergyPower decimal.Decimal `json:"energy_power"`
	Time        int64           `json:"time"`
	Ins         string          `json:"ins"`
}

type NftMinerCoordinate struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type NftMinerMetaverseInfoResponse struct {
	Count        int64           `json:"count"`
	BlockReward  float64         `json:"block_reward"`
	HalveNumber  int64           `json:"halve_number"`
	StakeAmounts string          `json:"stake_amounts"`
	EnergyPower  decimal.Decimal `json:"energy_power"`
	RewardAmount string          `json:"reward_amount"`
	StakedCount  int64           `json:"staked_count"`
	StakingCount int64           `json:"staking_count"`
	Region       int64           `json:"region"`
}

type NftMinerStakeInfoResponse struct {
	ID          int64  `json:"id"`
	NftId       int64  `json:"nft_id"`
	TxHash      string `json:"tx_hash"`
	BlockId     int64  `json:"block_id"`
	Time        int64  `json:"time"`
	StakeAmount int64  `json:"stake_amount"` //starking
	Cycle       int64  `json:"cycle"`
	EnergyPower string `json:"energy_power"`
	StakeStatus bool   `json:"stake_status"`
}

type EcosystemDetailInfoResponse struct {
	EcosystemId int64  `json:"ecosystem_id"`
	Ecosystem   string `json:"ecosystem"`
	//Logo            string `json:"logo"`
	LogoHash        string `json:"logo_hash"`
	BlockId         int64  `json:"block_id"`
	Hash            string `json:"hash"`
	Creator         string `json:"creator"`
	EcoType         int    `json:"eco_type"`
	EcoTag          int    `json:"eco_tag"`
	EcoCascade      int    `json:"eco_cascade"`
	MultiFee        bool   `json:"multi_fee"`
	EcoIntroduction string `json:"eco_introduction"`
	Time            int64  `json:"time"`
	GovernModel     int64  `json:"govern_model"`
	GovernCommittee string `json:"govern_committee"` //TODO: need add
	TokenSymbol     string `json:"token_symbol"`
	TotalAmount     string `json:"total_amount"`
	IsWithdraw      bool   `json:"is_withdraw"`
	Withdraw        string `json:"withdraw"`
	IsEmission      bool   `json:"is_emission"`
	Emission        string `json:"emission"`
	FeeModel        int    `json:"fee_model"`
	FeeModeAccount  string `json:"fee_mode_account"`

	FeeModelVmcost  sqldb.FeeModeFlag `json:"vmcost"`
	FeeModeStorage  sqldb.FeeModeFlag `json:"storage"`
	FeeModeElement  sqldb.FeeModeFlag `json:"element"`
	FeeModeExpedite sqldb.FeeModeFlag `json:"expedite"`

	WithholdingMode   int     `json:"withholding_mode"`
	IsCombustion      bool    `json:"is_combustion"`
	Combustion        string  `json:"combustion"`
	CombustionPercent int64   `json:"combustion_percent"`
	Circulations      string  `json:"circulations"`
	FollowFuel        float64 `json:"follow_fuel"`

	Registered       int               `json:"registered"`
	Country          int               `json:"country"`
	RegistrationNo   string            `json:"registration_no"`
	RegistrationType int               `json:"registration_type"`
	Social           map[string]string `json:"social"`
}

type EcosystemTxList struct {
	Hash     string `json:"hash"`
	BlockId  int64  `json:"block_id"`
	Time     int64  `json:"time"`
	Contract string `json:"contract"`
	Address  string `json:"address"`
	Status   int64  `json:"status"`
}

type EcosystemSearchResponse struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

type GeneralResponse struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	List  any   `json:"list"`
}

type EcosystemTokenSymbolList struct {
	Id             int64           `json:"id"`
	Account        string          `json:"account"`
	AccountName    string          `json:"account_name"`
	Amount         string          `json:"amount"`
	AccountedFor   decimal.Decimal `json:"accounted_for"`
	TokenSymbol    string          `json:"token_symbol"`
	FrontCommittee bool            `json:"front_committee"`
	Committee      bool            `json:"committee"`
}

type EcosystemMemberList struct {
	Id          int64  `json:"id"`
	Account     string `json:"account"`
	AccountName string `json:"account_name"`
	RolesName   string `json:"roles_name"`
	JoinTime    int64  `json:"join_time"`
}

type EcosystemAppInfo struct {
	AppId     int64             `json:"app_id"`
	Contracts []ContractsParams `json:"contracts,omitempty"`
	Page      []PageParams      `json:"page,omitempty"`
	Snippets  []SnippetsParams  `json:"snippets,omitempty"`
	Table     []TableParams     `json:"table,omitempty"`
	Params    []AppParams       `json:"params,omitempty"`
}

type EcosystemAppList struct {
	AppId      int64  `json:"app_id"`
	Name       string `json:"name"`
	Conditions string `json:"conditions"`
}

type AppInfo struct {
	Contracts ContractsParams `json:"contracts"`
	Page      PageParams      `json:"page"`
	Snippets  SnippetsParams  `json:"snippets"`
	Table     TableParams     `json:"table"`
	Params    AppParams       `json:"params"`
	Menu      MenuParams      `json:"menu"`
	Languages LanguagesParams `json:"languages"`
}

type ExportAppInfo struct {
	Conditions string `json:"conditions"`
	Name       string `json:"name"`
	Data       []any  `json:"data"`
}

type EcosystemTxRatioChart struct {
	Value float64 `json:"value"`
	Name  string  `json:"name"`
}

type EcosystemAttachmentResponse struct {
	Name     string `json:"name"`
	Id       int64  `json:"id"`
	Hash     string `json:"hash"`
	MimeType string `json:"mime_type"`
}

type TxDetailedInfoHeadResponse struct {
	BlockID      int64  `json:"block_id"`
	Hash         string `json:"hash"`
	ContractName string `json:"contract_name"`
	ContractCode string `json:"contract_code"`
	Params       string `json:"params"`
	Time         int64  `json:"time"`

	EcosystemName string `json:"ecosystem_name"`
	Ecosystem     int64  `json:"ecosystem"`
	LogoHash      string `json:"logo_hash"`
	TokenSymbol   string `json:"token_symbol"`
	Address       string `json:"address"`
	Size          string `json:"size"`
}

type AccountRatio struct {
	Account      string          `json:"account"`
	Amount       string          `json:"amount"`
	AccountedFor decimal.Decimal `json:"accounted_for"`
	StakeAmount  string          `json:"stake_amount,omitempty"`
}

//EcoSystem Chart response
type EcoTopTenHasTokenResponse struct {
	TokenSymbol string         `json:"token_symbol"`
	Name        string         `json:"name"` //ecosystem name
	List        []AccountRatio `json:"list"`
}

type EcoTopTenTxAmountResponse struct {
	EcoTopTenHasTokenResponse
}

type EcoCirculationsResponse struct {
	Circulations     string `json:"circulations"`
	StakeAmount      string `json:"stake_amount"`
	FreezeAmount     string `json:"freeze_amount"` //todo:need add
	NftBalanceSupply string `json:"nft_balance_supply"`
	BurningTokens    string `json:"burning_tokens"`
	Combustion       string `json:"combustion"`
	TokenSymbol      string `json:"token_symbol"`
	Name             string `json:"name"`
	SupplyToken      string `json:"supply_token"`
	Emission         string `json:"emission"`

	Change EcoCirculationsChangeResponse `json:"change"`
}

type EcoCirculationsChangeResponse struct {
	Time             []string `json:"time"`
	Circulations     []string `json:"circulations"`
	StakeAmount      []string `json:"stake_amount"`
	FreezeAmount     []string `json:"freeze_amount"` //todo:need add
	NftBalanceSupply []string `json:"nft_balance_supply"`
	BurningTokens    []string `json:"burning_tokens"`
	Combustion       []string `json:"combustion"`
	SupplyToken      []string `json:"supply_token"`
	Emission         []string `json:"emission"`
}

type EcoGasFeeResponse struct {
	GasFee      string `json:"gas_fee"`
	Combustion  string `json:"combustion"`
	TokenSymbol string `json:"token_symbol"`
	Name        string `json:"name"` //ecosystem name
}

type EcoGasFeeChangeResponse struct {
	TokenSymbol string   `json:"token_symbol"`
	Name        string   `json:"name"` //ecosystem name
	Time        []string `json:"time"`
	GasFee      []string `json:"gas_fee"`
	Combustion  []string `json:"combustion"`
}

type EcoTxAmountDiffResponse struct {
	TokenSymbol string   `json:"token_symbol"`
	Name        string   `json:"name"` //ecosystem name
	Time        []int64  `json:"time"`
	Amount      []string `json:"amount"`
}

type EcoTxGasFeeDiffResponse struct {
	TokenSymbol    string   `json:"token_symbol"`
	Name           string   `json:"name"` //ecosystem name
	Time           []int64  `json:"time"`
	EcoGasAmount   []string `json:"eco_gas_amount"`
	BasisGasAmount []string `json:"basis_gas_amount"`
}

type AccountTotalAmountChart struct {
	TotalAmount       decimal.Decimal `json:"total_amount"`
	Amount            decimal.Decimal `json:"amount" gorm:"column:amount"`
	AmountRatio       float64         `json:"amount_ratio"`
	StakeAmount       decimal.Decimal `json:"stake_amount" gorm:"column:stake_amount"`
	StakeAmountRatio  float64         `json:"stake_amount_ratio"`
	FreezeAmount      decimal.Decimal `json:"freeze_amount" gorm:"column:freeze_amount"` //TODO: need add
	FreezeAmountRatio float64         `json:"freeze_amount_ratio"`
	TokenSymbol       string          `json:"token_symbol"`
}

type AccountAmountChangePieChart struct {
	Outcome     decimal.Decimal `json:"outcome" gorm:"column:outcome"`
	Income      decimal.Decimal `json:"income" gorm:"column:income"`
	TokenSymbol string          `json:"token_symbol" gorm:"column:token_symbol"`
}

type AccountAmountChangeBarChart struct {
	TokenSymbol string   `json:"token_symbol"`
	Name        string   `json:"name"`
	Time        []int64  `json:"time"`
	Outcome     []string `json:"outcome"`
	Income      []string `json:"income"`
	Balance     []string `json:"balance"`
}

type AccountTxChart struct {
	Time []int64 `json:"time"`
	Tx   []int64 `json:"tx"`
}

//data chart
type GasFeeChangeResponse struct {
	Time   []int64  `json:"time"`
	GasFee []string `json:"gas_fee"`
}

type HonorNodeListResponse struct {
	NodeName        string          `json:"node_name"`
	KeyID           string          `json:"key_id"`
	City            string          `json:"city"`
	IconUrl         string          `json:"icon_url"`
	GasFee          string          `json:"gas_fee"`
	NodeBlocks      int64           `json:"node_blocks"`
	PkgAccountedFor decimal.Decimal `json:"pkg_accounted_for"`
	NodePosition    int64           `json:"node_position"`
}

type HonorNodeChartResponse struct {
	Total     int64                   `json:"total"`
	Page      int                     `json:"page"`
	Limit     int                     `json:"limit"`
	List      []HonorNodeListResponse `json:"list"`
	Name      []string                `json:"name"`
	NodeBlock []int64                 `json:"node_block"`
}

type CirculationsChangeResponse struct {
	Time         []int64  `json:"time"`
	Circulations []string `json:"circulations"`
	FreezeAmount []string `json:"freeze_amount"` //todo:need add
}

type CirculationsChartResponse struct {
	TotalCirculations string                     `json:"total_circulations"`
	Circulations      string                     `json:"circulations"`
	FreezeAmount      string                     `json:"freeze_amount"` //todo:need add
	Change            CirculationsChangeResponse `json:"change"`
}

type AccountChangeChartResponse struct {
	NowTotal        int64             `json:"now_total"`
	MaxAccountedFor decimal.Decimal   `json:"max_accounted_for"`
	MaxTime         int64             `json:"max_time"`
	MinAccountedFor decimal.Decimal   `json:"min_accounted_for"`
	MinTime         int64             `json:"min_time"`
	Time            []int64           `json:"time"`
	Total           []int64           `json:"total"`
	HasToken        []int64           `json:"has_token"`
	AccountedFor    []decimal.Decimal `json:"accounted_for"`
}

type ActiveReportInfo struct {
	Time          int64  `gorm:"column:time;not null"`
	ActiveAccount int64  `gorm:"column:active;not null"`
	Ratio         string `gorm:"column:ratio;type:varchar(30);not null"`
	RelativeRatio string `gorm:"column:relative_ratio;type:varchar(30);not null"`
	TxNumber      int64  `gorm:"column:tx_number;not null"`
	TxAmount      string `gorm:"column:tx_amount;type:decimal(30);not null"`
}

type DailyActiveChartResponse struct {
	MaxActive int64               `json:"max_active"`
	MaxTime   int64               `json:"max_time"`
	MinActive int64               `json:"min_active"`
	MinTime   int64               `json:"min_time"`
	Time      []int64             `json:"time"`
	Info      []DailyActiveReport `json:"info"`
}

type StakingAccountResponse struct {
	Account      string          `json:"account" gorm:"column:account"`
	StakeAmount  decimal.Decimal `json:"stake_amount" gorm:"column:stake_amount"`
	AccountedFor decimal.Decimal `json:"accounted_for"`
}

type BlockSizeListResponse struct {
	Id   int64 `json:"id"`
	Time int64 `json:"time"`
	Size int64 `json:"size"`
	Tx   int64 `json:"tx"`
}

type TransactionListResponse struct {
	Hash    string `json:"hash"`
	Block   int64  `json:"block"`
	Time    int64  `json:"time"`
	Address int64  `json:"address"`
	Name    string `json:"name"`
}

type DaysNumberResponse struct {
	Time   []int64 `json:"time"`
	Number []int64 `json:"number"`
}

type DaysAmountResponse struct {
	Time   []int64  `json:"time"`
	Amount []string `json:"amount"`
}

type NftMinerIntervalResponse struct {
	Time            string `json:"time,omitempty"`
	OneToTen        int64  `json:"one_to_ten"`
	TenToTwenty     int64  `json:"ten_to_twenty"`
	TwentyToThirty  int64  `json:"twenty_to_thirty"`
	ThirtyToForty   int64  `json:"thirty_to_forty"`
	FortyToFifty    int64  `json:"forty_to_fifty"`
	FiftyToSixty    int64  `json:"fifty_to_sixty"`
	SixtyToSeventy  int64  `json:"sixty_to_seventy"`
	SeventyToEighty int64  `json:"seventy_to_eighty"`
	EightyToNinety  int64  `json:"eighty_to_ninety"`
	NinetyToHundred int64  `json:"ninety_to_hundred"`
}

type NftMinerStakingChangeResponse struct {
	Time        []string `json:"time"`
	StakeAmount []int64  `json:"stake_amount"`
	Number      []int64  `json:"number"`
}

type NftMinerEnergyPowerChangeResponse struct {
	Time        []string `json:"time"`
	EnergyPower []string `json:"energy_power"`
}

type TokenEcosystemResponse struct {
	Emission        int64           `json:"emission"`
	EmissionRatio   decimal.Decimal `json:"emission_ratio"`
	UnEmission      int64           `json:"un_emission"`
	UnEmissionRatio decimal.Decimal `json:"un_emission_ratio"`
}

type EcosystemTxRatioResponse struct {
	Name  string          `json:"name"`
	Tx    int64           `json:"tx"`
	Ratio decimal.Decimal `json:"ratio"`
}

type EcosystemListResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	LogoHash string `json:"logo_hash"`
	Contract int64  `json:"contract"`
	Block    int64  `json:"block"`
	Hash     string `json:"hash"`
}

type EcosystemKeysRatioResponse struct {
	Id     int64           `json:"id"`
	Name   string          `json:"name"`
	Number int64           `json:"number"`
	Ratio  decimal.Decimal `json:"ratio"`
}

type MultiFeeEcosystemRatioResponse struct {
	MultiFee        int64           `json:"multi_fee"`
	MultiFeeRatio   decimal.Decimal `json:"multi_fee_ratio"`
	UnMultiFee      int64           `json:"un_multi_fee"`
	UnMultiFeeRatio decimal.Decimal `json:"un_multi_fee_ratio"`
}

type NewNftChangeChartResponse struct {
	Time  []int64           `json:"time"`
	New   []int64           `json:"new"`
	Stake []int64           `json:"stake"`
	Ratio []decimal.Decimal `json:"ratio"`
}

type GovernModelRatioResponse struct {
	DAOGovernance int64           `json:"dao_governance"`
	DAORatio      decimal.Decimal `json:"dao_ratio"`
	CreatorModel  int64           `json:"creator_model"`
	CreatorRatio  decimal.Decimal `json:"creator_ratio"`
}

type SearchHashResponse struct {
	IsTxHash bool `json:"is_tx_hash"`
}

type StorageCapacitysChart struct {
	Name             string   `json:"name,omitempty"` //ecosystem name
	Time             []int64  `json:"time"`
	StorageCapacitys []string `json:"storage_capacitys"`
}

type Positioning struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
	Val string `json:"val"` //value
}

//Node Related
type NodeMapResponse struct {
	CandidateTotal int64         `json:"candidate_total"`
	HonorTotal     int64         `json:"honor_total"`
	NodeList       []Positioning `json:"node_list"`
}

type NodeListResponse struct {
	Ranking        int64  `json:"ranking"`
	Id             int64  `json:"id"`
	IconUrl        string `json:"icon_url"`
	Name           string `json:"name"`
	Website        string `json:"website"`
	ApiAddress     string `json:"api_address"`
	Address        string `json:"address"`
	Packed         int64  `json:"packed"`
	PackedRate     string `json:"packed_rate"`
	Vote           string `json:"vote"`
	VoteRate       string `json:"vote_rate"`
	VoteTrend      int    `json:"vote_trend"` //0:unknown 1:Up 2:Down 3:Equal
	Staking        string `json:"staking"`
	FrontCommittee bool   `json:"front_committee"`
	Committee      bool   `json:"committee"`
}

type VotingResponse struct {
	TotalResult
	VoteInfo
}

type VoteInfo struct {
	Id           int64   `json:"id"`
	Title        string  `json:"title"`
	Created      int64   `json:"created"`
	VotedRate    int     `json:"voted_rate"`
	ResultRate   float64 `json:"result_rate"`
	RejectedRate float64 `json:"rejected_rate"`
}

type NodeDetailResponse struct {
	NodeListResponse
	StakeRate string `json:"stake_rate"`
	Account   string `json:"account"`
}

type NodeVoteResponse struct {
	VoteInfo
	Result int `json:"result"` //vote result
	Status int `json:"status"` //node vote status
}

type GasFee struct {
	Amount      string `json:"amount"`
	TokenSymbol string `json:"token_symbol"`
}

type NodeBlockListResponse struct {
	BlockId   int64  `json:"block_id"`
	Time      int64  `json:"time"`
	Tx        int32  `json:"tx"`
	EcoNumber int    `json:"eco_number"`
	GasFee1   GasFee `json:"gas_fee_1"` //IBXC
	GasFee2   GasFee `json:"gas_fee_2"`
	GasFee3   GasFee `json:"gas_fee_3"`
	GasFee4   GasFee `json:"gas_fee_4"`
	GasFee5   GasFee `json:"gas_fee_5"`
}

type VoteStatus struct {
	Agree      int64 `json:"agree"`
	Rejected   int64 `json:"rejected"`
	DidNotVote int64 `json:"did_not_vote"`
}

type DaoVoteChartResponse struct {
	NotEnoughVotes int64 `json:"not_enough_votes"`
	Rejected       int64 `json:"rejected"`
	No             int64 `json:"no"`
	Accept         int64 `json:"accept"`
	List           []struct {
		VoteStatus
		Total int64 `json:"total"`
		Time  int64 `json:"time"`
	} `json:"list"`
}

type DaoVoteDetailResponse struct {
	VoteStatus
	VoteInfo
	Progress int `json:"progress"`

	TypeDecision string `json:"type_decision"`
	Description  string `json:"description"`
	//DescriptionText string `json:"description_text"`
	VotingSubject struct {
		ContractAccept   string `json:"contract_accept"`
		Arguments        string `json:"arguments"`
		ContractOfReject string `json:"contract_of_reject"`
	} `json:"voting_subject"`

	Voting struct {
		Type            string   `json:"type"`
		Status          string   `json:"status"`
		CountTypeVoters int      `json:"count_type_voters"`
		VoteCountType   string   `json:"vote_count_type"`
		Filled          string   `json:"filled"` //full data
		Decision        string   `json:"decision"`
		DateStart       int64    `json:"date_start"`
		DateEnd         int64    `json:"date_end"`
		Quorum          int      `json:"quorum"`
		Volume          int      `json:"volume"`
		Participants    int64    `json:"participants"`
		Creator         string   `json:"creator"`
		Member          []string `json:"member"`
	} `json:"voting"`
}

type NodeVoteChangeResponse struct {
	Time []int64  `json:"time"`
	Vote []string `json:"vote"`
}

type NodeStakingChangeResponse struct {
	Time    []string `json:"time"`
	Staking []string `json:"staking"`
}

type NodeChangeResponse struct {
	HonorNodeChart
	NewHonor []string `json:"new_honor"` //need add
}

type RegionChangeResponse struct {
	Time []any   `json:"time"`
	List [][]any `json:"list"`
}

type NodeVoteHistoryResponse struct {
	Vote int64 `json:"vote"`
	NodeStakingHistoryResponse
}

type NodeStakingHistoryResponse struct {
	Id      int64  `json:"id"`
	TxHash  string `json:"tx_hash"`
	Address string `json:"address"`
	Time    int64  `json:"time"`
	Events  int64  `json:"events"`
	Amount  string `json:"amount"`
}

type EcoTxListResponse struct {
	LogTransaction
}

type NodePositionInfo struct {
	Id            int64  `json:"id"`
	NodeName      string `json:"node_name"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	NodePosition  string `json:"node_position"`
	Block         string `json:"block"`
	NodeBlock     string `json:"node_block"`
	ConsensusMode int32  `json:"consensus_mode"`
	ApiAddress    string `json:"api_address"`
}

type HonorNodeMapResponse struct {
	//LatestBlock         string `json:"latest_block"`
	//LatestNodePosition  string `json:"latest_node_position"`
	//LatestConsensusMode int32  `json:"latest_consensus_mode"`
	NodePositionInfo
	List []NodePositionInfo `json:"list"`
}

type NftMinerRegionListResponse struct {
	Region        string `json:"region"`
	Total         int64  `json:"total"`
	StakingNumber int64  `json:"staking_number"`
	StakingAmount string `json:"staking_amount"`
}
