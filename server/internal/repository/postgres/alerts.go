package postgres

import (
	"strings"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type alertRepo struct {
	db *sqlx.DB
}

func NewAlertRepo(db *sqlx.DB) *alertRepo {
	return &alertRepo{
		db: db,
	}
}

func (r *alertRepo) CreateAlert(alert models.CreateAlert) (uuid.UUID, error) {
	query := `
	INSERT INTO alerts (
		id,
		agent_id,
		type,
		message,
		severity
	) VALUES ($1, $2, $3, $4, $5);
	`
	if _, err := r.db.Exec(query,
		alert.Id,
		alert.AgentId,
		alert.Type,
		alert.Message,
		alert.Severity,
	); err != nil {
		return uuid.Nil, err
	}

	return alert.Id, nil
}

func (r *alertRepo) GetAlertByID(id uuid.UUID) (models.Alert, error) {
	var alert models.Alert
	query := `
	SELECT
		id,
		agent_id,
		type,
		message,
		severity,
		is_read,
		created_at
	FROM alerts
	WHERE id = $1
	`

	if err := r.db.QueryRow(query, id).Scan(
		&alert.Id,
		&alert.AgentId,
		&alert.Type,
		&alert.Message,
		&alert.Severity,
		&alert.IsRead,
		&alert.CreatedAt,
	); err != nil {
		return models.Alert{}, err
	}

	return alert, nil
}

func (r *alertRepo) ListAlerts(filter models.FilterAlertDB) ([]models.AlertResponse, int, error) {

	baseQuery := `
	SELECT 
		al.id,
		a.name AS agent_name,
		al.type,
		al.message,
		al.severity,
		al.is_read,
		al.created_at
	FROM alerts al
	LEFT JOIN agents a ON al.agent_id = a.id
	WHERE TRUE
	`

	countQuery := `
	SELECT COUNT(al.id)
	FROM alerts al
	LEFT JOIN agents a ON al.agent_id = a.id
	WHERE TRUE
	`

	conditions := []string{}

	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// FILTERS
	// Agent filter
	if filter.AgentId != nil {
		conditions = append(conditions, "agent_id = :agentId")
		params["agentId"] = *filter.AgentId
	}

	// Severity filter
	if filter.Severity != nil {
		conditions = append(conditions, "severity = :severity")
		params["severity"] = strings.ToLower(*filter.Severity)
	}

	// IsRead filter
	if filter.IsRead != nil {
		conditions = append(conditions, "is_read = :isRead")
		params["isRead"] = *filter.IsRead
	}

	// From filter
	if filter.From != nil {
		conditions = append(conditions, "al.created_at >= :from")
		params["from"] = *filter.From
	}

	// To filter
	if filter.To != nil {
		conditions = append(conditions, "al.created_at < :to")
		params["to"] = *filter.To
	}

	// Build WHERE clause
	if len(conditions) > 0 {
		where := " AND " + strings.Join(conditions, " AND ")
		baseQuery += where
		countQuery += where
	}

	// ORDER, LIMIT, OFFSET
	if filter.Limit == 0 {
		baseQuery += " ORDER BY al.is_read ASC, al.created_at DESC"
	} else {
		baseQuery += " ORDER BY al.is_read ASC, al.created_at DESC LIMIT :limit OFFSET :offset"
	}

	// Execute the main query
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var alerts []models.AlertResponse

	for rows.Next() {
		var a models.AlertResponse

		if err := rows.Scan(
			&a.Id,
			&a.AgentName,
			&a.Type,
			&a.Message,
			&a.Severity,
			&a.IsRead,
			&a.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		alerts = append(alerts, a)
	}

	// COUNT QUERY
	countSQL, args, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, err
	}

	countSQL = r.db.Rebind(countSQL)

	var total int
	if err := r.db.Get(&total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

func (r *alertRepo) MarkAlertRead(id uuid.UUID) error {
	query := `
		UPDATE alerts 
		SET 
			is_read = TRUE
		WHERE id = $1
	`

	// Execute the query
	row, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return errNoRowsAffected
	}

	return nil
}
