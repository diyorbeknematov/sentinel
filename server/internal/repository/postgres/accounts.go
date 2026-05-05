package postgres

import (
	"database/sql"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type accountRepo struct {
	db *sqlx.DB
}

func NewAccountRepo(db *sqlx.DB) *accountRepo {
	return &accountRepo{
		db: db,
	}
}

// CREATE
func (r *accountRepo) CreateAccount(account models.CreateAccount) (uuid.UUID, error) {
	query := `
	INSERT INTO accounts (
		id,
		username,
		password,
		email,
		api_key
	) VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.db.Exec(query,
		account.Id,
		account.Username,
		account.Password,
		account.Email,
		account.APIKey,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return account.Id, nil
}

// GET BY API-KEY
func (r *accountRepo) GetByAPIKey(apiKey string) (models.Account, error) {
	var account models.Account

	query := `
	SELECT 
		id,
		username,
		password,
		email,
		api_key
	FROM accounts
	WHERE api_key = $1;
	`

	if err := r.db.QueryRow(query, apiKey).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.APIKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, errNoRowsAffected
		}
		return models.Account{}, err
	}

	return account, nil
}

// GET BY USERNAME
func (r *accountRepo) GetByUsername(username string) (models.Account, error) {
	var account models.Account

	query := `
	SELECT 
		id,
		username,
		password,
		email,
		api_key
	FROM accounts
	WHERE username = $1;
	`

	if err := r.db.QueryRow(query, username).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.APIKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, errNoRowsAffected
		}
		return models.Account{}, err
	}

	return account, nil
}

func (r *accountRepo) GetByEmail(email string) (models.Account, error) {
	var account models.Account

	query := `
	SELECT 
		id,
		username,
		password,
		email,
		api_key
	FROM accounts
	WHERE email = $1;
	`

	if err := r.db.QueryRow(query, email).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.APIKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, errNoRowsAffected
		}
		return models.Account{}, err
	}

	return account, nil
}

// GET BY ID
func (r *accountRepo) GetByID(id uuid.UUID) (models.Account, error) {
	var account models.Account

	query := `
	SELECT 
		id,
		username,
		password,
		email,
		api_key
	FROM accounts
	WHERE id = $1;
	`

	if err := r.db.QueryRow(query, id).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.APIKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, errNoRowsAffected
		}
		return models.Account{}, err
	}

	return account, nil
}

// UPDATE ACCOUNT
func (r *accountRepo) UpdateAccount(req models.UpdateAccountDB) error {
	query := `
	UPDATE accounts 
	SET
		username = $2,
		password = $3
	WHERE id = $1;
	`

	res, err := r.db.Exec(query,
		req.Id,
		req.Username,
		req.Password,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errNoRowsAffected
	}

	return nil
}

// DELETE
func (r *accountRepo) DeleteAccount(id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1;`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errNoRowsAffected
	}

	return nil
}
