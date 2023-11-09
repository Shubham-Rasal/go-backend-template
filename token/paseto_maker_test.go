package token

import (
	"log"
	"testing"
	"time"

	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	key := util.RandomString(32)
	maker, err := NewPasetoMaker(key)
	log.Println(key)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	maker1, err := NewPasetoMaker("dsfdfsd")
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
	require.WithinDuration(t, payload.ExpiresAt, time.Now().Add(duration), time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
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

func TestInvalidPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
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
