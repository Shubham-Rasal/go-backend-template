package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// implements maker interface
type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (Maker, error) {

	if len(symetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}, nil

}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	if len(username) == 0 {
		return "", errors.New("username cannot be empty")
	}

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symetricKey, payload, nil)

}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil

}
