package products_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("GET /products/{id}", func() {
	var (
		product *types.Product
		otherCompanyProduct *types.Product
		rec *httptest.ResponseRecorder // Declare rec here
	)

	BeforeEach(func() {
		product = &types.Product{ID: 1, CompanyID: company.ID, Name: "Test Product"}
		otherCompanyProduct = &types.Product{ID: 2, CompanyID: 99, Name: "Other Company Product"}
		rec = httptest.NewRecorder() // Initialize rec here
	})

	Context("when authenticated as admin", func() {
		Context("and product exists", func() {
			It("should return 200 OK with the product", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(product, true, nil)

				req := newAuthenticatedRequest(http.MethodGet, "/products/1", nil, adminUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusOK))
				var respProduct types.Product
				Expect(json.NewDecoder(rec.Body).Decode(&respProduct)).To(Succeed())
				Expect(respProduct.ID).To(Equal(product.ID))
			})
		})

		Context("and product from another company exists", func() {
			It("should return 200 OK with the product", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), otherCompanyProduct.ID).Return(otherCompanyProduct, true, nil)

				req := newAuthenticatedRequest(http.MethodGet, "/products/2", nil, adminUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusOK))
				var respProduct types.Product
				Expect(json.NewDecoder(rec.Body).Decode(&respProduct)).To(Succeed())
				Expect(respProduct.ID).To(Equal(otherCompanyProduct.ID))
			})
		})

		Context("and product does not exist", func() {
			It("should return 404 Not Found", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(nil, false, nil)

				req := newAuthenticatedRequest(http.MethodGet, "/products/1", nil, adminUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("and repository returns an error", func() {
			It("should return 500 Internal Server Error", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(nil, false, errors.New("db error"))

				req := newAuthenticatedRequest(http.MethodGet, "/products/1", nil, adminUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("and invalid product ID", func() {
			It("should return 404 Not Found", func() {
				req := newAuthenticatedRequest(http.MethodGet, "/products/invalid", nil, adminUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("when authenticated as normal user", func() {
		Context("and product exists in own company", func() {
			It("should return 200 OK with the product", func() {
				product.CompanyID = normalUser.CompanyID // Ensure product belongs to normal user's company
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(product, true, nil)

				req := newAuthenticatedRequest(http.MethodGet, "/products/1", nil, normalUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusOK))
				var respProduct types.Product
				Expect(json.NewDecoder(rec.Body).Decode(&respProduct)).To(Succeed())
				Expect(respProduct.ID).To(Equal(product.ID))
			})
		})

		Context("and product exists in another company", func() {
			It("should return 403 Forbidden", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), otherCompanyProduct.ID).Return(otherCompanyProduct, true, nil)

				req := newAuthenticatedRequest(http.MethodGet, "/products/2", nil, normalUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("and product does not exist", func() {
			It("should return 404 Not Found", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(nil, false, nil)

				req := newAuthenticatedRequest(http.MethodGet, "/products/1", nil, normalUser)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("when unauthenticated", func() {
		It("should return 401 Unauthorized", func() {
			req := newAuthenticatedRequest(http.MethodGet, "/products/1", nil, nil)
				router.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})