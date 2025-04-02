package service

import (
	"context"

	"server/internal/dto"
	"server/internal/entity"
	"server/internal/repository"

	"github.com/google/uuid"
)

type AttachmentsService interface {
	GetAllAttachments(ctx context.Context) ([]dto.AttachmentsResponse, error)
	GetAttachmentByID(ctx context.Context, id string) (dto.AttachmentsResponse, error)
	GetAttachmentsByTransactionID(ctx context.Context, transactionID string) ([]dto.AttachmentsResponse, error)
	CreateAttachment(ctx context.Context, attachment dto.AttachmentsRequest) (dto.AttachmentsResponse, error)
	UpdateAttachment(ctx context.Context, id string, attachment dto.AttachmentsRequest) (dto.AttachmentsResponse, error)
	DeleteAttachment(ctx context.Context, id string) (dto.AttachmentsResponse, error)
}

type attachmentsService struct {
	txManager          repository.TxManager
	categoryRepository repository.AttachmentsRepository
}

func NewAttachmentsService(txManager repository.TxManager, categoryRepository repository.AttachmentsRepository) AttachmentsService {
	return &attachmentsService{
		txManager:          txManager,
		categoryRepository: categoryRepository,
	}
}

func (attachment_serv *attachmentsService) GetAllAttachments(ctx context.Context) ([]dto.AttachmentsResponse, error) {
	attachments, err := attachment_serv.categoryRepository.GetAllAttachments(ctx, nil)
	if err != nil {
		return nil, err
	}

	var attachmentsResponse []dto.AttachmentsResponse
	for _, attachment := range attachments {
		attachmentResponse := dto.AttachmentsResponse{
			ID:            attachment.ID.String(),
			Image:         attachment.Image,
			TransactionID: attachment.TransactionID.String(),
			CreatedAt:     attachment.CreatedAt.String(),
		}
		attachmentsResponse = append(attachmentsResponse, attachmentResponse)
	}

	return attachmentsResponse, nil
}

func (attachment_serv *attachmentsService) GetAttachmentByID(ctx context.Context, id string) (dto.AttachmentsResponse, error) {
	attachment, err := attachment_serv.categoryRepository.GetAttachmentByID(ctx, nil, id)
	if err != nil {
		return dto.AttachmentsResponse{}, err
	}

	attachmentResponse := dto.AttachmentsResponse{
		ID:            attachment.ID.String(),
		Image:         attachment.Image,
		TransactionID: attachment.TransactionID.String(),
		CreatedAt:     attachment.CreatedAt.String(),
	}

	return attachmentResponse, nil
}

func (attachment_serv *attachmentsService) GetAttachmentsByTransactionID(ctx context.Context, transactionID string) ([]dto.AttachmentsResponse, error) {
	attachments, err := attachment_serv.categoryRepository.GetAttachmentsByTransactionID(ctx, nil, transactionID)
	if err != nil {
		return nil, err
	}

	var attachmentsResponse []dto.AttachmentsResponse
	for _, attachment := range attachments {
		attachmentResponse := dto.AttachmentsResponse{
			ID:            attachment.ID.String(),
			Image:         attachment.Image,
			TransactionID: attachment.TransactionID.String(),
			CreatedAt:     attachment.CreatedAt.String(),
		}
		attachmentsResponse = append(attachmentsResponse, attachmentResponse)
	}

	return attachmentsResponse, nil
}

func (attachment_serv *attachmentsService) CreateAttachment(ctx context.Context, attachment dto.AttachmentsRequest) (dto.AttachmentsResponse, error) {
	TransactionUUID, err := uuid.Parse(attachment.TransactionID)
	if err != nil {
		return dto.AttachmentsResponse{}, err
	}

	attachmentEntity := entity.Attachments{
		Image:         attachment.Image,
		TransactionID: TransactionUUID,
	}

	newAttachment, err := attachment_serv.categoryRepository.CreateAttachment(ctx, nil, attachmentEntity)
	if err != nil {
		return dto.AttachmentsResponse{}, err
	}

	attachmentResponse := dto.AttachmentsResponse{
		ID:            newAttachment.ID.String(),
		Image:         newAttachment.Image,
		TransactionID: newAttachment.TransactionID.String(),
		CreatedAt:     newAttachment.CreatedAt.String(),
	}

	return attachmentResponse, nil
}

func (attachment_serv *attachmentsService) UpdateAttachment(ctx context.Context, id string, attachment dto.AttachmentsRequest) (dto.AttachmentsResponse, error) {
	attachmentEntity, err := attachment_serv.categoryRepository.GetAttachmentByID(ctx, nil, id)
	if err != nil {
		return dto.AttachmentsResponse{}, err
	}

	if attachment.Image != "" {
		attachmentEntity.Image = attachment.Image
	}

	newAttachment, err := attachment_serv.categoryRepository.UpdateAttachment(ctx, nil, attachmentEntity)
	if err != nil {
		return dto.AttachmentsResponse{}, err
	}

	attachmentResponse := dto.AttachmentsResponse{
		ID:            newAttachment.ID.String(),
		Image:         newAttachment.Image,
		TransactionID: newAttachment.TransactionID.String(),
		CreatedAt:     newAttachment.CreatedAt.String(),
	}

	return attachmentResponse, nil
}

func (attachment_serv *attachmentsService) DeleteAttachment(ctx context.Context, id string) (dto.AttachmentsResponse, error) {
	attachmentEntity, err := attachment_serv.categoryRepository.GetAttachmentByID(ctx, nil, id)
	if err != nil {
		return dto.AttachmentsResponse{}, err
	}

	deletedAttachment, err := attachment_serv.categoryRepository.DeleteAttachment(ctx, nil, attachmentEntity)
	if err != nil {
		return dto.AttachmentsResponse{}, err
	}

	attachmentResponse := dto.AttachmentsResponse{
		ID:            deletedAttachment.ID.String(),
		Image:         deletedAttachment.Image,
		TransactionID: deletedAttachment.TransactionID.String(),
		CreatedAt:     deletedAttachment.CreatedAt.String(),
	}

	return attachmentResponse, nil
}