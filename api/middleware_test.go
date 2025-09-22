package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Viczdera/bank-gorm/token"
	"github.com/Viczdera/bank-gorm/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthHeader(t *testing.T, request *http.Request, tokenMaker token.Maker, authType string, username string, time time.Duration) {
	accessToken, err := tokenMaker.CreateToken(username, time)
	require.NoError(t, err)

	//setauth header
	authorizationHeader := fmt.Sprintf("%s %s", authType, accessToken)
	request.Header.Set(AUTH_HEADER_KEY, authorizationHeader)
}
func TestAuthMiddleware(t *testing.T) {

	testCases := []struct {
		name           string
		setupAuth      func(*testing.T, *http.Request, token.Maker)
		requiredChecks func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, util.RandomOwner(), time.Minute)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, "unsupported_type", util.RandomOwner(), time.Minute)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationHeader",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, "", util.RandomOwner(), time.Minute)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, "", util.RandomOwner(), -time.Minute)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			//create dummy server
			server := newTestServer(t, nil)

			//dummy auth path
			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				})

			//recorder to record the request
			recoder := httptest.NewRecorder()

			//make requests
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recoder, request)

			//verify reesults
			testCase.requiredChecks(t, recoder)
		})
	}

}
