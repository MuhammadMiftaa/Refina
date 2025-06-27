package views

import (
	"errors"

	"gorm.io/gorm"
)

func ViewUserWallets(db *gorm.DB) error {
	if viewExist := db.Migrator().HasTable("view_user_wallets"); !viewExist {
		queryCreateUserWalletsView := `
		CREATE OR REPLACE VIEW view_user_wallets AS
		SELECT wallets.id, users.id AS user_id,
			wallets.number AS wallet_number, wallets.balance AS wallet_balance,
			wallets.name AS wallet_name, wallet_types.name AS wallet_type_name,
			wallet_types.type AS wallet_type
		FROM wallets
		LEFT JOIN users ON users.id = wallets.user_id AND users.deleted_at IS NULL
		LEFT JOIN wallet_types ON wallet_types.id = wallets.wallet_type_id AND wallet_types.deleted_at IS NULL
		WHERE wallets.deleted_at IS NULL;
	`

		if err := db.Exec(queryCreateUserWalletsView).Error; err != nil {
			return errors.New("failed to create user wallets view")
		}
	}

	return nil
}

func ViewUserInvestments(db *gorm.DB) error {
	if viewExist := db.Migrator().HasTable("view_user_investments"); !viewExist {
		queryCreateUserInvestmentsView := `		
		CREATE OR REPLACE VIEW view_user_investments AS
		SELECT investments.id, users.id AS user_id,
			investment_types.name AS investment_type,
			investments.name AS investment_name,
			investments.amount AS investment_amount,
			investments.quantity AS investment_quantity,
			investment_types.unit AS investment_unit,
			investments.investment_date AS investment_date
		FROM investments
		LEFT JOIN users ON users.id = investments.user_id AND users.deleted_at IS NULL
		LEFT JOIN investment_types ON investment_types.id = investments.investment_type_id AND investment_types.deleted_at IS NULL
		WHERE investments.deleted_at IS NULL;
	`
		if err := db.Exec(queryCreateUserInvestmentsView).Error; err != nil {
			return errors.New("failed to create user investments view")
		}
	}

	return nil
}

func ViewUserTransactions(db *gorm.DB) error {
	if viewExist := db.Migrator().HasTable("view_user_transactions"); !viewExist {
		queryCreateUserTransactionsView := `
		CREATE OR REPLACE VIEW view_user_transactions AS
		SELECT transactions.id AS id, users.id AS user_id,
			wallets.id AS wallet_id, wallets.number AS wallet_number, 
			wallet_types.name AS wallet_type, wallets.balance AS wallet_balance,
			categories.name AS category_name, categories.type AS category_type,
			transactions.amount, transactions.transaction_date, transactions.description,
			wallet_types.type AS wallet_type_name
		FROM transactions
		LEFT JOIN wallets ON wallets.id = transactions.wallet_id AND transactions.deleted_at IS NULL
		LEFT JOIN users ON users.id = wallets.user_id AND users.deleted_at IS NULL
		LEFT JOIN wallet_types ON wallet_types.id = wallets.wallet_type_id AND wallet_types.deleted_at IS NULL
		LEFT JOIN categories ON categories.id = transactions.category_id AND categories.deleted_at IS NULL
		WHERE transactions.deleted_at IS NULL;
	`
		if err := db.Exec(queryCreateUserTransactionsView).Error; err != nil {
			return errors.New("failed to create user transactions view")
		}
	}

	return nil
}
