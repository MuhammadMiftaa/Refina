package repository

import (
	"server/internal/entity"

	"gorm.io/gorm"
)

type InvestmentsRepository interface {
	GetAllInvestments() ([]entity.Investments, error)
	GetInvestmentByID(id string) (entity.Investments, error)
	GetInvestmentsByUserID(id string) ([]entity.Investments, error)
	CreateInvestment(investment entity.Investments) (entity.Investments, error)
	UpdateInvestment(investment entity.Investments) (entity.Investments, error)
	DeleteInvestment(investment entity.Investments) (entity.Investments, error)
}

type investmentsRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) InvestmentsRepository {
	return &investmentsRepository{db}
}

func (investment_repo *investmentsRepository) GetAllInvestments() ([]entity.Investments, error) {
	var investments []entity.Investments
	err := investment_repo.db.Find(&investments).Error
	if err != nil {
		return nil, err
	}

	return investments, nil
}

func (investment_repo *investmentsRepository) GetInvestmentByID(id string) (entity.Investments, error) {
	var investment entity.Investments
	err := investment_repo.db.First(&investment, "id = ?", id).Error
	if err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}

func (investment_repo *investmentsRepository) GetInvestmentsByUserID(id string) ([]entity.Investments, error) {
	var investments []entity.Investments
	err := investment_repo.db.Find(&investments, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return investments, nil
}

func (investment_repo *investmentsRepository) CreateInvestment(investment entity.Investments) (entity.Investments, error) {
	err := investment_repo.db.Create(&investment).Error
	if err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}

func (investment_repo *investmentsRepository) UpdateInvestment(investment entity.Investments) (entity.Investments, error) {
	err := investment_repo.db.Save(&investment).Error
	if err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}

func (investment_repo *investmentsRepository) DeleteInvestment(investment entity.Investments) (entity.Investments, error) {
	err := investment_repo.db.Delete(&investment).Error
	if err != nil {
		return entity.Investments{}, err
	}

	return investment, nil
}
