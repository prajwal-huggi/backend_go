package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user UserModel) error
	GetUser(ctx context.Context, id int) (UserModel, error)
	GetUsers(ctx context.Context) ([]UserModel, error)
	UpdateUser(ctx context.Context, user UserModel) error
	DeleteUser(ctx context.Context, id int) error
}
