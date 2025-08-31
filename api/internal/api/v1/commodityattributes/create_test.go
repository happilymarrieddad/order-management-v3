package commodityattributes_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Create Commodity Attribute Handler", func() {
	var createPayload map[string]interface{}

	BeforeEach(func() {
		createPayload = map[string]interface{}{
			"name":          "Test Attribute",
			"commodityType": types.CommodityTypeProduce,
		}
	})

	Context("with a valid request", func() {
		It("should create a commodity attribute successfully", func() {
			adminUser := &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}

			body, _ := json.Marshal(createPayload)

			mockCommodityAttributesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, ca *types.CommodityAttribute) error {
					Expect(ca.Name).To(Equal(createPayload["name"]))
					Expect(ca.CommodityType).To(Equal(createPayload["commodityType"]))
					ca.ID = 123 // Simulate DB assigning an ID
					return nil
				},
			)

			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, adminUser.ID)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedAttribute types.CommodityAttribute
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAttribute)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAttribute.ID).To(Equal(int64(123)))
			Expect(returnedAttribute.Name).To(Equal("Test Attribute"))
		})
	})

	Context("Unauthorized/Forbidden access", func() {
		// Dummy handler to be wrapped by the middleware for testing purposes
		dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			commodityattributes.Create(w, r)
		})

		It("should return 401 if user ID is not found in context", func() {
			body, _ := json.Marshal(createPayload)
			// No user ID passed to createRequestWithRepo
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil)

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user ID not found in context"))
		})

		It("should return 401 if user is not found in repository", func() {
			var nonExistentUserID int64 = 999
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonExistentUserID).Return(nil, false, nil)

			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, nonExistentUserID)

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
			Expect(rr.Body.String()).To(ContainSubstring("user not found"))
		})

		It("should return 500 if repository returns an error when getting user", func() {
			var someUserID int64 = 123
			dbErr := errors.New("user repo error")
			mockUsersRepo.EXPECT().Get(gomock.Any(), someUserID).Return(nil, false, dbErr)

			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, someUserID)

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("failed to get user information"))
		})

		It("should return 403 if user does not have admin role", func() {
			nonAdminUser := &types.User{ID: 2, Roles: types.Roles{types.RoleUser}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), nonAdminUser.ID).Return(nonAdminUser, true, nil)

			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, nonAdminUser.ID)

			middleware.AuthUserAdminRequired(dummyHandler).ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("only administrators can perform this action"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a malformed JSON body", func() {
			// Mock the user lookup to return an admin user for this context
			adminUser := &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()

			body := []byte(`{"name": "bad json",`)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, adminUser.ID)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a missing required field (name)", func() {
			// Mock the user lookup to return an admin user for this context
			adminUser := &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()

			delete(createPayload, "name")
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, adminUser.ID)
			commodityattributes.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a missing required field (commodityType)", func() {
			// Mock the user lookup to return an admin user for this context
			adminUser := &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()

			delete(createPayload, "commodityType")
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, adminUser.ID)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("with invalid field values", func() {
		BeforeEach(func() {
			// Mock the user lookup to return an admin user for this context
			adminUser := &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()
		})

		It("should return 400 if the name is too short", func() {
			createPayload["name"] = "a"
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, 1)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("failed on the 'min' tag"))
		})

		It("should return 400 if the name is too long", func() {
			createPayload["name"] = strings.Repeat("a", 256)
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, 1)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("failed on the 'max' tag"))
		})

		It("should return 400 for an invalid commodityType value", func() {
			createPayload["commodityType"] = 999 // Invalid enum value
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, 1)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("failed on the 'oneof' tag"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			// Mock the user lookup to return an admin user for this context
			adminUser := &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()

			body, _ := json.Marshal(createPayload)
			dbErr := errors.New("unexpected database error")

			mockCommodityAttributesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, adminUser.ID)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})

		It("should return 409 Conflict for a duplicate commodity attribute name", func() {
			// Mock the user lookup to return an admin user for this context
			adminUser := &types.User{ID: 1, Roles: types.Roles{types.RoleAdmin}}
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.ID).Return(adminUser, true, nil).AnyTimes()

			body, _ := json.Marshal(createPayload)
			// Simulate a unique constraint violation error
			uniqueConstraintErr := errors.New(`pq: duplicate key value violates unique constraint "commodity_attributes_name_key"`)

			mockCommodityAttributesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uniqueConstraintErr)

			req := createRequestWithRepo("POST", "/api/v1/commodity-attributes", body, nil, adminUser.ID)
			commodityattributes.Create(rr, req)

			Expect(rr.Code).To(Equal(http.StatusConflict))
			Expect(rr.Body.String()).To(ContainSubstring("Commodity attribute with this name already exists"))
		})
	})
})
