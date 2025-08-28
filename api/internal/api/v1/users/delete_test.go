package users_test

import (
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Delete User Handler", func() {
	Context("with a valid request", func() {
		It("should delete a user successfully and return 204", func() {
			mockUsersRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

			req := createRequestWithRepo("DELETE", "/api/v1/users/1", nil, map[string]string{"id": "1"})
			users.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})

		It("should return 204 even if the user to delete does not exist", func() {
			// The Delete operation is idempotent. If the resource is already gone, it's still a success.
			mockUsersRepo.EXPECT().Delete(gomock.Any(), int64(99)).Return(nil)

			req := createRequestWithRepo("DELETE", "/api/v1/users/99", nil, map[string]string{"id": "99"})
			users.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-numeric user ID", func() {
			req := createRequestWithRepo("DELETE", "/api/v1/users/invalid", nil, map[string]string{"id": "invalid"})
			users.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 if the repository fails to delete", func() {
			dbErr := errors.New("delete database error")
			mockUsersRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(dbErr)

			req := createRequestWithRepo("DELETE", "/api/v1/users/1", nil, map[string]string{"id": "1"})
			users.Delete(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
