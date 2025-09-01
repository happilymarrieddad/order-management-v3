package companies_test

import (
	"bytes"
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
	var (
		payload      companies.CreateCompanyPayload
		payloadBytes []byte
		newCompany   *types.Company
		err          error
	)

	BeforeEach(func() {
		payload = companies.CreateCompanyPayload{
			Name:      "Test Corp",
			AddressID: 1,
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		newCompany = &types.Company{
			ID:        1,
			Name:      payload.Name,
			AddressID: payload.AddressID,
		}
	})

	Context("when creation is successful", func() {
		It("should return 201 Created with the new company", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(&types.Address{ID: payload.AddressID}, true, nil)
			mockCompaniesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Do(
				func(ctx context.Context, comp *types.Company) {
					comp.ID = newCompany.ID
				}).Return(nil)

			req := newAuthenticatedRequest("POST", "/companies", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedCompany types.Company
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCompany.ID).To(Equal(newCompany.ID))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a validation error (missing name)", func() {
			payload.Name = "" // Make the payload invalid
			body, _ := json.Marshal(payload)
			req := newAuthenticatedRequest("POST", "/companies", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'name' is required."))
		})
	})

	Context("when a dependency is not found", func() {
		It("should return 400 if the address does not exist", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, nil)

			req := newAuthenticatedRequest("POST", "/companies", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 Internal Server Error", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(&types.Address{ID: payload.AddressID}, true, nil)
			mockCompaniesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db insert failed"))
			req := newAuthenticatedRequest("POST", "/companies", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to create company"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("POST", "/companies", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("POST", "/companies", bytes.NewBuffer(payloadBytes), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
