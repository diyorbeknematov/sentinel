package postgres

import (
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(user models.CreateUser) (uuid.UUID, error) {
	id := uuid.New()

	query := `
	INSERT INTO users (
		id,
		username,
		password,
		role
	) VALUES ($1, $2, $3, $4);
	`

	if _, err := r.db.Exec(query,
		id,
		user.UserName,
		user.Password,
		user.Role,
	); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userRepo) GetUserByUserName(userName string) (models.User, error) {
	var user models.User
	query := `
	SELECT 
		id,
		username,
		password,
		role
	FROM users
	WHERE username = $1;
	`

	if err := r.db.QueryRow(query, userName).Scan(
		&user.Id,
		&user.UserName,
		&user.Password,
		&user.Role,
	); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepo) GetUserByID(id uuid.UUID) (models.User, error) {
	var user models.User
	query := `
	SELECT 
		id,
		username,
		password,
		role
	FROM users
	WHERE id = $1;
	`

	if err := r.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.UserName,
		&user.Password,
		&user.Role,
	); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepo) UpdateUser(req models.UpdateUser) error {
	query := `
	UPDATE users 
	SET
		usernam = $2,
		password = $3
	WHERE id = $1;
	`
	row, err := r.db.Exec(query,
		req.Id,
		req.UserName,
		req.Password,
	)
	if err != nil {
		return nil
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

func (r *userRepo) UpdateUserRole(req models.UpdateRole) error {
	query := `
	UPDATE users 
	SET 
		role = $2
	WHERE id = $1
	`

	row, err := r.db.Exec(query,
		req.Id,
		req.Role,
	)
	if err != nil {
		return nil
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

func (r *userRepo) DeleteUser(id uuid.UUID) error {
	query := ` DELETE FROM users WHERE id = $1; `

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
