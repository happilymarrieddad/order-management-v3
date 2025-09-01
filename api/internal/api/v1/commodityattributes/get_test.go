package commodityattributes_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Commodity Attribute Handler", func() {
	Context("when the commodity attribute exists", func() {
		It("should return the commodity attribute for an authenticated user", func() {
			attr := &types.CommodityAttribute{ID: 1, Name: "Color", CommodityType: types.CommodityTypeProduce}
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(attr, true, nil)

			req := newAuthenticatedRequest("GET", "/commodity-attributes/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedAttr types.CommodityAttribute
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAttr)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAttr.ID).To(Equal(attr.ID))
			Expect(returnedAttr.Name).To(Equal(attr.Name))
		})
	})

	Context("when the commodity attribute does not exist", func() {
		It("should return 404 Not Found", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			req := newAuthenticatedRequest("GET", "/commodity-attributes/999", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("commodity attribute not found"))
		})
	})

	Context("when the ID is invalid", func() {
		It("should return 404 Not Found from the router", func() {
			req := newAuthenticatedRequest("GET", "/commodity-attributes/abc", nil, basicUser)
			router.ServeHTTP(rr, req)
			// This is a 404 because the route `/{id:[0-9]+}` does not match
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("db went boom")
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, false, dbErr)

			req := newAuthenticatedRequest("GET", "/commodity-attributes/1", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to get commodity attribute"))
		})
	})
})
