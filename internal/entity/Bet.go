package entity

import "time"

type Bet struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Amount     float64   `json:"amount"`
	CrashPoint float64   `json:"crash_point"`
	CreatedAt  time.Time `json:"created_at"`
}
