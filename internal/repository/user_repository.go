package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prajwal-huggi/backend_go/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.UserModel) error {

	query := `
	INSERT INTO users (name, email, role, password)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Role, user.Password)
	return err
}
func (r *UserRepository) GetUser(ctx context.Context, id int) (domain.UserModel, error) {

	query := `SELECT id, name, email, role, password FROM users WHERE id=$1`

	row := r.db.QueryRow(ctx, query, id)

	var user domain.UserModel

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)
	return user, err
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]domain.UserModel, error) {

	query := `SELECT id, name, email, role, password FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.UserModel

	for rows.Next() {
		var user domain.UserModel
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user domain.UserModel) error {

	query := `
	UPDATE users
	SET name = $1,
	    email = $2
	WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.ID)

	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)

	return err

}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.UserModel, error) {

	query := `SELECT id, name, email, password FROM users WHERE email=$1`

	var user domain.UserModel
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	return user, err
}

func (r *UserRepository) UpdateRefreshToken(
	ctx context.Context,
	userID int,
	token string,
	expiry time.Time,
) error {

	query := `
	UPDATE users
	SET refresh_token = $1,
	    refresh_token_expiry = $2
	WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, token, expiry, userID)

	return err
}
