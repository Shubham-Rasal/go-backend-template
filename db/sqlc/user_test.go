package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/stretchr/testify/require"
)

func createDummyUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomString(4),
		Password: util.RandomString(6),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	return user
}

func TestCreateUser(t *testing.T) {

	arg := CreateUserParams{
		Username: util.RandomString(4),
		Password: util.RandomString(6),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

}

func TestGetUser(t *testing.T) {

	user1 := createDummyUser(t)

	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)

	//get by username
	user3, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user3)
	require.Equal(t, user1.ID, user3.ID)
	require.Equal(t, user1.Username, user3.Username)
	require.Equal(t, user1.Email, user3.Email)
	require.Equal(t, user1.Password, user3.Password)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createDummyUser(t)
	}
	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, account := range users {
		require.NotEmpty(t, account)
	}
}

func TestUpdateEmail(t *testing.T) {
	// Create a dummy user
	user := createDummyUser(t)

	// Generate a new email
	newEmail := util.RandomEmail()

	// Update the user's email
	arg := UpdateEmailParams{
		ID:    user.ID,
		Email: newEmail,
	}
	err := testQueries.UpdateEmail(context.Background(), arg)
	require.NoError(t, err)

	// Retrieve the updated user from the database
	updatedUser, err := testQueries.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)

	// Verify that the user's email has been updated
	require.Equal(t, newEmail, updatedUser.Email)
}

func TestUpdatePassword(t *testing.T) {
	// Create a dummy user
	user := createDummyUser(t)

	// Generate a new password
	newPassword := util.RandomString(6)

	// Update the user's password
	arg := UpdatePasswordParams{
		ID:       user.ID,
		Password: newPassword,
	}
	err := testQueries.UpdatePassword(context.Background(), arg)
	require.NoError(t, err)

	// Retrieve the updated user from the database
	updatedUser, err := testQueries.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)

	// Verify that the user's password has been updated
	require.Equal(t, newPassword, updatedUser.Password)
}

func TestDeleteUser(t *testing.T) {
	// Create a dummy user
	user := createDummyUser(t)

	// Delete the user from the database
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	// Try to retrieve the deleted user from the database
	deletedUser, err := testQueries.GetUserById(context.Background(), user.ID)

	// Verify that the user has been deleted
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedUser)
}
