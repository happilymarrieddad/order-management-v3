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
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Get Commodity Attribute Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		targetAttribute *types.CommodityAttribute
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetAttribute = &types.CommodityAttribute{ID: 1, Name: "Test Attribute", CommodityType: types.CommodityTypeProduce}
	})

	performRequest := func(attributeID string, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodGet, "/commodity-attributes/"+attributeID, url.Values{}, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should get a commodity attribute successfully for an admin", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(targetAttribute, true, nil)

			performRequest(strconv.FormatInt(targetAttribute.ID, 10), adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.CommodityAttribute
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.ID).To(Equal(targetAttribute.ID))
		})

		It("should get a commodity attribute successfully for a normal user", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(targetAttribute, true, nil)

			performRequest(strconv.FormatInt(targetAttribute.ID, 10), normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid ID", func() {
			performRequest("invalid-id", adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Repository Errors", func() {
		It("should return 404 if the commodity attribute is not found", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(nil, false, nil)
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(nil, false, dbErr)
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
