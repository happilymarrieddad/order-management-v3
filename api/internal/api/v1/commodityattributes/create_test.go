package commodityattributes_test

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
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Create Commodity Attribute Endpoint", func() {
	var (
		rec     *httptest.ResponseRecorder
		payload commodityattributes.CreateCommodityAttributePayload
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		payload = commodityattributes.CreateCommodityAttributePayload{
			Name:          "Test Attribute",
			CommodityType: types.CommodityTypeProduce,
		}
	})

	performRequest := func(payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPost, "/commodity-attributes", url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should create a commodity attribute successfully for an admin", func() {
			mockCommodityAttributesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, ca *types.CommodityAttribute) error {
				ca.ID = 3 // Simulate ID generation
				return nil
			})

			performRequest(payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusCreated))
			var result types.CommodityAttribute
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Name).To(Equal(payload.Name))
			Expect(result.CommodityType).To(Equal(payload.CommodityType))
			Expect(result.ID).To(BeNumerically(">", 0))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			performRequest(payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with a malformed JSON body", func() {
			rec, err := testutils.PerformRequest(router, http.MethodPost, "/commodity-attributes", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if a required field is missing", func() {
			payload.Name = ""
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if an invalid commodity type is provided", func() {
			payload.CommodityType = 999 // Invalid type
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Repository Errors", func() {
		It("should return 500 on commodity attribute creation db error", func() {
			dbErr := errors.New("db error")
			mockCommodityAttributesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})