package repository

import (
	"context"
	"time"
	"todo_api/internal/models"
	"todo_api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
只要是下面情況都可以放這裡。
1. 操作 users table
2. 查 user
3. 建 user
4. 更新 user 欄位
5. 刪 user
*/
func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	utils.PerformOperation(ctx)

	var query string = `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, created_at, updated_at
	`

	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
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
		SELECT id, email, password, image_url, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	utils.PerformOperation(ctx)

	var query string = `
		SELECT id, email, password, COALESCE(image_url, ''), created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserImage(
	ctx context.Context,
	pool *pgxpool.Pool,
	id string,
	imageURL string,
) (*models.User, error) {
	// 這裡沿用 service 傳進來的 ctx，讓這次「上傳圖片 -> 更新 DB」屬於同一條 request 流程
	utils.PerformOperation(ctx)

	var query string = `
		UPDATE users
		SET image_url = $1,
		    updated_at = NOW()
		WHERE id = $2
		RETURNING id, email, password, image_url, created_at, updated_at
	`

	var user models.User
	err := pool.QueryRow(ctx, query, imageURL, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.ImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
