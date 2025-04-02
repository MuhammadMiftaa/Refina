package service

import (
	"context"

	"server/internal/dto"
	"server/internal/entity"
	"server/internal/repository"

	"github.com/google/uuid"
)

type CategoriesService interface {
	GetAllCategories(ctx context.Context) ([]dto.CategoriesResponse, error)
	GetCategoryByID(ctx context.Context, id string) (dto.CategoriesResponse, error)
	GetCategoriesByType(ctx context.Context, typeCategory string) ([]dto.CategoriesResponse, error)
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

func (category_serv *categoriesService) GetAllCategories(ctx context.Context) ([]dto.CategoriesResponse, error) {
	categories, err := category_serv.categoryRepository.GetAllCategories(ctx, nil)
	if err != nil {
		return nil, err
	}

	var groupedCategories []dto.CategoriesResponse
	for _, category := range categories {
		if category.ParentID == nil {
			exists := false
			for _, group := range groupedCategories {
				if group.GroupName == category.Name {
					exists = true
					break
				}
			}
			if !exists {
				groupName := category.Name
				groupedCategories = append(groupedCategories, dto.CategoriesResponse{
					GroupName: groupName,
					Type:      dto.CategoryType(category.Type),
					Category:  []dto.Category{},
				})
			}
		} else {
			if category.Parent != nil {
				for i, group := range groupedCategories {
					if group.GroupName == category.Parent.Name {
						groupedCategories[i].Category = append(groupedCategories[i].Category, dto.Category{
							ID:   category.ID.String(),
							Name: category.Name,
						})
					}
				}
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
		var categories []dto.Category
		categories = append(categories, dto.Category{
			ID:   category.ID.String(),
			Name: category.Name,
		})
		response = dto.CategoriesResponse{
			GroupName: category.Parent.Name,
			Type:      dto.CategoryType(category.Type),
			Category:  categories,
		}
	} else {
		response = dto.CategoriesResponse{
			GroupName: category.Name,
		}
	}

	return response, nil
}

func (category_serv *categoriesService) GetCategoriesByType(ctx context.Context, typeCategory string) ([]dto.CategoriesResponse, error) {
	categories, err := category_serv.categoryRepository.GetCategoriesByType(ctx, nil, typeCategory)
	if err != nil {
		return nil, err
	}

	var groupedCategories []dto.CategoriesResponse
	for _, category := range categories {
		if category.ParentID == nil {
			exists := false
			for _, group := range groupedCategories {
				if group.GroupName == category.Name {
					exists = true
					break
				}
			}
			if !exists {
				groupName := category.Name
				groupedCategories = append(groupedCategories, dto.CategoriesResponse{
					GroupName: groupName,
					Category:  []dto.Category{},
					Type:      dto.CategoryType(category.Type),
				})
			}
		} else {
			if category.Parent != nil {
				for i, group := range groupedCategories {
					if group.GroupName == category.Parent.Name {
						groupedCategories[i].Category = append(groupedCategories[i].Category, dto.Category{
							ID:   category.ID.String(),
							Name: category.Name,
						})
					}
				}
			}
		}
	}

	return groupedCategories, nil
}

func (category_serv *categoriesService) CreateCategory(ctx context.Context, category dto.CategoriesRequest) (dto.CategoriesResponse, error) {
	var newCategory entity.Categories
	var parentName string
	var err error
	if category.ParentID != "" {
		parent, err := category_serv.categoryRepository.GetCategoryByID(ctx, nil, category.ParentID)
		if err != nil {
			return dto.CategoriesResponse{}, err
		}

		parentName = parent.Name

		newCategory, err = category_serv.categoryRepository.CreateCategory(ctx, nil, entity.Categories{
			ParentID: &parent.ID,
			Name:     category.Name,
			Type:     parent.Type,
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
	if category.ParentID != "" {
		var categories []dto.Category
		categories = append(categories, dto.Category{
			ID:   newCategory.ID.String(),
			Name: newCategory.Name,
		})
		response = dto.CategoriesResponse{
			GroupName: parentName,
			Type:      dto.CategoryType(newCategory.Type),
			Category:  categories,
		}
	} else {
		response = dto.CategoriesResponse{
			GroupName: newCategory.Name,
			Type:      dto.CategoryType(newCategory.Type),
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
	if category.Type != "" && category.ParentID == "" {
		existCategory.Type = entity.CategoryType(category.Type)
	}

	newCategory, err := category_serv.categoryRepository.UpdateCategory(ctx, nil, existCategory)
	if err != nil {
		return dto.CategoriesResponse{}, err
	}

	var response dto.CategoriesResponse
	if existCategory.ParentID != nil {
		var categories []dto.Category
		categories = append(categories, dto.Category{
			ID:   newCategory.ID.String(),
			Name: newCategory.Name,
		})
		response = dto.CategoriesResponse{
			GroupName: existCategory.Parent.Name,
			Type:      dto.CategoryType(newCategory.Type),
			Category:  categories,
		}
	} else {
		response = dto.CategoriesResponse{
			GroupName: newCategory.Name,
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
	if existCategory.ParentID != nil {
		var categories []dto.Category
		categories = append(categories, dto.Category{
			ID:   deletedCategory.ID.String(),
			Name: deletedCategory.Name,
		})
		response = dto.CategoriesResponse{
			GroupName: existCategory.Parent.Name,
			Type:      dto.CategoryType(deletedCategory.Type),
			Category:  categories,
		}
	} else {
		response = dto.CategoriesResponse{
			GroupName: deletedCategory.Name,
		}
	}
	return response, nil
}
