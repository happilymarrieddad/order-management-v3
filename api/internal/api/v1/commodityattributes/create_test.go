package commodityattributes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodityattributes"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Create Commodity Attribute Handler", func() {
	var (
		payload      commodityattributes.CreateCommodityAttributePayload
		payloadBytes []byte
		newAttribute *types.CommodityAttribute
		err          error
	)

	BeforeEach(func() {
		payload = commodityattributes.CreateCommodityAttributePayload{
			Name:          "Test Attribute",
			CommodityType: types.CommodityTypeProduce,
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		newAttribute = &types.CommodityAttribute{
			ID:            1,
			Name:          payload.Name,
			CommodityType: payload.CommodityType,
		}
	})

	Context("when creation is successful", func() {
		It("should return 201 Created with the new attribute", func() {
			// Simulate the repo's Create method populating the ID of the passed-in struct.
			mockCommodityAttributesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Do(
				func(ctx context.Context, attr *types.CommodityAttribute) {
					attr.ID = newAttribute.ID
				}).Return(nil)

			req := newAuthenticatedRequest("POST", "/commodity-attributes", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedAttr types.CommodityAttribute
			err := json.Unmarshal(rr.Body.Bytes(), &returnedAttr)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedAttr.ID).To(Equal(newAttribute.ID))
			Expect(returnedAttr.Name).To(Equal(newAttribute.Name))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a malformed JSON body", func() {
			req := newAuthenticatedRequest("POST", "/commodity-attributes", bytes.NewBufferString(`{]`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a validation error (missing name)", func() {
			payload.Name = "" // Make the payload invalid
			body, _ := json.Marshal(payload)
			req := newAuthenticatedRequest("POST", "/commodity-attributes", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'name' is required."))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 Internal Server Error", func() {
			mockCommodityAttributesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db insert failed"))
			req := newAuthenticatedRequest("POST", "/commodity-attributes", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to create commodity attribute"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("POST", "/commodity-attributes", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("POST", "/commodity-attributes", bytes.NewBuffer(payloadBytes), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
