package companies_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/testutils"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/companies"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Create Company Endpoint", func() {
	var (
		rec     *httptest.ResponseRecorder
		payload companies.CreateCompanyPayload
		address *types.Address
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		address = &types.Address{ID: 1, Line1: "123 Test St"}
		payload = companies.CreateCompanyPayload{
			Name:      "Test Company",
			AddressID: address.ID,
		}
	})

	performRequest := func(payload interface{}, user *types.User) {
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())
		rec, err = testutils.PerformRequest(router, http.MethodPost, "/companies", url.Values{}, bytes.NewBuffer(body), user, mockGlobalRepo)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Happy Path", func() {
		It("should create a company successfully for an admin", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockCompaniesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, c *types.Company) error {
				c.ID = 3 // Simulate ID generation
				return nil
			})

			performRequest(payload, adminUser)

			Expect(rec.Code).To(Equal(http.StatusCreated))
			var result types.Company
			Expect(json.NewDecoder(rec.Body).Decode(&result)).To(Succeed())
			Expect(result.Name).To(Equal(payload.Name))
			Expect(result.AddressID).To(Equal(payload.AddressID))
			Expect(result.ID).To(BeNumerically(">", 0))
		})
	})

	Context("Authorization and Authentication", func() {
		It("should fail if the user is not authenticated", func() {
			performRequest(payload, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if not an admin", func() {
			performRequest(payload, normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("Invalid Input", func() {
		It("should fail with a malformed JSON body", func() {
			rec, err := testutils.PerformRequest(router, http.MethodPost, "/companies", url.Values{}, bytes.NewBuffer([]byte(`{`)), adminUser, mockGlobalRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should fail if a required field is missing", func() {
			payload.Name = ""
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Dependency and Repository Errors", func() {
		It("should fail if the address does not exist", func() {
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, nil)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 on address validation db error", func() {
			dbErr := errors.New("db error")
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(nil, false, dbErr)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 500 on company creation db error", func() {
			dbErr := errors.New("db error")
			mockAddressesRepo.EXPECT().Get(gomock.Any(), payload.AddressID).Return(address, true, nil)
			mockCompaniesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			performRequest(payload, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})