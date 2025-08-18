package service

import (
	"context"

	"server/internal/utils"
	"server/internal/repository"
	"server/internal/types/dto"
	"server/internal/types/entity"
)

type WalletTypesService interface {
	GetAllWalletTypes(ctx context.Context) ([]dto.WalletTypesResponse, error)
	GetWalletTypeByID(ctx context.Context, id string) (dto.WalletTypesResponse, error)
	CreateWalletType(ctx context.Context, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error)
	UpdateWalletType(ctx context.Context, id string, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error)
	DeleteWalletType(ctx context.Context, id string) (dto.WalletTypesResponse, error)
}

type walletTypesService struct {
	txManage        repository.TxManager
	walletTypesRepo repository.WalletTypesRepository
}

func NewWalletTypesService(txManager repository.TxManager, walletTypesRepo repository.WalletTypesRepository) WalletTypesService {
	return &walletTypesService{
		txManage:        txManager,
		walletTypesRepo: walletTypesRepo,
	}
}

func (walletTypeServ *walletTypesService) GetAllWalletTypes(ctx context.Context) ([]dto.WalletTypesResponse, error) {
	walletTypes, err := walletTypeServ.walletTypesRepo.GetAllWalletTypes(ctx, nil)
	if err != nil {
		return nil, err
	}

	var walletTypesResponse []dto.WalletTypesResponse
	for _, walletType := range walletTypes {
		walletTypeResponse := utils.ConvertToResponseType(walletType).(dto.WalletTypesResponse)
		walletTypesResponse = append(walletTypesResponse, walletTypeResponse)
	}

	return walletTypesResponse, nil
}

func (walletTypeServ *walletTypesService) GetWalletTypeByID(ctx context.Context, id string) (dto.WalletTypesResponse, error) {
	walletType, err := walletTypeServ.walletTypesRepo.GetWalletTypeByID(ctx, nil, id)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletType).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}

func (walletTypeServ *walletTypesService) CreateWalletType(ctx context.Context, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error) {
	walletTypeEntity := entity.WalletTypes{
		Name:        walletType.Name,
		Type:        entity.WalletType(walletType.Type),
		Description: walletType.Description,
	}

	walletTypeEntity, err := walletTypeServ.walletTypesRepo.CreateWalletType(ctx, nil, walletTypeEntity)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletTypeEntity).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}

func (walletTypeServ *walletTypesService) UpdateWalletType(ctx context.Context, id string, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error) {
	walletTypeEntity := entity.WalletTypes{
		Name:        walletType.Name,
		Type:        entity.WalletType(walletType.Type),
		Description: walletType.Description,
	}

	walletTypeEntity, err := walletTypeServ.walletTypesRepo.UpdateWalletType(ctx, nil, walletTypeEntity)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletTypeEntity).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}

func (walletTypeServ *walletTypesService) DeleteWalletType(ctx context.Context, id string) (dto.WalletTypesResponse, error) {
	walletTypeEntity, err := walletTypeServ.walletTypesRepo.GetWalletTypeByID(ctx, nil, id)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeEntity, err = walletTypeServ.walletTypesRepo.DeleteWalletType(ctx, nil, walletTypeEntity)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletTypeEntity).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}
