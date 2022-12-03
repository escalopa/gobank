package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/escalopa/gobank/db/mock"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	user, _ := createRandomUser(t)
	account := createRandomAccount(user.Username)

	testCases := []struct {
		name      string
		accountId int64
		testCaseBase
	}{
		{
			name:      "OK",
			accountId: account.ID,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(account, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchAccount(t, recorder.Body, account)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
				},
			},
		},
		{
			name:      "ErrNoRows",
			accountId: account.ID,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.Account{}, sql.ErrNoRows)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
				},
			},
		},
		{
			name:      "InternalError",
			accountId: account.ID,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.Account{}, sql.ErrConnDone)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
				},
			},
		},
		{
			name:      "BadRequest",
			accountId: 0,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
				},
			},
		},
		// {
		// 	name:      "Unauthorized",
		// 	accountId: 0,
		// 	testCaseBase: testCaseBase{
		// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 		},
		// 	},
		// },
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/accounts/%d", tc.accountId)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			runServerTest(t, tc, req)
		})
	}
}
