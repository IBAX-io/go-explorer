/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/block"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/transaction"
	"github.com/IBAX-io/go-ibax/packages/types"
)

// blocks is storing block data
type blocks struct {
	Header            types.BlockData
	PrevHeader        *types.BlockData
	PrevRollbacksHash []byte
	MrklRoot          []byte
	BinData           []byte
	Transactions      []*transaction.Transaction
	SysUpdate         bool
	GenBlock          bool // it equals true when we are generating a new block
	Notifications     []types.Notifications
}

// InfoBlock is model
type InfoBlock struct {
	Hash           []byte `gorm:"not null"`
	EcosystemID    int64  `gorm:"not null default 0"`
	KeyID          int64  `gorm:"not null default 0"`
	NodePosition   string `gorm:"not null default 0"`
	BlockID        int64  `gorm:"not null"`
	Time           int64  `gorm:"not null"`
	CurrentVersion string `gorm:"not null"`
	Sent           int8   `gorm:"not null"`
	RollbacksHash  []byte `gorm:"not null"`
	ConsensusMode  int32  `gorm:"not null"`
}

// TableName returns name of table
func (ib *InfoBlock) TableName() string {
	return "info_block"
}

// Get is retrieving model from database
func (ib *InfoBlock) Get() (bool, error) {
	return isFound(conf.GetDbConn().Conn().Last(ib))
}

func DealTransactionBlockTxDetial(mc *Block) (int64, *BlockDetailedInfoHex, *[]BlockTxDetailedInfoHex, error) {
	var (
		ret []BlockTxDetailedInfoHex
		ts  int64
	)

	//bk := &Block{}
	rt, err := GetBlocksDetailedInfoHexByScanOut(mc)
	if err == nil {
		for j := 0; j < int(rt.Tx); j++ {
			bh := BlockTxDetailedInfoHex{}
			bh.BlockID = rt.Header.BlockId
			bh.ContractName = rt.Transactions[j].ContractName
			bh.Hash = rt.Transactions[j].Hash
			bh.KeyID = rt.Transactions[j].KeyID
			//bh.Params = rt.Transactions[j].Params
			bh.Time = MsToSeconds(rt.Transactions[j].Time)
			bh.Type = rt.Transactions[j].Type
			if bh.Time == 0 {
				bh.Time = MsToSeconds(rt.Time)
			}
			if bh.KeyID == "" {
				bh.KeyID = rt.KeyID
			}
			if bh.Ecosystem == 0 {
				bh.Ecosystem = 1
			}
			bh.Ecosystemname = rt.Transactions[j].Ecosystemname
			bh.Ecosystem = rt.Transactions[j].Ecosystem
			if bh.Ecosystem == 1 {
				bh.Token_symbol = SysTokenSymbol
				if bh.Ecosystemname == "" {
					bh.Ecosystemname = SysEcosystemName
				}
			} else {
				bh.Token_symbol = rt.Transactions[j].TokenSymbol
			}
			//dlen := unsafe.Sizeof(rt.Transactions[j])
			bh.Size = rt.Transactions[j].Size
			ts += bh.Size
			//lg1, err := json.Marshal(rt.Transactions[j].Params)
			//if err == nil {
			//	bh.Params = string(lg1)
			//}

			ret = append(ret, bh)
		}
	}

	return ts, rt, &ret, err
}

func GetTransactionTxDetial(mc *Block) (int64, *BlockDetailedInfoHex, error) {
	var (
		ts int64
	)

	//bk := &Block{}
	rt, err := GetBlocksDetailedInfoHexByScanOut(mc)
	if err == nil {
		for j := 0; j < int(rt.Tx); j++ {
			bh := BlockTxDetailedInfoHex{}
			bh.Size = rt.Transactions[j].Size
			ts += bh.Size
		}
	}

	return ts, rt, err
}

