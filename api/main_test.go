package api

//same package as crud code
//main_test used to set up connection with query object

import (
	"os"
	"testing"
	"time"

	db "github.com/Viczdera/bank-gorm/db/sqlc"
	"github.com/Viczdera/bank-gorm/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		AccessTokenSymmetrickey: util.RandomString(32),
		AccessTokenDuration:     time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

// testing
func TestMain(t *testing.M) {

	gin.SetMode(gin.TestMode)

	//run and report back to test runner via the o.exit command
	os.Exit(t.Run())

}
