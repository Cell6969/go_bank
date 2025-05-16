package gapi

import (
	"testing"
	"time"

	db "github.com/Cell6969/go_bank/db/sqlc"
	"github.com/Cell6969/go_bank/util"
	"github.com/Cell6969/go_bank/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenKey:      util.RandomString(32),
		TokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}
