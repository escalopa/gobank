package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/escalopa/gobank/db/mock"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	user1, _ := createRandomUser(t)
	user2, _ := createRandomUser(t)

	account1 := createRandomAccount(user1.Username)
	account2 := createRandomAccount(user2.Username)

	amount := util.RandomInteger(1, 1000)
	account1.Currency, account2.Currency = util.EGP, util.EGP

	arg := createTransferReq{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	testCases := []struct {
		name          string
		transferArg   createTransferReq
		FromAccountID int64
		ToAccountID   int64
		testCaseBase
	}{
		{
			name:        "OK",
			transferArg: arg,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParam{
							FromAccountID: arg.FromAccountID,
							ToAccountID:   arg.ToAccountID,
							Amount:        arg.Amount,
						})).
						Times(1)

					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account2, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user1.Username, time.Minute)
				},
			},
		},
		{
			name: "BadRequest-Eq(IDS)",
			transferArg: createTransferReq{
				FromAccountID: account1.ID,
				ToAccountID:   account1.ID,
				Amount:        amount,
			},
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user1.Username, time.Minute)
				},
			},
		},
		{
			name:        "BadRequest-Binding",
			transferArg: createTransferReq{},
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user1.Username, time.Minute)
				},
			},
		},
		{
			name:        "InternalError",
			transferArg: arg,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParam{
							FromAccountID: arg.FromAccountID,
							ToAccountID:   arg.ToAccountID,
							Amount:        arg.Amount,
						})).
						Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)

					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account2, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, user1.Username, time.Minute)
				},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			data, err := json.Marshal(tc.transferArg)
			require.NoError(t, err)

			url := "/api/transfers"

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			runServerTest(t, tc, req)
		})
	}
}
