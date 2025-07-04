package views

import (
	"fmt"

	"gorm.io/gorm"
)

func ViewUserWalletsGroupByType(db *gorm.DB) error {
	queryDropView := `DROP VIEW IF EXISTS view_user_wallets_group_by_type;`
	if err := db.Exec(queryDropView).Error; err != nil {
		return fmt.Errorf("failed to drop existing view: %w", err)
	}

	queryCreateUserWalletsGroupByTypeView := `
		CREATE OR REPLACE VIEW view_user_wallets_group_by_type AS
		SELECT 
			users.id AS user_id,
			wallet_types.type AS type,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', wallets.id,
					'name', wallets.name,
					'number', wallets.number,
					'balance', wallets.balance
				)
			) AS wallets 
		FROM wallets
		JOIN users ON users.id = wallets.user_id AND users.deleted_at IS NULL
		JOIN wallet_types ON wallet_types.id = wallets.wallet_type_id AND wallet_types.deleted_at IS NULL
		WHERE wallets.deleted_at IS NULL
		GROUP BY users.id, wallet_types.type;
	`

	if err := db.Exec(queryCreateUserWalletsGroupByTypeView).Error; err != nil {
		return err
	}

	return nil
}

func ViewCategoryGroupByType(db *gorm.DB) error {
	queryDropView := `DROP VIEW IF EXISTS view_category_group_by_type;`
	if err := db.Exec(queryDropView).Error; err != nil {
		return fmt.Errorf("failed to drop existing view: %w", err)
	}

	queryCreateCategoryGroupByTypeView := `
		CREATE OR REPLACE VIEW view_category_group_by_type AS
		SELECT 
			parent.name AS group_name,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', child.id,
					'name', child.name
				)
				ORDER BY child.name
			) AS category,
			parent.type AS type
		FROM categories parent
		LEFT JOIN categories child ON child.parent_id = parent.id
		WHERE parent.parent_id IS NULL AND parent.deleted_at IS NULL
		GROUP BY parent.name, parent.type;

	`

	if err := db.Exec(queryCreateCategoryGroupByTypeView).Error; err != nil {
		return err
	}

	return nil
}
