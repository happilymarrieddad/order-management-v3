package addresses_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Address Handler", func() {
	Context("when the address exists", func() {
		It("should return the address successfully", func() {
			addressID := int64(123)
			expectedAddress := &types.Address{ID: addressID, Line1: "123 Found St"}

			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(expectedAddress, true, nil)

			req := createRequestWithRepo("GET", "/api/v1/addresses/123", nil, map[string]string{"id": "123"})
			addresses.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedAddress types.Address
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAddress.ID).To(Equal(addressID))
			Expect(returnedAddress.Line1).To(Equal("123 Found St"))
		})
	})

	Context("when the address does not exist", func() {
		It("should return 404 Not Found", func() {
			addressID := int64(404)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, nil)

			req := createRequestWithRepo("GET", "/api/v1/addresses/404", nil, map[string]string{"id": "404"})
			addresses.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			req := createRequestWithRepo("GET", "/api/v1/addresses/abc", nil, map[string]string{"id": "abc"})
			addresses.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid address ID"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			addressID := int64(500)
			dbErr := errors.New("database connection lost")

			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, dbErr)

			req := createRequestWithRepo("GET", "/api/v1/addresses/500", nil, map[string]string{"id": "500"})
			addresses.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
