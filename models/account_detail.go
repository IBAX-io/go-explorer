package models

import (
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type AccountDetail struct {
	ID          int64           `gorm:"primary_key;not null"`
	Account     string          `gorm:"not null"`
	Ecosystem   int64           `gorm:"primary_key;not null"`
	JoinTime    int64           `gorm:"not null"`
	LogoHash    string          `gorm:"not null"` //ecosystem logo hash
	Amount      decimal.Decimal `gorm:"not null"`
	StakeAmount decimal.Decimal `gorm:"not null"` //Contains nft_miner_stake AND candidate_referendum AND candidate_substitute
	OutputValue decimal.Decimal `gorm:"not null"`
	TotalAmount decimal.Decimal `gorm:"not null"` //Contains Amount+OutputValue+StakeAmount
}

// TableName returns name of table
func (a *AccountDetail) TableName() string {
	return `account_detail`
}

func AccountDetailTableExist() bool {
	var p AccountDetail
	if !HasTableOrView(p.TableName()) {
		return false
	}
	return true
}

const CreateAccountDetailTable = `
-- ----------------------------
-- Table structure for account_detail
-- ----------------------------
DROP TABLE IF EXISTS "public"."account_detail";
CREATE TABLE "public"."account_detail" (
  "id" int8 NOT NULL DEFAULT '0'::bigint,
  "account" char(24) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::bpchar,
  "ecosystem" int8 NOT NULL DEFAULT '1'::bigint,
  "logo_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "join_time" int8 NOT NULL DEFAULT '0'::bigint,
  "amount" numeric(30,0) NOT NULL DEFAULT '0'::numeric,
  "stake_amount" numeric(30,0) NOT NULL DEFAULT '0'::numeric,
  "output_value" numeric(30,0) NOT NULL DEFAULT '0'::numeric,
  "total_amount" numeric(30,0) NOT NULL DEFAULT '0'::numeric
)
;

-- ----------------------------
-- Primary Key structure for table account_detail
-- ----------------------------
ALTER TABLE "public"."account_detail" ADD CONSTRAINT "account_detail_pkey" PRIMARY KEY ("id", "ecosystem");
`

const AccountDetailSQL = `
-- ----------------------------
-- Trigger Function for table account_detail; Update total_amount
-- ----------------------------
CREATE OR REPLACE FUNCTION update_account_detail_total()
RETURNS TRIGGER AS $$
DECLARE
	logo_id BIGINT;
BEGIN
		NEW.total_amount = NEW.stake_amount + NEW.amount + NEW.output_value;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ----------------------------
-- Trigger for table account_detail,before inserting or updating amount or stake_amount or output_value
-- ----------------------------
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_account_detail_total_trigger' AND tgrelid = 'account_detail'::regclass) THEN
		CREATE TRIGGER "update_account_detail_total_trigger" BEFORE INSERT OR UPDATE OF "amount", "stake_amount", "output_value" ON "public"."account_detail"
		FOR EACH ROW
		EXECUTE PROCEDURE "public"."update_account_detail_total"();
	END IF;
END$$;


-- ----------------------------
-- The function inserts account_detail id, account, amount fields
-- ----------------------------
CREATE OR REPLACE FUNCTION insert_delete_account_detail(ecosystem_id BIGINT, account_id BIGINT,account CHARACTER, amount NUMERIC, action_type INT,stake NUMERIC)
RETURNS VOID AS $$
BEGIN
		-- DELETE
		IF action_type = 1 THEN
			-- If Exist utxo table
			IF EXISTS (SELECT 1 FROM spent_info WHERE input_tx_hash is null AND output_key_id = account_id AND ecosystem = ecosystem_id) THEN
			-- UPDATE
				UPDATE account_detail SET amount = 0, account = '', stake_amount = 0 WHERE id = account_id AND ecosystem = ecosystem_id;
			ELSE
			-- DELETE
				DELETE FROM account_detail WHERE id = account_id AND ecosystem = ecosystem_id;
			END IF;
		ELSE
			-- INSERT
			INSERT INTO account_detail(id,account,ecosystem,amount,stake_amount) VALUES (account_id,account,ecosystem_id,amount,stake);
		END IF;
END;
$$ LANGUAGE plpgsql;


-- ----------------------------
-- Trigger Function; Insert OR Update account, amount, stake_amount fields to account_detail table
-- ----------------------------
CREATE OR REPLACE FUNCTION insert_delete_keys()
RETURNS TRIGGER AS $$
DECLARE
	logo_id BIGINT;
	action_type INT;
	stake NUMERIC;
BEGIN
		stake = 0;
		IF TG_OP = 'DELETE' THEN
			action_type = 1;
		ELSE
			action_type = 2;
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = '1_keys' AND column_name = 'lock') THEN
				stake:=to_number(coalesce(NULLIF(NEW.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999999999')+ 
				to_number(coalesce(NULLIF(NEW.lock->>'candidate_referendum',''),'0'),'999999999999999999999999999999') + 
				to_number(coalesce(NULLIF(NEW.lock->>'candidate_substitute',''),'0'),'999999999999999999999999999999');
			END IF;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM account_detail WHERE id = NEW.id AND ecosystem = NEW.ecosystem) THEN
			PERFORM insert_delete_account_detail(NEW.ecosystem,NEW.id, NEW.account, NEW.amount, action_type, stake);
		ELSE
			UPDATE account_detail SET amount = NEW.amount, account = NEW.account,stake_amount = stake WHERE id = NEW.id AND ecosystem = NEW.ecosystem;
		END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ----------------------------
-- Trigger Function; Insert OR Update account, amount, stake_amount fields to account_detail table
-- ----------------------------
CREATE OR REPLACE FUNCTION update_keys_amount()
RETURNS TRIGGER AS $$
DECLARE
	stake NUMERIC;
BEGIN
		stake = 0;
		IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = '1_keys' AND column_name = 'lock') THEN
			stake:=to_number(coalesce(NULLIF(NEW.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999999999')+ 
			to_number(coalesce(NULLIF(NEW.lock->>'candidate_referendum',''),'0'),'999999999999999999999999999999') + 
			to_number(coalesce(NULLIF(NEW.lock->>'candidate_substitute',''),'0'),'999999999999999999999999999999');
		END IF;
		
		IF NOT EXISTS (SELECT 1 FROM account_detail WHERE id = NEW.id AND ecosystem = NEW.ecosystem) THEN
			PERFORM insert_delete_account_detail(NEW.ecosystem,NEW.id, NEW.account, NEW.amount, 2,stake);
		ELSE
			UPDATE account_detail SET amount = NEW.amount, account = NEW.account,stake_amount = stake WHERE id = NEW.id AND ecosystem = NEW.ecosystem;
		END IF;
		RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ----------------------------
-- Trigger for table "1_keys", BEFORE INSERT OR DELETE; Insert to account_detail table
-- ----------------------------
DO $$
DECLARE
  	row record;
	create_key_id BIGINT;
	block_time BIGINT;
	join_time BIGINT;
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'insert_delete_keys_trigger' AND tgrelid = '1_keys'::regclass) THEN
		SELECT output_key_id INTO create_key_id FROM spent_info WHERE block_id = 1 AND ecosystem = 1 AND type = 1 LIMIT 1;
		SELECT time INTO block_time FROM block_chain WHERE id = 1;
		
		FOR row IN
				SELECT id,account,ecosystem,amount,pub FROM "1_keys"
		LOOP
			join_time = 0; -- default
			IF (row.id = -110277540701013350 OR row.id = create_key_id) AND length(row.pub) = 64 AND row.ecosystem = 1 THEN
				join_time = block_time;
			END IF;
			
			IF NOT EXISTS (SELECT 1 FROM account_detail WHERE id = row.id AND ecosystem = row.ecosystem) THEN
				INSERT INTO account_detail(id,account,amount,ecosystem,join_time)	VALUES (row.id, row.account, row.amount, row.ecosystem, join_time);
			ELSE
				UPDATE account_detail SET account = row.account, amount = row.amount WHERE id = row.id AND ecosystem = row.ecosystem;
			END IF;
		END LOOP;
	
		CREATE TRIGGER "insert_delete_keys_trigger" BEFORE INSERT OR DELETE ON "public"."1_keys"
		FOR EACH ROW
		EXECUTE PROCEDURE "public"."insert_delete_keys"();
	END IF;
END$$;

-- ----------------------------
-- Trigger for table "1_keys"; Invoke the trigger after updating the keys table amount field
-- ----------------------------
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_keys_amount_trigger' AND tgrelid = '1_keys'::regclass) THEN
		CREATE TRIGGER "update_keys_amount_trigger" BEFORE UPDATE OF "amount" ON "public"."1_keys"
		FOR EACH ROW
		EXECUTE PROCEDURE "public"."update_keys_amount"();
	END IF;
END$$;


-- ----------------------------
-- The function inserts or update account_detail output_value fields
-- ----------------------------
CREATE OR REPLACE FUNCTION update_account_detail_output_value(ecosystem_id BIGINT, account_id BIGINT, output_amount NUMERIC, action_type INT)
RETURNS VOID AS $$
BEGIN
		IF action_type = 1 THEN
			-- DELETE
			UPDATE account_detail SET output_value = output_value + output_amount WHERE id = account_id AND ecosystem = ecosystem_id;
		ELSE
			-- INSERT
			INSERT INTO account_detail(id,ecosystem,output_value) VALUES (account_id, ecosystem_id, output_amount);
		END IF;
END;
$$ LANGUAGE plpgsql;


-- ----------------------------
-- Trigger Function; inserts or update account_detail output_value fields
-- ----------------------------
CREATE OR REPLACE FUNCTION update_spent_info_output_value()
RETURNS TRIGGER AS $$
DECLARE
	action_type INT;
	output_amount NUMERIC;
BEGIN
		IF TG_OP = 'DELETE' THEN
			action_type = 1;
			IF OLD.input_tx_hash IS NULL THEN
				output_amount = OLD.output_value * -1;
			ELSE
				output_amount = OLD.output_value;
			END IF;	
			PERFORM update_account_detail_output_value(OLD.ecosystem, OLD.output_key_id, output_amount, action_type);
		ELSIF TG_OP = 'INSERT' THEN
			action_type = 2;
			if NEW.input_tx_hash IS NULL THEN  -- INSERT
				IF NOT EXISTS (SELECT 1 FROM account_detail WHERE id = NEW.output_key_id AND ecosystem = NEW.ecosystem) THEN
					PERFORM update_account_detail_output_value(NEW.ecosystem, NEW.output_key_id, NEW.output_value, action_type);
				ELSE
					UPDATE account_detail SET output_value = output_value + NEW.output_value WHERE id = NEW.output_key_id AND ecosystem = NEW.ecosystem;
				END IF;
			END IF;
		ELSEIF TG_OP = 'UPDATE' THEN
			if NEW.input_tx_hash IS NOT NULL THEN
				UPDATE account_detail SET output_value = output_value - OLD.output_value WHERE id = OLD.output_key_id AND ecosystem = OLD.ecosystem;
			END IF;
		END IF;
		RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ----------------------------
-- Trigger for table "spent_info"; Invoke the trigger after updating or insert the spent_info table input_tx_hash field; Synchronize all UTXO account balances to the accountdetail table
-- ----------------------------
DO $$
DECLARE
  row record;
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_account_detail_output_value_trigger' AND tgrelid = 'spent_info'::regclass) THEN
		FOR row IN
				SELECT output_key_id,ecosystem,sum(output_value) AS output_value FROM spent_info WHERE input_tx_hash IS NULL GROUP BY output_key_id,ecosystem
		LOOP
			IF NOT EXISTS (SELECT 1 FROM account_detail WHERE id = row.output_key_id AND ecosystem = row.ecosystem) THEN
				INSERT INTO account_detail(id, output_value,ecosystem)	VALUES (row.output_key_id, row.output_value, row.ecosystem);
			ELSE
				UPDATE account_detail SET output_value = row.output_value WHERE id = row.output_key_id AND ecosystem = row.ecosystem;
			END IF;
		END LOOP;
		
		CREATE TRIGGER "update_account_detail_output_value_trigger" AFTER INSERT OR UPDATE OF "input_tx_hash" OR DELETE ON "public"."spent_info"
		FOR EACH ROW
		EXECUTE PROCEDURE "public"."update_spent_info_output_value"();
	END IF;
END$$;


-- ----------------------------
-- The function update account_detail logo_hash fields
-- ----------------------------
CREATE OR REPLACE FUNCTION update_account_detail_logo_hash(ecosystem_id BIGINT, logo_id BIGINT, action_type INT)
RETURNS VOID AS $$
DECLARE
	logo VARCHAR;
BEGIN
		IF action_type = 2 THEN
			SELECT COALESCE((SELECT hash FROM "1_binaries" WHERE id = logo_id),'') INTO logo;
			IF logo <> '' THEN
				UPDATE account_detail SET logo_hash = logo WHERE ecosystem = ecosystem_id;
			END IF;
		ELSE
			UPDATE account_detail SET logo_hash = '' WHERE ecosystem = ecosystem_id;
		END IF;
END;
$$ LANGUAGE plpgsql;


-- -----------------------------
-- Trigger Function; inserts or update account_detail logo_hash fields
-- ----------------------------
CREATE OR REPLACE FUNCTION update_ecosystem_info()
RETURNS TRIGGER AS $$
DECLARE
	logo_id BIGINT;
	action_type INT;
BEGIN
		IF TG_OP = 'DELETE' THEN
			action_type = 1;
		ELSE
			action_type = 2;
			-- get logo id
			logo_id := COALESCE(cast(NEW.info->>'logo' as numeric),0);
		END IF;
		PERFORM update_account_detail_logo_hash(NEW.id, logo_id, action_type);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ----------------------------
-- Trigger for table "1_ecosystems"; Invoke the trigger BEFORE updating or delete the "1_ecosystems" table info field
-- ----------------------------
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_ecosystems_info_trigger' AND tgrelid = '1_ecosystems'::regclass) THEN
		CREATE TRIGGER "update_ecosystems_info_trigger" BEFORE UPDATE OF "info" OR DELETE ON "public"."1_ecosystems"
		FOR EACH ROW
		EXECUTE PROCEDURE "public"."update_ecosystem_info"();
	END IF;
END$$;


-- -----------------------------
-- Trigger Function; update account_detail join_time fields
-- ----------------------------
CREATE OR REPLACE FUNCTION update_account_detail_join_time()
RETURNS TRIGGER AS $$
DECLARE
  key_id BIGINT;
	ecosystem_id BIGINT;
	join_t BIGINT;
BEGIN
		IF TG_OP = 'DELETE' THEN
			IF OLD.table_name = '1_keys' AND OLD.data = '' AND OLD.table_id SIMILAR TO '-?[0-9]+,[0-9]+' THEN
				key_id := CAST(SPLIT_PART(OLD.table_id, ',', 1)AS BIGINT);
				ecosystem_id := CAST(SPLIT_PART(OLD.table_id, ',', 2)AS BIGINT);
				UPDATE account_detail SET join_time = 0 WHERE id = key_id AND ecosystem = ecosystem_id;
			END IF;
		ELSIF TG_OP = 'INSERT' THEN
			IF NEW.table_name = '1_keys' AND NEW.data = '' AND NEW.table_id SIMILAR TO '-?[0-9]+,[0-9]+' THEN
				key_id := CAST(SPLIT_PART(NEW.table_id, ',', 1)AS BIGINT);
				ecosystem_id := CAST(SPLIT_PART(NEW.table_id, ',', 2)AS BIGINT);
				SELECT COALESCE((SELECT timestamp/1000 FROM log_transactions WHERE hash = NEW.tx_hash),0) INTO join_t;
				if NOT EXISTS (SELECT 1 FROM "1_keys" WHERE id = key_id AND ecosystem = ecosystem_id AND LENGTH(pub) = 64) THEN
					join_t = 0;
				END IF;
				UPDATE account_detail SET join_time = join_t WHERE id = key_id AND ecosystem = ecosystem_id;
			END IF;
		END IF;
		RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ----------------------------
-- Trigger for table "rollback_tx"; Invoke the trigger BEFORE insert the "rollback_tx"
-- ----------------------------
DO $$
DECLARE
  row record;
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_account_detail_join_time_trigger' AND tgrelid = 'rollback_tx'::regclass) THEN
		FOR row IN
			SELECT CAST(SPLIT_PART(table_id, ',', 1)AS BIGINT) AS key_id, CAST(SPLIT_PART(table_id, ',', 2)AS BIGINT) AS ecosystem,bk.time FROM rollback_tx rt LEFT JOIN block_chain AS bk ON(bk.id = rt.block_id) 
			WHERE rt.table_name = '1_keys' AND rt.data = '' AND rt.table_id SIMILAR TO '-?[0-9]+,[0-9]+'
		LOOP
			if NOT EXISTS (SELECT 1 FROM "1_keys" WHERE id = row.key_id AND ecosystem = row.ecosystem AND LENGTH(pub) = 64) THEN
				row.time = 0;
			END IF;
			
			IF NOT EXISTS (SELECT 1 FROM account_detail WHERE id = row.key_id AND ecosystem = row.ecosystem) THEN
				INSERT INTO account_detail(id, join_time,ecosystem)	VALUES (row.key_id, row.time, row.ecosystem);
			ELSE
				UPDATE account_detail SET join_time = row.time WHERE id = row.key_id AND ecosystem = row.ecosystem;
			END IF;
		END LOOP;
		
		CREATE TRIGGER "update_account_detail_join_time_trigger" BEFORE INSERT ON "public"."rollback_tx"
		FOR EACH ROW
		EXECUTE PROCEDURE "public"."update_account_detail_join_time"();
	END IF;
END$$;

-- ----------------------------
-- Trigger for table "rollback_tx"; Invoke the trigger AFTER delete the "rollback_tx"
-- ----------------------------
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'delete_account_detail_join_time_trigger' AND tgrelid = 'rollback_tx'::regclass) THEN
		CREATE TRIGGER "delete_account_detail_join_time_trigger" AFTER DELETE ON "public"."rollback_tx"
		FOR EACH ROW
		EXECUTE PROCEDURE "public"."update_account_detail_join_time"();
	END IF;
END$$;
`

func UpdateAccountDetail() {
	RealtimeWG.Add(1)
	defer func() {
		RealtimeWG.Done()
	}()
	a := &AccountDetail{}
	err := a.updateAccount()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("update account by account detail")
		return
	}
}

