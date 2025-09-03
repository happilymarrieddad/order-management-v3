package addresses_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update Address Handler", func() {
	var (
		payload      addresses.UpdateAddressPayload
		payloadBytes []byte
		existingAddr *types.Address
		err          error
	)

	BeforeEach(func() {
		payload = addresses.UpdateAddressPayload{
			Line1:   utils.Ref("456 Updated Ave"),
			City:    utils.Ref("Newville"),
			Country: utils.Ref("CAN"),
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		existingAddr = &types.Address{
			ID:         1,
			Line1:      "123 Original St",
			City:       "Oldtown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "USA",
		}
	})

	Context("when update is successful", func() {
		It("should return 200 OK with the updated address", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingAddr, true, nil)
			// The repo's Update method is responsible for geocoding and persistence.
			mockAddressesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPut, "/addresses/1", url.Values{}, bytes.NewBuffer(payloadBytes), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedAddr types.Address
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAddr.ID).To(Equal(existingAddr.ID))
			Expect(returnedAddr.Line1).To(Equal(*payload.Line1)) // Check that a field was updated
			Expect(returnedAddr.Country).To(Equal(*payload.Country))
		})
	})

	Context("when the address to update is not found", func() {
		It("should return 404 Not Found", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPut, "/addresses/999", url.Values{}, bytes.NewBuffer(payloadBytes), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a malformed JSON body", func() {
			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPut, "/addresses/1", url.Values{}, bytes.NewBufferString(`{]`), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 404 for a non-integer ID", func() {
			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPut, "/addresses/abc", url.Values{}, bytes.NewBuffer(payloadBytes), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 on update failure", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingAddr, true, nil)
			mockAddressesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db update failed"))

			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPut, "/addresses/1", url.Values{}, bytes.NewBuffer(payloadBytes), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to update address"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPut, "/addresses/1", url.Values{}, bytes.NewBuffer(payloadBytes), basicUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPut, "/addresses/1", url.Values{}, bytes.NewBuffer(payloadBytes), nil, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})