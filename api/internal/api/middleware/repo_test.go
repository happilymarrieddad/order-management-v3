package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("RepoMiddleware", func() {
	var (
		ctrl           *gomock.Controller
		mockGlobalRepo *mock_repos.MockGlobalRepo
		rr             *httptest.ResponseRecorder
		nextHandler    http.Handler
		wasCalled      bool
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockGlobalRepo = mock_repos.NewMockGlobalRepo(ctrl)
		rr = httptest.NewRecorder()
		wasCalled = false

		// This dummy handler will be called if the middleware passes the request on.
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wasCalled = true
			// Test that the repo can be retrieved from the context and is the correct instance.
			repoFromCtx := middleware.GetRepo(r.Context())
			Expect(repoFromCtx).To(BeIdenticalTo(mockGlobalRepo))
			w.WriteHeader(http.StatusOK)
		})
	})

	It("should inject the repo into the context and call the next handler", func() {
		repoMiddleware := middleware.RepoMiddleware(mockGlobalRepo)
		req := httptest.NewRequest("GET", "/", nil)
		repoMiddleware(nextHandler).ServeHTTP(rr, req)
		Expect(wasCalled).To(BeTrue(), "The next handler should have been called")
		Expect(rr.Code).To(Equal(http.StatusOK))
	})
})
