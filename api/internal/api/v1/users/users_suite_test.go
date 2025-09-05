package users_test

import (
	"context"
	"io"
	"net/http"
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
	RunSpecs(t, "Users Handler Suite")
}

var (
	mockCtrl          *gomock.Controller
	mockGlobalRepo    *mock_repos.MockGlobalRepo
	mockUsersRepo     *mock_repos.MockUsersRepo
	mockCompaniesRepo *mock_repos.MockCompaniesRepo
	mockAddressesRepo *mock_repos.MockAddressesRepo
	router            *mux.Router
	adminUser         *types.User
	normalUser        *types.User
	company           *types.Company
)

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockGlobalRepo = mock_repos.NewMockGlobalRepo(mockCtrl)
	mockUsersRepo = mock_repos.NewMockUsersRepo(mockCtrl)
	mockCompaniesRepo = mock_repos.NewMockCompaniesRepo(mockCtrl)
	mockAddressesRepo = mock_repos.NewMockAddressesRepo(mockCtrl)

	// Set up the mock chain
	mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo).AnyTimes()
	mockGlobalRepo.EXPECT().Companies().Return(mockCompaniesRepo).AnyTimes()
	mockGlobalRepo.EXPECT().Addresses().Return(mockAddressesRepo).AnyTimes()

	// Set up the router
	router = mux.NewRouter()
	users.AddRoutes(router)

	// Set up common test data
	company = &types.Company{ID: 1, Name: "Test Company"}
	normalUser = &types.User{ID: 1, CompanyID: company.ID, Roles: types.Roles{types.RoleUser}}
	adminUser = &types.User{ID: 2, CompanyID: company.ID, Roles: types.Roles{types.RoleAdmin}}
})

var _ = AfterEach(func() {
	mockCtrl.Finish()
})

func newAuthenticatedRequest(method, url string, body io.Reader, user *types.User) *http.Request {
	req, err := http.NewRequest(method, url, body)
	Expect(err).ToNot(HaveOccurred())

	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
	if user != nil {
		ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, user)
		return req.WithContext(ctxWithAuth)
	}
	return req.WithContext(ctxWithRepo)
}
