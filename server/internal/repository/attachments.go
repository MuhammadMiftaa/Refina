package repository

import (
	"context"
	"errors"

	"server/internal/types/entity"

	"gorm.io/gorm"
)

type AttachmentsRepository interface {
	GetAllAttachments(ctx context.Context, tx Transaction) ([]entity.Attachments, error)
	GetAttachmentByID(ctx context.Context, tx Transaction, id string) (entity.Attachments, error)
	GetAttachmentsByTransactionID(ctx context.Context, tx Transaction, transactionID string) ([]entity.Attachments, error)
	CreateAttachment(ctx context.Context, tx Transaction, attachment entity.Attachments) (entity.Attachments, error)
	UpdateAttachment(ctx context.Context, tx Transaction, attachment entity.Attachments) (entity.Attachments, error)
	DeleteAttachment(ctx context.Context, tx Transaction, attachment entity.Attachments) (entity.Attachments, error)
}

type attachmentsRepository struct {
	db *gorm.DB
}

func NewAttachmentsRepository(db *gorm.DB) AttachmentsRepository {
	return &attachmentsRepository{db}
}

// Helper untuk mendapatkan DB instance (transaksi atau biasa)
func (attachments_repo *attachmentsRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx) // Type assertion ke GORM transaction
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return attachments_repo.db.WithContext(ctx), nil
}

func (attachments_repo *attachmentsRepository) GetAllAttachments(ctx context.Context, tx Transaction) ([]entity.Attachments, error) {
	db, err := attachments_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var attachments []entity.Attachments
	if err := db.Find(&attachments).Error; err != nil {
		return nil, err
	}

	return attachments, nil
}

func (attachments_repo *attachmentsRepository) GetAttachmentByID(ctx context.Context, tx Transaction, id string) (entity.Attachments, error) {
	db, err := attachments_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Attachments{}, err
	}

	var attachment entity.Attachments
	if err := db.First(&attachment, "id = ?", id).Error; err != nil {
		return entity.Attachments{}, err
	}

	return attachment, nil
}

func (attachments_repo *attachmentsRepository) GetAttachmentsByTransactionID(ctx context.Context, tx Transaction, transactionID string) ([]entity.Attachments, error) {
	db, err := attachments_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var attachments []entity.Attachments
	if err := db.Where("transaction_id = ?", transactionID).Find(&attachments).Error; err != nil {
		return nil, err
	}

	return attachments, nil
}

func (attachments_repo *attachmentsRepository) CreateAttachment(ctx context.Context, tx Transaction, attachment entity.Attachments) (entity.Attachments, error) {
	db, err := attachments_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Attachments{}, err
	}

	if err := db.Create(&attachment).Error; err != nil {
		return entity.Attachments{}, err
	}

	return attachment, nil
}

func (attachments_repo *attachmentsRepository) UpdateAttachment(ctx context.Context, tx Transaction, attachment entity.Attachments) (entity.Attachments, error) {
	db, err := attachments_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Attachments{}, err
	}

	if err := db.Save(&attachment).Error; err != nil {
		return entity.Attachments{}, err
	}

	return attachment, nil
}

func (attachments_repo *attachmentsRepository) DeleteAttachment(ctx context.Context, tx Transaction, attachment entity.Attachments) (entity.Attachments, error) {
	db, err := attachments_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Attachments{}, err
	}

	if err := db.Delete(&attachment).Error; err != nil {
		return entity.Attachments{}, err
	}

	return attachment, nil
}
