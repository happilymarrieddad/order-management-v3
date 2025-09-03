package addresses_test

import (
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
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Get Address Endpoint", func() {
	var (
		rec     *httptest.ResponseRecorder
		address *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		address = &types.Address{ID: 1, Line1: "123 Main St"}
	})

	performRequest := func(addressID string, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodGet, "/addresses/"+addressID, url.Values{}, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should get an address successfully for an admin", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), address.ID).Return(address, true, nil)

			performRequest(strconv.FormatInt(address.ID, 10), adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.Address
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.ID).To(Equal(address.ID))
		})

		It("should get an address successfully for a normal user", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), address.ID).Return(address, true, nil)

			performRequest(strconv.FormatInt(address.ID, 10), normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(strconv.FormatInt(address.ID, 10), nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid address ID", func() {
			performRequest("invalid-id", adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Repository Errors", func() {
		It("should return 404 if the address is not found", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), address.ID).Return(nil, false, nil)
			performRequest(strconv.FormatInt(address.ID, 10), adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockAddressesRepo.EXPECT().Get(gomock.Any(), address.ID).Return(nil, false, dbErr)
			performRequest(strconv.FormatInt(address.ID, 10), adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
