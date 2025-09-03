package locations_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Create Location Endpoint", func() {
	var (
		rec     *httptest.ResponseRecorder
		payload locations.CreateLocationPayload
		address *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		address = &types.Address{ID: 1, Line1: "123 Test St"}
		payload = locations.CreateLocationPayload{
			Name:      "Test Location",
			CompanyID: company.ID, // Use company from suite
			AddressID: address.ID,
		}
	})

	performRequest := func(payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		req := newAuthenticatedRequest(http.MethodPost, "/locations", body, user)
		router.ServeHTTP(rec, req)
	}

	Context("Happy Path", func() {
		It("should create a location successfully for an admin", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(company, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

			performRequest(payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusCreated))
		})

		It("should create a location successfully for a normal user in their own company", func() {
			// Adjust payload to match normalUser's company
			payload.CompanyID = normalUser.CompanyID

			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(company, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

			performRequest(payload, normalUser)

			Expect(rec.Code).To(Equal(http.StatusCreated))
		})
	})

	Context("Error Paths", func() {
		It("should fail if not authenticated", func() {
			performRequest(payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		// Removed: "should fail if a non-admin tries to create a location"
		// This test is no longer relevant as normal users are now allowed to create locations.

		It("should fail with a malformed JSON body", func() {
			req := newAuthenticatedRequest(http.MethodPost, "/locations", []byte(`{`), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if the company does not exist", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(nil, false, nil)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if the address does not exist", func() {
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(company, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, nil)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(company, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

			performRequest(payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should fail if a normal user tries to create a location for another company", func() {
			// Set payload company ID to a different company than normalUser's
			payload.CompanyID = normalUser.CompanyID + 999

			// No mock expectations for CompaniesRepo.Get or AddressesRepo.Get here,
			// as the request should be forbidden before those checks.

			performRequest(payload, normalUser)

			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})
})
