package commodities_test

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

var _ = Describe("Find Commodities Endpoint", func() {
	var (
		rec   *httptest.ResponseRecorder
		com1 *types.Commodity
		com2 *types.Commodity
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		com1 = &types.Commodity{ID: 1, Name: "Commodity A", CommodityType: types.CommodityTypeProduce}
		com2 = &types.Commodity{ID: 2, Name: "Commodity B", CommodityType: types.CommodityTypeProduce}
	})

	performRequest := func(queryParams url.Values, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodGet, "/commodities/find?"+queryParams.Encode(), nil, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should find commodities successfully for an admin", func() {
			queryParams := url.Values{}
			expectedOpts := &repos.FindCommoditiesOpts{
				Limit:  10,
				Offset: 0,
			}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Commodity{com1, com2}, int64(2), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Commodity]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 2))
			Expect(result.Data).To(HaveLen(2))
		})

		It("should find commodities successfully for a normal user", func() {
			queryParams := url.Values{}
			expectedOpts := &repos.FindCommoditiesOpts{
				Limit:  10,
				Offset: 0,
			}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Commodity{com1, com2}, int64(2), nil)

			performRequest(queryParams, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should apply limit and offset", func() {
			queryParams := url.Values{}
			queryParams.Set("limit", "1")
			queryParams.Set("offset", "1")

			expectedOpts := &repos.FindCommoditiesOpts{
				Limit:  1,
				Offset: 1,
			}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Commodity{com2}, int64(2), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Commodity]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 2))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(com2.ID))
		})

		It("should filter by ids", func() {
			queryParams := url.Values{}
			queryParams.Add("id", strconv.FormatInt(com1.ID, 10))

			expectedOpts := &repos.FindCommoditiesOpts{
				IDs:    []int64{com1.ID},
				Limit:  10,
				Offset: 0,
			}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Commodity{com1}, int64(1), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Commodity]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 1))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(com1.ID))
		})

		It("should filter by name", func() {
			queryParams := url.Values{}
			queryParams.Add("name", com1.Name)

			expectedOpts := &repos.FindCommoditiesOpts{
				Names:  []string{com1.Name},
				Limit:  10,
				Offset: 0,
			}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Commodity{com1}, int64(1), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Commodity]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 1))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(com1.ID))
		})

		It("should filter by commodity_type", func() {
			queryParams := url.Values{}
			queryParams.Add("commodity_type", strconv.Itoa(int(com1.CommodityType)))

			expectedOpts := &repos.FindCommoditiesOpts{
				CommodityType: com1.CommodityType,
				Limit:         10,
				Offset:        0,
			}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Commodity{com1}, int64(1), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Commodity]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 1))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(com1.ID))
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
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
