package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update User Handler", func() {
	var (
		payload users.UpdateUserPayload
		body    []byte
		err     error
	)

	BeforeEach(func() {
		payload = users.UpdateUserPayload{
			FirstName: "Jane",
			LastName:  "Doe",
			AddressID: int64(2),
		}
		body, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("as an admin", func() {
		It("should update another user successfully", func() {
			userID := int64(2) // Use a different ID than the admin user to avoid mock collision

			// Mock the Get call to find the existing user
			existingUser := &types.User{ID: userID, FirstName: "John"}
			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(existingUser, true, nil)

			// Mock the address existence check
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(&types.Address{ID: payload.AddressID}, true, nil)

			// Mock the Update call
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, user *types.User) error {
					Expect(user.ID).To(Equal(userID))
					Expect(user.FirstName).To(Equal(payload.FirstName))
					return nil
				},
			)

			req := newAuthenticatedRequest("PUT", "/users/2", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedUser types.User
			err := json.Unmarshal(rr.Body.Bytes(), &returnedUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedUser.FirstName).To(Equal("Jane"))
		})
	})

	Context("as a non-admin user", func() {
		It("should update their own profile successfully", func() {
			// The basicUser (ID 2) is updating their own profile.
			userID := basicUser.ID

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(basicUser, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(&types.Address{ID: payload.AddressID}, true, nil)
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			req := newAuthenticatedRequest("PUT", "/users/2", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			var returnedUser types.User
			err := json.Unmarshal(rr.Body.Bytes(), &returnedUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedUser.FirstName).To(Equal(payload.FirstName))
		})

		It("should forbid updating another user's profile", func() {
			// The basicUser (ID 2) is trying to update another user (ID 3).
			targetUserID := int64(3)
			targetUser := &types.User{ID: targetUserID, CompanyID: 2} // Same company, different user

			// The handler will fetch the target user to check ownership
			mockUsersRepo.EXPECT().Get(gomock.Any(), targetUserID).Return(targetUser, true, nil)

			req := newAuthenticatedRequest("PUT", "/users/3", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("user not authorized to update this user"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a missing required field", func() {
			payload.FirstName = "" // Make the payload invalid
			body, _ := json.Marshal(payload)
			req := newAuthenticatedRequest("PUT", "/users/2", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'firstname' is required."))
		})

		It("should return 404 for a non-integer ID", func() {
			req := newAuthenticatedRequest("PUT", "/users/abc", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when a dependency is not found", func() {
		It("should return 404 if the user to update is not found", func() {
			userID := int64(404)
			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(nil, false, nil)

			req := newAuthenticatedRequest("PUT", "/users/404", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 400 if the new address is not found", func() {
			userID := int64(2)

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(&types.User{ID: userID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, nil)

			req := newAuthenticatedRequest("PUT", "/users/2", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 on final update error", func() {
			userID := int64(2)
			dbErr := errors.New("update failed")

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(&types.User{ID: userID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(&types.Address{ID: payload.AddressID}, true, nil)
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)

			req := newAuthenticatedRequest("PUT", "/users/2", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("unable to update user"))
		})
	})

})
