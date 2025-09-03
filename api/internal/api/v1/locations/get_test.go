package locations_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Get Location Endpoint", func() {
	var (
		rec      *httptest.ResponseRecorder
		location *types.Location
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		location = &types.Location{ID: 1, Name: "Test Location", CompanyID: company.ID, AddressID: 1}
	})

	performRequest := func(locationID int64, user *types.User) {
		req := newAuthenticatedRequest(http.MethodGet, "/locations/"+strconv.FormatInt(locationID, 10), nil, user)
		router.ServeHTTP(rec, req)
	}

	Context("Happy Path", func() {
		It("should get a location successfully for a non-admin in their own company", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, location.ID).Return(location, true, nil)

			performRequest(location.ID, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var returnedLocation types.Location
			err := json.NewDecoder(rec.Body).Decode(&returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.ID).To(Equal(location.ID))
		})

		It("should get a location successfully for an admin for any company", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(0), location.ID).Return(location, true, nil)

			performRequest(location.ID, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Error Paths", func() {
		It("should fail if not authenticated", func() {
			performRequest(location.ID, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 404 if a non-admin tries to get a location for another company", func() {
			otherLocation := &types.Location{ID: 2, Name: "Other Location", CompanyID: 99}
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, otherLocation.ID).Return(nil, false, nil)

			performRequest(otherLocation.ID, normalUser)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with an invalid ID", func() {
			req := newAuthenticatedRequest(http.MethodGet, "/locations/invalid-id", nil, normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound)) // Gorilla Mux returns 404 for non-matching routes
		})

		It("should return 404 if the location is not found", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, location.ID).Return(nil, false, nil)

			performRequest(location.ID, normalUser)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, location.ID).Return(nil, false, dbErr)

			performRequest(location.ID, normalUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
