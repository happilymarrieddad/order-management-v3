package companies_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Delete Company Endpoint", func() {
	var (
		rec *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
	})

	performRequest := func(companyID string, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodDelete, "/companies/"+companyID, url.Values{}, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should delete a company successfully for an admin", func() {
			mockCompaniesRepo.EXPECT().Delete(gomock.Any(), company.ID).Return(nil)

			performRequest("1", adminUser)

			Expect(rec.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest("1", nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			performRequest("1", normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid company ID", func() {
			performRequest("invalid-id", adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Repository Errors", func() {
		It("should return 404 if the company is not found", func() {
			mockCompaniesRepo.EXPECT().Delete(gomock.Any(), company.ID).Return(types.NewNotFoundError("company"))
			performRequest("1", adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Delete(gomock.Any(), company.ID).Return(dbErr)
			performRequest("1", adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
