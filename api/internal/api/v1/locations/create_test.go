package locations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Create Location Handler", func() {
	var (
		payload      locations.CreateLocationPayload
		payloadBytes []byte
		newLocation  *types.Location
		err          error
	)

	BeforeEach(func() {
		payload = locations.CreateLocationPayload{
			Name:      "Main Warehouse",
			CompanyID: 1,
			AddressID: 2,
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		newLocation = &types.Location{
			ID:        1,
			Name:      payload.Name,
			CompanyID: payload.CompanyID,
			AddressID: payload.AddressID,
		}
	})

	Context("when creation is successful", func() {
		It("should return 201 Created with the new location", func() {
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Do(
				func(ctx context.Context, loc *types.Location) {
					loc.ID = newLocation.ID
				}).Return(nil)

			req := newAuthenticatedRequest("POST", "/locations", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedLocation types.Location
			err := json.Unmarshal(rr.Body.Bytes(), &returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.ID).To(Equal(newLocation.ID))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a validation error (missing name)", func() {
			payload.Name = "" // Make the payload invalid
			body, _ := json.Marshal(payload)
			req := newAuthenticatedRequest("POST", "/locations", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'name' is required."))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 Internal Server Error", func() {
			mockLocationsRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db insert failed"))
			req := newAuthenticatedRequest("POST", "/locations", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to create location"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("POST", "/locations", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("POST", "/locations", bytes.NewBuffer(payloadBytes), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
