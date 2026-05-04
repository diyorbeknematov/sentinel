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
	query := `
	INSERT INTO nginxlogs (
		id,
		agent_id,
		ip_address,
		method,
		path,
		status,
		bytes,
		user_agent,
		log_time
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := r.db.Exec(query,
		log.Id,
		log.AgentId,
		log.IPAddress,
		log.Method,
		log.Path,
		log.Status,
		log.Bytes,
		log.UserAgent,
		log.LogTime,
	); err != nil {
		return uuid.Nil, err
	}

	return log.Id, nil
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
		bytes,
		user_agent,
		log_time,
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
		&log.Bytes,
		&log.UserAgent,
		&log.LogTime,
		&log.RecordedAt,
	); err != nil {
		return models.NginxLog{}, err
	}

	return log, nil
}

func (r *nginxLogRepo) ListNginxLogs(filter models.FilterNginxLogDB) ([]models.NginxLogResponse, int, error) {

	baseQuery := `
	SELECT
		nl.id,
		a.name AS agent_name,
		nl.ip_address,
		nl.method,
		nl.path,
		nl.status,
		nl.bytes,
		nl.user_agent,
		nl.log_time,
		nl.recorded_at
	FROM nginxlogs nl
	LEFT JOIN agents a ON nl.agent_id = a.id
	WHERE TRUE
	`

	countQuery := `
	SELECT COUNT(nl.id)
	FROM nginxlogs nl
	LEFT JOIN agents a ON nl.agent_id = a.id
	WHERE TRUE
	`

	conditions := []string{}
	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	if filter.AgentId != uuid.Nil {
		conditions = append(conditions, "agent_id = :agent_id")
		params["agent_id"] = filter.AgentId
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

	if len(conditions) > 0 {
		where := " AND " + strings.Join(conditions, " AND ")
		baseQuery += where
		countQuery += where
	}

	baseQuery += " ORDER BY recorded_at DESC LIMIT :limit OFFSET :offset"

	// 🔹 MAIN QUERY
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.NginxLogResponse
	for rows.Next() {
		var l models.NginxLogResponse

		err := rows.Scan(
			&l.Id,
			&l.AgentName,
			&l.IPAddress,
			&l.Method,
			&l.Path,
			&l.Status,
			&l.Bytes,
			&l.UserAgent,
			&l.LogTime,
			&l.RecordedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		logs = append(logs, l)
	}

	// 🔹 COUNT QUERY
	countQuery, args, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, err
	}

	countQuery = sqlx.Rebind(sqlx.DOLLAR, countQuery)

	var total int
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
