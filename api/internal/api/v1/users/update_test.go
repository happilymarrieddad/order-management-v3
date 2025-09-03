package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Update User Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		payload    users.UpdateUserPayload
		targetUser *types.User
		address    *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		targetUser = &types.User{ID: 1, CompanyID: normalUser.CompanyID, FirstName: "Old", LastName: "User", AddressID: 10}
		address = &types.Address{ID: 20, Line1: "New Address"}

		payload = users.UpdateUserPayload{
			FirstName: "Updated",
			LastName:  "User",
			AddressID: address.ID,
		}
	})

	Context("Happy Path", func() {
		It("should update a user successfully for a normal user updating themselves", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, u *types.User) error {
				Expect(u.FirstName).To(Equal(payload.FirstName))
				Expect(u.LastName).To(Equal(payload.LastName))
				Expect(u.AddressID).To(Equal(payload.AddressID))
				return nil
			})

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var result types.User
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.FirstName).To(Equal(payload.FirstName))
			Expect(result.Password).To(BeEmpty())
		})

		It("should update a user successfully for an admin updating any user", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), adminUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), nil)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if a normal user tries to update another user", func() {
			otherUser := &types.User{ID: 99, CompanyID: normalUser.CompanyID}
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, otherUser.ID).Return(otherUser, true, nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/99", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with an invalid user ID", func() {
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/invalid-id", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should fail with a malformed JSON body", func() {
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer([]byte(`{"first_name":`)), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if a required field is missing", func() {
			payload.FirstName = ""
			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should return 404 if the user to update is not found", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(nil, false, nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on get user db error", func() {
			dbErr := errors.New("db error")
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(nil, false, dbErr)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 400 if the address does not exist", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, nil)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 on get address db error", func() {
			dbErr := errors.New("db error")
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, dbErr)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on update user db error", func() {
			mockUsersRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, targetUser.ID).Return(targetUser, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			dbErr := errors.New("db error")
			mockUsersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())
			req := newAuthenticatedRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})