package models

import "time"

type User struct {
	ID           int64      `json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
}
