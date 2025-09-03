package addresses_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Create Address Endpoint", func() {
	var (
		rec     *httptest.ResponseRecorder
		payload addresses.CreateAddressPayload
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		payload = addresses.CreateAddressPayload{
			Line1:      "123 Main St",
			City:       "Anytown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "USA",
		}
	})

	performRequest := func(payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPost, "/addresses", url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should create an address successfully for an admin", func() {
			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, addr *types.Address) (*types.Address, error) {
				addr.ID = 1
				return addr, nil
			})

			performRequest(payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusCreated))
			var result types.Address
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Line1).To(Equal(payload.Line1))
			Expect(result.ID).To(BeNumerically(">", 0))
		})

		It("should create an address successfully for a normal user", func() {
			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, addr *types.Address) (*types.Address, error) {
				addr.ID = 1
				return addr, nil
			})

			performRequest(payload, normalUser)

			Expect(rec.Code).To(Equal(http.StatusCreated))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with a malformed JSON body", func() {
			rec, err := testutils.PerformRequest(router, http.MethodPost, "/addresses", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if a required field is missing", func() {
			payload.Line1 = ""
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Repository Errors", func() {
		It("should return 500 on address creation db error", func() {
			dbErr := errors.New("db error")
			mockAddressesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, dbErr)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
