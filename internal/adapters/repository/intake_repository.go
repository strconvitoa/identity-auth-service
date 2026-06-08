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
        INSERT INTO intakes (id, name, email, phone, issue, org_id,summary,lang)
        VALUES ($1, $2, $3, $4, $5, $6,$7,$8)
        RETURNING id, name, email, phone, issue, org_id, summary, lang
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
		intake.Summary,
		intake.Lang,
	).Scan(
		&intake.ID,
		&intake.Name,
		&intake.Email,
		&intake.Phone,
		&intake.Issue,
		&intake.OrgID,
		&intake.Summary,
		&intake.Lang,
	)

	if err != nil {
		return domain.Intake{}, err
	}

	return intake, nil
}
