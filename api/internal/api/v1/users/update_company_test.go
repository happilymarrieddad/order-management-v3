package users_test

import (
	"bytes"
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

var _ = Describe("Update User Company Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		payload    users.UpdateUserCompanyPayload
		targetUser *types.User
		newCompany *types.Company
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetUser = &types.User{ID: 1, CompanyID: company.ID, FirstName: "Test", LastName: "User"}
		newCompany = &types.Company{ID: 99, Name: "New Company"}

		payload = users.UpdateUserCompanyPayload{
			CompanyID: newCompany.ID,
		}
	})

	Context("Happy Path", func() {
		It("should update a user's company successfully for an admin", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(newCompany, true, nil)
			mockUsersRepo.EXPECT().GetIncludeInvisible(gomock.Any(), targetUser.ID).Return(targetUser, true, nil)
			mockUsersRepo.EXPECT().UpdateUserCompany(gomock.Any(), targetUser.ID, newCompany.ID).Return(nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), nil)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid user ID", func() {
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/invalid-id/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with a malformed JSON body", func() {
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer([]byte(`{\"company_id\":`)), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if a required field is missing", func() {
			payload.CompanyID = 0 // Invalid CompanyID
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should return 404 if the new company does not exist", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(nil, false, nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on get company db error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(nil, false, dbErr)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 404 if the user to update is not found", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(newCompany, true, nil)
			mockUsersRepo.EXPECT().GetIncludeInvisible(gomock.Any(), targetUser.ID).Return(nil, false, nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on get user db error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(newCompany, true, nil)
			mockUsersRepo.EXPECT().GetIncludeInvisible(gomock.Any(), targetUser.ID).Return(nil, false, dbErr)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on update user company db error", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(newCompany, true, nil)
			mockUsersRepo.EXPECT().GetIncludeInvisible(gomock.Any(), targetUser.ID).Return(targetUser, true, nil)
			dbErr := errors.New("db error")
			mockUsersRepo.EXPECT().UpdateUserCompany(gomock.Any(), targetUser.ID, newCompany.ID).Return(dbErr)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1/company", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