func GetBlocksDetailedInfoHexByScanOut(mc *Block) (*BlockDetailedInfoHex, error) {
	var (
		transize int64
	)
	result := BlockDetailedInfoHex{}
	blck, err := block.UnmarshallBlock(bytes.NewBuffer(mc.Data))
	if err != nil {
		return &result, err
	}

	txDetailedInfoCollection := make([]TxDetailedInfoHex, 0, len(blck.Transactions))
	for _, tx := range blck.Transactions {
		txDetailedInfo := TxDetailedInfoHex{
			Hash: hex.EncodeToString(tx.Hash()),
		}

		if tx.IsSmartContract() {
			if tx.SmartContract().TxSmart.UTXO != nil {
				txDetailedInfo.ContractName = UtxoTx
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.UTXO)
				txDetailedInfo.Params = string(dataBytes)
			} else if tx.SmartContract().TxSmart.TransferSelf != nil {
				txDetailedInfo.ContractName = UtxoTransferSelf
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.TransferSelf)
				txDetailedInfo.Params = string(dataBytes)
			} else {
				txDetailedInfo.ContractName, txDetailedInfo.Params = GetMineParam(tx.SmartContract().TxSmart.EcosystemID, tx.SmartContract().TxContract.Name, tx.SmartContract().TxData, tx.Hash())
			}
			txDetailedInfo.KeyID = converter.AddressToString(tx.KeyID())
			txDetailedInfo.Time = MsToSeconds(tx.Timestamp())
			txDetailedInfo.Type = int64(tx.Type())
			txDetailedInfo.Size = int64(len(tx.FullData))
			transize += txDetailedInfo.Size
		}

		if txDetailedInfo.Time == 0 {
			txDetailedInfo.Time = MsToSeconds(mc.Time)
		}
		if txDetailedInfo.KeyID == "" {
			txDetailedInfo.KeyID = converter.AddressToString(mc.KeyID)
		}

		if tx.IsSmartContract() {
			txDetailedInfo.Ecosystem = tx.SmartContract().TxSmart.EcosystemID
			if txDetailedInfo.Ecosystem == 0 {
				txDetailedInfo.Ecosystem = 1
			}
			txDetailedInfo.TokenSymbol, txDetailedInfo.Ecosystemname = Tokens.Get(txDetailedInfo.Ecosystem), EcoNames.Get(txDetailedInfo.Ecosystem)

		} else {
			if txDetailedInfo.Ecosystem == 0 {
				txDetailedInfo.Ecosystem = 1
			}
			if txDetailedInfo.Ecosystem == 1 {
				txDetailedInfo.TokenSymbol = SysTokenSymbol
				if txDetailedInfo.Ecosystemname == "" {
					txDetailedInfo.Ecosystemname = SysEcosystemName
				}
			}
		}

		txDetailedInfoCollection = append(txDetailedInfoCollection, txDetailedInfo)

		//log.WithFields(log.Fields{"block_id": blockModel.ID, "tx hash": txDetailedInfo.Hash, "contract_name": txDetailedInfo.ContractName, "key_id": txDetailedInfo.KeyID, "time": txDetailedInfo.Time, "type": txDetailedInfo.Type, "params": txDetailedInfoCollection}).Debug("Block Transactions Information")
	}

	header := BlockHeaderInfoHex{
		BlockId:      blck.Header.BlockId,
		Time:         blck.Header.Timestamp,
		EcosystemId:  blck.Header.EcosystemId,
		KeyId:        converter.AddressToString(blck.Header.KeyId),
		NodePosition: blck.Header.NodePosition,
		Sign:         hex.EncodeToString(blck.Header.Sign),
		BlockHash:    hex.EncodeToString(blck.Header.BlockHash),
		Version:      blck.Header.Version,
	}
	//prehash
	if mc.ID > 1 {
		var bk Block
		pfb, err := bk.GetId(mc.ID - 1)
		if err != nil {
			return &result, err
		}
		if pfb {
			header.PreHash = hex.EncodeToString(bk.Hash)
		}

	}
	if header.EcosystemId == 0 {
		header.EcosystemId = 1
	}

	bdi := BlockDetailedInfoHex{
		Header:        header,
		Hash:          hex.EncodeToString(mc.Hash),
		EcosystemID:   mc.EcosystemID,
		NodePosition:  mc.NodePosition,
		KeyID:         converter.AddressToString(mc.KeyID),
		Time:          mc.Time,
		Tx:            mc.Tx,
		RollbacksHash: hex.EncodeToString(mc.RollbacksHash),
		MerkleRoot:    hex.EncodeToString(blck.MerkleRoot),
		BinData:       hex.EncodeToString(blck.BinData),
		SysUpdate:     blck.SysUpdate,
		GenBlock:      blck.GenBlock,
		//StopCount:     blck.s,
		BlockSize:    ToCapacityString(int64(len(mc.Data))),
		TxTotalSize:  ToCapacityString(transize),
		Transactions: txDetailedInfoCollection,
	}

	if bdi.EcosystemID == 0 {
		bdi.EcosystemID = 1
	}
	return &bdi, nil
}

