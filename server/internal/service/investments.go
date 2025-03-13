package service

import (
	"context"
	"errors"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/repository"
)

type InvestmentsService interface {
	GetAllInvestments(ctx context.Context) ([]entity.Investments, error)
	GetInvestmentByID(ctx context.Context, id string) (entity.Investments, error)
	GetInvestmentsByUserID(ctx context.Context, id string) ([]entity.Investments, error)
	CreateInvestment(ctx context.Context, investment entity.InvestmentsRequest) (entity.Investments, error)
	UpdateInvestment(ctx context.Context, id string, investment entity.InvestmentsRequest) (entity.Investments, error)
	DeleteInvestment(ctx context.Context, id string) (entity.Investments, error)
}

type investmentsService struct {
	txManager             repository.TxManager
	investmentsRepository repository.InvestmentsRepository
}

func NewInvestmentService(txManager repository.TxManager, investmentsRepository repository.InvestmentsRepository) InvestmentsService {
	return &investmentsService{
		txManager:             txManager,
		investmentsRepository: investmentsRepository,
	}
}

func (investment_serv *investmentsService) GetAllInvestments(ctx context.Context) ([]entity.Investments, error) {
	investments, err := investment_serv.investmentsRepository.GetAllInvestments(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to get investments")
	}

	return investments, nil
}

func (investment_serv *investmentsService) GetInvestmentByID(ctx context.Context, id string) (entity.Investments, error) {
	investment, err := investment_serv.investmentsRepository.GetInvestmentByID(ctx, nil, id)
	if err != nil {
		return entity.Investments{}, errors.New("investment not found")
	}

	return investment, nil
}

func (investment_serv *investmentsService) GetInvestmentsByUserID(ctx context.Context, id string) ([]entity.Investments, error) {
	investments, err := investment_serv.investmentsRepository.GetInvestmentsByUserID(ctx, nil, id)
	if err != nil {
		return nil, errors.New("failed to get investments")
	}
	return investments, nil
}

func (investment_serv *investmentsService) CreateInvestment(ctx context.Context, investment entity.InvestmentsRequest) (entity.Investments, error) {
	userID, err := helper.ParseUUID(investment.UserID)
	if err != nil {
		return entity.Investments{}, errors.New("invalid user id")
	}

	investmentTypeID, err := helper.ParseUUID(investment.InvestmentTypeID)
	if err != nil {
		return entity.Investments{}, err
	}

	newInvestment, err := investment_serv.investmentsRepository.CreateInvestment(ctx, nil, entity.Investments{
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

func (investment_serv *investmentsService) UpdateInvestment(ctx context.Context, id string, investment entity.InvestmentsRequest) (entity.Investments, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(ctx, nil, id)
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

	investmentUpdated, err := investment_serv.investmentsRepository.UpdateInvestment(ctx, nil, existingInvestment)
	if err != nil {
		return entity.Investments{}, errors.New("failed to update investment")
	}

	return investmentUpdated, nil
}

func (investment_serv *investmentsService) DeleteInvestment(ctx context.Context, id string) (entity.Investments, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(ctx, nil, id)
	if err != nil {
		return entity.Investments{}, errors.New("investment not found")
	}

	investmentDeleted, err := investment_serv.investmentsRepository.DeleteInvestment(ctx, nil, existingInvestment)
	if err != nil {
		return entity.Investments{}, errors.New("failed to delete investment")
	}

	return investmentDeleted, nil
}
