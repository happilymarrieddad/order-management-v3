package commodityattributes_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update Commodity Attribute Handler", func() {
	var updatePayload map[string]interface{}
	var existingAttribute *types.CommodityAttribute
	var adminUser *types.User

	BeforeEach(func() {
		adminUser = &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
		existingAttribute = &types.CommodityAttribute{
			ID:            123,
			Name:          "Original Name",
			CommodityType: types.CommodityTypeProduce,
		}
		updatePayload = map[string]interface{}{
			"name":          "Updated Name",
			"commodityType": types.CommodityTypeProduce,
		}
	})

	Context("with a valid request", func() {
		It("should update a commodity attribute successfully", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), existingAttribute.ID).Return(existingAttribute, true, nil)
			mockCommodityAttributesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, ca *types.CommodityAttribute) error {
					Expect(ca.ID).To(Equal(existingAttribute.ID))
					Expect(ca.Name).To(Equal(updatePayload["name"]))
					Expect(ca.CommodityType).To(Equal(updatePayload["commodityType"]))
					return nil
				},
			)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			var returnedAttribute types.CommodityAttribute
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAttribute)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAttribute.ID).To(Equal(existingAttribute.ID))
			Expect(returnedAttribute.Name).To(Equal(updatePayload["name"]))
			Expect(returnedAttribute.CommodityType).To(Equal(updatePayload["commodityType"]))
		})
	})

	Context("Unauthorized/Forbidden access", func() {
		// Dummy handler to be wrapped by the middleware for testing purposes
		dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			commodityattributes.Update(w, r)
		})

		It("should return 401 if user ID is not found in context", func() {
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/123", body, map[string]string{"id": "123"}) // No userID passed

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user ID not found in context"))
		})

		It("should return 401 if user is not found in repository", func() {
			var nonExistentUserID int64 = 999
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonExistentUserID).Return(nil, false, nil)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/123", body, map[string]string{"id": "123"}, nonExistentUserID)

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user not found"))
		})

		It("should return 500 if repository returns an error when getting user", func() {
			var someUserID int64 = 123
			dbErr := errors.New("user repo error")
			mockUsersRepo.EXPECT().Get(gomock.Any(), someUserID).Return(nil, false, dbErr)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/123", body, map[string]string{"id": "123"}, someUserID)

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("failed to get user information"))
		})

		It("should return 403 if user does not have admin role", func() {
			nonAdminUser := &types.User{ID: 2, Roles: types.Roles{types.RoleUser}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonAdminUser.ID).Return(nonAdminUser, true, nil)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/123", body, map[string]string{"id": "123"}, nonAdminUser.ID)

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("only administrators can perform this action"))
		})
	})

	Context("with an invalid request", func() {
		BeforeEach(func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()
		})

		It("should return 400 for a malformed JSON body", func() {
			body := []byte(`{"name": "bad json",`)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid request body"))
		})

		It("should return 400 for a missing required field (name)", func() {
			delete(updatePayload, "name")
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Key: 'UpdateCommodityAttributePayload.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
		})

		It("should return 400 for a missing required field (commodityType)", func() {
			delete(updatePayload, "commodityType")
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Key: 'UpdateCommodityAttributePayload.CommodityType' Error:Field validation for 'CommodityType' failed on the 'required' tag"))
		})

		It("should return 400 for a non-integer ID in path", func() {
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/abc", body, map[string]string{"id": "abc"}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid commodity attribute ID"))
		})
	})

	Context("when the commodity attribute does not exist", func() {
		It("should return 404 Not Found", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), existingAttribute.ID).Return(nil, false, nil)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("commodity attribute not found"))
		})
	})

	Context("when the repository encounters an error", func() {
		BeforeEach(func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()
		})

		It("should return 500 for a generic database error during Get", func() {
			dbErr := errors.New("database get error")
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), existingAttribute.ID).Return(nil, false, dbErr)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})

		It("should return 500 for a generic database error during Update", func() {
			dbErr := errors.New("database update error")
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), existingAttribute.ID).Return(existingAttribute, true, nil)
			mockCommodityAttributesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})

		It("should return 409 Conflict for a duplicate commodity attribute name", func() {
			mockCommodityAttributesRepo.EXPECT().Get(gomock.Any(), existingAttribute.ID).Return(existingAttribute, true, nil)
			uniqueConstraintErr := errors.New(`pq: duplicate key value violates unique constraint "commodity_attributes_name_key"`)
			mockCommodityAttributesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(uniqueConstraintErr)

			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/commodity-attributes/"+strconv.FormatInt(existingAttribute.ID, 10), body, map[string]string{"id": strconv.FormatInt(existingAttribute.ID, 10)}, adminUser.ID)
			commodityattributes.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusConflict))
			Expect(rr.Body.String()).To(ContainSubstring("Commodity attribute with this name already exists"))
		})
	})
})
