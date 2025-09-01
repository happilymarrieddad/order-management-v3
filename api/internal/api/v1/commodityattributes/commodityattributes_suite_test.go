package commodityattributes_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	. "github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestCommodityAttributes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commodity Attributes Suite")
}

var (
	mockCtrl                    *gomock.Controller
	mockGlobalRepo              *mock_repos.MockGlobalRepo
	mockCommodityAttributesRepo *mock_repos.MockCommodityAttributesRepo
	mockUsersRepo               *mock_repos.MockUsersRepo
	rr                          *httptest.ResponseRecorder
	router                      *mux.Router
	ctx                         context.Context
	adminUser                   *types.User
	basicUser                   *types.User
)

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockGlobalRepo = mock_repos.NewMockGlobalRepo(mockCtrl)
	mockCommodityAttributesRepo = mock_repos.NewMockCommodityAttributesRepo(mockCtrl)
	mockUsersRepo = mock_repos.NewMockUsersRepo(mockCtrl)

	// Set up the mock chain
	mockGlobalRepo.EXPECT().CommodityAttributes().Return(mockCommodityAttributesRepo).AnyTimes()
	mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo).AnyTimes()

	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	AddRoutes(router)
	ctx = context.Background()

	adminUser = &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
	basicUser = &types.User{ID: 2, Roles: types.Roles{types.RoleUser}}
})

var _ = AfterEach(func() {
	mockCtrl.Finish()
})

// newAuthenticatedRequest creates a new HTTP request with the mocked repository
// and a user injected into the context. If user is nil, the request is unauthenticated.
func newAuthenticatedRequest(method, url string, body io.Reader, user *types.User) *http.Request {
	req, err := http.NewRequest(method, url, body)
	Expect(err).NotTo(HaveOccurred())

	// Inject the mocked GlobalRepo into the request's context
	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)

	// If a user is provided, add UserID to context and mock UsersRepo for role checks
	if user != nil {
		ctxWithRepo = middleware.AddUserIDToContext(ctxWithRepo, user.ID)
		mockUsersRepo.EXPECT().Get(gomock.Any(), user.ID).Return(user, true, nil).AnyTimes()
	}

	return req.WithContext(ctxWithRepo)
}
