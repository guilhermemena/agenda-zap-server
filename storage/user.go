package storage

import (
	"context"

	"github.com/guilhermemena/agenda-zap-server/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	DB *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{DB: db}
}

func (s *UserStorage) Create(ctx context.Context, u *types.User) (*types.User, error) {
	query := `
	INSERT INTO users (first_name, last_name, email, password)
	VALUES ($1, $2, $3, $4)
	RETURNING id, first_name, last_name, email, password, created_at, updated_at
`

	var newUser types.User
	err := s.DB.QueryRow(ctx, query, u.FirstName, u.LastName, u.Email, u.Password).Scan(
		&newUser.ID, &newUser.FirstName, &newUser.LastName, &newUser.Email, &newUser.Password, &newUser.CreatedAt, &newUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}
