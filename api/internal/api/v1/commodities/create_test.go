package commodities_test

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
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodities"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Create Commodity Endpoint", func() {
	var (
		rec     *httptest.ResponseRecorder
		payload commodities.CreateCommodityPayload
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		payload = commodities.CreateCommodityPayload{
			Name:          "Test Commodity",
			CommodityType: types.CommodityTypeProduce,
		}
	})

	performRequest := func(payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPost, "/commodities", url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should create a commodity successfully for an admin", func() {
			mockCommoditiesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, com *types.Commodity) error {
				com.ID = 1
				return nil
			})

			performRequest(payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusCreated))
			var result types.Commodity
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Name).To(Equal(payload.Name))
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
			rec, err := testutils.PerformRequest(router, http.MethodPost, "/commodities", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if a required field is missing", func() {
			payload.Name = ""
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Repository Errors", func() {
		It("should return 500 on commodity creation db error", func() {
			dbErr := errors.New("db error")
			mockCommoditiesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
