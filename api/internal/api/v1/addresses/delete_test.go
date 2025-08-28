package addresses_test

import (
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Delete Address Handler", func() {
	Context("with a valid request", func() {
		It("should delete the address successfully", func() {
			addressID := int64(123)
			mockAddressesRepo.EXPECT().Delete(gomock.Any(), addressID).Return(nil)

			req := createRequestWithRepo("DELETE", "/api/v1/addresses/123", nil, map[string]string{"id": "123"})
			addresses.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			req := createRequestWithRepo("DELETE", "/api/v1/addresses/abc", nil, map[string]string{"id": "abc"})
			addresses.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid address ID"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			addressID := int64(500)
			dbErr := errors.New("foreign key constraint fails")

			mockAddressesRepo.EXPECT().Delete(gomock.Any(), addressID).Return(dbErr)

			req := createRequestWithRepo("DELETE", "/api/v1/addresses/500", nil, map[string]string{"id": "500"})
			addresses.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
