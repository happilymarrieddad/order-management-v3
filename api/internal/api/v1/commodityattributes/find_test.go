package commodityattributes_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Commodity Attributes Handler", func() {
	Context("when commodity attributes exist", func() {
		It("should return a list of commodity attributes with default pagination", func() {
			foundCommodityAttributes := []*types.CommodityAttribute{
				{ID: 1, Name: "Color", CommodityType: types.CommodityTypeUnknown},
				{ID: 2, Name: "Size", CommodityType: types.CommodityTypeUnknown},
			}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.CommodityAttributeFindOpts{Limit: 10, Offset: 0})).Return(foundCommodityAttributes, int64(2), nil)

			// The endpoint is a POST to /find with an empty body for defaults.
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", []byte(`{}`), nil, 1)
			commodityattributes.Find(rr, req)

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
			foundCommodityAttributes := []*types.CommodityAttribute{{ID: 3, Name: "Weight", CommodityType: types.CommodityTypeUnknown}}
			opts := &repos.CommodityAttributeFindOpts{Limit: 5, Offset: 5}
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(opts)).Return(foundCommodityAttributes, int64(1), nil)

			// Send the pagination options in the request body.
			body, _ := json.Marshal(opts)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", body, nil, 1) // Added userID: 1
			commodityattributes.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
		})
	})

	Context("when no commodity attributes exist", func() {
		It("should return an empty list", func() {
			// The handler should still apply default limits even for an empty result.
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.CommodityAttributeFindOpts{Limit: 10, Offset: 0})).Return([]*types.CommodityAttribute{}, int64(0), nil)

			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", []byte(`{}`), nil, 1) // Added userID: 1
			commodityattributes.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			dbErr := errors.New("find query failed")
			mockCommodityAttributesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", []byte(`{}`), nil, 1) // Added userID: 1
			commodityattributes.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})

	Context("Unauthorized/Forbidden access", func() {
		// Dummy handler to be wrapped by the middleware for testing purposes
		dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		It("should return 401 if user ID is not found in context", func() {
			body, _ := json.Marshal(repos.CommodityAttributeFindOpts{})
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", body, nil) // No userID passed

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user ID not found in context"))
		})

		It("should return 401 if user is not found in repository", func() {
			var nonExistentUserID int64 = 999
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonExistentUserID).Return(nil, false, nil)

			body, _ := json.Marshal(repos.CommodityAttributeFindOpts{})
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", body, nil, nonExistentUserID)

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user not found"))
		})

		It("should return 500 if repository returns an error when getting user", func() {
			var someUserID int64 = 123
			dbErr := errors.New("user repo error")
			mockUsersRepo.EXPECT().Get(gomock.Any(), someUserID).Return(nil, false, dbErr)

			body, _ := json.Marshal(repos.CommodityAttributeFindOpts{})
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", body, nil, someUserID)

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("failed to get user information"))
		})

		It("should return 403 if user does not have admin role", func() {
			nonAdminUser := &types.User{ID: 2, Roles: types.Roles{types.RoleUser}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonAdminUser.ID).Return(nonAdminUser, true, nil)

			body, _ := json.Marshal(repos.CommodityAttributeFindOpts{})
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes/find", body, nil, nonAdminUser.ID)

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("only administrators can perform this action"))
		})
	})
})
