package locations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update Location Handler", func() {
	var (
		payload      locations.UpdateLocationPayload
		payloadBytes []byte
		existingLoc  *types.Location
		err          error
	)

	BeforeEach(func() {
		payload = locations.UpdateLocationPayload{
			Name:      utils.Ref("Updated Warehouse Name"),
			AddressID: utils.Ref(int64(3)),
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		existingLoc = &types.Location{ID: 1, Name: "Old Name", CompanyID: 1, AddressID: 2}
	})

	Context("as an admin", func() {
		It("should update any location successfully", func() {
			// adminUser is in company 1, updating a location in company 1
			mockLocationsRepo.EXPECT().Get(gomock.Any(), existingLoc.ID).Return(existingLoc, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(&types.Address{ID: *payload.AddressID}, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, loc *types.Location) error {
					Expect(loc.ID).To(Equal(existingLoc.ID))
					Expect(loc.Name).To(Equal(*payload.Name))
					Expect(loc.AddressID).To(Equal(*payload.AddressID))
					return nil
				},
			)

			req := newAuthenticatedRequest("PUT", "/locations/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			var returnedLocation types.Location
			err := json.Unmarshal(rr.Body.Bytes(), &returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.Name).To(Equal(*payload.Name))
		})
	})

	Context("as a non-admin user", func() {
		It("should update a location in their own company successfully", func() {
			// basicUser is in company 2. We'll try to update a location in company 2.
			ownLocation := &types.Location{ID: 5, Name: "Basic's Warehouse", CompanyID: basicUser.CompanyID, AddressID: 6}

			mockLocationsRepo.EXPECT().Get(gomock.Any(), ownLocation.ID).Return(ownLocation, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(&types.Address{ID: *payload.AddressID}, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			req := newAuthenticatedRequest("PUT", "/locations/5", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			var returnedLocation types.Location
			err := json.Unmarshal(rr.Body.Bytes(), &returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.Name).To(Equal(*payload.Name))
		})

		It("should forbid updating a location in another company", func() {
			// basicUser is in company 2, trying to update a location in company 1.
			otherLocation := &types.Location{ID: 1, Name: "Admin's Warehouse", CompanyID: 1, AddressID: 2}
			mockLocationsRepo.EXPECT().Get(gomock.Any(), otherLocation.ID).Return(otherLocation, true, nil)

			req := newAuthenticatedRequest("PUT", "/locations/1", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("user not authorized to update this location"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("PUT", "/locations/1", bytes.NewBuffer(payloadBytes), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("when a dependency is not found", func() {
		It("should return 404 if the location to update is not found", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(404)).Return(nil, false, nil)

			req := newAuthenticatedRequest("PUT", "/locations/404", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 400 if the new address is not found", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), existingLoc.ID).Return(existingLoc, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(nil, false, nil)

			req := newAuthenticatedRequest("PUT", "/locations/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("new address not found"))
		})
	})

	Context("with an invalid ID", func() {
		It("should return 404 for a non-integer ID", func() {
			req := newAuthenticatedRequest("PUT", "/locations/abc", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository fails on update", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("db update failed")

			mockLocationsRepo.EXPECT().Get(gomock.Any(), existingLoc.ID).Return(existingLoc, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(&types.Address{ID: *payload.AddressID}, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)

			req := newAuthenticatedRequest("PUT", "/locations/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to update location"))
		})
	})
})
