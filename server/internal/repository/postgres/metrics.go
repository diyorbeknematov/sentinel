package postgres

import (
	"strings"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type metricRepo struct {
	db *sqlx.DB
}

func NewMetricRepo(db *sqlx.DB) *metricRepo {
	return &metricRepo{
		db: db,
	}
}

func (r *metricRepo) CreateMetric(metric models.CreateMetric) (uuid.UUID, error) {
	id := uuid.New()

	query := `
	INERT INTO metrics (
		id,
		agent_id,
		cpu,
		ram,
		disk
	) VALUES ($1, $2, $3, $4, $5);
	`
	if _, err := r.db.Exec(query,
		id,
		metric.AgentId,
		metric.CPU,
		metric.RAM,
		metric.Disk,
	); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *metricRepo) GetMetricsByID(id uuid.UUID) (models.Metric, error) {
	var metric models.Metric
	query := `
	SELECT 
		id,
		agent_id,
		cpu,
		ram,
		disk,
		recorded_at
	FROM metrics
	WHERE id = $1
	`

	if err := r.db.QueryRow(query, id).Scan(
		&metric.Id,
		&metric.AgentId,
		&metric.CPU,
		&metric.RAM,
		&metric.Disk,
		&metric.RecordedAt,
	); err != nil {
		return models.Metric{}, err
	}

	return metric, nil
}

func (r *metricRepo) ListMetrics(filter models.FilterMetrics) ([]models.Metric, int, error) {
	baseQuery := `
	SELECT 
		id,
		agent_id,
		cpu,
		ram,
		disk,
		recorded_at
	FROM metrics
	WHERE TRUE 
	`
	countQuery := `SELECT COUNT(id) FROM metrics WHERE TRUE `

	conditions := []string{}
	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// Add search question
	if filter.AgentId != uuid.Nil {
		conditions = append(conditions, "agent_id = :agentId")
		params["agentId"] = filter.AgentId
	}

	if !filter.From.IsZero() {
		conditions = append(conditions, "recorded_at >= :from")
		params["from"] = filter.From
	}

	if !filter.To.IsZero() {
		conditions = append(conditions, "recorded_at < :to")
		params["to"] = filter.To
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	baseQuery += " ORDER BY recorded_at DESC LIMIT :limit OFFSET :offset"

	// Execute the main query
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var m models.Metric
		if err := rows.Scan(
			&m.Id,
			&m.AgentId,
			&m.CPU,
			&m.RAM,
			&m.Disk,
			&m.RecordedAt,
		); err != nil {

			return nil, 0, err
		}

		metrics = append(metrics, m)
	}

	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, err
	}
	countQuery = r.db.Rebind(countQuery)

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, err
	}

	return metrics, total, nil
}
