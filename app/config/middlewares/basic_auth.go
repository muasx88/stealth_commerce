package middlewares

import (
	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/internals/response"

	"github.com/gin-gonic/gin"
)

func BasicAuth(u domain.BuyerUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := c.GetHeader("secret_key")
		if secretKey == "" {
			response.Unauthorized(c, "secret_key is not provided")
			return
		}

		buyer, err := u.DetailBySecretKey(c.Request.Context(), secretKey)
		if err != nil {
			response.Unauthorized(c, "secret_key invalid")
			return
		}

		//TODO: store to cache
		//

		c.Set("buyer", buyer)
		c.Next()
	}
}
