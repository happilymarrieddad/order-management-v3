package commodities_test

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
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodities"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
)

var _ = Describe("Update Commodity Endpoint", func() {
	var (
		rec             *httptest.ResponseRecorder
		payload         commodities.UpdateCommodityPayload
		targetCommodity *types.Commodity
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetCommodity = &types.Commodity{ID: 1, Name: "Old Commodity"}

		payload = commodities.UpdateCommodityPayload{
			Name: utils.Ref("Updated Commodity"),
		}
	})

	performRequest := func(commodityID string, payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPut, "/commodities/"+commodityID, url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should update a commodity successfully for an admin", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), targetCommodity.ID).Return(targetCommodity, true, nil)
			mockCommoditiesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, com *types.Commodity) error {
				Expect(com.Name).To(Equal(utils.Deref(payload.Name)))
				return nil
			})

			performRequest(strconv.FormatInt(targetCommodity.ID, 10), payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.Commodity
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Name).To(Equal(utils.Deref(payload.Name)))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(strconv.FormatInt(targetCommodity.ID, 10), payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			performRequest(strconv.FormatInt(targetCommodity.ID, 10), payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid commodity ID", func() {
			performRequest("invalid-id", payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with a malformed JSON body", func() {
			rec, err := testutils.PerformRequest(router, http.MethodPut, "/commodities/1", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if no fields are provided for update", func() {
			payload.Name = nil
			payload.CommodityType = nil
			performRequest(strconv.FormatInt(targetCommodity.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should return 404 if the commodity to update is not found", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), targetCommodity.ID).Return(nil, false, nil)
			performRequest(strconv.FormatInt(targetCommodity.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on get commodity db error", func() {
			dbErr := errors.New("db error")
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), targetCommodity.ID).Return(nil, false, dbErr)
			performRequest(strconv.FormatInt(targetCommodity.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on update commodity db error", func() {
			dbErr := errors.New("db error")
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), targetCommodity.ID).Return(targetCommodity, true, nil)
			mockCommoditiesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest(strconv.FormatInt(targetCommodity.ID, 10), payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
