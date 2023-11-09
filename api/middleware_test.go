package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Shubham-Rasal/blog-backend/token"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func addAuth(t *testing.T, tokenMaker token.Maker, request *http.Request, authType string, username string, duration time.Duration) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	request.Header.Set("Authorization", authType+" "+token)
}

func TestAuthMiddleware(t *testing.T) {

	testcases := []struct {
		name          string
		setupAuth     func(t *testing.T, tokenMaker token.Maker, request *http.Request)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuth(t, tokenMaker, request, "Bearer", "shubham", time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name: "No Token",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "Invalid Auth Type",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuth(t, tokenMaker, request, "Invalid", "shubham", time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "Expired Token",
			setupAuth: func(t *testing.T, tokenMaker token.Maker, request *http.Request) {
				addAuth(t, tokenMaker, request, "Invalid", "shubham", -time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
	}

	for i := range testcases {
		testcase := testcases[i]
		t.Run(testcase.name, func(t *testing.T) {

			//setup a server
			server, err := newTestServer(t, nil)
			require.NoError(t, err)

			//add a auth route
			server.router.Post("/test-auth", Protected(server.tokenMaker), func(ctx *fiber.Ctx) error {
				return ctx.Status(http.StatusOK).JSON(fiber.Map{})
			})

			//construct a new request
			request := httptest.NewRequest(http.MethodPost, "/test-auth", nil)
			testcase.setupAuth(t, server.tokenMaker, request)

			//make a request
			response, err := server.router.Test(request)
			require.NoError(t, err)
			testcase.checkResponse(t, response)

		})
	}

}
