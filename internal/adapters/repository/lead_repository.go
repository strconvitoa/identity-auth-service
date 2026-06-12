package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/strconvitoa/martian-service/internal/core/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresLeadRepository struct {
	db *pgxpool.Pool
}

func NewPostgresLeadRepository(db *pgxpool.Pool) *PostgresLeadRepository {
	return &PostgresLeadRepository{
		db: db,
	}
}

func (r *PostgresLeadRepository) Save(lead domain.Lead) (domain.Lead, error) {
	slead := domain.Lead{}
	query := `
        INSERT INTO leads (id, name, email, phone, issue, org_id,summary,lang,status,area,description,urgency,quality,conflict,retention)
        VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
        RETURNING id, name, email, phone, issue, org_id, summary, lang,status,area,description,urgency,quality,conflict,retention,created_at
    `
	err := r.db.QueryRow(
		context.Background(),
		query,
		lead.ID,
		lead.Name,
		lead.Email,
		lead.Phone,
		lead.Issue,
		lead.OrgID,
		lead.Summary,
		lead.Lang,
		lead.Status, lead.Area, lead.Description, lead.Urgency, lead.Quality, lead.Conflict, lead.Retention,
	).Scan(
		&slead.ID,
		&slead.Name,
		&slead.Email,
		&slead.Phone,
		&slead.Issue,
		&slead.OrgID,
		&slead.Summary,
		&slead.Lang, &slead.Status, &slead.Area, &slead.Description, &slead.Urgency, &slead.Quality, &slead.Conflict, &slead.Retention, &slead.Created,
	)

	if err != nil {
		return domain.Lead{}, err
	}

	return slead, nil
}
func (r *PostgresLeadRepository) UpdateLeadStatus(id string, status string) (domain.Lead, error) {
	var slead domain.Lead
	query := `
		UPDATE leads 
		SET status = $1
		WHERE id = $2
		RETURNING id, org_id, name,email,phone,issue,summary,lang,status,area,description, urgency,quality,conflict,retention, created_at;
	`
	err := r.db.QueryRow(context.Background(), query, status, id).Scan(
		&slead.ID,
		&slead.OrgID,
		&slead.Name,
		&slead.Email,
		&slead.Phone,
		&slead.Issue,
		&slead.Summary,
		&slead.Lang, &slead.Status, &slead.Area, &slead.Description, &slead.Urgency, &slead.Quality, &slead.Conflict, &slead.Retention, &slead.Created,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Lead{}, fmt.Errorf("lead with id %s not found", id)
		}
		return domain.Lead{}, fmt.Errorf("failed to update lead status: %w", err)
	}

	return slead, nil
}
func (r *PostgresLeadRepository) SelectByStatus(org_id string, status string) ([]domain.Lead, error) {
	leads := []domain.Lead{}
	query := `
        SELECT id, name,email,phone,issue,org_id,summary,lang, status, area, description,urgency,quality,conflict,retention,created_at 
        FROM leads 
        WHERE org_id = $1 AND status = $2
    `
	rows, err := r.db.Query(context.Background(), query, org_id, status)
	if err != nil {
		return nil, fmt.Errorf("failed to select leads by status: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var lead domain.Lead
		err := rows.Scan(
			&lead.ID,
			&lead.Name,
			&lead.Email,
			&lead.Phone,
			&lead.Issue,
			&lead.OrgID,
			&lead.Summary,
			&lead.Lang,
			&lead.Status,
			&lead.Area,
			&lead.Description,
			&lead.Urgency,
			&lead.Quality,
			&lead.Conflict,
			&lead.Retention,
			&lead.Created,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead row: %w", err)
		}
		leads = append(leads, lead)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return leads, nil
}
