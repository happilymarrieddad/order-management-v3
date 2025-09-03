package users_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Delete User Endpoint", func() {
	var (
		rec *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
	})

	Context("Happy Path", func() {
		It("should delete a user successfully for an admin", func() {
			mockUsersRepo.EXPECT().Delete(gomock.Any(), int64(2)).Return(nil)

			req := newAuthenticatedRequest(http.MethodDelete, "/users/2", nil, adminUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("Error Paths", func() {
		It("should fail if not authenticated", func() {
			req := newAuthenticatedRequest(http.MethodDelete, "/users/2", nil, nil)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			req := newAuthenticatedRequest(http.MethodDelete, "/users/2", nil, normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})

		It("should fail with an invalid user ID", func() {
			req := newAuthenticatedRequest(http.MethodDelete, "/users/invalid-id", nil, adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error", func() {
			dbErr := errors.New("db error")
			mockUsersRepo.EXPECT().Delete(gomock.Any(), int64(2)).Return(dbErr)
			req := newAuthenticatedRequest(http.MethodDelete, "/users/2", nil, adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})