func GetBlocksTransactionInfoByBlockInfo(mc *Block) (*BlockDetailedInfoHex, error) {
	var (
		transize int64
	)
	result := BlockDetailedInfoHex{}
	blck, err := block.UnmarshallBlock(bytes.NewBuffer(mc.Data))
	if err != nil {
		return &result, err
	}

	for _, tx := range blck.Transactions {
		if tx.IsSmartContract() {
			transize += int64(len(tx.FullData))
		}
	}

	header := BlockHeaderInfoHex{
		BlockId:       blck.Header.BlockId,
		Time:          blck.Header.Timestamp,
		EcosystemId:   blck.Header.EcosystemId,
		KeyId:         converter.AddressToString(blck.Header.KeyId),
		NodePosition:  blck.Header.NodePosition,
		Sign:          hex.EncodeToString(blck.Header.Sign),
		BlockHash:     hex.EncodeToString(blck.Header.BlockHash),
		Version:       blck.Header.Version,
		ConsensusMode: blck.Header.ConsensusMode,
	}
	for _, value := range HonorNodes {
		if value.NodePosition == header.NodePosition && value.ConsensusMode == header.ConsensusMode {
			header.IconUrl = value.IconUrl
			header.Address = value.Country + "-" + value.City
			break
		}
	}
	if header.IconUrl == "" {
		header.IconUrl = conf.GetEnvConf().Url.Base + "default.png"
	}

	//prehash
	if mc.ID > 1 {
		var bk Block
		pfb, err := bk.GetId(mc.ID - 1)
		if err != nil {
			return &result, err
		}
		if pfb {
			header.PreHash = hex.EncodeToString(bk.Hash)
		}

	}
	if header.EcosystemId == 0 {
		header.EcosystemId = 1
	}

	bdi := BlockDetailedInfoHex{
		Header:        header,
		Hash:          hex.EncodeToString(mc.Hash),
		EcosystemID:   mc.EcosystemID,
		NodePosition:  mc.NodePosition,
		KeyID:         converter.AddressToString(mc.KeyID),
		Time:          mc.Time,
		Tx:            mc.Tx,
		RollbacksHash: hex.EncodeToString(mc.RollbacksHash),
		MerkleRoot:    hex.EncodeToString(blck.MerkleRoot),
		BinData:       hex.EncodeToString(blck.BinData),
		SysUpdate:     blck.SysUpdate,
		GenBlock:      blck.GenBlock,
		//StopCount:     blck.s,
		BlockSize:   ToCapacityString(int64(len(mc.Data))),
		TxTotalSize: ToCapacityString(transize),
	}

	if bdi.EcosystemID == 0 {
		bdi.EcosystemID = 1
	}
	return &bdi, nil
}

func GetBlocksTransactionListByBlockInfo(mc *Block) (*[]TxDetailedInfoResponse, error) {
	blck, err := block.UnmarshallBlock(bytes.NewBuffer(mc.Data))
	if err != nil {
		return nil, err
	}

	txDetailedInfoCollection := make([]TxDetailedInfoResponse, 0, len(blck.Transactions))
	for _, tx := range blck.Transactions {
		txDetailedInfo := TxDetailedInfoResponse{
			Hash: hex.EncodeToString(tx.Hash()),
		}
		if tx.IsSmartContract() {
			if tx.SmartContract().TxSmart.UTXO != nil {
				txDetailedInfo.ContractName = UtxoTx
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.UTXO)
				txDetailedInfo.Params = string(dataBytes)
			} else if tx.SmartContract().TxSmart.TransferSelf != nil {
				txDetailedInfo.ContractName = UtxoTransferSelf
				dataBytes, _ := json.Marshal(tx.SmartContract().TxSmart.TransferSelf)
				txDetailedInfo.Params = string(dataBytes)
			} else {
				txDetailedInfo.ContractName, _ = GetMineParam(tx.SmartContract().TxSmart.EcosystemID, tx.SmartContract().TxContract.Name, tx.SmartContract().TxData, tx.Hash())
			}
			txDetailedInfo.KeyID = converter.AddressToString(tx.KeyID())
			txDetailedInfo.Time = MsToSeconds(tx.Timestamp())
			txDetailedInfo.Type = int64(tx.Type())
			txDetailedInfo.Size = ToCapacityString(int64(len(tx.FullData)))
		}

		if txDetailedInfo.Time == 0 {
			txDetailedInfo.Time = mc.Time
		}
		if txDetailedInfo.KeyID == "" {
			txDetailedInfo.KeyID = converter.AddressToString(mc.KeyID)
		}

		if tx.IsSmartContract() {
			txDetailedInfo.Ecosystem = tx.SmartContract().TxSmart.EcosystemID
			if txDetailedInfo.Ecosystem == 0 {
				txDetailedInfo.Ecosystem = 1
			}
			txDetailedInfo.TokenSymbol, txDetailedInfo.Ecosystemname = Tokens.Get(txDetailedInfo.Ecosystem), EcoNames.Get(txDetailedInfo.Ecosystem)
		} else {
			if txDetailedInfo.Ecosystem == 0 {
				txDetailedInfo.Ecosystem = 1
			}
			if txDetailedInfo.Ecosystem == 1 {
				txDetailedInfo.TokenSymbol = SysTokenSymbol
				if txDetailedInfo.Ecosystemname == "" {
					txDetailedInfo.Ecosystemname = SysEcosystemName
				}
			}
		}

		txDetailedInfoCollection = append(txDetailedInfoCollection, txDetailedInfo)
	}

	return &txDetailedInfoCollection, nil
}
