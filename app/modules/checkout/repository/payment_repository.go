package repository

import (
	"context"

	"github.com/muasx88/stealth_commerce/app/domain"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
)

type paymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

// Detail implements domain.BuyerRepository.
func (r paymentRepository) CreatePayment(ctx context.Context, payload domain.PaymentPayload) (res domain.Payment, err error) {
	query := `
		INSERT INTO PAYMENTS (order_id, status)
		VALUES ($1, $2)
		RETURNING id, order_id, status, created_date
	`
	var payment domain.Payment
	err = r.db.QueryRow(ctx, query, payload.OrderId, domain.PAYMENT_SUCCEED).Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.Status,
		&payment.CreatedDate,
	)

	if err != nil {
		log.Error("error add payment data", err.Error())
		return res, err
	}

	log.Info("success add payment data to database")
	res = payment

	return
}
