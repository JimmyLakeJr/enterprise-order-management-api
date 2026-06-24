package service

import (
	"context"
	"math"
	"mime/multipart"
	"net/url"
	"strings"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/repository"
	"enterprise-order-management-api/internal/storage"

	"github.com/jackc/pgx/v5"
)

type ProductService interface {
	Create(ctx context.Context, req dto.ProductRequest) (*dto.ProductResponse, error)
	FindByID(ctx context.Context, id int64) (*dto.ProductResponse, error)
	List(ctx context.Context, query dto.ProductListQuery) ([]dto.ProductResponse, response.Meta, error)
	AdminList(ctx context.Context, query dto.ProductListQuery) ([]dto.ProductResponse, response.Meta, error)
	UploadImage(ctx context.Context, file *multipart.FileHeader) (*dto.UploadResponse, error)
	Update(ctx context.Context, id int64, req dto.ProductRequest) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) (*dto.ProductResponse, error)
}

type productService struct {
	products   repository.ProductRepository
	categories repository.CategoryRepository
	storage    *storage.LocalFileStorage
}

func NewProductService(products repository.ProductRepository, categories repository.CategoryRepository, fileStorage *storage.LocalFileStorage) ProductService {
	return &productService{products: products, categories: categories, storage: fileStorage}
}

func (s *productService) Create(ctx context.Context, req dto.ProductRequest) (*dto.ProductResponse, error) {
	if err := validateProductRequest(req); err != nil {
		return nil, err
	}
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
	product, err := s.products.FindActiveByID(ctx, id)
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

func (s *productService) AdminList(ctx context.Context, query dto.ProductListQuery) ([]dto.ProductResponse, response.Meta, error) {
	if query.Status == "" {
		query.Status = "all"
	}
	if query.Status != "all" && query.Status != "active" && query.Status != "inactive" {
		return nil, response.Meta{}, apperror.BadRequest("Invalid product status filter")
	}
	query.Admin = true
	return s.List(ctx, query)
}

func (s *productService) Update(ctx context.Context, id int64, req dto.ProductRequest) (*dto.ProductResponse, error) {
	if err := validateProductRequest(req); err != nil {
		return nil, err
	}
	if err := s.ensureActiveCategory(ctx, req.CategoryID); err != nil {
		return nil, err
	}

	existing, err := s.products.FindActiveByID(ctx, id)
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

func (s *productService) UploadImage(ctx context.Context, file *multipart.FileHeader) (*dto.UploadResponse, error) {
	if s.storage == nil {
		return nil, apperror.New(500, "UPLOAD_DISABLED", "Upload storage is not configured")
	}
	uploadedURL, err := s.storage.SaveImage(file, "products/images", 5*1024*1024)
	if err != nil {
		return nil, err
	}
	return &dto.UploadResponse{URL: uploadedURL}, nil
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

func (s *productService) Restore(ctx context.Context, id int64) (*dto.ProductResponse, error) {
	product, err := s.products.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil || product.IsActive {
		return nil, apperror.NotFound("Inactive product not found")
	}
	category, err := s.categories.FindActiveByID(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, apperror.Conflict("Hãy khôi phục danh mục trước khi khôi phục sản phẩm")
	}
	if err := s.products.Restore(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.NotFound("Inactive product not found")
		}
		return nil, err
	}
	product, err = s.products.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToProductResponse(product)
	return &res, nil
}

func (s *productService) ensureActiveCategory(ctx context.Context, categoryID int64) error {
	category, err := s.categories.FindActiveByID(ctx, categoryID)
	if err != nil {
		return err
	}
	if category == nil {
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

func validateProductRequest(req dto.ProductRequest) error {
	if req.Price < 0 {
		return apperror.BadRequest("Product price must be greater than or equal to 0")
	}
	if req.Stock < 0 {
		return apperror.BadRequest("Product stock must be greater than or equal to 0")
	}
	if req.CategoryID <= 0 {
		return apperror.BadRequest("Category is required")
	}
	if imageURL := strings.TrimSpace(req.ImageURL); imageURL != "" {
		if strings.HasPrefix(imageURL, "/uploads/") {
			return nil
		}
		parsed, err := url.ParseRequestURI(imageURL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return apperror.BadRequest("Product image URL is invalid")
		}
	}
	return nil
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
