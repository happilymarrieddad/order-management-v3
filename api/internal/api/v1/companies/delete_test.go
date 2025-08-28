package companies_test

import (
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Delete Company Handler", func() {
	Context("with a valid request", func() {
		It("should delete the company successfully", func() {
			companyID := int64(123)
			mockCompaniesRepo.EXPECT().Delete(gomock.Any(), companyID).Return(nil)

			req := createRequestWithRepo("DELETE", "/api/v1/companies/123", nil, map[string]string{"id": "123"})
			companies.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			req := createRequestWithRepo("DELETE", "/api/v1/companies/abc", nil, map[string]string{"id": "abc"})
			companies.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid company ID"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			companyID := int64(500)
			dbErr := errors.New("foreign key constraint fails")

			mockCompaniesRepo.EXPECT().Delete(gomock.Any(), companyID).Return(dbErr)

			req := createRequestWithRepo("DELETE", "/api/v1/companies/500", nil, map[string]string{"id": "500"})
			companies.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
