package http

import (
	"net/http"
	"strconv"

	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/internals/response"
	"github.com/muasx88/stealth_commerce/app/internals/validation"
	"github.com/muasx88/stealth_commerce/app/utils/helper"

	"github.com/gin-gonic/gin"
)

type handler struct {
	cartUsecase    domain.CartUsecase
	productUsecase domain.ProductUsecase
}

func NewCartHandler(
	cartUsecase domain.CartUsecase,
	productUsecase domain.ProductUsecase,
) *handler {
	return &handler{
		cartUsecase:    cartUsecase,
		productUsecase: productUsecase,
	}
}

func (h *handler) GetAll(c *gin.Context) {
	search := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	ctx := c.Request.Context()

	filter := domain.PageQueryString{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	buyer := helper.GetBuyerSession(c)

	res, err := h.cartUsecase.GetAll(ctx, buyer.Id, filter)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.OK(c, "carts", res)
}

func (h *handler) Add(c *gin.Context) {
	var payload domain.CartCreatePayload
	var err error

	if err = c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ok, errs := validation.ValidateStruct(payload)
	if !ok {
		response.BadRequest(c, "Bad request", errs)
		return
	}

	ctx := c.Request.Context()

	// check product exists
	product, err := h.productUsecase.ProductDetail(ctx, payload.ProductId)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		if statusCode == http.StatusNotFound {
			response.BadRequest(c, "product not found", nil)
			return
		} else {
			response.ServerError(c, err.Error())
			return
		}
	}

	// validate product qty
	if payload.Qty > product.Qty {
		response.BadRequest(c, "invalid quantity", nil)
		return
	}

	buyer := helper.GetBuyerSession(c)
	payload.BuyerId = buyer.Id

	res, err := h.cartUsecase.Add(ctx, payload)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.Created(c, "Cart Created", res)
}

func (h *handler) Update(c *gin.Context) {
	var payload domain.CartUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ok, errs := validation.ValidateStruct(payload)
	if !ok {
		response.BadRequest(c, "Bad request", errs)
		return
	}

	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	buyer := helper.GetBuyerSession(c)
	buyerId := buyer.Id

	// check if cart exists
	cart, err := h.cartUsecase.Detail(ctx, int64(id), buyerId)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	// check product exists
	product, _ := h.productUsecase.ProductDetail(ctx, cart.ProductId)

	// validate product qty
	if payload.Qty > product.Qty {
		response.BadRequest(c, "invalid quantity", nil)
		return
	}

	res, err := h.cartUsecase.Update(ctx, int64(id), buyer.Id, payload)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.OK(c, "OK", res)
}

func (h *handler) Delete(c *gin.Context) {
	var err error

	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	buyer := helper.GetBuyerSession(c)

	_, err = h.cartUsecase.Detail(ctx, int64(id), buyer.Id)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	ids := make([]int64, 1)
	ids = append(ids, int64(id)) // insert id to [id] slice

	err = h.cartUsecase.Delete(ctx, ids, buyer.Id)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.Deleted(c, "OK")
}
