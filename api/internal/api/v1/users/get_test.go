package users_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Get User Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		targetUser *types.User
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetUser = &types.User{ID: 2, CompanyID: company.ID, FirstName: "Target", LastName: "User"}
	})

	Context("Happy Path", func() {
		It("should get a user successfully for a normal user in their own company", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)

			req := newAuthenticatedRequest(http.MethodGet, "/users/2", nil, normalUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.User
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.ID).To(Equal(targetUser.ID))
			Expect(result.Password).To(BeEmpty())
		})

		It("should get a user successfully for an admin", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)

			req := newAuthenticatedRequest(http.MethodGet, "/users/2", nil, adminUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Error Paths", func() {
		It("should fail if not authenticated", func() {
			req := newAuthenticatedRequest(http.MethodGet, "/users/2", nil, nil)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if a normal user tries to get a user from another company", func() {
			otherCompany := &types.Company{ID: 99, Name: "Other Company"}
			otherUser := &types.User{ID: 3, CompanyID: otherCompany.ID, FirstName: "Other"}
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, otherUser.ID).Return(nil, false, nil)
			req := newAuthenticatedRequest(http.MethodGet, "/users/3", nil, normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with an invalid user ID", func() {
			req := newAuthenticatedRequest(http.MethodGet, "/users/invalid-id", nil, normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 404 if the user is not found", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(nil, false, nil)
			req := newAuthenticatedRequest(http.MethodGet, "/users/2", nil, normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(nil, false, dbErr)

			req := newAuthenticatedRequest(http.MethodGet, "/users/2", nil, normalUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
