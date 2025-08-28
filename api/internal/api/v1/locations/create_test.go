package locations_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Create Location Handler", func() {
	var createPayload map[string]interface{}

	BeforeEach(func() {
		createPayload = map[string]interface{}{
			"company_id": int64(1),
			"address_id": int64(2),
			"name":       "Main Warehouse",
		}
	})

	Context("with a valid request", func() {
		It("should create a location successfully", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			addressID := createPayload["address_id"].(int64)

			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(&types.Company{ID: companyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(&types.Address{ID: addressID}, true, nil)
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, loc *types.Location) error {
					Expect(loc.CompanyID).To(Equal(companyID))
					Expect(loc.AddressID).To(Equal(addressID))
					Expect(loc.Name).To(Equal("Main Warehouse"))
					loc.ID = 123 // Simulate DB assigning an ID
					return nil
				},
			)

			req := createRequestWithRepo("POST", "/api/v1/locations", body, nil)
			locations.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedLocation types.Location
			err := json.Unmarshal(rr.Body.Bytes(), &returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.ID).To(Equal(int64(123)))
			Expect(returnedLocation.Name).To(Equal("Main Warehouse"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a missing required field", func() {
			delete(createPayload, "name")
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/locations", body, nil)
			locations.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when a dependency is not found", func() {
		It("should return 400 if the company does not exist", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(nil, false, nil)

			req := createRequestWithRepo("POST", "/api/v1/locations", body, nil)
			locations.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("company not found"))
		})

		It("should return 400 if the address does not exist", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			addressID := createPayload["address_id"].(int64)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(&types.Company{ID: companyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, nil)

			req := createRequestWithRepo("POST", "/api/v1/locations", body, nil)
			locations.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})
	})

	Context("when the repository fails", func() {
		It("should return 400 if the location name already exists for the company", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			addressID := createPayload["address_id"].(int64)

			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(&types.Company{ID: companyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(&types.Address{ID: addressID}, true, nil)
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repos.ErrLocationNameExists)

			req := createRequestWithRepo("POST", "/api/v1/locations", body, nil)
			locations.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(repos.ErrLocationNameExists.Error()))
		})

		It("should return 500 for a generic database error", func() {
			body, _ := json.Marshal(createPayload)
			dbErr := errors.New("unexpected db error")

			mockCompaniesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&types.Company{ID: 1}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&types.Address{ID: 2}, true, nil)
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

			req := createRequestWithRepo("POST", "/api/v1/locations", body, nil)
			locations.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