func (a *AccountDetail) updateAccount() (err error) {
	var list []AccountDetail
	err = GetDB(nil).Where("id <> 0 AND account = ''").Limit(1000).Find(&list).Error
	if err != nil {
		return
	}
	for _, v := range list {
		err = GetDB(nil).Model(AccountDetail{}).Where("id = ? AND ecosystem = ?", v.ID, v.Ecosystem).Update("account", converter.AddressToString(v.ID)).Error
		if err != nil {
			return
		}
	}
	return
}

func GetEcosystemDetailMemberList(page, limit int, order string, ecosystem int64) (*GeneralResponse, error) {
	var (
		rets GeneralResponse
	)
	if order == "" {
		order = "join_time desc"
	}
	rets.Limit = limit
	rets.Page = page

	err := GetDB(nil).Model(AccountDetail{}).
		Where("ecosystem = ?", ecosystem).
		Count(&rets.Total).Error
	if err != nil {
		return nil, err
	}

	var ret []EcosystemMemberList
	err = GetDB(nil).Table("account_detail AS ad").
		Select(`id,account,(SELECT array_to_string(array(SELECT rs.role->>'name' FROM "1_roles_participants" as rs 
	WHERE rs.ecosystem=ad.ecosystem and rs.member->>'account' = ad.account AND rs.deleted = 0),' / '))AS roles_name,join_time`).
		Where("ecosystem = ?", ecosystem).
		Order(order).Offset((page - 1) * limit).Limit(limit).
		Find(&ret).Error
	if INameReady {
		for k, v := range ret {
			ie := &IName{}
			f, err := ie.Get(v.Account)
			if err == nil && f {
				ret[k].AccountName = ie.Name
			}
		}
	}
	rets.List = ret

	return &rets, err
}
