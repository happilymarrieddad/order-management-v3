package commodityattributes_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get Commodity Attribute Handler", func() {
	Context("when the commodity attribute exists", func() {
		It("should return the commodity attribute successfully", func() {
			commodityAttributeID := int64(123)
			expectedCommodityAttribute := &types.CommodityAttribute{ID: commodityAttributeID, Name: "Color", CommodityType: types.CommodityTypeUnknown}

			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), commodityAttributeID).Return(expectedCommodityAttribute, true, nil)

			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/123", nil, map[string]string{"id": "123"}, 1)
			commodityattributes.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCommodityAttribute types.CommodityAttribute
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCommodityAttribute)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCommodityAttribute.ID).To(Equal(commodityAttributeID))
			Expect(returnedCommodityAttribute.Name).To(Equal("Color"))
		})
	})

	Context("when the commodity attribute does not exist", func() {
		It("should return 404 Not Found", func() {
			commodityAttributeID := int64(404)
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), commodityAttributeID).Return(nil, false, nil)

			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/404", nil, map[string]string{"id": "404"}, 1) // Added userID: 1
			commodityattributes.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("commodity attribute not found"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/abc", nil, map[string]string{"id": "abc"}, 1) // Added userID: 1
			commodityattributes.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid commodity attribute ID"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			commodityAttributeID := int64(500)
			dbErr := errors.New("database connection lost")

			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), commodityAttributeID).Return(nil, false, dbErr)

			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/500", nil, map[string]string{"id": "500"}, 1) // Added userID: 1
			commodityattributes.Get(rr, req)

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
			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/123", nil, map[string]string{"id": "123"}) // No userID passed

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user ID not found in context"))
		})

		It("should return 401 if user is not found in repository", func() {
			var nonExistentUserID int64 = 999
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonExistentUserID).Return(nil, false, nil)

			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/123", nil, map[string]string{"id": "123"}, nonExistentUserID)

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user not found"))
		})

		It("should return 500 if repository returns an error when getting user", func() {
			var someUserID int64 = 123
			dbErr := errors.New("user repo error")
			mockUsersRepo.EXPECT().Get(gomock.Any(), someUserID).Return(nil, false, dbErr)

			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/123", nil, map[string]string{"id": "123"}, someUserID)

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("failed to get user information"))
		})

		It("should return 403 if user does not have admin role", func() {
			nonAdminUser := &types.User{ID: 2, Roles: types.Roles{types.RoleUser}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonAdminUser.ID).Return(nonAdminUser, true, nil)

			req := createRequestWithRepo("GET", "/api/v1/commodity-attributes/123", nil, map[string]string{"id": "123"}, nonAdminUser.ID)

			// Apply the middleware directly
			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("only administrators can perform this action"))
		})
	})
})