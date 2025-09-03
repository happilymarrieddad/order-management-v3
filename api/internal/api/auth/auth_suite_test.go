package auth_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/auth"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var (
	mockCtrl       *gomock.Controller
	mockGlobalRepo *mock_repos.MockGlobalRepo
	mockUsersRepo  *mock_repos.MockUsersRepo
	router         *mux.Router
)

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockGlobalRepo = mock_repos.NewMockGlobalRepo(mockCtrl)
	mockUsersRepo = mock_repos.NewMockUsersRepo(mockCtrl)

	// Set up the mock chain
	mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo).AnyTimes()

	// Set up the router for auth handlers
	router = mux.NewRouter()
	router.HandleFunc("/login", auth.Login).Methods("POST")
})

var _ = AfterEach(func() {
	mockCtrl.Finish()
})

// newAuthenticatedRequest creates a new http.Request with the mocked GlobalRepo
// and an optional authenticated user in the context.
func newAuthenticatedRequest(method, url string, body []byte, user *types.User) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
	if user != nil {
		ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, user)
		return req.WithContext(ctxWithAuth)
	}
	return req.WithContext(ctxWithRepo)
}
