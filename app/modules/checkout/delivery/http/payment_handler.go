package http

import (
	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/internals/response"
	"github.com/muasx88/stealth_commerce/app/internals/validation"
	"github.com/muasx88/stealth_commerce/app/utils/helper"

	"github.com/gin-gonic/gin"
)

type paymentHandler struct {
	paymentUsecase domain.PaymentUsecase
}

func NewPaymentHandler(u domain.PaymentUsecase) *paymentHandler {
	return &paymentHandler{
		paymentUsecase: u,
	}
}

func (h *paymentHandler) CreatePayment(c *gin.Context) {
	var payload domain.PaymentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ok, errs := validation.ValidateStruct(payload)
	if !ok {
		response.BadRequest(c, "Bad request", errs)
		return
	}

	buyer := helper.GetBuyerSession(c)

	res, err := h.paymentUsecase.CreatePayment(c.Request.Context(), buyer.Id, payload)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.Created(c, "Payment Succeed", res)
}
