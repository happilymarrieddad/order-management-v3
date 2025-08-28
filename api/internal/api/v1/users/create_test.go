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

var _ = Describe("Create User Handler", func() {
	var createPayload map[string]interface{}

	BeforeEach(func() {
		createPayload = map[string]interface{}{
			"email":            "test@example.com",
			"password":         "password123",
			"confirm_password": "password123",
			"first_name":       "Test",
			"last_name":        "User",
			"company_id":       int64(1),
			"address_id":       int64(2),
		}
	})

	Context("with a valid request", func() {
		It("should create a user successfully", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			addressID := createPayload["address_id"].(int64)
			email := createPayload["email"].(string)

			// Expectations must be in the order they are called in the handler.
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(&types.Company{ID: companyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(&types.Address{ID: addressID}, true, nil)
			mockUsersRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, user *types.User) error {
					Expect(user.Email).To(Equal(createPayload["email"]))
					Expect(user.FirstName).To(Equal(createPayload["first_name"]))
					Expect(user.Password).To(Equal(createPayload["password"]))
					user.ID = 999 // Simulate DB assigning an ID
					return nil
				},
			)

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)

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
			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("invalid request body"))
		})

		It("should return 400 for a missing required field (e.g., email)", func() {
			delete(createPayload, "email")
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("validation failed"))
		})

		It("should return 400 for mismatched passwords", func() {
			createPayload["confirm_password"] = "differentpassword"
			body, _ := json.Marshal(createPayload)
			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("validation failed"))
		})
	})

	Context("when a dependency check fails", func() {
		It("should return 400 if the user email already exists", func() {
			body, _ := json.Marshal(createPayload)
			email := createPayload["email"].(string)
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(&types.User{Email: email}, true, nil)

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("user with that email already exists"))
		})

		It("should return 500 if GetByEmail fails", func() {
			body, _ := json.Marshal(createPayload)
			email := createPayload["email"].(string)
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, false, errors.New("db error"))

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 400 if the company does not exist", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			email := createPayload["email"].(string)

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(nil, false, nil)

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("company not found"))
		})

		It("should return 500 if checking company existence fails", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			email := createPayload["email"].(string)

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(nil, false, errors.New("company db error"))

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 400 if the address does not exist", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			addressID := createPayload["address_id"].(int64)
			email := createPayload["email"].(string)

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(&types.Company{ID: companyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, nil)

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring("address not found"))
		})

		It("should return 500 if checking address existence fails", func() {
			body, _ := json.Marshal(createPayload)
			companyID := createPayload["company_id"].(int64)
			addressID := createPayload["address_id"].(int64)
			email := createPayload["email"].(string)

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), companyID).Return(&types.Company{ID: companyID}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), addressID).Return(nil, false, errors.New("address db error"))

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when the final creation fails", func() {
		It("should return 500 for a generic database error on user creation", func() {
			body, _ := json.Marshal(createPayload)
			email := createPayload["email"].(string)
			dbErr := errors.New("unique constraint violation")

			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, false, nil)
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&types.Company{ID: 1}, true, nil)
			mockAddressesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&types.Address{ID: 2}, true, nil)
			mockUsersRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

			req := createRequestWithRepo("POST", "/api/v1/users", body, nil)
			users.Create(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring(dbErr.Error()))
		})
	})
})
