package response

import (
	"net/http"

	"github.com/muasx88/stealth_commerce/app/domain"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, message string, v interface{}) {
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: message, Status: http.StatusOK, Data: v})
}

func Created(c *gin.Context, message string, v interface{}) {
	c.JSON(http.StatusCreated, domain.SuccessResponse{Message: message, Status: http.StatusCreated, Data: v})
}

func Deleted(c *gin.Context, message string) {
	c.JSON(http.StatusNoContent, domain.SuccessResponse{Message: message, Status: http.StatusNoContent})
}

func BadRequest(c *gin.Context, msg string, extInfo interface{}) {
	c.JSON(http.StatusBadRequest, domain.ResponseError{Message: msg, Status: http.StatusBadRequest, Errors: extInfo})
}

func Unauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ResponseError{Message: msg, Status: http.StatusUnauthorized})
}

func ServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, domain.ResponseError{Message: msg, Status: http.StatusInternalServerError})
}

func MsgWithCode(c *gin.Context, code int, msg string) {
	c.JSON(code, domain.ResponseError{Message: msg, Status: code})
}
