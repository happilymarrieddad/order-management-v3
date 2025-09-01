package commodities_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Commodity Handler", func() {
	Context("when the commodity exists", func() {
		It("should return the commodity for an authenticated user", func() {
			commodity := &types.Commodity{ID: 1, Name: "Potatoes", CommodityType: types.CommodityTypeProduce}
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(commodity, true, nil)

			req := newAuthenticatedRequest("GET", "/commodities/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCommodity types.Commodity
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCommodity)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCommodity.ID).To(Equal(commodity.ID))
			Expect(returnedCommodity.Name).To(Equal(commodity.Name))
		})
	})

	Context("when the commodity does not exist", func() {
		It("should return 404 Not Found", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			req := newAuthenticatedRequest("GET", "/commodities/999", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("commodity not found"))
		})
	})

	Context("when the ID is invalid", func() {
		It("should return 404 Not Found from the router", func() {
			req := newAuthenticatedRequest("GET", "/commodities/abc", nil, basicUser)
			router.ServeHTTP(rr, req)
			// This is a 404 because the route `/{id:[0-9]+}` does not match
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("db went boom")
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, false, dbErr)

			req := newAuthenticatedRequest("GET", "/commodities/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to get commodity"))
		})
	})
})
