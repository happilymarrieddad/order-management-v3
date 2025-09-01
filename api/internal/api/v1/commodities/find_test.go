package commodities_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Commodities Handler", func() {
	Context("when commodities exist", func() {
		It("should return a list of commodities for an admin user", func() {
			foundCommodities := []*types.Commodity{
				{ID: 1, Name: "Potatoes"},
				{ID: 2, Name: "Apples"},
			}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.FindCommoditiesOpts{Limit: 10, Offset: 0})).Return(foundCommodities, int64(2), nil)

			req := newAuthenticatedRequest("POST", "/commodities/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))

			dataBytes, _ := json.Marshal(result.Data)
			var returnedCommodities []types.Commodity
			json.Unmarshal(dataBytes, &returnedCommodities)
			Expect(returnedCommodities).To(HaveLen(2))
			Expect(returnedCommodities[0].Name).To(Equal("Potatoes"))
		})

		It("should return a list of commodities with custom pagination", func() {
			foundCommodities := []*types.Commodity{{ID: 3, Name: "Carrots"}}
			opts := &repos.FindCommoditiesOpts{Limit: 5, Offset: 5}
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(opts)).Return(foundCommodities, int64(1), nil)

			body, _ := json.Marshal(opts)
			req := newAuthenticatedRequest("POST", "/commodities/find", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
		})
	})

	Context("when no commodities exist", func() {
		It("should return an empty list for an admin user", func() {
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Eq(&repos.FindCommoditiesOpts{Limit: 10, Offset: 0})).Return([]*types.Commodity{}, int64(0), nil)

			req := newAuthenticatedRequest("POST", "/commodities/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 Internal Server Error", func() {
			mockCommoditiesRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), errors.New("find query failed"))
			req := newAuthenticatedRequest("POST", "/commodities/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to find commodities"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("POST", "/commodities/find", bytes.NewBufferString(`{}`), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("POST", "/commodities/find", bytes.NewBufferString(`{}`), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
