package commodities_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/commodities"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update Commodity Handler", func() {
	var (
		payload           commodities.UpdateCommodityPayload
		payloadBytes      []byte
		existingCommodity *types.Commodity
		err               error
	)

	BeforeEach(func() {
		payload = commodities.UpdateCommodityPayload{
			Name:          utils.Ref("Russet Potatoes"),
			CommodityType: utils.Ref(types.CommodityTypeProduce),
		}
		payloadBytes, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		existingCommodity = &types.Commodity{
			ID:            1,
			Name:          "Potatoes",
			CommodityType: types.CommodityTypeProduce,
		}
	})

	Context("when update is successful", func() {
		It("should return 200 OK with the updated commodity", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingCommodity, true, nil)
			mockCommoditiesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			req := newAuthenticatedRequest("PUT", "/commodities/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedCommodity types.Commodity
			err := json.Unmarshal(rr.Body.Bytes(), &returnedCommodity)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCommodity.ID).To(Equal(existingCommodity.ID))
			Expect(returnedCommodity.Name).To(Equal(*payload.Name))
			Expect(returnedCommodity.CommodityType).To(Equal(*payload.CommodityType))
		})
	})

	Context("when the commodity to update is not found", func() {
		It("should return 404 Not Found", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			req := newAuthenticatedRequest("PUT", "/commodities/999", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("with invalid input", func() {
		It("should return 400 for a malformed JSON body", func() {
			req := newAuthenticatedRequest("PUT", "/commodities/1", bytes.NewBufferString(`{]`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 404 for a non-integer ID", func() {
			req := newAuthenticatedRequest("PUT", "/commodities/abc", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

	})

	Context("when the repository fails", func() {
		It("should return 500 on update failure", func() {
			mockCommoditiesRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(existingCommodity, true, nil)
			mockCommoditiesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db update failed"))

			req := newAuthenticatedRequest("PUT", "/commodities/1", bytes.NewBuffer(payloadBytes), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to update commodity"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("PUT", "/commodities/1", bytes.NewBuffer(payloadBytes), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("PUT", "/commodities/1", bytes.NewBuffer(payloadBytes), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
