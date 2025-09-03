package commodityattributes_test

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

var _ = Describe("Find Commodity Attributes Endpoint", func() {
	var (
		rec *httptest.ResponseRecorder
		attr1 *types.CommodityAttribute
		attr2 *types.CommodityAttribute
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		attr1 = &types.CommodityAttribute{ID: 1, Name: "Attribute A", CommodityType: types.CommodityTypeProduce}
		attr2 = &types.CommodityAttribute{ID: 2, Name: "Attribute B", CommodityType: types.CommodityTypeProduce}
	})

	performRequest := func(queryParams url.Values, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodGet, "/commodity-attributes/find?"+queryParams.Encode(), url.Values{}, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should find commodity attributes successfully for an admin", func() {
			expectedAttributes := []*types.CommodityAttribute{attr1, attr2}
			expectedOpts := &repos.CommodityAttributeFindOpts{
				Limit:  10,
				Offset: 0,
			}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return(expectedAttributes, int64(len(expectedAttributes)), nil)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.CommodityAttribute]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", len(expectedAttributes)))
			Expect(result.Data).To(HaveLen(len(expectedAttributes)))
			Expect(result.Data[0].ID).To(Equal(attr1.ID))
		})

		It("should find commodity attributes successfully for a normal user", func() {
			expectedAttributes := []*types.CommodityAttribute{attr1, attr2}
			expectedOpts := &repos.CommodityAttributeFindOpts{
				Limit:  10,
				Offset: 0,
			}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return(expectedAttributes, int64(len(expectedAttributes)), nil)

			performRequest(url.Values{}, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.CommodityAttribute]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", len(expectedAttributes)))
			Expect(result.Data).To(HaveLen(len(expectedAttributes)))
		})

		It("should apply limit and offset", func() {
			expectedAttributes := []*types.CommodityAttribute{attr2}
			expectedOpts := &repos.CommodityAttributeFindOpts{
				Limit:  1,
				Offset: 1,
			}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return(expectedAttributes, int64(2), nil)

			params := url.Values{}
			params.Set("limit", "1")
			params.Set("offset", "1")
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.CommodityAttribute]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 2))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(attr2.ID))
		})

		It("should filter by multiple IDs", func() {
			ids := []int64{1, 2}
			expectedAttributes := []*types.CommodityAttribute{attr1, attr2}
			expectedOpts := &repos.CommodityAttributeFindOpts{
				Limit:  10,
				Offset: 0,
				IDs: ids,
			}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return(expectedAttributes, int64(2), nil)

			params := url.Values{}
			for _, id := range ids {
				params.Add("id", strconv.FormatInt(id, 10))
			}
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should filter by multiple commodity types", func() {
			commodityTypes := []types.CommodityType{types.CommodityTypeProduce}
			expectedAttributes := []*types.CommodityAttribute{
				{ID: 1, Name: "Produce Attribute", CommodityType: types.CommodityTypeProduce},
			}
			expectedOpts := &repos.CommodityAttributeFindOpts{
				Limit:  10,
				Offset: 0,
				CommodityTypes: commodityTypes,
			}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return(expectedAttributes, int64(1), nil)

			params := url.Values{}
			for _, ct := range commodityTypes {
				params.Add("commodity_types", ct.String())
			}
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should handle invalid limit parameter gracefully", func() {
			params := url.Values{}
			params.Add("limit", "invalid")
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should handle invalid offset parameter gracefully", func() {
			params := url.Values{}
			params.Add("offset", "invalid")
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Error Paths", func() {
		It("should fail if not authenticated", func() {
			performRequest(url.Values{}, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
