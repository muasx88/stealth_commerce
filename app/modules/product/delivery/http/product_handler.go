package http

import (
	"strconv"

	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/internals/response"
	"github.com/muasx88/stealth_commerce/app/internals/validation"

	"github.com/gin-gonic/gin"
)

type handler struct {
	usecase domain.ProductUsecase
}

func NewProductHandler(productUsecase domain.ProductUsecase) *handler {
	return &handler{
		usecase: productUsecase,
	}
}

func (h *handler) GetAll(c *gin.Context) {
	search := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	// sortBy := ctx.DefaultQuery("sort_by", "name")
	// sort := ctx.DefaultQuery("sort", "DESC")

	filter := domain.PageQueryString{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	res, err := h.usecase.ProductList(c.Request.Context(), filter)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.OK(c, "products", res)
}

func (h *handler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	product, err := h.usecase.ProductDetail(c.Request.Context(), int64(id))
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.OK(c, "product", product)
}

func (h *handler) Add(c *gin.Context) {
	var payload domain.ProductCreatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ok, errs := validation.ValidateStruct(payload)
	if !ok {
		response.BadRequest(c, "Bad request", errs)
		return
	}

	res, err := h.usecase.AddProduct(c.Request.Context(), payload)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.Created(c, "OK", res)
}

func (h *handler) Update(c *gin.Context) {
	var payload domain.ProductUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	ok, errs := validation.ValidateStruct(payload)
	if !ok {
		response.BadRequest(c, "Bad request", errs)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	res, err := h.usecase.UpdateProduct(c.Request.Context(), int64(id), payload)
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.OK(c, "OK", res)
}

func (h *handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.usecase.DeleteProduct(c.Request.Context(), int64(id))
	if err != nil {
		statusCode := response.GetStatusCode(err)
		response.MsgWithCode(c, statusCode, err.Error())
		return
	}

	response.Deleted(c, "OK")
}
