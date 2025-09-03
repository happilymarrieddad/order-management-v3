package testutils

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	mock_repos "github.com/happilymarrieddad/order-management-v3/api/internal/repos/mocks"
)

// PerformRequest is a generic helper function to perform an HTTP request for API tests.
// It constructs the request, sets up authentication, and serves it to the router.
func PerformRequest(
	router *mux.Router,
	method string,
	path string,
	params url.Values,
	body io.Reader,
	user *types.User,
	mockGlobalRepo *mock_repos.MockGlobalRepo, // Add mockGlobalRepo as a parameter
) (*httptest.ResponseRecorder, error) {
	urlStr := path
	if len(params) > 0 {
		urlStr += "?" + params.Encode()
	}

	req, err := newAuthenticatedRequest(method, urlStr, body, user, mockGlobalRepo)
	if err != nil {
		return nil, err
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec, nil
}

// newAuthenticatedRequest creates a new http.Request with the mocked GlobalRepo
// and an optional authenticated user in the context.
func newAuthenticatedRequest(method, url string, body io.Reader, user *types.User, mockGlobalRepo *mock_repos.MockGlobalRepo) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
	if user != nil {
		ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, user)
		return req.WithContext(ctxWithAuth), nil
	}
	return req.WithContext(ctxWithRepo), nil
}
