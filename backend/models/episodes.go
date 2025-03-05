package models

import (
	"time"
)

type Episode struct {
	ID      int64     `json:"id"`
	Season  int64     `json:"season"`
	Episode int64     `json:"episode"` // Episode number
	Title   string    `json:"title"`
	AirDate time.Time `json:"air_date"`
}
