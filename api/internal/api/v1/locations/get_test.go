package locations_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Location Handler", func() {
	Context("when the location exists", func() {
		It("should return the location with its company and address", func() {
			locationID := int64(123)
			companyID := int64(1)
			addressID := int64(2)

			expectedLocation := &types.Location{ID: locationID, Name: "Found Location", CompanyID: companyID, AddressID: addressID}
			expectedCompany := &types.Company{ID: companyID, Name: "Test Co"}
			expectedAddress := &types.Address{ID: addressID, Line1: "123 Test St"}

			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(expectedLocation, true, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(expectedCompany, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(expectedAddress, true, nil)

			req := createRequestWithRepo("GET", "/api/v1/locations/123", nil, map[string]string{"id": "123"})
			locations.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedLocation types.Location
			err := json.Unmarshal(rr.Body.Bytes(), &returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.ID).To(Equal(locationID))
			Expect(returnedLocation.Name).To(Equal("Found Location"))
			Expect(returnedLocation.Company).NotTo(BeNil())
			Expect(returnedLocation.Company.Name).To(Equal("Test Co"))
			Expect(returnedLocation.Address).NotTo(BeNil())
			Expect(returnedLocation.Address.Line1).To(Equal("123 Test St"))
		})
	})

	Context("when the location does not exist", func() {
		It("should return 404 Not Found", func() {
			locationID := int64(404)
			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(nil, false, nil)

			req := createRequestWithRepo("GET", "/api/v1/locations/404", nil, map[string]string{"id": "404"})
			locations.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			locationID := int64(500)
			dbErr := errors.New("db connection lost")
			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(nil, false, dbErr)

			req := createRequestWithRepo("GET", "/api/v1/locations/500", nil, map[string]string{"id": "500"})
			locations.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
