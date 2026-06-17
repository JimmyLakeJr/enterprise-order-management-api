package service

import (
	"context"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5"
)

type CategoryService interface {
	Create(ctx context.Context, req dto.CategoryRequest) (*dto.CategoryResponse, error)
	List(ctx context.Context) ([]dto.CategoryResponse, error)
	FindByID(ctx context.Context, id int64) (*dto.CategoryResponse, error)
	Update(ctx context.Context, id int64, req dto.CategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id int64) error
}

type categoryService struct {
	categories repository.CategoryRepository
}

func NewCategoryService(categories repository.CategoryRepository) CategoryService {
	return &categoryService{categories: categories}
}

func (s *categoryService) Create(ctx context.Context, req dto.CategoryRequest) (*dto.CategoryResponse, error) {
	exists, err := s.categories.ExistsByNameOtherCategory(ctx, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.Conflict("Category name already exists")
	}

	category := categoryFromRequest(req)
	if err := s.categories.Create(ctx, category); err != nil {
		return nil, err
	}
	res := ToCategoryResponse(category)
	return &res, nil
}

func (s *categoryService) List(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.categories.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.CategoryResponse, 0, len(categories))
	for i := range categories {
		res = append(res, ToCategoryResponse(&categories[i]))
	}
	return res, nil
}

func (s *categoryService) FindByID(ctx context.Context, id int64) (*dto.CategoryResponse, error) {
	category, err := s.categories.FindActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, apperror.NotFound("Category not found")
	}

	res := ToCategoryResponse(category)
	return &res, nil
}

func (s *categoryService) Update(ctx context.Context, id int64, req dto.CategoryRequest) (*dto.CategoryResponse, error) {
	current, err := s.categories.FindActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, apperror.NotFound("Category not found")
	}

	exists, err := s.categories.ExistsByNameOtherCategory(ctx, req.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.Conflict("Category name already exists")
	}

	category := categoryFromRequest(req)
	category.ID = id
	if err := s.categories.Update(ctx, category); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.NotFound("Category not found")
		}
		return nil, err
	}
	category, err = s.categories.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToCategoryResponse(category)
	return &res, nil
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	if err := s.categories.SoftDelete(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return apperror.NotFound("Category not found")
		}
		return err
	}
	return nil
}

func categoryFromRequest(req dto.CategoryRequest) *model.Category {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	return &model.Category{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    isActive,
	}
}

func ToCategoryResponse(category *model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}
}
