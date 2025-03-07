package repository

import (
	"server/internal/entity"

	"gorm.io/gorm"
)

type CategoriesRepository interface {
	GetAllCategories() ([]entity.Categories, error)
	GetCategoryByID(id string) (entity.Categories, error)
	CreateCategory(category entity.Categories) (entity.Categories, error)
	UpdateCategory(category entity.Categories) (entity.Categories, error)
	DeleteCategory(category entity.Categories) (entity.Categories, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoriesRepository {
	return &categoryRepository{db}
}

func (category_repo *categoryRepository) GetAllCategories() ([]entity.Categories, error) {
	var categories []entity.Categories
	err := category_repo.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (category_repo *categoryRepository) GetCategoryByID(id string) (entity.Categories, error) {
	var category entity.Categories
	err := category_repo.db.First(&category, "id = ?", id).Error
	if err != nil {
		return entity.Categories{}, err
	}

	return category, nil
}

func (category_repo *categoryRepository) CreateCategory(category entity.Categories) (entity.Categories, error) {
	err := category_repo.db.Create(&category).Error
	if err != nil {
		return entity.Categories{}, err
	}
	return category, nil
}

func (category_repo *categoryRepository) UpdateCategory(category entity.Categories) (entity.Categories, error) {
	err := category_repo.db.Save(&category).Error
	if err != nil {
		return entity.Categories{}, err
	}
	return category, nil
}

func (category_repo *categoryRepository) DeleteCategory(category entity.Categories) (entity.Categories, error) {
	err := category_repo.db.Delete(&category).Error
	if err != nil {
		return entity.Categories{}, err
	}
	return category, nil
}
