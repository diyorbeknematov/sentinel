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
	query := `
	INSERT INTO metrics (
		id,
		agent_id,
		cpu,
		ram,
		disk,
		log_time
	) VALUES ($1, $2, $3, $4, $5, $6);
	`
	if _, err := r.db.Exec(query,
		metric.Id,
		metric.AgentId,
		metric.CPU,
		metric.RAM,
		metric.Disk,
		metric.LogTime,
	); err != nil {
		return uuid.Nil, err
	}

	return metric.Id, nil
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
		log_time,
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
		&metric.LogTime,
		&metric.RecordedAt,
	); err != nil {
		return models.Metric{}, err
	}

	return metric, nil
}

func (r *metricRepo) ListMetrics(filter models.FilterMetricsDB) ([]models.MetricResponse, int, error) {
	baseQuery := `
	SELECT 
		m.id,
		a.name AS agent_name,
		m.cpu,
		m.ram,
		m.disk,
		m.log_time,
		m.recorded_at
	FROM metrics m
	LEFT JOIN agents a ON m.agent_id = a.id
	WHERE TRUE 
	`
	countQuery := `
	SELECT 
		COUNT(m.id) 
	FROM metrics m
	LEFT JOIN agents a ON m.agent_id = a.id
	WHERE TRUE `

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

	var metrics []models.MetricResponse
	for rows.Next() {
		var m models.MetricResponse
		if err := rows.Scan(
			&m.Id,
			&m.AgentName,
			&m.CPU,
			&m.RAM,
			&m.Disk,
			&m.LogTime,
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
