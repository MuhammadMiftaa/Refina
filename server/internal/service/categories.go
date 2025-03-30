package service

import (
	"context"

	"server/internal/dto"
	"server/internal/entity"
	"server/internal/repository"

	"github.com/google/uuid"
)

type CategoriesService interface {
	GetAllCategories(ctx context.Context) (map[string][]string, error)
	GetCategoryByID(ctx context.Context, id string) (dto.CategoriesResponse, error)
	GetCategoriesByType(ctx context.Context, typeCategory string) (map[string][]string, error)
	CreateCategory(ctx context.Context, category dto.CategoriesRequest) (dto.CategoriesResponse, error)
	UpdateCategory(ctx context.Context, id string, category dto.CategoriesRequest) (dto.CategoriesResponse, error)
	DeleteCategory(ctx context.Context, id string) (dto.CategoriesResponse, error)
}

type categoriesService struct {
	txManager          repository.TxManager
	categoryRepository repository.CategoriesRepository
}

func NewCategoriesService(txManager repository.TxManager, categoryRepository repository.CategoriesRepository) CategoriesService {
	return &categoriesService{
		txManager:          txManager,
		categoryRepository: categoryRepository,
	}
}

func (category_serv *categoriesService) GetAllCategories(ctx context.Context) (map[string][]string, error) {
	categories, err := category_serv.categoryRepository.GetAllCategories(ctx, nil)
	if err != nil {
		return nil, err
	}

	groupedCategories := make(map[string][]string)
	for _, category := range categories {
		if category.ParentID == nil {
			if _, exists := groupedCategories[category.Name]; !exists {
				groupedCategories[category.Name] = []string{}
			}
		} else {
			if category.Parent != nil {
				groupedCategories[category.Parent.Name] = append(groupedCategories[category.Parent.Name], category.Name)
			}
		}
	}

	return groupedCategories, nil
}

func (category_serv *categoriesService) GetCategoryByID(ctx context.Context, id string) (dto.CategoriesResponse, error) {
	category, err := category_serv.categoryRepository.GetCategoryByID(ctx, nil, id)
	if err != nil {
		return dto.CategoriesResponse{}, err
	}

	var response dto.CategoriesResponse

	if category.Parent != nil {
		response = dto.CategoriesResponse{
			ID:          category.ID.String(),
			Category:    category.Parent.Name,
			SubCategory: category.Name,
			Type:        dto.CategoryType(category.Type),
		}
	} else {
		response = dto.CategoriesResponse{
			ID:          category.ID.String(),
			Category:    category.Name,
			SubCategory: "",
			Type:        dto.CategoryType(category.Type),
		}
	}

	return response, nil
}

func (category_serv *categoriesService) GetCategoriesByType(ctx context.Context, typeCategory string) (map[string][]string, error) {
	categories, err := category_serv.categoryRepository.GetCategoriesByType(ctx, nil, typeCategory)
	if err != nil {
		return nil, err
	}

	groupedCategories := make(map[string][]string)
	for _, category := range categories {
		if category.ParentID == nil {
			if _, exists := groupedCategories[category.Name]; !exists {
				groupedCategories[category.Name] = []string{}
			}
		} else {
			if category.Parent != nil {
				groupedCategories[category.Parent.Name] = append(groupedCategories[category.Parent.Name], category.Name)
			}
		}
	}

	return groupedCategories, nil
}

func (category_serv *categoriesService) CreateCategory(ctx context.Context, category dto.CategoriesRequest) (dto.CategoriesResponse, error) {
	var newCategory entity.Categories
	var err error
	if category.ParentID != "" {
		parent, err := category_serv.categoryRepository.GetCategoryByID(ctx, nil, category.ParentID)
		if err != nil {
			return dto.CategoriesResponse{}, err
		}

		newCategory, err = category_serv.categoryRepository.CreateCategory(ctx, nil, entity.Categories{
			ParentID: &parent.ID,
			Name:     category.Name,
			Type:     entity.CategoryType(category.Type),
		})
		if err != nil {
			return dto.CategoriesResponse{}, err
		}
	} else {
		newCategory, err = category_serv.categoryRepository.CreateCategory(ctx, nil, entity.Categories{
			Name: category.Name,
			Type: entity.CategoryType(category.Type),
		})
		if err != nil {
			return dto.CategoriesResponse{}, err
		}
	}

	var response dto.CategoriesResponse
	if newCategory.Parent != nil {
		response = dto.CategoriesResponse{
			ID:          newCategory.ID.String(),
			Category:    newCategory.Parent.Name,
			SubCategory: newCategory.Name,
			Type:        dto.CategoryType(newCategory.Type),
		}
	} else {
		response = dto.CategoriesResponse{
			ID:          newCategory.ID.String(),
			Category:    newCategory.Name,
			SubCategory: "",
			Type:        dto.CategoryType(newCategory.Type),
		}
	}

	return response, nil
}

func (category_serv *categoriesService) UpdateCategory(ctx context.Context, id string, category dto.CategoriesRequest) (dto.CategoriesResponse, error) {
	existCategory, err := category_serv.categoryRepository.GetCategoryByID(ctx, nil, id)
	if err != nil {
		return dto.CategoriesResponse{}, err
	}

	if category.ParentID != "" {
		_, err := category_serv.categoryRepository.GetCategoryByID(ctx, nil, category.ParentID)
		if err != nil {
			return dto.CategoriesResponse{}, err
		}

		parentUUID, err := uuid.Parse(category.ParentID)
		if err != nil {
			return dto.CategoriesResponse{}, err
		}
		existCategory.ParentID = &parentUUID
	}
	if category.Name != "" {
		existCategory.Name = category.Name
	}
	if category.Type != "" {
		existCategory.Type = entity.CategoryType(category.Type)
	}

	newCategory, err := category_serv.categoryRepository.UpdateCategory(ctx, nil, existCategory)
	if err != nil {
		return dto.CategoriesResponse{}, err
	}
	
	var response dto.CategoriesResponse
	if newCategory.Parent != nil {
		response = dto.CategoriesResponse{
			ID:          newCategory.ID.String(),
			Category:    newCategory.Parent.Name,
			SubCategory: newCategory.Name,
			Type:        dto.CategoryType(newCategory.Type),
		}
	} else {
		response = dto.CategoriesResponse{
			ID:          newCategory.ID.String(),
			Category:    newCategory.Name,
			SubCategory: "",
			Type:        dto.CategoryType(newCategory.Type),
		}
	}

	return response, nil
}

func (category_serv *categoriesService) DeleteCategory(ctx context.Context, id string) (dto.CategoriesResponse, error) {
	existCategory, err := category_serv.categoryRepository.GetCategoryByID(ctx, nil, id)
	if err != nil {
		return dto.CategoriesResponse{}, err
	}

	deletedCategory, err := category_serv.categoryRepository.DeleteCategory(ctx, nil, existCategory)
	if err != nil {
		return dto.CategoriesResponse{}, err
	}

	var response dto.CategoriesResponse
	if deletedCategory.Parent != nil {
		response = dto.CategoriesResponse{
			ID:          deletedCategory.ID.String(),
			Category:    deletedCategory.Parent.Name,
			SubCategory: deletedCategory.Name,
			Type:        dto.CategoryType(deletedCategory.Type),
		}
	} else {
		response = dto.CategoriesResponse{
			ID:          deletedCategory.ID.String(),
			Category:    deletedCategory.Name,
			SubCategory: "",
			Type:        dto.CategoryType(deletedCategory.Type),
		}
	}
	return response, nil
}
