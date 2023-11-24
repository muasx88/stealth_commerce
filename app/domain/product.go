package domain

import (
	"context"
	"time"
)

type ProductCreatePayload struct {
	Name        string `json:"name" validate:"required"`
	Sku         string `json:"sku" validate:"required,min=5,max=15"`
	Price       int64  `json:"price" validate:"required,numeric"`
	Qty         int64  `json:"qty" validate:"required,numeric"`
	Description string `json:"description"`
}

type ProductUpdatePayload struct {
	Name        string `json:"name" validate:"required"`
	Price       int64  `json:"price" validate:"required,numeric"`
	Description string `json:"description"`
}

type Product struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Sku         string    `json:"sku" validate:"required,min=5,max=15"`
	Price       int64     `json:"price" validate:"required,numeric"`
	Qty         int64     `json:"qty" validate:"required,numeric"`
	Description string    `json:"description"`
	CreatedDate time.Time `json:"created_date"`
}

type ProductUsecase interface {
	ProductList(ctx context.Context, filter PageQueryString) ([]Product, error)
	ProductDetail(ctx context.Context, id int64) (Product, error)
	AddProduct(ctx context.Context, payload ProductCreatePayload) (Product, error)
	UpdateProduct(ctx context.Context, id int64, payload ProductUpdatePayload) (Product, error)
	UpdateProductQty(ctx context.Context, id int64, payload int64) error
	DeleteProduct(ctx context.Context, id int64) error
}

type ProductRepository interface {
	ProductList(ctx context.Context, filter PageQueryString) ([]Product, error)
	ProductDetail(ctx context.Context, id int64) (Product, error)
	ProductBySku(ctx context.Context, sku string) (Product, error)
	AddProduct(ctx context.Context, payload ProductCreatePayload) (Product, error)
	UpdateProduct(ctx context.Context, id int64, payload ProductUpdatePayload) (Product, error)
	UpdateProductQty(ctx context.Context, id int64, payload int64) error
	DeleteProduct(ctx context.Context, id int64) error
}
