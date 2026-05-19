package postgres

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var metadataType []byte

type appLogRepo struct {
	db *sqlx.DB
}

func NewAppLogRepo(db *sqlx.DB) *appLogRepo {
	return &appLogRepo{
		db: db,
	}
}

func (r *appLogRepo) CreateAppLog(log models.CreateAppLog) (uuid.UUID, error) {
	query := `
	INSERT INTO applogs (
		id,
		agent_id,
		service_name,
		event,
		level,
		message,
		metadata,
		log_time
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	if _, err := r.db.Exec(query,
		log.Id,
		log.AgentId,
		log.ServiceName,
		log.Event,
		log.Level,
		log.Message,
		log.Metadata,
		log.LogTime,
	); err != nil {
		return uuid.Nil, errors.New("failed to create app log: " + err.Error())
	}

	return log.Id, nil
}

func (r *appLogRepo) GetLogByID(id uuid.UUID) (models.Log, error) {
	var log models.Log
	query := `
	SELECT
		id,
		agent_id,
		service_name,
		event,
		level,
		message,
		metadata,
		log_time,
		recorded_at
	FROM applogs
	WHERE id = $1;
	`

	if err := r.db.QueryRow(query, id).Scan(
		&log.Id,
		&log.AgentId,
		&log.ServiceName,
		&log.Event,
		&log.Level,
		&log.Message,
		&metadataType,
		&log.LogTime,
		&log.RecordedAt,
	); err != nil {
		return models.Log{}, errors.New("failed to get app log: " + err.Error())
	}

	if err := json.Unmarshal(metadataType, &log.Metadata); err != nil {
		return models.Log{}, errors.New("failed to unmarshal metadata: " + err.Error())
	}

	return log, nil
}

func (r *appLogRepo) ListLogs(filter models.FilterAppLogDB) ([]models.AppLogResponse, int, error) {

	baseQuery := `
	SELECT 
		ap.id,
		a.name AS agent_name,
		ap.service_name,
		ap.event,
		ap.level,
		ap.message,
		ap.metadata,
		ap.log_time,
		ap.recorded_at
	FROM applogs ap
	LEFT JOIN agents a ON ap.agent_id = a.id
	WHERE TRUE
	`

	countQuery := `
	SELECT COUNT(ap.id)
	FROM applogs ap
	LEFT JOIN agents a ON ap.agent_id = a.id
	WHERE TRUE
	`

	conditions := []string{}

	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// 🔹 Agent filter
	if filter.AgentId != uuid.Nil {
		conditions = append(conditions, "agent_id = :agentId")
		params["agentId"] = filter.AgentId
	}

	// 🔹 Level filter
	if filter.Level != "" {
		conditions = append(conditions, "level = :level")
		params["level"] = filter.Level
	}

	// 🔹 From filter
	if !filter.From.IsZero() {
		conditions = append(conditions, "log_time >= :from")
		params["from"] = filter.From
	}

	// 🔹 To filter
	if !filter.To.IsZero() {
		conditions = append(conditions, "log_time < :to")
		params["to"] = filter.To
	}

	// 🔹 Apply WHERE conditions
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// 🔹 Final query
	baseQuery += " ORDER BY recorded_at DESC LIMIT :limit OFFSET :offset"

	// MAIN QUERY
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, errors.New("failed to list app logs: " + err.Error())
	}
	defer rows.Close()

	var logs []models.AppLogResponse

	for rows.Next() {
		var l models.AppLogResponse

		if err := rows.Scan(
			&l.Id,
			&l.AgentName,
			&l.ServiceName,
			&l.Event,
			&l.Level,
			&l.Message,
			&metadataType,
			&l.LogTime,
			&l.RecordedAt,
		); err != nil {
			return nil, 0, errors.New("failed to scan app log: " + err.Error())
		}

		if err := json.Unmarshal(metadataType, &l.Metadata); err != nil {
			return nil, 0, errors.New("failed to unmarshal metadata: " + err.Error())
		}

		logs = append(logs, l)
	}

	// COUNT QUERY
	countSQL, args, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, errors.New("failed to prepare count query: " + err.Error())
	}

	countSQL = r.db.Rebind(countSQL)

	var total int
	if err := r.db.Get(&total, countSQL, args...); err != nil {
		return nil, 0, errors.New("failed to get log count: " + err.Error())
	}

	return logs, total, nil
}
