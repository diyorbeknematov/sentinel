package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	UserName string    `json:"username"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

type CreateUser struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type UpdateUser struct {
	Id       uuid.UUID `json:"-"`
	UserName string    `json:"username" validate:"required"`
	Password string    `json:"password" validate:"required"`
}

type UpdateRole struct {
	Id   uuid.UUID `json:"-"`
	Role uuid.UUID `json:"role" validate:"required"`
}
