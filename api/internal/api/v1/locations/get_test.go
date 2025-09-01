package locations_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Location Handler", func() {
	Context("when the location exists", func() {
		It("should return the location for an authenticated user", func() {
			location := &types.Location{ID: 1, Name: "Test Location"}
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(location, true, nil)

			req := newAuthenticatedRequest("GET", "/locations/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedLocation types.Location
			err := json.Unmarshal(rr.Body.Bytes(), &returnedLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation.ID).To(Equal(location.ID))
		})
	})

	Context("when the location does not exist", func() {
		It("should return 404 Not Found", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			req := newAuthenticatedRequest("GET", "/locations/999", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the ID is invalid", func() {
		It("should return 404 Not Found from the router", func() {
			req := newAuthenticatedRequest("GET", "/locations/abc", nil, basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("db went boom")
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, false, dbErr)

			req := newAuthenticatedRequest("GET", "/locations/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to get location"))
		})
	})
})
