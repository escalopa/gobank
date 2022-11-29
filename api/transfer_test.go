package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/escalopa/go-bank/db/mock"
	db "github.com/escalopa/go-bank/db/sqlc"
	"github.com/escalopa/go-bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1, account2, amount := createRandomAccount(), createRandomAccount(), util.RandomInteger(1, 1000)
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
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParam{
							FromAccountID: arg.FromAccountID,
							ToAccountID:   arg.ToAccountID,
							Amount:        arg.Amount,
						})).
						Times(1)

					store.EXPECT().
						GetAccount(gomock.Any(), gomock.All()).
						Times(2)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				},
			},
		},
		{
			name: "BadRequest-Eq(IDS)",
			transferArg: createTransferReq{
				FromAccountID: 1,
				ToAccountID:   1,
				Amount:        10,
			},
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(2)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
		},
		{
			name:        "BadRequest-Binding",
			transferArg: createTransferReq{},
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Any()).
						Times(0)

					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
		},
		{
			name:        "InternalError",
			transferArg: arg,
			testCaseBase: testCaseBase{
				buildStubsMethod: func(store *mockdb.MockStore) {
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParam{
							FromAccountID: arg.FromAccountID,
							ToAccountID:   arg.ToAccountID,
							Amount:        arg.Amount,
						})).
						Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)

					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(2)
				},
				checkResponseMethod: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
