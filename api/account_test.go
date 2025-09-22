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

	mockdb "github.com/Viczdera/bank-gorm/db/mock"
	db "github.com/Viczdera/bank-gorm/db/sqlc"
	"github.com/Viczdera/bank-gorm/token"
	"github.com/Viczdera/bank-gorm/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := genRamdomAccount()

	testCases := []struct {
		name           string
		accountID      int64
		setupAuth      func(*testing.T, *http.Request, token.Maker, string)
		buildStubs     func(store *mockdb.MockStore)
		requiredChecks func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)). //can be called with any context (hence Any)
												Times(1). //and specific account arguements (hence Eq)
												Return(account, nil)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0, // invalid ID to simulate bad request
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "ResponseBody",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)). //can be called with any context (hence Any)
												Times(1). //and specific account arguements (hence Eq)
												Return(account, nil)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response body
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var resAccount db.Account
				err = json.Unmarshal(body, &resAccount)
				require.NoError(t, err)
				require.Equal(t, account, resAccount)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)). //can be called with any context (hence Any)
												Times(1). //and specific account arguements (hence Eq)
												Return(db.Account{}, sql.ErrNoRows)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)). //can be called with any context (hence Any)
												Times(1). //and specific account arguements (hence Eq)
												Return(db.Account{}, sql.ErrConnDone)
			},
			requiredChecks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {

			//new controller instance
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			//build stubs
			testCase.buildStubs(store)

			//create and start test server and test request. dont really have to start a new server
			//but can use a recorder feature of the http test package
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)

			//check error
			require.NoError(t, err)

			testCase.setupAuth(t, req, server.tokenMaker, account.Owner)
			server.router.ServeHTTP(recorder, req)

			testCase.requiredChecks(t, recorder)

		})
	}

}

func genRamdomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomBalance(),
	}
}
