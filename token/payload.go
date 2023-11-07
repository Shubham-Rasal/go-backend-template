package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil

}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (payload *Payload) JWTPayload() JWTClaims {
	// Create claims while leaving out some of the optional fields
	return JWTClaims{
		payload.Username,
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			Issuer:    "blog",
			ID:        payload.ID.String(),
		},
	}
}
