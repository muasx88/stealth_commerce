package repository

import (
	"context"
	"fmt"

	"github.com/muasx88/stealth_commerce/app/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type cartRepository struct {
	db *pgxpool.Pool
}

type cartParam struct {
	id        *int64
	buyerId   *int64
	productId *int64
}

func NewCartRepository(db *pgxpool.Pool) domain.CartRepository {
	return &cartRepository{
		db: db,
	}
}

// GetByIds implements domain.CartRepository.
func (r cartRepository) GetByIds(ctx context.Context, ids []int64, buyerId int64) (res []domain.Cart, err error) {
	query := `
		SELECT c.id, c.buyer_id, c.product_id, p.name as product_name, c.qty, c.note, c.created_date, c.updated_date
		FROM carts c
		LEFT JOIN products p on c.product_id = p.id
		WHERE c.id = ANY($1) AND c.buyer_id = $2;
	`

	rows, err := r.db.Query(ctx, query, ids, buyerId)

	if err != nil {
		log.Error("error get carts list", err.Error())
		return res, err
	}
	defer rows.Close()

	cartList := make([]domain.Cart, 0)
	for rows.Next() {
		var cart domain.Cart
		err := rows.Scan(
			&cart.Id,
			&cart.BuyerId,
			&cart.ProductId,
			&cart.ProductName,
			&cart.Qty,
			&cart.Note,
			&cart.CreatedDate,
			&cart.UpdatedDate,
		)

		if err != nil {
			log.Error("Error scan rows", err)
			return res, err
		}

		cartList = append(cartList, cart)
	}

	res = cartList
	return
}

// GetAll implements domain.CartRepository.
func (r cartRepository) GetAll(ctx context.Context, buyerId int64, filter domain.PageQueryString) (res []domain.Cart, err error) {
	var whereClause string

	whereClause += " WHERE c.buyer_id = @buyer_id"
	if filter.Search != "" {
		whereClause += " AND p.name ~* @search"
	}

	query := fmt.Sprintf(`
		SELECT c.id, c.buyer_id, c.product_id, p.name as product_name, c.qty, c.note, c.created_date, c.updated_date
		FROM carts c
		LEFT JOIN products p on c.product_id = p.id
		%s
		ORDER BY c.updated_date DESC
		LIMIT @limit OFFSET @offset;
	`, whereClause)

	rows, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"buyer_id": buyerId,
			"search":   filter.Search,
			"limit":    filter.PageSize,
			"offset":   (filter.Page - 1) * filter.PageSize, // get the offset
		},
	)

	if err != nil {
		log.Error("error get carts list", err.Error())
		return res, err
	}
	defer rows.Close()

	cartList := make([]domain.Cart, 0)
	for rows.Next() {
		var cart domain.Cart
		err := rows.Scan(
			&cart.Id,
			&cart.BuyerId,
			&cart.ProductId,
			&cart.ProductName,
			&cart.Qty,
			&cart.Note,
			&cart.CreatedDate,
			&cart.UpdatedDate,
		)

		if err != nil {
			log.Error("Error scan rows", err)
			return res, err
		}

		cartList = append(cartList, cart)
	}

	res = cartList
	return
}

// getDetail implements domain.CartRepository.
func (r cartRepository) getDetail(ctx context.Context, isById bool, param cartParam) (res domain.Cart, err error) {
	var whereClause string
	var errNotFound error

	whereClause += " WHERE 1=1"
	if isById {
		whereClause += " AND (c.id = @id AND c.buyer_id = @buyer_id)"
		errNotFound = domain.ErrCartNotFound
	} else {
		whereClause += " AND (c.product_id = @product_id AND c.buyer_id = @buyer_id)"
		errNotFound = domain.ErrProductNotFound
	}

	query := fmt.Sprintf(`
		SELECT c.id, c.buyer_id, c.product_id, p.name as product_name, c.qty, c.note, c.created_date, c.updated_date
		FROM carts c
		LEFT JOIN products p on c.product_id = p.id
		%s
	`, whereClause)

	var cart domain.Cart
	err = r.db.QueryRow(ctx, query,
		pgx.NamedArgs{
			"id":         &param.id,
			"buyer_id":   &param.buyerId,
			"product_id": &param.productId,
		},
	).
		Scan(
			&cart.Id,
			&cart.BuyerId,
			&cart.ProductId,
			&cart.ProductName,
			&cart.Qty,
			&cart.Note,
			&cart.CreatedDate,
			&cart.UpdatedDate,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, errNotFound
		} else {
			log.Error("error get cart by product in database", err.Error())
			return res, err
		}
	}

	res = cart
	return
}

// DetailByProduct implements domain.CartRepository.
func (r cartRepository) DetailByProduct(ctx context.Context, buyerId int64, productId int64) (res domain.Cart, err error) {
	param := cartParam{
		productId: &productId,
		buyerId:   &buyerId,
	}

	res, err = r.getDetail(ctx, false, param)
	return
}

// Detail implements domain.CartRepository.
func (r cartRepository) Detail(ctx context.Context, id int64, buyerId int64) (res domain.Cart, err error) {
	param := cartParam{
		id:      &id,
		buyerId: &buyerId,
	}

	res, err = r.getDetail(ctx, true, param)
	return
}

// Add implements domain.CartRepository.
func (r cartRepository) Add(ctx context.Context, payload domain.CartCreatePayload) (res domain.Cart, err error) {
	query := `
		INSERT INTO carts (product_id, buyer_id, qty, note, created_date, updated_date)
		VALUES (
			@product_id,
			@buyer_id,
			@qty,
			@note,
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP
		)
		RETURNING id, product_id, buyer_id, qty, note, created_date, updated_date
	`

	var cart domain.Cart
	err = r.db.QueryRow(ctx, query,
		pgx.NamedArgs{
			"product_id": payload.ProductId,
			"buyer_id":   payload.BuyerId,
			"qty":        payload.Qty,
			"note":       payload.Note,
		},
	).Scan(
		&cart.Id,
		&cart.ProductId,
		&cart.BuyerId,
		&cart.Qty,
		&cart.Note,
		&cart.CreatedDate,
		&cart.UpdatedDate,
	)

	if err != nil {
		log.Error("error add cart data", err.Error())
		return res, err
	}

	log.Info("success add cart data to database")
	res = cart

	return
}

// Update implements domain.CartRepository.
func (r cartRepository) Update(ctx context.Context, id int64, buyerId int64, payload domain.CartUpdatePayload) (res domain.Cart, err error) {
	query := `
		UPDATE carts
		SET qty = $1, note = $2, updated_date = CURRENT_TIMESTAMP
		WHERE id = $3 AND buyer_id = $4
		RETURNING id, product_id, buyer_id, qty, note, created_date, updated_date
	`

	var cart domain.Cart
	err = r.db.QueryRow(ctx, query,
		payload.Qty,
		payload.Note,
		id,
		buyerId,
	).Scan(
		&cart.Id,
		&cart.ProductId,
		&cart.BuyerId,
		&cart.Qty,
		&cart.Note,
		&cart.CreatedDate,
		&cart.UpdatedDate,
	)

	if err != nil {
		log.Error("error update cart data", err.Error())
		return res, err
	}

	log.Info("success update cart data to database")
	res = cart

	return
}

// Delete implements domain.CartRepository.
func (r cartRepository) Delete(ctx context.Context, ids []int64, buyerId int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM carts WHERE id = ANY($1) AND buyer_id = $2", ids, buyerId)
	if err != nil {
		log.Error("error delete carts in database", err.Error())
		return err
	}

	return nil
}
