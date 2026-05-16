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
