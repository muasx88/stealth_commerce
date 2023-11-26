package repository

import (
	"context"

	"github.com/muasx88/stealth_commerce/app/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) domain.OrderRepository {
	return &orderRepository{db: db}
}

// CreateOrder implements domain.OrderRepository.
func (r orderRepository) CreateOrder(ctx context.Context, buyerId int64, payload domain.Order) (res domain.OrderResponse, err error) {
	tx, err := r.db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		log.Error("error while begin transaction", err.Error())
		return res, err
	}

	query := `
		INSERT INTO orders (order_number, buyer_id, grand_total, created_date)
		VALUES (
			@order_number,
			@buyer_id,
			@grand_total,
			CURRENT_TIMESTAMP
		)
		RETURNING id, order_number, buyer_id, grand_total, status, created_date;
	`

	var order domain.OrderResponse
	err = tx.QueryRow(ctx, query,
		pgx.NamedArgs{
			"order_number": payload.OrderNumber,
			"buyer_id":     buyerId,
			"grand_total":  payload.GrandTotal,
		},
	).Scan(
		&order.Id,
		&order.OrderNumber,
		&order.BuyerId,
		&order.GrandTotal,
		&order.Status,
		&order.CreatedDate,
	)

	if err != nil {
		log.Error("error save order data", err.Error())
		return res, err
	}

	query = `
		INSERT INTO order_details (order_id, product_id, qty, price, total)
		VALUES (
			@order_id,
			@product_id,
			@qty,
			@price,
			@total
		)
	`
	for _, detail := range payload.Detail {
		_, err = tx.Exec(ctx, query,
			pgx.NamedArgs{
				"order_id":   order.Id,
				"product_id": detail.ProductId,
				"qty":        detail.Qty,
				"price":      detail.Price,
				"total":      detail.Total,
			},
		)

		if err != nil {
			log.Error("error save order detail data", err.Error())
			return res, err
		}
	}

	tx.Commit(ctx)

	log.Info("success store order data to database")

	res = order
	return
}

func (r orderRepository) DetailByOrderNumber(ctx context.Context, buyerId int64, orderNumber string) (res domain.OrderResponse, err error) {
	query := `
		SELECT id, order_number, grand_total, status, created_date
		FROM orders
		WHERE buyer_id = $1 AND order_number = $2
	`

	var result domain.OrderResponse
	err = r.db.QueryRow(ctx, query, buyerId, orderNumber).
		Scan(
			&result.Id,
			&result.OrderNumber,
			&result.GrandTotal,
			&result.Status,
			&result.CreatedDate,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, domain.ErrOrderNotFound
		} else {
			return res, err
		}
	}

	res = result
	return
}

func (r orderRepository) UpdateStatus(ctx context.Context, buyerId, orderId int64, status string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE buyer_id = $2 and id = $3
	`

	_, err := r.db.Exec(
		ctx,
		query,
		status,
		buyerId,
		orderId,
	)

	if err != nil {
		log.Errorf("error update order status: %s", err.Error())
		return err
	}

	log.Info("success update status order")
	return nil
}
