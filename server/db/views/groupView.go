package views

import "gorm.io/gorm"

func ViewUserWalletsGroupByType(db *gorm.DB) error {
	if viewExist := db.Migrator().HasTable("view_user_wallets_group_by_type"); !viewExist {
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
	}

	return nil
}
