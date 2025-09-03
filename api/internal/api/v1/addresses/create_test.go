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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Create Address Handler", func() {
	var (
		payload      addresses.CreateAddressPayload
		payloadBytes []byte
		newAddress   *types.Address
		err          error
	)

	BeforeEach(func() {
		payload = addresses.CreateAddressPayload{
			Line1:      "123 Main St",
			City:       "Anytown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "USA",
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		newAddress = &types.Address{
			ID:         1,
			Line1:      payload.Line1,
			City:       payload.City,
			State:      payload.State,
			PostalCode: payload.PostalCode,
			Country:    payload.Country,
			GlobalCode: "849VCWC8+R9",
		}
	})

	Context("when creation is successful", func() {
		It("should return 201 Created with the new address", func() {
			// The repo's Create method is responsible for geocoding and persistence.
			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(newAddress, nil)

			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPost, "/addresses", url.Values{}, bytes.NewBuffer(payloadBytes), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedAddr types.Address
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAddr.ID).To(Equal(newAddress.ID))
			Expect(returnedAddr.GlobalCode).To(Equal(newAddress.GlobalCode))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a malformed JSON body", func() {
			// Use testutils.PerformRequest
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPost, "/addresses", url.Values{}, bytes.NewBufferString(`{]`), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a validation error (missing line1)", func() {
			// This payload is missing the required 'line_1' field.
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPost, "/addresses", url.Values{}, bytes.NewBufferString(`{"city": "Anytown", "state": "CA", "postal_code": "12345", "country": "USA"}`), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'line1' is required."))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 Internal Server Error", func() {
			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("db insert failed"))

			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPost, "/addresses", url.Values{}, bytes.NewBuffer(payloadBytes), adminUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to create address"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden", func() {
			var reqErr error
			rr, reqErr = testutils.PerformRequest(router, http.MethodPost, "/addresses", url.Values{}, bytes.NewBuffer(payloadBytes), basicUser, mockGlobalRepo)
			Expect(reqErr).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})
	})
})