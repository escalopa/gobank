package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/escalopa/gobank/db/mock"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdateUser(t *testing.T) {
	var arg db.UpdateUserParams

	testCases := []struct {
		name    string
		userArg db.UpdateUserParams
		testCaseBase
	}{
		{
			name:    "OK",
			userArg: arg,
			testCaseBase: testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						UpdateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {

				},
			},
		},
	}

	_ = testCases
}
