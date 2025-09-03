package companies_test

import (
	"bytes"
	"context" // Added context import
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

var _ = Describe("Update Company Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		payload    companies.UpdateCompanyPayload
		targetCompany *types.Company
		newAddress *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetCompany = &types.Company{ID: 1, Name: "Old Company", AddressID: 10}
		newAddress = &types.Address{ID: 20, Line1: "New Address"}

		payload = companies.UpdateCompanyPayload{
			Name:      utils.Ref("Updated Company"),
			AddressID: utils.Ref(newAddress.ID),
		}
	})

	performRequest := func(companyID string, payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPut, "/companies/"+companyID, url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should update a company successfully for an admin", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(targetCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(newAddress, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, c *types.Company) error {
				Expect(c.Name).To(Equal(*payload.Name))
				Expect(c.AddressID).To(Equal(*payload.AddressID))
				return nil
			})

			performRequest("1", payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.Company
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Name).To(Equal(*payload.Name))
		})

		It("should update only the name if address_id is not provided", func() {
			payload.AddressID = nil
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(targetCompany, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, c *types.Company) error {
				Expect(c.Name).To(Equal(*payload.Name))
				Expect(c.AddressID).To(Equal(targetCompany.AddressID)) // Should remain unchanged
				return nil
			})

			performRequest("1", payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should update only the address_id if name is not provided", func() {
			payload.Name = nil
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(targetCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(newAddress, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, c *types.Company) error {
				Expect(c.Name).To(Equal(targetCompany.Name)) // Should remain unchanged
				Expect(c.AddressID).To(Equal(*payload.AddressID))
				return nil
			})

			performRequest("1", payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			// No mock expectation for CompaniesRepo.Get here, as the request should be stopped by middleware.
			performRequest("1", payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			// No mock expectation for CompaniesRepo.Get here, as the request should be stopped by middleware.
			performRequest("1", payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid company ID", func() {
			performRequest("invalid-id", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with a malformed JSON body", func() {
			rec, err := testutils.PerformRequest(router, http.MethodPut, "/companies/1", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if both name and address_id are missing", func() {
			payload.Name = nil
			payload.AddressID = nil
			// No mock expectation for CompaniesRepo.Get here, as the custom validation in handler will return early.
			performRequest("1", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should return 404 if the company to update is not found", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(nil, false, nil)
			performRequest("1", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on get company db error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(nil, false, dbErr)
			performRequest("1", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 400 if the new address does not exist", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(targetCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(nil, false, nil)
			performRequest("1", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 on get address db error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(targetCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(nil, false, dbErr)
			performRequest("1", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on update company db error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), targetCompany.ID).Return(targetCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(newAddress, true, nil)
			mockCompaniesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest("1", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})