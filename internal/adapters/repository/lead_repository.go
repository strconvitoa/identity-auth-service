package repository

import (
	"context"
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
	savedlead := domain.Lead{}
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
		&savedlead.ID,
		&savedlead.Name,
		&savedlead.Email,
		&savedlead.Phone,
		&savedlead.Issue,
		&savedlead.OrgID,
		&savedlead.Summary,
		&savedlead.Lang, &savedlead.Status, &savedlead.Area, &savedlead.Description, &savedlead.Urgency, &savedlead.Quality, &savedlead.Conflict, &savedlead.Retention, &savedlead.Created,
	)

	if err != nil {
		return domain.Lead{}, err
	}

	return savedlead, nil
}
func (r *PostgresLeadRepository) SelectByStatus(org_id string, status string) ([]domain.Lead, error) {
	// 1. Initialize a slice to hold the results.
	// Using an empty slice instead of nil ensures you return `[]` instead of `null` in JSON.
	leads := []domain.Lead{}

	// 2. Define the query
	query := `
        SELECT id, name,email,phone,issue,org_id,summary,lang, status, area, description,urgency,quality,conflict,retention,created_at 
        FROM leads 
        WHERE org_id = $1 AND status = $2
    `

	// 3. Execute the query
	rows, err := r.db.Query(context.Background(), query, org_id, status)
	if err != nil {
		return nil, fmt.Errorf("failed to select leads by status: %w", err)
	}
	defer rows.Close()

	// 4. Iterate through the result set
	for rows.Next() {
		var lead domain.Lead
		// Destination arguments must match the columns in your SELECT statement exactly
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
