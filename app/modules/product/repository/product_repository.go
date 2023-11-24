package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/muasx88/stealth_commerce/app/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r productRepository) ProductList(ctx context.Context, filter domain.PageQueryString) (res []domain.Product, err error) {
	var whereClause string

	whereClause += " WHERE 1=1"

	if filter.Search != "" {
		whereClause += " AND name ~* @search"
	}

	query := fmt.Sprintf(`
		SELECT id, name, sku, price, qty, created_date
		FROM products
		%s
		ORDER BY id DESC
		LIMIT @limit OFFSET @offset;
	`, whereClause)

	rows, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"search": filter.Search,
			"limit":  filter.PageSize,
			"offset": (filter.Page - 1) * filter.PageSize, // get the offset
		},
	)

	if err != nil {
		log.Error("error get product list", err.Error())
		return res, err
	}

	defer rows.Close()

	productList := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Sku,
			&product.Price,
			&product.Qty,
			&product.CreatedDate,
		)

		if err != nil {
			log.Error("error scan rows", err.Error())
			return res, err
		}

		productList = append(productList, product)
	}

	res = productList
	return
}

func (r productRepository) ProductBySku(ctx context.Context, sku string) (res domain.Product, err error) {
	query := fmt.Sprintf(`
		SELECT id, name, sku, price, qty, description, created_date
		FROM products
		WHERE sku = $1
	`)

	var result domain.Product
	err = r.db.QueryRow(ctx, query, sku).
		Scan(
			&result.Id,
			&result.Name,
			&result.Sku,
			&result.Price,
			&result.Qty,
			&result.Description,
			&result.CreatedDate,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, domain.ErrProductNotFound
		} else {
			return res, err
		}
	}

	res = result
	return
}

func (r productRepository) ProductDetail(ctx context.Context, id int64) (res domain.Product, err error) {
	query := fmt.Sprintf(`
		SELECT id, name, sku, price, qty, description, created_date
		FROM products
		WHERE id = $1
	`)

	var result domain.Product
	err = r.db.QueryRow(ctx, query, id).
		Scan(
			&result.Id,
			&result.Name,
			&result.Sku,
			&result.Price,
			&result.Qty,
			&result.Description,
			&result.CreatedDate,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, domain.ErrProductNotFound
		} else {
			return res, err
		}
	}

	res = result
	return
}

func (r productRepository) AddProduct(ctx context.Context, payload domain.ProductCreatePayload) (res domain.Product, err error) {
	query := `
		INSERT INTO products (name, sku, price, qty, description, created_date)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, name, sku, price, qty, description, created_date
	`

	var result domain.Product
	err = r.db.QueryRow(
		ctx,
		query,
		payload.Name,
		payload.Sku,
		payload.Price,
		payload.Qty,
		payload.Description,
	).Scan(
		&result.Id,
		&result.Name,
		&result.Sku,
		&result.Price,
		&result.Qty,
		&result.Description,
		&result.CreatedDate,
	)

	if err != nil {
		log.Error("error insert product", err.Error())
		return res, errors.New("error insert product")
	}

	res = result
	return
}

func (r productRepository) UpdateProduct(ctx context.Context, id int64, payload domain.ProductUpdatePayload) (res domain.Product, err error) {
	query := `
		UPDATE products 
		SET name = $1, price = $2, description = $3
		WHERE id = $4
		RETURNING id, name, sku, price, qty, description, created_date
	`
	var result domain.Product
	err = r.db.QueryRow(
		ctx,
		query,
		payload.Name,
		payload.Price,
		payload.Description,
		id,
	).Scan(
		&result.Id,
		&result.Name,
		&result.Sku,
		&result.Price,
		&result.Qty,
		&result.Description,
		&result.CreatedDate,
	)

	if err != nil {
		log.Error("error update product", err.Error())
		return res, errors.New("error update product")
	}

	res = result
	return
}

func (r productRepository) UpdateProductQty(ctx context.Context, id int64, payload int64) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET  qty = $1 WHERE id = $2", payload, id)
	if err != nil {
		log.Error("error update product qty. product id:", id)
		return err
	}

	return nil
}

func (r productRepository) DeleteProduct(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Error("error delete product", err.Error())
		return err
	}

	return nil
}
