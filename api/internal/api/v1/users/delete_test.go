package users_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("DELETE /users/{id}", func() {
	Context("when a normal user deletes themselves", func() {
		It("should return 204 No Content", func() {
			targetUser := &types.User{ID: normalUser.ID, CompanyID: normalUser.CompanyID}

			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, normalUser.ID).Return(targetUser, true, nil)
			mockUsersRepo.EXPECT().Delete(gomock.Any(), normalUser.ID).Return(nil)

			req := newAuthenticatedRequest("DELETE", "/users/1", nil, normalUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("when an admin deletes another user in the same company", func() {
		It("should return 204 No Content", func() {
			targetUser := &types.User{ID: normalUser.ID, CompanyID: adminUser.CompanyID}

			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.CompanyID, normalUser.ID).Return(targetUser, true, nil)
			mockUsersRepo.EXPECT().Delete(gomock.Any(), normalUser.ID).Return(nil)

			req := newAuthenticatedRequest("DELETE", "/users/1", nil, adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("when a normal user tries to delete another user", func() {
		It("should return 403 Forbidden", func() {
			targetUser := &types.User{ID: adminUser.ID, CompanyID: normalUser.CompanyID} // A different user

			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, adminUser.ID).Return(targetUser, true, nil)

			req := newAuthenticatedRequest("DELETE", "/users/2", nil, normalUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("when an admin tries to delete a user in another company", func() {
		It("should return 403 Forbidden", func() {
			otherCompany := &types.Company{ID: 99, Name: "Other Company"}
			targetUser := &types.User{ID: 3, CompanyID: otherCompany.ID}

			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)

			req := newAuthenticatedRequest("DELETE", "/users/3", nil, adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("when the user to be deleted does not exist", func() {
		It("should return 404 Not Found", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.CompanyID, gomock.Any()).Return(nil, false, nil)

			req := newAuthenticatedRequest("DELETE", "/users/999", nil, adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the database delete fails", func() {
		It("should return 500 Internal Server Error", func() {
			targetUser := &types.User{ID: normalUser.ID, CompanyID: adminUser.CompanyID}
			dbErr := errors.New("db delete error")

			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.CompanyID, normalUser.ID).Return(targetUser, true, nil)
			mockUsersRepo.EXPECT().Delete(gomock.Any(), normalUser.ID).Return(dbErr)

			req := newAuthenticatedRequest("DELETE", "/users/1", nil, adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})