package postgres

import (
	"strings"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

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
		user_id,
		event,
		level,
		message,
		log_time
	) VALUES ($1, $2, $3, $4, $5, $6, &7);
	`
	if _, err := r.db.Exec(query,
		log.Id,
		log.AgentId,
		log.UserId,
		log.Event,
		log.Level,
		log.Message,
		log.LogTime,
	); err != nil {
		return uuid.Nil, err
	}

	return log.Id, nil
}

func (r *appLogRepo) GetLogByID(id uuid.UUID) (models.Log, error) {
	var log models.Log
	query := `
	SELECT
		id,
		agent_id,
		event,
		level,
		message,
		log_time,
		recorded_at
	FROM applogs
	WHERE id = $1;
	`

	if err := r.db.QueryRow(query, id).Scan(
		&log.Id,
		&log.AgentId,
		&log.Event,
		&log.Level,
		&log.Message,
		&log.LogTime,
		&log.RecordedAt,
	); err != nil {
		return models.Log{}, err
	}

	return log, nil
}

func (r *appLogRepo) ListLogs(filter models.FilterAppLogDB) ([]models.AppLogResponse, int, error) {

	baseQuery := `
	SELECT 
		ap.id,
		a.name AS agent_name,
		ap.user_id,
		ap.event,
		ap.level,
		ap.message,
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
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.AppLogResponse

	for rows.Next() {
		var l models.AppLogResponse

		if err := rows.Scan(
			&l.Id,
			&l.AgentName,
			&l.UserId,
			&l.Event,
			&l.Level,
			&l.Message,
			&l.LogTime,
			&l.RecordedAt,
		); err != nil {
			return nil, 0, err
		}

		logs = append(logs, l)
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

	return logs, total, nil
}
