package model

import "time"

type Product struct {
	ID          int64     `db:"id"`
	CategoryID  int64     `db:"category_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       int64     `db:"price"`
	Stock       int       `db:"stock"`
	ImageURL    string    `db:"image_url"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	Category *Category `db:"-"`
}
