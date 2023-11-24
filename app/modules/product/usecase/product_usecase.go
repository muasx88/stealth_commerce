package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/muasx88/stealth_commerce/app/domain"
)

type productUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(product domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{repo: product}
}
func (u productUsecase) ProductList(ctx context.Context, filter domain.PageQueryString) (res []domain.Product, err error) {
	res, err = u.repo.ProductList(ctx, filter)
	if err != nil {
		return res, err
	}

	return
}

func (u productUsecase) ProductDetail(ctx context.Context, id int64) (res domain.Product, err error) {
	res, err = u.repo.ProductDetail(ctx, id)
	if err != nil {
		return res, err
	}

	return
}

func (u productUsecase) AddProduct(ctx context.Context, payload domain.ProductCreatePayload) (res domain.Product, err error) {
	payload.Sku = strings.ToUpper(payload.Sku) // uppercase

	product, _ := u.repo.ProductBySku(ctx, payload.Sku)
	if product != (domain.Product{}) {
		return res, fmt.Errorf("product with sku %s already exists", product.Sku)
	}

	res, err = u.repo.AddProduct(ctx, payload)
	if err != nil {
		return res, err
	}

	return
}

func (u productUsecase) UpdateProductQty(ctx context.Context, id int64, payload int64) error {
	return u.repo.UpdateProductQty(ctx, id, payload)
}

func (u productUsecase) UpdateProduct(ctx context.Context, id int64, payload domain.ProductUpdatePayload) (res domain.Product, err error) {

	product, err := u.repo.ProductDetail(ctx, id)
	if product == (domain.Product{}) {
		return res, errors.New(fmt.Sprintf("product with id %d not found", id))
	}

	res, err = u.repo.UpdateProduct(ctx, id, payload)
	if err != nil {
		return res, err
	}

	return
}

func (u productUsecase) DeleteProduct(ctx context.Context, id int64) (err error) {
	product, err := u.repo.ProductDetail(ctx, id)
	if product == (domain.Product{}) {
		return fmt.Errorf("product with id %d not found", id)
	}

	err = u.repo.DeleteProduct(ctx, id)
	return err
}
