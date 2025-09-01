package companies_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Companies Handler", func() {
	Context("when companies exist", func() {
		It("should return a list of companies for an admin user", func() {
			foundCompanies := []*types.Company{
				{ID: 1, Name: "Alpha Co"},
				{ID: 2, Name: "Beta Co"},
			}
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.CompanyFindOpts{Limit: 10, Offset: 0})).Return(foundCompanies, int64(2), nil)

			req := newAuthenticatedRequest("POST", "/companies/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))

			dataBytes, _ := json.Marshal(result.Data)
			var returnedCompanies []types.Company
			json.Unmarshal(dataBytes, &returnedCompanies)
			Expect(returnedCompanies).To(HaveLen(2))
		})
	})

	Context("when no companies exist", func() {
		It("should return an empty list for an admin user", func() {
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return([]*types.Company{}, int64(0), nil)

			req := newAuthenticatedRequest("POST", "/companies/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 Internal Server Error", func() {
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), errors.New("find query failed"))
			req := newAuthenticatedRequest("POST", "/companies/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to find companies"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("POST", "/companies/find", bytes.NewBufferString(`{}`), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("POST", "/companies/find", bytes.NewBufferString(`{}`), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
