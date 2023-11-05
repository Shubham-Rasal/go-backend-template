package db

import (
	"context"
	"testing"
	"time"

	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomPost(t *testing.T) Post {
	userDetails := CreateUserParams{
		Username: util.RandomUserName(),
		Role:     util.RandomRole(),
	}

	user, err := testQueries.CreateUser(context.Background(), userDetails)
	require.NoError(t, err)

	arg := CreatePostParams{
		Title:     util.RandomString(7),
		Body:      util.RandomString(7),
		UserID:    int32(user.ID),
		Status:    util.RandomString(7),
	}

	post, err := testQueries.CreatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	require.Equal(t, arg.Title, post.Title)
	require.Equal(t, arg.Body, post.Body)
	require.Equal(t, arg.UserID, post.UserID)
	require.Equal(t, arg.Status, post.Status)

	return post
}

func TestCreatePost(t *testing.T) {
	createRandomPost(t)
}

func TestGetPost(t *testing.T) {
	post1 := createRandomPost(t)
	post2, err := testQueries.GetPost(context.Background(), post1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.ID, post2.ID)
	require.Equal(t, post1.Title, post2.Title)
	require.Equal(t, post1.Body, post2.Body)
	require.Equal(t, post1.UserID, post2.UserID)
	require.Equal(t, post1.Status, post2.Status)
	require.WithinDuration(t, post1.CreatedAt, post2.CreatedAt, time.Second)
	require.Equal(t, post1.Likes, post2.Likes)
}

func TestListPosts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPost(t)
	}

	arg := ListPostsParams{
		Limit:  5,
		Offset: 5,
	}

	posts, err := testQueries.ListPosts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, posts, 5)

	for _, post := range posts {
		require.NotEmpty(t, post)
	}
}

func TestLikePost(t *testing.T) {
	post := createRandomPost(t)
	err := testQueries.LikePost(context.Background(), post.ID)
	require.NoError(t, err)

	// Check if the post was actually liked
	likedPost, err := testQueries.GetPost(context.Background(), post.ID)
	require.NoError(t, err)
	require.Equal(t, post.ID, likedPost.ID)
	require.Equal(t, post.Title, likedPost.Title)
	require.Equal(t, post.Likes+1, likedPost.Likes)
}

// func TestUpdatePost(t *testing.T) {
// 	post1 := createRandomPost(t)

// 	arg := UpdatePostParams{
// 		ID:      post1.ID,
// 		Title:   util.RandomTitle(),
// 		Body:    util.RandomBody(),
// 		Status:  util.RandomStatus(),
// 		Likes:   util.RandomLikes(),
// 		Updated: time.Now(),
// 	}

// 	err := testQueries.UpdatePost(context.Background(), arg)
// 	require.NoError(t, err)

// 	post2, err := testQueries.GetPost(context.Background(), post1.ID)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, post2)

// 	require.Equal(t, post1.ID, post2.ID)
// 	require.Equal(t, arg.Title, post2.Title)
// 	require.Equal(t, arg.Body, post2.Body)
// 	require.Equal(t, post1.UserID, post2.UserID)
// 	require.Equal(t, arg.Status, post2.Status)
// 	require.WithinDuration(t, arg.Updated, post2.Updated, time.Second)
// 	require.Equal(t, arg.Likes, post2.Likes)
// }

// func TestDeletePost(t *testing.T) {
// 	post1 := createRandomPost(t)

// 	err := testQueries.DeletePost(context.Background(), post1.ID)
// 	require.NoError(t, err)

// 	post2, err := testQueries.GetPost(context.Background(), post1.ID)
// 	require.Error(t, err)
// 	require.Empty(t, post2)
// }
