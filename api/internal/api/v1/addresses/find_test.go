package addresses_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/addresses"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Addresses Handler", func() {
	Context("when addresses exist", func() {
		It("should return a list of addresses with default pagination", func() {
			foundAddresses := []*types.Address{
				{ID: 1, Line1: "123 A St"},
				{ID: 2, Line1: "456 B St"},
			}
			// The handler should apply default limit/offset when none are provided.
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.AddressFindOpts{Limit: 10, Offset: 0})).Return(foundAddresses, int64(2), nil)

			// The endpoint is a POST to /find with an empty body for defaults.
			req := createRequestWithRepo("POST", "/api/v1/addresses/find", []byte(`{}`), nil)
			addresses.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))

			dataBytes, _ := json.Marshal(result.Data)
			var returnedAddresses []types.Address
			json.Unmarshal(dataBytes, &returnedAddresses)
			Expect(returnedAddresses).To(HaveLen(2))
			Expect(returnedAddresses[0].Line1).To(Equal("123 A St"))
		})

		It("should return a list of addresses with custom pagination", func() {
			foundAddresses := []*types.Address{{ID: 3, Line1: "789 C St"}}
			opts := &repos.AddressFindOpts{Limit: 5, Offset: 5}
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(opts)).Return(foundAddresses, int64(1), nil)

			// Send the pagination options in the request body.
			body, _ := json.Marshal(opts)
			req := createRequestWithRepo("POST", "/api/v1/addresses/find", body, nil)
			addresses.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
		})
	})

	Context("when no addresses exist", func() {
		It("should return an empty list", func() {
			// The handler should still apply default limits even for an empty result.
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.AddressFindOpts{Limit: 10, Offset: 0})).Return([]*types.Address{}, int64(0), nil)

			req := createRequestWithRepo("POST", "/api/v1/addresses/find", []byte(`{}`), nil)
			addresses.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("find query failed")
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			req := createRequestWithRepo("POST", "/api/v1/addresses/find", []byte(`{}`), nil)
			addresses.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
