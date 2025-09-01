package users_test

import (
	"encoding/json"
	"errors"
	"net/http"

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

			req := newAuthenticatedRequest("GET", "/users/123", nil, basicUser)
			router.ServeHTTP(rr, req)

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

			req := newAuthenticatedRequest("GET", "/users/404", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			Expect(rr.Body.String()).To(ContainSubstring("user not found"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 404 for a non-integer ID", func() {
			req := newAuthenticatedRequest("GET", "/users/abc", nil, basicUser)
			router.ServeHTTP(rr, req)

			// This is a router-level 404 because the route `/{id:[0-9]+}` does not match
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 for a generic database error", func() {
			userID := int64(500)
			dbErr := errors.New("database connection lost")

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(nil, false, dbErr)

			req := newAuthenticatedRequest("GET", "/users/500", nil, basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to get user"))
		})
	})

})
