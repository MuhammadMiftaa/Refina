package service

import (
	"context"
	"errors"

	"server/internal/repository"
	"server/internal/types/dto"
	"server/internal/types/entity"
	"server/internal/types/view"
	helper "server/internal/utils"
)

type InvestmentsService interface {
	GetAllInvestments(ctx context.Context) ([]dto.InvestmentsResponse, error)
	GetInvestmentByID(ctx context.Context, id string) (dto.InvestmentsResponse, error)
	GetInvestmentsByUserID(ctx context.Context, token string) ([]view.ViewUserInvestments, error)
	CreateInvestment(ctx context.Context, investment dto.InvestmentsRequest) (dto.InvestmentsResponse, error)
	UpdateInvestment(ctx context.Context, id string, investment dto.InvestmentsRequest) (dto.InvestmentsResponse, error)
	DeleteInvestment(ctx context.Context, id string) (dto.InvestmentsResponse, error)
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

func (investment_serv *investmentsService) GetAllInvestments(ctx context.Context) ([]dto.InvestmentsResponse, error) {
	investments, err := investment_serv.investmentsRepository.GetAllInvestments(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to get investments")
	}

	var investmentsResponse []dto.InvestmentsResponse
	for _, investment := range investments {
		investmentResponse := helper.ConvertToResponseType(investment).(dto.InvestmentsResponse)
		investmentsResponse = append(investmentsResponse, investmentResponse)
	}

	return investmentsResponse, nil
}

func (investment_serv *investmentsService) GetInvestmentByID(ctx context.Context, id string) (dto.InvestmentsResponse, error) {
	investment, err := investment_serv.investmentsRepository.GetInvestmentByID(ctx, nil, id)
	if err != nil {
		return dto.InvestmentsResponse{}, errors.New("investment not found")
	}

	investmentResponse := helper.ConvertToResponseType(investment).(dto.InvestmentsResponse)

	return investmentResponse, nil
}

func (investment_serv *investmentsService) GetInvestmentsByUserID(ctx context.Context, token string) ([]view.ViewUserInvestments, error) {
	userData, err := helper.VerifyToken(token[7:])
	if err != nil {
		return nil, errors.New("invalid token")
	}

	investments, err := investment_serv.investmentsRepository.GetInvestmentsByUserID(ctx, nil, userData.ID)
	if err != nil {
		return nil, errors.New("failed to get investments")
	}

	return investments, nil
}

func (investment_serv *investmentsService) CreateInvestment(ctx context.Context, investment dto.InvestmentsRequest) (dto.InvestmentsResponse, error) {
	userID, err := helper.ParseUUID(investment.UserID)
	if err != nil {
		return dto.InvestmentsResponse{}, errors.New("invalid user id")
	}

	investmentTypeID, err := helper.ParseUUID(investment.InvestmentTypeID)
	if err != nil {
		return dto.InvestmentsResponse{}, err
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
		return dto.InvestmentsResponse{}, err
	}

	investmentResponse := helper.ConvertToResponseType(newInvestment).(dto.InvestmentsResponse)

	return investmentResponse, nil
}

func (investment_serv *investmentsService) UpdateInvestment(ctx context.Context, id string, investment dto.InvestmentsRequest) (dto.InvestmentsResponse, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(ctx, nil, id)
	if err != nil {
		return dto.InvestmentsResponse{}, errors.New("investment not found")
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
		return dto.InvestmentsResponse{}, errors.New("failed to update investment")
	}

	investmentResponse := helper.ConvertToResponseType(investmentUpdated).(dto.InvestmentsResponse)

	return investmentResponse, nil
}

func (investment_serv *investmentsService) DeleteInvestment(ctx context.Context, id string) (dto.InvestmentsResponse, error) {
	existingInvestment, err := investment_serv.investmentsRepository.GetInvestmentByID(ctx, nil, id)
	if err != nil {
		return dto.InvestmentsResponse{}, errors.New("investment not found")
	}

	investmentDeleted, err := investment_serv.investmentsRepository.DeleteInvestment(ctx, nil, existingInvestment)
	if err != nil {
		return dto.InvestmentsResponse{}, errors.New("failed to delete investment")
	}

	investmentResponse := helper.ConvertToResponseType(investmentDeleted).(dto.InvestmentsResponse)

	return investmentResponse, nil
}
