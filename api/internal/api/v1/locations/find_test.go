package locations_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Locations Handler", func() {
	Context("when locations exist", func() {
		It("should return a list of locations", func() {
			findOpts := &repos.LocationFindOpts{
				CompanyIDs: []int64{1},
				Limit:      10,
				Offset:     0,
			}
			body, _ := json.Marshal(findOpts)

			foundLocations := []*types.Location{
				{ID: 1, Name: "Location A", CompanyID: 1},
				{ID: 2, Name: "Location B", CompanyID: 1},
			}
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(findOpts)).Return(foundLocations, int64(2), nil)

			req := createRequestWithRepo("POST", "/api/v1/locations/find", body, nil)
			locations.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))

			dataBytes, _ := json.Marshal(result.Data)
			var returnedLocations []types.Location
			json.Unmarshal(dataBytes, &returnedLocations)
			Expect(returnedLocations).To(HaveLen(2))
			Expect(returnedLocations[0].Name).To(Equal("Location A"))
		})
	})

	Context("when no locations are found", func() {
		It("should return an empty list", func() {
			findOpts := &repos.LocationFindOpts{Names: []string{"Non-existent"}}
			body, _ := json.Marshal(findOpts)

			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Eq(findOpts)).Return([]*types.Location{}, int64(0), nil)

			req := createRequestWithRepo("POST", "/api/v1/locations/find", body, nil)
			locations.Find(rr, req)

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
			mockLocationsRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			req := createRequestWithRepo("POST", "/api/v1/locations/find", []byte(`{}`), nil)
			locations.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
