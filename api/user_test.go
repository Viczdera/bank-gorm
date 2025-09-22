package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Viczdera/bank-gorm/db/mock"
	db "github.com/Viczdera/bank-gorm/db/sqlc"
	"github.com/Viczdera/bank-gorm/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	//check if pass is hashed
	err := util.CheckPassword(e.password, arg.PasswordHashed)
	if err != nil {
		return false
	}

	e.arg.PasswordHashed = arg.PasswordHashed
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		arg:      arg,
		password: password,
	}
}

func TestCreateUserAPI(t *testing.T) {

	user, password := genRandomUser()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		body           gin.H
		buildStubs     func(store *mockdb.MockStore)
		requiredChecks func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					PasswordHashed: hashedPassword,
					FullName:       user.FullName,
					Email:          user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). //can be called with any context (hence Any)
														Times(1). //and specific account arguements (hence Eq)
														Return(user, nil)
			},
			requiredChecks: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"username":  user.Username,
				"email":     12345,
				"password":  password,
				"full_name": user.FullName,
			}, // invalid email to simulate bad request
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			requiredChecks: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ResponseBody",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()). //can be called with any context (hence Any)
											Times(1). //and specific account arguements (hence Eq)
											Return(user, nil)
			},
			requiredChecks: func(recorder *httptest.ResponseRecorder) {
				//check response body
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var resUser db.User
				err = json.Unmarshal(body, &resUser)
				require.NoError(t, err)
				require.Equal(t, user, resUser)
			},
		},

		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()). //can be called with any context (hence Any)
											Times(1). //and specific account arguements (hence Eq)
											Return(db.User{}, sql.ErrConnDone)
			},
			requiredChecks: func(recorder *httptest.ResponseRecorder) {
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

			url := "/users"
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			//check error
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			testCase.requiredChecks(recorder)

		})
	}

}

func genRandomUser() (db.User, string) {
	password := util.RandomString(10)

	user := db.User{
		Username: util.RandomString(5),
		Email:    util.RandomEmail(),
		FullName: util.RandomString(10),
	}

	return user, password

}
