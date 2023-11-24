package domain

type APIResponse struct {
	Status  StatusResponse `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

type SuccessResponse struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Token   *string     `json:"token,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Code    string      `json:"code,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}
type StatusResponse struct {
	StatusCode int      `json:"statusCode"`
	Code       string   `json:"code"`
	Message    string   `json:"message"`
	Detail     []string `json:"detail"`
}
