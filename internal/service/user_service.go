package service

import (
	"context"
	"fmt"

	"github.com/prajwal-huggi/backend_go/internal/domain"
)

type UserService struct{
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService{
	return &UserService{repo: repo}
}

func (s *UserService)CreateUser(ctx context.Context, user domain.UserModel) error{
	if user.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService)GetUser(ctx context.Context, id int) (domain.UserModel, error){
	return s.repo.GetUser(ctx, id)
}

func (s *UserService)GetUsers(ctx context.Context) ([]domain.UserModel, error){
	return s.repo.GetUsers(ctx)
}

func (s *UserService)UpdateUser(ctx context.Context, user domain.UserModel) error{
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService)DeleteUser(ctx context.Context, id int) error{
	return s.repo.DeleteUser(ctx, id)
}
