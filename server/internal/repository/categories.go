package repository

import (
	"context"
	"errors"

	"server/internal/entity"

	"gorm.io/gorm"
)

type CategoriesRepository interface {
	GetAllCategories(ctx context.Context, tx Transaction) ([]entity.Categories, error)
	GetCategoryByID(ctx context.Context, tx Transaction, id string) (entity.Categories, error)
	GetCategoriesByType(ctx context.Context, tx Transaction, typeCategory string) ([]entity.Categories, error)
	CreateCategory(ctx context.Context, tx Transaction, category entity.Categories) (entity.Categories, error)
	UpdateCategory(ctx context.Context, tx Transaction, category entity.Categories) (entity.Categories, error)
	DeleteCategory(ctx context.Context, tx Transaction, category entity.Categories) (entity.Categories, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoriesRepository {
	return &categoryRepository{db}
}

// Helper untuk mendapatkan DB instance (transaksi atau biasa)
func (category_repo *categoryRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx) // Type assertion ke GORM transaction
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return category_repo.db.WithContext(ctx), nil
}

func (category_repo *categoryRepository) GetAllCategories(ctx context.Context, tx Transaction) ([]entity.Categories, error) {
	db, err := category_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var categories []entity.Categories
	if err := db.Preload("Parent").Preload("Children").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (category_repo *categoryRepository) GetCategoryByID(ctx context.Context, tx Transaction, id string) (entity.Categories, error) {
	db, err := category_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Categories{}, err
	}

	var category entity.Categories
	if err := db.Preload("Parent").Preload("Children").First(&category, "id = ?", id).Error; err != nil {
		return entity.Categories{}, err
	}

	return category, nil
}

func (category_repo *categoryRepository) GetCategoriesByType(ctx context.Context, tx Transaction, typeCategory string) ([]entity.Categories, error) {
	db, err := category_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var categories []entity.Categories
	if err := db.Preload("Parent").Preload("Children").Where("type = ?", typeCategory).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (category_repo *categoryRepository) CreateCategory(ctx context.Context, tx Transaction, category entity.Categories) (entity.Categories, error) {
	db, err := category_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Categories{}, err
	}

	if err := db.Create(&category).Error; err != nil {
		return entity.Categories{}, err
	}

	return category, nil
}

func (category_repo *categoryRepository) UpdateCategory(ctx context.Context, tx Transaction, category entity.Categories) (entity.Categories, error) {
	db, err := category_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Categories{}, err
	}

	if err := db.Save(&category).Error; err != nil {
		return entity.Categories{}, err
	}

	return category, nil
}

func (category_repo *categoryRepository) DeleteCategory(ctx context.Context, tx Transaction, category entity.Categories) (entity.Categories, error) {
	db, err := category_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Categories{}, err
	}
	
	if err := db.Delete(&category).Error; err != nil {
		return entity.Categories{}, err
	}
	
	return category, nil
}
