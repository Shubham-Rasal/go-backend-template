package db

import (
	"context"
	"testing"

	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/stretchr/testify/require"
)

func TestLikeTrasaction(t *testing.T) {

	store := NewStore(testDB)

	//create 5 users
	var accounts []Account
	for i := 0; i < 5; i++ {
		user := createDummyUser(t)
		arg := CreateAccountParams{
			Username: user.Username,
			Role:     util.RandomRole(),
			UserID:   int32(user.ID),
		}
		account, err := store.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, account)
		accounts = append(accounts, account)
	}

	//create 5 posts
	var posts []Post
	for i := 0; i < 5; i++ {
		arg := CreatePostParams{
			Title:  util.RandomString(7),
			Body:   util.RandomString(70),
			UserID: int32(accounts[i].UserID),
			Status: util.RandomString(7),
		}
		post, err := store.CreatePost(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, post)
		posts = append(posts, post)
	}

	errors := make(chan error)

	//like the posts
	for i := 0; i < 5; i++ {
		index := i
		go func() {
			arg := LikePostParams{
				UserID: int32(accounts[index].UserID),
				PostID: int32(posts[index].ID),
			}
			err := store.LikeTx(context.Background(), arg)
			errors <- err
		}()
	}

	for i := 0; i < 5; i++ {
		err := <-errors
		require.NoError(t, err)
	}

}
