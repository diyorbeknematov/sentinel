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

func (r *agentRepo) CreateAgent(agent models.CreateAgentDB) (uuid.UUID, error) {
	query := `
		INSERT INTO agents (
		id,
		account_id,
		name,
		ip_address,
		last_seen
		) VALUES ($1, $2, $3, $4, $5);
	`

	if _, err := r.db.Exec(query,
		agent.Id,
		agent.AccountID,
		agent.Name,
		agent.IPAddress,
		agent.LastSeen,
	); err != nil {
		return uuid.Nil, err
	}

	return agent.Id, nil
}

func (r *agentRepo) GetAgentByID(id uuid.UUID) (models.Agent, error) {
	var agent models.Agent

	query := `
	SELECT 
		id,
		account_id,
		name,
		ip_address,
		last_seen,
		created_at
	FROM agents
	WHERE id = $1;
	`
	if err := r.db.QueryRow(query, id).Scan(
		&agent.Id,
		&agent.AccountID,
		&agent.Name,
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

func (r *agentRepo) ListAgents(filter models.FilterAgentDB) ([]models.Agent, int, error) {
	baseQuery := `
	SELECT
		id,
		account_id,
		name,
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
	if filter.AccountID != uuid.Nil {
		conditions = append(conditions, "account_id = :accountId")
		params["accountId"] = filter.AccountID
	}

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
	if filter.Limit > 0 {
		baseQuery += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"
	} else {
		baseQuery += " ORDER BY created_at DESC"
	}

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
			&a.AccountID,
			&a.Name,
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
