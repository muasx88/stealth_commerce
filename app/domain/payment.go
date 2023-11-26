package domain

import (
	"context"
	"time"
)

var (
	PAYMENT_SUCCEED = "SUCCEED"
	PAYMENT_FAILED  = "FAILED"
	PAYMENT_ERROR   = "FAILED"
)

type PaymentPayload struct {
	OrderId     int64  `json:"order_id,omitempty"`
	OrderNumber string `json:"order_number" validate:"required"`
	GrandTotal  int    `json:"grand_total" validate:"required,numeric"`
}

type Payment struct {
	Id          int64     `json:"id"`
	OrderId     int64     `json:"order_id"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
}

type PaymentUsecase interface {
	CreatePayment(ctx context.Context, buyerId int64, payload PaymentPayload) (Payment, error)
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payload PaymentPayload) (Payment, error)
}
