package api

import (
	"bytes"
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

func createRandomUser() db.User {
	return db.User{
		Username: util.RandomUserName(),
		Password: util.RandomString(6),
		Email:    util.RandomEmail(),
		ID:       int64(util.RandomInt(1, 1000)),
	}
}

func TestGetUserAPI(t *testing.T) {

	user := createRandomUser()

	testcases := []struct {
		name          string
		username      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatch(t, user, response)
			},
		},
		{
			name:     "Bad Request",
			username: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:     "Not Found",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:     "Internal Server Error",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/users/%s", testcase.username)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			fmt.Println("req : ", testcase.name)
			response, err := server.router.Test(req)
			require.NoError(t, err)
			testcase.checkResponse(t, response)
		})
	}

}

func requireBodyMatch(t *testing.T, user db.User, response *http.Response) {
	u, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	var gotUser db.User
	err = json.Unmarshal(u, &gotUser)
	require.NoError(t, err)
}

func TestCreateUserAPI(t *testing.T) {
	user := createRandomUser()
	arg := db.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	testcases := []struct {
		name          string
		body          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, arg.Username, arg.Password, arg.Email),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusCreated, response.StatusCode)
				requireBodyMatch(t, user, response)
			},
		},
		{
			name: "Bad Request",
			body: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "Internal Server Error",
			body: fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, arg.Username, arg.Password, arg.Email),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.User{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/users")

			//create the body
			body := []byte(testcase.body)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			response, err := server.router.Test(req)
			require.NoError(t, err)
			testcase.checkResponse(t, response)
		})
	}
}
