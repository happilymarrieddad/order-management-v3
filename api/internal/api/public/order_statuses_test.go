package public_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Order Statuses Endpoint", func() {
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

	Context("GET /order-statuses", func() {
		It("should return a list of all order statuses", func() {
			req := createRequest("GET", "/order-statuses")
			rr := executeRequest(req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))

			var orderStatuses []string
			err := json.Unmarshal(rr.Body.Bytes(), &orderStatuses)
			Expect(err).NotTo(HaveOccurred())

			// Assert that the returned order statuses contain expected values
			Expect(orderStatuses).To(ConsistOf(
				"pending_acceptance",
				"pending_booking",
				"hold",
				"booked",
				"shipped_in_transit",
				"delivered",
				"ready_to_invoice",
				"invoiced",
				"rejected",
				"cancelled",
				"hold_for_pod",
				"order_template",
				"paid_in_full",
			))
		})
	})

	Context("Unsupported Methods for /order-statuses", func() {
		It("should return 405 Method Not Allowed for POST", func() {
			req := createRequest("POST", "/order-statuses")
			rr := executeRequest(req)
			Expect(rr.Code).To(Equal(http.StatusMethodNotAllowed))
		})

		It("should return 405 Method Not Allowed for PUT", func() {
			req := createRequest("PUT", "/order-statuses")
			rr := executeRequest(req)
			Expect(rr.Code).To(Equal(http.StatusMethodNotAllowed))
		})

		It("should return 405 Method Not Allowed for DELETE", func() {
			req := createRequest("DELETE", "/order-statuses")
			rr := executeRequest(req)
			Expect(rr.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})
})
