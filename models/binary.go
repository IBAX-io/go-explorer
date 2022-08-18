/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"path"
	"strings"

	"github.com/IBAX-io/go-ibax/packages/converter"
)

const BinaryTableSuffix = "_binaries"

// Binary represents record of {prefix}_binaries table
type Binary struct {
	Ecosystem int64
	ID        int64
	AppId     int64
	Name      string
	Data      []byte
	Hash      string
	MimeType  string
	Account   string
}

// SetTablePrefix is setting table prefix
func (b *Binary) SetTablePrefix(prefix string) {
	b.Ecosystem = converter.StrToInt64(prefix)
}

// SetTableName sets name of table
func (b *Binary) SetTableName(tableName string) {
	ecosystem, _ := converter.ParseName(tableName)
	b.Ecosystem = ecosystem
}

// TableName returns name of table
func (b *Binary) TableName() string {
	if b.Ecosystem == 0 {
		b.Ecosystem = 1
	}
	return `1_binaries`
}

// GetByID is retrieving model from db by id
func (b *Binary) GetByID(id int64) (bool, error) {
	return isFound(GetDB(nil).Where("id=?", id).First(b))
}

func (b *Binary) GetByHash(hash string) (bool, error) {
	return isFound(GetDB(nil).Select("id,hash,name,mime_type,ecosystem,app_id").Where("hash=?", hash).First(b))
}

func (b *Binary) GetByIdHash(id int64) (bool, error) {
	return isFound(GetDB(nil).Select("hash").Where("id=?", id).First(b))
}

func (b *Binary) GetByPng(d *Binary) (bool, error) {
	return isFound(GetDB(nil).Where("ecosystem = ? and app_id=? and hash = ? and mime_type !=?", d.Ecosystem, d.AppId, d.Hash, d.MimeType).First(b))
}

func (b *Binary) GetByJpeg() string {
	file := ""
	fileSuffix := path.Ext(b.Name) //
	if fileSuffix == "" {
		if b.MimeType == "image/jpeg" || b.MimeType == "image/jpg" || b.MimeType == "image/png" {
			ns := strings.Split(b.MimeType, "/")
			if len(ns) == 2 {
				//file = strconv.FormatInt(b.ID, 10) + "-" + b.Name + "." + ns[1]
				file = b.Hash + "." + ns[1]
				return file
			} else {
				return file
			}
		} else {
			file = b.Hash
			return file
		}
	} else {
		file = b.Hash + fileSuffix
		return file
	}

	//file = strconv.FormatInt(b.ID, 10) + "-" + b.Name
	//return file
}

func (b *Binary) FindAppNameByEcosystem(page, limit int, order string, ecosystem int64, where map[string]any) (GeneralResponse, error) {
	var by []Binary
	var rets GeneralResponse
	var total int64
	if order == "" {
		order = "id desc"
	}
	if len(where) == 0 {
		if err := GetDB(nil).Table(b.TableName()).Where("ecosystem = ?", ecosystem).Count(&total).Error; err != nil {
			return rets, err
		}

		if err := GetDB(nil).Select("id,name,hash,mime_type").Where("ecosystem = ?", ecosystem).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&by).Error; err != nil {
			return rets, err
		}
	} else {
		where["ecosystem ="] = ecosystem
		cond, vals, err := WhereBuild(where)
		if err != nil {
			return rets, err
		}
		if err := GetDB(nil).Table(b.TableName()).Where(cond, vals...).Count(&total).Error; err != nil {
			return rets, err
		}

		if err := GetDB(nil).Select("id,name,hash,mime_type").Where(cond, vals...).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&by).Error; err != nil {
			return rets, err
		}

	}
	rets.Page = page
	rets.Limit = limit
	rets.Total = total
	ret := make([]EcosystemAttachmentResponse, len(by))
	for i, value := range by {
		ret[i].Id = value.ID
		ret[i].Name = value.Name
		ret[i].Hash = value.Hash
		ret[i].MimeType = value.MimeType
	}
	rets.List = ret
	return rets, nil
}

func CompareHash(data []byte, urlHash string) bool {
	urlHash = strings.ToLower(urlHash)

	var hash []byte
	switch len(urlHash) {
	case 32:
		h := md5.Sum(data)
		hash = h[:]
	case 64:
		hash = crypto.Hash(data)
	}

	return hex.EncodeToString(hash) == urlHash
}
