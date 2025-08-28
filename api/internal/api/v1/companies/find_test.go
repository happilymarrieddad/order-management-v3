package companies_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Companies Handler", func() {
	Context("when companies exist", func() {
		It("should return a list of companies with default pagination", func() {
			foundCompanies := []*types.Company{
				{ID: 1, Name: "Company A"},
				{ID: 2, Name: "Company B"},
			}
			// The handler should apply default limit/offset when none are provided.
			// The test sends an empty JSON body, so the handler should use defaults.
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.CompanyFindOpts{Limit: 10, Offset: 0})).Return(foundCompanies, int64(2), nil)

			// The endpoint is a POST to /find with an empty body for defaults.
			req := createRequestWithRepo("POST", "/api/v1/companies/find", []byte(`{}`), nil)
			companies.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))

			// We need to remarshal and unmarshal the data to compare it properly
			dataBytes, _ := json.Marshal(result.Data)
			var returnedCompanies []types.Company
			json.Unmarshal(dataBytes, &returnedCompanies)
			Expect(returnedCompanies).To(HaveLen(2))
			Expect(returnedCompanies[0].Name).To(Equal("Company A"))
		})

		It("should return a list of companies with custom pagination", func() {
			foundCompanies := []*types.Company{{ID: 3, Name: "Company C"}}
			opts := &repos.CompanyFindOpts{Limit: 5, Offset: 5}
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(opts)).Return(foundCompanies, int64(1), nil)

			// Send the pagination options in the request body.
			body, _ := json.Marshal(opts)
			req := createRequestWithRepo("POST", "/api/v1/companies/find", body, nil)
			companies.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
		})
	})

	Context("when no companies exist", func() {
		It("should return an empty list", func() {
			// The handler should still apply default limits even for an empty result.
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.CompanyFindOpts{Limit: 10, Offset: 0})).Return([]*types.Company{}, int64(0), nil)

			req := createRequestWithRepo("POST", "/api/v1/companies/find", []byte(`{}`), nil)
			companies.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("find query failed")
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			req := createRequestWithRepo("POST", "/api/v1/companies/find", []byte(`{}`), nil)
			companies.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
