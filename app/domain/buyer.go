package domain

import (
	"context"
	"time"
)

type Buyer struct {
	Id       int64  `json:"id"`
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
	// SecretKey   string    `json:"secret_key,omitempty"`
	CreatedDate time.Time `json:"created_date.omitempty"`
}

type BuyerUsecase interface {
	Detail(ctx context.Context, id int64) (Buyer, error)
	DetailBySecretKey(ctx context.Context, secretKey string) (Buyer, error)
}

type BuyerRepository interface {
	Detail(ctx context.Context, id int64) (Buyer, error)
	DetailBySecretKey(ctx context.Context, secretKey string) (Buyer, error)
}
