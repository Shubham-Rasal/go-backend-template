package db

import (
	"context"
	"testing"

	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username: util.RandomUserName(),
		Role:     util.RandomRole(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Role, user.Role)
}

func TestGetUser(t *testing.T) {
	arg := CreateUserParams{
		Username: util.RandomUserName(),
		Role:     util.RandomRole(),
	}

	user1, err := testQueries.CreateUser(context.Background(), arg)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Role, user2.Role)
}

func TestDeleteUser(t *testing.T) {
	arg := CreateUserParams{
		Username: util.RandomUserName(),
		Role:     util.RandomRole(),
	}

	user1, err := testQueries.CreateUser(context.Background(), arg)
	err = testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		arg := CreateUserParams{
			Username: util.RandomUserName(),
			Role:     util.RandomRole(),
		}
		_, err := testQueries.CreateUser(context.Background(), arg)
		require.NoError(t, err)
	}
	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestUpdateUser(t *testing.T) {
	arg := CreateUserParams{
		Username: util.RandomUserName(),
		Role:     util.RandomRole(),
	}

	user1, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user1)

	arg2 := UpdateReputationParams{
		ID:         user1.ID,		
		Reputation: int32(util.RandomInt(0, 100)),
	}

	err = testQueries.UpdateReputation(context.Background(), arg2)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, arg2.Reputation, user2.Reputation)
}
