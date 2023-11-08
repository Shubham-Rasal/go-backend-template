package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	mockdb "github.com/Shubham-Rasal/blog-backend/db/mock"
	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetUserAPI(t *testing.T) {

	user := crateRandomUser()

	testcases := []struct {
		name          string
		userId        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:   "OK",
			userId: int64(user.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatch(t, user, response)
			},
		},
		{
			name:   "Bad Request",
			userId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:   "Not Found",
			userId: int64(user.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:   "Internal Server Error",
			userId: int64(user.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
	}

	for i := range testcases {
		testcase := testcases[i]

		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			//create a random user
			testcase.buildStubs(store)

			//creates a serer as normals
			server , err := newTestServer(t, store)
			require.NoError(t, err)

			url := fmt.Sprintf("/users/%d", testcase.userId)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			response, err := server.router.Test(req)
			require.NoError(t, err)
			testcase.checkResponse(t, response)
		})
	}

}

func crateRandomUser() db.User {
	return db.User{
		Username: util.RandomUserName(),
		Role:     util.RandomRole(),
		ID:       int64(util.RandomInt(1, 1000)),
	}
}

func requireBodyMatch(t *testing.T, user db.User, response *http.Response) {
	u, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	var gotUser db.User
	err = json.Unmarshal(u, &gotUser)
	require.NoError(t, err)
}
