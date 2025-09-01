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
			Name:      utils.Ref("Updated Corp"),
			AddressID: utils.Ref(int64(3)),
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		existingComp = &types.Company{
			ID:        1,
			Name:      "Original Corp",
			AddressID: 2,
		}
	})

	Context("as an admin", func() {
		It("should update any company successfully", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingComp, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(&types.Address{ID: *payload.AddressID}, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.Name).To(Equal(*payload.Name))
			Expect(returnedCompany.AddressID).To(Equal(*payload.AddressID))
		})
	})

	Context("as a non-admin user", func() {
		It("should update their own company successfully", func() {
			// basicUser is in company 2. We'll try to update company 2.
			ownCompany := &types.Company{ID: basicUser.CompanyID, Name: "Basic User Co", AddressID: 4}
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), basicUser.CompanyID).Return(ownCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(&types.Address{ID: *payload.AddressID}, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			req := newAuthenticatedRequest("PUT", "/companies/2", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.Name).To(Equal(*payload.Name))
		})

		It("should forbid updating another company", func() {
			// basicUser is in company 2, trying to update company 1.
			otherCompany := &types.Company{ID: 1, Name: "Other Co"}
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(otherCompany, true, nil)

			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("user not authorized to update this company"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
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

	Context("when a dependency is not found", func() {
		It("should return 400 if the new address is not found", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingComp, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(nil, false, nil)

			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 on update failure", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingComp, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(&types.Address{ID: *payload.AddressID}, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db update failed"))

			req := newAuthenticatedRequest("PUT", "/companies/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to update company"))
		})
	})
})
