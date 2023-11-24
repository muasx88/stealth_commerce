package usecase

import (
	"context"

	"github.com/muasx88/stealth_commerce/app/domain"

	log "github.com/sirupsen/logrus"
)

type cartUsecase struct {
	repo domain.CartRepository
}

func NewCartUsecase(repo domain.CartRepository) domain.CartUsecase {
	return &cartUsecase{repo: repo}
}

// GetAll implements domain.CartUsecase.
func (u cartUsecase) GetAll(ctx context.Context, buyerId int64, filter domain.PageQueryString) (res []domain.Cart, err error) {
	res, err = u.repo.GetAll(ctx, buyerId, filter)
	if err != nil {
		return nil, err
	}

	return
}

// GetByIds implements domain.CartUsecase.
func (u cartUsecase) GetByIds(ctx context.Context, ids []int64, buyerId int64) (res []domain.Cart, err error) {
	res, err = u.repo.GetByIds(ctx, ids, buyerId)
	if err != nil {
		return nil, err
	}

	return
}

// Detail implements domain.CartUsecase.
func (u cartUsecase) Detail(ctx context.Context, id int64, buyerId int64) (res domain.Cart, err error) {
	res, err = u.repo.Detail(ctx, id, buyerId)
	if err != nil {
		return res, err
	}

	return
}

// DetailByProduct implements domain.CartUsecase.
func (u cartUsecase) DetailByProduct(ctx context.Context, id int64, buyerId int64) (res domain.Cart, err error) {
	res, err = u.repo.DetailByProduct(ctx, id, buyerId)
	if err != nil {
		return res, err
	}

	return
}

// Add implements domain.CartUsecase.
func (u cartUsecase) Add(ctx context.Context, payload domain.CartCreatePayload) (res domain.Cart, err error) {

	// check if there is existing product in cart
	// then do the update action
	productCart, ok := u.checkProductInCart(ctx, payload.ProductId, payload.BuyerId)
	if ok {
		updatePayload := domain.CartUpdatePayload{
			Qty:  payload.Qty + productCart.Qty,
			Note: payload.Note,
		}

		res, err = u.Update(ctx, productCart.Id, payload.BuyerId, updatePayload)
		return
	}

	res, err = u.repo.Add(ctx, payload)
	if err != nil {
		log.Error("error add qty prodct in cart -> ", err.Error())
		return
	}

	return res, nil
}

// check the existing product in cart
func (u cartUsecase) checkProductInCart(ctx context.Context, productId int64, buyerId int64) (res domain.Cart, ok bool) {
	productCart, _ := u.repo.DetailByProduct(ctx, buyerId, productId)
	if productCart == (domain.Cart{}) {
		return res, false
	}

	res = productCart

	return res, true
}

// Update implements domain.CartUsecase.
func (u cartUsecase) Update(ctx context.Context, id int64, buyerId int64, payload domain.CartUpdatePayload) (res domain.Cart, err error) {
	res, err = u.repo.Update(ctx, id, buyerId, payload)
	if err != nil {
		log.Error("error update qty prodct in cart -> ", err.Error())
		return res, err
	}

	return
}

// Delete implements domain.CartUsecase.
func (u cartUsecase) Delete(ctx context.Context, ids []int64, buyerId int64) error {
	return u.repo.Delete(ctx, ids, buyerId)
}
