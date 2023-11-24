package usecase

import (
	"context"

	"github.com/muasx88/stealth_commerce/app/domain"
)

type buyertUsecase struct {
	repo domain.BuyerRepository
}

func NewProductUsecase(buyer domain.BuyerRepository) domain.BuyerUsecase {
	return &buyertUsecase{repo: buyer}
}

func (u buyertUsecase) Detail(ctx context.Context, id int64) (res domain.Buyer, err error) {
	res, err = u.repo.Detail(ctx, id)
	if err != nil {
		return res, err
	}

	return
}

func (u buyertUsecase) DetailBySecretKey(ctx context.Context, secretKey string) (res domain.Buyer, err error) {
	res, err = u.repo.DetailBySecretKey(ctx, secretKey)
	if err != nil {
		return res, err
	}

	return
}
