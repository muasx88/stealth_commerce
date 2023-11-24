package repository

import (
	"context"
	"fmt"

	"github.com/muasx88/stealth_commerce/app/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) domain.AuthRepository {
	return &authRepository{db: db}
}

// Detail implements domain.AdminRepository.
func (r authRepository) Detail(ctx context.Context, payload domain.LoginPayload) (res domain.Admin, err error) {
	query := fmt.Sprintf(`
		SELECT id, username, password, fullname, created_date
		FROM admin
		WHERE username = $1
	`)

	var result domain.Admin
	err = r.db.QueryRow(ctx, query, payload.Username).
		Scan(
			&result.Id,
			&result.Username,
			&result.Password,
			&result.Fullname,
			&result.CreatedDate,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, domain.ErrUserNotFound
		} else {
			return res, err
		}
	}

	res = result
	return
}
