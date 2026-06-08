package repository

import (
	"context"

	"github.com/strconvitoa/martian-service/internal/core/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrgRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrgRepository(db *pgxpool.Pool) *PostgresOrgRepository {
	return &PostgresOrgRepository{
		db: db,
	}
}

func (r *PostgresOrgRepository) Save(Org domain.Org) (domain.Org, error) {
	query := `
        INSERT INTO orgs (id, name, trade_name)
        VALUES ($1, $2, $3)
    `

	_, err := r.db.Exec(context.Background(), query,
		Org.ID,
		Org.Name,
		Org.TradeName,
	)

	if err != nil {
		return domain.Org{}, err
	}

	return Org, nil
}
func (r *PostgresOrgRepository) Exists(org domain.Org) (bool, error) {
	// The query returns true if a row exists, otherwise it returns no rows/null
	query := `SELECT EXISTS(SELECT 1 FROM orgs WHERE name = $1 AND trade_name = $2)`

	var exists bool
	// Using QueryRowContext is idiomatic for handling timeouts and cancellations
	err := r.db.QueryRow(context.Background(), query, org.Name, org.TradeName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (r *PostgresOrgRepository) ExistsByID(id string) (bool, error) {
	// The query returns true if a row exists, otherwise it returns no rows/null
	query := `SELECT EXISTS(SELECT 1 FROM orgs WHERE id = $1)`

	var exists bool
	// Using QueryRowContext is idiomatic for handling timeouts and cancellations
	err := r.db.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
