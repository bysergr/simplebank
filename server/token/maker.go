package token

import "time"

type tokenType string

const (
	RefreshToken tokenType = "refresh"
	AccessToken  tokenType = "access"
)

type Maker interface {
	CreateToken(username string, tokenType tokenType, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
