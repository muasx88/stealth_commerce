package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/muasx88/stealth_commerce/app/domain"
	"github.com/muasx88/stealth_commerce/app/utils/helper"

	log "github.com/sirupsen/logrus"
)

var (
	wg sync.WaitGroup
)

type orderUsecase struct {
	orderRepo       domain.OrderRepository
	cartUsecase     domain.CartUsecase
	prouductUsecase domain.ProductUsecase
	mu              sync.Mutex
	ProductData     []domain.Product
	CardtData       []int64
}

func NewOrderUseOrderUsecase(orderRepo domain.OrderRepository, c domain.CartUsecase, p domain.ProductUsecase) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:       orderRepo,
		cartUsecase:     c,
		prouductUsecase: p,
	}
}

// CreateOrder implements domain.OrderUseOrderUsecase.
func (u *orderUsecase) CreateOrder(ctx context.Context, buyerId int64, payload domain.OrderPayload) (res domain.OrderResponse, err error) {
	var errValiateCart error
	var errOrder error

	//get carts info
	carts, _ := u.getCartsData(ctx, payload.CartIds, buyerId)

	fmt.Println(carts, "the carts")
	if len(carts) < 1 {
		return res, errors.New("cart items empty")
	}

	var order domain.Order
	grandTotal := 0

	for _, cart := range carts {
		wg.Add(1)
		go func(cart domain.Cart) {
			defer wg.Done()

			u.mu.Lock()
			product, err := u.getProductData(ctx, cart.ProductId)
			if err != nil {
				errValiateCart = err
				return
			}

			// validate qty product
			if product.Qty < cart.Qty {
				errValiateCart = err
				return
			}

			total := int(cart.Qty * product.Price)

			item := domain.OrderDetail{
				ProductId: cart.ProductId,
				Qty:       cart.Qty,
				Price:     int(product.Price),
				Total:     total,
			}

			grandTotal += total
			order.Detail = append(order.Detail, item)

			//store the decreased for temporary
			product.Qty -= cart.Qty
			u.ProductData = append(u.ProductData, product)

			//store the cart ids for temp
			u.CardtData = append(u.CardtData, cart.Id)

			u.mu.Unlock()

		}(cart)
	}

	wg.Wait()

	if errValiateCart != nil {
		return res, fmt.Errorf("error create order. %s", errValiateCart.Error())
	}

	order.OrderNumber = helper.GenerateOrderNumber()
	order.BuyerId = buyerId
	order.GrandTotal = grandTotal

	wg.Add(1)
	go func() {
		result, err := u.saveOrder(ctx, buyerId, order, &wg)
		if err != nil {
			errOrder = err
			return
		}

		res = result
	}()
	wg.Wait()

	if errOrder != nil {
		return res, errOrder
	}

	return res, nil
}

func (u *orderUsecase) saveOrder(ctx context.Context, buyerId int64, orderData domain.Order, wg *sync.WaitGroup) (res domain.OrderResponse, err error) {
	defer wg.Done()

	u.mu.Lock()
	res, err = u.orderRepo.CreateOrder(ctx, buyerId, orderData)
	if err != nil {
		log.Error(err.Error())
		return res, errors.New("error save order")
	}

	//updating product qty
	for _, product := range u.ProductData {
		productErr := u.prouductUsecase.UpdateProductQty(ctx, product.Id, product.Qty)
		if productErr != nil {
			//TODO: still need to fix it
			log.Error("error update product qty", productErr.Error())
		}
	}

	// delete the cart
	u.cartUsecase.Delete(ctx, u.CardtData, buyerId)

	u.mu.Unlock()

	log.Info("success save order")
	return
}

func (u *orderUsecase) getProductData(ctx context.Context, productId int64) (res domain.Product, err error) {
	log.Info("getting product data")
	res, err = u.prouductUsecase.ProductDetail(ctx, productId)
	if err != nil {
		log.Error("error get product data", err.Error())
		return res, err
	}

	log.Error("success get product data")
	return
}

func (u *orderUsecase) getCartsData(ctx context.Context, ids []int64, buyerId int64) (res []domain.Cart, err error) {
	log.Info("getting carts data")
	res, err = u.cartUsecase.GetByIds(ctx, ids, buyerId)
	if err != nil {
		log.Error("error get carts data", err.Error())
		return res, err
	}

	log.Info("success get carts data")
	return
}
