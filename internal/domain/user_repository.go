package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user UserModel) error
	GetUser(ctx context.Context, id int) (UserModel, error)
	GetUsers(ctx context.Context) ([]UserModel, error)
	UpdateUser(ctx context.Context, user UserModel) error
	DeleteUser(ctx context.Context, id int) error

	UpdateRefreshToken(ctx context.Context, userID int, token string, expiry time.Time) error
	GetUserByEmail(ctx context.Context, email string) (UserModel, error)
}
