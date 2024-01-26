package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID        `json:"id,omitempty"`
	Username  string           `json:"username,omitempty"`
	Issuer    string           `json:"issuer,omitempty"`
	Subject   string           `json:"subject,omitempty"`
	ExpiresAt time.Time        `json:"expires_at"`
	NotBefore time.Time        `json:"not_before"`
	IssuedAt  time.Time        `json:"issued_at"`
	Audience  jwt.ClaimStrings `json:"audience,omitempty"`
	TokenType tokenType        `json:"token_type,omitempty"`
}

func NewPayload(username string, tokenType tokenType, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		ExpiresAt: time.Now().Add(duration),
		NotBefore: time.Now(),
		IssuedAt:  time.Now(),
		TokenType: tokenType,
	}

	return payload, err
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}

	if time.Now().Before(payload.NotBefore) {
		return ErrExpiredToken
	}

	return nil
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.ExpiresAt), nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.NotBefore), nil
}

func (payload *Payload) GetIssuer() (string, error) {
	return payload.Issuer, nil
}

func (payload *Payload) GetSubject() (string, error) {
	return payload.Subject, nil
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return payload.Audience, nil
}
