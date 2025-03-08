package repository

import (
	"errors"

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

	GetUserWallets(id string) ([]entity.UserWallet, error)
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

func (user_repo *usersRepository) GetUserWallets(id string) ([]entity.UserWallet, error) {
	var userWallet []entity.UserWallet
	err := user_repo.db.Table("users").Select("wallets.id, users.id AS user_id, users.name, users.email, wallets.number AS wallet_number, wallets.balance AS wallet_balance, wallet_types.name AS wallet_name, wallet_types.type AS wallet_type").
		Joins("LEFT JOIN wallets ON users.id = wallets.user_id AND wallets.deleted_at IS NULL").
		Joins("LEFT JOIN wallet_types ON wallets.wallet_type_id = wallet_types.id AND wallet_types.deleted_at IS NULL").
		Where("users.id = ?", id).
		Where("users.deleted_at IS NULL").
		Find(&userWallet).Error
	if err != nil {
		return nil, errors.New("user wallet not found")
	}

	return userWallet, nil
}
