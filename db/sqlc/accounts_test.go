package db

import (
	"context"
	"testing"

	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {

	user := createDummyUser(t)

	arg := CreateAccountParams{
		Username: user.Username,
		Role:     util.RandomRole(),
		UserID:   int32(user.ID),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Username, account.Username)
	require.Equal(t, arg.Role, account.Role)
}

func TestGetAccount(t *testing.T) {
	user := createDummyUser(t)

	arg := CreateAccountParams{
		Username: user.Username,
		Role:     util.RandomRole(),
		UserID:   int32(user.ID),
	}

	account1, err := testQueries.CreateAccount(context.Background(), arg)
	account2, err := testQueries.GetAccount(context.Background(), account1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.Role, account2.Role)
}

func TestDeleteAccount(t *testing.T) {
	user := createDummyUser(t)

	arg := CreateAccountParams{
		Username: user.Username,
		Role:     util.RandomRole(),
		UserID:   int32(user.ID),
	}

	account1, err := testQueries.CreateAccount(context.Background(), arg)
	err = testQueries.DeleteAccount(context.Background(), account1.UserID)
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.UserID)
	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		user := createDummyUser(t)

		arg := CreateAccountParams{
			Username: user.Username,
			Role:     util.RandomRole(),
			UserID:   int32(user.ID),
		}
		_, err := testQueries.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateUser(t *testing.T) {
	user := createDummyUser(t)

	arg := CreateAccountParams{
		Username: user.Username,
		Role:     util.RandomRole(),
		UserID:   int32(user.ID),
	}

	account1, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	arg2 := UpdateReputationParams{
		UserID:     int32(account1.UserID),
		Reputation: int32(util.RandomInt(0, 100)),
	}

	err = testQueries.UpdateReputation(context.Background(), arg2)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), int32(account1.UserID))
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, arg2.Reputation, account2.Reputation)
}
