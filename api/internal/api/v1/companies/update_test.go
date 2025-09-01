package companies_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update Company Handler", func() {
	var (
		payload      companies.UpdateCompanyPayload
		payloadBytes []byte
		existingComp *types.Company
		err          error
	)

	BeforeEach(func() {
		payload = companies.UpdateCompanyPayload{
			Name: utils.Ref("Updated Corp"),
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		existingComp = &types.Company{
			ID:   1,
			Name: "Original Corp",
		}
	})

	Context("when update is successful", func() {
		It("should return 200 OK with the updated company", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingComp, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.Name).To(Equal(*payload.Name))
		})
	})

	Context("when the company to update is not found", func() {
		It("should return 404 Not Found", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			req := newAuthenticatedRequest("PUT", "/companies/999", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a malformed JSON body", func() {
			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBufferString(`{]`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 404 for a non-integer ID", func() {
			req := newAuthenticatedRequest("PUT", "/companies/abc", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 on update failure", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingComp, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db update failed"))

			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to update company"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
