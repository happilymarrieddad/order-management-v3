package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Create User Endpoint", func() {
	var (
		rec     *httptest.ResponseRecorder
		payload users.CreateUserPayload
		address *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()

		// Setup standard objects
		address = &types.Address{ID: 1, Line1: "123 Test St"}

		payload = users.CreateUserPayload{
			Email:           "new.user@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "New",
			LastName:        "User",
			CompanyID:       company.ID,
			AddressID:       address.ID,
		}
	})

	Context("Happy Path", func() {
		It("should create a user successfully for a non-admin in their own company", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(company, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockUsersRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, u *types.User) error {
				u.ID = 3 // Simulate ID generation
				return nil
			})

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusCreated))
			var result types.User
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Email).To(Equal(payload.Email))
			Expect(result.Password).To(BeEmpty())
		})

		It("should create a user successfully for an admin in another company", func() {
			payload.CompanyID = 99 // Different company
			otherCompany := &types.Company{ID: 99, Name: "Admin-Created Company"}

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(otherCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockUsersRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusCreated))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), nil) // No user in context
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if a non-admin tries to create a user in another company", func() {
			payload.CompanyID = 99 // Different company
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with a malformed JSON body", func() {
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(`{"email":`)), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if passwords do not match", func() {
			payload.ConfirmPassword = "wrongpassword"
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if a required field is missing", func() {
			payload.FirstName = ""
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should fail if the user email already exists", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(&types.User{}, true, nil)
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if the company does not exist", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(nil, false, nil)
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if the address does not exist", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(company, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, nil)
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 on user creation db error", func() {
			dbErr := errors.New("db error")
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(company, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockUsersRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
