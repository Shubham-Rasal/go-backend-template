package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/Shubham-Rasal/blog-backend/util"
)

func newTestServer(t *testing.T, store db.Store) (*Server, error) {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		AccessDuration:    time.Minute,
	}

	server, err := NewServer(config, store)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
