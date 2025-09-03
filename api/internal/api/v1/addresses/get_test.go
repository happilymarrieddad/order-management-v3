package addresses_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
)

var _ = Describe("Get Address Handler", func() {
	Context("when the address exists", func() {
		It("should return the address successfully", func() {
			addressID := int64(123)
			expectedAddress := &types.Address{ID: addressID, Line1: "123 Found St"}

			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(expectedAddress, true, nil)

			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodGet, "/addresses/123", url.Values{}, nil, basicUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedAddress types.Address
			err = json.Unmarshal(rr.Body.Bytes(), &returnedAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAddress.ID).To(Equal(addressID))
			Expect(returnedAddress.Line1).To(Equal("123 Found St"))
		})
	})

	Context("when the address does not exist", func() {
		It("should return 404 Not Found", func() {
			addressID := int64(404)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, nil)

			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodGet, "/addresses/404", url.Values{}, nil, basicUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 404 for a non-integer ID", func() {
			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodGet, "/addresses/abc", url.Values{}, nil, basicUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			// This is a router-level 404 because the route `/{id:[0-9]+}` does not match
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			addressID := int64(500)
			dbErr := errors.New("database connection lost")

			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, dbErr)

			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodGet, "/addresses/500", url.Values{}, nil, basicUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to get address"))
		})
	})
})