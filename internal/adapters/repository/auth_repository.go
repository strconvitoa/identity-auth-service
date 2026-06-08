package repository

import (
	"context"

	"github.com/strconvitoa/martian-service/internal/core/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuthRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAuthRepository(db *pgxpool.Pool) *PostgresAuthRepository {
	return &PostgresAuthRepository{
		db: db,
	}
}

func (r *PostgresAuthRepository) Insert(user domain.Auth) (domain.Auth, error) {
	query := `
        INSERT INTO auth (id, expires_at, user_id,email, token) 
        VALUES ($1, $2, $3, $4,$5)
    `

	_, err := r.db.Exec(
		context.Background(),
		query,
		user.ID,
		user.ExpiresAt,
		user.UserID,
		user.Email,
		user.Token,
	)
	if err != nil {
		return domain.Auth{}, err
	}

	return user, nil
}

func (r *PostgresAuthRepository) Get(auth domain.Auth) (domain.Auth, error) {
	return domain.Auth{}, nil
}
func (r *PostgresAuthRepository) Find(auth domain.Auth) (domain.Auth, error) {
	return domain.Auth{}, nil
}

func (r *PostgresAuthRepository) Exists(auth domain.Auth) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM auth 
			WHERE email = $1
			AND token = $2
			AND expires_at > NOW()
		)`

	var exists bool
	err := r.db.QueryRow(context.Background(), query, auth.Email, auth.Token).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresAuthRepository) Delete(id string) (bool, error) {
	query := `DELETE FROM auth WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		// Return false instead of nil because nil is not a boolean
		return false, err
	}

	// Return true if at least one row was deleted, otherwise false
	return true, nil
}
func (r *PostgresAuthRepository) DeleteByEmail(email string) error {
	query := `DELETE FROM auth WHERE email = $1`

	_, err := r.db.Exec(context.Background(), query, email)
	if err != nil {
		// Return false instead of nil because nil is not a boolean
		return err
	}
	// Return true if at least one row was deleted, otherwise false
	return nil
}
