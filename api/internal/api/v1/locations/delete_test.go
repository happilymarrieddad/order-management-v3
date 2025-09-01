package locations_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Delete Location Handler", func() {
	Context("when deletion is successful", func() {
		It("should return 204 No Content for an admin user", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, true, nil)
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

			req := newAuthenticatedRequest("DELETE", "/locations/1", nil, adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("when the location does not exist", func() {
		It("should return 404 Not Found", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(999)).Return(nil, false, nil)

			req := newAuthenticatedRequest("DELETE", "/locations/999", nil, adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository fails", func() {
		It("should return 500 on delete failure", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, true, nil)
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(errors.New("db delete failed"))

			req := newAuthenticatedRequest("DELETE", "/locations/1", nil, adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to delete location"))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("DELETE", "/locations/1", nil, basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("DELETE", "/locations/1", nil, nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
