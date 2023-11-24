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
