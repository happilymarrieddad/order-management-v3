package locations_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

var _ = Describe("Update Location Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		payload    locations.UpdateLocationPayload
		address    *types.Address
		newAddress *types.Address
		location   *types.Location
		locationID int64
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		address = &types.Address{ID: 1, Line1: "123 Test St"}
		newAddress = &types.Address{ID: 2, Line1: "456 New St"}
		locationID = 1
		location = &types.Location{ID: locationID, Name: "Old Name", CompanyID: company.ID, AddressID: address.ID}
		payload = locations.UpdateLocationPayload{
			Name:      utils.Ref("New Name"),
			AddressID: utils.Ref(newAddress.ID),
		}
	})

	performRequest := func(locID int64, payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		req := newAuthenticatedRequest(http.MethodPut, "/locations/"+strconv.FormatInt(locID, 10), body, user)
		router.ServeHTTP(rec, req)
	}

	Context("Happy Path", func() {
		It("should update a location successfully for a non-admin in their own company", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(location, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(newAddress, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, loc *types.Location) error {
				Expect(loc.Name).To(Equal(*payload.Name))
				Expect(loc.AddressID).To(Equal(*payload.AddressID))
				return nil
			})

			performRequest(locationID, payload, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should update a location successfully for an admin", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(0), locationID).Return(location, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(newAddress, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			performRequest(locationID, payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should update only the name", func() {
			payload.AddressID = nil
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(location, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, loc *types.Location) error {
				Expect(loc.Name).To(Equal(*payload.Name))
				Expect(loc.AddressID).To(Equal(address.ID)) // Should not change
				return nil
			})

			performRequest(locationID, payload, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should update only the address", func() {
			payload.Name = nil
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(location, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(newAddress, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, loc *types.Location) error {
				Expect(loc.Name).To(Equal("Old Name")) // Should not change
				Expect(loc.AddressID).To(Equal(*payload.AddressID))
				return nil
			})

			performRequest(locationID, payload, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Error Paths", func() {
		It("should fail if not authenticated", func() {
			performRequest(locationID, payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if a non-admin tries to update a location for another company", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(nil, false, nil)
			performRequest(locationID, payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with an invalid ID", func() {
			req := newAuthenticatedRequest(http.MethodPut, "/locations/invalid-id", nil, normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail if the location is not found", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(nil, false, nil)
			performRequest(locationID, payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail if the new address is not found", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(location, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(nil, false, nil)
			performRequest(locationID, payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 on a database error during get", func() {
			dbErr := errors.New("db error")
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(nil, false, dbErr)
			performRequest(locationID, payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on a database error during update", func() {
			dbErr := errors.New("db error")
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(location, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), *payload.AddressID).Return(newAddress, true, nil)
			mockLocationsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest(locationID, payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should fail with a malformed JSON body", func() {
			req := newAuthenticatedRequest(http.MethodPut, "/locations/"+strconv.FormatInt(locationID, 10), []byte(`{`), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
