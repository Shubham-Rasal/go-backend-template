package token

import (
	"testing"
	"time"

	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker("12dghkeghkehgkdjfhgsdljfsldjfslkdfjlsdfjlsdjfldskjflsdjflskdjflskdjflsdjf")
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	maker1, err := NewJWTMaker("df")
	require.Error(t, err)
	require.Empty(t, maker1)

	username := util.RandomString(6)
	duration := time.Minute

	// Create token
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify token
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.NotZero(t, payload.ID)
	require.NotZero(t, payload.IssuedAt)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(duration), time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomString(6)
	duration := time.Minute

	// Create token
	token, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify token
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrExpiredToken.Error())
	require.Empty(t, payload)

}

func TestInvalidJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomString(6)
	duration := time.Minute

	// Create token
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify token
	payload, err := maker.VerifyToken(token + "12")
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrInvalidToken.Error())
	require.Empty(t, payload)

}

func TestInvalidJWTTokenAlgorithm(t *testing.T) {

	payload, err := NewPayload(util.RandomString(6), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload.JWTPayload())
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	// Verify token
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())

}
