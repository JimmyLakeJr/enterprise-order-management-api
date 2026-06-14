package service

import (
	"context"
	"math"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5"
)

type ProductService interface {
	Create(ctx context.Context, req dto.ProductRequest) (*dto.ProductResponse, error)
	FindByID(ctx context.Context, id int64) (*dto.ProductResponse, error)
	List(ctx context.Context, query dto.ProductListQuery) ([]dto.ProductResponse, response.Meta, error)
	Update(ctx context.Context, id int64, req dto.ProductRequest) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id int64) error
}

type productService struct {
	products   repository.ProductRepository
	categories repository.CategoryRepository
}

func NewProductService(products repository.ProductRepository, categories repository.CategoryRepository) ProductService {
	return &productService{products: products, categories: categories}
}

func (s *productService) Create(ctx context.Context, req dto.ProductRequest) (*dto.ProductResponse, error) {
	if err := s.ensureActiveCategory(ctx, req.CategoryID); err != nil {
		return nil, err
	}

	product := productFromRequest(req)
	if err := s.products.Create(ctx, product); err != nil {
		return nil, err
	}

	created, err := s.products.FindByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}
	response := ToProductResponse(created)
	return &response, nil
}

func (s *productService) FindByID(ctx context.Context, id int64) (*dto.ProductResponse, error) {
	product, err := s.products.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, apperror.NotFound("Product not found")
	}

	response := ToProductResponse(product)
	return &response, nil
}

func (s *productService) List(ctx context.Context, query dto.ProductListQuery) ([]dto.ProductResponse, response.Meta, error) {
	normalizeProductListQuery(&query)

	products, total, err := s.products.List(ctx, query)
	if err != nil {
		return nil, response.Meta{}, err
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		responses = append(responses, ToProductResponse(&products[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))
	return responses, response.Meta{
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) Update(ctx context.Context, id int64, req dto.ProductRequest) (*dto.ProductResponse, error) {
	if err := s.ensureActiveCategory(ctx, req.CategoryID); err != nil {
		return nil, err
	}

	existing, err := s.products.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, apperror.NotFound("Product not found")
	}

	product := productFromRequest(req)
	product.ID = id

	if err := s.products.Update(ctx, product); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.NotFound("Product not found")
		}
		return nil, err
	}

	updated, err := s.products.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := ToProductResponse(updated)
	return &response, nil
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	if err := s.products.SoftDelete(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return apperror.NotFound("Product not found")
		}
		return err
	}
	return nil
}

func (s *productService) ensureActiveCategory(ctx context.Context, categoryID int64) error {
	category, err := s.categories.FindByID(ctx, categoryID)
	if err != nil {
		return err
	}
	if category == nil || !category.IsActive {
		return apperror.BadRequest("Category is invalid or inactive")
	}
	return nil
}

func normalizeProductListQuery(query *dto.ProductListQuery) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 10
	}
}

func productFromRequest(req dto.ProductRequest) *model.Product {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return &model.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    isActive,
	}
}

func ToProductResponse(product *model.Product) dto.ProductResponse {
	res := dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
		IsActive:    product.IsActive,
	}
	if product.Category != nil {
		category := ToCategoryResponse(product.Category)
		res.Category = &category
	}
	return res
}
