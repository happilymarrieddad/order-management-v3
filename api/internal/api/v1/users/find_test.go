package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Find Users Handler", func() {
	Context("when users exist", func() {
		It("should return a list of users with default pagination", func() {
			foundUsers := []*types.User{
				{ID: 1, Email: "a@b.com"},
				{ID: 2, Email: "c@d.com"},
			}
			expectedOpts := &repos.UserFindOpts{Limit: 10, Offset: 0}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return(foundUsers, int64(2), nil)

			req := newAuthenticatedRequest("POST", "/users/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(2)))
		})

		It("should return a list of users with custom pagination", func() {
			foundUsers := []*types.User{{ID: 3, Email: "e@f.com"}}
			opts := &repos.UserFindOpts{Limit: 5, Offset: 5}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Eq(opts)).Return(foundUsers, int64(1), nil)

			body, _ := json.Marshal(opts)
			req := newAuthenticatedRequest("POST", "/users/find", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(1)))
		})
	})

	Context("when no users exist", func() {
		It("should return an empty list", func() {
			opts := &repos.UserFindOpts{Emails: []string{"notfound@example.com"}}
			body, _ := json.Marshal(opts)

			// The handler should still apply default limits even for an empty result.
			expectedOpts := &repos.UserFindOpts{Emails: []string{"notfound@example.com"}, Limit: 10}
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Eq(expectedOpts)).Return([]*types.User{}, int64(0), nil)

			req := newAuthenticatedRequest("POST", "/users/find", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var result types.FindResult
			err := json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(int64(0)))
			Expect(result.Data).To(BeEmpty())
		})
	})

	Context("with an invalid request", func() {
		It("should return 400 for a malformed JSON body", func() {
			body := []byte(`{"emails": ["bad@json.com"]`) // Malformed JSON
			req := newAuthenticatedRequest("POST", "/users/find", bytes.NewBuffer(body), adminUser)
			router.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the repository encounters an error", func() {
		It("should return 500 if the repository fails", func() {
			dbErr := errors.New("find database error")
			mockUsersRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, int64(0), dbErr)

			req := newAuthenticatedRequest("POST", "/users/find", bytes.NewBufferString(`{}`), adminUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when the user is not an admin", func() {
		It("should return 403 Forbidden for a non-admin user", func() {
			req := newAuthenticatedRequest("POST", "/users/find", bytes.NewBufferString(`{}`), basicUser)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
			Expect(rr.Body.String()).To(ContainSubstring("forbidden"))
		})

		It("should return 401 Unauthorized for an unauthenticated user", func() {
			req := newAuthenticatedRequest("POST", "/users/find", bytes.NewBufferString(`{}`), nil)
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
