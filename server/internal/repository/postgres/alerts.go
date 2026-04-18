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
	id := uuid.New()
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
		id,
		alert.AgentId,
		alert.Type,
		alert.Message,
		alert.Severity,
	); err != nil {
		return uuid.Nil, err
	}

	return id, nil
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

func (r *alertRepo) ListAlerts(filter models.FilterAlert) ([]models.Alert, int, error) {
	baseQuery := `
	SELECT 
		id,
		agent_id,
		type,
		message,
		severity,
		is_read,
		created_at
	FROM alerts
	WHERE TRUE 
	`

	countQuery := `SELECT COUNT(id) FROM alerts WHERE TRUE `

	conditions := []string{}
	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// Add search condition
	if filter.AgentId != nil {
		conditions = append(conditions, "agent_id = :agentId")
		params["agentId"] = filter.AgentId
	}

	if filter.Type != nil {
		conditions = append(conditions, "type ILIKE :type")
		params["type"] = "%" + *filter.Type + "%"
	}

	if filter.Severity != nil {
		conditions = append(conditions, "severity ILIKE :severity")
		params["severity"] = "%" + *filter.Severity + "%"
	}

	if filter.IsRead != nil {
		conditions = append(conditions, "is_read = :isRead")
		params["isRead"] = filter.IsRead
	}

	if filter.From != nil {
		conditions = append(conditions, "created_at >= :from")
		params["from"] = filter.From
	}

	if filter.To != nil {
		conditions = append(conditions, "created_at < :to")
		params["to"] = filter.To
	}

	// Add WHERE clause if contions exist
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += countQuery
	}

	// order + pagination
	baseQuery += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"

	// Execute the main query
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(
			&a.Id,
			&a.AgentId,
			&a.Type,
			&a.Message,
			&a.Severity,
			&a.CreatedAt,
		); err != nil {

			return nil, 0, err
		}
		alerts = append(alerts, a)
	}

	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, err
	}

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

func (r *alertRepo) MarkAlertRead(isRead models.MarkAlertRead) error {
	query := `
		UPDATE alerts 
		SET 
			is_read = $2
		WHERE id = $1
	`

	// Execute the query
	row, err := r.db.Exec(query,
		isRead.Id,
		isRead.IsRead,
	)
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
