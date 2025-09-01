package commodityattributes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Commodity Attributes Handler", func() {
	Context("when commodity attributes exist", func() {
		It("should return a list of attributes for an admin user", func() {
			foundCommodityAttributes := []*types.CommodityAttribute{
				{ID: 1, Name: "Color"},
				{ID: 2, Name: "Size"},
			}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.CommodityAttributeFindOpts{Limit: 10, Offset: 0})).Return(foundCommodityAttributes, int64(2), nil)

			req := newAuthenticatedRequest("POST", "/commodity-attributes/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))

			// We need to remarshal and unmarshal the data to compare it properly
			dataBytes, _ := json.Marshal(result.Data)
			var returnedCommodityAttributes []types.CommodityAttribute
			json.Unmarshal(dataBytes, &returnedCommodityAttributes)
			Expect(returnedCommodityAttributes).To(HaveLen(2))
			Expect(returnedCommodityAttributes[0].Name).To(Equal("Color"))
		})

		It("should return a list of commodity attributes with custom pagination", func() {
			foundCommodityAttributes := []*types.CommodityAttribute{{ID: 3, Name: "Weight"}}
			opts := &repos.CommodityAttributeFindOpts{Limit: 5, Offset: 5}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(opts)).Return(foundCommodityAttributes, int64(1), nil)

			body, _ := json.Marshal(opts)
			req := newAuthenticatedRequest("POST", "/commodity-attributes/find", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
		})
	})

	Context("when no commodity attributes exist", func() {
		It("should return an empty list for an admin user", func() {
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.CommodityAttributeFindOpts{Limit: 10, Offset: 0})).Return([]*types.CommodityAttribute{}, int64(0), nil)

			req := newAuthenticatedRequest("POST", "/commodity-attributes/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 Internal Server Error", func() {
			dbErr := errors.New("find query failed")
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			req := newAuthenticatedRequest("POST", "/commodity-attributes/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to find commodity attributes"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden", func() {
			req := newAuthenticatedRequest("POST", "/commodity-attributes/find", bytes.NewBufferString(`{}`), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})
	})
})
