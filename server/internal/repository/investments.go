package repository

import (
	"context"
	"errors"

	"server/internal/types/entity"
	"server/internal/types/view"

	"gorm.io/gorm"
)

type InvestmentsRepository interface {
	GetAllInvestments(ctx context.Context, tx Transaction) ([]entity.Investments, error)
	GetInvestmentByID(ctx context.Context, tx Transaction, id string) (entity.Investments, error)
	GetInvestmentsByUserID(ctx context.Context, tx Transaction, id string) ([]view.ViewUserInvestments, error)
	CreateInvestment(ctx context.Context, tx Transaction, investment entity.Investments) (entity.Investments, error)
	UpdateInvestment(ctx context.Context, tx Transaction, investment entity.Investments) (entity.Investments, error)
	DeleteInvestment(ctx context.Context, tx Transaction, investment entity.Investments) (entity.Investments, error)
}

type investmentsRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) InvestmentsRepository {
	return &investmentsRepository{db}
}

// Helper untuk mendapatkan DB instance (transaksi atau biasa)
func (investment_repo *investmentsRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx) // Type assertion ke GORM transaction
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return investment_repo.db.WithContext(ctx), nil
}

func (investment_repo *investmentsRepository) GetAllInvestments(ctx context.Context, tx Transaction) ([]entity.Investments, error) {
	db, err := investment_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var investments []entity.Investments
	if err := db.Find(&investments).Error; err != nil {
		return nil, err
	}

	return investments, nil
}

func (investment_repo *investmentsRepository) GetInvestmentByID(ctx context.Context, tx Transaction, id string) (entity.Investments, error) {
	db, err := investment_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Investments{}, err
	}

	var investment entity.Investments
	if err := db.First(&investment, "id = ?", id).Error; err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}

func (investment_repo *investmentsRepository) GetInvestmentsByUserID(ctx context.Context, tx Transaction, id string) ([]view.ViewUserInvestments, error) {
	db, err := investment_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var userInvestments []view.ViewUserInvestments
	err = db.Table("view_user_investments").Where("user_id = ?", id).Find(&userInvestments).Error
	if err != nil {
		return nil, errors.New("user investments not found")
	}

	return userInvestments, nil
}

func (investment_repo *investmentsRepository) CreateInvestment(ctx context.Context, tx Transaction, investment entity.Investments) (entity.Investments, error) {
	db, err := investment_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Investments{}, err
	}

	if err := db.Create(&investment).Error; err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}

func (investment_repo *investmentsRepository) UpdateInvestment(ctx context.Context, tx Transaction, investment entity.Investments) (entity.Investments, error) {
	db, err := investment_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Investments{}, err
	}

	if err := db.Omit("InvestmentTypes", "User").Save(&investment).Error; err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}

func (investment_repo *investmentsRepository) DeleteInvestment(ctx context.Context, tx Transaction, investment entity.Investments) (entity.Investments, error) {
	db, err := investment_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Investments{}, err
	}

	if err := db.Delete(&investment).Error; err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}
