package repository

import (
	"context"
	"errors"

	"server/internal/helper/data"
	"server/internal/types/entity"

	"gorm.io/gorm"
)

type ReportsRepository interface {
	GetAllReports(ctx context.Context, tx Transaction) ([]entity.Reports, error)
	GetReportByID(ctx context.Context, tx Transaction, id string) (entity.Reports, error)
	GetReportByUserID(ctx context.Context, tx Transaction, user_id string) ([]entity.Reports, error)
	GetProcessedReportByUserID(ctx context.Context, tx Transaction, user_id string) ([]entity.Reports, error)
	CreateReport(ctx context.Context, tx Transaction, report entity.Reports) (entity.Reports, error)
	UpdateReport(ctx context.Context, tx Transaction, report entity.Reports) (entity.Reports, error)
	DeleteReport(ctx context.Context, tx Transaction, report entity.Reports) (entity.Reports, error)
}
type reportsRepository struct {
	db *gorm.DB
}

func NewReportsRepository(db *gorm.DB) ReportsRepository {
	return &reportsRepository{db}
}

// Helper to get DB instance (transactional or regular)
func (reports_repo *reportsRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx) // Type assertion ke GORM transaction
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return reports_repo.db.WithContext(ctx), nil
}

func (reports_repo *reportsRepository) GetAllReports(ctx context.Context, tx Transaction) ([]entity.Reports, error) {
	db, err := reports_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var reports []entity.Reports
	if err := db.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (reports_repo *reportsRepository) GetReportByID(ctx context.Context, tx Transaction, id string) (entity.Reports, error) {
	db, err := reports_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Reports{}, err
	}

	var report entity.Reports
	if err := db.Where("id = ?", id).First(&report).Error; err != nil {
		return entity.Reports{}, err
	}
	return report, nil
}

func (reports_repo *reportsRepository) GetReportByUserID(ctx context.Context, tx Transaction, user_id string) ([]entity.Reports, error) {
	db, err := reports_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var reports []entity.Reports
	if err := db.Where("user_id = ?", user_id).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (reports_repo *reportsRepository) GetProcessedReportByUserID(ctx context.Context, tx Transaction, user_id string) ([]entity.Reports, error) {
	db, err := reports_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var reports []entity.Reports
	if err := db.Where("user_id = ? AND status = ?", user_id, data.REPORT_STATUS_PROCESSING).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (reports_repo *reportsRepository) CreateReport(ctx context.Context, tx Transaction, report entity.Reports) (entity.Reports, error) {
	db, err := reports_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Reports{}, err
	}

	if err := db.Create(&report).Error; err != nil {
		return entity.Reports{}, err
	}
	return report, nil
}

func (reports_repo *reportsRepository) UpdateReport(ctx context.Context, tx Transaction, report entity.Reports) (entity.Reports, error) {
	db, err := reports_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Reports{}, err
	}

	if err := db.Save(&report).Error; err != nil {
		return entity.Reports{}, err
	}
	return report, nil
}

func (reports_repo *reportsRepository) DeleteReport(ctx context.Context, tx Transaction, report entity.Reports) (entity.Reports, error) {
	db, err := reports_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Reports{}, err
	}

	if err := db.Delete(&report).Error; err != nil {
		return entity.Reports{}, err
	}
	return report, nil
}
