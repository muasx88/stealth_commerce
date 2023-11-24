package domain

import (
	"context"
	"time"
)

type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AccessToken struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type Admin struct {
	Id          int64     `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password,omitempty"`
	Fullname    string    `json:"fullname"`
	CreatedDate time.Time `json:"created_date"`
}

type AuthUsecase interface {
	Login(ctx context.Context, payload LoginPayload) (AccessToken, error)
}

type AuthRepository interface {
	Detail(ctx context.Context, payload LoginPayload) (Admin, error)
}
