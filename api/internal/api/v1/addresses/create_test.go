package addresses_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

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

			req := newAuthenticatedRequest("POST", "/addresses", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

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
			req := newAuthenticatedRequest("POST", "/addresses", bytes.NewBufferString(`{]`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a validation error (missing line1)", func() {
			// This payload is missing the required 'line_1' field.
			req := newAuthenticatedRequest("POST", "/addresses", bytes.NewBufferString(`{"city": "Anytown", "state": "CA", "postal_code": "12345", "country": "USA"}`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'line1' is required."))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 Internal Server Error", func() {
			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("db insert failed"))
			req := newAuthenticatedRequest("POST", "/addresses", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to create address"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden", func() {
			req := newAuthenticatedRequest("POST", "/addresses", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})
	})
})
