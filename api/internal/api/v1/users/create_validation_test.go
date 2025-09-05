package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Create User Validation", func() {
	var (
		rec     *httptest.ResponseRecorder
		payload users.CreateUserPayload
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		payload = users.CreateUserPayload{
			Email:           "valid.email@example.com",
			Password:        "a-valid-password",
			ConfirmPassword: "a-valid-password",
			FirstName:       "ValidFirst",
			LastName:        "ValidLast",
			CompanyID:       1,
			AddressID:       1,
		}
	})

	DescribeTable("validation errors for user payload",
		func(mutator func(p *users.CreateUserPayload), expectedValidationField string, expectedBodyField string) {
			// For validation tests that are expected to fail before hitting the db,
			// we can satisfy the early mock expectations here.
			mockUsersRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, false, nil).AnyTimes()
			mockCompaniesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&types.Company{ID: 1}, true, nil).AnyTimes()
			mockAddressesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&types.Address{ID: 1}, true, nil).AnyTimes()

			mutator(&payload)

			// Explicitly validate the payload to see if it fails as expected
			validationErr := types.Validate(&payload)
			Expect(validationErr).To(HaveOccurred())
			Expect(validationErr.Error()).To(ContainSubstring(fmt.Sprintf("Field validation for '%s' failed", expectedValidationField)))

			body, err := json.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())

			req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
			// Optionally, check if the error message contains the expected field
			Expect(rec.Body.String()).To(ContainSubstring(fmt.Sprintf("Field '%s'", expectedBodyField)))
		},
				Entry("missing first_name", func(p *users.CreateUserPayload) { p.FirstName = "" }, "FirstName", "firstname"),
		Entry("short first_name", func(p *users.CreateUserPayload) { p.FirstName = "a" }, "FirstName", "firstname"),
		Entry("missing last_name", func(p *users.CreateUserPayload) { p.LastName = "" }, "LastName", "lastname"),
		Entry("short last_name", func(p *users.CreateUserPayload) { p.LastName = "b" }, "LastName", "lastname"),
		Entry("missing email", func(p *users.CreateUserPayload) { p.Email = "" }, "Email", "email"),
		Entry("invalid email", func(p *users.CreateUserPayload) { p.Email = "invalid-email" }, "Email", "email"),
		Entry("missing password", func(p *users.CreateUserPayload) { p.Password = "" }, "Password", "password"),
		Entry("short password", func(p *users.CreateUserPayload) { p.Password = "1234567" }, "Password", "password"),
		Entry("missing company_id", func(p *users.CreateUserPayload) { p.CompanyID = 0 }, "CompanyID", "companyid"),
		Entry("missing address_id", func(p *users.CreateUserPayload) { p.AddressID = 0 }, "AddressID", "addressid"),
	)

	It("should fail if the JSON keys are not in snake_case", func() {
		// Marshal a map with camelCase keys
		payloadMap := map[string]interface{}{
			"email":           payload.Email,
			"password":        payload.Password,
			"confirmPassword": payload.ConfirmPassword, // This key is camelCase
			"firstName":       payload.FirstName,
			"lastName":        payload.LastName,
			"companyId":       payload.CompanyID,
			"addressId":       payload.AddressID,
		}

		body, err := json.Marshal(payloadMap)
		Expect(err).NotTo(HaveOccurred())

		req := newAuthenticatedRequest(http.MethodPost, "/users", bytes.NewBuffer(body), normalUser)
		router.ServeHTTP(rec, req)

		Expect(rec.Code).To(Equal(http.StatusBadRequest))
		Expect(rec.Body.String()).To(ContainSubstring("confirmpassword")) // Expecting error about the snake_case field
	})
})