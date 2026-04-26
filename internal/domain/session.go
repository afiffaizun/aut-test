package domain

import "time"

type Session struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	RefreshTokenID string    `json:"refresh_token_id"`
	UserAgent      string    `json:"user_agent"`
	IPAddress      string    `json:"ip_address"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	IsRevoked      bool      `json:"is_revoked"`
}
