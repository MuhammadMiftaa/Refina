package service

import (
	"errors"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/repository"
)

type InvestmentsService interface {
	GetAllInvestments() ([]entity.Investments, error)
	GetInvestmentByID(id string) (entity.Investments, error)
	GetInvestmentsByUserID(id string) ([]entity.Investments, error)
	CreateInvestment(investment entity.InvestmentsRequest) (entity.Investments, error)
	UpdateInvestment(id string, investment entity.InvestmentsRequest) (entity.Investments, error)
	DeleteInvestment(id string) (entity.Investments, error)
}

type investmentsService struct {
	investmentsRepository repository.InvestmentsRepository
}

func NewInvestmentService(investmentsRepository repository.InvestmentsRepository) InvestmentsService {
	return &investmentsService{investmentsRepository}
}

func (investment_serv *investmentsService) GetAllInvestments() ([]entity.Investments, error) {
	investments, err := investment_serv.investmentsRepository.GetAllInvestments()
	if err != nil {
		return nil, errors.New("failed to get investments")
	}

	return investments, nil
}

func (investment_serv *investmentsService) GetInvestmentByID(id string) (entity.Investments, error) {
	investment, err := investment_serv.investmentsRepository.GetInvestmentByID(id)
	if err != nil {
		return entity.Investments{}, errors.New("investment not found")
	}

	return investment, nil
}

func (investment_serv *investmentsService) GetInvestmentsByUserID(id string) ([]entity.Investments, error) {
	investments, err := investment_serv.investmentsRepository.GetInvestmentsByUserID(id)
	if err != nil {
		return nil, errors.New("failed to get investments")
	}
	return investments, nil
}

func (investment_serv *investmentsService) CreateInvestment(investment entity.InvestmentsRequest) (entity.Investments, error) {
	userID, err := helper.ParseUUID(investment.UserID)
	if err != nil {
		return entity.Investments{}, errors.New("invalid user id")
	}

	investmentTypeID, err := helper.ParseUUID(investment.InvestmentTypeID)
	if err != nil {
		return entity.Investments{}, err
	}

	newInvestment, err := investment_serv.investmentsRepository.CreateInvestment(entity.Investments{
		UserID:           userID,
		InvestmentTypeID: investmentTypeID,
		Name:             investment.Name,
		Amount:           investment.Amount,
		Quantity:         investment.Quantity,
		InvestmentDate:   investment.InvestmentDate,
		Description:      investment.Description,
	})
	if err != nil {
		return entity.Investments{}, err
	}

	return newInvestment, nil
}

func (investment_serv *investmentsService) UpdateInvestment(id string, investment entity.InvestmentsRequest) (entity.Investments, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(id)
	if err != nil {
		return entity.Investments{}, errors.New("investment not found")
	}

	if investment.Name != "" {
		existingInvestment.Name = investment.Name
	}
	if investment.Amount != 0 {
		existingInvestment.Amount = investment.Amount
	}
	if investment.Quantity != 0 {
		existingInvestment.Quantity = investment.Quantity
	}
	if !investment.InvestmentDate.IsZero() {
		existingInvestment.Description = investment.Description
	}

	investmentUpdated, err := investment_serv.investmentsRepository.UpdateInvestment(existingInvestment)
	if err != nil {
		return entity.Investments{}, errors.New("failed to update investment")
	}

	return investmentUpdated, nil
}

func (investment_serv *investmentsService) DeleteInvestment(id string) (entity.Investments, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(id)
	if err != nil {
		return entity.Investments{}, errors.New("investment not found")
	}

	investmentDeleted, err := investment_serv.investmentsRepository.DeleteInvestment(existingInvestment)
	if err != nil {
		return entity.Investments{}, errors.New("failed to delete investment")
	}

	return investmentDeleted, nil
}
