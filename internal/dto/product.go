package dto

type ProductRequest struct {
	CategoryID  int64  `json:"category_id" validate:"required,gt=0"`
	Name        string `json:"name" validate:"required,min=2,max=150"`
	Description string `json:"description" validate:"max=1000"`
	Price       int64  `json:"price" validate:"gte=0"`
	Stock       int    `json:"stock" validate:"gte=0"`
	ImageURL    string `json:"image_url" validate:"omitempty,url,max=1000"`
	IsActive    *bool  `json:"is_active"`
}

type ProductResponse struct {
	ID          int64             `json:"id"`
	CategoryID  int64             `json:"category_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       int64             `json:"price"`
	Stock       int               `json:"stock"`
	ImageURL    string            `json:"image_url"`
	IsActive    bool              `json:"is_active"`
	Category    *CategoryResponse `json:"category,omitempty"`
}

type ProductListQuery struct {
	Page       int
	Limit      int
	Search     string
	CategoryID int64
	MinPrice   int64
	MaxPrice   int64
}
