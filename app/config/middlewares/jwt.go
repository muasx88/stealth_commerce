package middlewares

import (
	"fmt"
	"strings"

	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/internals/config"
	"github.com/muasx88/stealth_commerce/app/internals/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type JWTClaim struct {
	Users domain.Admin
	Exp   int64 `json:"exp"`
	jwt.StandardClaims
}

func JwtProtected() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			log.Error("invalid token header")
			response.Unauthorized(c, "invalid token")
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		secret := []byte(config.Config.JWTKey)
		// validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil {
			log.Error("error token", err.Error())
			response.Unauthorized(c, "invalid token")
			return
		}

		if !token.Valid {
			response.Unauthorized(c, "invalid token")
			return
		}

		claims, ok := token.Claims.(*JWTClaim)
		if !ok || !token.Valid {
			log.Error("error token")
			response.Unauthorized(c, "invalid token")
			return
		}

		c.Set("users", claims)
		c.Next()
	}

}
