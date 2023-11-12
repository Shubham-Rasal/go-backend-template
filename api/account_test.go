package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Shubham-Rasal/blog-backend/db/mock"
	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/Shubham-Rasal/blog-backend/token"
	"github.com/Shubham-Rasal/blog-backend/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {

	user := createRandomUser()
	account := crateRandomAccount(user)

	testcases := []struct {
		name          string
		userId        int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, tokenMaker token.Maker, request *http.Request)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:   "OK",
			userId: int64(account.UserID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.UserID)).Times(1).Return(account, nil)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(account.Username)).Times(1).Return(user, nil)
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuth(t, tokenMaker, request, "Bearer", user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireAccountBodyMatch(t, account, response)
			},
		},
		{
			name:   "Bad Request",
			userId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuth(t, tokenMaker, request, "Bearer", user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:   "Not Found",
			userId: int64(account.UserID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.UserID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(account.Username)).Times(1).Return(user, nil)
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuth(t, tokenMaker, request, "Bearer", user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:   "Internal Server Error",
			userId: int64(account.UserID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.UserID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(account.Username)).Times(1).Return(user, nil)

			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuth(t, tokenMaker, request, "Bearer", user.Username, time.Minute)
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
			server, err := newTestServer(t, store)
			require.NoError(t, err)

			url := fmt.Sprintf("/accounts/%d", testcase.userId)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			require.NotEmpty(t, req)
			testcase.setupAuth(t, server.tokenMaker, req)

			response, err := server.router.Test(req)
			require.NoError(t, err)
			testcase.checkResponse(t, response)
		})
	}

}

func crateRandomAccount(user db.User) db.Account {
	return db.Account{
		Username: user.Username,
		Role:     util.RandomRole(),
		ID:       int64(util.RandomInt(1, 1000)),
		UserID:   int32(user.ID),
	}
}

func requireAccountBodyMatch(t *testing.T, user db.Account, response *http.Response) {
	u, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	var gotUser db.Account
	err = json.Unmarshal(u, &gotUser)
	require.NoError(t, err)
}
