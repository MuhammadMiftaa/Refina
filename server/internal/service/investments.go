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
	return investment_serv.investmentsRepository.GetAllInvestments()
}

func (investment_serv *investmentsService) GetInvestmentByID(id string) (entity.Investments, error) {
	return investment_serv.investmentsRepository.GetInvestmentByID(id)
}

func (investment_serv *investmentsService) GetInvestmentsByUserID(id string) ([]entity.Investments, error) {
	return investment_serv.investmentsRepository.GetInvestmentsByUserID(id)
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

	return investment_serv.investmentsRepository.CreateInvestment(entity.Investments{
		UserID:            userID,
		InvestmentTypeID: investmentTypeID,
		Name:              investment.Name,
		Amount:            investment.Amount,
		Quantity:          investment.Quantity,
		InvestmentDate:    investment.InvestmentDate,
		Description:       investment.Description,
	})
}

func (investment_serv *investmentsService) UpdateInvestment(id string, investment entity.InvestmentsRequest) (entity.Investments, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(id)
	if err != nil {
		return entity.Investments{}, errors.New("investment not found")
	}

	existingInvestment.Name = investment.Name
	existingInvestment.Amount = investment.Amount
	existingInvestment.Quantity = investment.Quantity
	existingInvestment.Description = investment.Description

	return investment_serv.investmentsRepository.UpdateInvestment(existingInvestment)
}

func (investment_serv *investmentsService) DeleteInvestment(id string) (entity.Investments, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(id)
	if err != nil {
		return entity.Investments{}, errors.New("investment not found")
	}

	return investment_serv.investmentsRepository.DeleteInvestment(existingInvestment)
}
