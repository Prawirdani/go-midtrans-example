package repository

import (
	"context"
	"database/sql"

	"github.com/prawirdani/go-midtrans-example/internal/entity"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

type userRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(dbConn *sql.DB) UserRepository {
	return &userRepository{
		dbConn: dbConn,
	}
}

func (r *userRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	users := make([]entity.User, 0)

	rows, err := r.dbConn.QueryContext(ctx, selectUserQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	query := selectUserQuery + " WHERE id = ?"
	err := r.dbConn.QueryRowContext(ctx, query, id).
		Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return user, nil
}

const selectUserQuery = "SELECT id, first_name, last_name, email, phone FROM users"
