package locations_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/v1/locations"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

var _ = Describe("Delete Location Endpoint", func() {
	var (
		rec        *httptest.ResponseRecorder
		locationID int64
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		locationID = 1
	})

	performRequest := func(locID int64, user *types.User) {
		req := newAuthenticatedRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locID, 10), nil, user)
		router.ServeHTTP(rec, req)
	}

	Context("Router-level Tests", func() {
		It("should delete a location successfully for an admin", func() {
			// Admin path doesn't do a Get first, but the non-admin path does.
			// Since the route is admin-only, we only test the admin path.
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(nil)

			performRequest(locationID, adminUser)

			Expect(rec.Code).To(Equal(http.StatusNoContent))
		})

		It("should fail if not authenticated", func() {
			performRequest(locationID, nil)
			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail if a non-admin tries to delete a location", func() {
			performRequest(locationID, normalUser)
			Expect(rec.Code).To(Equal(http.StatusForbidden))
		})

		It("should fail with an invalid ID", func() {
			req := newAuthenticatedRequest(http.MethodDelete, "/locations/invalid-id", nil, adminUser)
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 404 if the location is not found on the Delete for admin", func() {
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(types.NewNotFoundError("location"))

			performRequest(locationID, adminUser)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on a database error during delete", func() {
			dbErr := errors.New("db error")
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(dbErr)
			performRequest(locationID, adminUser)
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("Direct Handler Tests", func() {
		var w *httptest.ResponseRecorder
		var r *http.Request

		BeforeEach(func() {
			w = httptest.NewRecorder()
		})

		It("should delete a location successfully for an admin (direct call)", func() {
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(nil)

			// Manually construct request and context
			vars := map[string]string{"id": strconv.FormatInt(locationID, 10)}
			r = httptest.NewRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locationID, 10), nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, adminUser)
			r = r.WithContext(ctxWithAuth)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusNoContent))
		})

		It("should return 400 for invalid location ID format (direct call)", func() {
			vars := map[string]string{"id": "invalid-id"}
			r = httptest.NewRequest(http.MethodDelete, "/locations/invalid-id", nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, adminUser)
			r = r.WithContext(ctxWithAuth)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 401 if not authenticated (direct call)", func() {
			vars := map[string]string{"id": strconv.FormatInt(locationID, 10)}
			r = httptest.NewRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locationID, 10), nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			// No auth user in context
			r = r.WithContext(ctxWithRepo)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 403 if a normal user tries to delete another company's location (direct call)", func() {
			otherLocation := &types.Location{ID: 99, CompanyID: 999} // Different company
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(otherLocation, true, nil)

			vars := map[string]string{"id": strconv.FormatInt(locationID, 10)}
			r = httptest.NewRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locationID, 10), nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, normalUser)
			r = r.WithContext(ctxWithAuth)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusForbidden))
		})

		It("should return 404 if location not found during ownership check (direct call)", func() {
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(nil, false, nil)

			vars := map[string]string{"id": strconv.FormatInt(locationID, 10)}
			r = httptest.NewRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locationID, 10), nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, normalUser)
			r = r.WithContext(ctxWithAuth)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on database error during ownership check (direct call)", func() {
			dbErr := errors.New("db error")
			mockLocationsRepo.EXPECT().Get(gomock.Any(), normalUser.CompanyID, locationID).Return(nil, false, dbErr)

			vars := map[string]string{"id": strconv.FormatInt(locationID, 10)}
			r = httptest.NewRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locationID, 10), nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, normalUser)
			r = r.WithContext(ctxWithAuth)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return 404 if location not found on final delete (direct call)", func() {
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(types.NewNotFoundError("location"))
			// No Get expectation for admin path

			vars := map[string]string{"id": strconv.FormatInt(locationID, 10)}
			r = httptest.NewRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locationID, 10), nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, adminUser)
			r = r.WithContext(ctxWithAuth)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 500 on database error on final delete (direct call)", func() {
			dbErr := errors.New("db error")
			mockLocationsRepo.EXPECT().Delete(gomock.Any(), locationID).Return(dbErr)
			// No Get expectation for admin path

			vars := map[string]string{"id": strconv.FormatInt(locationID, 10)}
			r = httptest.NewRequest(http.MethodDelete, "/locations/"+strconv.FormatInt(locationID, 10), nil)
			r = mux.SetURLVars(r, vars)
			ctxWithRepo := context.WithValue(r.Context(), middleware.RepoKey, mockGlobalRepo)
			ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, adminUser)
			r = r.WithContext(ctxWithAuth)

			locations.Delete(w, r)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})