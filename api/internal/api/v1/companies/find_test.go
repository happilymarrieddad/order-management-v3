package companies_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Find Companies Endpoint", func() {
	var (
		rec *httptest.ResponseRecorder
		comp1 *types.Company
		comp2 *types.Company
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		comp1 = &types.Company{ID: 1, Name: "Company A", AddressID: 1}
		comp2 = &types.Company{ID: 2, Name: "Company B", AddressID: 2}
	})

	performRequest := func(queryParams url.Values, user *types.User) {
		var err error
		rec, err = testutils.PerformRequest(router, http.MethodGet, "/companies/find?"+queryParams.Encode(), url.Values{}, nil, user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should find companies successfully for an admin", func() {
			queryParams := url.Values{}
			expectedOpts := &repos.CompanyFindOpts{
				Limit:  10,
				Offset: 0,
			}
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Company{comp1, comp2}, int64(2), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Company]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 2))
			Expect(result.Data).To(HaveLen(2))
			Expect(result.Data[0].ID).To(Equal(comp1.ID))
		})

		It("should apply limit and offset", func() {
			queryParams := url.Values{}
			queryParams.Set("limit", "1")
			queryParams.Set("offset", "1")

			expectedOpts := &repos.CompanyFindOpts{
				Limit:  1,
				Offset: 1,
			}
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Company{comp2}, int64(2), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Company]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 2))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(comp2.ID))
		})

		It("should filter by name", func() {
			queryParams := url.Values{}
			queryParams.Set("name", "Company A")

			expectedOpts := &repos.CompanyFindOpts{
				Names: []string{"Company A"},
				Limit: 10,
				Offset: 0,
			}
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Company{comp1}, int64(1), nil)

			performRequest(queryParams, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Company]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", 1))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(comp1.ID))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(url.Values{}, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			performRequest(url.Values{}, normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Error Paths", func() {
		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockCompaniesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
