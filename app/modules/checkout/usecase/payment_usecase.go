package usecase

import (
	"context"
	"fmt"

	"github.com/muasx88/stealth_commerce/app/domain"
)

type paymentUsecase struct {
	repo         domain.PaymentRepository
	orderUsecase domain.OrderUsecase
}

func NewPaymentUsecase(repo domain.PaymentRepository, orderUsecase domain.OrderUsecase) domain.PaymentUsecase {
	return &paymentUsecase{
		repo:         repo,
		orderUsecase: orderUsecase,
	}
}

func (u paymentUsecase) CreatePayment(ctx context.Context, buyerId int64, payload domain.PaymentPayload) (res domain.Payment, err error) {
	order, err := u.orderUsecase.DetailByOrderNumber(ctx, buyerId, payload.OrderNumber)
	if err != nil {
		return res, err
	}

	if payload.GrandTotal != order.GrandTotal {
		return res, fmt.Errorf("grand total invalid, must equals to %d", order.GrandTotal)
	}

	if order.Status != domain.ORDER_PENDING {
		return res, fmt.Errorf("error order has been %s", order.Status)
	}

	payload.OrderId = order.Id
	res, err = u.repo.CreatePayment(ctx, payload)
	if err != nil {
		return res, err
	}

	// update order status if payment succeed
	u.orderUsecase.UpdateStatus(ctx, buyerId, order.Id, domain.ORDER_COMPLETED)

	return
}
