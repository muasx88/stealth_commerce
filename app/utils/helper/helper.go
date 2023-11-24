package helper

import (
	"math/rand"
	"strings"
	"time"

	"github.com/muasx88/stealth_commerce/app/config/middlewares"
	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/internals/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateJwt(userData domain.Admin) (string, error) {
	secret := config.Config.JWTKey
	expirationTime := time.Now().Add(4 * time.Hour)

	claims := &middlewares.JWTClaim{
		Users: domain.Admin{
			Fullname: userData.Fullname,
		},
		Exp: expirationTime.Unix(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GenerateOrderNumber() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	maxLength := 10
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomStr := make([]byte, maxLength)
	for i := range randomStr {
		randomStr[i] = characters[random.Intn(len(characters))]
	}

	return strings.ToUpper(string(randomStr))
}

func GetBuyerSession(c *gin.Context) domain.Buyer {
	buyer := c.MustGet("buyer").(domain.Buyer)
	return buyer
}

func UniqueIntSlice(slice []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
