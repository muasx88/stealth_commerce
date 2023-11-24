package domain

import (
	"context"
	"time"
)

type CartCreatePayload struct {
	ProductId int64  `json:"product_id" validate:"required"`
	BuyerId   int64  `json:"buyer_id"`
	Qty       int64  `json:"qty" validate:"required"`
	Note      string `json:"note"`
}

type CartUpdatePayload struct {
	Qty  int64  `json:"qty" validate:"required"`
	Note string `json:"note"`
}

type Cart struct {
	Id          int64     `json:"id"`
	ProductId   int64     `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	BuyerId     int64     `json:"buyer_id"`
	Qty         int64     `json:"qty"`
	Note        string    `json:"note"`
	CreatedDate time.Time `json:"created_date,omitempty"`
	UpdatedDate time.Time `json:"updated_date,omitempty"`
}

type CartUsecase interface {
	GetAll(ctx context.Context, buyerId int64, filter PageQueryString) ([]Cart, error)
	GetByIds(ctx context.Context, ids []int64, buyerId int64) ([]Cart, error)
	Detail(ctx context.Context, id int64, buyerId int64) (Cart, error)
	Add(ctx context.Context, payload CartCreatePayload) (Cart, error)
	Update(ctx context.Context, id int64, buyerId int64, payload CartUpdatePayload) (Cart, error)
	Delete(ctx context.Context, ids []int64, buyerId int64) error
}

type CartRepository interface {
	GetAll(ctx context.Context, buyerId int64, filter PageQueryString) ([]Cart, error)
	GetByIds(ctx context.Context, ids []int64, buyerId int64) ([]Cart, error)
	Detail(ctx context.Context, id int64, buyerId int64) (Cart, error)
	DetailByProduct(ctx context.Context, buyerId int64, productId int64) (Cart, error)
	Add(ctx context.Context, payload CartCreatePayload) (Cart, error)
	Update(ctx context.Context, id int64, buyerId int64, payload CartUpdatePayload) (Cart, error)
	Delete(ctx context.Context, ids []int64, buyerId int64) error
}
