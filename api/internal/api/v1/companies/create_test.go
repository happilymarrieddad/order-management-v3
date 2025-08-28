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

var _ = Describe("Create Company Handler", func() {
	var createPayload map[string]interface{}

	BeforeEach(func() {
		createPayload = map[string]interface{}{
			"name":       "Test Company",
			"address_id": int64(1),
		}
	})

	Context("with a valid request", func() {
		It("should create a company successfully", func() {
			body, _ := json.Marshal(createPayload)

			mockCompaniesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, company *types.Company) error {
					Expect(company.Name).To(Equal(createPayload["name"]))
					Expect(company.AddressID).To(Equal(createPayload["address_id"]))
					company.ID = 123 // Simulate DB assigning an ID
					return nil
				},
			)

			req := createRequestWithRepo("POST", "/api/v1/companies", body, nil)
			companies.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.ID).To(Equal(int64(123)))
			Expect(returnedCompany.Name).To(Equal("Test Company"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a malformed JSON body", func() {
			body := []byte(`{"name": "bad json",`)
			req := createRequestWithRepo("POST", "/api/v1/companies", body, nil)
			companies.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a missing required field (name)", func() {
			delete(createPayload, "name")
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/companies", body, nil)
			companies.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a missing required field (address_id)", func() {
			delete(createPayload, "address_id")
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/companies", body, nil)
			companies.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			body, _ := json.Marshal(createPayload)
			dbErr := errors.New("unexpected database error")

			mockCompaniesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

			req := createRequestWithRepo("POST", "/api/v1/companies", body, nil)
			companies.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
