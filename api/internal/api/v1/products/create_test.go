package products_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/products"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("POST /products", func() {
	var (
		pld products.CreateProductPayload
		product *types.Product
	)

	BeforeEach(func() {
		pld = products.CreateProductPayload{
			CompanyID:   company.ID,
			CommodityID: 1,
			Attributes: []*types.ProductAttributeValue{
				{CommodityAttributeID: 1, Value: "Red"},
			},
		}
		product = &types.Product{ID: 1, CompanyID: pld.CompanyID, CommodityID: pld.CommodityID, Name: "Red Apple"}
	})

	Context("when authenticated as admin", func() {
		Context("and valid payload", func() {
			It("should create a product", func() {
				mockCommoditiesRepo.EXPECT().Get(gomock.Any(), pld.CommodityID).Return(&types.Commodity{ID: pld.CommodityID}, true, nil)
				mockProductsRepo.EXPECT().Create(gomock.Any(), gomock.Any(), pld.Attributes).DoAndReturn(func(_ any, p *types.Product, _ []*types.ProductAttributeValue) error {
					*p = *product
					return nil
				})

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPost, "/products", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusCreated))
				var respProduct types.Product
				Expect(json.NewDecoder(rr.Body).Decode(&respProduct)).To(Succeed())
				Expect(respProduct.ID).To(Equal(product.ID))
			})
		})

		Context("and invalid payload", func() {
			It("should return 400 Bad Request", func() {
				pld.CompanyID = 0 // Invalid
				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPost, "/products", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("and commodity not found", func() {
			It("should return 400 Bad Request", func() {
				mockCommoditiesRepo.EXPECT().Get(gomock.Any(), pld.CommodityID).Return(nil, false, nil)

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPost, "/products", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("and repository error", func() {
			It("should return 500 Internal Server Error", func() {
				mockCommoditiesRepo.EXPECT().Get(gomock.Any(), pld.CommodityID).Return(&types.Commodity{ID: pld.CommodityID}, true, nil)
				mockProductsRepo.EXPECT().Create(gomock.Any(), gomock.Any(), pld.Attributes).Return(errors.New("db error"))

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPost, "/products", bytes.NewReader(body), adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("when authenticated as normal user", func() {
		Context("and valid payload for own company", func() {
			It("should create a product", func() {
				pld.CompanyID = normalUser.CompanyID // Own company
				mockCommoditiesRepo.EXPECT().Get(gomock.Any(), pld.CommodityID).Return(&types.Commodity{ID: pld.CommodityID}, true, nil)
				mockProductsRepo.EXPECT().Create(gomock.Any(), gomock.Any(), pld.Attributes).DoAndReturn(func(_ any, p *types.Product, _ []*types.ProductAttributeValue) error {
					*p = *product
					return nil
				})

				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPost, "/products", bytes.NewReader(body), normalUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("and valid payload for other company", func() {
			It("should return 403 Forbidden", func() {
				pld.CompanyID = 999 // Other company
				body, _ := json.Marshal(pld)
				req := newAuthenticatedRequest(http.MethodPost, "/products", bytes.NewReader(body), normalUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusForbidden))
			})
		})
	})

	Context("when unauthenticated", func() {
		It("should return 401 Unauthorized", func() {
			body, _ := json.Marshal(pld)
			req := newAuthenticatedRequest(http.MethodPost, "/products", bytes.NewReader(body), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
