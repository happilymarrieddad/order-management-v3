package companies_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Company Handler", func() {
	Context("when the company exists", func() {
		It("should return the company for an authenticated user", func() {
			company := &types.Company{ID: 1, Name: "Test Co"}
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(company, true, nil)

			req := newAuthenticatedRequest("GET", "/companies/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.ID).To(Equal(company.ID))
		})
	})

	Context("when the company does not exist", func() {
		It("should return 404 Not Found", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			req := newAuthenticatedRequest("GET", "/companies/999", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the ID is invalid", func() {
		It("should return 404 Not Found from the router", func() {
			req := newAuthenticatedRequest("GET", "/companies/abc", nil, basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("db went boom")
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, false, dbErr)

			req := newAuthenticatedRequest("GET", "/companies/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to get company"))
		})
	})
})
