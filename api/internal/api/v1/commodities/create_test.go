package commodities_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodities"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Create Commodity Handler", func() {
	var (
		payload      commodities.CreateCommodityPayload
		payloadBytes []byte
		newCommodity *types.Commodity
		err          error
	)

	BeforeEach(func() {
		payload = commodities.CreateCommodityPayload{
			Name:          "Potatoes",
			CommodityType: types.CommodityTypeProduce,
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		newCommodity = &types.Commodity{
			ID:            1,
			Name:          payload.Name,
			CommodityType: payload.CommodityType,
		}
	})

	Context("when creation is successful", func() {
		It("should return 201 Created with the new commodity", func() {
			// Simulate the repo's Create method populating the ID of the passed-in struct.
			mockCommoditiesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Do(
				func(ctx context.Context, comm *types.Commodity) {
					comm.ID = newCommodity.ID
				}).Return(nil)

			req := newAuthenticatedRequest("POST", "/commodities", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedCommodity types.Commodity
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCommodity)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCommodity.ID).To(Equal(newCommodity.ID))
			Expect(returnedCommodity.Name).To(Equal(newCommodity.Name))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a malformed JSON body", func() {
			req := newAuthenticatedRequest("POST", "/commodities", bytes.NewBufferString(`{]`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 for a validation error (missing name)", func() {
			payload.Name = "" // Make the payload invalid
			body, _ := json.Marshal(payload)
			req := newAuthenticatedRequest("POST", "/commodities", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'name' is required."))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 Internal Server Error", func() {
			mockCommoditiesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db insert failed"))
			req := newAuthenticatedRequest("POST", "/commodities", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to create commodity"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden", func() {
			req := newAuthenticatedRequest("POST", "/commodities", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})
	})
})
