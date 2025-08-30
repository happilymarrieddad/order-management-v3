package commodityattributes_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
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
	mockCtrl                  *gomock.Controller
	mockGlobalRepo            *mock_repos.MockGlobalRepo
	mockCommodityAttributesRepo *mock_repos.MockCommodityAttributesRepo
	mockUsersRepo             *mock_repos.MockUsersRepo
	rr                        *httptest.ResponseRecorder
	ctx                       context.Context
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
	ctx = context.Background()
})

var _ = AfterEach(func() {
	mockCtrl.Finish()
})

// createRequestWithRepo creates a new HTTP request with the mocked repository
// injected into the context, and optionally a user ID.
func createRequestWithRepo(method, url string, body []byte, vars map[string]string, userID ...int64) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	Expect(err).NotTo(HaveOccurred())

	// Inject the mocked GlobalRepo into the request\'s context
	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
	req = req.WithContext(ctxWithRepo)

	// Optionally add UserID to context and mock UsersRepo for admin check
	if len(userID) > 0 {
		req = req.WithContext(middleware.AddUserIDToContext(req.Context(), userID[0]))
		// Mock the UsersRepo to return an admin user for the given userID
		mockUsersRepo.EXPECT().Get(gomock.Any(), userID[0]).Return(&types.User{
			ID:    userID[0],
			Roles: types.Roles{types.RoleAdmin}, // Ensure the user has the admin role
		}, true, nil).AnyTimes()
	}

	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}

	return req
}
