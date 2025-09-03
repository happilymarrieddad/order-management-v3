package commodityattributes_test

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
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

var _ = Describe("Update Commodity Attribute Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		payload    commodityattributes.UpdateCommodityAttributePayload
		targetAttribute *types.CommodityAttribute
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetAttribute = &types.CommodityAttribute{ID: 1, Name: "Old Attribute", CommodityType: types.CommodityTypeProduce}

		payload = commodityattributes.UpdateCommodityAttributePayload{
			Name: utils.Ref("Updated Attribute"),
		}
	})

	performRequest := func(attributeID string, payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPut, "/commodity-attributes/"+attributeID, url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should update a commodity attribute successfully for an admin", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(targetAttribute, true, nil)
			mockCommodityAttributesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, ca *types.CommodityAttribute) error {
				Expect(ca.Name).To(Equal(*payload.Name))
				return nil
			})

			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.CommodityAttribute
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Name).To(Equal(*payload.Name))
		})

		It("should update only the name if other fields are not provided", func() {
			// No other fields in payload for now
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(targetAttribute, true, nil)
			mockCommodityAttributesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, ca *types.CommodityAttribute) error {
				Expect(ca.Name).To(Equal(*payload.Name))
				return nil
			})

			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid ID", func() {
			performRequest("invalid-id", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with a malformed JSON body", func() {
			rec, err := testutils.PerformRequest(router, http.MethodPut, "/commodity-attributes/1", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if no fields are provided for update", func() {
			payload.Name = nil
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should return 404 if the commodity attribute to update is not found", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(nil, false, nil)
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on get commodity attribute db error", func() {
			dbErr := errors.New("db error")
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(nil, false, dbErr)
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on update commodity attribute db error", func() {
			dbErr := errors.New("db error")
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), targetAttribute.ID).Return(targetAttribute, true, nil)
			mockCommodityAttributesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest(strconv.FormatInt(targetAttribute.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
