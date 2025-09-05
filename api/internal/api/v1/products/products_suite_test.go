package products_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/products"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestProducts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Products Handler Suite")
}

var (
	mockCtrl          *gomock.Controller
	mockGlobalRepo    *mock_repos.MockGlobalRepo
	mockProductsRepo  *mock_repos.MockProductsRepo
	mockCommoditiesRepo *mock_repos.MockCommoditiesRepo
	router            *mux.Router
	adminUser         *types.User
	normalUser        *types.User
	company           *types.Company
)

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockGlobalRepo = mock_repos.NewMockGlobalRepo(mockCtrl)
	mockProductsRepo = mock_repos.NewMockProductsRepo(mockCtrl)
	mockCommoditiesRepo = mock_repos.NewMockCommoditiesRepo(mockCtrl)

	// Set up the mock chain
	mockGlobalRepo.EXPECT().Products().Return(mockProductsRepo).AnyTimes()
	mockGlobalRepo.EXPECT().Commodities().Return(mockCommoditiesRepo).AnyTimes()

	// Set up the router
	router = mux.NewRouter()
	products.AddRoutes(router)

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
