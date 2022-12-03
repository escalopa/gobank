package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

func TestListAccount(t *testing.T) {
	user, _ := createRandomUser(t)

	accounts := []db.Account{
		createRandomAccount(user.Username),
		createRandomAccount(user.Username),
		createRandomAccount(user.Username),
	}

	arg := listAccountReq{
		PageID:   2,
		PageSize: 5,
	}

	testCases := []struct {
		name              string
		listAccountReqArg listAccountReq
		testCaseBase
	}{
		{
			name:              "OK",
			listAccountReqArg: arg,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{
							Owner:  user.Username,
							Limit:  arg.PageSize,
							Offset: (arg.PageID - 1) * arg.PageSize,
						})).
						Times(1).
						Return(accounts, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchAccounts(t, recorder.Body, accounts)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
				},
			},
		},
		{
			name:              "BadRequest",
			listAccountReqArg: listAccountReq{},
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Any()).
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
		{
			name:              "InternalError",
			listAccountReqArg: arg,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Any()).
						Times(1).
						Return([]db.Account{}, sql.ErrTxDone)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
				},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(tc.listAccountReqArg)
		require.NoError(t, err)

		url := fmt.Sprintf("/api/accounts?page_id=%d&page_size=%d", tc.listAccountReqArg.PageID, tc.listAccountReqArg.PageSize)
		reader := bytes.NewReader(buf.Bytes())

		req, err := http.NewRequest(http.MethodGet, url, reader)
		require.NoError(t, err)

		runServerTest(t, tc, req)
	}
}
