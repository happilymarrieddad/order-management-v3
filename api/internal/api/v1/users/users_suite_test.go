package users_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestUsers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Suite")
}

var (
	mockCtrl          *gomock.Controller
	mockGlobalRepo    *mock_repos.MockGlobalRepo
	mockUsersRepo     *mock_repos.MockUsersRepo
	mockCompaniesRepo *mock_repos.MockCompaniesRepo
	mockAddressesRepo *mock_repos.MockAddressesRepo
	rr                *httptest.ResponseRecorder
	ctx               context.Context
)

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockGlobalRepo = mock_repos.NewMockGlobalRepo(mockCtrl)
	mockUsersRepo = mock_repos.NewMockUsersRepo(mockCtrl)
	mockCompaniesRepo = mock_repos.NewMockCompaniesRepo(mockCtrl)
	mockAddressesRepo = mock_repos.NewMockAddressesRepo(mockCtrl)

	// Set up the mock chain for all required repositories
	mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo).AnyTimes()
	mockGlobalRepo.EXPECT().Companies().Return(mockCompaniesRepo).AnyTimes()
	mockGlobalRepo.EXPECT().Addresses().Return(mockAddressesRepo).AnyTimes()

	rr = httptest.NewRecorder()
	ctx = context.Background()
})

var _ = AfterEach(func() {
	mockCtrl.Finish()
})

// createRequestWithRepo creates a new HTTP request with the mocked repository
// injected into the context.
func createRequestWithRepo(method, url string, body []byte, vars map[string]string) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	Expect(err).NotTo(HaveOccurred())

	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
	req = req.WithContext(ctxWithRepo)

	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}

	return req
}
