package repository

import (
	"context"

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
	INSERT INTO users (name, email, role)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Role)
	return err
}
func (r *UserRepository) GetUser(ctx context.Context, id int) (domain.UserModel, error) {

	query := `SELECT id, name, email, role FROM users WHERE id=$1`

	row := r.db.QueryRow(ctx, query, id)

	var user domain.UserModel

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	return user, err
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]domain.UserModel, error) {

	query := `SELECT id, name, email, role FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.UserModel

	for rows.Next() {
		var user domain.UserModel
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
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
