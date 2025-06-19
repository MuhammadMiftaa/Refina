package repository

import (
	"errors"

	"server/internal/dto"
	"server/internal/entity"

	"gorm.io/gorm"
)

type UsersRepository interface {
	GetAllUsers() ([]entity.Users, error)
	GetUserByID(id string) (entity.Users, error)
	GetUserByEmail(email string) (entity.Users, error)
	CreateUser(user entity.Users) (entity.Users, error)
	UpdateUser(user entity.Users) (entity.Users, error)
	DeleteUser(user entity.Users) (entity.Users, error)

	GetUserWallets(id string) ([]dto.ViewUserWallets, error)
	GetUserInvestments(id string) ([]dto.ViewUserInvestments, error)
	GetUserTransactions(id string) ([]dto.ViewUserTransactions, error)
}

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &usersRepository{db}
}

func (user_repo *usersRepository) GetAllUsers() ([]entity.Users, error) {
	var users []entity.Users
	err := user_repo.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (user_repo *usersRepository) GetUserByID(id string) (entity.Users, error) {
	var user entity.Users
	err := user_repo.db.First(&user, "id = ?", id).Error
	if err != nil {
		return entity.Users{}, errors.New("user not found")
	}

	return user, nil
}

func (user_repo *usersRepository) GetUserByEmail(email string) (entity.Users, error) {
	var user entity.Users
	err := user_repo.db.First(&user, "email = ?", email).Error
	if err != nil {
		return entity.Users{}, errors.New("user not found")
	}

	return user, nil
}

func (user_repo *usersRepository) CreateUser(user entity.Users) (entity.Users, error) {
	err := user_repo.db.Create(&user).Error
	if err != nil {
		return entity.Users{}, errors.New("failed to create user")
	}

	return user, nil
}

func (user_repo *usersRepository) UpdateUser(user entity.Users) (entity.Users, error) {
	err := user_repo.db.Save(&user).Error
	if err != nil {
		return entity.Users{}, errors.New("failed to update user")
	}

	return user, nil
}

func (user_repo *usersRepository) DeleteUser(user entity.Users) (entity.Users, error) {
	err := user_repo.db.Delete(&user).Error
	if err != nil {
		return entity.Users{}, errors.New("failed to delete user")
	}

	return user, nil
}

func (user_repo *usersRepository) GetUserWallets(id string) ([]dto.ViewUserWallets, error) {
	if viewExist := user_repo.db.Migrator().HasTable("view_user_wallets"); !viewExist {
		queryCreateUserWalletsView := `
		CREATE OR REPLACE VIEW view_user_wallets AS
		SELECT wallets.id, users.id AS user_id,
			wallets.number AS wallet_number, wallets.balance AS wallet_balance,
			wallets.name AS wallet_name, wallet_types.name AS wallet_type_name,
			wallet_types.type AS wallet_type
		FROM users
		LEFT JOIN wallets ON users.id = wallets.user_id AND wallets.deleted_at IS NULL
		LEFT JOIN wallet_types ON wallets.wallet_type_id = wallet_types.id AND wallet_types.deleted_at IS NULL
		WHERE users.deleted_at IS NULL;
	`

		if err := user_repo.db.Exec(queryCreateUserWalletsView).Error; err != nil {
			return nil, errors.New("failed to create user wallets view")
		}
	}

	var userWallets []dto.ViewUserWallets
	err := user_repo.db.Table("view_user_wallets").Where("user_id = ?", id).Find(&userWallets).Error
	if err != nil {
		return nil, errors.New("user wallets not found")
	}

	return userWallets, nil
}

func (user_repo *usersRepository) GetUserInvestments(id string) ([]dto.ViewUserInvestments, error) {
	if viewExist := user_repo.db.Migrator().HasTable("view_user_investments"); !viewExist {
		queryCreateUserInvestmentsView := `		
		CREATE OR REPLACE VIEW view_user_investments AS
		SELECT investments.id, users.id AS user_id,
			investment_types.name AS investment_type,
			investments.name AS investment_name,
			investments.amount AS investment_amount,
			investments.quantity AS investment_quantity,
			investment_types.unit AS investment_unit,
			investments.investment_date AS investment_date
		FROM users
		LEFT JOIN investments ON users.id = investments.user_id AND investments.deleted_at IS NULL
		LEFT JOIN investment_types ON investments.investment_type_id = investment_types.id AND investment_types.deleted_at IS NULL
		WHERE users.deleted_at IS NULL;
	`
		if err := user_repo.db.Exec(queryCreateUserInvestmentsView).Error; err != nil {
			return nil, errors.New("failed to create user investments view")
		}
	}

	var userInvestments []dto.ViewUserInvestments
	err := user_repo.db.Table("view_user_investments").Where("user_id = ?", id).Find(&userInvestments).Error
	if err != nil {
		return nil, errors.New("user investments not found")
	}

	return userInvestments, nil
}

func (user_repo *usersRepository) GetUserTransactions(id string) ([]dto.ViewUserTransactions, error) {
	if viewExist := user_repo.db.Migrator().HasTable("view_user_transactions"); !viewExist {
		queryCreateUserTransactionsView := `
		CREATE OR REPLACE VIEW view_user_transactions AS
		SELECT transactions.id AS id, users.id AS user_id,
			wallets.id AS wallet_id, wallets.number AS wallet_number, 
			wallet_types.name AS wallet_type, wallets.balance AS wallet_balance,
			categories.name AS category_name, categories.type AS category_type,
			transactions.amount, transactions.transaction_date, transactions.description,
			attachments.image
		FROM users
		LEFT JOIN wallets ON users.id = wallets.user_id AND wallets.deleted_at IS NULL
		LEFT JOIN wallet_types ON wallets.wallet_type_id = wallet_types.id AND wallet_types.deleted_at IS NULL
		LEFT JOIN transactions ON wallets.id = transactions.wallet_id AND transactions.deleted_at IS NULL
		LEFT JOIN categories ON transactions.category_id = categories.id AND categories.deleted_at IS NULL
		LEFT JOIN attachments ON transactions.id = attachments.transaction_id AND attachments.deleted_at IS NULL
		WHERE users.deleted_at IS NULL;
	`
		if err := user_repo.db.Exec(queryCreateUserTransactionsView).Error; err != nil {
			return nil, errors.New("failed to create user transactions view")
		}
	}

	var userTransactions []dto.ViewUserTransactions
	err := user_repo.db.Table("view_user_transactions").Where("user_id = ?", id).Find(&userTransactions).Error
	if err != nil {
		return nil, errors.New("user transactions not found")
	}

	return userTransactions, nil
}
