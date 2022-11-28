package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/escalopa/go-bank/db/mock"
	db "github.com/escalopa/go-bank/db/sqlc"
	"github.com/escalopa/go-bank/util"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user, password := createRadomUser(t)
	uniqueViolationError := &pq.Error{Code: "23505"}

	arg := createUserReq{
		Username:        user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		Password:        password,
		PasswordConfirm: password,
	}

	tc := []struct {
		name    string
		userArg createUserReq
		testCaseBase
	}{
		{
			name:    "OK",
			userArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(user, nil)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusCreated, recorder.Code)
					requireBodyMatchUser(t, recorder.Body, user)
				},
			},
		},
		{
			name: "BadRequest",
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
		},
		{
			name:    "SQLUserNameViolation",
			userArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, uniqueViolationError)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusForbidden, recorder.Code)
				},
			},
		},
		{
			name:    "SQLEmailViolation",
			userArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, uniqueViolationError)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusForbidden, recorder.Code)
				},
			},
		},
		{
			name:    "InternalErrorDB",
			userArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, sql.ErrConnDone)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
		},
	}

	for i := 0; i < len(tc); i++ {
		tci := tc[i]

		url := "/api/users"

		data, err := json.Marshal(tci.userArg)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
		require.NoError(t, err)

		runServerTest(t, tci, req)
	}

}

func requireBodyMatchUser(t *testing.T, b io.Reader, user db.User) {
	data, err := io.ReadAll(b)
	require.NoError(t, err)

	var userReceived db.User
	err = json.Unmarshal(data, &userReceived)
	require.NoError(t, err)
}

func createRadomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)

	hashPassword, err := util.GenerateHashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return user, password

}
