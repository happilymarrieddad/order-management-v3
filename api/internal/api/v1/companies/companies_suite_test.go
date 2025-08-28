package companies_test

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

func TestCompanies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Companies Suite")
}

var (
	mockCtrl          *gomock.Controller
	mockGlobalRepo    *mock_repos.MockGlobalRepo
	mockCompaniesRepo *mock_repos.MockCompaniesRepo
	rr                *httptest.ResponseRecorder
	ctx               context.Context
)

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockGlobalRepo = mock_repos.NewMockGlobalRepo(mockCtrl)
	mockCompaniesRepo = mock_repos.NewMockCompaniesRepo(mockCtrl)

	// Set up the mock chain: GlobalRepo -> Companies() -> MockCompaniesRepo
	mockGlobalRepo.EXPECT().Companies().Return(mockCompaniesRepo).AnyTimes()

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

	// Inject the mocked GlobalRepo into the request's context
	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
	req = req.WithContext(ctxWithRepo)

	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}

	return req
}
