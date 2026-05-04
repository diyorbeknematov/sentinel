package models

import "github.com/google/uuid"

type Account struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	Email    string    `json:"email"`
	APIKey   string    `json:"api_key,omitempty"`
}

type CreateAccount struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password" db:"password"`
	Email    string    `json:"email" db:"email"`
	APIKey   string    `json:"api_key,omitempty" db:"api_key"`
}

type UpdateAccount struct {
	Id       string `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateAccountDB struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}
