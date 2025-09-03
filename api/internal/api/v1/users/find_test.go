package users_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Find Users Endpoint", func() {
	var (
		rec *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
	})

	performRequest := func(queryParams url.Values, user *types.User) {
		req := newAuthenticatedRequest(http.MethodGet, "/users/find?"+queryParams.Encode(), nil, user)
		router.ServeHTTP(rec, req)
	}

	Context("Happy Path", func() {
		It("should find users successfully for an admin", func() {
			expectedUsers := []*types.User{
				{ID: 1, CompanyID: company.ID, Email: "user1@example.com"},
				{ID: 2, CompanyID: company.ID, Email: "user2@example.com"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(expectedUsers, int64(len(expectedUsers)), nil)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.User]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", len(expectedUsers)))
			Expect(result.Data).To(HaveLen(len(expectedUsers)))
			Expect(result.Data[0].Email).To(Equal(expectedUsers[0].Email))
		})

		It("should find users successfully for a normal user (limited to their company)", func() {
			expectedUsers := []*types.User{
				{ID: 1, CompanyID: normalUser.CompanyID, Email: "user1@example.com"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.UserFindOpts) ([]*types.User, int64, error) {
				Expect(opts.CompanyID).To(Equal(normalUser.CompanyID))
				return expectedUsers, int64(len(expectedUsers)), nil
			})

			performRequest(url.Values{}, normalUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.User]
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Total).To(BeNumerically("==", len(expectedUsers)))
			Expect(result.Data).To(HaveLen(len(expectedUsers)))
		})

		It("should apply limit and offset", func() {
			expectedUsers := []*types.User{
				{ID: 1, CompanyID: company.ID, Email: "filtered@example.com"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.UserFindOpts) ([]*types.User, int64, error) {
				Expect(opts.Emails).To(ContainElement("filtered@example.com"))
				Expect(opts.Limit).To(Equal(5))
				Expect(opts.Offset).To(Equal(10))
				return expectedUsers, int64(len(expectedUsers)), nil
			})

			params := url.Values{}
			params.Add("email", "filtered@example.com")
			params.Add("limit", "5")
			params.Add("offset", "10")
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should filter by company_id for admin users", func() {
			otherCompanyID := int64(99)
			expectedUsers := []*types.User{
				{ID: 3, CompanyID: otherCompanyID, Email: "user3@example.com"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.UserFindOpts) ([]*types.User, int64, error) {
				Expect(opts.CompanyID).To(Equal(otherCompanyID))
				return expectedUsers, int64(len(expectedUsers)), nil
			})

			params := url.Values{}
			params.Add("company_id", strconv.FormatInt(otherCompanyID, 10))
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should filter by multiple IDs", func() {
			ids := []int64{1, 2}
			expectedUsers := []*types.User{
				{ID: 1, CompanyID: company.ID, Email: "user1@example.com"},
				{ID: 2, CompanyID: company.ID, Email: "user2@example.com"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.UserFindOpts) ([]*types.User, int64, error) {
				Expect(opts.IDs).To(ConsistOf(ids))
				return expectedUsers, int64(len(expectedUsers)), nil
			})

			params := url.Values{}
			for _, id := range ids {
				params.Add("id", strconv.FormatInt(id, 10))
			}
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should filter by multiple emails", func() {
			emails := []string{"user1@example.com", "user2@example.com"}
			expectedUsers := []*types.User{
				{ID: 1, CompanyID: company.ID, Email: "user1@example.com"},
				{ID: 2, CompanyID: company.ID, Email: "user2@example.com"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.UserFindOpts) ([]*types.User, int64, error) {
				Expect(opts.Emails).To(ConsistOf(emails))
				return expectedUsers, int64(len(expectedUsers)), nil
			})

			params := url.Values{}
			for _, email := range emails {
				params.Add("email", email)
			}
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should filter by multiple first names", func() {
			firstNames := []string{"John", "Jane"}
			expectedUsers := []*types.User{
				{ID: 1, CompanyID: company.ID, FirstName: "John"},
				{ID: 2, CompanyID: company.ID, FirstName: "Jane"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.UserFindOpts) ([]*types.User, int64, error) {
				Expect(opts.FirstNames).To(ConsistOf(firstNames))
				return expectedUsers, int64(len(expectedUsers)), nil
			})

			params := url.Values{}
			for _, name := range firstNames {
				params.Add("first_name", name)
			}
			performRequest(params, adminUser)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should filter by multiple last names", func() {
			lastNames := []string{"Doe", "Smith"}
			expectedUsers := []*types.User{
				{ID: 1, CompanyID: company.ID, LastName: "Doe"},
				{ID: 2, CompanyID: company.ID, LastName: "Smith"},
			}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, opts *repos.UserFindOpts) ([]*types.User, int64, error) {
				Expect(opts.LastNames).To(ConsistOf(lastNames))
				return expectedUsers, int64(len(expectedUsers)), nil
			})

			params := url.Values{}
			for _, name := range lastNames {
				params.Add("last_name", name)
			}
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
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			performRequest(url.Values{}, adminUser)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
