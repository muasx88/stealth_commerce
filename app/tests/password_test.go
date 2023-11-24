package test

import (
	"testing"

	"github.com/muasx88/stealth_commerce/app/utils"
)

func TestHashPassword(t *testing.T) {
	password := "admin123"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hashedPassword)
}
