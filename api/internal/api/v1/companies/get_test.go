package companies_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Company Handler", func() {
	Context("when the company exists", func() {
		It("should return the company successfully", func() {
			companyID := int64(123)
			expectedCompany := &types.Company{ID: companyID, Name: "Found Corp", AddressID: 1}

			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(expectedCompany, true, nil)

			req := createRequestWithRepo("GET", "/api/v1/companies/123", nil, map[string]string{"id": "123"})
			companies.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.ID).To(Equal(companyID))
			Expect(returnedCompany.Name).To(Equal("Found Corp"))
		})
	})

	Context("when the company does not exist", func() {
		It("should return 404 Not Found", func() {
			companyID := int64(404)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(nil, false, nil)

			req := createRequestWithRepo("GET", "/api/v1/companies/404", nil, map[string]string{"id": "404"})
			companies.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("company not found"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			req := createRequestWithRepo("GET", "/api/v1/companies/abc", nil, map[string]string{"id": "abc"})
			companies.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid company ID"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			companyID := int64(500)
			dbErr := errors.New("database connection lost")

			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(nil, false, dbErr)

			req := createRequestWithRepo("GET", "/api/v1/companies/500", nil, map[string]string{"id": "500"})
			companies.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
