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
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Get Commodity Endpoint", func() {
	var (
		rec       *httptest.ResponseRecorder
		commodity *types.Commodity
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		commodity = &types.Commodity{ID: 1, Name: "Test Commodity"}
	})

	performRequest := func(commodityID string, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodGet, "/commodities/"+commodityID, url.Values{}, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should get a commodity successfully for an admin", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), commodity.ID).Return(commodity, true, nil)

			performRequest(strconv.FormatInt(commodity.ID, 10), adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.Commodity
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.ID).To(Equal(commodity.ID))
		})

		It("should get a commodity successfully for a normal user", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), commodity.ID).Return(commodity, true, nil)

			performRequest(strconv.FormatInt(commodity.ID, 10), normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(strconv.FormatInt(commodity.ID, 10), nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid commodity ID", func() {
			performRequest("invalid-id", adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Repository Errors", func() {
		It("should return 404 if the commodity is not found", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), commodity.ID).Return(nil, false, nil)
			performRequest(strconv.FormatInt(commodity.ID, 10), adminUser)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), commodity.ID).Return(nil, false, dbErr)
			performRequest(strconv.FormatInt(commodity.ID, 10), adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
