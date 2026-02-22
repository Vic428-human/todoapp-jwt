package repository

import (
	"context"
	"time"
	"todo_api/internal/models"
	"todo_api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 引用 "todo_api/internal/utils"

func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	utils.PerformOperation(ctx)

	var query string = `
		INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at, updated_at
	`
	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	utils.PerformOperation(ctx)

	var query string = `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
