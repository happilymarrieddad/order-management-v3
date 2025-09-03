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
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Find Addresses Endpoint", func() {
	var (
		rec   *httptest.ResponseRecorder
		addr1 *types.Address
		addr2 *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		addr1 = &types.Address{ID: 1, Line1: "123 Main St"}
		addr2 = &types.Address{ID: 2, Line1: "456 Oak Ave"}
	})

	performRequest := func(queryParams url.Values, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodGet, "/addresses/find?"+queryParams.Encode(), nil, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should find addresses successfully for an admin", func() {
			queryParams := url.Values{}
			expectedOpts := &repos.AddressFindOpts{
				Limit:  10,
				Offset: 0,
			}
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Address{addr1, addr2}, int64(2), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Address]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 2))
			Expect(result.Data).To(HaveLen(2))
		})

		It("should find addresses successfully for a normal user", func() {
			queryParams := url.Values{}
			expectedOpts := &repos.AddressFindOpts{
				Limit:  10,
				Offset: 0,
			}
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Address{addr1, addr2}, int64(2), nil)

			performRequest(queryParams, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should apply limit and offset", func() {
			queryParams := url.Values{}
			queryParams.Set("limit", "1")
			queryParams.Set("offset", "1")

			expectedOpts := &repos.AddressFindOpts{
				Limit:  1,
				Offset: 1,
			}
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Address{addr2}, int64(2), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Address]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 2))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(addr2.ID))
		})

		It("should filter by ids", func() {
			queryParams := url.Values{}
			queryParams.Add("id", strconv.FormatInt(addr1.ID, 10))

			expectedOpts := &repos.AddressFindOpts{
				IDs:    []int64{addr1.ID},
				Limit:  10,
				Offset: 0,
			}
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Address{addr1}, int64(1), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Address]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 1))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(addr1.ID))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(url.Values{}, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("Error Paths", func() {
		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockAddressesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
