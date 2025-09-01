package users_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	"github.com/happilymarrieddad/order-management-v3/api/types"
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
	router            *mux.Router
	adminUser         *types.User
	basicUser         *types.User
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
	router = mux.NewRouter()
	users.AddRoutes(router)

	adminUser = &types.User{ID: 1, CompanyID: 1, Roles: types.Roles{types.RoleAdmin}}
	basicUser = &types.User{ID: 2, CompanyID: 2, Roles: types.Roles{types.RoleUser}}
})

var _ = AfterEach(func() {
	mockCtrl.Finish()
})

// newAuthenticatedRequest creates a new HTTP request with the mocked repository
// injected into the context.
func newAuthenticatedRequest(method, url string, body io.Reader, user *types.User) *http.Request {
	req, err := http.NewRequest(method, url, body)
	Expect(err).NotTo(HaveOccurred())

	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)

	if user != nil {
		ctxWithRepo = middleware.AddUserIDToContext(ctxWithRepo, user.ID)
		// Mock the UsersRepo to return the user for role checks
		mockUsersRepo.EXPECT().Get(gomock.Any(), user.ID).Return(user, true, nil).AnyTimes()
	}

	return req.WithContext(ctxWithRepo)
}
