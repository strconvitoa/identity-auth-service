package repository

import (
	"context"

	"github.com/strconvitoa/martian-service/internal/core/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIntakeRepository struct {
	db *pgxpool.Pool
}

func NewPostgresIntakeRepository(db *pgxpool.Pool) *PostgresIntakeRepository {
	return &PostgresIntakeRepository{
		db: db,
	}
}

func (r *PostgresIntakeRepository) Save(intake domain.Intake) (domain.Intake, error) {

	query := `
        INSERT INTO intakes (id, name, email, phone, issue, org_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, name, email, phone, issue, org_id
    `
	err := r.db.QueryRow(
		context.Background(),
		query,
		intake.ID,
		intake.Name,
		intake.Email,
		intake.Phone,
		intake.Issue,
		intake.OrgID,
	).Scan(
		&intake.ID,
		&intake.Name,
		&intake.Email,
		&intake.Phone,
		&intake.Issue,
		&intake.OrgID,
	)

	if err != nil {
		return domain.Intake{}, err
	}

	return intake, nil
}
