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

var _ = Describe("Create Address Handler", func() {
	var createPayload map[string]interface{}

	BeforeEach(func() {
		createPayload = map[string]interface{}{
			"line_1":      "123 Main St",
			"city":        "Anytown",
			"state":       "CA",
			"postal_code": "12345",
		}
	})

	Context("with a valid request", func() {
		It("should create an address successfully", func() {
			body, _ := json.Marshal(createPayload)

			createdAddress := &types.Address{
				ID:         1,
				Line1:      "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
			}

			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(createdAddress, nil)

			req := createRequestWithRepo("POST", "/api/v1/addresses", body, nil)
			addresses.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedAddress types.Address
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAddress.ID).To(Equal(int64(1)))
			Expect(returnedAddress.Line1).To(Equal("123 Main St"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a malformed JSON body", func() {
			body := []byte(`{"line_1": "bad json",`)
			req := createRequestWithRepo("POST", "/api/v1/addresses", body, nil)
			addresses.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a missing required field (line_1)", func() {
			delete(createPayload, "line_1")
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/addresses", body, nil)
			addresses.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			body, _ := json.Marshal(createPayload)
			dbErr := errors.New("unexpected database error")

			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, dbErr)

			req := createRequestWithRepo("POST", "/api/v1/addresses", body, nil)
			addresses.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
