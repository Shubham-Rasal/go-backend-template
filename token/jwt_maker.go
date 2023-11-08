package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const secretSize = 32

// implements maker interface
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secret string) (Maker, error) {

	if len(secret) < secretSize {
		return nil, fmt.Errorf("secret key smaller than %d", secretSize)
	}

	return &JWTMaker{
		secret,
	}, nil

}

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {

	//first create a payload with necessary fields
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	//hash the payload with an algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload.JWTPayload())

	//sign it with the secretkey and send it to the user
	return token.SignedString([]byte(maker.secretKey))

}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		//check if the algorithm is correct
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		//return the secret key
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, keyFunc)
	if err != nil {

		fmt.Println(err.Error())
		fmt.Println(jwt.ErrTokenInvalidClaims.Error())
		if strings.Contains(err.Error(), jwt.ErrTokenInvalidClaims.Error()) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        uuid.MustParse(payload.ID),
		Username:  payload.Username,
		IssuedAt:  time.Unix(int64(payload.IssuedAt.Time.Unix()), 0),
		ExpiresAt: time.Unix(int64(payload.ExpiresAt.Time.Unix()), 0),
	}, nil

}
