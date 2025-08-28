package companies_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update Company Handler", func() {
	var updatePayload map[string]interface{}

	BeforeEach(func() {
		updatePayload = map[string]interface{}{
			"name":       "Updated Company Name",
			"address_id": int64(2),
		}
	})

	Context("with a valid request", func() {
		It("should update the company successfully", func() {
			companyID := int64(123)
			body, _ := json.Marshal(updatePayload)

			// Mock the Get call to find the existing company
			existingCompany := &types.Company{ID: companyID, Name: "Old Name", AddressID: 1}
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(existingCompany, true, nil)

			// Mock the Update call
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, company *types.Company) error {
					Expect(company.ID).To(Equal(companyID))
					Expect(company.Name).To(Equal(updatePayload["name"]))
					Expect(company.AddressID).To(Equal(updatePayload["address_id"]))
					return nil
				},
			)

			req := createRequestWithRepo("PUT", "/api/v1/companies/123", body, map[string]string{"id": "123"})
			companies.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.Name).To(Equal("Updated Company Name"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/companies/abc", body, map[string]string{"id": "abc"})
			companies.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a missing required field", func() {
			delete(updatePayload, "name")
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/companies/123", body, map[string]string{"id": "123"})
			companies.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the target company does not exist", func() {
		It("should return 404 Not Found", func() {
			companyID := int64(404)
			body, _ := json.Marshal(updatePayload)

			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(nil, false, nil)

			req := createRequestWithRepo("PUT", "/api/v1/companies/404", body, map[string]string{"id": "404"})
			companies.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 on Get error", func() {
			companyID := int64(500)
			body, _ := json.Marshal(updatePayload)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(nil, false, errors.New("get error"))

			req := createRequestWithRepo("PUT", "/api/v1/companies/500", body, map[string]string{"id": "500"})
			companies.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on Update error", func() {
			companyID := int64(123)
			body, _ := json.Marshal(updatePayload)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(&types.Company{ID: companyID}, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update error"))

			req := createRequestWithRepo("PUT", "/api/v1/companies/123", body, map[string]string{"id": "123"})
			companies.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
