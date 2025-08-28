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

var _ = Describe("Update Address Handler", func() {
	var updatePayload map[string]interface{}

	BeforeEach(func() {
		updatePayload = map[string]interface{}{
			"line_1":      "456 Updated Ave",
			"city":        "Newville",
			"state":       "NY",
			"postal_code": "54321",
		}
	})

	Context("with a valid request", func() {
		It("should update the address successfully", func() {
			addressID := int64(123)
			body, _ := json.Marshal(updatePayload)

			// Mock the Get call to find the existing address
			existingAddress := &types.Address{ID: addressID, Line1: "Old St"}
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(existingAddress, true, nil)

			// Mock the Update call
			mockAddressesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			req := createRequestWithRepo("PUT", "/api/v1/addresses/123", body, map[string]string{"id": "123"})
			addresses.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedAddress types.Address
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAddress.Line1).To(Equal("456 Updated Ave"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/addresses/abc", body, map[string]string{"id": "abc"})
			addresses.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a missing required field", func() {
			delete(updatePayload, "city")
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/addresses/123", body, map[string]string{"id": "123"})
			addresses.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the target address does not exist", func() {
		It("should return 404 Not Found", func() {
			addressID := int64(404)
			body, _ := json.Marshal(updatePayload)

			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, nil)

			req := createRequestWithRepo("PUT", "/api/v1/addresses/404", body, map[string]string{"id": "404"})
			addresses.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 on Get error", func() {
			addressID := int64(500)
			body, _ := json.Marshal(updatePayload)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, errors.New("get error"))

			req := createRequestWithRepo("PUT", "/api/v1/addresses/500", body, map[string]string{"id": "500"})
			addresses.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on Update error", func() {
			addressID := int64(123)
			body, _ := json.Marshal(updatePayload)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(&types.Address{ID: addressID}, true, nil)
			mockAddressesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update error"))

			req := createRequestWithRepo("PUT", "/api/v1/addresses/123", body, map[string]string{"id": "123"})
			addresses.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
