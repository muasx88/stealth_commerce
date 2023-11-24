package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/utils"
	"github.com/muasx88/stealth_commerce/app/utils/helper"
)

type authUsecase struct {
	repo domain.AuthRepository
}

func NewAuthUsecase(repo domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{repo: repo}
}

// Login implements domain.AuthUsecase.
func (u authUsecase) Login(ctx context.Context, payload domain.LoginPayload) (res domain.AccessToken, err error) {
	user, err := u.repo.Detail(ctx, payload)
	if err != nil {
		return res, err
	}

	err = utils.CheckPassword(payload.Password, user.Password)
	if err != nil {
		return res, errors.New("invalid email or password")
	}

	token, errToken := helper.GenerateJwt(user)
	if errToken != nil {
		return res, fmt.Errorf("error generate token: %s", errToken.Error())
	}

	res.Type = "Bearer"
	res.Token = token

	return res, nil
}
