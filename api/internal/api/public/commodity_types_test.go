package public_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commodity Types Endpoint", func() {
	// createRequest creates a new HTTP request for testing.
	createRequest := func(method, path string) *http.Request {
		req, err := http.NewRequest(method, path, nil)
		Expect(err).NotTo(HaveOccurred())
		return req
	}

	// executeRequest executes the request against the router and returns the recorder.
	executeRequest := func(req *http.Request) *httptest.ResponseRecorder {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		return rr
	}

	Context("GET /commodity-types", func() {
		It("should return a list of all commodity types", func() {
			req := createRequest("GET", "/commodity-types")
			rr := executeRequest(req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))

			var commodityTypes []string
			err := json.Unmarshal(rr.Body.Bytes(), &commodityTypes)
			Expect(err).NotTo(HaveOccurred())

			// Assert that the returned commodity types contain expected values
			Expect(commodityTypes).To(ConsistOf("produce"))
		})
	})

	Context("Unsupported Methods for /commodity-types", func() {
		It("should return 405 Method Not Allowed for POST", func() {
			req := createRequest("POST", "/commodity-types")
			rr := executeRequest(req)
			Expect(rr.Code).To(Equal(http.StatusMethodNotAllowed))
		})

		It("should return 405 Method Not Allowed for PUT", func() {
			req := createRequest("PUT", "/commodity-types")
			rr := executeRequest(req)
			Expect(rr.Code).To(Equal(http.StatusMethodNotAllowed))
		})

		It("should return 405 Method Not Allowed for DELETE", func() {
			req := createRequest("DELETE", "/commodity-types")
			rr := executeRequest(req)
			Expect(rr.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})
})
