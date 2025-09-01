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

var _ = Describe("Create User Handler", func() {
	var (
		payload users.CreateUserPayload
		body    []byte
		err     error
	)

	BeforeEach(func() {
		payload = users.CreateUserPayload{
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Test",
			LastName:        "User",
			CompanyID:       int64(1),
			AddressID:       int64(2),
		}
		body, err = json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("with a valid request", func() {
		It("should create a user successfully", func() {

			// Expectations must be in the order they are called in the handler.
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(&types.Company{ID: payload.CompanyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(&types.Address{ID: payload.AddressID}, true, nil)
			mockUsersRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, user *types.User) error {
					Expect(user.Email).To(Equal(payload.Email))
					Expect(user.FirstName).To(Equal(payload.FirstName))
					Expect(user.Password).To(Equal(payload.Password))
					user.ID = 999 // Simulate DB assigning an ID
					return nil
				},
			)

			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusCreated))

			var returnedUser types.User
			err := json.Unmarshal(rr.Body.Bytes(), &returnedUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedUser.ID).To(Equal(int64(999)))
			Expect(returnedUser.Email).To(Equal("test@example.com"))
			Expect(returnedUser.Password).To(BeEmpty(), "Password should never be in the response")
		})
	})

	Context("with an invalid request body", func() {
		It("should return 400 for a malformed JSON body", func() {
			body := []byte(`{"email": "bad json",`)
			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid request body"))
		})

		It("should return 400 for a missing required field (e.g., email)", func() {
			payload.Email = ""
			body, _ := json.Marshal(payload)
			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'email' is required."))
		})

		It("should return 400 for mismatched passwords", func() {
			payload.ConfirmPassword = "differentpassword"
			body, _ := json.Marshal(payload)
			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("Field 'confirmpassword' has an invalid value."))
		})
	})

	Context("when a dependency check fails", func() {
		It("should return 400 if the user email already exists", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(&types.User{Email: payload.Email}, true, nil)

			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("user with that email already exists"))
		})

		It("should return 500 if GetByEmail fails", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, errors.New("db error"))

			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 400 if the company does not exist", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(nil, false, nil)

			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("company not found"))
		})

		It("should return 500 if checking company existence fails", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(nil, false, errors.New("company db error"))

			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 400 if the address does not exist", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(&types.Company{ID: payload.CompanyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, nil)

			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})

		It("should return 500 if checking address existence fails", func() {
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), payload.CompanyID).Return(&types.Company{ID: payload.CompanyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, errors.New("address db error"))

			req := newAuthenticatedRequest("POST", "/users", bytes.NewBuffer(body), basicUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
