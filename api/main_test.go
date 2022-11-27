package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mockdb "github.com/escalopa/go-bank/db/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

type testCase interface {
	buildStubs(store *mockdb.MockStore)
	checkResponse(t *testing.T, response *httptest.ResponseRecorder)
}

type testCaseBase struct {
	buildStubsMethod    func(store *mockdb.MockStore)
	checkResponseMethod func(t *testing.T, response *httptest.ResponseRecorder)
}

func (tcb testCaseBase) buildStubs(store *mockdb.MockStore) { tcb.buildStubsMethod(store) }
func (tcb testCaseBase) checkResponse(t *testing.T, response *httptest.ResponseRecorder) {
	tcb.checkResponseMethod(t, response)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func runServerTest(t *testing.T, tc testCase, req *http.Request) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	tc.buildStubs(store)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	server.router.ServeHTTP(recorder, req)
	tc.checkResponse(t, recorder)
}
