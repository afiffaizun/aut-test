package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token *RefreshToken) (string, error)
	GetRefreshToken(ctx context.Context, tokenID string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
