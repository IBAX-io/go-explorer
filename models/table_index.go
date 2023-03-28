package models

import "fmt"

const (
	IndexBTree  = "btree"
	IndexHash   = "Hash"
	IndexGiST   = "GiST"
	IndexGIN    = "gin"
	IndexSPGiST = "SP-GiST"
	IndexBRIN   = "BRIN"
)

func ParseIndexMethod(method string) int {
	switch method {
	case IndexBTree:
		return 0
	case IndexHash:
		return 1
	case IndexGiST:
		return 2
	case IndexGIN:
		return 3
	case IndexSPGiST:
		return 4
	case IndexBRIN:
		return 5
	}
	return 0
}

func FormatIndexMethod(method int) string {
	switch method {
	case 0:
		return IndexBTree
	case 1:
		return IndexHash
	case 2:
		return IndexGiST
	case 3:
		return IndexGIN
	case 4:
		return IndexSPGiST
	case 5:
		return IndexBRIN
	}
	return IndexBTree
}

func extensionExist(extensionName string) (bool, error, string) {
	var extname string
	f, err := isFound(GetDB(nil).Table("pg_extension").Select("extname").Where("extname = ?", extensionName).Take(&extname))
	return f, err, extname
}

func deleteExtension(extensionName string) error {
	f, err, _ := extensionExist(extensionName)
	if err != nil {
		return err
	}
	if f {
		err = GetDB(nil).Exec(fmt.Sprintf("DROP EXTENSION %s", extensionName)).Error
		if err != nil {
			return fmt.Errorf("delete extension failed:%s", err.Error())
		}
	}
	return nil
}

func createExtension(extensionName string) error {
	f, err, _ := extensionExist(extensionName)
	if err != nil {
		return err
	}
	if !f {
		err = GetDB(nil).Exec(fmt.Sprintf("CREATE EXTENSION %s", extensionName)).Error
		if err != nil {
			return fmt.Errorf("create extension failed:%s", err.Error())
		}
	}
	return nil
}

func tableIndexExist(tableName, indexName string) (bool, error, string) {
	var name string
	f, err := isFound(GetDB(nil).Table("pg_stat_user_indexes").Select("indexrelname").
		Where("relname = ? AND indexrelname = ?", tableName, indexName).Take(&indexName))
	return f, err, name
}

func createTableIndex(tableName, indexName, rows string, method int) error {
	f, err, _ := tableIndexExist(tableName, indexName)
	if err != nil {
		return err
	}
	if f {
		return nil
	}
	switch method {
	case 0, 1, 3:
		err = GetDB(nil).Exec(fmt.Sprintf(`
CREATE INDEX "%s"
             ON "%s" using %s (%s)
`, indexName, tableName, FormatIndexMethod(method), rows)).Error
		if err != nil {
			return fmt.Errorf("create table index failed:%s", err.Error())
		}
	default:
		return fmt.Errorf("not support index method:%s", FormatIndexMethod(method))
	}
	return nil
}

func deleteTableIndex(tableName, indexName string) error {
	f, err, _ := tableIndexExist(tableName, indexName)
	if err != nil {
		return err
	}
	if !f {
		return nil
	}

	err = GetDB(nil).Exec(fmt.Sprintf(`
DROP INDEX %s
`, indexName)).Error
	if err != nil {
		return fmt.Errorf("create table index failed:%s", err.Error())
	}
	return nil
}

// CreateIndexMain Indexes cannot be added to the ecological table, which will cause synchronization failure
func CreateIndexMain() error {
	var sp SpentInfo
	err := createTableIndex(sp.TableName(), sp.TableName()+"_input_tx_hash_output_key_id_ecosystem_idx",
		`"input_tx_hash","output_key_id","ecosystem"`, ParseIndexMethod(IndexBTree))
	if err != nil {
		return err
	}

	var sph SpentInfoHistory
	err = createTableIndex(sph.TableName(), sph.TableName()+"_sender_id_recipient_id_ecosystem_idx",
		`"sender_id","recipient_id","ecosystem"`, ParseIndexMethod(IndexBTree))
	if err != nil {
		return err
	}

	var lg LogTransaction
	err = createTableIndex(lg.TableName(), lg.TableName()+"_ecosystem_id_block_timestamp_idx",
		`"ecosystem_id","block","timestamp"`, ParseIndexMethod(IndexBTree))
	if err != nil {
		return err
	}

	//must be set Database UTF-8 Format
	err = createExtension("btree_gin")
	if err != nil {
		return err
	}

	err = createExtension("pg_trgm")
	if err != nil {
		return err
	}

	var txr TransactionRelation
	err = createTableIndex(txr.TableName(), txr.TableName()+"_sender_ids_recipient_ids_idx",
		`"sender_ids" GIN_TRGM_OPS,"recipient_ids" GIN_TRGM_OPS`, ParseIndexMethod(IndexGIN))
	if err != nil {
		return err
	}

	return nil
}
