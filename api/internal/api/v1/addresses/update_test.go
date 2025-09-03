package addresses_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

var _ = Describe("Update Address Endpoint", func() {
	var (
		rec           *httptest.ResponseRecorder
		payload       addresses.UpdateAddressPayload
		targetAddress *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetAddress = &types.Address{ID: 1, Line1: "Old Address"}

		payload = addresses.UpdateAddressPayload{
			Line1: utils.Ref("Updated Address"),
		}
	})

	performRequest := func(addressID string, payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPut, "/addresses/"+addressID, url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should update an address successfully for an admin", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), targetAddress.ID).Return(targetAddress, true, nil)
			mockAddressesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, addr *types.Address) error {
				Expect(addr.Line1).To(Equal(*payload.Line1))
				return nil
			})

			performRequest(strconv.FormatInt(targetAddress.ID, 10), payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.Address
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Line1).To(Equal(utils.Deref(payload.Line1)))
		})

		It("should update an address successfully for a normal user", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), targetAddress.ID).Return(targetAddress, true, nil)
			mockAddressesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, addr *types.Address) error {
				Expect(addr.Line1).To(Equal(*payload.Line1))
				return nil
			})

			performRequest(strconv.FormatInt(targetAddress.ID, 10), payload, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(strconv.FormatInt(targetAddress.ID, 10), payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid address ID", func() {
			performRequest("invalid-id", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with a malformed JSON body", func() {
			rec, err := testutils.PerformRequest(router, http.MethodPut, "/addresses/1", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if no fields are provided for update", func() {
			payload.Line1 = nil
			performRequest(strconv.FormatInt(targetAddress.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should return 404 if the address to update is not found", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), targetAddress.ID).Return(nil, false, nil)
			performRequest(strconv.FormatInt(targetAddress.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on get address db error", func() {
			dbErr := errors.New("db error")
			mockAddressesRepo.EXPECT().Get(gomock.Any(), targetAddress.ID).Return(nil, false, dbErr)
			performRequest(strconv.FormatInt(targetAddress.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on update address db error", func() {
			dbErr := errors.New("db error")
			mockAddressesRepo.EXPECT().Get(gomock.Any(), targetAddress.ID).Return(targetAddress, true, nil)
			mockAddressesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest(strconv.FormatInt(targetAddress.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
