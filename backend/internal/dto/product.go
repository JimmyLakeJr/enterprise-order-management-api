package dto

import "time"

type CreateProductRequest struct {
	CategoryID  int64  `json:"category_id" validate:"required,gt=0"`
	Name        string `json:"name" validate:"required,min=2,max=150"`
	Description string `json:"description" validate:"omitempty,max=2000"`
	Price       *int64 `json:"price" validate:"required,gte=0"`
	Stock       *int   `json:"stock" validate:"required,gte=0"`
	ImageURL    string `json:"image_url" validate:"omitempty,url"`
}

type UpdateProductRequest struct {
	CategoryID  *int64  `json:"category_id" validate:"omitempty,gt=0"`
	Name        *string `json:"name" validate:"omitempty,min=2,max=150"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	Price       *int64  `json:"price" validate:"omitempty,gte=0"`
	Stock       *int    `json:"stock" validate:"omitempty,gte=0"`
	ImageURL    *string `json:"image_url" validate:"omitempty,url"`
	IsActive    *bool   `json:"is_active" validate:"omitempty"`
}

type ProductResponse struct {
	ID          int64             `json:"id"`
	CategoryID  int64             `json:"category_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       int64             `json:"price"`
	Stock       int               `json:"stock"`
	ImageURL    string            `json:"image_url"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
