package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthHeader(
	t *testing.T,
	request *http.Request,
	maker token.Maker,
	authHeaderType string,
	username string,
) {
	token, payload, err := maker.CreateToken(username)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authHeader := fmt.Sprintf("%s %s", authHeaderType, token)
	request.Header.Add(authorizationHeaderKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {

	testCases := []struct {
		name string
		testCaseBase
	}{
		{
			name: "OK",
			testCaseBase: testCaseBase{
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, authorizationTypeBearer, util.RandomString(6))
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				}},
		},
		{
			name: "Unauthorized",
			testCaseBase: testCaseBase{
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
		},
		{
			name: "InvalidHeaderFormat",
			testCaseBase: testCaseBase{
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, "", util.RandomString(6))
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				}},
		},
		{
			name: "UnsupportedAuthType",
			testCaseBase: testCaseBase{
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					addAuthHeader(t, req, maker, "OAuth", util.RandomString(6))
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				}},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"

			server.router.GET(authPath, authMiddleware(server.tm), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			tc.setupAuth(t, req, server.tm)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)

		})

	}
}
