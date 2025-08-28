package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	jwtpkg "github.com/happilymarrieddad/order-management-v3/api/internal/jwt"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

var _ = Describe("AuthMiddleware", func() {
	var (
		rr          *httptest.ResponseRecorder
		nextHandler http.Handler
		wasCalled   bool
	)

	BeforeEach(func() {
		rr = httptest.NewRecorder()
		wasCalled = false

		// This dummy handler will be called only if authentication is successful.
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wasCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Set the JWT secret for the test environment.
		os.Setenv("JWT_SECRET", "test-secret-for-middleware")
	})

	AfterEach(func() {
		os.Unsetenv("JWT_SECRET")
	})

	Context("with a valid token", func() {
		It("should call the next handler and add userID to the context", func() {
			user := &types.User{ID: 123}
			token, err := jwtpkg.GenerateToken(user)
			Expect(err).NotTo(HaveOccurred())

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-App-Token", token)

			// This test handler also checks the context value for correctness.
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wasCalled = true
				ctxUserID := r.Context().Value(middleware.UserIDKey)
				Expect(ctxUserID).NotTo(BeNil())
				Expect(ctxUserID).To(Equal(user.ID))
				w.WriteHeader(http.StatusOK)
			})

			middleware.AuthMiddleware(testHandler).ServeHTTP(rr, req)

			Expect(wasCalled).To(BeTrue())
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})

	Context("with an invalid or missing token", func() {
		It("should return 401 Unauthorized for a missing token", func() {
			req := httptest.NewRequest("GET", "/", nil) // No token header
			middleware.AuthMiddleware(nextHandler).ServeHTTP(rr, req)
			Expect(wasCalled).To(BeFalse())
			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 401 Unauthorized for a malformed token", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-App-Token", "this.is.not.a.valid.token")
			middleware.AuthMiddleware(nextHandler).ServeHTTP(rr, req)
			Expect(wasCalled).To(BeFalse())
			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 401 Unauthorized for a token signed with the wrong secret", func() {
			// Manually create a token with a different secret to avoid env var caching issues.
			claims := jwtv5.MapClaims{
				"sub": "123",
				"exp": time.Now().Add(time.Hour * 24).Unix(),
				"iat": time.Now().Unix(),
			}
			tokenObj := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
			token, err := tokenObj.SignedString([]byte("a-different-secret"))
			Expect(err).NotTo(HaveOccurred())

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-App-Token", token)
			middleware.AuthMiddleware(nextHandler).ServeHTTP(rr, req)
			Expect(wasCalled).To(BeFalse())
			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 401 Unauthorized for an expired token", func() {
			// Manually create a token with an expiration time in the past.
			claims := jwtv5.MapClaims{
				"sub": "789",
				"exp": time.Now().Add(-1 * time.Hour).Unix(), // Expired one hour ago
				"iat": time.Now().Add(-2 * time.Hour).Unix(),
			}
			tokenObj := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)

			// Sign with the correct secret that the middleware will use for verification.
			secret := os.Getenv("JWT_SECRET")
			Expect(secret).NotTo(BeEmpty())
			token, err := tokenObj.SignedString([]byte(secret))
			Expect(err).NotTo(HaveOccurred())

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-App-Token", token)
			middleware.AuthMiddleware(nextHandler).ServeHTTP(rr, req)
			Expect(wasCalled).To(BeFalse())
			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
