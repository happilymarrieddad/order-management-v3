package users_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/users"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Users Handler", func() {
	Context("with a valid request", func() {
		It("should find users and return them with a total count", func() {
			opts := repos.UserFindOpts{Limit: 10, Offset: 0, Emails: []string{"find@example.com"}}
			body, _ := json.Marshal(opts)

			foundUsers := []*types.User{{ID: 1, Email: "find@example.com"}}
			totalCount := int64(1)

			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(foundUsers, totalCount, nil)

			req := createRequestWithRepo("POST", "/api/v1/users/find", body, nil)
			users.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)

			Expect(response).To(HaveKeyWithValue("total", float64(1))) // JSON numbers are float64
			Expect(response["data"]).To(HaveLen(1))
		})

		It("should return an empty data slice when no users are found", func() {
			opts := repos.UserFindOpts{Emails: []string{"notfound@example.com"}}
			body, _ := json.Marshal(opts)

			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return([]*types.User{}, int64(0), nil)

			req := createRequestWithRepo("POST", "/api/v1/users/find", body, nil)
			users.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)
			Expect(response).To(HaveKeyWithValue("total", float64(0)))
			Expect(response["data"]).To(BeEmpty())
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a malformed JSON body", func() {
			body := []byte(`{"emails": ["bad@json.com"]`) // Malformed JSON
			req := createRequestWithRepo("POST", "/api/v1/users/find", body, nil)
			users.Find(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 if the repository fails", func() {
			opts := repos.UserFindOpts{}
			body, _ := json.Marshal(opts)
			dbErr := errors.New("find database error")

			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			req := createRequestWithRepo("POST", "/api/v1/users/find", body, nil)
			users.Find(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
