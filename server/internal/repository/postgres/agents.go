package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	errNoRowsAffected = errors.New("no rows affected")
)

type agentRepo struct {
	db *sqlx.DB
}

func NewAgentRepo(db *sqlx.DB) *agentRepo {
	return &agentRepo{
		db: db,
	}
}

func (r *agentRepo) CreateAgent(agent models.CreateAgent) (uuid.UUID, error) {
	id := uuid.New()

	query := `
		INSERT INTO agents (
		id,
		name,
		api_key,
		ip_address,
		last_seen
		) VALUES ($1, $2, $3, $4, $5);
	`

	if _, err := r.db.Exec(query,
		id,
		agent.Name,
		agent.APIKey,
		agent.IPAddress,
		agent.LastSeen,
	); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *agentRepo) GetAgentByID(id uuid.UUID) (models.Agent, error) {
	var agent models.Agent

	query := `
	SELECT 
		id,
		name,
		api_key,
		ip_address,
		last_seen,
		created_at
	FROM agents
	WHERE id = $1;
	`
	if err := r.db.QueryRow(query, id).Scan(
		&agent.Id,
		&agent.Name,
		&agent.APIKey,
		&agent.IPAddress,
		&agent.LastSeen,
		&agent.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Agent{}, err
		}
	}

	return agent, nil
}

func (r *agentRepo) GetAgentByAPIKey(apiKey string) (models.Agent, error) {
	var agent models.Agent

	query := `
	SELECT
		id,
		name,
		api_key,
		ip_address,
		last_seen,
		created_at
	FROM agents
	WHERE api_key = $1;
	`

	if err := r.db.QueryRow(query, apiKey).Scan(
		&agent.Id,
		&agent.Name,
		&agent.APIKey,
		&agent.IPAddress,
		&agent.LastSeen,
		&agent.CreatedAt,
	); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return models.Agent{}, err
		}

		return models.Agent{}, err
	}

	return agent, nil
}

func (r *agentRepo) ListAgents(filter models.FilterAgent) ([]models.Agent, int, error) {
	baseQuery := `
	SELECT
		id,
		name,
		api_key,
		ip_address,
		last_seen,
		created_at
	FROM agents WHERE TRUE 
	`
	countQuery := `SELECT COUNT(id) FROM agents WHERE TRUE `

	conditions := []string{}

	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// Add search condition
	if filter.Name != "" {
		conditions = append(conditions, "name ILIKE :name")
		params["name"] = "%" + filter.Name + "%"
	}

	if !filter.From.IsZero() {
		conditions = append(conditions, "last_seen >= :from")
		params["from"] = filter.From
	}

	if !filter.To.IsZero() {
		conditions = append(conditions, "last_seen < :to")
		params["to"] = filter.To
	}

	if len(conditions) > 0 {
		condStr := " AND " + strings.Join(conditions, " AND ")
		baseQuery += condStr
		countQuery += condStr
	}

	// Pagination
	baseQuery += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"

	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var a models.Agent
		if err := rows.Scan(
			&a.Id,
			&a.Name,
			&a.APIKey,
			&a.IPAddress,
			&a.LastSeen,
			&a.CreatedAt,
		); err != nil {

			return nil, 0, err
		}
		agents = append(agents, a)
	}

	// Execute the count query
	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, err
	}
	countQuery = r.db.Rebind(countQuery)

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

func (r *agentRepo) UpdateLastSeen(lastSeen models.UpdateLastSeen) error {
	query := `
	UPDATE agents 
	SET 
		last_seen = $2
	WHERE id = $1
	`

	// Execute the query
	row, err := r.db.Exec(query,
		lastSeen.Id,
		lastSeen.LastSeen,
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

func (r *agentRepo) DeleteAgent(id uuid.UUID) error {
	query := `DELETE FROM agents WHERE id = $1;`

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
