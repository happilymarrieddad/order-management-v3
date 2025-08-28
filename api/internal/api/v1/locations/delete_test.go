package locations_test

import (
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Delete Location Handler", func() {
	Context("with a valid request", func() {
		It("should delete the location successfully", func() {
			locationID := int64(123)
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(nil)

			req := createRequestWithRepo("DELETE", "/api/v1/locations/123", nil, map[string]string{"id": "123"})
			locations.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			req := createRequestWithRepo("DELETE", "/api/v1/locations/abc", nil, map[string]string{"id": "abc"})
			locations.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			locationID := int64(500)
			dbErr := errors.New("db delete failed")
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(dbErr)

			req := createRequestWithRepo("DELETE", "/api/v1/locations/500", nil, map[string]string{"id": "500"})
			locations.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
