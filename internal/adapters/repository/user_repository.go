package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/strconvitoa/martian-service/internal/core/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Save(user domain.User) (domain.User, error) {
	svdusr := domain.User{}
	query := `
        INSERT INTO users (id, name, password, email, org_id, role, status,phone) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, name, email, org_id, role, status,phone
    `
	err := r.db.QueryRow(
		context.Background(),
		query,
		user.ID,
		user.Name,
		user.Password,
		user.Email,
		user.OrgID,
		user.Role,
		user.Status,
		user.Phone,
	).Scan(
		&svdusr.ID,
		&svdusr.Name,
		&svdusr.Email,
		&svdusr.OrgID,
		&svdusr.Role,
		&svdusr.Status,
		&svdusr.Phone,
	)

	if err != nil {
		return domain.User{}, err
	}

	return svdusr, nil
}
func (r *PostgresUserRepository) Exists(email string) (bool, error) {
	// The query returns true if a row exists, otherwise it returns no rows/null
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	// Using QueryRowContext is idiomatic for handling timeouts and cancellations
	err := r.db.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (r *PostgresUserRepository) CredentialsMatch(email string, password string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND password = $2)`

	var match bool
	err := r.db.QueryRow(context.Background(), query, email, password).Scan(&match)
	return match, err
}
func (r *PostgresUserRepository) Get(email string) (domain.User, error) {
	var user domain.User
	query := `
        SELECT id, name, email, org_id, role, status,phone 
        FROM users 
        WHERE email = $1
    `
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.OrgID,
		&user.Role,
		&user.Status,
		&user.Phone,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no user was found gracefully
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}
func (r *PostgresUserRepository) SelectPassword(email string) (string, error) {
	var pword string
	query := `
        SELECT password
        FROM users 
        WHERE email = $1
    `
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&pword,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return pword, errors.New("user not found")
		}
		return pword, err
	}

	return pword, nil
}
func (r *PostgresUserRepository) UpdatePassword(user domain.User) (domain.User, error) {
	query := `
        UPDATE users 
        SET password = $1,status = 'active'
        WHERE email = $2 AND id = $3
        RETURNING id, name, email, org_id, role, status, phone
    `

	// Create a fresh, empty user instance to hold the exact database return values
	var updatedUser domain.User

	// Scan into the new updatedUser struct variables
	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Password,
		user.Email,
		user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.OrgID,
		&updatedUser.Role,
		&updatedUser.Status,
		&updatedUser.Phone,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	// Returns the clean object directly matching the SQL RETURNING clause
	return updatedUser, nil
}
func (r *PostgresUserRepository) UpdateStatus(user domain.User) (domain.User, error) {
	// 1. Update the RETURNING clause to include all user columns
	query := `
        UPDATE users 
        SET status = $1
        WHERE email = $2 AND id = $3
        RETURNING id, name, password, email, org_id, role, status
    `

	// 2. Expand the Scan method to receive all returned columns
	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Status,
		user.Email,
		user.ID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.OrgID,
		&user.Role,
		&user.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}
func (r *PostgresUserRepository) DeleteByEmail(email string) error {
	query := `DELETE FROM users WHERE email = $1`

	// pgx uses Exec instead of ExecContext, passing ctx as the first argument
	result, err := r.db.Exec(context.Background(), query, email)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// pgx returns a CommandTag object, which has a RowsAffected() method
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with email: %s", email)
	}

	return nil
}

func (r *PostgresUserRepository) SelectAllByOrgID(orgID string) ([]domain.User, error) {
	query := `
        SELECT id, name, email,phone,role,status, org_id 
        FROM users 
        WHERE org_id = $1
    `
	rows, err := r.db.Query(context.Background(), query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.Status, &u.OrgID)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
