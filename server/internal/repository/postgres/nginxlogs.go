package postgres

import (
	"strings"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type nginxLogRepo struct {
	db *sqlx.DB
}

func NewNGINXLogRepo(db *sqlx.DB) *nginxLogRepo {
	return &nginxLogRepo{
		db: db,
	}
}

func (r *nginxLogRepo) CreateNginxLog(log models.CreateNginxLog) (uuid.UUID, error) {
	id := uuid.New()

	query := `
	INSERT INTO nginxlogs (
		id,
		agent_id,
		ip_address,
		method,
		path,
		status,
		response_time
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := r.db.Exec(query,
		id,
		log.AgentId,
		log.IPAddress,
		log.Method,
		log.Path,
		log.Status,
		log.ResponseTime,
	); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *nginxLogRepo) GetNginxLogByID(id uuid.UUID) (models.NginxLog, error) {
	var log models.NginxLog

	query := `
	SELECT
		id,
		agent_id,
		ip_address,
		method,
		path,
		status,
		response_time,
		recorded_at
	FROM nginxlogs
	WHERE id = $1
	`

	if err := r.db.QueryRow(query, id).Scan(
		&log.Id,
		&log.AgentId,
		&log.IPAddress,
		&log.Method,
		&log.Path,
		&log.Status,
		&log.ResponseTime,
		&log.RecordedAt,
	); err != nil {
		return models.NginxLog{}, err
	}

	return log, nil
}

func (r *nginxLogRepo) ListNginxLogs(filter models.FilterNginxLog) ([]models.NginxLog, int, error) {
	baseQuery := `
	SELECT
		id,
		agent_id,
		ip_address,
		method,
		path,
		status,
		response_time,
		recorded_at
	FROM nginxlogs
	WHERE TRUE 
	`

	countQuery := `SELECT COUNT(id) FROM nginxlogs WHER TRUE `

	conditions := []string{}
	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// Add search condition
	if filter.AgentId != uuid.Nil {
		conditions = append(conditions, "agent_id = :agentId")
		params["agentId"] = filter.AgentId
	}

	if filter.Method != "" {
		conditions = append(conditions, "method = :method")
		params["method"] = filter.Method
	}

	if filter.Status != 0 {
		conditions = append(conditions, "status = :status")
		params["status"] = filter.Status
	}

	if !filter.From.IsZero() {
		conditions = append(conditions, "recorded_at >= :from")
		params["from"] = filter.From
	}

	if !filter.To.IsZero() {
		conditions = append(conditions, "recorded_at < :to")
		params["to"] = filter.To
	}

	// ADD WHERE clause if conditions exist
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	baseQuery += " ORDER BY recorded_at DESC LIMIT :limit OFFSET :offset"

	// Exectue the main query
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.NginxLog
	for rows.Next() {
		var l models.NginxLog
		if err := rows.Scan(
			&l.Id,
			&l.AgentId,
			&l.IPAddress,
			&l.Method,
			&l.Path,
			&l.Status,
			&l.ResponseTime,
			&l.RecordedAt,
		); err != nil {
			return nil, 0, err
		}

		logs = append(logs, l)
	}

	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, err
	}

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
