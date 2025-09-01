package companies_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestCompanies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Companies Suite")
}

var (
	mockCtrl          *gomock.Controller
	mockGlobalRepo    *mock_repos.MockGlobalRepo
	mockCompaniesRepo *mock_repos.MockCompaniesRepo
	mockUsersRepo     *mock_repos.MockUsersRepo
	rr                *httptest.ResponseRecorder
	router            *mux.Router
	adminUser         *types.User
	basicUser         *types.User
)

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockGlobalRepo = mock_repos.NewMockGlobalRepo(mockCtrl)
	mockCompaniesRepo = mock_repos.NewMockCompaniesRepo(mockCtrl)
	mockUsersRepo = mock_repos.NewMockUsersRepo(mockCtrl)

	// Set up the mock chain: GlobalRepo -> Companies() -> MockCompaniesRepo
	mockGlobalRepo.EXPECT().Companies().Return(mockCompaniesRepo).AnyTimes()
	mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo).AnyTimes()

	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	companies.AddRoutes(router)

	adminUser = &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
	basicUser = &types.User{ID: 2, Roles: types.Roles{types.RoleUser}}
})

var _ = AfterEach(func() {
	mockCtrl.Finish()
})

func newAuthenticatedRequest(method, url string, body io.Reader, user *types.User) *http.Request {
	req, err := http.NewRequest(method, url, body)
	Expect(err).NotTo(HaveOccurred())

	// Inject the mocked GlobalRepo into the request's context
	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)

	if user != nil {
		ctxWithRepo = middleware.AddUserIDToContext(ctxWithRepo, user.ID)
		mockUsersRepo.EXPECT().Get(gomock.Any(), user.ID).Return(user, true, nil).AnyTimes()
	}
	return req.WithContext(ctxWithRepo)
}
