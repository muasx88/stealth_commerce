package http

import (
	"net/http"

	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/internals/response"
	"github.com/muasx88/stealth_commerce/app/internals/validation"

	"github.com/gin-gonic/gin"
)

type handler struct {
	usecase domain.AuthUsecase
}

func NewAuthHandler(usecase domain.AuthUsecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) Login(c *gin.Context) {
	var payload domain.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ok, errs := validation.ValidateStruct(payload)
	if !ok {
		response.BadRequest(c, "Bad request", errs)
		return
	}

	res, err := h.usecase.Login(c.Request.Context(), payload)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		if statusCode == http.StatusNotFound {
			response.MsgWithCode(c, statusCode, err.Error())
			return
		} else {
			response.ServerError(c, err.Error())
			return
		}
	}

	response.OK(c, "login success", res)
}
