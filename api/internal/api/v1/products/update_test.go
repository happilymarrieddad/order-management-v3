package products_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/products"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("PUT /products/{id}", func() {
	var (
		pld     products.UpdateProductPayload
		product *types.Product
	)

	BeforeEach(func() {
		product = &types.Product{ID: 1, CompanyID: company.ID, CommodityID: 1, Name: "Old Product Name"}
		pld = products.UpdateProductPayload{
			CommodityID: utils.Ref[int64](2),
			Attributes: []*types.ProductAttributeValue{
				{CommodityAttributeID: 1, Value: "Blue"},
			},
		}
	})

	Context("when authenticated as admin", func() {
		Context("and valid payload", func() {
			It("should update the product", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(product, true, nil)
				mockCommoditiesRepo.EXPECT().Get(gomock.Any(), *pld.CommodityID).Return(&types.Commodity{ID: *pld.CommodityID}, true, nil)
				mockProductsRepo.EXPECT().Update(gomock.Any(), gomock.Any(), pld.Attributes).Return(nil)

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPut, "/products/1", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})

		Context("and product does not exist", func() {
			It("should return 404 Not Found", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(nil, false, nil)

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPut, "/products/1", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("and invalid payload", func() {
			It("should return 400 Bad Request", func() {
				pld.CommodityID = utils.Ref[int64](0) // Invalid
				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPut, "/products/1", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("and commodity not found", func() {
			It("should return 400 Bad Request", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(product, true, nil)
				mockCommoditiesRepo.EXPECT().Get(gomock.Any(), *pld.CommodityID).Return(nil, false, nil)

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPut, "/products/1", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("and repository error", func() {
			It("should return 500 Internal Server Error", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(product, true, nil)
				mockCommoditiesRepo.EXPECT().Get(gomock.Any(), *pld.CommodityID).Return(&types.Commodity{ID: *pld.CommodityID}, true, nil)
				mockProductsRepo.EXPECT().Update(gomock.Any(), gomock.Any(), pld.Attributes).Return(errors.New("db error"))

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPut, "/products/1", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("when unauthenticated", func() {
		It("should return 401 Unauthorized", func() {
			body, _ := json.Marshal(pld)
			req := newAuthenticatedRequest(http.MethodPut, "/products/1", bytes.NewReader(body), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
