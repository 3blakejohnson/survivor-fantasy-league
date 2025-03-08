package models

import (
	"time"
)

type League struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Name       string    `json:"name"`
	OwnerID    int64     `json:"owner_id"`
	InviteCode string    `json:"invite_code"`
	Season     int64     `json:"season"`
}
