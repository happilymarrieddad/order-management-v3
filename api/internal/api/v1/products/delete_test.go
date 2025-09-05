package products_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("DELETE /products/{id}", func() {
	var (
		product *types.Product
	)

	BeforeEach(func() {
		product = &types.Product{ID: 1, CompanyID: company.ID, Name: "Test Product"}
	})

	Context("when authenticated as admin", func() {
		Context("and product exists", func() {
			It("should return 204 No Content", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(product, true, nil)
				mockProductsRepo.EXPECT().Delete(gomock.Any(), product.ID).Return(nil)

				req := newAuthenticatedRequest(http.MethodDelete, "/products/1", nil, adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("and product does not exist", func() {
			It("should return 404 Not Found", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(nil, false, nil)

				req := newAuthenticatedRequest(http.MethodDelete, "/products/1", nil, adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("and repository returns an error", func() {
			It("should return 500 Internal Server Error", func() {
				mockProductsRepo.EXPECT().Get(gomock.Any(), product.ID).Return(product, true, nil)
				mockProductsRepo.EXPECT().Delete(gomock.Any(), product.ID).Return(errors.New("db error"))

				req := newAuthenticatedRequest(http.MethodDelete, "/products/1", nil, adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			})
		})

				Context("and invalid product ID", func() {
			It("should return 404 Not Found", func() {
				req := newAuthenticatedRequest(http.MethodDelete, "/products/invalid", nil, adminUser)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("when unauthenticated", func() {
		It("should return 401 Unauthorized", func() {
			req := newAuthenticatedRequest(http.MethodDelete, "/products/1", nil, nil)
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})
