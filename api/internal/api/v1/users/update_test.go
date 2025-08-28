package users_test

import (
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
	var updatePayload map[string]interface{}

	BeforeEach(func() {
		updatePayload = map[string]interface{}{
			"first_name": "Jane",
			"last_name":  "Doe",
			"address_id": int64(2),
		}
	})

	Context("with a valid request", func() {
		It("should update a user successfully", func() {
			userID := int64(1)
			body, _ := json.Marshal(updatePayload)
			addressID := updatePayload["address_id"].(int64)

			// Mock the Get call to find the existing user
			existingUser := &types.User{ID: userID, FirstName: "John"}
			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(existingUser, true, nil)

			// Mock the address existence check
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(&types.Address{ID: addressID}, true, nil)

			// Mock the Update call
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, user *types.User) error {
					Expect(user.ID).To(Equal(userID))
					Expect(user.FirstName).To(Equal(updatePayload["first_name"]))
					return nil
				},
			)

			req := createRequestWithRepo("PUT", "/api/v1/users/1", body, map[string]string{"id": "1"})
			users.Update(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var returnedUser types.User
			err := json.Unmarshal(rr.Body.Bytes(), &returnedUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedUser.FirstName).To(Equal("Jane"))
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a missing required field", func() {
			delete(updatePayload, "first_name")
			body, _ := json.Marshal(updatePayload)
			req := createRequestWithRepo("PUT", "/api/v1/users/1", body, map[string]string{"id": "1"})
			users.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when a dependency is not found", func() {
		It("should return 404 if the user to update is not found", func() {
			userID := int64(404)
			body, _ := json.Marshal(updatePayload)
			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(nil, false, nil)

			req := createRequestWithRepo("PUT", "/api/v1/users/404", body, map[string]string{"id": "404"})
			users.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 400 if the new address is not found", func() {
			userID := int64(1)
			body, _ := json.Marshal(updatePayload)
			addressID := updatePayload["address_id"].(int64)

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(&types.User{ID: userID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, nil)

			req := createRequestWithRepo("PUT", "/api/v1/users/1", body, map[string]string{"id": "1"})
			users.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 on final update error", func() {
			userID := int64(1)
			body, _ := json.Marshal(updatePayload)
			addressID := updatePayload["address_id"].(int64)
			dbErr := errors.New("update failed")

			mockUsersRepo.EXPECT().Get(gomock.Any(), userID).Return(&types.User{ID: userID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(&types.Address{ID: addressID}, true, nil)
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)

			req := createRequestWithRepo("PUT", "/api/v1/users/1", body, map[string]string{"id": "1"})
			users.Update(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
