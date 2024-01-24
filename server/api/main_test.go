package api

import (
	db "github.com/bysergr/simple-bank/db/sqlc"
	"github.com/bysergr/simple-bank/token"
	"github.com/bysergr/simple-bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	return NewServer(store, maker, config)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
