package domain

import (
	"context"
	"time"
)

var (
	ORDER_PENDING   = "PENDING"
	ORDER_COMPLETED = "COMPLETED"
	ORDER_CANCELED  = "CANCELED"
)

type OrderPayload struct {
	CartIds []int64 `json:"cart_ids" validate:"min=1,dive,required"`
}

type OrderResponse struct {
	Id          int64     `json:"id"`
	BuyerId     int64     `json:"buyer_id"`
	OrderNumber string    `json:"order_number"`
	GrandTotal  int       `json:"grand_total"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
}

type Order struct {
	OrderNumber string        `json:"order_number"`
	BuyerId     int64         `json:"buyer_id"`
	GrandTotal  int           `json:"grand_total"`
	Detail      []OrderDetail `json:"detail"`
}

type OrderDetail struct {
	ProductId int64 `json:"product_id"`
	Qty       int64 `json:"qty"`
	Price     int   `json:"price"`
	Total     int   `json:"total"`
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, buyerId int64, payload OrderPayload) (OrderResponse, error)
	DetailByOrderNumber(ctx context.Context, buyerId int64, orderNumber string) (OrderResponse, error)
	UpdateStatus(ctx context.Context, buyerId int64, orderId int64, status string) error
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, buyerId int64, payload Order) (OrderResponse, error)
	DetailByOrderNumber(ctx context.Context, buyerId int64, orderNumber string) (OrderResponse, error)
	UpdateStatus(ctx context.Context, buyerId int64, orderId int64, status string) error
}
