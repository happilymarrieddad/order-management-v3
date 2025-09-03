package locations_test

import (
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

var _ = Describe("Find Locations Endpoint", func() {
	var (
		rec  *httptest.ResponseRecorder
		loc1 *types.Location
		loc2 *types.Location
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		loc1 = &types.Location{ID: 1, Name: "Location A", CompanyID: company.ID, AddressID: 1}
		loc2 = &types.Location{ID: 2, Name: "Location B", CompanyID: company.ID, AddressID: 2}
	})

	performRequest := func(user *types.User, queryParams url.Values) {
		var body []byte
		if queryParams != nil {
			body = []byte(queryParams.Encode())
		}
		req := newAuthenticatedRequest(http.MethodGet, "/locations/find?"+queryParams.Encode(), body, user)
		router.ServeHTTP(rec, req)
	}

	Context("when the user is an admin", func() {
		It("should find locations within their own company", func() {
			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{adminUser.CompanyID},
				Limit:      10, // Default limit
				Offset:     0,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{loc1, loc2}, int64(2), nil)

			performRequest(adminUser, nil)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))
			Expect(result.Data).To(HaveLen(2))
			Expect(result.Data[0].ID).To(Equal(loc1.ID))
			Expect(result.Data[1].ID).To(Equal(loc2.ID))
		})

		It("should not find locations from other companies", func() {
			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{adminUser.CompanyID},
				Limit:      10,
				Offset:     0,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{}, int64(0), nil)

			performRequest(adminUser, nil)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})

		It("should apply limit and offset", func() {
			queryParams := url.Values{}
			queryParams.Set("limit", "1")
			queryParams.Set("offset", "1")

			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{adminUser.CompanyID},
				Limit:      1,
				Offset:     1,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{loc2}, int64(2), nil)

			performRequest(adminUser, queryParams)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(loc2.ID))
		})

		It("should filter by name", func() {
			queryParams := url.Values{}
			queryParams.Set("name", "Location A")

			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{adminUser.CompanyID},
				Names:      []string{"Location A"},
				Limit:      10,
				Offset:     0,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{loc1}, int64(1), nil)

			performRequest(adminUser, queryParams)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
			Expect(result.Data).To(HaveLen(1))
			Expect(result.Data[0].ID).To(Equal(loc1.ID))
		})
	})

	Context("when the user is a normal user", func() {
		It("should find locations within their own company", func() {
			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{normalUser.CompanyID},
				Limit:      10,
				Offset:     0,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{loc1, loc2}, int64(2), nil)

			performRequest(normalUser, nil)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))
			Expect(result.Data).To(HaveLen(2))
			Expect(result.Data[0].ID).To(Equal(loc1.ID))
			Expect(result.Data[1].ID).To(Equal(loc2.ID))
		})

		It("should not find locations from other companies", func() {
			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{normalUser.CompanyID},
				Limit:      10,
				Offset:     0,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{}, int64(0), nil)

			performRequest(normalUser, nil)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("Error Paths", func() {
		It("should return 401 if not authenticated", func() {
			performRequest(nil, nil) // No user
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 500 on a database error", func() {
			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{adminUser.CompanyID},
				Limit:      10,
				Offset:     0,
			}
			dbErr := errors.New("db error")
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return(nil, int64(0), dbErr)

			performRequest(adminUser, nil)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should handle invalid limit parameter gracefully", func() {
			queryParams := url.Values{}
			queryParams.Set("limit", "invalid")

			// Expect default limit to be used
			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{adminUser.CompanyID},
				Limit:      10,
				Offset:     0,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{loc1}, int64(1), nil)

			performRequest(adminUser, queryParams)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
			Expect(result.Data).To(HaveLen(1))
		})

		It("should handle invalid offset parameter gracefully", func() {
			queryParams := url.Values{}
			queryParams.Set("offset", "invalid")

			// Expect default offset to be used
			expectedOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{adminUser.CompanyID},
				Limit:      10,
				Offset:     0,
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.Location{loc1}, int64(1), nil)

			performRequest(adminUser, queryParams)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.FindResult[types.Location] // Changed to generic FindResult
			err := json.NewDecoder(rec.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
			Expect(result.Data).To(HaveLen(1))
		})
	})
})
