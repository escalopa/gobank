package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/escalopa/go-bank/db/mock"
	db "github.com/escalopa/go-bank/db/sqlc"
	"github.com/escalopa/go-bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := createRandomAccount()

	tc := []struct {
		name      string
		accountId int64
		testCaseBase
	}{
		{
			name:      "OK",
			accountId: account.ID,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(account, nil)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchAccount(t, recorder.Body, account)
				},
			},
		},
		{
			name:      "ErrNoRows",
			accountId: account.ID,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.Account{}, sql.ErrNoRows)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				}},
		},
		{
			name:      "InternalError",
			accountId: account.ID,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.Account{}, sql.ErrConnDone)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				}},
		},
		{
			name:      "BadRequest",
			accountId: 0,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
		},
	}

	for i := 0; i < len(tc); i++ {
		tci := tc[i]

		t.Run(tci.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/accounts/%d", tci.accountId)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			runServerTest(t, tci, req)
		})
	}
}

func TestCreateAccount(t *testing.T) {
	account := createRandomAccount()

	arg := createAccountReq{
		Owner:    account.Owner,
		Currency: account.Currency,
	}

	tc := []struct {
		name       string
		accountArg createAccountReq
		testCaseBase
	}{
		{
			name:       "OK",
			accountArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
							Owner:    arg.Owner,
							Balance:  0,
							Currency: arg.Currency,
						})).
						Times(1).
						Return(account, nil)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusCreated, recorder.Code)
					requireBodyMatchAccount(t, recorder.Body, account)
				},
			},
		},
		{
			name:       "BadRequest",
			accountArg: createAccountReq{},
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
		},
		{
			name:       "InternalError",
			accountArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.Account{}, sql.ErrConnDone)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
		},
	}

	for i := 0; i < len(tc); i++ {
		tci := tc[i]

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(tci.accountArg)
		require.NoError(t, err)

		url := "/api/accounts"
		reader := bytes.NewReader(buf.Bytes())

		req, err := http.NewRequest(http.MethodPost, url, reader)
		require.NoError(t, err)

		runServerTest(t, tci, req)
	}
}

func TestListAccount(t *testing.T) {
	accounts := []db.Account{
		createRandomAccount(),
		createRandomAccount(),
		createRandomAccount(),
	}

	arg := listAccountReq{
		PageID:   2,
		PageSize: 5,
	}

	tc := []struct {
		name              string
		listAccountReqArg listAccountReq
		testCaseBase
	}{
		{
			name:              "OK",
			listAccountReqArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{
							Limit:  arg.PageSize,
							Offset: (arg.PageID - 1) * arg.PageSize,
						})).
						Times(1).
						Return(accounts, nil)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchAccounts(t, recorder.Body, accounts)
				},
			},
		},
		{
			name:              "BadRequest",
			listAccountReqArg: listAccountReq{},
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
		},
		{
			name:              "InternalError",
			listAccountReqArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Any()).
						Times(1).
						Return([]db.Account{}, sql.ErrTxDone)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
		},
	}

	for i := 0; i < len(tc); i++ {
		tci := tc[i]

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(tci.listAccountReqArg)
		require.NoError(t, err)

		url := fmt.Sprintf("/api/accounts?page_id=%d&page_size=%d", tci.listAccountReqArg.PageID, tci.listAccountReqArg.PageSize)
		reader := bytes.NewReader(buf.Bytes())

		req, err := http.NewRequest(http.MethodGet, url, reader)
		require.NoError(t, err)

		runServerTest(t, tci, req)
	}
}

func createRandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInteger(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, b *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(b)
	require.NoError(t, err)

	var accountReceived db.Account
	err = json.Unmarshal(data, &accountReceived)
	require.NoError(t, err)

	require.Equal(t, accountReceived, account)
}

func requireBodyMatchAccounts(t *testing.T, b *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(b)
	require.NoError(t, err)

	var accountsReceived []db.Account
	err = json.Unmarshal(data, &accountsReceived)
	require.NoError(t, err)

	require.Equal(t, accountsReceived, accounts)
}
