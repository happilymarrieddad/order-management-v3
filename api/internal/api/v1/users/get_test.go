package users_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get User Handler", func() {

	Context("when the user exists", func() {
		It("should return the user successfully", func() {
			userID := int64(123)
			expectedUser := &types.User{ID: userID, Email: "found@example.com"}

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(expectedUser, true, nil)

			req := createRequestWithRepo("GET", "/api/v1/users/123", nil, map[string]string{"id": "123"})
			users.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedUser types.User
			err := json.Unmarshal(rr.Body.Bytes(), &returnedUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedUser.ID).To(Equal(userID))
			Expect(returnedUser.Email).To(Equal("found@example.com"))
			Expect(returnedUser.Password).To(BeEmpty())
		})
	})

	Context("when the user does not exist", func() {
		It("should return 404 Not Found", func() {
			userID := int64(404)
			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(nil, false, nil)

			req := createRequestWithRepo("GET", "/api/v1/users/404", nil, map[string]string{"id": "404"})
			users.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("user not found"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a non-integer ID", func() {
			req := createRequestWithRepo("GET", "/api/v1/users/abc", nil, map[string]string{"id": "abc"})
			users.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid user ID"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			userID := int64(500)
			dbErr := errors.New("database connection lost")

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(nil, false, dbErr)

			req := createRequestWithRepo("GET", "/api/v1/users/500", nil, map[string]string{"id": "500"})
			users.Get(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
