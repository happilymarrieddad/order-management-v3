package addresses_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
)

var _ = Describe("Find Addresses Handler", func() {
	Context("when addresses exist", func() {
		It("should return a list of addresses for an admin user", func() {
			foundAddresses := []*types.Address{
				{ID: 1, Line1: "123 A St"},
				{ID: 2, Line1: "456 B St"},
			}
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.AddressFindOpts{Limit: 10, Offset: 0})).Return(foundAddresses, int64(2), nil)

			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodPost, "/addresses/find", url.Values{}, bytes.NewBufferString(`{}`), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult[[]*types.Address] // Changed to generic FindResult
			err = json.NewDecoder(rr.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))
			Expect(result.Data).To(HaveLen(2))
			Expect(result.Data[0].Line1).To(Equal("123 A St"))
		})

		It("should return a list of addresses with custom pagination", func() {
			foundAddresses := []*types.Address{{ID: 3, Line1: "789 C St"}}
			opts := &repos.AddressFindOpts{Limit: 5, Offset: 5}
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(opts)).Return(foundAddresses, int64(1), nil)

			body, _ := json.Marshal(opts)
			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodPost, "/addresses/find", url.Values{}, bytes.NewBuffer(body), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult[[]*types.Address] // Changed to generic FindResult
			err = json.NewDecoder(rr.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
		})
	})

	Context("when no addresses exist", func() {
		It("should return an empty list", func() {
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.AddressFindOpts{Limit: 10, Offset: 0})).Return([]*types.Address{}, int64(0), nil)

			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodPost, "/addresses/find", url.Values{}, bytes.NewBufferString(`{}`), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult[[]*types.Address] // Changed to generic FindResult
			err = json.NewDecoder(rr.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("find query failed")
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodPost, "/addresses/find", url.Values{}, bytes.NewBufferString(`{}`), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to find addresses"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodPost, "/addresses/find", url.Values{}, bytes.NewBufferString(`{}`), basicUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			// Use testutils.PerformRequest
			var err error
			rr, err = testutils.PerformRequest(router, http.MethodPost, "/addresses/find", url.Values{}, bytes.NewBufferString(`{}`), nil, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})