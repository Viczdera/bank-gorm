package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Viczdera/bank-gorm/db/mock"
	db "github.com/Viczdera/bank-gorm/db/sqlc"
	"github.com/Viczdera/bank-gorm/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	fromAccount := db.Account{ID: 1, Owner: "alice", Currency: "USD", Balance: 1000}
	toAccount := db.Account{ID: 2, Owner: "bob", Currency: "USD", Balance: 500}
	amount := int64(100)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(*testing.T, *http.Request, token.Maker, string)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account": fromAccount.ID,
				"to_account":   toAccount.ID,
				"currency":     fromAccount.Currency,
				"amount":       amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Validate from account
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)
				// Validate to account
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(toAccount, nil)
				// Transfer
				arg := db.TransferTxParams{
					FromAccount: fromAccount.ID,
					ToAccount:   toAccount.ID,
					Amount:      amount,
				}
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TransferTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account": 99,
				"to_account":   toAccount.ID,
				"currency":     fromAccount.Currency,
				"amount":       amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(99))).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "CurrencyMismatch",
			body: gin.H{
				"from_account": fromAccount.ID,
				"to_account":   toAccount.ID,
				"currency":     "EUR",
				"amount":       amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, owner string) {
				addAuthHeader(t, request, tokenMaker, AUTH_TYPE, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker, fromAccount.Owner)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
