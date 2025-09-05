package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("PUT /users/{id}/company", func() {
	var (
		pld users.UpdateUserCompanyPayload
	)

	BeforeEach(func() {
		pld = users.UpdateUserCompanyPayload{CompanyID: company.ID}
	})

	Context("when an admin updates a user in their own company", func() {
		It("should return 204 No Content", func() {
			targetUser := &types.User{ID: normalUser.ID, CompanyID: adminUser.CompanyID}
			newCompany := &types.Company{ID: pld.CompanyID}

			mockUsersRepo.EXPECT().Get(gomock.Any(), int64(0), targetUser.ID).Return(targetUser, true, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(newCompany, true, nil)
			mockUsersRepo.EXPECT().UpdateUserCompany(gomock.Any(), targetUser.ID, newCompany.ID).Return(nil)

			body, _ := json.Marshal(pld)
			req := newAuthenticatedRequest("PUT", "/users/1/company", bytes.NewReader(body), adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("when an admin updates a user in another company", func() {
		It("should return 204 No Content", func() {
			otherCompanyUser := &types.User{ID: 3, CompanyID: 99}
			newCompany := &types.Company{ID: pld.CompanyID}

			mockUsersRepo.EXPECT().Get(gomock.Any(), int64(0), otherCompanyUser.ID).Return(otherCompanyUser, true, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(newCompany, true, nil)
			mockUsersRepo.EXPECT().UpdateUserCompany(gomock.Any(), otherCompanyUser.ID, newCompany.ID).Return(nil)

			body, _ := json.Marshal(pld)
			req := newAuthenticatedRequest("PUT", "/users/3/company", bytes.NewReader(body), adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("when the target user does not exist", func() {
		It("should return 404 Not Found", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), int64(0), gomock.Any()).Return(nil, false, nil)

			body, _ := json.Marshal(pld)
			req := newAuthenticatedRequest("PUT", "/users/999/company", bytes.NewReader(body), adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the new company does not exist", func() {
		It("should return 400 Bad Request", func() {
			pld.CompanyID = 999
			targetUser := &types.User{ID: normalUser.ID, CompanyID: adminUser.CompanyID}

			mockUsersRepo.EXPECT().Get(gomock.Any(), int64(0), targetUser.ID).Return(targetUser, true, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), pld.CompanyID).Return(nil, false, nil)

			body, _ := json.Marshal(pld)
			req := newAuthenticatedRequest("PUT", "/users/1/company", bytes.NewReader(body), adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the request body is invalid", func() {
		It("should return 400 Bad Request", func() {
			req := newAuthenticatedRequest("PUT", "/users/1/company", bytes.NewReader([]byte(`{`)), adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the database update fails", func() {
		It("should return 500 Internal Server Error", func() {
			dbErr := errors.New("db update error")
			targetUser := &types.User{ID: normalUser.ID, CompanyID: adminUser.CompanyID}
			newCompany := &types.Company{ID: pld.CompanyID}

			mockUsersRepo.EXPECT().Get(gomock.Any(), int64(0), targetUser.ID).Return(targetUser, true, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), newCompany.ID).Return(newCompany, true, nil)
			mockUsersRepo.EXPECT().UpdateUserCompany(gomock.Any(), targetUser.ID, newCompany.ID).Return(dbErr)

			body, _ := json.Marshal(pld)
			req := newAuthenticatedRequest("PUT", "/users/1/company", bytes.NewReader(body), adminUser)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
