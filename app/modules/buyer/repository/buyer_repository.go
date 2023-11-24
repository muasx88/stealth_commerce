package repository

import (
	"context"

	"github.com/muasx88/stealth_commerce/app/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type buyerRepository struct {
	db *pgxpool.Pool
}

func NewBuyerRepository(db *pgxpool.Pool) domain.BuyerRepository {
	return &buyerRepository{db: db}
}

// Detail implements domain.BuyerRepository.
func (r *buyerRepository) Detail(ctx context.Context, id int64) (res domain.Buyer, err error) {
	query := `
		SELECT id, fullname, address, created_date
		FROM buyers
		WHERE id = $1
	`
	var buyer domain.Buyer
	err = r.db.QueryRow(ctx, query, id).
		Scan(
			&buyer.Id,
			&buyer.Fullname,
			&buyer.Address,
			&buyer.CreatedDate,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, domain.ErrBuyerNotFound
		} else {
			return res, err
		}
	}

	res = buyer
	return
}

// DetailBySecretKey implements domain.BuyerRepository.
func (r *buyerRepository) DetailBySecretKey(ctx context.Context, secretKey string) (res domain.Buyer, err error) {
	query := `
		SELECT id, fullname, address, created_date
		FROM buyers
		WHERE secret_key = $1
	`
	var buyer domain.Buyer
	err = r.db.QueryRow(ctx, query, secretKey).
		Scan(
			&buyer.Id,
			&buyer.Fullname,
			&buyer.Address,
			&buyer.CreatedDate,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, domain.ErrBuyerNotFound
		} else {
			return res, err
		}
	}

	res = buyer
	return
}
