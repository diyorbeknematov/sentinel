package models

import "time"

type Token struct {
	Token     string    `json:"token"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Login struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Register struct {
	Role     string `json:"role" validate:"required"`
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}