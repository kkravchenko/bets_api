package entity

type BetRequest struct {
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	CrashPoint float64 `json:"crash_point"`
}
