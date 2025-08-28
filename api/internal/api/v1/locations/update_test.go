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

var _ = Describe("Update Location Handler", func() {
	var updatePayload map[string]interface{}

	BeforeEach(func() {
		updatePayload = map[string]interface{}{
			"name":       "Updated Warehouse Name",
			"address_id": int64(3),
		}
	})

	Context("with a valid request", func() {
		It("should update the location successfully", func() {
			locationID := int64(123)
			body, _ := json.Marshal(updatePayload)
			newAddressID := updatePayload["address_id"].(int64)

			existingLocation := &types.Location{ID: locationID, Name: "Old Name", CompanyID: 1, AddressID: 2}
			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(existingLocation, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), newAddressID).Return(&types.Address{ID: newAddressID}, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, loc *types.Location) error {
					Expect(loc.ID).To(Equal(locationID))
					Expect(loc.Name).To(Equal(updatePayload["name"]))
					Expect(loc.AddressID).To(Equal(newAddressID))
					return nil
				},
			)

			req := createRequestWithRepo("PUT", "/api/v1/locations/123", body, map[string]string{"id": "123"})
			locations.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			var returnedLocation types.Location
			err := json.Unmarshal(rr.Body.Bytes(), &returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.Name).To(Equal("Updated Warehouse Name"))
		})
	})

	Context("when a dependency is not found", func() {
		It("should return 404 if the location to update is not found", func() {
			locationID := int64(404)
			body, _ := json.Marshal(updatePayload)
			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(nil, false, nil)

			req := createRequestWithRepo("PUT", "/api/v1/locations/404", body, map[string]string{"id": "404"})
			locations.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 400 if the new address is not found", func() {
			locationID := int64(123)
			body, _ := json.Marshal(updatePayload)
			newAddressID := updatePayload["address_id"].(int64)

			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(&types.Location{ID: locationID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), newAddressID).Return(nil, false, nil)

			req := createRequestWithRepo("PUT", "/api/v1/locations/123", body, map[string]string{"id": "123"})
			locations.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("new address not found"))
		})
	})

	Context("when the repository fails", func() {
		It("should return 400 if the new name is a duplicate", func() {
			locationID := int64(123)
			body, _ := json.Marshal(updatePayload)
			newAddressID := updatePayload["address_id"].(int64)

			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(&types.Location{ID: locationID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), newAddressID).Return(&types.Address{ID: newAddressID}, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(repos.ErrLocationNameExists)

			req := createRequestWithRepo("PUT", "/api/v1/locations/123", body, map[string]string{"id": "123"})
			locations.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 for a generic database error", func() {
			locationID := int64(123)
			body, _ := json.Marshal(updatePayload)
			newAddressID := updatePayload["address_id"].(int64)
			dbErr := errors.New("db update failed")

			mockLocationsRepo.EXPECT().Get(gomock.Any(), locationID).Return(&types.Location{ID: locationID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), newAddressID).Return(&types.Address{ID: newAddressID}, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)

			req := createRequestWithRepo("PUT", "/api/v1/locations/123", body, map[string]string{"id": "123"})
			locations.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
