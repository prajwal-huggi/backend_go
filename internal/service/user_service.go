package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prajwal-huggi/backend_go/internal/auth"
	"github.com/prajwal-huggi/backend_go/internal/domain"
	// "github.com/prajwal-huggi/backend_go/util"
)

type UserService struct{
	repo domain.UserRepository
	jwt *auth.JWTService
}

func NewUserService(repo domain.UserRepository, jwt *auth.JWTService) *UserService{
	return &UserService{
		repo: repo,
		jwt: jwt,
	}
}

func (s *UserService)CreateUser(ctx context.Context, user domain.UserModel) error{
	if user.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.Password = hashedPassword
	
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

func (s *UserService) Login(ctx context.Context, email, password string) (string, string, error) {

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	// fmt.Println("User from DB:", user)

	err = auth.CheckPassword(password, user.Password)
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}
	// fmt.Println("Password matched for user:", user.Email)

	accessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	expiry := time.Now().Add(7 * 24 * time.Hour)

	err = s.repo.UpdateRefreshToken(ctx, user.ID, refreshToken, expiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {

	token, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	newAccessToken, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
