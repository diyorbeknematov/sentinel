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
	id := uuid.New()

	query := `
	INERT INTO applogs (
		id,
		agent_id,
		user_id,
		type,
		level,
		message,
		ip_address
	) VALUES ($1, $2, $3, $4, $5, $6)
	`
	if _, err := r.db.Exec(query,
		id,
		log.AgentId,
		log.UserId,
		log.Type,
		log.Level,
		log.Message,
		log.IPAddress,
	); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *appLogRepo) GetLogByID(id uuid.UUID) (models.Log, error) {
	var log models.Log
	query := `
	SELECT
		id,
		agent_id,
		type,
		level,
		message,
		ip_address,
		recorded_at
	FROM applogs
	WHER id = $1;
	`

	if err := r.db.QueryRow(query, id).Scan(
		&log.Id,
		&log.AgentId,
		&log.Type,
		&log.Level,
		&log.Message,
		&log.IPAddress,
		&log.RecordedAt,
	); err != nil {
		return models.Log{}, err
	}

	return log, nil
}

func (r *appLogRepo) ListLogs(filter models.FilterAppLog) ([]models.Log, int, error) {
	baseQuery := `
	SELECT 
		id,
		agent_id,
		user_id,
		type,
		level,
		message,
		ip_address,
		recorded_at
	FROM applogs
	WHERE TRUE 
	`

	countQuery := `SELECT COUNT(id) FROM applogs WHERE TRUE `

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

	if filter.UserId != "" {
		conditions = append(conditions, "user_id = :userId")
		params["userId"] = filter.UserId
	}

	if filter.Type != "" {
		conditions = append(conditions, "type ILIKE :type")
		params["type"] = "%" + filter.Type + "%"
	}

	if filter.Level != "" {
		conditions = append(conditions, "level ILIKE :level")
		params["type"] = "%" + filter.Level + "%"
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

	var logs []models.Log
	for rows.Next() {
		var l models.Log
		if err := rows.Scan(
			&l.Id,
			&l.AgentId,
			&l.UserId,
			&l.Type,
			&l.Level,
			&l.Message,
			&l.IPAddress,
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
