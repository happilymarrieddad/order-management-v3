package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/auth"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("Login Handler", func() {
	var (
		ctrl           *gomock.Controller
		mockGlobalRepo *mock_repos.MockGlobalRepo
		mockUsersRepo  *mock_repos.MockUsersRepo
		rr             *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockGlobalRepo = mock_repos.NewMockGlobalRepo(ctrl)
		mockUsersRepo = mock_repos.NewMockUsersRepo(ctrl)
		rr = httptest.NewRecorder()

		// Set up the mock chain: globalRepo.Users() -> mockUsersRepo
		mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	// Helper function to create a request with the repo in its context
	createRequestWithRepo := func(body []byte) *http.Request {
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		// This key must match the one used in the middleware to inject the repo
		ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
		return req.WithContext(ctxWithRepo)
	}

	Context("when the request is valid", func() {
		It("should return a JWT token on successful authentication", func() {
			password := "password123"
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			Expect(err).NotTo(HaveOccurred())

			user := &types.User{
				ID:       1,
				Email:    "test@example.com",
				Password: string(hashedPassword),
			}

			creds := map[string]string{"email": "test@example.com", "password": password}
			body, _ := json.Marshal(creds)

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(user, true, nil)

			req := createRequestWithRepo(body)
			auth.Login(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var response map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(HaveKey("token"))
			Expect(response["token"]).NotTo(BeEmpty())
		})
	})

	Context("when the user is not found", func() {
		It("should return 401 Unauthorized", func() {
			creds := map[string]string{"email": "notfound@example.com", "password": "password123"}
			body, _ := json.Marshal(creds)

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), "notfound@example.com").Return(nil, false, nil)

			req := createRequestWithRepo(body)
			auth.Login(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("when the password is incorrect", func() {
		It("should return 401 Unauthorized", func() {
			password := "password123"
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			Expect(err).NotTo(HaveOccurred())

			user := &types.User{ID: 1, Email: "test@example.com", Password: string(hashedPassword)}
			creds := map[string]string{"email": "test@example.com", "password": "wrongpassword"}
			body, _ := json.Marshal(creds)

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(user, true, nil)

			req := createRequestWithRepo(body)
			auth.Login(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("when the database returns an error", func() {
		It("should return 500 Internal Server Error", func() {
			creds := map[string]string{"email": "test@example.com", "password": "password123"}
			body, _ := json.Marshal(creds)

			dbErr := errors.New("database connection lost")
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(nil, false, dbErr)

			req := createRequestWithRepo(body)
			auth.Login(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when the request body is invalid", func() {
		It("should return 400 Bad Request for malformed JSON", func() {
			body := []byte(`{"email": "test@example.com",`) // Malformed JSON
			req := createRequestWithRepo(body)
			auth.Login(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
