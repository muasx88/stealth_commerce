package domain

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Item is not found")
	ErrConflict            = errors.New("your Item already exist")
	ErrBadParamInput       = errors.New("given Param is not valid")
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrBuyerNotFound   = errors.New("buyer not found")
	ErrCartNotFound    = errors.New("cart not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrOrderNotFound   = errors.New("order not found")
)

var ErrorMap = map[error]int{
	ErrInternalServerError: http.StatusInternalServerError,
	ErrNotFound:            http.StatusNotFound,
	ErrProductNotFound:     http.StatusNotFound,
	ErrBuyerNotFound:       http.StatusNotFound,
	ErrCartNotFound:        http.StatusNotFound,
	ErrOrderNotFound:       http.StatusNotFound,
	ErrUserNotFound:        http.StatusNotFound,
	ErrConflict:            http.StatusConflict,
}

type IError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}
