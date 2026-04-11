package postgres

import (
	"log/slog"
	"strings"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type logRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewLogRepo(db *sqlx.DB, logger *slog.Logger) *logRepo {
	return &logRepo{
		db:     db,
		logger: logger,
	}
}

func (r *logRepo) CreateLog(log models.CreateLog) (uuid.UUID, error) {
	id := uuid.New()

	query := `
	INERT INTO logs (
		id,
		agent_id,
		type,
		level,
		message,
		ip_address
	) VALUES ($1, $2, $3, $4, $5, $6)
	`
	if _, err := r.db.Exec(query,
		id,
		log.AgentId,
		log.Type,
		log.Level,
		log.Message,
		log.IPAddress,
	); err != nil {
		r.logger.Error(err.Error())
		return uuid.Nil, err
	}

	return id, nil
}

func (r *logRepo) CreateLogsBatch(logs []models.CreateLog) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	query := `
	INSERT INTO logs (
		id,
		agent_id,
		type,
		level,
		message,
		ip_address
	)VALUES($1, $2, $3, $4, $5, %6)
	`
	for _, log := range logs {
		id := uuid.New()

		if _, err := tx.Exec(query,
			id,
			log.AgentId,
			log.Type,
			log.Level,
			log.Message,
			log.IPAddress,
		); err != nil {
			r.logger.Error(err.Error())
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}

func (r *logRepo) ListLogs(filter models.FilterLog) ([]models.Log, int, error) {
	baseQuery := `
	SELECT 
		id,
		agent_id,
		type,
		level,
		message,
		ip_address,
		recorded_at
	FROM logs
	WHER TRUE 
	`

	countQuery := `SELECT COUNT(id) FROM logs WHERE TRUE `

	conditions := []string{}
	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offest,
	}

	// Add search condition
	if filter.AgentId != uuid.Nil {
		conditions = append(conditions, "agent_id = :agentId")
		params["agentId"] = filter.AgentId
	}

	if filter.Type != "" {
		conditions = append(conditions, "type ILIKE :type")
		params["type"] = "%" + filter.Type + "%"
	}

	if filter.Level != "" {
		conditions = append(conditions, "level ILIKE :level")
		params["type"] = "%" + filter.Level + "%"
	}

	if filter.IPAddress != "" {
		conditions = append(conditions, "ip_address = :ipAddrss")
		params["ipAddress"] = filter.IPAddress
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
		r.logger.Error(err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var l models.Log
		if err := rows.Scan(
			&l.Id,
			&l.AgentId,
			&l.Type,
			&l.Level,
			&l.Message,
			&l.IPAddress,
			&l.RecordedAt,
		); err != nil {
			r.logger.Error(err.Error())
			return nil, 0, err
		}

		logs = append(logs, l)
	}

	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, 0, err
	}

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		r.logger.Error(err.Error())
		return nil, 0, err
	}

	return logs, total, nil
}
