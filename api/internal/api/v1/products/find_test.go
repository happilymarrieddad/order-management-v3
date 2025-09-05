package products_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Products Endpoint", func() {
	var (
		rec *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
	})

	performRequest := func(queryParams url.Values, user *types.User) {
		req := newAuthenticatedRequest(http.MethodGet, "/products/find?"+queryParams.Encode(), nil, user)
		router.ServeHTTP(rec, req)
	}

	Context("Happy Path", func() {
		It("should find products successfully for an admin, scoped to their company", func() {
			expectedProducts := []*types.Product{
				{ID: 1, CompanyID: adminUser.CompanyID, Name: "Product 1"},
				{ID: 2, CompanyID: adminUser.CompanyID, Name: "Product 2"},
			}
			mockProductsRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.ProductFindOpts) ([]*types.Product, int64, error) {
				Expect(opts.CompanyID).To(Equal(adminUser.CompanyID))
				return expectedProducts, int64(len(expectedProducts)), nil
			})

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Product]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", len(expectedProducts)))
			Expect(result.Data).To(HaveLen(len(expectedProducts)))
			Expect(result.Data[0].Name).To(Equal(expectedProducts[0].Name))
		})

		It("should find products successfully for a normal user, scoped to their company", func() {
			expectedProducts := []*types.Product{
				{ID: 1, CompanyID: normalUser.CompanyID, Name: "Product 1"},
			}
			mockProductsRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.ProductFindOpts) ([]*types.Product, int64, error) {
				Expect(opts.CompanyID).To(Equal(normalUser.CompanyID))
				return expectedProducts, int64(len(expectedProducts)), nil
			})

			performRequest(url.Values{}, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Product]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", len(expectedProducts)))
			Expect(result.Data).To(HaveLen(len(expectedProducts)))
		})

		It("should apply limit and offset", func() {
			expectedProducts := []*types.Product{
				{ID: 1, CompanyID: company.ID, Name: "Filtered Product"},
			}
			mockProductsRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.ProductFindOpts) ([]*types.Product, int64, error) {
				Expect(opts.Name).To(Equal("Filtered Product"))
				Expect(opts.Limit).To(Equal(5))
				Expect(opts.Offset).To(Equal(10))
				return expectedProducts, int64(len(expectedProducts)), nil
			})

			params := url.Values{}
			params.Add("name", "Filtered Product")
			params.Add("limit", "5")
			params.Add("offset", "10")
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should filter by name", func() {
			productName := "Specific Product"
			expectedProducts := []*types.Product{
				{ID: 1, CompanyID: company.ID, Name: productName},
			}
			mockProductsRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.ProductFindOpts) ([]*types.Product, int64, error) {
				Expect(opts.Name).To(Equal(productName))
				return expectedProducts, int64(len(expectedProducts)), nil
			})

			params := url.Values{}
			params.Add("name", productName)
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should return StatusBadRequest for invalid limit parameter", func() {
			params := url.Values{}
			params.Add("limit", "invalid")
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return StatusBadRequest for invalid offset parameter", func() {
			params := url.Values{}
			params.Add("offset", "invalid")
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Error Paths", func() {
		It("should fail if not authenticated", func() {
			performRequest(url.Values{}, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockProductsRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})