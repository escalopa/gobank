package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mockdb "github.com/escalopa/gobank/db/mock"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type testCase interface {
	buildStubsMethod(store *mockdb.MockStore)
	checkResponseMethod(t *testing.T, response *httptest.ResponseRecorder)
	setupAuthMethod(t *testing.T, request *http.Request, maker token.Maker)
}

type testCaseBase struct {
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, response *httptest.ResponseRecorder)
	setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
}

func (tcb testCaseBase) buildStubsMethod(store *mockdb.MockStore) { tcb.buildStubs(store) }
func (tcb testCaseBase) checkResponseMethod(t *testing.T, response *httptest.ResponseRecorder) {
	tcb.checkResponse(t, response)
}
func (tcb testCaseBase) setupAuthMethod(t *testing.T, request *http.Request, maker token.Maker) {
	tcb.setupAuth(t, request, maker)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func runServerTest(t *testing.T, tc testCase, req *http.Request) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	tc.buildStubsMethod(store)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	tc.setupAuthMethod(t, req, server.tokenMaker)
	server.router.ServeHTTP(recorder, req)
	tc.checkResponseMethod(t, recorder)
}

func newTestServer(t *testing.T, store db.Store) *Server {

	testConfig := util.Config{}
	testConfig.TokenSymmetricKey = util.RandomString(32)

	server, err := NewServer(testConfig, store)
	require.NoError(t, err)
	return server
}
